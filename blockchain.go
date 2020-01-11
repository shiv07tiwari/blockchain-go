package main

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
)
const dbFile = "blockchain4.db"
const blocksBucket = "blocks"

// Contains tip and databse
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// To access all the blocks one by one without loading them all iun memory.
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	// A read only to get the Hash of last block
	bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	// Adding a new block
	newBlock := NewBlock(data, lastHash)

	// Updating the tip and putting new block in Bucket
	 bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		 b.Put(newBlock.Hash, newBlock.Serialize())
		 b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})

}

// Create a new Blockchain
func NewBlockchain() *Blockchain {

	var tip []byte

	// Open the database, create if doesn't exist already
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	// Making a genesis block and updating it in the bucket
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		 // if bucket doesnt exist, make one
		if (b == nil) {
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	bc := Blockchain{tip, db}
	return &bc
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block
	fmt.Println("CurrentHash", i.currentHash)

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}

