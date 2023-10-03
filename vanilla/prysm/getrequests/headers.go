package prysm

import (
	"erreq/service"
)

type HeaderByIDJSON struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                struct {
		Root      string `json:"root"`
		Canonical bool   `json:"canonical"`
		Header    struct {
			Message struct {
				Slot          string `json:"slot"`
				ProposerIndex string `json:"proposer_index"`
				ParentRoot    string `json:"parent_root"`
				StateRoot     string `json:"state_root"`
				BodyRoot      string `json:"body_root"`
			} `json:"message"`
			Signature string `json:"signature"`
		} `json:"header"`
	} `json:"data"`
}

type HeadersJSON struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                []struct {
		Root      string `json:"root"`
		Canonical bool   `json:"canonical"`
		Header    struct {
			Message struct {
				Slot          string `json:"slot"`
				ProposerIndex string `json:"proposer_index"`
				ParentRoot    string `json:"parent_root"`
				StateRoot     string `json:"state_root"`
				BodyRoot      string `json:"body_root"`
			} `json:"message"`
			Signature string `json:"signature"`
		} `json:"header"`
	} `json:"data"`
}

type HeaderByID struct {
	service.GetRequest
}

type Headers struct {
	service.GetRequest
}

func SpawnHeaderByID() service.Get {
	return &HeaderByID{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/beacon/headers",
			Response: &HeaderByIDJSON{},
		},
	}
}

func SpawnHeaders() service.Get {
	return &Headers{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/beacon/headers",
			Response: &HeadersJSON{},
		},
	}
}
