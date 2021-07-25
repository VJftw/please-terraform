package terraformregistryprovider

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/VJftw/please-terraform/pkg/plz"
	getter "github.com/hashicorp/go-getter"
)

const (
	OptsDataFile = ".please_terraformregistryprovider"
)

type Build struct {
	*TerraformRegistryProvider
	Out string `long:"out" description:"The output path of the module"`
	Pkg string `long:"pkg" description:"The Please package being built against."`
}

func (p *Build) Save() error {
	fileBytes, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(p.Out, DataFile), fileBytes, 0644)
}

func (c *Build) Execute(args []string) error {
	repoRoot := plz.MustAbsPlzOut(c.Pkg)
	c.AbsPath = filepath.Join(repoRoot, "gen", c.Pkg, c.Out)
	generatedURL := fmt.Sprintf("%s/v1/providers/%s/%s/%s/download/%s/%s",
		c.Registry,
		c.Namespace,
		c.Type,
		c.Version,
		c.OS,
		c.Arch,
	)
	parsedURL, err := url.Parse(generatedURL)
	if err != nil {
		return fmt.Errorf("could not parse generated URL '%s': %w", generatedURL, err)
	}
	log.Printf("Obtaining Provider Download URL from '%s'", parsedURL)

	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return fmt.Errorf("could not get URL: %w", err)
	}
	defer resp.Body.Close()

	providerDownloadResponse := &ProviderDownloadResponse{}
	if err := json.NewDecoder(resp.Body).Decode(providerDownloadResponse); err != nil {
		return fmt.Errorf("could not unmarshal data from '%s'", parsedURL)
	}

	getterURL := providerDownloadResponse.DownloadURL

	unpackedLayoutDir := fmt.Sprintf("%s/%s/%s/%s/%s",
		parsedURL.Host,
		c.Namespace,
		c.Type,
		c.Version,
		strings.ToLower(fmt.Sprintf("%s_%s", c.OS, c.Arch)),
	)

	err = getter.GetAny(filepath.Join(c.Out, unpackedLayoutDir), getterURL)
	if err != nil {
		return fmt.Errorf("could not download provider: %w", err)
	}

	return c.Save()
}

// ProviderDownloadResponse represents the data returned by https://www.terraform.io/docs/internals/provider-registry-protocol.html#find-a-provider-package
type ProviderDownloadResponse struct {
	Protocols           []string                            `json:"protocols"`
	OS                  string                              `json:"os"`
	Arch                string                              `json:"arch"`
	Filename            string                              `json:"filename"`
	DownloadURL         string                              `json:"download_url"`
	SHASumsURL          string                              `json:"shasums_url"`
	SHASumsSignatureURL string                              `json:"shasums_signature_url"`
	SHASum              string                              `json:"shasum"`
	SigningKeys         ProviderDownloadResponseSigningKeys `json:"signing_keys"`
}

type ProviderDownloadResponseSigningKeys struct {
	GPGPublicKeys []ProviderDownloadResponseGPGPublicKey `json:"gpg_public_keys"`
}

type ProviderDownloadResponseGPGPublicKey struct {
	KeyID          string `json:"key_id"`
	ASCIIArmor     string `json:"ascii_armor"`
	TrustSignature string `json:"trust_signature"`
	Source         string `json:"Hashicorp"`
	SourceURL      string `json:"source_url"`
}

func LoadOptsDataFromDir(path string) (*Build, error) {
	providerOptsFilePath := filepath.Join(path, OptsDataFile)
	providerOptsBytes, err := os.ReadFile(providerOptsFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read c from '%s': %w", providerOptsFilePath, err)
	}
	providerOpts := &Build{}
	if err := json.Unmarshal(providerOptsBytes, &providerOpts); err != nil {
		return nil, fmt.Errorf("could not unmarshal c from '%s': %w", providerOptsFilePath, err)
	}

	return providerOpts, nil
}

func IsProviderDir(path string) bool {
	providerOptsFilePath := filepath.Join(path, OptsDataFile)
	if _, err := os.Stat(providerOptsFilePath); os.IsNotExist(err) {
		return false
	}

	return true
}
