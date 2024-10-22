package substreams

import (
	"context"
	"fmt"
	"os"

	"go.mau.fi/zerozap"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
	"github.com/streamingfast/cli"
	"github.com/streamingfast/logging"
	sink "github.com/streamingfast/substreams-sink"
	pbsubstreamsrpc "github.com/streamingfast/substreams/pb/sf/substreams/rpc/v2"

	v1 "github.com/helloacai/spindle/pb/contract/v1"
	"github.com/helloacai/spindle/pkg/db"
	"github.com/helloacai/spindle/pkg/log"
	"github.com/helloacai/spindle/pkg/manager"
)

const (
	substreamsAPIToken       = "SUBSTREAMS_API_TOKEN"
	expectedOutputModuleType = "contract.v1.EventsCalls"
)

var zlogger = zap.New(zerozap.New(log.Logger))

type Tracer struct{}

func (t *Tracer) Enabled() bool {
	return false
}

func AddFlags(cmd *cobra.Command) {
	sink.AddFlagsToSet(cmd.Flags())
}

func init() {
	if len(os.Getenv(substreamsAPIToken)) == 0 {
		// try to get token from secrets
		b, err := os.ReadFile("/etc/secrets/" + substreamsAPIToken)
		if err != nil {
			panic(err)
		}
		fmt.Println("GOT TOKEN: " + string(b))
		os.Setenv(substreamsAPIToken, string(b))
	}
}

// blocking call
func Listen(cmd *cobra.Command) {
	logging.InstantiateLoggers()

	// hard coded for now
	endpoint := "polygon.streamingfast.io:443"
	manifestPath := "substreams/spindle-v0.1.0.spkg"
	outputModuleName := "map_events_calls"

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

	sinker.Run(context.Background(), fetchCursor(), sink.NewSinkerHandlers(handleBlockScopedData, handleBlockUndoSignal))
}

func fetchCursor() *sink.Cursor {
	cursorString, err := db.GetCursor()
	if err != nil {
		log.Err(err).Msg("error getting cursor")
		return sink.NewBlankCursor()
	} else if len(cursorString) == 0 {
		return sink.NewBlankCursor()
	}

	cursor, err := sink.NewCursor(cursorString)
	if err != nil {
		log.Err(err).Msg("error recovering cursor, using blank cursor")
		return sink.NewBlankCursor()
	}
	return cursor
}

var seenBlocks = 0

func handleBlockScopedData(ctx context.Context, data *pbsubstreamsrpc.BlockScopedData, isLive *bool, cursor *sink.Cursor) error {
	seenBlocks++

	// apparently we get all blocks to this method so we need to skip ones we don't care about
	if len(data.Output.MapOutput.TypeUrl) == 0 {
		if seenBlocks > 10 { // save cursor if we've seen a bunch of blocks
			seenBlocks = 0
			return db.SaveCursor(cursor.String())
		}
		return nil
	}

	eventsCalls := &v1.EventsCalls{}
	if err := data.Output.MapOutput.UnmarshalTo(eventsCalls); err != nil {
		return fmt.Errorf("unable to unmarshal events calls : %w (typeurl: %s)", err, data.Output.MapOutput.TypeUrl)
	}

	if err := manager.Handle(log.WithContext(ctx), eventsCalls); err != nil {
		return err
	}

	return db.SaveCursor(cursor.String())
}

func handleBlockUndoSignal(ctx context.Context, undoSignal *pbsubstreamsrpc.BlockUndoSignal, cursor *sink.Cursor) error {
	// TODO: rewind if needed
	log.Info().Str("lastValidBlock", undoSignal.LastValidBlock.String()).
		Msg("NOT Rewinding changes back to block %s (unimplemented)\n")

	return db.SaveCursor(cursor.String())
}
