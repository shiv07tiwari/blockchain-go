package data

import (
	"encoding/json"
	"io/ioutil"
)

// Genesis defines the genesis contract of the blockchain
type Genesis struct {
	Balances map[string]uint `json:"balances"`
	Symbol   string          `json:"symbol"`
}

func loadGenesis(path string) (Genesis, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return Genesis{}, err
	}
	var ret Genesis
	err = json.Unmarshal(content, &ret)
	if err != nil {
		return Genesis{}, err
	}
	return ret, nil
}
