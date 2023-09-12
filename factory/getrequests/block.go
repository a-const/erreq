package getrequests

import (
	"brreq/service"
)

type BlockByIDJSON struct {
	Version             string `json:"version"`
	ExecutionOptimistic bool   `json:"execution_optimistic"`
	Finalized           bool   `json:"finalized"`
	Data                struct {
		Message struct {
			Slot          string `json:"slot"`
			ProposerIndex string `json:"proposer_index"`
			ParentRoot    string `json:"parent_root"`
			StateRoot     string `json:"state_root"`
			Body          struct {
				RandaoReveal string `json:"randao_reveal"`
				Eth1Data     struct {
					DepositRoot  string `json:"deposit_root"`
					DepositCount string `json:"deposit_count"`
					BlockHash    string `json:"block_hash"`
				} `json:"eth1_data"`
				Graffiti          string `json:"graffiti"`
				ProposerSlashings []struct {
					SignedHeader1 struct {
						Message struct {
							Slot          string `json:"slot"`
							ProposerIndex string `json:"proposer_index"`
							ParentRoot    string `json:"parent_root"`
							StateRoot     string `json:"state_root"`
							BodyRoot      string `json:"body_root"`
						} `json:"message"`
						Signature string `json:"signature"`
					} `json:"signed_header_1"`
					SignedHeader2 struct {
						Message struct {
							Slot          string `json:"slot"`
							ProposerIndex string `json:"proposer_index"`
							ParentRoot    string `json:"parent_root"`
							StateRoot     string `json:"state_root"`
							BodyRoot      string `json:"body_root"`
						} `json:"message"`
						Signature string `json:"signature"`
					} `json:"signed_header_2"`
				} `json:"proposer_slashings"`
				AttesterSlashings []struct {
					Attestation1 struct {
						AttestingIndices []string `json:"attesting_indices"`
						Signature        string   `json:"signature"`
						Data             struct {
							Slot            string `json:"slot"`
							Index           string `json:"index"`
							BeaconBlockRoot string `json:"beacon_block_root"`
							Source          struct {
								Epoch string `json:"epoch"`
								Root  string `json:"root"`
							} `json:"source"`
							Target struct {
								Epoch string `json:"epoch"`
								Root  string `json:"root"`
							} `json:"target"`
						} `json:"data"`
					} `json:"attestation_1"`
					Attestation2 struct {
						AttestingIndices []string `json:"attesting_indices"`
						Signature        string   `json:"signature"`
						Data             struct {
							Slot            string `json:"slot"`
							Index           string `json:"index"`
							BeaconBlockRoot string `json:"beacon_block_root"`
							Source          struct {
								Epoch string `json:"epoch"`
								Root  string `json:"root"`
							} `json:"source"`
							Target struct {
								Epoch string `json:"epoch"`
								Root  string `json:"root"`
							} `json:"target"`
						} `json:"data"`
					} `json:"attestation_2"`
				} `json:"attester_slashings"`
				Attestations []struct {
					AggregationBits string `json:"aggregation_bits"`
					Signature       string `json:"signature"`
					Data            struct {
						Slot            string `json:"slot"`
						Index           string `json:"index"`
						BeaconBlockRoot string `json:"beacon_block_root"`
						Source          struct {
							Epoch string `json:"epoch"`
							Root  string `json:"root"`
						} `json:"source"`
						Target struct {
							Epoch string `json:"epoch"`
							Root  string `json:"root"`
						} `json:"target"`
					} `json:"data"`
				} `json:"attestations"`
				Deposits []struct {
					Proof []string `json:"proof"`
					Data  struct {
						Pubkey                string `json:"pubkey"`
						WithdrawalCredentials string `json:"withdrawal_credentials"`
						Amount                string `json:"amount"`
						Signature             string `json:"signature"`
					} `json:"data"`
				} `json:"deposits"`
				VoluntaryExits []struct {
					Message struct {
						Epoch          string `json:"epoch"`
						ValidatorIndex string `json:"validator_index"`
					} `json:"message"`
					Signature string `json:"signature"`
				} `json:"voluntary_exits"`
			} `json:"body"`
		} `json:"message"`
		Signature string `json:"signature"`
	} `json:"data"`
}

type BlockByID struct {
	service.GetRequest
}

func SpawnBlockByID() service.Get {
	return &BlockByID{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:3500/eth/v2/beacon/blocks",
			Response: &BlockByIDJSON{},
		},
	}
}
