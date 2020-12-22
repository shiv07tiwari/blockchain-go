package data

import (
	"encoding/json"
	"io/ioutil"
)

// Genesis defines the genesis contract of the blockchain
type Genesis struct {
	Balances map[string]int `json:"balances"`
}

// GenerateGenesis generates the GenTx.
func GenerateGenesis() (Tx, error) {

	// Path containing the genesis contract file
	genesisFilePath := "/home/shivansh_tiwari/go/src/blockchain-go/data/genesis.json"

	// ioutil.ReadFile returns a byte array and error, if any
	content, err := ioutil.ReadFile(genesisFilePath)

	// Handle the error
	if err != nil {
		return Tx{}, err
	}

	var ret Genesis

	// Parse the JSON encoded data of the file to ret
	err = json.Unmarshal(content, &ret)

	// Handle the error
	if err != nil {
		return Tx{}, err
	}

	var tx Tx
	// For the account-amount pair in genesis JSON, create a new CoinBaseTx
	for account, amount := range ret.Balances {
		tx = CoinbaseTx(account, "genesis", amount)
	}

	// return the genesis transaction
	return tx, nil
}
