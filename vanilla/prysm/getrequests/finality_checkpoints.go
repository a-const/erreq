package prysm

import (
	"erreq/service"
)

type FinalityCheckpointsJSON struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                struct {
		PreviousJustified struct {
			Epoch string `json:"epoch"`
			Root  string `json:"root"`
		} `json:"previous_justified"`
		CurrentJustified struct {
			Epoch string `json:"epoch"`
			Root  string `json:"root"`
		} `json:"current_justified"`
		Finalized struct {
			Epoch string `json:"epoch"`
			Root  string `json:"root"`
		} `json:"finalized"`
	} `json:"data"`
}

type FinalityCheckpoints struct {
	service.GetRequest
}

func SpawnFinalityCheckpoints() service.Get {
	return &FinalityCheckpoints{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/beacon/states",
			Response: &FinalityCheckpointsJSON{},
		},
	}
}
