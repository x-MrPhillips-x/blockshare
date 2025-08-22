package blockchain

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RideChainer represents the ride related chain behavior
type RideChainer interface {
	SubmitPendingRideTx(tx RideTx) (RideTx, error)
	BecomeValidator(driverUUID string) error
	StakeTokens(amount int, driverUUID string) error
	SlashValidator(driverUUID string, slasher string, reason string) error
	VerifyDriver(driverUUID string, validator string, results string) error
	GetDriverStake(driverUUID string) int
	IsValidator(driverUUID string) bool
	RewardValidator(validatorUUID string, amount int) error
	ApproveRideTx(tx RideTx, validatorUUID string) (string, error)
	RequestDriverVerification(driverUUID, requestedBy string) error
	SubmitPickupProof(tx RideTx, pickupCode string) error
	SubmitDropoff(tx RideTx, dropoffLocation LatLng) error
	HasActiveRide(driverUUID string) bool
}

// RideChain represents the entire blockchain composed of rideTx
type RideChain struct {
	DriverStakes map[string]int // driverUUID → amount
	TokenLedger  *TokenLedger
	Validators   map[string]bool // driverUUID -> isValidator
	// PendingRideTxs map of riderUUID -> RideTx
	PendingRideTxs       map[string]RideTx
	RideApprovals        map[string]map[string]bool // txID → validatorUUID → approval
	ApprovalQuorum       int                        // min approvals required
	PendingVerifications map[string]DriverVerificationRequest
	minValidatorStake    int
}

func NewRideChain(ledgeFileLocation string) (*RideChain, error) {
	ledger, err := LoadTokenLedgerFromFile(ledgeFileLocation)
	if err != nil {
		return nil, err
	}
	return &RideChain{
		TokenLedger:          ledger,
		DriverStakes:         make(map[string]int),
		Validators:           make(map[string]bool),
		PendingRideTxs:       make(map[string]RideTx),
		RideApprovals:        make(map[string]map[string]bool),
		ApprovalQuorum:       1, // for now there is only genesis validator
		PendingVerifications: make(map[string]DriverVerificationRequest),
		minValidatorStake:    10,
	}, nil
}

// SubmitPendingRideTx adds a active RideTx to the pendingRideTx queue
// once the rideTx is complete this RideTx will move to AwaitingApproval
func (rc *RideChain) SubmitPendingRideTx(tx RideTx) (RideTx, error) {
	if err := ValidateRideTx(tx); err != nil {
		return RideTx{}, err
	}

	rc.PendingRideTxs[tx.DriverUUID] = tx

	fmt.Printf("Ride submitted: %v\n", tx)
	return tx, nil
}

// ValidateRideTx
// also adds the RiderPaymentReceived
func ValidateRideTx(tx RideTx) error {
	// 1. Required field checks
	// todo tx.TxID is not given until sumission to mempool,
	// should we already have this by now?
	if tx.DriverUUID == "" || tx.RiderUUID == "" {
		return errors.New("missing core identifiers")
	}
	if tx.PaidAmount <= 0 {
		return errors.New("invalid or missing PaidAmount")
	}

	// 2. Location validity
	if tx.PickupLocation == "" {
		return errors.New("invalid pickup location")
	}

	if tx.DropOffLocation == "" {
		return errors.New("invalid drop off location")
	}

	// 3. Event history lifecycle (simple sanity check)
	if len(tx.RideTxEvts) == 0 {
		return errors.New("no ride events recorded")
	}

	// Check logical flow
	// todo add dropoffEvt
	// todo paidEvt
	var hasRideRequestedEvt, hasDriverAcceptedEvt, hasRiderPaymentRecievedEvt bool
	for _, evt := range tx.RideTxEvts {
		switch evt.EventType {
		case RideRequested:
			hasRideRequestedEvt = true
		case DriverAccepted:
			hasDriverAcceptedEvt = true
		case RiderPaymentRecieved:
			hasRiderPaymentRecievedEvt = true
		}
	}

	if !hasRideRequestedEvt || !hasDriverAcceptedEvt || !hasRiderPaymentRecievedEvt {
		return errors.New("ride transaction event flow incomplete")
	}

	// timestamp sanity because these operations should take a few minutes
	if tx.TimeRequested.After(time.Now().Add(3 * time.Minute)) {
		return errors.New("invalid or future timestamp")
	}

	return nil
}

