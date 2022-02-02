package channel

import (
	"github.com/spf13/cobra"
)

func ChannelCmdChain() *cobra.Command {
	channelRootCmd := &cobra.Command{
		Use:   "channel",
		Short: "interact with channel endpoints on a mazzaroth gateway node",
	}

	channelRootCmd.AddCommand(
		lookup(),
		list(),
		exec())

	return channelRootCmd
}
