package main

import (
	"context"
	"encoding/json"
	propcounter "erreq/prop-counter"
	"erreq/service"
	vanilla "erreq/vanilla"
	"fmt"
	"os"
	"os/signal"
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
	genesis              = vanilla.SpawnGetRequest("prysm_genesis")
	validators           = vanilla.SpawnGetRequest("prysm_validators")
	root                 = vanilla.SpawnGetRequest("prysm_root")
	fork                 = vanilla.SpawnGetRequest("prysm_fork")
	finalityCheckpoints  = vanilla.SpawnGetRequest("prysm_finality_checkpoints")
	validatorByID        = vanilla.SpawnGetRequest("prysm_validator_by_id")
	blockByID            = vanilla.SpawnGetRequest("prysm_block_by_id")
	validatorBalances    = vanilla.SpawnGetRequest("prysm_validator_balances")
	syncCommittees       = vanilla.SpawnGetRequest("prysm_sync_committees")
	rewardsBlocks        = vanilla.SpawnGetRequest("prysm_rewards_blocks")
	randao               = vanilla.SpawnGetRequest("prysm_randao")
	lightclientUpdates   = vanilla.SpawnGetRequest("prysm_lightclient_updates")
	lightclientBootstrap = vanilla.SpawnGetRequest("prysm_lightclient_bootstrap")
	headerByID           = vanilla.SpawnGetRequest("prysm_header_by_id")
	headers              = vanilla.SpawnGetRequest("prysm_headers")
	depositSnapshot      = vanilla.SpawnGetRequest("prysm_deposit_snapshot")
	blobSidecars         = vanilla.SpawnGetRequest("prysm_blob_sidecars")
	blindedBlocks        = vanilla.SpawnGetRequest("prysm_blinded_blocks")
	attestations         = vanilla.SpawnGetRequest("prysm_attestations")

	// Analyze
	ctr = propcounter.NewCounter()

	// Geth
	bbn = vanilla.SpawnPostRequest("geth_block_by_number")
)

type UI struct {
	writter *uilive.Writer
	wg      sync.WaitGroup
	spin    *chin.Chin
}

func (ui *UI) gracefulShutdown(ctx context.Context) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	for range sigchan {
		break
	}
	ui.spin.Stop()
	ui.wg.Wait()

	select {
	case <-ctx.Done():
		log.Infoln("Shutdown...")
		os.Exit(0)
	case <-time.After(time.Millisecond * 200):
		log.Infoln("Shutdown...")
		os.Exit(1)
	}
}

