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
	PickupLocation   LatLng    `json:"pickupLocation"`

	// PickUpPlaceDetails represent the place details at the pickup location
	PickUpPlaceDetails PlaceDetails `json:"pickUpPlaceDetails"`

	// todo we need to get this from the computedRoutes.destination if possible
	// the destination needs to be converted to lat,lng
	DropoffLocation  LatLng    `json:"dropoffLocation"`
	Passengers       int       `json:"passengers"` // todo deprecate put inside vehicle
	Luggage          int       `json:"luggage"`    // todo deprecate put inside vehicle
	PaidAmount       int       `json:"paidAmount"`
	PickupCode       string    `json:"pickupCode"`
	RouteHash        string    `json:"routeHash"` // optional full route hash
	PickupConfirmed  bool      `json:"pickUpConfirmed"`
	DropoffConfirmed bool      `json:"dropOffConfirmed"`
	DropoffTime      time.Time `json:"dropOffTime"`
	DriverLocation   LatLng    `json:"driverLocation"` // This is needed for the google embedded map string
	DriverAccepted   bool      `json:"driverAccepted"`

	// StripeSessionId represents successful payment from rider
	// TODO see what opportunities we have with buttons to
	// perform refund and blah blah
	StripeSessionId string `json:"stripeSessionId"`

	// RideTxEvts capture the events of the RideTx (i.e request, accept, paid, arrived...)
	RideTxEvts []RideTxEvt `json:"rideTxEvts"`

	// Vehicle details
	Vehicle Vehicle `json:"vehicle"`

	// ComputedRoute should have all the information on pickup/dropoff times/miles away
	ComputedRoute ComputedRoute `json:"computedRoute"`
}

// PlaceDetails present details for example type of place is grocery_store,
// meaning driver may have passenger with groceries
type PlaceDetails struct {
	Types            []string `json:"types"` // https://developers.google.com/maps/documentation/places/web-service/place-types
	PlaceId          string   `json:"placeId"`
	DisplayName      string   `json:"displayName"` // name displayed on marquee of place
	FormattedAddress string   `json:"formattedAddress"`
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

type ComputedRoute struct {
	UUID             string  `json:"uuid"`
	Brand            string  `json:"brand"`      // todo deprecate part of vehicle
	Year             string  `json:"year"`       // todo deprecate part of vehicle
	Model            string  `json:"model"`      // todo deprecate part of vehicle
	Color            string  `json:"color"`      // todo deprecate part of vehicle
	Img              string  `json:"img"`        // todo deprecate part of vehicle
	Price            int     `json:"price"`      // cents
	TrustScore       float64 `json:"trustScore"` // todo deprecate should be able to show this with the driver
	Departure        string  `json:"departure"`
	EstimatedArrival string  `json:"arrival"`
	Destination      string  `json:"destination"`
	MilesAway        int     `json:"milesAway"`
	MinutesAway      int     `json:"minutesAway"`
	TravelTime       int     `json:"travelTime"`
	TravelMiles      int     `json:"travelMiles"`
}
