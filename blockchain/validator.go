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

type DriverVerificationRequest struct {
	DriverUUID  string
	RequestedBy string
	Timestamp   time.Time
	Status      string
}
