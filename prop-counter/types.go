package propcounter

import (
	"brreq/vanilla"
)

type ProposersJSON struct {
	From        int             `json:"from"`
	To          int             `json:"to"`
	MaxProposed int             `json:"max_proposed"`
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
	Ctr    map[int]int
}

var (
	getBlockByID  = vanilla.SpawnGetRequest("block_by_id")
	getValidators = vanilla.SpawnGetRequest("validators")
	//getValidatorByID = vanila.SpawnGetRequest("validator_by_id")
	//writter          = uilive.New()
	emptyContract = "0x0000000000000000000000000000000000000000"
)
