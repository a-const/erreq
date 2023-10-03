package prysm

import (
	"erreq/service"
)

type PeersJSON struct {
	Data []struct {
		PeerID             string `json:"peer_id"`
		Enr                string `json:"enr"`
		LastSeenP2PAddress string `json:"last_seen_p2p_address"`
		State              string `json:"state"`
		Direction          string `json:"direction"`
	} `json:"data"`
}

type PeerByIDJSON struct {
	Data struct {
		PeerID             string `json:"peer_id"`
		Enr                string `json:"enr"`
		LastSeenP2PAddress string `json:"last_seen_p2p_address"`
		State              string `json:"state"`
		Direction          string `json:"direction"`
	} `json:"data"`
}
type Peers struct {
	service.GetRequest
}

func SpawnPeers() service.Get {
	return &Peers{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/node/peers",
			Response: &PeersJSON{},
		},
	}
}

func SpawnPeerByID() service.Get {
	return &Peers{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/node/peers",
			Response: &PeerByIDJSON{},
		},
	}
}
