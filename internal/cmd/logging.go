package cmd

import (
	log "github.com/VJftw/please-terraform/internal/logging"
	"github.com/rs/zerolog"
)

// LoggingOpts represents the available logging options for command line tools.
type LoggingOpts struct {
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
}

func configureLoggingFromOpts(opts *LoggingOpts) {
	defaultLevel := zerolog.ErrorLevel

	loggingVerbosity := len(opts.Verbose)

	newLevelInt := int(defaultLevel) - loggingVerbosity

	if newLevelInt < int(zerolog.TraceLevel) {
		newLevelInt = int(zerolog.TraceLevel)
	}

	newLevel := zerolog.Level(newLevelInt)

	zerolog.SetGlobalLevel(newLevel)
	log.Logger = log.Logger.Level(newLevel)
}
