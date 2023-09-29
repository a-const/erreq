package prysm

import (
	"brreq/service"
)

type IdentityJSON struct {
	Data struct {
		PeerID             string   `json:"peer_id"`
		Enr                string   `json:"enr"`
		P2PAddresses       []string `json:"p2p_addresses"`
		DiscoveryAddresses []string `json:"discovery_addresses"`
		Metadata           struct {
			SeqNumber string `json:"seq_number"`
			Attnets   string `json:"attnets"`
		} `json:"metadata"`
	} `json:"data"`
}

type Identity struct {
	service.GetRequest
}

func SpawnIdentity() service.Get {
	return &Identity{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/node/identity",
			Response: &IdentityJSON{},
		},
	}
}
