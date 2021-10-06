package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/VJftw/please-terraform/internal/logging"
	"github.com/VJftw/please-terraform/pkg/please"
)

var log = logging.NewLogger()

type Command struct {
	Local    *LocalCommand    `command:"local"`
	Registry *RegistryCommand `command:"registry"`
}

type Module struct {
	Aliases []string `long:"aliases" description:""`
	Pkg     string   `long:"pkg" description:""`
	Strip   []string `long:"strip" description:""`
	Deps    []string `long:"deps" description:""`
	Out     string   `long:"out" description:""`

	AbsolutePath string
}

// Load returns a Module loaded from the given directory.
func Load(directory string) (*Module, error) {
	fileBytes, err := os.ReadFile(filepath.Join(directory, ".please_terraform.json"))
	if err != nil {
		return nil, fmt.Errorf("could not read saved module: %w", err)
	}

	m := &Module{}
	if err := json.Unmarshal(fileBytes, m); err != nil {
		return nil, fmt.Errorf("could not unmarshal saved module: %w", err)
	}

	if m.AbsolutePath == "" {
		return nil, errors.New("absolute path missing from loaded module")
	}

	return m, nil
}

// Save saves the Module data to be re-used in other builds.
func (m *Module) Save() error {
	// TODO: add APIVersion and Kind to Module (e.g. terraform.please.build/v1/Module)
	fileBytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal module: %w", err)
	}

	return os.WriteFile(filepath.Join(m.Out, ".please_terraform.json"), fileBytes, 0644)
}

// Build prepares the module to be used by Please.
func (m *Module) Build() error {
	m.AbsolutePath = filepath.Join(please.MustAbsPlzOut(m.Pkg), "gen", m.Pkg, m.Out)

	if err := m.StripDirs(); err != nil {
		return err
	}

	for _, dep := range m.Deps {
		depModule, err := Load(dep)
		if err != nil {
			return err
		}
		if err := depModule.UpdateReferences(m.Out); err != nil {
			return err
		}
	}

	return m.Save()
}

// Strip strips the configured directories from the module.
func (m *Module) StripDirs() error {
	for _, stripDir := range m.Strip {
		stripPath := filepath.Join(m.Out, stripDir)
		if err := os.RemoveAll(stripPath); err != nil {
			return fmt.Errorf("could not remove directory '%s': %w", stripDir, err)
		}
	}

	return nil
}

// UpdateReferences replaces all references to this Module's Aliases to it's absolute path in the given directory.
func (m *Module) UpdateReferences(directory string) error {
	versionRE := regexp.MustCompile(`version\s*=\s*"[^\n"]*"`)

	err := filepath.Walk(directory, func(path string, fi os.FileInfo, err error) error {
		for _, alias := range m.Aliases {
			if filepath.Ext(path) == ".tf" {
				log.Debug().
					Str("path", path).
					Str("alias", alias).
					Msg("replacing module sources")
				pattern := fmt.Sprintf(`source\s*=\s*"%s"`, alias)
				re := regexp.MustCompile(pattern)
				tfContents, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("could not read '%s': %w", path, err)
				}

				newContents := re.ReplaceAll(tfContents, []byte(fmt.Sprintf("source = \"%s\"", m.AbsolutePath)))
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
