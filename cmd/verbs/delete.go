package verbs

import (
	"encoding/hex"
	"fmt"

	"github.com/kochavalabs/crypto"
	"github.com/kochavalabs/mazzaroth-go"
	"github.com/kochavalabs/mazzaroth-xdr/go-xdr/xdr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Delete(resource string) *cobra.Command {
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

			channelId, err := xdr.IDFromHexString(viper.GetString(channelId))
			if err != nil {
				return err
			}

			pk, err := crypto.FromHex(viper.GetString(privateKey))
			if err != nil {
				return err
			}

			tx, err := mazzaroth.Transaction(sender, channelId).Contract(mazzaroth.GenerateNonce(), maxBlockExpirationRange).Delete().Sign(pk)
			if err != nil {
				return err
			}

			id, _, err := client.TransactionSubmit(cmd.Context(), tx)
			if err != nil {
				return err
			}

			fmt.Println("transaction id:", hex.EncodeToString(id[:]))
			return nil
		},
	}

	return channel
}