func main() {
	ctx := context.Background()
	ui := &UI{
		writter: uilive.New(),
		spin:    chin.New(),
	}

	app := cli.App{}
	app.Name = "erreq"
	app.Usage = "Little helper for your curl requests to beacon-chain rpc"
	app.Version = "0.3"
	app.Before = func(ctx *cli.Context) error {
		ui.spin = ui.spin.WithWait(&ui.wg)
		go ui.spin.Start()
		go ui.gracefulShutdown(ctx.Context)
		return nil
	}
	app.After = func(ctx *cli.Context) error {
		ui.spin.Stop()
		ui.wg.Wait()
		return nil
	}

	appFlags := make([]cli.Flag, 0, 1)
	delayFlag := &cli.Int64Flag{
		Name:  "d",
		Value: 0,
		Usage: "delay for refresh",
	}
	appFlags = append(appFlags, delayFlag)

	vanillaAction := func(ctx *cli.Context, obj service.Get, port string, params ...string) error {
		ui.writter.Start()

		delay := ctx.Int64(delayFlag.Name)
		if delay > 0 {
			ticker := time.NewTicker(time.Second * time.Duration(delay))
			for {
				rsp := obj.Request(params, port)
				indented, err := json.MarshalIndent(rsp, " ", "   ")
				if err != nil {
					log.Errorf("can't indent json data! err: %s", err)
				}
				fmt.Fprintf(ui.writter, "\n%s\nDelay = %ds. Working...", indented, delay)
				//fmt.Printf("\n%s\nDelay = %ds. Working...", indented, delay)
				<-ticker.C
			}
		}
		rsp := obj.Request(params, port)
		indented, err := json.MarshalIndent(rsp, " ", "    ")
		if err != nil {
			log.Errorf("can't indent json data! err: %s", err)
		}
		fmt.Printf("%s", indented)
		return nil
	}

	beaconPortFlag := &cli.StringFlag{
		Name:  "p",
		Value: "3500",
		Usage: "port for rpc requests",
	}
	gethPortFlag := &cli.StringFlag{
		Name:  "p",
		Value: "8545",
		Usage: "port for rpc requests",
	}

	//////////////////////////////
	//	 (vanilla) prysm/Node	//
	//////////////////////////////

	//
	// Peers
	//

	peerFlags := make([]cli.Flag, 0, 1)
	peerByIDFlag := &cli.StringFlag{
		Name:  "id",
		Usage: "get peer by its ID",
	}
	peerFlags = append(peerFlags, peerByIDFlag, beaconPortFlag)

	peerCommand := &cli.Command{
		Name:  "peers",
		Flags: peerFlags,
		Action: func(ctx *cli.Context) error {
			if len(ctx.String(peerByIDFlag.Name)) > 0 {
				return vanillaAction(ctx, peerbyID, ctx.String("p"), ctx.String(peerByIDFlag.Name))
			}
			return vanillaAction(ctx, peers, ctx.String("p"))
		},
	}

	//
	// Sync
	//

	syncCommand := &cli.Command{
		Name:  "syncing",
		Flags: []cli.Flag{beaconPortFlag},
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, syncing, ctx.String("p"))
		},
	}

	//
	// Identity
	//

	identityCommand := &cli.Command{
		Name:  "identity",
		Flags: []cli.Flag{beaconPortFlag},
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, identity, ctx.String("p"))
		},
	}

	//
	// Peer count
	//

	peerCountCommand := &cli.Command{
		Name:  "peer_count",
		Flags: []cli.Flag{beaconPortFlag},
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, peerCount, ctx.String("p"))
		},
	}

	//
	// Version
	//

	versionCommand := &cli.Command{
		Name:  "peer_count",
		Flags: []cli.Flag{beaconPortFlag},
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, version, ctx.String("p"))
		},
	}

	//////////////////////////////
	//  (vanilla) prysm/Beacon	//
	//////////////////////////////

	//
	// Genesis
	//

	genesisCommand := &cli.Command{
		Name:  "genesis",
		Flags: []cli.Flag{beaconPortFlag},
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, genesis, ctx.String("p"))
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
	stateFlags = append(stateFlags, stateIDFlag, beaconPortFlag)

	blockFlags := make([]cli.Flag, 0, 1)
	blockIDFlag := &cli.StringFlag{
		Name:     "id",
		Usage:    "block ID",
		Required: true,
	}

	blockFlags = append(blockFlags, blockIDFlag, beaconPortFlag)

	validatorsCommand := &cli.Command{
		Name:  "validators",
		Flags: append(stateFlags, validatorIDFlag),
		Action: func(ctx *cli.Context) error {
			if len(ctx.String("v")) > 0 {
				return vanillaAction(ctx, validatorByID, ctx.String(stateIDFlag.Name), ctx.String("p"), "validators", ctx.String(validatorIDFlag.Name))
			}
			return vanillaAction(ctx, validators, ctx.String("p"), ctx.String("s"), "validators")
		},
	}

	//
	// Root
	//

	rootCommand := &cli.Command{
		Name:  "root",
		Flags: stateFlags,
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, root, ctx.String("p"), ctx.String(stateIDFlag.Name), "root")
		},
	}

	//
	// Fork
	//

	forkCommand := &cli.Command{
		Name:  "fork",
		Flags: stateFlags,
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, fork, ctx.String("p"), ctx.String(stateIDFlag.Name), "fork")
		},
	}

	//
	// Finality checkpoints
	//

	finalityCheckpointsCommand := &cli.Command{
		Name:  "finality_checkpoints",
		Flags: stateFlags,
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, finalityCheckpoints, ctx.String("p"), ctx.String(stateIDFlag.Name), "finality_checkpoints")
		},
	}

	//
	// Block by id
	//

	blockCommand := &cli.Command{
		Name:  "block",
		Flags: blockFlags,
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, blockByID, ctx.String("p"), ctx.String(blockIDFlag.Name))

		},
	}

	//
	// Validator balances
	//

	validatorBalancesCommand := &cli.Command{
		Name:  "validator_balances",
		Flags: stateFlags,
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, validatorBalances, ctx.String("p"), ctx.String(stateIDFlag.Name), "validator_balances")

		},
	}

	//
	// Sync committees
	//

	syncCommitteesCommand := &cli.Command{
		Name:  "sync_committees",
		Flags: stateFlags,
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, syncCommittees, ctx.String("p"), ctx.String(stateIDFlag.Name), "sync_committees")

		},
	}

	//
	// Reward blocks
	//

	rewardsBlocksCommand := &cli.Command{
		Name:  "rewards_blocks",
		Flags: blockFlags,
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, rewardsBlocks, ctx.String("p"), ctx.String(blockIDFlag.Name))

		},
	}

	//
	// Randao
	//

	randaoCommand := &cli.Command{
		Name:  "randao",
		Flags: stateFlags,
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, randao, ctx.String("p"), ctx.String(stateIDFlag.Name), "randao")

		},
	}

	//
	// LightClient updates
	//

	lightclientUpdatesCommand := &cli.Command{
		Name: "lightclient_updates",
		Flags: []cli.Flag{
			beaconPortFlag,
		},
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, lightclientUpdates, ctx.String("p"))

		},
	}

	//
	// LightClient bootstrap
	//

	lightclientBootstrapCommand := &cli.Command{
		Name: "lightclient_bootstrap",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "r",
				Usage:    "block root",
				Required: true,
			},
			beaconPortFlag,
		},
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, lightclientBootstrap, ctx.String("p"), ctx.String("r"))

		},
	}

	//
	// Header by id
	//

	headerByIDCommand := &cli.Command{
		Name:  "header",
		Flags: blockFlags,
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, headerByID, ctx.String("p"), ctx.String(blockIDFlag.Name))

		},
	}

	//
	// Headers
	//

	headersCommand := &cli.Command{
		Name: "headers",
		Flags: []cli.Flag{
			beaconPortFlag,
		},
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, headers, ctx.String("p"))

		},
	}

	//
	// Deposit snapshot
	//

	depositSnapshotCommand := &cli.Command{
		Name: "deposit_snapshot",
		Flags: []cli.Flag{
			beaconPortFlag,
		},
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, depositSnapshot, ctx.String("p"))

		},
	}

	//
	// Blob sidecars
	//

	blobSidecarsCommand := &cli.Command{
		Name:  "blob_sidecars",
		Flags: blockFlags,
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, blobSidecars, ctx.String("p"), ctx.String(blockIDFlag.Name))

		},
	}

	//
	// Blinded blocks
	//

	blindedBlocksCommand := &cli.Command{
		Name:  "blinded_blocks",
		Flags: blockFlags,
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, blindedBlocks, ctx.String("p"), ctx.String(blockIDFlag.Name))

		},
	}

	//
	// Attestations
	//

	attestationsCommand := &cli.Command{
		Name:  "attestations",
		Flags: blockFlags,
		Action: func(ctx *cli.Context) error {
			return vanillaAction(ctx, attestations, ctx.String("p"), ctx.String(blockIDFlag.Name), "attestations")

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
	proposerCountFlags = append(proposerCountFlags, filenameFlag, fromFlag, toFlag, beaconPortFlag)

	proposerCountCommand := &cli.Command{
		Name:        "prop-count",
		Description: "Retrieves proposed blocks and other info for every validator",
		Flags:       proposerCountFlags,
		Action: func(ctx *cli.Context) error {
			ctr.Count(ctx.String(fromFlag.Name), ctx.String(toFlag.Name), ctx.String(filenameFlag.Name), ctx.String("p"))
			return nil
		},
	}

	//////////////////////////////
	//		(vanilla) Geth	    //
	//////////////////////////////

	//
	// Block by number
	//

	blockByNumberFlags := make([]cli.Flag, 0, 1)
	numberFlag := &cli.Int64Flag{
		Name:  "n",
		Usage: "block to retreive",
	}
	blockByNumberFlags = append(blockByNumberFlags, numberFlag, gethPortFlag)

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
			}, ctx.String("p"))
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
		validatorBalancesCommand,
		syncCommitteesCommand,
		rewardsBlocksCommand,
		randaoCommand,
		lightclientUpdatesCommand,
		lightclientBootstrapCommand,
		headerByIDCommand,
		headersCommand,
		depositSnapshotCommand,
		blobSidecarsCommand,
		blindedBlocksCommand,
		attestationsCommand,
		// Analyse
		proposerCountCommand,
		// Geth
		bbnCommand,
	}
	app.Flags = appFlags

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Fatalf("can't start app! err: %s", err)
	}

}
