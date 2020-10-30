package data

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Snapshot data type
type Snapshot [32]byte

// State contains current state of the app
type State struct {
	Balances map[string]uint
	txPool   []Tx
	dbFile   *os.File
	snapshot [32]byte
}

// NewStateFromDisk reconstruct the blockchain from disk dB
func NewStateFromDisk() (*State, error) {

	genesisFilePath := "/home/shivansh_tiwari/go/src/blockchain-go/data/genesis.json"
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

	txDbFilePath := "/home/shivansh_tiwari/go/src/blockchain-go/data/tx.db"
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Tx, 0), f, Snapshot{}}

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

// Persist sync the transactions with on disk dB
func (s *State) Persist() (Snapshot, error) {
	// make a copy of the loop
	mem := make([]Tx, len(s.txPool))
	copy(mem, s.txPool)

	for i := 0; i < len(mem); i++ {
		txJSON, err := json.Marshal(s.txPool[i])
		if err != nil {
			return Snapshot{}, err
		}

		fmt.Println("Persisting the new Transaction into the disk")
		if _, err = s.dbFile.Write(append(txJSON, '\n')); err != nil {
			return Snapshot{}, err
		}

		// now remove the synced transaction
		err = s.Snapshot()
		if err != nil {
			return Snapshot{}, err
		}
		s.txPool = append(s.txPool[:i], s.txPool[i+1:]...)
	}
	return s.snapshot, nil

}

// apply a transaction on the state
func (s *State) apply(tx Tx) error {

	// skipping the validation during reward
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

// Close close the dB connection
func (s *State) Close() error {
	return s.dbFile.Close()
}

// Snapshot returns the hash value of complete dB file
func (s *State) Snapshot() error {
	_, err := s.dbFile.Seek(0, 0)
	if err != nil {
		return err
	}
	txData, err := ioutil.ReadAll(s.dbFile)
	if err != nil {
		return err
	}
	s.snapshot = sha256.Sum256(txData)
	return nil
}
