package main

import (
	"context"

	"github.com/VJftw/please-terraform/internal/cmd"
	"github.com/VJftw/please-terraform/pkg/module"
	"github.com/VJftw/please-terraform/pkg/root"
)

type opts struct {
	BuildModule    *BuildModule    `command:"build_module"`
	DownloadModule *DownloadModule `command:"download_module"`
	Root           *root.Command   `command:"root"`
}

func main() {
	opts := &opts{}
	cmd.MustParseFlags(opts)
}

type BuildModule struct {
	Name string `long:"name" required:"true" description:"The Please name of the Terraform module."`
	Pkg  string `long:"pkg" required:"true" description:"The Please package of the Terraform module."`
	Srcs string `long:"srcs" requrired:"true" description:"Space separated src files that make up the Terraform module. These will be flattened."`
	Out  string `long:"out" required:"true" description:"The directory to write the processed Terraform module to."`

	Aliases []string `long:"aliases" required:"false" description:"The aliases for the Terraform module that will be replaced in other Terraform configuration."`
	Strip   []string `long:"strip" required:"false" description:"The directories to strip from the Terraform module."`
	Deps    []string `long:"deps" required:"false" description:"Other Terraform modules that this Terraform module depends on."`
}

func (c *BuildModule) Execute(args []string) error {
	moduleBuilder, err := module.NewBuilder(
		module.WithBuilderPleaseTargetName(c.Name),
		module.WithBuilderPleaseTargetPackage(c.Pkg),
		module.WithBuilderSpaceSeparatedSrcs(c.Srcs),
		module.WithBuilderOutputDir(c.Out),

		module.WithBuilderAliases(c.Aliases),
		module.WithBuilderStrip(c.Strip),
		module.WithBuilderDeps(c.Deps),
	)
	if err != nil {
		return err
	}

	return moduleBuilder.Build(context.Background())
}

type DownloadModule struct {
	Registry   string `long:"registry" description:""`
	Namespace  string `long:"namespace" description:""`
	ModuleName string `long:"module_name"`
	Version    string `long:"version" description:""`
	System     string `long:"system" description:""`
	Out        string `long:"out" description:""`
}

func (c *DownloadModule) Execute(args []string) error {
	moduleDownloader, err := module.NewDownloader(
		module.WithDownloaderRegistry(c.Registry),
		module.WithDownloaderNamespace(c.Namespace),
		module.WithDownloaderModuleName(c.ModuleName),
		module.WithDownloaderVersion(c.Version),
		module.WithDownloaderSystem(c.System),
		module.WithDownloaderOutDir(c.Out),
	)
	if err != nil {
		return err
	}

	return moduleDownloader.Download(context.Background())
}
