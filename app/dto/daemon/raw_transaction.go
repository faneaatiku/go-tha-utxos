package daemon

type RawTransactionOutput map[string]float64

type SignRawTransactionResponse struct {
	Hex      string `json:"hex"`
	Complete bool   `json:"complete"`
}

type RawTransactionInput struct {
	Txid string `json:"txid"`
	Vout int    `json:"vout"`
}

type RpcSignRawTransactionResponse struct {
	Result SignRawTransactionResponse `json:"result"`
}
