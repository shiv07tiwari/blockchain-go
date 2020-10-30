package data

import (
	"crypto/sha256"
	"encoding/json"
)

// Block data structure
type Block struct {
	Header BlockHeader
	TXs    []Tx
}

// BlockHeader is block's meta data
type BlockHeader struct {
	Parent Snapshot
	Time   uint64
}

// BlockData .
type BlockData struct {
	Key   Snapshot `json:"snapshot"`
	Value Block    `json:"block"`
}

// Hash creates hash of the block
// now this is an optimized way of hashing as compared to snapshot of
// complete dB
func (b Block) Hash() (Snapshot, error) {
	blockJSON, err := json.Marshal(b)
	return sha256.Sum256(blockJSON), err
}

// NewBlock creates and returns a new block
func NewBlock(snapshot Snapshot, time uint64, Txs []Tx) Block {
	return Block{BlockHeader{snapshot, time}, Txs}
}
