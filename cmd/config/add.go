package config

import (
	"errors"

	"github.com/kochavalabs/m8/internal/cfg"
	"github.com/kochavalabs/m8/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultChannelId          = `0000000000000000000000000000000000000000000000000000000000000000`
	defaultGatewayNodeAddress = `http://localhost:6299`
	channelIdLength           = 32
)

func add() *cobra.Command {
	add := &cobra.Command{
		Use:   "add",
		Short: "add items to mazzaroth resources",
	}
	add.AddCommand(addChannel())
	return add
}

func addChannel() *cobra.Command {
	channel := &cobra.Command{
		Use:   "channel",
		Short: "add a channel to the mazzaroth cli cfg",
		RunE: func(cmd *cobra.Command, args []string) error {
			var config *cfg.Configuration
			v := viper.Get("cfg")
			if v != nil {
				config = v.(*cfg.Configuration)
			} else {
				config = &cfg.Configuration{}
			}

			channelCfg, err := tui.ChannelPrompt()
			if err != nil {
				return err
			}

			if config.ContainsChannel(channelCfg.Channel.ChannelID, channelCfg.Channel.ChannelAlias) {
				return errors.New("channel already exists with the same id or alias")
			}

			config.Channels = append(config.Channels, channelCfg)
			cfg.ToFile(viper.GetString(cfgPath), config)
			return nil
		},
	}
	return channel
}
