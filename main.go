package main

import (
	factory "brreq/factory"
	"brreq/service"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/adhocore/chin"
	"github.com/gosuri/uilive"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	//Node
	peers     = factory.SpawnGetRequest("peers")
	peerbyID  = factory.SpawnGetRequest("peerbyid")
	syncing   = factory.SpawnGetRequest("syncing")
	identity  = factory.SpawnGetRequest("identity")
	peerCount = factory.SpawnGetRequest("peer_count")
	version   = factory.SpawnGetRequest("version")

	//Beacon
	genesis             = factory.SpawnGetRequest("genesis")
	validators          = factory.SpawnGetRequest("validators")
	root                = factory.SpawnGetRequest("root")
	fork                = factory.SpawnGetRequest("fork")
	finalityCheckpoints = factory.SpawnGetRequest("finality_checkpoints")
	validatorByID       = factory.SpawnGetRequest("validator_by_id")
	blockByID           = factory.SpawnGetRequest("block_by_id")
)

func main() {
	// Screen refresher
	writer := uilive.New()
	writer.Start()

	// Spinner
	var wg sync.WaitGroup
	s := chin.New().WithWait(&wg)
	go s.Start()

	app := cli.App{}
	app.Name = "brreq"
	app.Usage = "Little helper for your curl requests to beacon-chain rpc"
	app.Version = "0.1"

	appFlags := make([]cli.Flag, 0, 1)
	delayFlag := &cli.Int64Flag{
		Name:  "d",
		Value: 0,
		Usage: "delay for refresh",
	}
	appFlags = append(appFlags, delayFlag)

	action := func(ctx *cli.Context, wr *uilive.Writer, obj service.Get, params ...string) error {

		delay := ctx.Int64(delayFlag.Name)
		if delay > 0 {
			ticker := time.NewTicker(time.Second * time.Duration(delay))
			for {
				rsp := obj.Request(params)
				indented, err := json.MarshalIndent(rsp, " ", "   ")
				if err != nil {
					log.Errorf("can't indent json data! err: %s", err)
				}
				fmt.Fprintf(wr, "\n%s\nDelay = %ds. Working...", indented, delay)
				<-ticker.C
			}
		}

		rsp := obj.Request(params)
		indented, err := json.MarshalIndent(rsp, " ", "    ")
		if err != nil {
			log.Errorf("can't indent json data! err: %s", err)
		}
		fmt.Printf("%s", indented)
		return nil
	}

	//////////////////////////////
	//			Node			//
	//////////////////////////////

	//
	// Peers
	//

	peerFlags := make([]cli.Flag, 0, 1)
	peerByIDFlag := &cli.StringFlag{
		Name:  "id",
		Usage: "get peer by its ID",
	}
	peerFlags = append(peerFlags, peerByIDFlag)

	peerCommand := &cli.Command{
		Name:  "peers",
		Flags: peerFlags,
		Action: func(ctx *cli.Context) error {
			if len(ctx.String("id")) > 0 {
				return action(ctx, writer, peerbyID, ctx.String("id"))
			}
			return action(ctx, writer, peers)
		},
	}

	//
	// Sync
	//

	syncCommand := &cli.Command{
		Name: "syncing",
		Action: func(ctx *cli.Context) error {
			return action(ctx, writer, syncing)
		},
	}

	//
	// Identity
	//

	identityCommand := &cli.Command{
		Name: "identity",
		Action: func(ctx *cli.Context) error {
			return action(ctx, writer, identity)
		},
	}

	//
	// Peer count
	//

	peerCountCommand := &cli.Command{
		Name: "peer_count",
		Action: func(ctx *cli.Context) error {
			return action(ctx, writer, peerCount)
		},
	}

	//
	// Version
	//

	versionCommand := &cli.Command{
		Name: "peer_count",
		Action: func(ctx *cli.Context) error {
			return action(ctx, writer, version)
		},
	}

	//////////////////////////////
	//			Beacon			//
	//////////////////////////////

	//
	// Genesis
	//

	genesisCommand := &cli.Command{
		Name: "genesis",
		Action: func(ctx *cli.Context) error {
			return action(ctx, writer, genesis)
		},
	}

	//
	// Validators
	//

	stateFlags := make([]cli.Flag, 0, 1)
	stateIDFlag := &cli.StringFlag{
		Name:     "id",
		Usage:    "state ID",
		Required: true,
	}
	validatorIDFlag := &cli.StringFlag{
		Name:  "v",
		Usage: "validator ID",
	}
	stateFlags = append(stateFlags, stateIDFlag)

	validatorsCommand := &cli.Command{
		Name:  "validators",
		Flags: append(stateFlags, validatorIDFlag),
		Action: func(ctx *cli.Context) error {
			if len(ctx.String("v")) > 0 {
				return action(ctx, writer, validatorByID, ctx.String("id"), "validators", ctx.String("v"))
			}
			return action(ctx, writer, validators, ctx.String("id"), "validators")
		},
	}

	//
	// Root
	//

	rootCommand := &cli.Command{
		Name:  "root",
		Flags: stateFlags,
		Action: func(ctx *cli.Context) error {
			return action(ctx, writer, root, ctx.String("id"), "root")
		},
	}

	//
	// Fork
	//

	forkCommand := &cli.Command{
		Name:  "fork",
		Flags: stateFlags,
		Action: func(ctx *cli.Context) error {
			return action(ctx, writer, fork, ctx.String("id"), "fork")
		},
	}

	//
	// Finality checkpoints
	//

	finalityCheckpointsCommand := &cli.Command{
		Name:  "finality_checkpoints",
		Flags: stateFlags,
		Action: func(ctx *cli.Context) error {
			return action(ctx, writer, finalityCheckpoints, ctx.String("id"), "finality_checkpoints")
		},
	}

	//
	// Block by id
	//

	blockFlags := make([]cli.Flag, 0, 1)
	blockIDFlag := &cli.StringFlag{
		Name:     "id",
		Usage:    "block ID",
		Required: true,
	}

	blockFlags = append(blockFlags, blockIDFlag)

	blockCommand := &cli.Command{
		Name:  "block",
		Flags: blockFlags,
		Action: func(ctx *cli.Context) error {
			return action(ctx, writer, blockByID, ctx.String("id"))

		},
	}

	app.Commands = []*cli.Command{
		//Node
		peerCommand,
		syncCommand,
		identityCommand,
		peerCountCommand,
		versionCommand,
		//Beacon
		genesisCommand,
		validatorsCommand,
		rootCommand,
		forkCommand,
		finalityCheckpointsCommand,
		blockCommand,
	}
	app.Flags = appFlags

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("can't start app! err: %s", err)
	}

}
