package module

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/VJftw/please-terraform/internal/logging"
	getter "github.com/hashicorp/go-getter"
)

// RegistryCommand represents a Terraform Registry Module.
type RegistryCommand struct {
	*Module
	Registry  string `long:"registry" description:""`
	Namespace string `long:"namespace" description:""`
	Name      string `long:"name" description:""`
	Version   string `long:"version" description:""`
	Provider  string `long:"provider" description:""`
}

// Execute builds a Terraform Registry Module as a Terraform Module.
func (c *RegistryCommand) Execute(args []string) error {
	downloadURL, err := c.GetDownloadURL()
	if err != nil {
		return err
	}

	if err := c.Download(downloadURL); err != nil {
		return err
	}

	return c.Build()
}

// GetDownloadURL returns the `hashicorp/go-getter` compatible URI.
func (c *RegistryCommand) GetDownloadURL() (string, error) {
	address := fmt.Sprintf("%s/v1/modules/%s/%s/%s/%s/download",
		c.Registry,
		c.Namespace,
		c.Name,
		c.Provider,
		c.Version,
	)

	getterURL, err := url.Parse(address)
	if err != nil {
		return "", fmt.Errorf("could not parse '%s' as URL: %w", address, err)
	}

	logging.Logger.Info().Str("url", getterURL.String()).Msg("retrieving download url from registry")
	resp, err := http.Get(getterURL.String())
	if err != nil {
		return "", fmt.Errorf("could not get download URL: %w", err)
	}

	downloadURL := resp.Header.Get("X-Terraform-Get")

	return downloadURL, nil
}

// Download retrieves the configured Terraform Module from the configured Terraform Registry.
func (c *RegistryCommand) Download(downloadURL string) error {
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
