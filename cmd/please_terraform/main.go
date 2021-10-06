package main

import (
	"github.com/VJftw/please-terraform/internal/cmd"
	"github.com/VJftw/please-terraform/pkg/module"
	"github.com/VJftw/please-terraform/pkg/provider"
	"github.com/VJftw/please-terraform/pkg/workspace"
)

type Opts struct {
	Module    *module.Command    `command:"module"`
	Provider  *provider.Command  `command:"provider"`
	Workspace *workspace.Command `command:"workspace"`
}

func main() {
	opts := &Opts{}
	cmd.MustParseFlags(opts)
}
