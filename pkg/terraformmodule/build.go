package terraformmodule

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Build struct {
	Name            string            `long:"name" description:""`
	ModuleDirectory string            `long:"module_directory" description:"The path to the module directory"`
	Out             string            `long:"out" description:"The path to output the built module to"`
	DependencyPaths map[string]string `long:"dependency_paths" description:"The dependencies for this module"`
	Pkg             string            `long:"pkg"`
	Srcs            []string          `long:"srcs"`
	fullPath        string
}

func (c *Build) Execute(args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	RepoRoot := filepath.Dir(filepath.Dir(strings.TrimSuffix(filepath.Dir(cwd), c.Pkg)))
	_ = RepoRoot

	// Dependants
	jsonOpts, err := json.Marshal(c)
	if err != nil {
		return err
	}
	optsFilePath := filepath.Join(c.Out, ".please_terraformregistrymodule")
	if err := os.WriteFile(optsFilePath, jsonOpts, 0666); err != nil {
		return fmt.Errorf("could not write '%s': %w", optsFilePath, err)
	}

	// Dependencies
	log.Printf("Dependencies: %+v", c.DependencyPaths)
	// version is a reserved word in Terraform
	versionRE := regexp.MustCompile(`version\s*=\s*"[^\n"]*"`)
	// hcl2Parser := hcl2parse.NewParser()
	modules := []*Build{}
	providers := []*terraformregistryprovider.Build{}
	for dep, dep_out := range c.DependencyPaths {
		switch {
		case IsModuleDir(dep):
			c, err := LoadOptsDataFromDir(dep)
			if err != nil {
				return err
			}
			c.fullPath = filepath.Join(RepoRoot, dep_out)
			modules = append(modules, c)
		case terraformregistryprovider.IsProviderDir(dep):
			c, err := terraformregistryprovider.LoadOptsDataFromDir(dep)
			if err != nil {
				return err
			}
			providers = append(providers, c)
		}

	}
	return nil
}
