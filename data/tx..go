package data

type Tx struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

func (t *Tx) isReward() bool {
	return t.Data == "reward"
}
