package prysm

import (
	"erreq/service"
)

type VersionJSON struct {
	Data struct {
		Version string `json:"version"`
	} `json:"data"`
}

type Version struct {
	service.GetRequest
}

func SpawnVersion() service.Get {
	return &Version{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/node/version",
			Response: &VersionJSON{},
		},
	}
}
