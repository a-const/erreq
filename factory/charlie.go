package fabric

import (
	get "brreq/factory/getrequests"
	"brreq/service"
)

func SpawnGetRequest(endpoint string) service.Get {
	switch endpoint {
	case "peers":
		return get.SpawnPeers()
	case "syncing":
		return get.SpawnSyncing()
	case "identity":
		return get.SpawnIdentity()
	case "peer_count":
		return get.SpawnPeerCount()
	case "version":
		return get.SpawnVersion()
	case "peerbyid":
		return get.SpawnPeerByID()
	case "genesis":
		return get.SpawnGenesis()
	case "validators":
		return get.SpawnValidators()
	case "root":
		return get.SpawnRoot()
	case "fork":
		return get.SpawnFork()
	case "finality_checkpoints":
		return get.SpawnFinalityCheckpoints()
	case "validator_by_id":
		return get.SpawnValidatorByID()
	case "block_by_id":
		return get.SpawnBlockByID()
	}
	return nil
}
