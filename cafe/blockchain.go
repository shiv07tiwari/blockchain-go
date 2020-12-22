package main

import (
	"blockchain-go/data"
)

// Blockchain struct
type Blockchain struct {
	// Variable stores the state of the blockchain
	state *data.State
	// Stores the error, if any while creating a new state
	err error
}

// NewBlockChain initiates the blockchain
func NewBlockChain() (*Blockchain, error) {

	// ReRun all the transactions in the database to create a new state.
	state, err := data.NewStateFromDisk()
	return &Blockchain{state, err}, err
}

// GetState returns the generates state of Blockchain
func (b *Blockchain) GetState() (*data.State, error) {
	return b.state, b.err
}
