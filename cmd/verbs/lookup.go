package verbs

import (
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kochavalabs/m8/internal/tui"
	"github.com/kochavalabs/mazzaroth-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// lookup
// blockheight
// block
// tx
// receipt

const (
	gold     = `#E3BD2D`
	darkGrey = `#353C3B`
	teal     = `#01A299`
	white    = `#FFFFFF`

	transactionId  = `tx-id`
	channelAddress = `channel-address`
	channelId      = `channel-id`
	header         = `header`
	blockid        = `block-id`
)

func Lookup(resource string) *cobra.Command {
	lookup := &cobra.Command{
		Use:   "lookup",
		Short: "look up items on a mazzaroth node",
	}
	// sub command chain by resource type
	switch resource {
	case "channel":
		lookup.AddCommand(lookupAbi(), lookupBlock(), lookupTx(), lookupReceipt())
		return lookup
	default:
		return lookup
	}
}

func lookupAbi() *cobra.Command {
	abi := &cobra.Command{
		Use:   "abi",
		Short: "return the application binary interfacem (ABI) for a channel",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}

			abi, err := client.ChannelAbi(cmd.Context(), viper.GetString(channelId))
			if err != nil {
				return err
			}

			v, err := json.MarshalIndent(abi, "", " ")
			if err != nil {
				return err
			}

			fmt.Println(string(v))
			return nil
		},
	}
	return abi
}

func lookupBlockHeight() *cobra.Command {

	blockHeight := &cobra.Command{
		Use:   "height",
		Short: "return the block height of a given channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}

			height, err := client.BlockHeight(cmd.Context(), viper.GetString(channelId))
			if err != nil {
				return err
			}

			v, err := json.MarshalIndent(height, "", " ")
			if err != nil {
				return err
			}

			fmt.Println(string(v))
			return nil
		},
	}
	return blockHeight
}

func lookupBlock() *cobra.Command {
	block := &cobra.Command{
		Use:   "block",
		Short: "look up items on a mazzaroth node",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}

			switch {
			// block header lookup
			case viper.GetBool(header):
				blockheader, err := client.BlockHeaderLookup(cmd.Context(), viper.GetString(channelId), viper.GetString(blockid))
				if err != nil {
					return err
				}

				v, err := json.MarshalIndent(blockheader, "", " ")
				if err != nil {
					return err
				}

				fmt.Println(string(v))
				return nil
			// block lookup
			default:

				block, err := client.BlockLookup(cmd.Context(), viper.GetString(channelId), viper.GetString(blockid))
				if err != nil {
					return err
				}

				v, err := json.MarshalIndent(block, "", " ")
				if err != nil {
					return err
				}

				fmt.Println(string(v))
				return nil
			}
		},
	}
	block.Flags().Bool(header, false, "option to return block header")
	block.Flags().String(blockid, "", "id of block")
	block.MarkFlagRequired(blockid)
	return block
}

func lookupTx() *cobra.Command {
	txLookup := &cobra.Command{
		Use:   "tx",
		Short: "look up items on a mazzaroth node",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr := viper.GetString(channelAddress)

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(addr))
			if err != nil {
				return err
			}

			channelId := viper.GetString(channelId)
			transactionId := viper.GetString(transactionId)

			txCmd := tui.TxLookup(cmd.Context(), client, channelId, transactionId)
			txModel := tui.NewTxModel(txCmd)

			if err := tea.NewProgram(txModel).Start(); err != nil {
				return err
			}
			return nil
		},
	}
	txLookup.Flags().String(transactionId, "", "id of the transaction being looked up")
	txLookup.MarkFlagRequired(transactionId)
	return txLookup
}

func lookupReceipt() *cobra.Command {
	rcpt := &cobra.Command{
		Use:   "rcpt",
		Short: "lookup a receipt for a given channel by transaction id",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}

			channelId := viper.GetString(channelId)
			transactionId := viper.GetString(transactionId)

			rcptCmd := tui.RcptLookup(cmd.Context(), client, channelId, transactionId)
			rcptModel := tui.NewRcptModel(rcptCmd)

			if err := tea.NewProgram(rcptModel).Start(); err != nil {
				return err
			}
			return nil
		},
	}
	rcpt.Flags().String(transactionId, "", "transaction id assoicated to the receipt being looked up")
	rcpt.MarkFlagRequired(transactionId)

	return rcpt
}
