package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/VJftw/please-terraform/internal/logging"
	"github.com/hashicorp/go-getter"
)

type RegistryCommand struct {
	*Provider
}

// Execute builds a Terraform Registry Module as a Terraform Module.
func (c *RegistryCommand) Execute(args []string) error {
	host, downloadURL, err := c.GetDownloadURL()
	if err != nil {
		return err
	}

	if err := c.Download(host, downloadURL); err != nil {
		return err
	}

	return c.Build()
}

// GetDownloadURL returns the host and `hashicorp/go-getter` compatible URI.
func (c *RegistryCommand) GetDownloadURL() (string, string, error) {
	address := fmt.Sprintf("%s/v1/providers/%s/%s/%s/download/%s/%s",
		c.Registry,
		c.Namespace,
		c.Type,
		c.Version,
		c.OS,
		c.Arch,
	)

	getterURL, err := url.Parse(address)
	if err != nil {
		return "", "", fmt.Errorf("could not parse '%s' as URL: %w", address, err)
	}

	logging.Logger.Info().Str("url", getterURL.String()).Msg("retrieving download url from registry")
	resp, err := http.Get(getterURL.String())
	if err != nil {
		return "", "", fmt.Errorf("could not get download URL: %w", err)
	}

	providerDownloadResponse := &ProviderDownloadResponse{}
	if err := json.NewDecoder(resp.Body).Decode(providerDownloadResponse); err != nil {
		return "", "", fmt.Errorf("could not unmarshal data from '%s'", getterURL.String())
	}

	return getterURL.Host, providerDownloadResponse.DownloadURL, nil
}

// Download retrieves the configured Terraform Module from the configured Terraform Registry.
func (c *RegistryCommand) Download(host, downloadURL string) error {
	unpackedLayoutDir := fmt.Sprintf("%s/%s/%s/%s/%s",
		host,
		c.Namespace,
		c.Type,
		c.Version,
		strings.ToLower(fmt.Sprintf("%s_%s", c.OS, c.Arch)),
	)
	logging.Logger.Info().Str("url", downloadURL).Msg("downloading")

	if err := getter.GetAny(filepath.Join(c.Out, unpackedLayoutDir), downloadURL); err != nil {
		return fmt.Errorf("could not get '%s': %w", downloadURL, err)
	}

	return nil
}

// ProviderDownloadResponse represents the data returned by https://www.terraform.io/docs/internals/provider-registry-protocol.html#find-a-provider-package
type ProviderDownloadResponse struct {
	Protocols           []string `json:"protocols"`
	OS                  string   `json:"os"`
	Arch                string   `json:"arch"`
	Filename            string   `json:"filename"`
	DownloadURL         string   `json:"download_url"`
	SHASumsURL          string   `json:"shasums_url"`
	SHASumsSignatureURL string   `json:"shasums_signature_url"`
	SHASum              string   `json:"shasum"`
}
