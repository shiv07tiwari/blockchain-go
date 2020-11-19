package data

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
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

	balances := make(map[string]uint)
	isFileEmpty := true

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
		var blockData BlockData
		json.Unmarshal(scanner.Bytes(), &blockData)

		for _, tx := range blockData.Value.TXs {
			isFileEmpty = false
			if err := state.apply(tx); err != nil {
				return nil, err
			}
		}
		state.snapshot = blockData.Key

	}

	// Persist the genesis Block if not already done
	if isFileEmpty {
		tx, err := GenerateGenesis()
		if err != nil {
			return nil, err
		}
		err = state.Add(tx)
		if err != nil {
			return nil, err
		}
		_, err = state.Persist()
		if err != nil {
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
	block := NewBlock(
		s.snapshot,
		uint64(time.Now().Unix()),
		s.txPool,
	)
	pow := CreateNewProof(&block)
	nonce, hash := pow.Mine()
	block.Nonce = nonce

	if pow.Validate() == false {
		return Snapshot{}, errors.New("Invalid Block")
	}

	blockData := BlockData{hash, block}

	blockDataJSON, err := json.Marshal(blockData)

	fmt.Println("Persisting new Block to Disk")

	_, err = s.dbFile.Write(append(blockDataJSON, '\n'))
	if err != nil {
		return Snapshot{}, err
	}
	s.snapshot = hash
	s.txPool = []Tx{}

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
