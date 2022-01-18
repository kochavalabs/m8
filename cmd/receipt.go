package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/kochavalabs/mazzaroth-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func receiptCmdChain() *cobra.Command {
	receiptRootCmd := &cobra.Command{
		Use:   "receipt",
		Short: "interact with receipt endpoints on a mazzaroth gateway node",
	}

	receiptLookupCmd := &cobra.Command{
		Use:   "lookup",
		Short: "lookup a receipt for a given channel by transaction id",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}

			receipt, err := client.ReceiptLookup(cmd.Context(), viper.GetString(channelId), viper.GetString(transactionId))
			if err != nil {
				return err
			}

			v, err := json.MarshalIndent(receipt, "", "\t")
			if err != nil {
				return err
			}

			fmt.Println(string(v))
			return nil
		},
	}
	receiptLookupCmd.Flags().String(transactionId, "", "transaction id assoicated to the receipt being looked up")
	receiptLookupCmd.MarkFlagRequired(transactionId)

	receiptRootCmd.AddCommand(receiptLookupCmd)
	return receiptRootCmd
}
