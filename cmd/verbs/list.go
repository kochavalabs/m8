package verbs

import (
	"encoding/json"
	"fmt"

	"github.com/kochavalabs/mazzaroth-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	headers = `headers`
	height  = `height`
	number  = `number`
)

func List() *cobra.Command {
	list := &cobra.Command{
		Use:   "list",
		Short: "list items for a given mazzaroth channel",
	}
	list.AddCommand(listBlocks())
	return list
}

func listBlocks() *cobra.Command {
	blocks := &cobra.Command{
		Use:   "blocks",
		Short: "list blocks or block headers for a given channel at a specific height",
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
	blocks.Flags().Bool(headers, false, "option to return list of block headers")
	blocks.Flags().Int(height, 0, "starting block height value")
	blocks.MarkFlagRequired(height)
	blocks.Flags().Int(number, 1, "number of blocks to list")
	blocks.MarkFlagRequired(number)
	return blocks
}
