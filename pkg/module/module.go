package module

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/VJftw/please-terraform/internal/logging"
	"github.com/VJftw/please-terraform/pkg/please"
)

var log = logging.NewLogger()

// Metadata represents a module's metadata.
type Metadata struct {
	Target  string
	Aliases []string
}

// Load returns a module's Metadata loaded from the given directory.
func Load(path string) (*Metadata, error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read saved module: %w", err)
	}

	m := &Metadata{}
	if err := json.Unmarshal(fileBytes, m); err != nil {
		return nil, fmt.Errorf("could not unmarshal saved module: %w", err)
	}

	return m, nil
}

// Save saves the Metadata data to be re-used in other workflows.
func (m *Metadata) Save(path string) error {
	fileBytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal module: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("could not create directory '%s': %w", dir, err)
	}

	if err := os.WriteFile(path, fileBytes, 0644); err != nil {
		return fmt.Errorf("could not write '%s': %w", path, err)
	}

	log.Debug().Str("file", path).Msg("saved module metadata")
	return nil
}

// StripDirs strips the configured directories from the module.
func (m *Metadata) StripDirs(out string, strip []string) error {
	for _, stripDir := range strip {
		stripPath := filepath.Join(out, stripDir)
		if err := os.RemoveAll(stripPath); err != nil {
			return fmt.Errorf("could not remove directory '%s': %w", stripDir, err)
		}
	}

	return nil
}

// ColocateModules colocates the given module paths to the given out directory.
func ColocateModules(metadataFilePath string, out string, modulePaths []string) error {
	log.Debug().Strs("modulePaths", modulePaths).Msg("colocating modules")
	if len(modulePaths) < 1 {
		log.Debug().Msg("no modules to colocate")
		return nil
	}

	modulesDir := filepath.Join(out, ".modules")
	if err := os.MkdirAll(modulesDir, 0750); err != nil {
		return fmt.Errorf("could not create modules dir '%s': %w", modulesDir, err)
	}

	for _, modulePath := range modulePaths {
		replace := fmt.Sprintf(".%c%s", filepath.Separator, filepath.Join(".modules", modulePath))

		moduleMeta, err := Load(filepath.Join(modulePath, metadataFilePath))
		if err != nil {
			return err
		}

		for _, alias := range moduleMeta.Aliases {
			log.Debug().Str("alias", alias).Str("path", replace).Msg("replacing in module")
			if err := please.ReplaceInDirectory(out, alias, replace); err != nil {
				return err
			}
		}

		if err := os.MkdirAll(filepath.Dir(filepath.Join(out, replace)), 0750); err != nil {
			return err
		}
		if err := please.Sync(modulePath, filepath.Join(out, replace), []string{}); err != nil {
			return err
		}
	}

	return nil
}
