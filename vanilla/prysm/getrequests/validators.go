package prysm

import (
	"brreq/service"
)

type ValidatorsJSON struct {
	Data []struct {
		Index     string `json:"index"`
		Balance   string `json:"balance"`
		Status    string `json:"status"`
		Validator struct {
			Pubkey                     string `json:"pubkey"`
			WithdrawalCredentials      string `json:"withdrawal_credentials"`
			Contract                   string `json:"contract"`
			EffectiveBalance           string `json:"effective_balance"`
			EffectiveActivity          string `json:"effective_activity"`
			Slashed                    bool   `json:"slashed"`
			ActivationEligibilityEpoch string `json:"activation_eligibility_epoch"`
			ActivationEpoch            string `json:"activation_epoch"`
			ExitEpoch                  string `json:"exit_epoch"`
			WithdrawableEpoch          string `json:"withdrawable_epoch"`
		} `json:"validator"`
	} `json:"data"`
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
}

type ValidatorByIDJSON struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                struct {
		Index     string `json:"index"`
		Balance   string `json:"balance"`
		Status    string `json:"status"`
		Validator struct {
			Pubkey                     string `json:"pubkey"`
			WithdrawalCredentials      string `json:"withdrawal_credentials"`
			EffectiveBalance           string `json:"effective_balance"`
			Slashed                    bool   `json:"slashed"`
			ActivationEligibilityEpoch string `json:"activation_eligibility_epoch"`
			ActivationEpoch            string `json:"activation_epoch"`
			ExitEpoch                  string `json:"exit_epoch"`
			WithdrawableEpoch          string `json:"withdrawable_epoch"`
		} `json:"validator"`
	} `json:"data"`
}

type Validators struct {
	service.GetRequest
}

type ValidatorByID struct {
	service.GetRequest
}

func SpawnValidators() service.Get {
	return &Validators{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/beacon/states",
			Response: &ValidatorsJSON{},
		},
	}
}

func SpawnValidatorByID() service.Get {
	return &ValidatorByID{
		GetRequest: service.GetRequest{
			Url:      "http://127.0.0.1:",
			Endpoint: "/eth/v1/beacon/states",
			Response: &ValidatorByIDJSON{},
		},
	}
}
