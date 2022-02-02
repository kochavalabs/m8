package resources

import (
	"github.com/kochavalabs/m8/cmd/verbs"
	"github.com/spf13/cobra"
)

func ChannelCmdChain() *cobra.Command {
	channelRootCmd := &cobra.Command{
		Use:   "channel",
		Short: "interact with channel endpoints on a mazzaroth gateway node",
	}
	channelRootCmd.PersistentFlags().String(channelId, "", "defaults to the active channel id in the cfg")
	channelRootCmd.PersistentFlags().String(channelAddress, "", "defaults to active channel address in the cfg")

	channelRootCmd.AddCommand(
		verbs.Lookup("channel"),
		verbs.Exec(),
		verbs.Pause())

	return channelRootCmd
}
