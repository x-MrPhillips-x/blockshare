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
	SubmitRideTx(tx RideTx) (string, error)
	BecomeValidator(driverUUID string) error
	StakeTokens(amount int, driverUUID string) error
	SlashValidator(driverUUID string, reason string) error
	VerifyDriver(driverUUID string, result bool) error
	GetDriverStake(driverUUID string) int
	IsValidator(driverUUID string) bool
	RewardValidator(validatorUUID string, amount int) error
	ApproveRideTx(txID, validatorUUID string) error
	RequestDriverVerification(driverUUID, requestedBy string) error
	SubmitPickupProof(txID string, pickupCode string) error
	SubmitDropoff(txID, dropoffLocation string) error
}

// RideChain represents the entire blockchain composed of rideTx
type RideChain struct {
	DriverStakes         map[string]int // driverUUID → amount
	TokenLedger          *TokenLedger
	Validators           map[string]bool // driverUUID -> isValidator
	ValidatorLog         []ValidatorLogEvt
	PendingRides         map[string]RideTx          // txID → RideTx
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
		PendingRides:         make(map[string]RideTx),
		RideApprovals:        make(map[string]map[string]bool),
		ApprovalQuorum:       1, // for now there is only genesis validator
		PendingVerifications: make(map[string]DriverVerificationRequest),
		minValidatorStake:    10,
	}, nil
}

// SubmitRideTx validates and stores the ride
func (rc *RideChain) SubmitRideTx(tx RideTx) (string, error) {
	if tx.DriverUUID == "" || tx.RiderUUID == "" {
		return "", errors.New("invalid ride: missing driver or rider ID")
	}
	if tx.PaidAmount <= 0 {
		return "", errors.New("invalid payment amount")
	}
	tx.Timestamp = time.Now()
	tx.TxID = generateRideHash(tx)

	rc.PendingRides[tx.TxID] = tx
	rc.RideApprovals[tx.TxID] = make(map[string]bool)

	fmt.Printf("Ride submitted: %s\n", tx.TxID)
	return tx.TxID, nil
}

// TODO guard with mutex
func (rc *RideChain) BecomeValidator(driverUUID string) error {
	stake := rc.TokenLedger.GetStake(driverUUID)

	// Genesis validator rule: allow bootstrapper with any stake
	// TODO should the genesis validator at some point need to
	// also stake the minValidatorStake?
	if len(rc.Validators) == 0 {
		rc.Validators[driverUUID] = true
		rc.logValidatorEvent(driverUUID, "let there be light! genesis validator created")
		return nil
	}

	if stake < rc.minValidatorStake {
		return fmt.Errorf("driver %s must stake at least %d tokens to become a validator", driverUUID, rc.minValidatorStake)
	}

	rc.Validators[driverUUID] = true
	rc.logValidatorEvent(driverUUID, "became a validator")

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
	rc.logValidatorEvent(driverUUID, fmt.Sprintf("was slashed %d tokens", slashedAmount))

	// Remove validator status
	delete(rc.Validators, driverUUID)
	rc.logValidatorEvent(driverUUID, "removed from validators")

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
	rc.logValidatorEvent(validator, fmt.Sprintf("Driver %s verified by validator %s with results %s\n", driverUUID, validator, results))
	return nil
}

func (rc *RideChain) GetDriverStake(driverUUID string) int {
	return rc.DriverStakes[driverUUID]
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

func (rc *RideChain) ApproveRideTx(txID, validatorUUID string) error {
	if !rc.IsValidator(validatorUUID) {
		return fmt.Errorf("%s is not a validator", validatorUUID)
	}
	tx, exists := rc.PendingRides[txID]
	if !exists {
		return fmt.Errorf("ride %s not found", txID)
	}
	if rc.RideApprovals[txID][validatorUUID] {
		return fmt.Errorf("validator %s already approved ride %s", validatorUUID, txID)
	}

	// Register approval
	rc.RideApprovals[txID][validatorUUID] = true
	rc.logValidatorEvent(validatorUUID, fmt.Sprintf("approved txID %s", txID))

	// Count approvals
	if len(rc.RideApprovals[txID]) >= rc.ApprovalQuorum {
		// Move to ledger
		RideLedger[txID] = tx
		delete(rc.PendingRides, txID)
		delete(rc.RideApprovals, txID)
		fmt.Printf("Ride %s approved and committed\n", txID)
		rc.logValidatorEvent(validatorUUID, fmt.Sprintf("approved txID and commited %s", txID))

	}

	return nil
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

func (rc *RideChain) SubmitPickupProof(txID string, pickupCode string) error {
	tx, exists := rc.PendingRides[txID]
	if !exists {
		return fmt.Errorf("ride %s not found", txID)
	}

	if tx.PickupConfirmed {
		return fmt.Errorf("pickup already confirmed for ride %s", txID)
	}

	if tx.PickupCode != pickupCode {
		return fmt.Errorf("invalid pickup code for ride %s", txID)
	}

	// Confirm pickup
	tx.PickupConfirmed = true
	rc.PendingRides[txID] = tx // save updated tx

	rc.logValidatorEvent(tx.DriverUUID, fmt.Sprintf("pickup code confirmed for ride %s", txID))

	fmt.Printf("Pickup code confirmed for ride %s\n", txID)
	return nil
}

func (rc *RideChain) SubmitDropoff(txID string, dropoffLocation LatLng) error {
	tx, exists := rc.PendingRides[txID]
	if !exists {
		return fmt.Errorf("ride %s not found", txID)
	}

	if !tx.PickupConfirmed {
		return fmt.Errorf("pickup not confirmed for ride %s", txID)
	}

	if tx.DropoffConfirmed {
		return fmt.Errorf("dropoff already submitted for ride %s", txID)
	}

	tx.DropoffLocation = dropoffLocation
	tx.DropoffConfirmed = true
	tx.DropoffTime = time.Now()

	rc.PendingRides[txID] = tx // update with drop-off

	rc.logValidatorEvent(tx.DriverUUID, fmt.Sprintf("dropoff submitted for ride %s", txID))

	fmt.Printf("Dropoff submitted for ride %s\n", txID)
	return nil
}
