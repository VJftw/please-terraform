package terraformmodule

type TerraformModule struct {
	// Build BuildOpts `command:"build"`
	Build *Build `command:"build"`
}

// func (c *TerraformRegistryModule) Execute(args []string) error {
// 	logging.Logger.Info().Msg("hello world")
// 	return nil
// }
