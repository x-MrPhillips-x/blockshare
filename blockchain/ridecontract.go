package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

// Core ride transaction model
type RideTx struct {
	TxID             string    `json:"txID"`
	DriverUUID       string    `json:"driverUUID"`
	RiderUUID        string    `json:"riderUUID"`
	TimeRequested    time.Time `json:"timeRequested"` // time rideTx construction began
	EstimatedPickup  time.Time `json:"estimatedPickup"`
	EstimatedDropoff time.Time `json:"estimatedDropoff"`

	PickupLocation  string `json:"pickUp"`  // general pickup "south nashville"
	DropOffLocation string `json:"dropOff"` // general dropoff "downtown aread"

	// the destination needs to be converted to lat,lng
	Passengers       int       `json:"passengers"` // todo deprecate put inside vehicle
	Luggage          int       `json:"luggage"`    // todo deprecate put inside vehicle
	PaidAmount       int       `json:"paidAmount"`
	RouteHash        string    `json:"routeHash"` // optional full route hash
	PickupConfirmed  bool      `json:"pickUpConfirmed"`
	DropoffConfirmed bool      `json:"dropOffConfirmed"`
	DropoffTime      time.Time `json:"dropOffTime"`
	DriverAccepted   bool      `json:"driverAccepted"`

	// RideTxEvts capture the events of the RideTx (i.e request, accept, paid, arrived...)
	RideTxEvts []RideTxEvt `json:"rideTxEvts"`

	// Vehicle details
	Vehicle Vehicle `json:"vehicle"`
}

type Vehicle struct {
	Brand string `json:"brand"`
	Year  string `json:"year"`
	Model string `json:"model"`
	Color string `json:"color"`
	Plate string `json:"plate"`
	Img   string `json:"img"`
	Seats int    `json:"seats"`
}

// Simulated ledger (in-memory for now)
var RideLedger = map[string]RideTx{}

// Generate a SHA-256 hash of the ride data (for TxID or chain anchoring)
func generateRideHash(tx RideTx) string {
	data, _ := json.Marshal(tx)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash[:])
}
