package run

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/VJftw/please-terraform/internal/logging"
)

var log = logging.NewLogger()

type Command struct {
	TerraformBinary string `long:"terraform_binary"`
	OS              string `long:"os"`
	Arch            string `long:"arch"`
	TerraformRoot   string `long:"terraform_root"`
	PreBinaries     string `long:"pre_binaries"`
	PostBinaries    string `long:"post_binaries"`
	ExtraArgs       string `long:"extra_args"`
	Plugins         string `long:"plugins"`
	Pkg             string `long:"pkg"`
	Name            string `long:"name"`
}

func (c *Command) Execute(args []string) error {

	log.Debug().Interface("flags", c).Msg("flags")

	// TODO: load repo root from TerraformRoot.
	// TODO: Write terraform binary location to TerraformRoot.

	repoRoot, err := os.Getwd()
	if err != nil {
		return err
	}

	c.TerraformRoot = filepath.Join(repoRoot, c.TerraformRoot)
	c.TerraformBinary = filepath.Join(repoRoot, c.TerraformBinary)

	pluginDir := filepath.Join(c.TerraformRoot, "_plugins")
	if err := os.MkdirAll(pluginDir, os.ModePerm); err != nil {
		return err
	}

	// TODO: add plugins caching support

	absPlzOut := filepath.Join(repoRoot, "plz-out")

	terraformWorkspace := filepath.Join(absPlzOut, "terraform", c.Pkg, c.Name)
	if err := os.MkdirAll(terraformWorkspace, os.ModePerm); err != nil {
		return err
	}

	// rsync replacement
	// remove files, skipping .terraform*, *.tfstate*, _plugins
	dirList, err := ioutil.ReadDir(terraformWorkspace)
	if err != nil {
		return fmt.Errorf("could not read directory '%s': %w", terraformWorkspace, err)
	}

	for _, f := range dirList {
		name := f.Name()
		switch {
		case strings.HasPrefix(name, ".terraform"):
			fallthrough
		case strings.Contains(name, ".tfstate"):
			fallthrough
		case strings.HasSuffix(name, "_plugins"):
			log.Debug().
				Str("path", name).
				Msg("skipping")
			continue
		default:
			log.Debug().
				Str("path", name).
				Msg("removing")
			if err := os.RemoveAll(name); err != nil {
				return fmt.Errorf("could not remove '%s': %w", name, err)
			}
		}
	}

	terraformRootFiles, err := ioutil.ReadDir(c.TerraformRoot)
	if err != nil {
		return fmt.Errorf("could not read directory '%s': %w", c.TerraformRoot, err)
	}

	for _, f := range terraformRootFiles {
		if f.Name() != "." && f.Name() != ".." && !f.IsDir() {
			if err := Copy(filepath.Join(c.TerraformRoot, f.Name()), filepath.Join(terraformWorkspace, f.Name())); err != nil {
				return err
			}
		}
	}

	// TODO: execute pre-binaries

	cmd := exec.Command("/bin/sh", "-xc", "terraform "+c.ExtraArgs)
	cmd.Dir = terraformWorkspace
	cmd.Env = []string{"PATH=" + filepath.Dir(c.TerraformBinary) + ":" + os.Getenv("PATH")}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// TODO: clean output

	// TODO: execute post-binaries

	return nil
}

func Copy(src, dst string) error {
	log.Debug().
		Str("src", src).Str("dest", dst).Msg("copying")
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