// TODO guard with mutex
func (rc *RideChain) BecomeValidator(driverUUID string) error {
	stake := rc.TokenLedger.GetStake(driverUUID)

	// Genesis validator rule: allow bootstrapper with any stake
	// TODO should the genesis validator at some point need to
	// also stake the minValidatorStake?
	if len(rc.Validators) == 0 {
		rc.Validators[driverUUID] = true
		// todo update RideTxEvts
		// rc.logValidatorEvent(driverUUID, "let there be light! genesis validator created")
		return nil
	}

	if stake < rc.minValidatorStake {
		return fmt.Errorf("driver %s must stake at least %d tokens to become a validator", driverUUID, rc.minValidatorStake)
	}

	rc.Validators[driverUUID] = true
	// todo update RideTxEvts
	// rc.logValidatorEvent(driverUUID, "became a validator")

	return nil
}

func (rc *RideChain) StakeTokens(amount int, driverUUID string) error {
	if err := rc.TokenLedger.Stake(driverUUID, amount); err != nil {
		return err
	}
	return rc.TokenLedger.SaveToFile()
}

// SlashValidator punishes bad validators by slashing their stake
func (rc *RideChain) SlashValidator(driverUUID string, slasher string, reason string) error {
	rc.TokenLedger.mu.Lock()
	defer rc.TokenLedger.mu.Unlock()

	// Only validators can slash
	if !rc.IsValidator(slasher) {
		return fmt.Errorf("unauthorized: %s is not a validator", slasher)
	}

	if !rc.IsValidator(driverUUID) {
		return fmt.Errorf("%s is not a validator", driverUUID)
	}

	// prevent getStake deadlock
	stake := rc.TokenLedger.Stakes[driverUUID]
	if stake <= 0 {
		return fmt.Errorf("validator %s has no stake to slash", driverUUID)
	}

	slashedAmount := stake / 2
	rc.TokenLedger.Stakes[driverUUID] -= slashedAmount
	// todo update RideTxEvts

	// rc.logValidatorEvent(driverUUID, fmt.Sprintf("was slashed %d tokens", slashedAmount))

	// Remove validator status
	delete(rc.Validators, driverUUID)

	// todo update RideTxEvts

	// rc.logValidatorEvent(driverUUID, "removed from validators")

	fmt.Printf("Validator %s was slashed by validator %s for %d tokens. Reason: %s\n",
		driverUUID, slasher, slashedAmount, reason)

	// prevent SaveToFile deadlock
	data, err := json.MarshalIndent(rc.TokenLedger, "", "  ")
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(rc.TokenLedger.filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(rc.TokenLedger.filename, data, 0644)
}

// VerifyDriver for now is a simple validation action
// we will want to get insurance and background check info from ride module
// here is an example of the links in TN
// - fetch from https://sor.tbi.tn.gov/api/search : results from here will probably need to be manually (by human) parsed
// - fetch https://verifyinsurance.revenue.tn.gov/assets/api/api.php : results from here can be automatically parsed
func (rc *RideChain) VerifyDriver(driverUUID string, validator string, results string) error {
	if !rc.IsValidator(validator) {
		return fmt.Errorf("%s is not a validator", validator)
	}

	// check pending verificaiton request
	request, exists := rc.PendingVerifications[driverUUID]
	if !exists {
		return fmt.Errorf("no pending verification for driver %s", driverUUID)
	}

	request.Status = "approved"
	rc.PendingVerifications[driverUUID] = request

	fmt.Printf("Driver %s verified by validator %s with results %s\n", driverUUID, validator, results)

	// todo update RideTxEvts

	// rc.logValidatorEvent(validator, fmt.Sprintf("Driver %s verified by validator %s with results %s\n", driverUUID, validator, results))
	return nil
}

func (rc *RideChain) GetDriverStake(driverUUID string) int {
	return rc.DriverStakes[driverUUID]
}

func (rc *RideChain) GetPendingRideTx() {
}

func (rc *RideChain) IsValidator(driverUUID string) bool {
	return rc.Validators[driverUUID]
}

func (rc *RideChain) RewardValidator(validatorUUID string, amount int) error {
	rc.TokenLedger.mu.Lock()
	defer rc.TokenLedger.mu.Unlock()

	if !rc.IsValidator(validatorUUID) {
		return fmt.Errorf("%s is not a validator", validatorUUID)
	}

	rc.TokenLedger.Balances[validatorUUID] += amount
	fmt.Printf("Validator %s rewarded %d tokens\n", validatorUUID, amount)
	return rc.TokenLedger.SaveToFile()
}

// ApproveRideTx approve and complete the RideTx after this
// the driver will be able to make trx again
func (rc *RideChain) ApproveRideTx(tx RideTx, validatorUUID string) (string, error) {
	if !rc.IsValidator(validatorUUID) {
		return "", fmt.Errorf("%s is not a validator", validatorUUID)
	}
	tx, exists := rc.PendingRideTxs[tx.DriverUUID]
	if !exists {
		return "", fmt.Errorf("ride %v not found", tx)
	}
	if rc.RideApprovals[tx.DriverUUID][validatorUUID] {
		return "", fmt.Errorf("validator %v already approved ride %v", validatorUUID, tx)
	}

	// Register approval
	rc.RideApprovals[tx.DriverUUID] = make(map[string]bool)
	rc.RideApprovals[tx.DriverUUID][validatorUUID] = true

	tx.RideTxEvts = append(tx.RideTxEvts, RideTxEvt{
		EventType: RideApproved,
		Timestamp: time.Now(),
	})

	// Count approvals
	if len(rc.RideApprovals[tx.DriverUUID]) >= rc.ApprovalQuorum {
		tx.TxID = generateRideHash(tx)

		// Move to ledger
		RideLedger[tx.TxID] = tx
		delete(rc.PendingRideTxs, tx.DriverUUID)
		delete(rc.RideApprovals, tx.DriverUUID)
		fmt.Printf("Ride %v approved and committed\n", tx)

		// rc.logValidatorEvent(validatorUUID, fmt.Sprintf("approved txID and commited %s", txID))

	}

	fmt.Printf("RideTx approved: %v\n", tx.TxID)

	return tx.TxID, nil
}

func (rc *RideChain) RequestDriverVerification(driverUUID, requestedBy string) error {
	if _, ok := rc.PendingVerifications[driverUUID]; ok {
		return fmt.Errorf("verification for driver %s already requested", driverUUID)
	}
	rc.PendingVerifications[driverUUID] = DriverVerificationRequest{
		DriverUUID:  driverUUID,
		RequestedBy: requestedBy,
		Timestamp:   time.Now(),
		Status:      "pending",
	}
	return nil
}

func (rc *RideChain) SubmitPickupProof(tx RideTx) error {
	tx, exists := rc.PendingRideTxs[tx.DriverUUID]
	if !exists {
		return fmt.Errorf("rideTx %v not found", tx)
	}

	if tx.PickupConfirmed {
		return fmt.Errorf("pickup already confirmed for rideTx %v", tx)
	}

	// Confirm pickup
	tx.PickupConfirmed = true
	rc.PendingRideTxs[tx.RiderUUID] = tx // save updated tx

	// todo update RideTxEvts

	// rc.logValidatorEvent(tx.DriverUUID, fmt.Sprintf("pickup code confirmed for ride %s", txID))

	fmt.Printf("Pickup code confirmed for rideTx %v\n", tx)
	return nil
}

func (rc *RideChain) SubmitDropoff(tx RideTx, dropoffLocation string) error {
	tx, exists := rc.PendingRideTxs[tx.RiderUUID]
	if !exists {
		return fmt.Errorf("rideTx %v not found", tx)
	}

	if !tx.PickupConfirmed {
		return fmt.Errorf("pickup not confirmed for ride %v", tx)
	}

	if tx.DropoffConfirmed {
		return fmt.Errorf("dropoff already submitted for ride %v", tx)
	}

	tx.DropOffLocation = dropoffLocation
	tx.DropoffConfirmed = true
	tx.DropoffTime = time.Now()

	rc.PendingRideTxs[tx.RiderUUID] = tx // update with drop-off

	// todo update RideTxEvts

	// rc.logValidatorEvent(tx.DriverUUID, fmt.Sprintf("dropoff submitted for ride %s", txID))

	fmt.Printf("Dropoff submitted for rideTx %v\n", tx)
	return nil
}

func (rc *RideChain) HasActiveRide(driverUUID string) bool {
	for _, tx := range rc.PendingRideTxs {
		if tx.DriverUUID == driverUUID && !tx.DropoffConfirmed {
			return true
		}
	}
	return false
}
