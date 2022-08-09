package module

// Opts represent the available options to this Module package as whole.
type Opts struct {
	// The file in which Please module metadata is stored in relative to each module.
	MetadataFile string `long:"metadata_file" default:".please/terraform/module.json" description:"The file in which Please module metadata is stored in relative to each module."`
}

// Command represents the `module` command and available subcommands.
type Command struct {
	Local    *CommandLocal    `command:"local"`
	Registry *CommandRegistry `command:"registry"`
}
