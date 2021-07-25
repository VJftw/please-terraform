package logging

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger = NewLogger()

func NewLogger() zerolog.Logger {
	return zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
