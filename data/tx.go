package data

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
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

	// ID of the reference Tx
	ID []byte
	// Index of the reference Tx
	Out int
	// Account string, connecting the Input and Output
	Sig string
}

// SetID for the transaction
func (tx *Tx) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	if err != nil {
		log.Panic("Error in setting Id")
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash
}

// CoinbaseTx which will be used in Genesis and Rewards
func CoinbaseTx(to, data string, amount int) Tx {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{[]byte{}, -1, to}
	txout := TxOutput{amount, to}

	tx := Tx{Snapshot{}, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return tx
}

// NewTransaction .
func NewTransaction(from, to string, amount int, state *State) Tx {
	var inputs []TxInput
	var outputs []TxOutput
	availabe, unspentOutputs := state.GetSpendableOutputs(from)

	if availabe < amount {
		log.Panic("Insufficient")
	}

	for txid, outs := range unspentOutputs {
		txID, err := hex.DecodeString(txid)

		if err != nil {
			return Tx{}
		}

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	if availabe > amount {
		outputs = append(outputs, TxOutput{availabe - amount, from})
	}

	tx := Tx{Snapshot{}, inputs, outputs}
	tx.SetID()

	return tx
}

// IsCoinbase .
func (tx *Tx) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// CanUnlock .
func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

// CanBeUnlocked .
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
