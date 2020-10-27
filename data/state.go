package data

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// State contains current state of the app
type State struct {
	Balances map[string]uint
	txPool   []Tx
	dbFile   *os.File
}

// NewStateFromDisk reconstruct the blockchain from disk dB
func NewStateFromDisk() (*State, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	genesisFilePath := filepath.Join(cwd, "data", "genesis.json")
	gen, err := loadGenesis(genesisFilePath)
	if err != nil {
		return nil, err
	}
	balances := make(map[string]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	// till here, the genesis contract is processed.
	// now, all the transactions will be replayed

	txDbFilePath := filepath.Join(cwd, "data", "tx.db")
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Tx, 0), f}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		var tx Tx
		json.Unmarshal(scanner.Bytes(), &tx)

		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}

	return state, nil
}

// Add a transaction to State
func (s *State) Add(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}
	s.txPool = append(s.txPool, tx)
	return nil
}

// apply a transaction on the state
func (s *State) apply(tx Tx) error {

	if tx.isReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}
	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}
