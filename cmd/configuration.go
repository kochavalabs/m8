package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/kochavalabs/crypto"
	"github.com/kochavalabs/mazzaroth-cli/internal/cfg"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func configurationCmdChain() *cobra.Command {
	cfgRootCmd := &cobra.Command{
		Use:   "cfg",
		Short: "mazzaroth cli configurations and preferences",
	}

	cfgInitCmd := &cobra.Command{
		Use:   "init",
		Short: "Initalize the mazzaroth cli configuration and preferences",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Bind Cobra flags with viper
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			// Environment variables are expected to be ALL CAPS
			viper.AutomaticEnv()
			viper.SetEnvPrefix("m8")
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCfg := &cfg.Configuration{
				Version:  version,
				User:     &cfg.UserCfg{},
				Channels: make([]*cfg.ChannelCfg, 0, 0),
			}

			// configuration directory prompt
			cfgDirPrompt := promptui.Prompt{
				Label: "Configuration directory",
				Validate: func(input string) error {
					p := path.Dir(input)
					if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
						return err
					}
					return nil
				},
				Default: path.Dir(viper.GetString(cfgPath)),
			}

			directory, err := cfgDirPrompt.Run()
			if err != nil {
				return err
			}

			// check if existing configuration exists
			if _, err := os.Stat(viper.GetString(cfgPath)); !errors.Is(err, os.ErrNotExist) {
				overwriteExistingPrompt := promptui.Prompt{
					Label:     "Overwrite existing config at " + directory,
					Default:   "n",
					IsConfirm: true,
				}
				if v, err := overwriteExistingPrompt.Run(); err != nil || strings.ToLower(v) == "n" {
					return nil
				}
			}

			// key generation
			keyGenPrompt := promptui.Prompt{
				Label:     "Generate key pair",
				Default:   "y",
				IsConfirm: true,
			}
			genKey, _ := keyGenPrompt.Run()

			switch strings.ToLower(genKey) {
			case "n":
				// prompt to set keys
				privKeyPrompt := promptui.Prompt{
					Label: "Add private key",
					Validate: func(input string) error {
						pub, err := crypto.FromHex(input)
						if err != nil {
							return err
						}
						if len(pub) != privKeylength {
							return errors.New("invalid private key length")
						}
						return nil
					},
					HideEntered: true,
				}
				priv, err := privKeyPrompt.Run()
				if err != nil {
					return err
				}

				privKey, err := crypto.FromHex(priv)
				if err != nil {
					return err
				}
				pubKey := privKey[32:]

				pubKeyPrompt := promptui.Prompt{
					Label: "Add public key",
					Validate: func(input string) error {
						pub, err := crypto.FromHex(input)
						if err != nil {
							return err
						}
						if len(pub) != pubKeyLength {
							return errors.New("invalid public key length")
						}
						return nil
					},
					Default: crypto.ToHex(pubKey),
				}
				pub, err := pubKeyPrompt.Run()
				if err != nil {
					return err
				}
				cliCfg.User.PrivateKey = priv
				cliCfg.User.PublicKey = pub
			default: // default case is y
				// Generate Key
				pub, priv, err := crypto.GenerateEd25519KeyPair()
				if err != nil {
					return err
				}
				cliCfg.User.PrivateKey = crypto.ToHex(priv)
				cliCfg.User.PublicKey = crypto.ToHex(pub)
			}

			// Add Channel Prompt
			addChannelPrompt := promptui.Prompt{
				Label:     "Setup a channel",
				Default:   "y",
				IsConfirm: true,
			}

			addChannel, err := addChannelPrompt.Run()
			if err != nil {
				// Prompt will return an error when confirmation is No
				// log error in case there are other critical errors during prompt execution
			}

			switch strings.ToLower(addChannel) {
			case "n":
			default: // default case is y
				channelCfg, err := channelPrompt()
				if err != nil {
					return err
				}

				cliCfg.Channels = append(cliCfg.Channels, channelCfg)
				cliCfg.User.ActiveChannel = channelCfg.Channel.ChannelAlias
			}

			if err := cfg.ToFile(viper.GetString(cfgPath), cliCfg); err != nil {
				fmt.Println("here")
				return err
			}

			return nil
		},
	}

	cfgShowCmd := &cobra.Command{
		Use:   "show",
		Short: "display the current cfg file",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := viper.Get("cfg").(*cfg.Configuration)
			if cfg == nil {
				return errors.New("missing configuration")
			}
			cfgYaml, err := yaml.Marshal(cfg)
			if err != nil {
				return err
			}
			fmt.Println(string(cfgYaml))
			return nil
		},
	}

	cfgSetCmd := &cobra.Command{
		Use:   "set",
		Short: "set or updates values in the mazzaroth cfg",
	}

	cfgSetChannelCmd := &cobra.Command{
		Use:   "channel",
		Short: "sets the active channel in the mazzaroth cfg",
		RunE: func(cmd *cobra.Command, args []string) error {
			var config *cfg.Configuration
			v := viper.Get("cfg")
			if v != nil {
				config = v.(*cfg.Configuration)
			} else {
				config = &cfg.Configuration{}
			}
			if ok := config.ContainsChannel("", viper.GetString(channelAlias)); !ok {
				return errors.New("no channel with the supplied channel alias found")
			}
			config.User.ActiveChannel = viper.GetString(channelAlias)

			cfg.ToFile(viper.GetString(cfgPath), config)

			return nil
		},
	}
	cfgSetChannelCmd.Flags().String(channelAlias, "", "channel alias for the channel to set as active")
	cfgSetChannelCmd.MarkFlagRequired(channelAlias)

	cfgAddCmd := &cobra.Command{
		Use:   "add",
		Short: "add elements to the mazzaroth cli cfg",
	}

	cfgAddChannelCmd := &cobra.Command{
		Use:   "channel",
		Short: "add a channel to the mazzaroth cli",
		RunE: func(cmd *cobra.Command, args []string) error {
			var config *cfg.Configuration
			v := viper.Get("cfg")
			if v != nil {
				config = v.(*cfg.Configuration)
			} else {
				config = &cfg.Configuration{}
			}
			fmt.Println(config)
			fmt.Println(config.User.ActiveChannel)
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

	cfgSetCmd.AddCommand(cfgSetChannelCmd)
	cfgAddCmd.AddCommand(cfgAddChannelCmd)
	cfgRootCmd.AddCommand(cfgInitCmd, cfgShowCmd, cfgSetCmd, cfgAddCmd)
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
