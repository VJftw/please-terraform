package module

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/VJftw/please-terraform/internal/logging"
	"github.com/hashicorp/go-getter"
)

// CommandRegistry represents the `module local` command and its flags.
type CommandRegistry struct {
	Name       string   `long:"name"`
	Registry   string   `long:"registry" description:""`
	Namespace  string   `long:"namespace" description:""`
	ModuleName string   `long:"module_name"`
	Version    string   `long:"version" description:""`
	Provider   string   `long:"provider" description:""`
	Aliases    []string `long:"aliases" description:""`
	Pkg        string   `long:"pkg" description:""`
	Strip      []string `long:"strip" description:""`
	Deps       []string `long:"deps" description:""`
	Out        string   `long:"out" description:""`

	Opts *Opts
}

// Execute builds a Terraform Registry Module as a Terraform Module.
func (c *CommandRegistry) Execute(args []string) error {
	m := &Metadata{
		Target:  fmt.Sprintf("//%s:%s", c.Pkg, c.Name),
		Aliases: c.Aliases,
	}

	downloadURL, err := c.GetDownloadURL()
	if err != nil {
		return err
	}

	if err := c.Download(downloadURL); err != nil {
		return err
	}

	// Strip directories
	if err := m.StripDirs(c.Out, c.Strip); err != nil {
		return fmt.Errorf("could not strip directories: %w", err)
	}

	m.Aliases = append(m.Aliases, []string{
		// Supports referencing by Please target.
		fmt.Sprintf("//%s:%s", c.Pkg, c.Name),
		// Supports referencing by Terraform Module Registry.
		fmt.Sprintf("%s/%s/%s", c.Namespace, c.ModuleName, c.Provider),
	}...)

	if filepath.Base(c.Pkg) == c.Name {
		// Supports referencing by default Please target for pkg.
		m.Aliases = append(m.Aliases, fmt.Sprintf("//%s", c.Pkg))
	}

	// colocate modules
	if err := ColocateModules(c.Opts.MetadataFile, c.Out, c.Deps); err != nil {
		return err
	}

	if err := m.Save(filepath.Join(c.Out, c.Opts.MetadataFile)); err != nil {
		return err
	}

	return nil
}

// GetDownloadURL returns the `hashicorp/go-getter` compatible URI.
func (c *CommandRegistry) GetDownloadURL() (string, error) {
	address := fmt.Sprintf("%s/v1/modules/%s/%s/%s/%s/download",
		c.Registry,
		c.Namespace,
		c.ModuleName,
		c.Provider,
		c.Version,
	)

	getterURL, err := url.Parse(address)
	if err != nil {
		return "", fmt.Errorf("could not parse '%s' as URL: %w", address, err)
	}

	log.Info().Str("url", getterURL.String()).Msg("retrieving download url from registry")
	resp, err := http.Get(getterURL.String())
	if err != nil {
		return "", fmt.Errorf("could not get download URL: %w", err)
	}

	downloadURL := resp.Header.Get("X-Terraform-Get")

	return downloadURL, nil
}

// Download retrieves the configured Terraform Module from the configured Terraform Registry.
func (c *CommandRegistry) Download(downloadURL string) error {
	logging.Logger.Info().Str("url", downloadURL).Msg("downloading")
	if err := getter.GetAny(c.Out, downloadURL); err != nil {
		return fmt.Errorf("could not get '%s': %w", downloadURL, err)
	}

	// purge .git directories for consistent hashes.
	purgeDirs := []string{".git"}
	for _, purgeDir := range purgeDirs {
		if err := os.RemoveAll(filepath.Join(c.Out, purgeDir)); err != nil {
			return fmt.Errorf("could not prune '%s' folder: %w", purgeDir, err)
		}
	}

	return nil
}
