package log

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

var Logger = zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stdout})

func Info() *zerolog.Event {
	return Logger.Info()
}

func Debug() *zerolog.Event {
	return Logger.Info()
}

func WithContext(ctx context.Context) context.Context {
	return Logger.WithContext(ctx)
}
