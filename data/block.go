package data

// Block data structure
type Block struct {
	// BlockHeader, which contains the Parent Hash and timestamp
	Header BlockHeader
	// All the transactions in the block
	TXs []Tx
	// Nonce value used in Proof of Work
	Nonce int
}

// BlockHeader is block's meta data
type BlockHeader struct {
	// Parent Hash Vaule to link the blockchain
	Parent Snapshot
	// TimeStamp of the Block
	Time uint64
}

// BlockData .
type BlockData struct {
	// This is the struct that will be stored in the dB.
	// Hash Value of the block
	Key Snapshot `json:"snapshot"`

	// Block object
	Value Block `json:"block"`
}

// NewBlock creates and returns a new block
func NewBlock(snapshot Snapshot, time uint64, Txs []Tx) Block {
	return Block{BlockHeader{snapshot, time}, Txs, 0}
}
