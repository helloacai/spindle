package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/streamingfast/cli"
	. "github.com/streamingfast/cli"
	"github.com/streamingfast/logging"
	sink "github.com/streamingfast/substreams-sink"
	pbsubstreamsrpc "github.com/streamingfast/substreams/pb/sf/substreams/rpc/v2"
	"go.mau.fi/zerozap"
	"go.uber.org/zap"

	v1 "github.com/helloacai/spindle/pb/contract/v1"
	"github.com/helloacai/spindle/pkg/log"
	"github.com/helloacai/spindle/pkg/manager"
	"github.com/helloacai/spindle/pkg/server"
)

var expectedOutputModuleType = "contract.v1.EventsCalls"

var zlogger = zap.New(zerozap.New(log.Logger))

type Tracer struct{}

func (t *Tracer) Enabled() bool {
	return false
}

func main() {
	logging.InstantiateLoggers()

	Run(
		"sinker",
		"acs thread manager sink",

		Command(sinkRunE,
			"sink <endpoint> <manifest> [<output_module>]",
			"Run the sinker code",
			RangeArgs(2, 3),
			Flags(func(flags *pflag.FlagSet) {
				sink.AddFlagsToSet(flags)
			}),
		),

		OnCommandErrorLogAndExit(zlogger),
	)
}

func sinkRunE(cmd *cobra.Command, args []string) error {
	endpoint := args[0]
	manifestPath := args[1]

	outputModuleName := sink.InferOutputModuleFromPackage
	if len(args) == 3 {
		outputModuleName = args[2]
	}

	sinker, err := sink.NewFromViper(
		cmd,
		expectedOutputModuleType,
		endpoint,
		manifestPath,
		outputModuleName,
		":", // This is the block range, in our case defined as Substreams module's start block and up forever
		zlogger,
		&Tracer{},
	)
	cli.NoError(err, "unable to create sinker: %s", err)

	sinker.OnTerminating(func(err error) {
		cli.NoError(err, "unexpected sinker error")

		log.Info().Msg("sink is terminating")
	})

	cursor := sink.NewBlankCursor() // TODO: save state and restart from cursor?

	// TODO: this is terrible.
	go func() {
		if err := server.Run(); err != nil {
			panic(err)
		}
	}()

	sinker.Run(context.Background(), cursor, sink.NewSinkerHandlers(handleBlockScopedData, handleBlockUndoSignal))
	return nil
}

func handleBlockScopedData(ctx context.Context, data *pbsubstreamsrpc.BlockScopedData, isLive *bool, cursor *sink.Cursor) error {
	// apparently we get all blocks to this method so we need to skip ones we don't care about
	if len(data.Output.MapOutput.TypeUrl) == 0 {
		return nil
	}

	eventsCalls := &v1.EventsCalls{}
	if err := data.Output.MapOutput.UnmarshalTo(eventsCalls); err != nil {
		return fmt.Errorf("unable to unmarshal events calls : %w (typeurl: %s)", err, data.Output.MapOutput.TypeUrl)
	}

	if err := manager.Handle(log.WithContext(ctx), eventsCalls); err != nil {
		return err
	}

	// TODO: save cursor to file if we want to restart from it on crash
	return nil
}

func handleBlockUndoSignal(ctx context.Context, undoSignal *pbsubstreamsrpc.BlockUndoSignal, cursor *sink.Cursor) error {
	// TODO: rewind if needed
	log.Info().Str("lastValidBlock", undoSignal.LastValidBlock.String()).
		Msg("NOT Rewinding changes back to block %s (unimplemented)\n")
	return nil
}
