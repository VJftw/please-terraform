package main

import (
	"github.com/VJftw/please-terraform/internal/cmd"
	"github.com/VJftw/please-terraform/pkg/module"
	"github.com/VJftw/please-terraform/pkg/root"
)

type opts struct {
	Module *module.Command `command:"module"`
	Root   *root.Command   `command:"root"`
}

func main() {
	opts := &opts{}
	cmd.MustParseFlags(opts)
}
