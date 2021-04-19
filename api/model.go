package api

type CloudflareRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

type Block struct {
	TxByHash     map[string][]byte `json:"-"`
	TxByIndex    [][]byte          `json:"-"`
	BlockBytes   []byte            `json:"-"`
	Number       string            `json:"number"`
	Transactions []*Transaction    `json:"transactions"`
}

type Response struct {
	ID    int    `json:"id"`
	Block *Block `json:"result"`
}
