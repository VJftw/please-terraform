package terraformregistrymodule

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/VJftw/please-terraform/internal/logging"
)

const DataFile = ".please_terraformregistrymodule.json"

type Command struct {
	Build *Build `command:"build"`
}

type TerraformRegistryModule struct {
	Registry  string   `long:"registry" description:""`
	Namespace string   `long:"namespace" description:""`
	Name      string   `long:"name" description:""`
	Version   string   `long:"version" description:""`
	Provider  string   `long:"provider" description:""`
	Aliases   []string `long:"aliases" description:"Aliases for the module in CSV"`
	AbsPath   string
}

func LoadFromDir(dir string) (*TerraformRegistryModule, error) {
	fileBytes, err := os.ReadFile(filepath.Join(dir, DataFile))
	if err != nil {
		return nil, err
	}

	provider := &TerraformRegistryModule{}
	if err := json.Unmarshal(fileBytes, provider); err != nil {
		return nil, err
	}

	return provider, nil
}

// UpdateReferences will enumerate the Terraform files in the given path
// to replace references to this module's aliases to the absolute path of this module.
// TODO: use // hcl2Parser := hcl2parse.NewParser()
func (m *TerraformRegistryModule) UpdateReferences(path string) error {
	versionRE := regexp.MustCompile(`version\s*=\s*"[^\n"]*"`)

	err := filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
		for _, alias := range m.Aliases {
			if filepath.Ext(path) == ".tf" {
				pattern := fmt.Sprintf(`source\s*=\s*"%s"`, alias)
				re := regexp.MustCompile(pattern)
				tfContents, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("could not read '%s': %w", path, err)
				}

				newContents := re.ReplaceAll(tfContents, []byte(fmt.Sprintf("source = \"%s\"", m.AbsPath)))
				newContents = versionRE.ReplaceAll(newContents, []byte{})
				if err := os.WriteFile(path, newContents, 0644); err != nil {
					return fmt.Errorf("could not write file '%s': %w", path, err)
				}
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("could not walk files: %w", err)
	}
	return nil
}

func IsModuleDir(path string) bool {
	logging.Logger.Debug().Str("path", path).Msg("is this a module directory?")
	moduleOptsFilePath := filepath.Join(path, DataFile)
	if _, err := os.Stat(moduleOptsFilePath); os.IsNotExist(err) {
		return false
	}

	return true
}
