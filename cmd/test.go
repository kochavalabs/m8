package cmd

import "github.com/spf13/cobra"

func testCmdChain() *cobra.Command {

	testRootCmd := &cobra.Command{
		Use:   "test",
		Short: "test channel contracts on mazzaroth nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load Manifiest file here
			return nil
		},
	}

	return testRootCmd
}
