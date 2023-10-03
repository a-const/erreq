package prysm

import (
	"erreq/service"
)

type RewardsBlocksJSON struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                struct {
		ProposerIndex     string `json:"proposer_index"`
		Total             string `json:"total"`
		Attestations      string `json:"attestations"`
		SyncAggregate     string `json:"sync_aggregate"`
		ProposerSlashings string `json:"proposer_slashings"`
		AttesterSlashings string `json:"attester_slashings"`
	} `json:"data"`
}

type RewardsBlocks struct {
	service.GetRequest
}

func SpawnRewardsBlocks() service.Get {
	return &RewardsBlocks{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/beacon/rewards/blocks",
			Response: &RewardsBlocksJSON{},
		},
	}
}
