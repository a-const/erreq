package prysm

import (
	"erreq/service"
)

type ValidatorBalancesJSON struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                []struct {
		Index   string `json:"index"`
		Balance string `json:"balance"`
	} `json:"data"`
}

type ValidatorBalances struct {
	service.GetRequest
}

func SpawnValidatorBalances() service.Get {
	return &ValidatorBalances{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/beacon/states",
			Response: &ValidatorBalancesJSON{},
		},
	}
}
