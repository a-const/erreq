package prysm

import (
	"erreq/service"
)

type AttestationsJSON struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                []struct {
		AggregationBits string `json:"aggregation_bits"`
		Signature       string `json:"signature"`
		Data            struct {
			Slot            string `json:"slot"`
			Index           string `json:"index"`
			BeaconBlockRoot string `json:"beacon_block_root"`
			Source          struct {
				Epoch string `json:"epoch"`
				Root  string `json:"root"`
			} `json:"source"`
			Target struct {
				Epoch string `json:"epoch"`
				Root  string `json:"root"`
			} `json:"target"`
		} `json:"data"`
	} `json:"data"`
}

type Attestations struct {
	service.GetRequest
}

func SpawnAttestations() service.Get {
	return &Attestations{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/beacon/blocks",
			Response: &AttestationsJSON{},
		},
	}
}
