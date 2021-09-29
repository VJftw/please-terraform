package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/VJftw/please-terraform/pkg/please"
)

type Command struct {
	Registry *RegistryCommand `command:"registry"`
}

type Provider struct {
	Aliases   []string `long:"aliases" description:""`
	Registry  string   `long:"registry" description:""`
	Namespace string   `long:"namespace" description:""`
	Pkg       string   `long:"pkg"`
	Type      string   `long:"type" description:""`
	Version   string   `long:"version" description:""`
	OS        string   `long:"os" description:""`
	Arch      string   `long:"arch" description:""`

	Out          string `long:"out" description:""`
	AbsolutePath string
}

// Load returns a Provider loaded from the given directory.
func Load(directory string) (*Provider, error) {
	fileBytes, err := os.ReadFile(filepath.Join(directory, ".please_terraform.json"))
	if err != nil {
		return nil, fmt.Errorf("could not read saved module: %w", err)
	}

	p := &Provider{}
	if err := json.Unmarshal(fileBytes, p); err != nil {
		return nil, fmt.Errorf("could not unmarshal saved module: %w", err)
	}

	return p, nil
}

// Save saves the Provider data to be re-used in other builds.
func (p *Provider) Save() error {
	fileBytes, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal module: %w", err)
	}

	return os.WriteFile(filepath.Join(p.Out, ".please_terraform.json"), fileBytes, 0644)
}

// Build prepares the module to be used by Please.
func (p *Provider) Build() error {
	p.AbsolutePath = filepath.Join(please.MustAbsPlzOut(p.Pkg), "gen", p.Pkg, p.Out)

	return p.Save()
}

func (p *Provider) UpdateReferences() error {
	return nil
}
