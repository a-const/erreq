package main

import (
	propcounter "brreq/prop-counter"
	"brreq/service"
	vanila "brreq/vanila"
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
	peers     = vanila.SpawnGetRequest("peers")
	peerbyID  = vanila.SpawnGetRequest("peerbyid")
	syncing   = vanila.SpawnGetRequest("syncing")
	identity  = vanila.SpawnGetRequest("identity")
	peerCount = vanila.SpawnGetRequest("peer_count")
	version   = vanila.SpawnGetRequest("version")

	//Beacon
	genesis             = vanila.SpawnGetRequest("genesis")
	validators          = vanila.SpawnGetRequest("validators")
	root                = vanila.SpawnGetRequest("root")
	fork                = vanila.SpawnGetRequest("fork")
	finalityCheckpoints = vanila.SpawnGetRequest("finality_checkpoints")
	validatorByID       = vanila.SpawnGetRequest("validator_by_id")
	blockByID           = vanila.SpawnGetRequest("block_by_id")

	ctr = propcounter.NewCounter()
)

func main() {
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

	action := func(ctx *cli.Context, obj service.Get, params ...string) error {
		writer := uilive.New()
		writer.Start()
		delay := ctx.Int64(delayFlag.Name)
		if delay > 0 {
			ticker := time.NewTicker(time.Second * time.Duration(delay))
			for {
				rsp := obj.Request(params)
				indented, err := json.MarshalIndent(rsp, " ", "   ")
				if err != nil {
					log.Errorf("can't indent json data! err: %s", err)
				}
				fmt.Fprintf(writer, "\n%s\nDelay = %ds. Working...", indented, delay)
				<-ticker.C
			}
		}

		rsp := obj.Request(params)
		indented, err := json.MarshalIndent(rsp, " ", "    ")
		if err != nil {
			log.Errorf("can't indent json data! err: %s", err)
		}
		fmt.Printf("%s", indented)
		writer.Stop()
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
				return action(ctx, peerbyID, ctx.String("id"))
			}
			return action(ctx, peers)
		},
	}

	//
	// Sync
	//

	syncCommand := &cli.Command{
		Name: "syncing",
		Action: func(ctx *cli.Context) error {
			return action(ctx, syncing)
		},
	}

	//
	// Identity
	//

	identityCommand := &cli.Command{
		Name: "identity",
		Action: func(ctx *cli.Context) error {
			return action(ctx, identity)
		},
	}

	//
	// Peer count
	//

	peerCountCommand := &cli.Command{
		Name: "peer_count",
		Action: func(ctx *cli.Context) error {
			return action(ctx, peerCount)
		},
	}

	//
	// Version
	//

	versionCommand := &cli.Command{
		Name: "peer_count",
		Action: func(ctx *cli.Context) error {
			return action(ctx, version)
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
			return action(ctx, genesis)
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
				return action(ctx, validatorByID, ctx.String("id"), "validators", ctx.String("v"))
			}
			return action(ctx, validators, ctx.String("id"), "validators")
		},
	}

	//
	// Root
	//

	rootCommand := &cli.Command{
		Name:  "root",
		Flags: stateFlags,
		Action: func(ctx *cli.Context) error {
			return action(ctx, root, ctx.String("id"), "root")
		},
	}

	//
	// Fork
	//

	forkCommand := &cli.Command{
		Name:  "fork",
		Flags: stateFlags,
		Action: func(ctx *cli.Context) error {
			return action(ctx, fork, ctx.String("id"), "fork")
		},
	}

	//
	// Finality checkpoints
	//

	finalityCheckpointsCommand := &cli.Command{
		Name:  "finality_checkpoints",
		Flags: stateFlags,
		Action: func(ctx *cli.Context) error {
			return action(ctx, finalityCheckpoints, ctx.String("id"), "finality_checkpoints")
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
			return action(ctx, blockByID, ctx.String("id"))

		},
	}

	proposerCountFlags := make([]cli.Flag, 0, 1)
	fromFlag := &cli.StringFlag{
		Name:     "from",
		Usage:    "slot from which app will count",
		Required: true,
	}
	toFlag := &cli.StringFlag{
		Name:     "to",
		Usage:    "last slot for counter",
		Required: true,
	}
	filenameFlag := &cli.StringFlag{
		Name:  "filename",
		Usage: "name for output file",
	}
	proposerCountFlags = append(proposerCountFlags, filenameFlag, fromFlag, toFlag)

	ctr := propcounter.NewCounter()

	proposerCountCommand := &cli.Command{
		Name:  "prop-count",
		Flags: proposerCountFlags,
		Action: func(ctx *cli.Context) error {
			ctr.Count(ctx.String("from"), ctx.String("to"), ctx.String("filename"))
			return nil
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
		// Counter
		proposerCountCommand,
	}
	app.Flags = appFlags

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("can't start app! err: %s", err)
	}

}
