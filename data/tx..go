package data

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

// Tx trasaction data structure
type Tx struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

// Reward is added for people who maintain and mine blockchain
func (t *Tx) isReward() bool {
	return t.Data == "reward"
}

// NewTx returns a new transaction
func NewTx(from, to string, value uint, data string) Tx {
	return Tx{from, to, value, data}
}

// Transaction references an output from previous transaction.
// Transaction is the input output transaction struct
type Transaction struct {
	ID      Snapshot
	Inputs  []TxInput
	Outputs []TxOutput
}

// TxInput struct
type TxInput struct {
	// Ref transaction id
	ID Snapshot
	// Output index. Kinda like 3rd from this id
	Out int
	// connects the input to ref output
	Sig string
}

// TxOutput struct. Token resides here, locked with the public key
type TxOutput struct {
	Value     int
	PublicKey string
}

// CoinbaseTx is the first transaction of the blockchain.
// It is also used during rewards.
func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}
	txin := TxInput{Snapshot{}, -1, data}
	txout := TxOutput{100, to}

	tx := Transaction{Snapshot{}, []TxInput{txin}, []TxOutput{txout}}
	tx.GenerateID()
	return &tx
}

// GenerateID hashes the Tx
func (t *Transaction) GenerateID() {
	var encoded bytes.Buffer
	var hash Snapshot
	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(t)

	if err != nil {
		return
	}
	hash = sha256.Sum256(encoded.Bytes())
	t.ID = hash
}

// IsCoinbase .
func (t *Transaction) IsCoinbase() bool {
	return len(t.Inputs) == 1 && len(t.Inputs[0].ID) == 0 && t.Inputs[0].Out == -1
}

// CanUnlock .
func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

// CanBeUnlocked .
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PublicKey == data
}
