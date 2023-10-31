package analytics

import (
	"erreq/service"
	geth "erreq/vanilla/geth/postrequests"
	prysm "erreq/vanilla/prysm/getrequests"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (c *Collector) getAllProposersNum(state string) int {
	resp := getValidators.Request([]string{state, "validators"}, c.Port)
	return len(resp.(*prysm.ValidatorsJSON).Data)
}

func mustBeType(t any) bool {
	switch t.(type) {
	case *service.ErrorHandler:
		return false
	default:
		return true
	}
}

func (c *Collector) getBurned(blockNum uint64) *big.Int {
	hexed := fmt.Sprintf("0x%x", blockNum)
	resp := getGethBlock.Request(map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_getBlockByNumber",
		"params":  []any{hexed, false},
	}, "8545")

	var (
		bfString = resp.(*geth.BlockByNumberJSON).Result.BaseFeePerGas
		guString = resp.(*geth.BlockByNumberJSON).Result.GasUsed
		burned   = big.NewInt(0)
	)

	bf, _ := new(big.Int).SetString(strings.TrimPrefix(bfString, "0x"), 16)
	gu, _ := new(big.Int).SetString(strings.TrimPrefix(guString, "0x"), 16)

	burned.Mul(bf, gu)
	return burned

}

func (c *Collector) getFork(state string) (string, string) {
	resp := getFork.Request([]string{state, "fork"}, c.Port)
	var (
		ver   string
		epoch string
	)
	switch t := resp.(type) {
	case *prysm.ForkJSON:
		ver = t.Data.CurrentVersion
		epoch = t.Data.Epoch
	case *service.ErrorHandler:
		log.Fatal("error retrieving fork version!")
	}
	return ver, epoch
}

func (c *Collector) getBlockFromSlot(slot int) string {
	var block any
	for block = getBlockByID.Request([]string{fmt.Sprintf("%d", slot)}, c.Port); !mustBeType(block); block = getBlockByID.Request([]string{fmt.Sprintf("%d", slot)}, c.Port) {
		log.Error("Can't retrieve eth1 block from beacon slot - requested slot doesn't exist. Trying out next slot")
		slot++
	}
	return block.(*prysm.BlockByIDJSON).Data.Message.Body.ExecutionPayload.BlockNumber
}

func (c *Collector) ActivitiesContractsEarnings_Bellatrix(from string, to string) {
	// Save earnings at "from" time to slice size of validators num at "to" time
	// Some indices could be created in requested period, we should reserve space for them
	balanceBefore := make([]uint64, c.Output.ProposersNum)

	resp := getValidators.Request([]string{from, "validators"}, c.Port)
	valNumBefore := c.Output.ProposersNum
	for i := 0; i < valNumBefore; i++ {
		balanceBefore[i] = mustParseUInt64(resp.(*prysm.ValidatorsJSON).Data[i].Balance)
	}

	resp = getValidators.Request([]string{to, "validators"}, c.Port)
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
		c.Output.Minted.Add(c.Output.Minted, new(big.Int).SetUint64(c.Output.Proposers[i].Earned))

		writter.Start()
		fmt.Fprintf(writter, "\n[Bellatrix][Adding other info] Contracts, activities and earnings...%d out of %d validator   \n", i, len(c.Output.Proposers))
		writter.Stop()
	}
}

func (c *Collector) firstWitdrawalsContractsActivity(from string, to string) []*ValidatorsEarningsCapellaJSON {
	var (
		firstWdFill int
	)

	valNum := c.Output.ProposersNum
	forkSlot := mustAtoi(c.Fork)*32 + 1
	forkBlock := mustAtoi(c.getBlockFromSlot(forkSlot))
	toBlock := forkBlock + (c.To - c.From)
	earnings := make([]*ValidatorsEarningsCapellaJSON, c.Output.ProposersNum)

	validators := getValidators.Request([]string{to, "validators"}, c.Port)
	if !mustBeType(validators) {
		log.Fatal("Error receiving list of validators")
	}

	for i, v := range validators.(*prysm.ValidatorsJSON).Data {
		earnings[i] = &ValidatorsEarningsCapellaJSON{
			IsCredentialsSet: false,
			Withdrawals:      []Withdrawal{},
		}
		writter.Start()
		fmt.Fprintf(writter, "\n[Capella][Searching for withdrawal credentials] Checked %d out of %d validators...   \n", i, len(earnings))
		writter.Stop()

		// Contract
		if v.Validator.Contract != emptyContract {
			c.Output.Proposers[i].Contract = v.Validator.Contract
		}
		//Activity
		actv := mustParseUInt64(v.Validator.EffectiveActivity)
		c.Output.Proposers[i].Activity = actv

		if v.Validator.WithdrawalCredentials[0:5] == "0x010" {
			earnings[i].IsCredentialsSet = true
			continue
		}
		firstWdFill++

	}

	writter.Start()
	fmt.Fprint(writter, "\n[Capella][Searching for witdrawal credentials] Done!   \n")
	writter.Stop()

	for bl := forkBlock; bl < toBlock && firstWdFill < valNum-1; bl++ {
		writter.Start()
		fmt.Fprintf(writter, "\n[Capella][Searching for first withdrawals] %d out of %d validators...   \n", firstWdFill, valNum)
		writter.Stop()

		hexed := fmt.Sprintf("0x%x", bl)
		block := getGethBlock.Request(map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "eth_getBlockByNumber",
			"params":  []any{hexed, false},
		}, "8545")
		if !mustBeType(block) {
			continue
		}

		for _, wd := range block.(*geth.BlockByNumberJSON).Result.Withdrawals {
			vIdx, err := strconv.ParseInt(strings.TrimPrefix(wd.ValidatorIndex, "0x"), 16, 64)
			if err != nil {
				log.Fatal("can't parse hex to int")
			}
			if len(earnings[vIdx].Withdrawals) == 0 {
				earnings[vIdx].Withdrawals = append(earnings[vIdx].Withdrawals, wd)
				firstWdFill++
			}
		}
	}
	writter.Start()
	fmt.Fprint(writter, "\n[Capella][Searching for first witdrawals] Done!   \n")
	writter.Stop()

	return earnings
}

