package propcounter

import (
	"brreq/service"
	"brreq/vanila"
	"brreq/vanila/getrequests"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type ProposersJSON struct {
	MaxProposed int             `json:"max_proposed"`
	Proposers   []*ProposerJSON `json:"proposers"`
}

type ProposerJSON struct {
	Index   int `json:"index"`
	Counter int `json:"counter"`
}

type Counter struct {
	ActivationEpoch int
	Ctr             map[int]int
}

var (
	getValidator  = vanila.SpawnGetRequest("validator_by_id")
	getValidators = vanila.SpawnGetRequest("validators")
	getBlockByID  = vanila.SpawnGetRequest("block_by_id")
)

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) prepare(stateID string, proposerIndex int64) {
	rsp := getValidators.Request([]string{stateID, "validators"})
	c.Ctr = make(map[int]int, len(rsp.(*getrequests.ValidatorsJSON).Data))
	rsp = getValidator.Request([]string{stateID, "validators", fmt.Sprintf("%d", proposerIndex)})
	ae, err := strconv.Atoi(rsp.(*getrequests.ValidatorByIDJSON).Data.Validator.ActivationEpoch)
	if err != nil {
		log.Fatalf("can't convert activationEpoch to int! err: %s", err)
	}
	c.ActivationEpoch = ae
}

func (c *Counter) Count(stateID string, proposerIndex int64, filename string, from int64, to string) {
	c.prepare(stateID, proposerIndex)
	headBlock := getBlockByID.Request([]string{stateID})
	headIndex, err := strconv.Atoi(headBlock.(*getrequests.BlockByIDJSON).Data.Message.Slot)
	if err != nil {
		log.Fatalf("can't convert validatorindex to int! err: %s", err)
	}
	log.Infof("Head index = %d Activation epoch = %d", headIndex, c.ActivationEpoch)

	for i := c.ActivationEpoch * 32; i <= headIndex; i++ {
		block := getBlockByID.Request([]string{fmt.Sprintf("%d", i)})
		switch t := block.(type) {
		case *getrequests.BlockByIDJSON:
			ind, err := strconv.Atoi(t.Data.Message.ProposerIndex)
			if err != nil {
				log.Fatal("Can't conver proposer index to int!")
			}
			log.Infof("Block: %d.  Proposer index %d.", i, ind)
			c.Ctr[ind] += 1
		case *service.ErrorHandler:
			continue
		}
	}

	output := &ProposersJSON{
		MaxProposed: -1,
		Proposers:   make([]*ProposerJSON, 0, len(c.Ctr)),
	}
	for i, value := range c.Ctr {
		if output.MaxProposed < value {
			output.MaxProposed = value
		}
		p := &ProposerJSON{
			Counter: value,
			Index:   i,
		}
		output.Proposers = append(output.Proposers, p)
	}
	jsonByte, err := json.Marshal(output)
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

	fmt.Printf("Done. Results in file: %s", fn)
}
