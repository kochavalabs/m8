package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/kochavalabs/mazzaroth-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func channelCmdChain() *cobra.Command {

	channelRootCmd := &cobra.Command{
		Use:   "channel",
		Short: "interact with channel endpoints on a mazzaroth gateway node",
	}

	channelAbiCmd := &cobra.Command{
		Use:   "abi",
		Short: "return the application binary interface for a channel",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(address)))
			if err != nil {
				return err
			}

			abi, err := client.ChannelAbi(cmd.Context(), viper.GetString(channelId))
			if err != nil {
				return err
			}

			v, err := json.MarshalIndent(abi, "", "\t")
			if err != nil {
				return err
			}

			fmt.Println(string(v))
			return nil
		},
	}

	channelConfigCmd := &cobra.Command{
		Use:   "config",
		Short: "return the configuration for a channel",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := mazzaroth.NewMazzarothClient(mazzaroth.WithAddress(viper.GetString(address)))
			if err != nil {
				return err
			}

			config, err := client.ChannelConfig(cmd.Context(), viper.GetString(channelId))
			if err != nil {
				return err
			}

			v, err := json.MarshalIndent(config, "", "\t")
			if err != nil {
				return err
			}

			fmt.Println(string(v))
			return nil
		},
	}

	channelRootCmd.AddCommand(channelAbiCmd, channelConfigCmd)
	return channelRootCmd
}
