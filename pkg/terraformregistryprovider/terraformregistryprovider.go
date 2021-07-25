package terraformregistryprovider

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
)

const DataFile = ".please_terraformregistryprovider.json"

type Command struct {
	Build *Build `command:"build"`
}

// TerraformRegistryProvider represents a Provider from a Terraform Registry.
type TerraformRegistryProvider struct {
	Registry  string   `long:"registry" description:""`
	Namespace string   `long:"namespace" description:""`
	Type      string   `long:"type" description:""`
	Version   string   `long:"version" description:""`
	OS        string   `long:"os" description:""`
	Arch      string   `long:"arch" description:""`
	Aliases   []string `long:"aliases" description:""`
	AbsPath   string
}

func LoadFromDir(dir string) (*TerraformRegistryProvider, error) {
	fileBytes, err := os.ReadFile(filepath.Join(dir, DataFile))
	if err != nil {
		return nil, err
	}

	provider := &TerraformRegistryProvider{}
	if err := json.Unmarshal(fileBytes, provider); err != nil {
		return nil, err
	}

	return provider, nil
}

// UpdateVersionReferences will enumerate the Terraform files in the given path
// to replace references to this provider's aliases to this provider's version.
// TODO: use // hcl2Parser := hcl2parse.NewParser()
func (p *TerraformRegistryProvider) UpdateVersionReferences() error {

	return nil
}

// CreateProviderInstallationFilesystemMirrorBlock creates the `filesystem_mirror` block
// for the `provider_installation` block for use with the Terraform CLI. This makes Terraform
// use this Terraform provider downloaded via Please.
// (https://www.terraform.io/docs/cli/config/config-file.html#provider-installation)
func (p *TerraformRegistryProvider) CreateProviderInstallationFilesystemMirrorBlock() {

}

func (o *TerraformRegistryProvider) IsUsedInString(haystack string) bool {
	sourceValueRe := regexp.MustCompile(`source\s*=\s*"(.*)"`)
	matchValues := sourceValueRe.FindStringSubmatch(haystack)
	if len(matchValues) < 2 {
		return false
	}
	for _, alias := range o.Aliases {
		if alias == matchValues[1] {
			return true
		}
	}

	return false
}
