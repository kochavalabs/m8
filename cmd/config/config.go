package config

import (
	"github.com/spf13/cobra"
)

func ConfigurationCmdChain() *cobra.Command {
	cfgRootCmd := &cobra.Command{
		Use:   "cfg",
		Short: "mazzaroth cli configurations and preferences",
	}

	cfgRootCmd.AddCommand(
		set(),
		add())
	return cfgRootCmd
}
