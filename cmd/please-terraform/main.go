package main

import (
	"github.com/VJftw/please-terraform/internal/cmd"
	"github.com/VJftw/please-terraform/pkg/terraformregistrymodule"
	"github.com/VJftw/please-terraform/pkg/terraformregistryprovider"
)

type Opts struct {
	// Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	TerraformRegistryModule   *terraformregistrymodule.Command   `command:"terraformregistrymodule"`
	TerraformRegistryProvider *terraformregistryprovider.Command `command:"terraformregistryprovider"`
	// TerraformModule           *terraformmodule.TerraformModule                     `command:"terraformmodule"`
}

func main() {
	opts := &Opts{}
	cmd.MustParseFlags(opts)

	// switch args[0] {
	// case "terraformregistrymodule_build":
	// 	err = terraformregistrymodule.Build(&globalOpts.TerraformRegistryModuleBuild.Build)
	// case "terraformregistryprovider_build":
	// 	err = terraformregistryprovider.Build(&globalOpts.TerraformRegistryProviderBuild)
	// default:
	// 	log.Fatalf("invalid argument: %v", args[0])
	// }

	// if err != nil {
	// 	log.Fatal(err)
	// }
}
