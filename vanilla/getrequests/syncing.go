package getrequests

import (
	"brreq/service"
)

type SyncingJSON struct {
	Data struct {
		HeadSlot     string `json:"head_slot"`
		SyncDistance string `json:"sync_distance"`
		IsSyncing    bool   `json:"is_syncing"`
		IsOptimistic bool   `json:"is_optimistic"`
		ElOffline    bool   `json:"el_offline"`
	} `json:"data"`
}

type Syncing struct {
	service.GetRequest
}

func SpawnSyncing() service.Get {
	return &Syncing{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:3500/eth/v1/node/syncing",
			Response: &SyncingJSON{},
		},
	}
}