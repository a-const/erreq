package analytics

import (
	"erreq/vanilla"
	"math/big"

	"github.com/gosuri/uilive"
)

type Withdrawal struct {
	Index          string `json:"index"`
	ValidatorIndex string `json:"validatorIndex"`
	Address        string `json:"address"`
	Amount         string `json:"amount"`
}

type ValidatorsEarningsCapellaJSON struct {
	IsCredentialsSet bool         `json:"is_credential_set"`
	Withdrawals      []Withdrawal `json:"withdrawals"`
}

type ProposersJSON struct {
	From         int             `json:"from"`
	To           int             `json:"to"`
	MaxProposed  int             `json:"max_proposed"`
	Burned       *big.Int        `json:"burned"`
	Minted       *big.Int        `json:"minted"`
	ProposersNum int             `json:"proposers_num"`
	Proposers    []*ProposerJSON `json:"proposers"`
}

type ProposerJSON struct {
	Index    int    `json:"index"`
	Counter  int    `json:"counter"`
	Earned   uint64 `json:"earned"`
	Contract string `json:"contract,omitempty"`
	Activity uint64 `json:"activity,omitempty"`
}

type Collector struct {
	Output *ProposersJSON
	Port   string
	From   int
	To     int
	Fork   string
}

var (
	getBlockByID  = vanilla.SpawnGetRequest("prysm_block_by_id")
	getValidators = vanilla.SpawnGetRequest("prysm_validators")
	getGethBlock  = vanilla.SpawnPostRequest("geth_block_by_number")
	getFork       = vanilla.SpawnGetRequest("prysm_fork")
	writter       = uilive.New()
	emptyContract = "0x0000000000000000000000000000000000000000"
)

const (
	capellaVersion   = "0x03"
	bellatrixVersion = "0x02"
	cross            = "0x00"
)
