package root

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/VJftw/please-terraform/pkg/please"
)

// CommandVirtualEnv represents the virtualenv subcommand.
type CommandVirtualEnv struct {
	TerraformBinary   string `long:"terraform_binary"`
	OS                string `long:"os"`
	Arch              string `long:"arch"`
	RootModule        string `long:"root_module"`
	VirtualEnvBaseDir string `long:"virtual_env_base_dir" default:"./plz-out/terraform/venvs"`

	PleaseOpts *please.Opts
}

// Execute executes the virtualenv subcommand.
func (c *CommandVirtualEnv) Execute(args []string) error {
	repoRoot := please.MustRepoRoot(c.PleaseOpts.PlzOutDir)
	if strings.HasPrefix(c.VirtualEnvBaseDir, "./") {
		c.VirtualEnvBaseDir = filepath.Join(repoRoot, c.VirtualEnvBaseDir)
	}

	// Use absolute path to Terraform binary as we're using this at
	// `plz run ...` time where the current working directory won't be in the
	// root of the repository.
	onlyStartsWithPlzOutRegex := regexp.MustCompile(
		fmt.Sprintf(`^%c?%s`, filepath.Separator, c.PleaseOpts.PlzOutDir),
	)
	if onlyStartsWithPlzOutRegex.MatchString(c.TerraformBinary) {
		c.TerraformBinary = filepath.Join(repoRoot, c.TerraformBinary)
		log.Debug().
			Str("path", c.TerraformBinary).
			Msg("updated terraform_binary to absolute path")
	}

	// Add the `terraform` binary to the sourcing user's path.
	fmt.Printf(`PATH="%s:$PATH"`+"\n", filepath.Dir(c.TerraformBinary))
	fmt.Println(`export PATH`)

	// We cannot run Terraform commands in the `plz-out/gen/<rule>` directory
	// as Terraform creates symlinks which Please warns us will be removed.
	// We also cannot use a `plz-out/terraform` directory as Terraform plans
	// include absolute paths to files referenced in Terraform configuration so
	// we would not be able to use a plan file on another computer which is
	// common in CI/CD.
	// Instead, we create a 'virtualenv' directory outside of the repository and
	// copy the generated terraform root. We then replace modules in the
	// Terraform root with their absolute paths under `plz-out/gen`.
	virtualEnvDir := filepath.Join(
		c.VirtualEnvBaseDir,
		c.RootModule,
	)

	if err := os.MkdirAll(virtualEnvDir, 0750); err != nil {
		return fmt.Errorf("could not create virtual env dir '%s': %w", virtualEnvDir, err)
	}

	if err := please.Sync(c.RootModule, virtualEnvDir, []string{
		`\.terraform.*`,
		`.*\.tfstate`,
	}); err != nil {
		return err
	}

	// add symbolic link to plz-out
	old := filepath.Join(repoRoot, c.PleaseOpts.PlzOutDir)
	new := filepath.Join(virtualEnvDir, c.PleaseOpts.PlzOutDir)
	if err := os.Symlink(old, new); err != nil {
		return fmt.Errorf("could not create symlink from '%s' to '%s': %w", old, new, err)
	}

	// set GIT_DIR
	fmt.Printf(`GIT_DIR=%s/.git`+"\n", repoRoot)
	fmt.Println(`export GIT_DIR`)

	// change user's working directory.
	log.Info().Str("path", virtualEnvDir).Msg("working directory")
	fmt.Printf(`cd "%s"`+"\n", virtualEnvDir)

	return nil
}
