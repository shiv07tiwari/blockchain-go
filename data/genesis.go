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
func GenerateGenesis(path string) (Tx, Genesis, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return Tx{}, Genesis{}, err
	}
	var ret Genesis
	err = json.Unmarshal(content, &ret)
	if err != nil {
		return Tx{}, Genesis{}, err
	}
	var tx Tx
	for account, amount := range ret.Balances {
		tx = Tx{account, account, amount, "Genesis Block"}
	}
	return tx, ret, nil
}
