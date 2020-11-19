package data

import (
	"encoding/json"
	"io/ioutil"
)

// Genesis defines the genesis contract of the blockchain
type Genesis struct {
	Balances map[string]uint `json:"balances"`
}

// GenerateGenesis generates the GenTx.
// TODO: Convert this to a CoinbaseTx
func GenerateGenesis() (Tx, error) {
	genesisFilePath := "/home/shivansh_tiwari/go/src/blockchain-go/data/genesis.json"
	content, err := ioutil.ReadFile(genesisFilePath)
	if err != nil {
		return Tx{}, err
	}
	var ret Genesis
	err = json.Unmarshal(content, &ret)
	if err != nil {
		return Tx{}, err
	}
	var tx Tx
	for account, amount := range ret.Balances {
		tx = Tx{account, account, amount, "reward"}
	}
	return tx, nil
}
