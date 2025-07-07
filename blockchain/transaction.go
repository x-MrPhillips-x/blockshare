package blockchain

import "time"

type Transaction struct {
	PickUp  LatLng `json:"pickUp"`
	DropOff string `json:"dropOff"`
	Price   int    `json:"price"`
}

type LatLng struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type RideTxEventType string

const (
	RideRequested        RideTxEventType = "RideRequested"
	DriverAccepted       RideTxEventType = "DriverAccepted"
	RideApproved         RideTxEventType = "RideApproved"
	PickupVerified       RideTxEventType = "PickupVerified"
	DropoffConfirmed     RideTxEventType = "DropoffConfirmed"
	InsuranceVerified    RideTxEventType = "InsuranceVerified"
	DriverValidated      RideTxEventType = "DriverValidated"
	RiderPaymentRecieved RideTxEventType = "RiderPaymentRecieved"

	// ValidationEvt represents validator finalizing this RideTx
	ValidationEvt RideTxEventType = "ValidationEvt"
)

type RideTxEvt struct {
	EventType RideTxEventType        `json:"eventType"`
	Timestamp time.Time              `json:"timestamp"`
	Validator string                 `json:"validator"`
	Metadata  map[string]interface{} `json:"metadata"` // or a typed struct
}
