package geth

import (
	"erreq/service"
)

type BlockByNumberJSON struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		BaseFeePerGas    string `json:"baseFeePerGas"`
		Difficulty       string `json:"difficulty"`
		ExtraData        string `json:"extraData"`
		GasLimit         string `json:"gasLimit"`
		GasUsed          string `json:"gasUsed"`
		Hash             string `json:"hash"`
		LogsBloom        string `json:"logsBloom"`
		Miner            string `json:"miner"`
		MixHash          string `json:"mixHash"`
		Nonce            string `json:"nonce"`
		Number           string `json:"number"`
		ParentHash       string `json:"parentHash"`
		ReceiptsRoot     string `json:"receiptsRoot"`
		Sha3Uncles       string `json:"sha3Uncles"`
		Size             string `json:"size"`
		StateRoot        string `json:"stateRoot"`
		Timestamp        string `json:"timestamp"`
		TotalDifficulty  string `json:"totalDifficulty"`
		Transactions     []any  `json:"transactions"`
		TransactionsRoot string `json:"transactionsRoot"`
		Uncles           []any  `json:"uncles"`
		Withdrawals      []struct {
			Index          string `json:"index"`
			ValidatorIndex string `json:"validatorIndex"`
			Address        string `json:"address"`
			Amount         string `json:"amount"`
		} `json:"withdrawals"`
		WithdrawalsRoot string `json:"withdrawalsRoot"`
	} `json:"result"`
}

type BlockByNumber struct {
	service.PostRequest
}

func SpawnBlockByNumer() service.Post {
	return &BlockByNumber{
		PostRequest: service.PostRequest{
			Url:      "http://127.0.0.1:",
			Response: &BlockByNumberJSON{},
		},
	}
}
