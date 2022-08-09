package module

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/VJftw/please-terraform/pkg/please"
)

// CommandLocal represents the `module local` command and its flags.
type CommandLocal struct {
	Name string `long:"name" required:"true" description:"The Please name of the Terraform module."`
	Pkg  string `long:"pkg" required:"true" description:"The Please package of the Terraform module."`
	Srcs string `long:"srcs" requrired:"true" description:"Space separated src files that make up the Terraform module. These will be flattened."`
	Out  string `long:"out" required:"true" description:"The directory to write the processed Terraform module to."`

	Aliases []string `long:"aliases" required:"false" description:"The aliases for the Terraform module that will be replaced in other Terraform configuration."`
	Strip   []string `long:"strip" required:"false" description:"The directories to strip from the Terraform module."`
	Deps    []string `long:"deps" required:"false" description:"Other Terraform modules that this Terraform module depends on."`

	Opts *Opts
}

// Execute builds a Terraform Registry Module as a Terraform Module.
func (c *CommandLocal) Execute(args []string) error {
	// set empty defaults
	if c.Aliases == nil {
		c.Aliases = []string{}
	}
	if c.Strip == nil {
		c.Strip = []string{}
	}
	if c.Deps == nil {
		c.Deps = []string{}
	}

	log.Info().
		Str("name", c.Name).
		Str("pkg", c.Pkg).
		Str("srcs", c.Srcs).
		Str("out", c.Out).
		Strs("aliases", c.Aliases).
		Strs("strip", c.Strip).
		Strs("deps", c.Deps).
		Msg("building module")

	log.Debug().Str("path", c.Out).Msg("ensuring out directory")
	if err := os.MkdirAll(c.Out, os.ModePerm); err != nil {
		return fmt.Errorf("could not ensure out directory '%s': %w", c.Out, err)
	}

	m := &Metadata{
		Target:  fmt.Sprintf("//%s:%s", c.Pkg, c.Name),
		Aliases: c.Aliases,
	}

	log.Debug().Strs("strip", c.Strip).Msg("stripping directories")
	// Strip directories
	if err := m.StripDirs(c.Out, c.Strip); err != nil {
		return fmt.Errorf("could not strip directories: %w", err)
	}

	log.Debug().Msg("adding aliases")
	m.Aliases = append(m.Aliases, []string{
		// Supports referencing by Please target.
		fmt.Sprintf("//%s:%s", c.Pkg, c.Name),
	}...)

	if filepath.Base(c.Pkg) == c.Name {
		// Supports referencing by default Please target for pkg.
		m.Aliases = append(m.Aliases, fmt.Sprintf("//%s", c.Pkg))
	}

	srcs := strings.Split(c.Srcs, " ")
	if srcs[0] != "" {
		log.Debug().Msg("copying files")
		for _, src := range srcs {
			// flatten
			if err := please.CopyFile(src, filepath.Join(c.Out, filepath.Base(src))); err != nil {
				return fmt.Errorf("could not copy file: %w", err)
			}
		}
	}

	log.Debug().Strs("deps", c.Deps).Msg("colocating modules")
	// colocate modules
	if err := ColocateModules(c.Opts.MetadataFile, c.Out, c.Deps); err != nil {
		return err
	}

	log.Debug().Str("path", c.Opts.MetadataFile).Msg("saving metadata")

	if err := m.Save(filepath.Join(c.Out, c.Opts.MetadataFile)); err != nil {
		return err
	}

	return nil
}
