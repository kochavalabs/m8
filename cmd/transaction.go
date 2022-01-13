package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/kochavalabs/mazzaroth-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func transactionCmdChain() *cobra.Command {
	transactionRootCmd := &cobra.Command{
		Use:   "tx",
		Short: "interact with receipt endpoints on a mazzaroth gateway node",
	}

	transactionLookupCmd := &cobra.Command{
		Use:   "lookup",
		Short: "lookup a transaction for a given channel by transaction id",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(address)))
			if err != nil {
				return err
			}

			receipt, err := client.TransactionLookup(cmd.Context(), viper.GetString(channelId), viper.GetString(transactionid))
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
	transactionLookupCmd.Flags().String(transactionid, "", "id of the transaction being looked up")
	transactionLookupCmd.MarkFlagRequired(transactionid)

	transactionCallCmd := &cobra.Command{
		Use:   "call",
		Short: "call functions on a mazzaroth channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			/*client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(address)))
			if err != nil {
				return err
			}

			sender, err := xdr.IDFromHexString(viper.GetString("pub-key"))
			if err != nil {
				return err
			}

			channelId, err := xdr.IDFromHexString(viper.GetString("active-channel"))
			if err != nil {
				return err
			}

			tx := mazzaroth.Transaction(sender, channelId)
			*/
			return nil
		},
	}
	transactionCallCmd.Flags().String(function, "", "the function to be called")
	//	transactionCallCmd.MarkFlagRequired(function)

	transactionCallCmd.Flags().StringSlice(args, []string{""}, "the args to pass within the function")

	transactionRootCmd.AddCommand(transactionLookupCmd, transactionCallCmd)
	return transactionRootCmd
}
