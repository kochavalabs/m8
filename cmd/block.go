package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/kochavalabs/mazzaroth-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func blockCmdChain() *cobra.Command {

	blockRootCmd := &cobra.Command{
		Use:   "block",
		Short: "interact with block endpoints on a mazzaroth gateway nmode",
	}

	blockHeightCmd := &cobra.Command{
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

			v, err := json.MarshalIndent(height, "", "\t")
			if err != nil {
				return err
			}

			fmt.Println(string(v))
			return nil
		},
	}

	blockListCmd := &cobra.Command{
		Use:   "list",
		Short: "list blocks or block headers for a given channel at a specific height",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}

			switch {
			// block header list
			case viper.GetBool(headers):
				blockheaders, err := client.BlockHeaderList(cmd.Context(),
					viper.GetString(channelId), viper.GetInt(height), viper.GetInt(number))
				if err != nil {
					return err
				}

				v, err := json.MarshalIndent(blockheaders, "", "\t")
				if err != nil {
					return err
				}

				fmt.Println(string(v))
				return nil
			// block list
			default:

				blocks, err := client.BlockList(cmd.Context(),
					viper.GetString(channelId), viper.GetInt(height), viper.GetInt(number))
				if err != nil {
					return err
				}

				v, err := json.MarshalIndent(blocks, "", "\t")
				if err != nil {
					return err
				}

				fmt.Println(string(v))
				return nil
			}
		},
	}
	blockListCmd.Flags().Bool(headers, false, "option to return list of block headers")
	blockListCmd.Flags().Int(height, 0, "starting block height value")
	blockListCmd.MarkFlagRequired(height)
	blockListCmd.Flags().Int(number, 1, "number of blocks to list")
	blockListCmd.MarkFlagRequired(number)

	blockLookupCmd := &cobra.Command{
		Use:   "lookup",
		Short: "return a block or block header for a given channel by id",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(channelAddress)))
			if err != nil {
				return err
			}

			switch {
			// block header lookup
			case viper.GetBool(headers):
				blockheader, err := client.BlockHeaderLookup(cmd.Context(), viper.GetString(channelId), viper.GetString(blockid))
				if err != nil {
					return err
				}

				v, err := json.MarshalIndent(blockheader, "", "\t")
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

				v, err := json.MarshalIndent(block, "", "\t")
				if err != nil {
					return err
				}

				fmt.Println(string(v))
				return nil
			}
		},
	}
	blockLookupCmd.Flags().Bool(header, false, "option to return block header")
	blockLookupCmd.Flags().String(blockid, "", "id of block")
	blockLookupCmd.MarkFlagRequired(blockid)

	blockRootCmd.AddCommand(blockHeightCmd, blockListCmd, blockLookupCmd)
	return blockRootCmd
}
