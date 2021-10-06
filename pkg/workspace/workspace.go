package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/VJftw/please-terraform/internal/logging"

	"github.com/VJftw/please-terraform/pkg/module"
	"github.com/VJftw/please-terraform/pkg/please"
)

var log = logging.NewLogger()

type Command struct {
	*Workspace
}

type Workspace struct {
	Pkg    string `long:"pkg"`
	Name   string `long:"name"`
	OS     string `long:"os" description:""`
	Arch   string `long:"arch" description:""`
	PkgDir string `long:"pkg_dir"`

	SrcsCSV     string   `long:"srcs"`
	VarFilesCSV string   `long:"var_files"`
	Deps        []string `long:"deps"`

	Out          string `long:"out" description:""`
	AbsolutePath string
}

func (c *Command) Execute(args []string) error {

	// make the out directory
	os.MkdirAll(c.Out, os.ModePerm)

	// shift srcs into outs, flattening the root module
	srcs := strings.Split(c.SrcsCSV, ",")
	for _, src := range srcs {
		dest := filepath.Join(c.Out, filepath.Base(src))
		log.Debug().Str("src", src).Str("dest", dest).Msg("moving file")
		if err := os.Rename(src, dest); err != nil {
			return fmt.Errorf("could not move file: %w", err)
		}
	}

	// configure modules
	// Terraform modules via Please work by determining the absolute path to the module source and updating the references to it.
	for _, dep := range c.Deps {
		depModule, err := module.Load(dep)
		if err != nil {
			return err
		}
		if err := depModule.UpdateReferences(c.Out); err != nil {
			return err
		}
	}

	// substitute build env vars
	if err := please.ReplaceInDirectory(c.Out, "$PKG", c.Pkg); err != nil {
		return err
	}

	if err := please.ReplaceInDirectory(c.Out, "$PKG_DIR", c.PkgDir); err != nil {
		return err
	}

	if err := please.ReplaceInDirectory(c.Out, "$NAME", c.Name); err != nil {
		return err
	}

	if err := please.ReplaceInDirectory(c.Out, "$ARCH", c.Arch); err != nil {
		return err
	}

	if err := please.ReplaceInDirectory(c.Out, "$OS", c.OS); err != nil {
		return err
	}

	// shift var files into outs
	varFiles := strings.Split(c.VarFilesCSV, ",")
	for i, varFile := range varFiles {
		os.Rename(varFile, filepath.Join(
			c.Out,
			fmt.Sprintf("%d-%s", i,
				strings.ReplaceAll(filepath.Base(varFile), ".tfvars", ".auto.tfvars"),
			)),
		)
	}

	return nil
}
