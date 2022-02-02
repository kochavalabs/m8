package resources

import (
	"errors"

	"github.com/kochavalabs/crypto"
	"github.com/kochavalabs/m8/cmd/verbs"
	"github.com/kochavalabs/m8/internal/cfg"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func ConfigurationCmdChain() *cobra.Command {
	cfgRootCmd := &cobra.Command{
		Use:   "cfg",
		Short: "mazzaroth cli configurations and preferences",
	}

	cfgRootCmd.AddCommand(
		verbs.Set(),
		verbs.Add())
	return cfgRootCmd
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
