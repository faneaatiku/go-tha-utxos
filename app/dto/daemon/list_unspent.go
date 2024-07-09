package daemon

type Unspent struct {
	Txid          string        `json:"txid"`
	Vout          int           `json:"vout"`
	Address       string        `json:"address"`
	Label         string        `json:"label"`
	ScriptPubKey  string        `json:"scriptPubKey"`
	Amount        float64       `json:"amount"`
	Confirmations int           `json:"confirmations"`
	Spendable     bool          `json:"spendable"`
	Solvable      bool          `json:"solvable"`
	Desc          string        `json:"desc"`
	ParentDescs   []interface{} `json:"parent_descs"`
	Safe          bool          `json:"safe"`
}

type ListUnspentRequest struct {
	MinimumAmount float64  `json:"minimumAmount"`
	MaximumAmount *float64 `json:"maximumAmount,omitempty"`
	MaximumCount  int      `json:"maximumCount"`
}
