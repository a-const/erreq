package prysm

import (
	"erreq/service"
)

type RandaoJSON struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                struct {
		Randao string `json:"randao"`
	} `json:"data"`
}

type Randao struct {
	service.GetRequest
}

func SpawnRandao() service.Get {
	return &Randao{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/beacon/states",
			Response: &RandaoJSON{},
		},
	}
}
