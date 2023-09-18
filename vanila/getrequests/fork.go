package getrequests

import (
	"brreq/service"
)

type ForkJSON struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                struct {
		PreviousVersion string `json:"previous_version"`
		CurrentVersion  string `json:"current_version"`
		Epoch           string `json:"epoch"`
	} `json:"data"`
}

type Fork struct {
	service.GetRequest
}

func SpawnFork() service.Get {
	return &Fork{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:3500/eth/v1/beacon/states",
			Response: &ForkJSON{},
		},
	}
}
