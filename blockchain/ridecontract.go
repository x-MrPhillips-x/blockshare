package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

// RideContract represents the ride related chain behavior
// type RideContract interface {
// 	SubmitRide(tx RideTx) error
// 	BecomeValidator(driverUUID string) error
// 	StakeTokens(amount int, driverUUID string) error
// 	SlashValidator(driverUUID string, reason string) error
// 	VerifyDriver(driverUUID string, result bool) error
// }

// Core ride transaction model
type RideTx struct {
	TxID             string    `json:"txID"`
	DriverUUID       string    `json:"driverUUID"`
	RiderUUID        string    `json:"riderUUID"`
	TimeRequested    time.Time `json:"timeRequested"`
	EstimatedPickup  time.Time `json:"estimatedPickup"`
	EstimatedDropoff time.Time `json:"estimatedDropoff"`
	PickupLocation   LatLng    `json:"pickupLocation"`
	DropoffLocation  LatLng    `json:"dropoffLocation"`
	Passengers       int       `json:"passengers"`
	Luggage          int       `json:"luggage"`
	PaidAmount       int       `json:"paidAmount"`
	PickupCode       string    `json:"pickupCode"`
	RouteHash        string    `json:"routeHash"` // optional full route hash
	Timestamp        time.Time `json:"timestamp"` // when submitted on chain
	PickupConfirmed  bool      `json:"pickUpConfirmed"`
	DropoffConfirmed bool      `json:"dropOffConfirmed"`
	DropoffTime      time.Time `json:"dropOffTime"`
}

// Simulated ledger (in-memory for now)
var RideLedger = map[string]RideTx{}

// Generate a SHA-256 hash of the ride data (for TxID or chain anchoring)
func generateRideHash(tx RideTx) string {
	data, _ := json.Marshal(tx)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash[:])
}
