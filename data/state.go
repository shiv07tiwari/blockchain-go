package data

import (
	"bufio"
	"encoding/hex"
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
	txPool   []Tx
	dbFile   *os.File
	Snapshot [32]byte
	Users    map[string]bool
}

// NewStateFromDisk reconstruct the blockchain from disk dB
func NewStateFromDisk() (*State, error) {

	isFileEmpty := true

	txDbFilePath := "/home/shivansh_tiwari/go/src/blockchain-go/data/tx.db"
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	users := make(map[string]bool)
	state := &State{make([]Tx, 0), f, Snapshot{}, users}

	for scanner.Scan() {
		isFileEmpty = false
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		var blockData BlockData
		json.Unmarshal(scanner.Bytes(), &blockData)

		for _, tx := range blockData.Value.TXs {
			for _, in := range tx.Inputs {
				if _, ok := state.Users[in.Sig]; !ok {
					users[in.Sig] = true
				}
			}
			for _, out := range tx.Outputs {
				if _, ok := state.Users[out.PubKey]; !ok {
					users[out.PubKey] = true
				}
			}
		}

		state.Snapshot = blockData.Key
	}
	// Persist the genesis Block if not already done
	if isFileEmpty {
		tx, err := GenerateGenesis()
		if err != nil {
			return nil, err
		}
		err = state.Add(tx)
		state.Users[tx.Inputs[0].Sig] = true
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
	s.txPool = append(s.txPool, tx)
	return nil
}

// Persist sync the transactions with on disk dB
func (s *State) Persist() (Snapshot, error) {
	block := NewBlock(
		s.Snapshot,
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
	s.Snapshot = hash
	s.txPool = []Tx{}

	return s.Snapshot, nil

}

// Close close the dB connection
func (s *State) Close() error {
	return s.dbFile.Close()
}

// GetSpendableOutputs .
func (s *State) GetSpendableOutputs(address string) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := s.FindUnspentTransactions(address)
	collected := 0

	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID[:])

		for outID, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				collected += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outID)
			}
		}
	}
	return collected, unspentOuts
}

// FindUnspentTransactions .
func (s *State) FindUnspentTransactions(address string) []Tx {
	var unspentTxs []Tx
	var blocks []BlockData

	spentTXOs := make(map[string][]int)

	scanner := bufio.NewScanner(s.dbFile)
	s.dbFile.Seek(0, 0)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil
		}
		var blockData BlockData
		json.Unmarshal(scanner.Bytes(), &blockData)
		blocks = append(blocks, blockData)
	}

	for _, blockData := range blocks {
		for _, tx := range blockData.Value.TXs {
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}
	}

	for _, blockData := range blocks {
		for _, tx := range blockData.Value.TXs {
			txID := hex.EncodeToString(tx.ID[:])
		Outputs:
			for outIDX, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIDX {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, tx)
				}
			}
		}
	}
	return unspentTxs
}
