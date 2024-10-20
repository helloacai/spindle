package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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
	"github.com/helloacai/spindle/pkg/db"
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

const (
	substreamsAPIToken = "SUBSTREAMS_API_TOKEN"
)

func setenv() {
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

func main() {
	setenv()

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

	var cursor *sink.Cursor
	cursorString, err := db.GetCursor()
	if err != nil {
		log.Err(err).Msg("error getting cursor")
		cursor = sink.NewBlankCursor()
	} else if len(cursorString) == 0 {
		cursor = sink.NewBlankCursor()
	} else {
		cursor, err = sink.NewCursor(cursorString)
		if err != nil {
			log.Err(err).Msg("error recovering cursor, using blank cursor")
			cursor = sink.NewBlankCursor()
		}
	}

	// TODO: this is terrible.
	go func() {
		if err := server.Run(); err != nil {
			panic(err)
		}
	}()

	// TODO: this is also terrible. workaround to make render not shut down our instances
	go func() {
		for {
			resp, err := http.Get("https://yelpagent.onrender.com")
			if err != nil {
				log.Err(err).Msg("YelpAgent keepalive error")
			} else {
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Err(err).Msg("YelpAgent keepalive body error")
				} else {
					log.Debug().Msg("YelpAgent keepalive: " + string(b))
				}
			}
			resp.Body.Close()

			resp, err = http.Get("https://coordinatingagent.onrender.com")
			if err != nil {
				log.Err(err).Msg("CoordinatingAgent keepalive error")
			} else {
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Err(err).Msg("CoordinatingAgent keepalive body error")
				} else {
					log.Debug().Msg("CoordinatingAgent keepalive: " + string(b))
				}
			}
			resp.Body.Close()

			resp, err = http.Get("https://gcalagent.onrender.com")
			if err != nil {
				log.Err(err).Msg("GCalAgent keepalive error")
			} else {
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Err(err).Msg("GCalAgent keepalive body error")
				} else {
					log.Debug().Msg("GCalAgent keepalive: " + string(b))
				}
			}
			resp.Body.Close()

			resp, err = http.Get("https://spindle.onrender.com/healthz")
			if err != nil {
				log.Err(err).Msg("spindle keepalive error")
			} else {
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Err(err).Msg("spindle keepalive body error")
				} else {
					log.Debug().Msg("spindle keepalive: " + string(b))
				}
			}
			resp.Body.Close()

			time.Sleep(10 * time.Second)
		}
	}()

	sinker.Run(context.Background(), cursor, sink.NewSinkerHandlers(handleBlockScopedData, handleBlockUndoSignal))
	return nil
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
