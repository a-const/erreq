package getrequests

import (
	"brreq/service"
)

type RootJSON struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                struct {
		Root string `json:"root"`
	} `json:"data"`
}

type Root struct {
	service.GetRequest
}

func SpawnRoot() service.Get {
	return &Root{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:3500/eth/v1/beacon/states",
			Response: &RootJSON{},
		},
	}
}
