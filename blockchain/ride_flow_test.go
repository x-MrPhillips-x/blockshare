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

	txID, err := rc.SubmitRideTx(RideTx{
		RiderUUID:  rider,
		DriverUUID: driver,
		PaidAmount: 100,
		PickupCode: "1931",
	})
	assert.Nil(t, err)

	err = rc.SubmitPickupProof(txID, "1931")
	assert.Nil(t, err)

	err = rc.SubmitDropoff(txID, LatLng{Lat: "36.1684", Lng: "86.8259"})
	assert.Nil(t, err)

	err = rc.ApproveRideTx(txID, driver)
	assert.Nil(t, err)

	_, ok := RideLedger[txID]
	assert.True(t, ok)

}
