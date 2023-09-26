package propcounter

import (
	geth "brreq/vanilla/geth/postrequests"
	prysm "brreq/vanilla/prysm/getrequests"
	"fmt"
	"math/big"
	"strings"
)

func (c *Counter) getAllProposersNum(state string) int {
	resp := getValidators.Request([]string{state, "validators"})
	return len(resp.(*prysm.ValidatorsJSON).Data)
}

func (c *Counter) getActivitiesContractsEarnings(from string, to string) {
	// Save earnings at "from" time to slice size of validators num at "to" time
	// Some indices could be created in requested period, we should reserve space for them
	balanceBefore := make([]uint64, c.getAllProposersNum(to))

	resp := getValidators.Request([]string{from, "validators"})
	valNumBefore := c.getAllProposersNum(from)
	for i := 0; i < valNumBefore; i++ {
		balanceBefore[i] = mustParseUInt64(resp.(*prysm.ValidatorsJSON).Data[i].Balance)
	}

	resp = getValidators.Request([]string{to, "validators"})
	for i := 0; i < len(resp.(*prysm.ValidatorsJSON).Data); i++ {
		// Contract
		if resp.(*prysm.ValidatorsJSON).Data[i].Validator.Contract != emptyContract {
			c.Output.Proposers[i].Contract = resp.(*prysm.ValidatorsJSON).Data[i].Validator.Contract
		}
		//Activity
		actv := mustParseUInt64(resp.(*prysm.ValidatorsJSON).Data[i].Validator.EffectiveActivity)
		c.Output.Proposers[i].Activity = actv
		// Earnings
		balanceAfter := mustParseUInt64(resp.(*prysm.ValidatorsJSON).Data[i].Balance)
		c.Output.Proposers[i].Earned = balanceAfter - balanceBefore[i]

		c.Output.Minted.Add(c.Output.Minted, big.NewInt(int64(c.Output.Proposers[i].Earned)))
		writter.Start()
		fmt.Fprintf(writter, "\n[Adding other info] Contracts, activities and earnings...%d out of %d validator   ", i, len(c.Output.Proposers))
		writter.Stop()
	}
}

func (c *Counter) getBurned(blockNum uint64) *big.Int {
	hexed := fmt.Sprintf("0x%x", blockNum)
	resp := getGethBlock.Request(map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_getBlockByNumber",
		"params":  []any{hexed, false},
	})

	var (
		bfString = resp.(*geth.BlockByNumberJSON).Result.BaseFeePerGas
		guString = resp.(*geth.BlockByNumberJSON).Result.GasUsed
		burned   = big.NewInt(0)
	)

	bf, _ := new(big.Int).SetString(strings.TrimPrefix(bfString, "0x"), 16)
	gu, _ := new(big.Int).SetString(strings.TrimPrefix(guString, "0x"), 16)

	burned = burned.Mul(bf, gu)
	return burned

}
