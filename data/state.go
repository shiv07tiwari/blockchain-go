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
	// In Memory pool of transactions
	txPool []Tx
	// Pointer to the database file
	dbFile *os.File
	// Last Hash value
	Snapshot [32]byte
	// List of users
	Users map[string]bool
}

// NewStateFromDisk reconstruct the blockchain from disk dB
func NewStateFromDisk() (*State, error) {

	// A helper variable to check if CoinbaseTx should be created
	isFileEmpty := true

	// Path of the database file
	txDbFilePath := "/home/shivansh_tiwari/go/src/blockchain-go/data/tx.db"

	// Open the file with permissions to read and append to the file.
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)

	// Handle the error, if any while opening the file
	if err != nil {
		return nil, err
	}

	// Create the scanner to scan the dB file
	scanner := bufio.NewScanner(f)
	// Users map, to be used in state
	users := make(map[string]bool)
	// object of state with initial values
	state := &State{make([]Tx, 0), f, Snapshot{}, users}

	// Iterate line by line over the dB file
	for scanner.Scan() {
		// Set the helper bool as false, i.e. CoinbaseTx is already created
		isFileEmpty = false

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		// Parse the JSON encoded block into blockData
		var blockData BlockData
		json.Unmarshal(scanner.Bytes(), &blockData)

		// Generate the users list
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

		// Update the latest Snapshot of the state
		state.Snapshot = blockData.Key
	}

	// Persist the genesis Block if not already done
	if isFileEmpty {

		// Returns the CoinbaseTx and error, if any
		tx, err := GenerateGenesis()
		if err != nil {
			return nil, err
		}

		// Add the coinbase Tx to the state
		err = state.Add(tx)

		// Update the user list
		state.Users[tx.Inputs[0].Sig] = true
		if err != nil {
			return nil, err
		}

		// Persist the state's transaction Pool
		_, err = state.Persist()
		if err != nil {
			return nil, err
		}
	}
	return state, nil
}

// Add a transaction to State's transaction pool
func (s *State) Add(tx Tx) error {
	s.txPool = append(s.txPool, tx)
	return nil
}

// Persist syncs the transactions with on disk dB
func (s *State) Persist() (Snapshot, error) {

	// Create a new block
	block := NewBlock(
		s.Snapshot,
		// Get current time
		uint64(time.Now().Unix()),
		s.txPool,
	)

	// Create a new Proof Of work for the block and Mine it
	pow := CreateNewProof(&block)
	nonce, hash := pow.Mine()

	// Store the nonce
	block.Nonce = nonce

	// Check if the block is valid or not
	if pow.Validate() == false {
		return Snapshot{}, errors.New("Invalid Block")
	}

	blockData := BlockData{hash, block}

	blockDataJSON, err := json.Marshal(blockData)

	fmt.Println("Persisting new Block to Disk")

	// Write the newBlock data to the file and clear the txPool, update the snapshot
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

// GetSpendableOutputs returns the spendable amount for a user and the list of unspent Outputs
func (s *State) GetSpendableOutputs(address string) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)

	// Get all the unspent transactions for a user
	unspentTxs := s.FindUnspentTransactions(address)
	collected := 0

	// Iterate over the unspent transactions
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID[:])

		// Iterate over all the outputs in the transaction
		for outID, out := range tx.Outputs {

			// If the user can unlock the output, add the amount and append the Output to unspent Outputs
			if out.CanBeUnlocked(address) {
				collected += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outID)
			}
		}
	}
	return collected, unspentOuts
}

// FindUnspentTransactions returns all the unspent transactions for an address
func (s *State) FindUnspentTransactions(address string) []Tx {

	// Array of unspent transactions
	var unspentTxs []Tx
	// Array of blocks in the state
	var blocks []BlockData

	// Map of spent transactions
	spentTXOs := make(map[string][]int)

	// Create a scanner to read the file and reset the pointer
	scanner := bufio.NewScanner(s.dbFile)
	s.dbFile.Seek(0, 0)

	// Iterate over the dB and append the blocks
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil
		}
		var blockData BlockData
		json.Unmarshal(scanner.Bytes(), &blockData)
		blocks = append(blocks, blockData)
	}

	// Iterate over each block
	for _, blockData := range blocks {
		// Iterate over all the Txs in a block
		for _, tx := range blockData.Value.TXs {

			// If transaction is not coinbase, put its reference output and index in spentTxs map
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

	// Iterate over each block
	for _, blockData := range blocks {
		for _, tx := range blockData.Value.TXs {
			txID := hex.EncodeToString(tx.ID[:])
		Outputs:
			// Iterate over each Output of the Tx
			for outIDX, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIDX {
							// If the output's index is in spentTx, continue
							continue Outputs
						}
					}
				}
				// Else, add the Tx in unspent Tx if the output can be unlocked by the user
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, tx)
				}
			}
		}
	}

	// Return all the unspent transactions
	return unspentTxs
}
