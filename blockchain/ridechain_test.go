package blockchain

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRideChain_BecomeValidator(t *testing.T) {
	rc, _ := NewRideChain("test/token_ledger.json")
	type args struct {
		driverUUID string
	}
	tests := []struct {
		name       string
		driverUUID string
		wantErr    error
	}{
		{
			name:       "success the genenisis validator test",
			driverUUID: "driver-123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rc.BecomeValidator(tt.driverUUID)
			if err != nil {
				assert.Equal(t, err.Error(), err.Error())
			} else {
				assert.Nil(t, tt.wantErr)
			}

			assert.True(t, rc.Validators[tt.driverUUID])
		})
	}
}

func TestRideChain_SlashValidator(t *testing.T) {
	rc, _ := NewRideChain("data/token_ledger.json")
	genesisUUID := "genesis-123"
	firstDriver := "driver-123"

	type args struct {
		driverUUID string
		slasher    string
		reason     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		prepare func()
	}{
		{
			name: "genesis validator fails slashing first driver, because slasher is not a validator",
			args: args{
				driverUUID: firstDriver,
				slasher:    genesisUUID,
				reason:     fmt.Sprintf("slashing %s for some reason", firstDriver),
			},
			wantErr: errors.New("unauthorized: genesis-123 is not a validator"),
		},
		{
			name: "genesis validator fails slashing first driver, because first driver is not a validator",
			args: args{
				driverUUID: firstDriver,
				slasher:    genesisUUID,
				reason:     fmt.Sprintf("slashing %s for some reason", firstDriver),
			},
			prepare: func() {
				rc.BecomeValidator(genesisUUID)
			},
			wantErr: errors.New("driver-123 is not a validator"),
		},
		{
			name: "success genesis validator fails slashing first driver",
			args: args{
				driverUUID: firstDriver,
				slasher:    genesisUUID,
				reason:     fmt.Sprintf("slashing %s for some reason", firstDriver),
			},
			prepare: func() {

				rc.BecomeValidator(genesisUUID)

				// fund first driver for testing
				rc.TokenLedger.Balances[firstDriver] = 20

				rc.StakeTokens(10, firstDriver)
				rc.BecomeValidator(firstDriver)
			},
		}, // TODO slashing genesis validator
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			err := rc.SlashValidator(tt.args.driverUUID, tt.args.slasher, tt.args.reason)
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.Nil(t, tt.wantErr)
			}
		})
	}
}
