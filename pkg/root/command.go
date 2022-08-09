package root

import "github.com/VJftw/please-terraform/internal/logging"

var log = logging.NewLogger()

// Command represents the root subcommand.
type Command struct {
	Build      *CommandBuild      `command:"build"`
	VirtualEnv *CommandVirtualEnv `command:"virtualenv"`
}