func (c *Collector) ActivitiesContractsEarnings_Capella(from string, to string) {
	earnings := c.firstWitdrawalsContractsActivity(from, to)
	fromBlock := mustAtoi(c.getBlockFromSlot(c.From))
	toBlock := fromBlock + (c.To - c.From)

	for bl := fromBlock; bl <= toBlock; bl++ {
		writter.Start()
		fmt.Fprintf(writter, "\n[Capella][Collecting witdrawals in selected range] %d out of %d blocks...   \n", bl, toBlock)
		writter.Stop()

		hexed := fmt.Sprintf("0x%x", bl)
		block := getGethBlock.Request(map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "eth_getBlockByNumber",
			"params":  []any{hexed, false},
		}, "8545")
		if !mustBeType(block) {
			continue
		}

		for _, wd := range block.(*geth.BlockByNumberJSON).Result.Withdrawals {
			vIdx, err := strconv.ParseInt(strings.TrimPrefix(wd.ValidatorIndex, "0x"), 16, 64)
			if err != nil {
				log.Fatal("can't parse hex to int")
			}
			if len(earnings[vIdx].Withdrawals) == 0 {
				continue
			}
			if earnings[vIdx].Withdrawals[0].Index != wd.Index {
				earnings[vIdx].Withdrawals = append(earnings[vIdx].Withdrawals, wd)
			}
		}

	}

	writter.Start()
	fmt.Fprintf(writter, "\n[Capella][Collecting witdrawals in selected range] Done!   \n")
	writter.Stop()

	var (
		balancesFrom []prysm.ValidatorDataJSON = make([]prysm.ValidatorDataJSON, c.Output.ProposersNum)
		balancesTo   []prysm.ValidatorDataJSON = make([]prysm.ValidatorDataJSON, c.Output.ProposersNum)
	)

	vf := getValidators.Request([]string{from, "validators"}, c.Port)
	copy(balancesFrom, vf.(*prysm.ValidatorsJSON).Data)
	vt := getValidators.Request([]string{to, "validators"}, c.Port)
	copy(balancesTo, vt.(*prysm.ValidatorsJSON).Data)

	if !mustBeType(vf) || !mustBeType(vt) {
		log.Fatal("error receiving validators info")
	}

	for ind := 0; ind < c.Output.ProposersNum; ind++ {
		writter.Start()
		fmt.Fprintf(writter, "\n[Capella][Filling output] %d out of %d validators...   \n", ind, len(c.Output.Proposers))
		writter.Stop()
		if earnings[ind].IsCredentialsSet {
			for eInd := 1; eInd < len(earnings[ind].Withdrawals); eInd++ {
				earned, err := strconv.ParseUint(strings.TrimPrefix(earnings[ind].Withdrawals[eInd].Amount, "0x"), 16, 64)
				if err != nil {
					log.Fatal("can't convert hex string to int")
				}
				c.Output.Proposers[ind].Earned += earned
			}

			c.Output.Proposers[ind].Earned += mustParseUInt64(balancesTo[ind].Balance) % uint64(8192*1e9)
			c.Output.Minted.Add(c.Output.Minted, new(big.Int).SetUint64(c.Output.Proposers[ind].Earned))
			log.Infof("Idx: %d, Contract: %s, Set: %+v", ind, balancesTo[ind].Validator.Contract, earnings[ind].IsCredentialsSet)
			continue
		}
		after := mustParseUInt64(balancesTo[ind].Balance)
		before := mustParseUInt64(balancesFrom[ind].Balance)
		c.Output.Proposers[ind].Earned += after - before
		c.Output.Minted.Add(c.Output.Minted, new(big.Int).SetUint64(c.Output.Proposers[ind].Earned))
	}

}

func (c *Collector) ActivitiesContractsEarnings_BellatrixCapella(from string, to string) {
	forkSlot := mustAtoi(c.Fork) * 32
	c.ActivitiesContractsEarnings_Bellatrix(from, fmt.Sprintf("%d", forkSlot))
	c.ActivitiesContractsEarnings_Capella(fmt.Sprintf("%d", forkSlot), to)
}

func (c *Collector) ActivitiesContractsEarnings(from string, to string) {
	var (
		ver   string
		toVer string
	)
	ver, _ = c.getFork(from)

	if toVer, c.Fork = c.getFork(to); ver != toVer {
		ver = cross
	}

	switch ver[0:4] {
	case bellatrixVersion:
		writter.Start()
		fmt.Fprint(writter, "\n[Fork in selected range defined] Bellatrix   \n")
		writter.Stop()
		c.ActivitiesContractsEarnings_Bellatrix(from, to)
		return
	case capellaVersion:
		writter.Start()
		fmt.Fprint(writter, "\n[Fork in selected range defined] Capella   \n")
		writter.Stop()
		c.ActivitiesContractsEarnings_Capella(from, to)
		return
	case cross:
		writter.Start()
		fmt.Fprint(writter, "\n[Fork in selected range defined] Bellatrix -> Capella   \n")
		writter.Stop()
		c.ActivitiesContractsEarnings_BellatrixCapella(from, to)
		return
	default:
		log.Fatal("unrecognizable fork version!")
	}

}
