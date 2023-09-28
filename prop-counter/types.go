package propcounter

import (
	"brreq/vanilla"
	"math/big"

	"github.com/gosuri/uilive"
)

type ProposersJSON struct {
	From        int             `json:"from"`
	To          int             `json:"to"`
	MaxProposed int             `json:"max_proposed"`
	Burned      *big.Int        `json:"burned"`
	Minted      *big.Int        `json:"minted"`
	Proposers   []*ProposerJSON `json:"proposers"`
}

type ProposerJSON struct {
	Index    int    `json:"index"`
	Counter  int    `json:"counter"`
	Earned   uint64 `json:"earned"`
	Contract string `json:"contract,omitempty"`
	Activity uint64 `json:"activity,omitempty"`
}

type Counter struct {
	Output *ProposersJSON
	Port   string
}

var (
	getBlockByID  = vanilla.SpawnGetRequest("prysm_block_by_id")
	getValidators = vanilla.SpawnGetRequest("prysm_validators")
	getGethBlock  = vanilla.SpawnPostRequest("geth_block_by_number")
	writter       = uilive.New()
	emptyContract = "0x0000000000000000000000000000000000000000"
)
