package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/kochavalabs/crypto"
	"github.com/kochavalabs/mazzaroth-go"
	"github.com/kochavalabs/mazzaroth-xdr/go-xdr/xdr"
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

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}

			tx, err := client.TransactionLookup(cmd.Context(), viper.GetString(channelId), viper.GetString(transactionId))
			if err != nil {
				return err
			}

			v, err := json.MarshalIndent(tx, "", "\t")
			if err != nil {
				return err
			}

			fmt.Println(string(v))
			return nil
		},
	}
	transactionLookupCmd.Flags().String(transactionId, "", "id of the transaction being looked up")
	transactionLookupCmd.MarkFlagRequired(transactionId)

	transactionCallCmd := &cobra.Command{
		Use:   "call",
		Short: "call functions on a mazzaroth channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}
			pk, err := crypto.FromHex(viper.GetString(privateKey))
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
			xdrArgs := make([]xdr.Argument, 0, 0)
			values := viper.GetStringSlice(arguments)
			for _, a := range values {
				xdrArgs = append(xdrArgs, xdr.Argument(a))
			}

			tx, err := mazzaroth.Transaction(sender, channelId).Call(mazzaroth.GenerateNonce(), defaultBlockExpirationNumber).Function(viper.GetString(function)).Arguments(xdrArgs...).Sign(pk)
			if err != nil {
				return err
			}

			id, rcpt, err := client.TransactionSubmit(cmd.Context(), tx)
			if err != nil {
				return err
			}

			if rcpt != nil {
				jsonRcpt, err := json.MarshalIndent(rcpt, "", "\t")
				if err != nil {
					return err
				}
				fmt.Println("receipt object:")
				fmt.Println(string(jsonRcpt))
			}

			txId := hex.EncodeToString(id[:])
			fmt.Println("transaction id: ", txId)
			return nil
		},
	}
	transactionCallCmd.Flags().String(function, "", "the function to be called")
	transactionCallCmd.MarkFlagRequired(function)

	transactionCallCmd.Flags().StringSlice(arguments, []string{""}, "the args to pass within the function")

	transactionRootCmd.AddCommand(transactionLookupCmd, transactionCallCmd)
	return transactionRootCmd
}
