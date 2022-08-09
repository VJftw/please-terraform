package logging

import (
	"os"

	"github.com/rs/zerolog"
)

// Logger is the default Logger.
var Logger = NewLogger()

// NewLogger returns a new logger.
func NewLogger() zerolog.Logger {
	return zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
