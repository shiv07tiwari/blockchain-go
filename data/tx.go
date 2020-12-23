package data

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
)

// Tx trasaction data structure
// Tx consists of inputs, outputs and transaction Id
type Tx struct {
	// TransactionId
	ID Snapshot `json:"id"`
	// List of inputs
	Inputs []TxInput `json:"inputs"`
	// List of outputs
	Outputs []TxOutput `json:"outputs"`
}

// TxOutput data structure
type TxOutput struct {
	// Value is the amount in the transaction.
	Value int
	// Account public key which is used to unlock the output by a particular address
	PubKey string
}

// TxInput data structure
// Input refers to the output transaction of previously made transactions.
type TxInput struct {

	// ID of the Tx which contains the reference Output
	ID []byte
	// Index of the reference Output
	Out int
	// Account string, connecting the Input and Output
	Sig string
}

// SetID encodes the complete Tx into a byte array
func (tx *Tx) SetID() {

	// Create a variable sized byte buffer
	var encoded bytes.Buffer
	var hash [32]byte

	// Get a new encoder and transmit the data of Tx
	encode := gob.NewEncoder(&encoded)
	encode.Encode(tx)

	// Get a fixed size byte array encrypted checksum
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash
}

// CoinbaseTx which will be used in Genesis and Rewards
func CoinbaseTx(to, data string, amount int) Tx {

	// The coinbase input will refer to a null output
	txin := TxInput{[]byte{}, -1, to}
	txout := TxOutput{amount, to}

	tx := Tx{Snapshot{}, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return tx
}

// NewTransaction creates a new transaction
func NewTransaction(from, to string, amount int, state *State) (Tx, error) {
	var inputs []TxInput
	var outputs []TxOutput

	// Get the spendable outputs for the sender
	availabe, unspentOutputs := state.GetSpendableOutputs(from)

	// If availabe amount for sender is less than the Tx amount, return Insufficient Error
	if availabe < amount {
		return Tx{}, errors.New("Insufficient Balance")
	}

	// Iterate over the unspent Outputs
	for txid, outs := range unspentOutputs {
		txID, err := hex.DecodeString(txid)

		if err != nil {
			return Tx{}, err
		}

		// Create the input referring to the unspent outputs
		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Create the Tx output
	outputs = append(outputs, TxOutput{amount, to})

	// Create a new output to return the remaining money back to the sender.
	if availabe > amount {
		outputs = append(outputs, TxOutput{availabe - amount, from})
	}

	// Create Tx, set its Id and return
	tx := Tx{Snapshot{}, inputs, outputs}
	tx.SetID()

	return tx, nil
}

// IsCoinbase is a helper function to tell if the Tx is coinbase or not
func (tx *Tx) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// CanUnlock checks if an input can be used by an address
func (in *TxInput) CanUnlock(address string) bool {
	return in.Sig == address
}

// CanBeUnlocked checks if an output can be used by an address
func (out *TxOutput) CanBeUnlocked(address string) bool {
	return out.PubKey == address
}
