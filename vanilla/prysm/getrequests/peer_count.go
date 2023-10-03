package prysm

import (
	"erreq/service"
)

type PeerCountJSON struct {
	Data struct {
		Disconnected  string `json:"disconnected"`
		Connecting    string `json:"connecting"`
		Connected     string `json:"connected"`
		Disconnecting string `json:"disconnecting"`
	} `json:"data"`
}

type PeerCount struct {
	service.GetRequest
}

func SpawnPeerCount() service.Get {
	return &PeerCount{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/node/peer_count",
			Response: &PeerCountJSON{},
		},
	}
}
