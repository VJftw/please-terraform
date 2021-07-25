package cmd

import (
	"os"
	"path"

	"github.com/VJftw/please-terraform/internal/logging"
	"github.com/jessevdk/go-flags"
)

// MustParseFlags parses the given application options from command line arguments.
func MustParseFlags(opts interface{}) {
	appName := path.Base(os.Args[0])
	flagParser := flags.NewNamedParser(appName, flags.Default)

	flagParser.EnvNamespaceDelimiter = "_"
	flagParser.NamespaceDelimiter = "_"

	loggingOpts := &LoggingOpts{}
	flagParser.AddGroup("logging options", "logging options", loggingOpts)

	flagParser.CommandHandler = func(cmd flags.Commander, args []string) error {
		configureLoggingFromOpts(loggingOpts)
		return cmd.Execute(args)
	}

	flagParser.AddGroup(appName+" options", "", opts)

	args, err := flagParser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			handleFlagsErr(flagsErr)
		}
		logging.Logger.Fatal().Err(err).Msg("encountered error")
	}

	if len(args) > 0 {
		logging.Logger.Fatal().Strs("extra-args", args).Msg("found unexpected extra arguments")
	}

}

func handleFlagsErr(err *flags.Error) {
	if err.Type == flags.ErrHelp {
		os.Exit(0)
	}
}
