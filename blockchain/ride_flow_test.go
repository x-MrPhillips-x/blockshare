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
		RideTxEvts:      rideTxEvts,
		PickupLocation:  "South Nashville",
		DropOffLocation: "Downtown Area",
	})
	assert.Nil(t, err)

	err = rc.SubmitPickupProof(tx)
	assert.Nil(t, err)

	err = rc.SubmitDropoff(tx, "Downtown Area")
	assert.Nil(t, err)

	// tx.TxID = generateRideHash(tx)
	// happens here now
	txID, err := rc.ApproveRideTx(tx, driver)
	assert.Nil(t, err)

	_, ok := RideLedger[txID]
	assert.True(t, ok)

	// todo assert RideTxEvts

}
