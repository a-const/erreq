package prysm

import (
	"erreq/service"
)

type LightClientBootstrapJSON struct {
	Version string `json:"version"`
	Data    struct {
		Header struct {
			Beacon struct {
				Slot          string `json:"slot"`
				ProposerIndex string `json:"proposer_index"`
				ParentRoot    string `json:"parent_root"`
				StateRoot     string `json:"state_root"`
				BodyRoot      string `json:"body_root"`
			} `json:"beacon"`
		} `json:"header"`
		CurrentSyncCommittee struct {
			Pubkeys         []string `json:"pubkeys"`
			AggregatePubkey string   `json:"aggregate_pubkey"`
		} `json:"current_sync_committee"`
		CurrentSyncCommitteeBranch []string `json:"current_sync_committee_branch"`
	} `json:"data"`
}

type LightClientBootstrap struct {
	service.GetRequest
}

func SpawnLightClientBootstrap() service.Get {
	return &LightClientBootstrap{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/beacon/light_client/bootstrap",
			Response: &LightClientBootstrapJSON{},
		},
	}
}
