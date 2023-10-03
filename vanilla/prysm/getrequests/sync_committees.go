package prysm

import (
	"erreq/service"
)

type SyncCommitteesJSON struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                struct {
		Validators          []string   `json:"validators"`
		ValidatorAggregates [][]string `json:"validator_aggregates"`
	} `json:"data"`
}

type SyncCommittees struct {
	service.GetRequest
}

func SpawnSyncCommittees() service.Get {
	return &SyncCommittees{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/beacon/states",
			Response: &SyncCommitteesJSON{},
		},
	}
}
