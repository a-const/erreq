package vanilla

import (
	"erreq/service"
	geth "erreq/vanilla/geth/postrequests"
	prysm "erreq/vanilla/prysm/getrequests"
)

func SpawnGetRequest(endpoint string) service.Get {
	switch endpoint {
	case "prysm_peers":
		return prysm.SpawnPeers()
	case "prysm_syncing":
		return prysm.SpawnSyncing()
	case "prysm_identity":
		return prysm.SpawnIdentity()
	case "prysm_peer_count":
		return prysm.SpawnPeerCount()
	case "prysm_version":
		return prysm.SpawnVersion()
	case "prysm_peerbyid":
		return prysm.SpawnPeerByID()
	case "prysm_genesis":
		return prysm.SpawnGenesis()
	case "prysm_validators":
		return prysm.SpawnValidators()
	case "prysm_root":
		return prysm.SpawnRoot()
	case "prysm_fork":
		return prysm.SpawnFork()
	case "prysm_finality_checkpoints":
		return prysm.SpawnFinalityCheckpoints()
	case "prysm_validator_by_id":
		return prysm.SpawnValidatorByID()
	case "prysm_block_by_id":
		return prysm.SpawnBlockByID()
	case "prysm_validator_balances":
		return prysm.SpawnValidatorBalances()
	case "prysm_sync_committees":
		return prysm.SpawnSyncCommittees()
	case "prysm_rewards_blocks":
		return prysm.SpawnRewardsBlocks()
	case "prysm_randao":
		return prysm.SpawnRandao()
	case "prysm_lightclient_updates":
		return prysm.SpawnLightClientUpdates()
	case "prysm_lightclient_bootstrap":
		return prysm.SpawnLightClientBootstrap()
	case "prysm_header_by_id":
		return prysm.SpawnHeaderByID()
	case "prysm_headers":
		return prysm.SpawnHeaders()
	case "prysm_deposit_snapshot":
		return prysm.SpawnDepositSnapshot()
	case "prysm_blob_sidecars":
		return prysm.SpawnBlobSidecars()
	case "prysm_blinded_blocks":
		return prysm.SpawnBlindedBlocks()
	case "prysm_attestations":
		return prysm.SpawnAttestations()
	}
	return nil
}

func SpawnPostRequest(endpoint string) service.Post {
	switch endpoint {
	case "geth_block_by_number":
		return geth.SpawnBlockByNumer()
	}
	return nil
}
