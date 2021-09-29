package main

import (
	"github.com/VJftw/please-terraform/internal/cmd"
	"github.com/VJftw/please-terraform/pkg/module"
	"github.com/VJftw/please-terraform/pkg/provider"
)

type Opts struct {
	Module   *module.Command   `command:"module"`
	Provider *provider.Command `command:"provider"`
}

func main() {
	opts := &Opts{}
	cmd.MustParseFlags(opts)
}
