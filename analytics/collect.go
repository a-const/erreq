package analytics

import (
	"encoding/json"
	"erreq/service"
	prysm "erreq/vanilla/prysm/getrequests"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"golang.org/x/exp/slices"

	log "github.com/sirupsen/logrus"
)

func NewCounter() *Collector {
	return &Collector{}
}

func (c *Collector) Log(msg string) {

}

// func (c *Collector) Logf(msg string, params...){

// }

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
		log.Fatal("can't convert string to uint64!")
	}
	return res
}

func (c *Collector) toFile(filename string) {
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

func (c *Collector) Collect(from string, to string, filename string, port string) {
	c.Port = port

	if to == "head" {
		headBlock := getBlockByID.Request([]string{to}, c.Port)
		c.To = mustAtoi(headBlock.(*prysm.BlockByIDJSON).Data.Message.Slot)
	} else {
		c.To = mustAtoi(to)
	}
	c.From = mustAtoi(from)

	c.Output = &ProposersJSON{
		MaxProposed:  -1,
		From:         c.From,
		To:           c.To,
		ProposersNum: c.getAllProposersNum(to),
		Burned:       big.NewInt(0),
		Minted:       big.NewInt(0),
	}
	c.Output.Proposers = make([]*ProposerJSON, c.Output.ProposersNum)

	for i := 0; i < c.Output.ProposersNum; i++ {
		c.Output.Proposers[i] = &ProposerJSON{
			Index: i,
		}
	}

	for i := c.From; i <= c.To; i++ {
		block := getBlockByID.Request([]string{fmt.Sprintf("%d", i)}, c.Port)
		switch t := block.(type) {
		case *prysm.BlockByIDJSON:
			ind := mustAtoi(t.Data.Message.ProposerIndex)
			c.Output.Proposers[ind].Index = ind
			c.Output.Proposers[ind].Counter += 1
			if c.Output.MaxProposed < c.Output.Proposers[ind].Counter {
				c.Output.MaxProposed = c.Output.Proposers[ind].Counter
			}
			c.Output.Burned.Add(c.Output.Burned, c.getBurned(mustParseUInt64(t.Data.Message.Body.ExecutionPayload.BlockNumber)))
			writter.Start()
			fmt.Fprintf(writter, "\n[Counting blocks] From: %d. To: %d. Current block: %d.  Proposer index %d.   \n", c.From, c.To, i, ind)
			writter.Stop()
		case *service.ErrorHandler:
			continue
		}
	}

	c.Output.Burned.Div(c.Output.Burned, big.NewInt(1e9))

	c.ActivitiesContractsEarnings(from, to)

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
