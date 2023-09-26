package vanilla

import (
	"brreq/service"
	gethPost "brreq/vanilla/geth/postrequests"
	prysmGet "brreq/vanilla/prysm/getrequests"
)

func SpawnGetRequest(endpoint string) service.Get {
	switch endpoint {
	case "prysm_peers":
		return prysmGet.SpawnPeers()
	case "prysm_syncing":
		return prysmGet.SpawnSyncing()
	case "prysm_identity":
		return prysmGet.SpawnIdentity()
	case "prysm_peer_count":
		return prysmGet.SpawnPeerCount()
	case "prysm_version":
		return prysmGet.SpawnVersion()
	case "prysm_peerbyid":
		return prysmGet.SpawnPeerByID()
	case "prysm_genesis":
		return prysmGet.SpawnGenesis()
	case "prysm_validators":
		return prysmGet.SpawnValidators()
	case "prysm_root":
		return prysmGet.SpawnRoot()
	case "prysm_fork":
		return prysmGet.SpawnFork()
	case "prysm_finality_checkpoints":
		return prysmGet.SpawnFinalityCheckpoints()
	case "prysm_validator_by_id":
		return prysmGet.SpawnValidatorByID()
	case "prysm_block_by_id":
		return prysmGet.SpawnBlockByID()
	}
	return nil
}

func SpawnPostRequest(endpoint string) service.Post {
	switch endpoint {
	case "geth_block_by_number":
		return gethPost.SpawnBlockByNumer()
	}
	return nil
}
