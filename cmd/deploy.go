package cmd

import "github.com/spf13/cobra"

func deployCmdChain() *cobra.Command {

	deployRootCmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy channel contracts to mazzaroth nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load Manifiest file here
			return nil
		},
	}

	return deployRootCmd
}
