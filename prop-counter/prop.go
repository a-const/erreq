package propcounter

import (
	"brreq/service"
	"brreq/vanilla/getrequests"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/exp/slices"

	log "github.com/sirupsen/logrus"
)

func NewCounter() *Counter {
	return &Counter{}
}

func mustAtoi(val string) int {
	res, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal("can't convert string to int!")
	}
	return res
}
func mustParseUInt64(val string) uint64 {
	res, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		log.Fatal("can't convert string to int!")
	}
	return res
}

// func (c *Counter) log(){

// }

func (c *Counter) getAllProposersNum() int {
	resp := getValidators.Request([]string{"head", "validators"})
	return len(resp.(*getrequests.ValidatorsJSON).Data)
}

func (c *Counter) getActivitiesContractsEarnings() {
	resp := getValidators.Request([]string{"head", "validators"})
	for i := 0; i < len(resp.(*getrequests.ValidatorsJSON).Data); i++ {
		if resp.(*getrequests.ValidatorsJSON).Data[i].Validator.Contract != emptyContract {
			c.Output.Proposers[i].Contract = resp.(*getrequests.ValidatorsJSON).Data[i].Validator.Contract
		}
		actv := mustParseUInt64(resp.(*getrequests.ValidatorsJSON).Data[i].Validator.EffectiveActivity)

		c.Output.Proposers[i].Activity = actv
		balance := mustParseUInt64(resp.(*getrequests.ValidatorsJSON).Data[i].Balance)

		c.Output.Proposers[i].Earned = balance % uint64(8192*1e9)
		// go func() {
		// 	ticker := time.NewTicker(time.Second * time.Duration(1))
		// 	for {
		// 		fmt.Fprintf(writter, "\n[Adding other info] Contracts, activities and earnings...%d out of %d validator", i, len(c.Output.Proposers))
		// 		<-ticker.C
		// 	}
		// }()

		//log.Infof("Writing contracts, activities and earnings...%d / %d", i, len(c.Output.Proposers))
		writter.Start()
		fmt.Fprintf(writter, "\n[Adding other info] Contracts, activities and earnings...%d out of %d validator   ", i, len(c.Output.Proposers))
		writter.Stop()
	}
}

func (c *Counter) Count(from string, to string, filename string) {
	var (
		toIndex   int
		fromIndex int
	)
	if to == "head" {
		headBlock := getBlockByID.Request([]string{to})
		toIndex = mustAtoi(headBlock.(*getrequests.BlockByIDJSON).Data.Message.Slot)
	} else {
		toIndex = mustAtoi(to)
	}
	fromIndex = mustAtoi(from)

	c.Output = &ProposersJSON{
		MaxProposed: -1,
		From:        fromIndex,
		To:          toIndex,
		Proposers:   make([]*ProposerJSON, c.getAllProposersNum()),
	}

	for i := 0; i < len(c.Output.Proposers); i++ {
		c.Output.Proposers[i] = &ProposerJSON{
			Index: i,
		}
	}

	//ticker := time.NewTicker(time.Second * time.Duration(1))
	for i := fromIndex; i <= toIndex; i++ {
		block := getBlockByID.Request([]string{fmt.Sprintf("%d", i)})
		switch t := block.(type) {
		case *getrequests.BlockByIDJSON:
			ind := mustAtoi(t.Data.Message.ProposerIndex)
			c.Output.Proposers[ind].Index = ind
			c.Output.Proposers[ind].Counter += 1
			if c.Output.MaxProposed < c.Output.Proposers[ind].Counter {
				c.Output.MaxProposed = c.Output.Proposers[ind].Counter
			}

			//if (<-ticker.C) {
			writter.Start()
			fmt.Fprintf(writter, "\n[Counting blocks] From: %d. To: %d. Current block: %d.  Proposer index %d.   ", fromIndex, toIndex, i, ind)
			writter.Stop()
			//}
			//log.Infof("From: %d. To: %d. Current block: %d.  Proposer index %d.", fromIndex, toIndex, i, ind)
		case *service.ErrorHandler:
			continue
		}
	}

	c.getActivitiesContractsEarnings()

	slices.SortStableFunc(c.Output.Proposers, func(a, b *ProposerJSON) int {
		if a.Counter > b.Counter {
			return -1
		}
		if a.Counter < b.Counter {
			return 1
		}
		return 0
	})
	c.toFile(filename)
}

func (c *Counter) toFile(filename string) {
	jsonByte, err := json.MarshalIndent(c.Output, " ", "   ")
	if err != nil {
		log.Fatalf("can't marshal results to json! %s", err)
	}

	fn := "proposers.json"
	if len(filename) > 0 {
		fn = filename
	}

	file, err := os.Create(fn)
	if err != nil {
		log.Fatalf("File is not created! err: %s", err)
	}
	if _, err = file.Write(jsonByte); err != nil {
		log.Fatalf("Data hasn't been written! err: %s", err)
	}

	fmt.Printf("\nDone. Results in file: %s\n", fn)
}
