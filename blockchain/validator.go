package blockchain

import "time"

// ValidatorProfile from Driver.TransactionsQueue, we can calculate proof of physical work
// validator staking can happen after completing x rides
// validator staking can happen after paying for y rides as a passenger.
// pass monthly kyc
// earn enough trust rating.
type ValidatorProfile struct {
	UUID        string
	RidesServed int
	AvgRating   float64
	Verified    bool // manual/off-chain KYC
	LastKYC     time.Time
}

type ValidatorLogEvt struct {
	ValidatorUUID string
	Evt           string // i.e genesis validator verifying first driver
	Timestamp     time.Time
}

type DriverVerificationRequest struct {
	DriverUUID  string
	RequestedBy string
	Timestamp   time.Time
	Status      string
}

func (rc *RideChain) logValidatorEvent(validatorUUID, evt string) {
	rc.ValidatorLog = append(rc.ValidatorLog, ValidatorLogEvt{
		ValidatorUUID: validatorUUID,
		Evt:           evt,
		Timestamp:     time.Now(),
	})
}
