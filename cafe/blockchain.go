package main

import "blockchain-go/data"

// Blockchain struct
type Blockchain struct {
	state *data.State
	err   error
}

// NewBlockChain initiates the blockchain
func NewBlockChain() (*Blockchain, error) {
	state, err := data.NewStateFromDisk()
	return &Blockchain{state, err}, err
}

// GetState returns the generates state of Blockchain
func (b *Blockchain) GetState() (*data.State, error) {
	return b.state, b.err
}
