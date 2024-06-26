package prysm

import (
	"erreq/service"
)

type GenesisJSON struct {
	Data struct {
		GenesisTime           string `json:"genesis_time"`
		GenesisValidatorsRoot string `json:"genesis_validators_root"`
		GenesisForkVersion    string `json:"genesis_fork_version"`
	} `json:"data"`
}

type Genesis struct {
	service.GetRequest
}

func SpawnGenesis() service.Get {
	return &Genesis{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/beacon/genesis",
			Response: &GenesisJSON{},
		},
	}
}
