package root

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/VJftw/please-terraform/pkg/module"
	"github.com/VJftw/please-terraform/pkg/please"
)

// CommandBuild represents the build subcommand.
type CommandBuild struct {
	Pkg      string   `long:"pkg"`
	Name     string   `long:"name"`
	OS       string   `long:"os"`
	Arch     string   `long:"arch"`
	Out      string   `long:"out"`
	PkgDir   string   `long:"pkg_dir"`
	Srcs     string   `long:"srcs"`
	VarFiles []string `long:"var_files"`
	Modules  []string `long:"modules"`

	ModuleOpts *module.Opts
}

// Execute executes the build subcommand.
func (c *CommandBuild) Execute(args []string) error {
	// make the out directory
	if err := os.MkdirAll(c.Out, os.ModePerm); err != nil {
		return err
	}

	// flatten
	srcs := strings.Split(c.Srcs, " ")
	for _, src := range srcs {
		// flatten
		if err := please.CopyFile(src, filepath.Join(c.Out, filepath.Base(src))); err != nil {
			return fmt.Errorf("could not move file: %w", err)
		}
	}

	// Substitute build env vars into srcs
	// This is useful for re-using a source file in multiple Terraform roots
	// such as templating a Terraform remote state configuration.
	err := filepath.Walk(c.Out, func(path string, fi os.FileInfo, err error) error {
		if filepath.Ext(path) == ".tf" {
			log.Debug().
				Str("path", path).
				Msg("substituing env vars")

			tfContents, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("could not read '%s': %w", path, err)
			}

			newContents := bytes.ReplaceAll(tfContents, []byte("$PKG"), []byte(c.Pkg))
			newContents = bytes.ReplaceAll(newContents, []byte("$PKG_DIR"), []byte(c.PkgDir))
			newContents = bytes.ReplaceAll(newContents, []byte("$NAME"), []byte(c.Name))
			newContents = bytes.ReplaceAll(newContents, []byte("$ARCH"), []byte(c.Arch))
			newContents = bytes.ReplaceAll(newContents, []byte("$OS"), []byte(c.OS))

			if err := os.WriteFile(path, newContents, 0644); err != nil {
				return fmt.Errorf("could not write file '%s': %w", path, err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Shift var files into outs so that they are auto-loaded.
	for i, varFile := range c.VarFiles {
		newName, err := AutoTFVarsName(i, varFile)
		if err != nil {
			return err
		}
		if err := please.CopyFile(varFile, filepath.Join(c.Out, newName)); err != nil {
			return fmt.Errorf("could not move var file '%s': %w", varFile, err)
		}
		log.Debug().Str("old_tfvars_name", varFile).Str("new_tfvars_name", newName).Msg("configured tfvars")
	}

	// colocate modules
	if err := module.ColocateModules(c.ModuleOpts.MetadataFile, c.Out, c.Modules); err != nil {
		return err
	}

	return nil
}

// AutoTFVarsName returns a tfvars file name that will be automatically be
// loaded by Terraform for given index and var file.
func AutoTFVarsName(i int, varFile string) (string, error) {
	baseName := filepath.Base(varFile)
	switch {
	case strings.HasSuffix(baseName, ".tfvars"):
		return fmt.Sprintf("%d-%s.auto.tfvars", i, strings.TrimSuffix(baseName, ".tfvars")), nil
	case strings.HasSuffix(baseName, ".tfvars.json"):
		return fmt.Sprintf("%d-%s.auto.tfvars.json", i, strings.TrimSuffix(baseName, ".tfvars.json")), nil
	}

	return "", fmt.Errorf("'%s' does not end in '.tfvars' or '.tfvars.json'", baseName)
}
