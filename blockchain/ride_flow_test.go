package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRideFlow_Happy_Path(t *testing.T) {
	driver := "genesis-123"
	rider := "rider-abc"
	rc, err := NewRideChain("test/token_ledger.json")
	assert.Nil(t, err)

	err = rc.BecomeValidator(driver)
	assert.Nil(t, err)

	// adds RideTx to pendingRideTxs
	rideTxEvts := []RideTxEvt{}
	rideTxEvts = append(rideTxEvts, RideTxEvt{
		EventType: RideRequested,
	}, RideTxEvt{
		EventType: DriverAccepted,
	}, RideTxEvt{
		EventType: RiderPaymentRecieved,
	})
	tx, err := rc.SubmitPendingRideTx(RideTx{
		RiderUUID:       rider,
		DriverUUID:      driver,
		PaidAmount:      100,
		PickupCode:      "1931",
		StripeSessionId: "someStripeSuccessString",
		ComputedRoute: ComputedRoute{
			Destination: "some destination hopefully not final😅",
		},
		RideTxEvts: rideTxEvts,
		PickupLocation: LatLng{
			Lat: "36.00000",
			Lng: "-86.00000",
		},
	})
	assert.Nil(t, err)

	err = rc.SubmitPickupProof(tx, "1931")
	assert.Nil(t, err)

	err = rc.SubmitDropoff(tx, LatLng{Lat: "36.1684", Lng: "86.8259"})
	assert.Nil(t, err)

	// tx.TxID = generateRideHash(tx)
	// happens here now
	txID, err := rc.ApproveRideTx(tx, driver)
	assert.Nil(t, err)

	_, ok := RideLedger[txID]
	assert.True(t, ok)

	// todo assert RideTxEvts

}
