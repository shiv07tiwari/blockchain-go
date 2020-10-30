package data

// Tx trasaction data structure
type Tx struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

func (t *Tx) isReward() bool {
	return t.Data == "reward"
}

// NewTx returns a new transaction
func NewTx(from, to string, value uint, data string) Tx {
	return Tx{from, to, value, data}
}
