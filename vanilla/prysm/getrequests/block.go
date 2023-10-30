package prysm

import (
	"erreq/service"
)

type BlockByIDJSON struct {
	Version string `json:"version"`
	Data    struct {
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
				ProposerSlashings []any  `json:"proposer_slashings"`
				AttesterSlashings []any  `json:"attester_slashings"`
				Attestations      []struct {
					AggregationBits string `json:"aggregation_bits"`
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
					Signature string `json:"signature"`
				} `json:"attestations"`
				Deposits          []any  `json:"deposits"`
				VoluntaryExits    []any  `json:"voluntary_exits"`
				ActivityChanges   []any  `json:"activity_changes"`
				TransactionsCount string `json:"transactions_count"`
				BaseFee           string `json:"base_fee"`
				SyncAggregate     struct {
					SyncCommitteeBits      string `json:"sync_committee_bits"`
					SyncCommitteeSignature string `json:"sync_committee_signature"`
				} `json:"sync_aggregate"`
				ExecutionPayload struct {
					ParentHash    string `json:"parent_hash"`
					FeeRecipient  string `json:"fee_recipient"`
					StateRoot     string `json:"state_root"`
					ReceiptsRoot  string `json:"receipts_root"`
					LogsBloom     string `json:"logs_bloom"`
					PrevRandao    string `json:"prev_randao"`
					BlockNumber   string `json:"block_number"`
					GasLimit      string `json:"gas_limit"`
					GasUsed       string `json:"gas_used"`
					Timestamp     string `json:"timestamp"`
					ExtraData     string `json:"extra_data"`
					BaseFeePerGas string `json:"base_fee_per_gas"`
					BlockHash     string `json:"block_hash"`
					Transactions  []any  `json:"transactions"`
					Withdrawals   []struct {
						Index          string `json:"index"`
						ValidatorIndex string `json:"validator_index"`
						Address        string `json:"address"`
						Amount         string `json:"amount"`
					} `json:"withdrawals"`
				} `json:"execution_payload"`
				BlsToExecutionChanges []any `json:"bls_to_execution_changes"`
			} `json:"body"`
		} `json:"message"`
		Signature string `json:"signature"`
	} `json:"data"`
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
}

type BlockByID struct {
	service.GetRequest
}

func SpawnBlockByID() service.Get {
	return &BlockByID{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v2/beacon/blocks",
			Response: &BlockByIDJSON{},
		},
	}
}
