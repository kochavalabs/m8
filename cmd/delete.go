package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kochavalabs/crypto"
	"github.com/kochavalabs/m8/internal/tui"
	"github.com/kochavalabs/mazzaroth-go"
	"github.com/kochavalabs/mazzaroth-xdr/go-xdr/xdr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func delete() *cobra.Command {
	delete := &cobra.Command{
		Use:   "delete",
		Short: "delete resources on a mazzaroth node",
	}
	delete.AddCommand(deleteChannel())
	return delete
}

func deleteChannel() *cobra.Command {
	channel := &cobra.Command{
		Use:   "channel",
		Short: "delete a channel contract",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}

			sender, err := xdr.IDFromHexString(viper.GetString(publicKey))
			if err != nil {
				return err
			}

			cId, err := xdr.IDFromHexString(viper.GetString(channelId))
			if err != nil {
				return err
			}

			pk, err := crypto.FromHex(viper.GetString(privateKey))
			if err != nil {
				return err
			}

			blockHeight, err := client.BlockHeight(cmd.Context(), viper.GetString(channelId))
			if err != nil {
				return err
			}

			tx, err := mazzaroth.Transaction(sender, cId).
				Contract(mazzaroth.GenerateNonce(), blockHeight.Height+maxBlockExpirationRange).
				Delete().Sign(pk)
			if err != nil {
				return err
			}

			channelCmd := tui.ChannelDelete(cmd.Context(), client, tx)
			channelModel := tui.NewChannelModel(channelCmd)

			if err := tea.NewProgram(channelModel).Start(); err != nil {
				return err
			}
			return nil
		},
	}

	return channel
}
