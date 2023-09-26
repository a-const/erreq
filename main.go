package main

import (
	propcounter "brreq/prop-counter"
	"brreq/service"
	vanilla "brreq/vanilla"
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
	peers     = vanilla.SpawnGetRequest("prysm_peers")
	peerbyID  = vanilla.SpawnGetRequest("prysm_peerbyid")
	syncing   = vanilla.SpawnGetRequest("prysm_syncing")
	identity  = vanilla.SpawnGetRequest("prysm_identity")
	peerCount = vanilla.SpawnGetRequest("prysm_peer_count")
	version   = vanilla.SpawnGetRequest("prysm_version")

	//Beacon
	genesis             = vanilla.SpawnGetRequest("prysm_genesis")
	validators          = vanilla.SpawnGetRequest("prysm_validators")
	root                = vanilla.SpawnGetRequest("prysm_root")
	fork                = vanilla.SpawnGetRequest("prysm_fork")
	finalityCheckpoints = vanilla.SpawnGetRequest("prysm_finality_checkpoints")
	validatorByID       = vanilla.SpawnGetRequest("prysm_validator_by_id")
	blockByID           = vanilla.SpawnGetRequest("prysm_block_by_id")

	// Analyze
	ctr = propcounter.NewCounter()

	// Geth
	bbn = vanilla.SpawnPostRequest("geth_block_by_number")
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

	vanillaAction := func(ctx *cli.Context, obj service.Get, params ...string) error {
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
	//	     (vanilla) Node		//
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
			if len(ctx.String(peerByIDFlag.Name)) > 0 {
				return vanillaAction(ctx, peerbyID, ctx.String(peerByIDFlag.Name))
			}
			return vanillaAction(ctx, peers)
		},
	}

	//
	// Sync
	//

	syncCommand := &cli.Command{
		Name: "syncing",
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, syncing)
		},
	}

	//
	// Identity
	//

	identityCommand := &cli.Command{
		Name: "identity",
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, identity)
		},
	}

	//
	// Peer count
	//

	peerCountCommand := &cli.Command{
		Name: "peer_count",
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, peerCount)
		},
	}

	//
	// Version
	//

	versionCommand := &cli.Command{
		Name: "peer_count",
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, version)
		},
	}

	//////////////////////////////
	//	   (vanilla) Beacon		//
	//////////////////////////////

	//
	// Genesis
	//

	genesisCommand := &cli.Command{
		Name: "genesis",
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, genesis)
		},
	}

	//
	// Validators
	//

	stateFlags := make([]cli.Flag, 0, 1)
	stateIDFlag := &cli.StringFlag{
		Name:     "s",
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
				return vanillaAction(ctx, validatorByID, ctx.String(stateIDFlag.Name), "validators", ctx.String(validatorIDFlag.Name))
			}
			return vanillaAction(ctx, validators, ctx.String("s"), "validators")
		},
	}

	//
	// Root
	//

	rootCommand := &cli.Command{
		Name:  "root",
		Flags: stateFlags,
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, root, ctx.String(stateIDFlag.Name), "root")
		},
	}

	//
	// Fork
	//

	forkCommand := &cli.Command{
		Name:  "fork",
		Flags: stateFlags,
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, fork, ctx.String(stateIDFlag.Name), "fork")
		},
	}

	//
	// Finality checkpoints
	//

	finalityCheckpointsCommand := &cli.Command{
		Name:  "finality_checkpoints",
		Flags: stateFlags,
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, finalityCheckpoints, ctx.String(stateIDFlag.Name), "finality_checkpoints")
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
			return vanillaAction(ctx, blockByID, ctx.String(blockIDFlag.Name))

		},
	}

	//////////////////////////////
	//	 		Analyze		    //
	//////////////////////////////

	//
	// Prop count
	//

	proposerCountFlags := make([]cli.Flag, 0, 1)
	fromFlag := &cli.StringFlag{
		Name:     "f",
		Usage:    "slot from which app will count",
		Required: true,
	}
	toFlag := &cli.StringFlag{
		Name:     "t",
		Usage:    "last slot for counter",
		Required: true,
	}
	filenameFlag := &cli.StringFlag{
		Name:  "filename",
		Usage: "name for output file",
	}
	proposerCountFlags = append(proposerCountFlags, filenameFlag, fromFlag, toFlag)

	proposerCountCommand := &cli.Command{
		Name:        "prop-count",
		Description: "Retrieves proposed blocks and other info for every validator",
		Flags:       proposerCountFlags,
		Action: func(ctx *cli.Context) error {
			ctr.Count(ctx.String(fromFlag.Name), ctx.String(toFlag.Name), ctx.String(filenameFlag.Name))
			return nil
		},
	}

	//////////////////////////////
	//	 		Geth		    //
	//////////////////////////////

	//
	// Block by number
	//

	blockByNumberFlags := make([]cli.Flag, 0, 1)
	numberFlag := &cli.Int64Flag{
		Name:  "n",
		Usage: "block to retreive",
	}
	blockByNumberFlags = append(blockByNumberFlags, numberFlag)

	bbnCommand := &cli.Command{
		Name:        "geth-block",
		Description: "Geth geth block by its number",
		Flags:       blockByNumberFlags,
		Action: func(ctx *cli.Context) error {
			hexed := fmt.Sprintf("0x%x", ctx.Int64(numberFlag.Name))
			rsp := bbn.Request(map[string]any{
				"jsonrpc": "2.0",
				"id":      1,
				"method":  "eth_getBlockByNumber",
				"params":  []any{hexed, false},
			})
			indented, err := json.MarshalIndent(rsp, " ", "    ")
			if err != nil {
				log.Errorf("can't indent json data! err: %s", err)
			}
			fmt.Printf("%s", indented)
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
		// Analyse
		proposerCountCommand,
		// Geth
		bbnCommand,
	}
	app.Flags = appFlags

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("can't start app! err: %s", err)
	}

}
