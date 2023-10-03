package prysm

import (
	"erreq/service"
)

type BlobSidecarsJSON struct {
	Data []struct {
		BlockRoot       string `json:"block_root"`
		Index           string `json:"index"`
		Slot            string `json:"slot"`
		BlockParentRoot string `json:"block_parent_root"`
		ProposerIndex   string `json:"proposer_index"`
		Blob            string `json:"blob"`
		KzgCommitment   string `json:"kzg_commitment"`
		KzgProof        string `json:"kzg_proof"`
	} `json:"data"`
}

type BlobSidecars struct {
	service.GetRequest
}

func SpawnBlobSidecars() service.Get {
	return &BlobSidecars{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/beacon/blob_sidecars",
			Response: &BlobSidecarsJSON{},
		},
	}
}
