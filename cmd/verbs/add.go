package verbs

import (
	"errors"

	"github.com/kochavalabs/crypto"
	"github.com/kochavalabs/m8/internal/cfg"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultChannelId          = `0000000000000000000000000000000000000000000000000000000000000000`
	defaultGatewayNodeAddress = `http://localhost:6299`
	channelIdLength           = 32
)

func Add() *cobra.Command {
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

			channelCfg, err := channelPrompt()
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

func channelPrompt() (*cfg.ChannelCfg, error) {
	channelAliasPrompt := promptui.Prompt{
		Label:   "Channel Alias",
		Default: "default-channel",
	}
	alias, err := channelAliasPrompt.Run()
	if err != nil {
		return nil, err
	}
	channelIDPrompt := promptui.Prompt{
		Label: "Channel id",
		Validate: func(input string) error {
			pub, err := crypto.FromHex(input)
			if err != nil {
				return err
			}
			if len(pub) != channelIdLength {
				return errors.New("invalid channel id length")
			}
			return nil
		},
		Default: defaultChannelId,
	}
	id, err := channelIDPrompt.Run()
	if err != nil {
		return nil, err
	}

	channelAddrPrompt := promptui.Prompt{
		Label:   "Channel Address",
		Default: defaultGatewayNodeAddress,
	}

	addr, err := channelAddrPrompt.Run()
	if err != nil {
		return nil, err
	}

	channelCfg := &cfg.ChannelCfg{
		Channel: &cfg.Channel{
			ChannelAddress: addr,
			ChannelID:      id,
			ChannelAlias:   alias,
		},
	}
	return channelCfg, nil
}
