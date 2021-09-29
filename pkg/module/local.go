package module

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type LocalCommand struct {
	*Module
	Name string `long:"name"`
	Srcs string `long:"srcs"`
}

// Execute builds a Terraform Registry Module as a Terraform Module.
func (c *LocalCommand) Execute(args []string) error {
	if err := os.MkdirAll(c.Out, os.ModePerm); err != nil {
		return err
	}
	srcs := strings.Split(c.Srcs, " ")
	for _, src := range srcs {
		// flatten
		if err := os.Rename(src, filepath.Join(c.Out, filepath.Base(src))); err != nil {
			return fmt.Errorf("could not move file: %w", err)
		}
	}

	return c.Build()
}
