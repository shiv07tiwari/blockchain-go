package data

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
)

// We force the network to work to add the block into the blockchain
// Work is heavy, and validation is simple
// Create a nonce that starts at 0
// Data + nonce hashed should meet the set of requirements.

const difficulty = 11

// ProofOfWork struct
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// Data struct
type Data struct {
	prevHash   Snapshot
	TXs        []Tx
	nonce      []byte
	difficulty []byte
}

// CreateNewProof creates a new ProofOfWork for the Block
func CreateNewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(265-difficulty))
	pow := &ProofOfWork{b, target}
	return pow
}

// InitData returns a new data block
func (pow *ProofOfWork) InitData(nonce int) Data {
	data := Data{pow.Block.Header.Parent, pow.Block.TXs, toHex(int64(nonce)),
		toHex(int64(difficulty))}

	return data
}

// Mine mines the block
func (pow *ProofOfWork) Mine() (int, Snapshot) {
	var intHash big.Int
	var hash Snapshot

	nonce := 0

	for nonce < 3 {
		data := pow.InitData(nonce)
		dataJSON, err := json.Marshal(data)
		fmt.Println(data)
		fmt.Println("and")
		fmt.Println(dataJSON)
		if err != nil {
			log.Panic(err)
		}
		hash = sha256.Sum256(dataJSON)
		fmt.Printf("\r%x", hash)
		fmt.Println("")
		intHash.SetBytes(hash[:])
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash
}

func toHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
