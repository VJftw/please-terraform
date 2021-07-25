package terraformregistrymodule

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/VJftw/please-terraform/internal/logging"
	"github.com/VJftw/please-terraform/pkg/plz"
	"github.com/VJftw/please-terraform/pkg/terraformregistryprovider"
	getter "github.com/hashicorp/go-getter"
)

type Build struct {
	*TerraformRegistryModule
	Out   string   `long:"out" description:"The path to output the built module to"`
	Deps  []string `long:"deps" description:"The dependencies for this module"`
	Strip []string `long:"strip" description:"The directories to strip from this module"`
	Pkg   string   `long:"pkg"`
}

func (c *Build) Execute(args []string) error {
	logging.Logger.Info().
		Strs("deps", c.Deps).
		Strs("strip", c.Strip).
		Str("pkg", c.Pkg).
		Str("out", c.Out).
		Msg("building with parameters")

	absPlzOut := plz.MustAbsPlzOut(c.Pkg)
	c.AbsPath = filepath.Join(absPlzOut, "gen", c.Pkg, c.Out)

	generatedURL := fmt.Sprintf("%s/v1/modules/%s/%s/%s/%s/download",
		c.Registry,
		c.Namespace,
		c.Name,
		c.Provider,
		c.Version,
	)
	parsedURL, err := url.Parse(generatedURL)
	if err != nil {
		return fmt.Errorf("could not parse generated URL '%s': %w", generatedURL, err)
	}
	logging.Logger.Info().Str("url", parsedURL.String()).Msg("retrieving getter download url from registry")

	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return fmt.Errorf("could not get module getter URL: %w", err)
	}
	getterURL := resp.Header.Get("X-Terraform-Get")

	logging.Logger.Info().Str("url", getterURL).Msg("downloading")
	err = getter.GetAny(c.Out, getterURL)
	if err != nil {
		return fmt.Errorf("could not get module: %w", err)
	}

	// strip
	logging.Logger.Info().Strs("paths", c.Strip).Msg("stripping directories")
	for _, stripDir := range c.Strip {
		stripPath := filepath.Join(c.Out, stripDir)
		if err := os.RemoveAll(stripPath); err != nil {
			return fmt.Errorf("could not remove directory '%s': %w", stripDir, err)
		}
	}

	// Provider dependencies, overwrite versions in SRCS to match.
	// Module dependencies, overwrite source+version to be absolute path of given module.
	modules := []*TerraformRegistryModule{}
	providers := []*terraformregistryprovider.TerraformRegistryProvider{}
	for _, dep := range c.Deps {
		switch {
		case IsModuleDir(dep):
			c, err := LoadFromDir(dep)
			if err != nil {
				return err
			}
			modules = append(modules, c)
		case terraformregistryprovider.IsProviderDir(dep):
			c, err := terraformregistryprovider.LoadFromDir(dep)
			if err != nil {
				return err
			}
			providers = append(providers, c)
		}

	}

	for _, p := range providers {
		if err := p.UpdateVersionReferences(); err != nil {
			return err
		}
	}

	for _, m := range modules {
		if err := m.UpdateReferences(c.Out); err != nil {
			return err
		}
	}

	missingDeps, err := c.FindMissingDependencies(providers)
	if err != nil {
		return err
	}
	_ = missingDeps
	// if len(missingDeps) > 0 {
	// 	return fmt.Errorf("missing deps: %v", missingDeps)
	// }

	return c.Save()
}

func (m *Build) Save() error {
	fileBytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(m.Out, DataFile), fileBytes, 0644)
}

func (c *Build) FindMissingDependencies(providers []*terraformregistryprovider.TerraformRegistryProvider) ([]string, error) {
	missingDeps := []string{}
	// TODO: this should only consider modules (at the moment it's considering providers as well)
	sourceRe := regexp.MustCompile(`source\s*=\s*"[^/\.][^\n"]*"`)
	err := filepath.Walk(c.Out, func(path string, fi os.FileInfo, err error) error {
		if filepath.Ext(path) == ".tf" {
			tfContents, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("could not read '%s': %w", path, err)
			}

			missingMatchesFromDeps := []string{}
			if matches := sourceRe.FindAllString(string(tfContents), -1); len(matches) > 0 {
				for _, m := range matches {
					if !isMatchInDeps(m, providers) {
						missingMatchesFromDeps = append(missingMatchesFromDeps, m)
					}
				}
			}
			if len(missingMatchesFromDeps) > 0 {
				missingDeps = append(missingMatchesFromDeps, missingDeps...)
				// return fmt.Errorf("missing dependencies for '%v' in '%s'", missingMatchesFromDeps, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return missingDeps, nil
}

func isMatchInDeps(match string, providers []*terraformregistryprovider.TerraformRegistryProvider) bool {
	for _, p := range providers {
		if p.IsUsedInString(match) {
			return true
		}
	}

	return false
}
