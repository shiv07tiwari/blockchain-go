package data

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Basic logic of proof of work :-
// We force the network to work to add the block into the blockchain
// Work is heavy, and validation is simple
// Create a nonce that starts at 0
// Data + nonce hashed should meet the set of requirements.

const difficulty = 16

// ProofOfWork struct
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// POWData struct
type POWData struct {
	PrevHash   Snapshot `json:"prevHash"`
	TXs        []Tx     `json:"transactions"`
	Nonce      int      `json:"nonce"`
	Difficulty int      `json:"difficulty"`
}

// CreateNewProof creates a new ProofOfWork for the Block
func CreateNewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)

	// Left shift the value of 1 as per the difficulty
	target.Lsh(target, uint(265-difficulty))
	pow := &ProofOfWork{b, target}
	return pow
}

// InitData returns a new data block
func (pow *ProofOfWork) InitData(nonce int) POWData {
	powData := POWData{pow.Block.Header.Parent, pow.Block.TXs, nonce,
		difficulty}

	return powData
}

// Mine mines the block
func (pow *ProofOfWork) Mine() (int, Snapshot) {
	var intHash big.Int
	var hash Snapshot

	nonce := 0

	// Iterate over for the value of Nonce
	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		dataJSON, err := json.Marshal(data)
		if err != nil {
			log.Panic(err)
		}

		hash = sha256.Sum256(dataJSON)
		intHash.SetBytes(hash[:])

		// Set the valid hash if the condition is satisfied
		if intHash.Cmp(pow.Target) == -1 {

			fmt.Println("Valid Hash : ")
			fmt.Printf("\r%x", hash)
			fmt.Println("")
			break
		} else {
			// Increment the Nonce
			nonce++
		}
	}
	return nonce, hash
}

// Validate validates a POW
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonce)
	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(dataJSON)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
