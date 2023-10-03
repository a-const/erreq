package prysm

import (
	"erreq/service"
)

type LightClientUpdatesJSON []struct {
	Version string `json:"version"`
	Data    struct {
		AttestedHeader struct {
			Beacon struct {
				Slot          string `json:"slot"`
				ProposerIndex string `json:"proposer_index"`
				ParentRoot    string `json:"parent_root"`
				StateRoot     string `json:"state_root"`
				BodyRoot      string `json:"body_root"`
			} `json:"beacon"`
		} `json:"attested_header"`
		NextSyncCommittee struct {
			Pubkeys         []string `json:"pubkeys"`
			AggregatePubkey string   `json:"aggregate_pubkey"`
		} `json:"next_sync_committee"`
		NextSyncCommitteeBranch []string `json:"next_sync_committee_branch"`
		FinalizedHeader         struct {
			Beacon struct {
				Slot          string `json:"slot"`
				ProposerIndex string `json:"proposer_index"`
				ParentRoot    string `json:"parent_root"`
				StateRoot     string `json:"state_root"`
				BodyRoot      string `json:"body_root"`
			} `json:"beacon"`
		} `json:"finalized_header"`
		FinalityBranch []string `json:"finality_branch"`
		SyncAggregate  struct {
			SyncCommitteeBits      string `json:"sync_committee_bits"`
			SyncCommitteeSignature string `json:"sync_committee_signature"`
		} `json:"sync_aggregate"`
		SignatureSlot string `json:"signature_slot"`
	} `json:"data"`
}

type LightClientUpdates struct {
	service.GetRequest
}

func SpawnLightClientUpdates() service.Get {
	return &LightClientUpdates{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/beacon/light_client/updates",
			Response: &[]LightClientUpdatesJSON{},
		},
	}
}
