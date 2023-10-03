package prysm

import (
	"erreq/service"
)

type DepositSnapshotJSON struct {
	Data struct {
		Finalized            []string `json:"finalized"`
		DepositRoot          string   `json:"deposit_root"`
		DepositCount         string   `json:"deposit_count"`
		ExecutionBlockHash   string   `json:"execution_block_hash"`
		ExecutionBlockHeight string   `json:"execution_block_height"`
	} `json:"data"`
}

type DepositSnapshot struct {
	service.GetRequest
}

func SpawnDepositSnapshot() service.Get {
	return &DepositSnapshot{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/beacon/deposit_snapshot",
			Response: &DepositSnapshotJSON{},
		},
	}
}
