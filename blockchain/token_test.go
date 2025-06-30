package blockchain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenLedger_Stake(t *testing.T) {
	tests := []struct {
		name             string
		ledger           *TokenLedger
		driverUUID       string
		amount           int
		wantErr          error
		fundDriver       bool
		remainingBalance int
	}{
		{
			name:       "driver is broke broke trying to stake 10 tokens",
			ledger:     NewTokenLedger(),
			driverUUID: "driver-123",
			amount:     10,
			wantErr:    fmt.Errorf("insufficient balance for driver driver-123"),
		},
		{
			name:             "driver stakes tokens successfully",
			ledger:           NewTokenLedger(),
			driverUUID:       "driver-123",
			amount:           10,
			fundDriver:       true,
			remainingBalance: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// todo fund driver
			if tt.fundDriver {
				tt.ledger.Balances[tt.driverUUID] = 11
			}

			err := tt.ledger.Stake(tt.driverUUID, tt.amount)
			if err != nil {
				assert.NotNil(t, tt.wantErr)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.Nil(t, tt.wantErr)
				assert.Equal(t, tt.remainingBalance, tt.ledger.Balances[tt.driverUUID])
			}

		})
	}
}
