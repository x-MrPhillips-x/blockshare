package blockchain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const ledgerFilePath = "data/token_ledger.json"

type TokenLedger struct {
	Balances map[string]int `json:"balances"` // driverUUID -> token balance
	Stakes   map[string]int `json:"stakes"`   // driverUUID -> staked tokens
	mu       sync.RWMutex   `json:"-"`
	filename string         `json:"-"`
}

func NewTokenLedger() *TokenLedger {
	return &TokenLedger{
		Balances: make(map[string]int),
		Stakes:   make(map[string]int),
	}
}

// Mint tokens to a driver (e.g. admin or faucet action)
// TODO define faucet action, what are your thoughts here
func (m *TokenLedger) Mint(driverUUID string, amount int) {
	m.Balances[driverUUID] += amount
}

// Stake tokens
func (m *TokenLedger) Stake(driverUUID string, amount int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	balance, ok := m.Balances[driverUUID]
	if !ok || balance < amount {
		return fmt.Errorf("insufficient balance for driver %s", driverUUID)
	}
	m.Balances[driverUUID] -= amount
	m.Stakes[driverUUID] += amount
	return nil
}

func (m *TokenLedger) GetStake(driverUUID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Stakes[driverUUID]
}

// Unstake tokens
func (m *TokenLedger) Unstake(driverUUID string, amount int) error {
	if m.Stakes[driverUUID] < amount {
		return fmt.Errorf("not enough tokens staked")
	}
	m.Stakes[driverUUID] -= amount
	m.Balances[driverUUID] += amount
	return nil
}

func (t *TokenLedger) SaveToFile() error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(t.filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(t.filename, data, 0644)
}

func LoadTokenLedgerFromFile(filename string) (*TokenLedger, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return &TokenLedger{
			Stakes:   make(map[string]int),
			Balances: make(map[string]int),
			filename: filename,
		}, nil
	}

	var ledger TokenLedger
	err = json.Unmarshal(data, &ledger)
	if err != nil {
		return nil, err
	}

	// Init missing fields
	if ledger.Stakes == nil {
		ledger.Stakes = make(map[string]int)
	}
	if ledger.Balances == nil {
		ledger.Balances = make(map[string]int)
	}
	ledger.filename = filename
	return &ledger, nil
}
