package cmd

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/kochavalabs/crypto"
	"github.com/kochavalabs/mazzaroth-cli/internal/cfg"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	configFileName = `/.m8`
	version        = `0.0.1`
)

func configurationCmdChain() *cobra.Command {
	cfgRootCmd := &cobra.Command{
		Use:   "cfg",
		Short: "returns mazzaroth cli configurations and preferences",
		RunE: func(cmd *cobra.Command, args []string) error {
			// todo
			return nil
		},
	}

	cfgInitCmd := &cobra.Command{
		Use:   "init",
		Short: "Initalize the mazzaroth cli configuration and preferences",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCfg := &cfg.Configuration{
				User: &cfg.UserCfg{},
			}

			// Default config location $HOME
			dirname, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			// Configuration Directory Prompt
			cfgDirPrompt := promptui.Prompt{
				Label: "Configuration Directory",
				Validate: func(input string) error {
					if _, err := os.Stat(input); errors.Is(err, os.ErrNotExist) {
						return err
					}
					return nil
				},
				Default: dirname,
			}
			directory, err := cfgDirPrompt.Run()
			if err != nil {
				return err
			}

			// Check if existing configuration exists
			if _, err := os.Stat(directory + configFileName); !errors.Is(err, os.ErrNotExist) {
				overwriteExistingPrompt := promptui.Prompt{
					Label:     "Overwrite existing config at " + directory,
					Default:   "n",
					IsConfirm: true,
				}
				if v, err := overwriteExistingPrompt.Run(); err != nil || strings.ToLower(v) == "n" {
					return nil
				}
			}

			// Key Generation
			keyGenPrompt := promptui.Prompt{
				Label:     "Generate Key Pair",
				Default:   "y",
				IsConfirm: true,
			}
			genKey, err := keyGenPrompt.Run()
			if err != nil {
				// Prompt will return an error when confirmation is No
				// log error in case there are other critical errors during prompt execution
			}

			switch strings.ToLower(genKey) {
			case "n":
				// Prompt to set keys
				privKeyPrompt := promptui.Prompt{
					Label: "Add Private Key",
					Validate: func(input string) error {
						// TODO :: Add private key format check
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

				// Generate a Public key for a default value from private key supplied
				pubKey, err := crypto.Ed25519PublicKeyFromPrivate(privKey)
				if err != nil {
					return err
				}

				pubKeyPrompt := promptui.Prompt{
					Label: "Add Public Key",
					Validate: func(input string) error {
						// TODO :: Add public key format check
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
				priv, pub, err := crypto.GenerateEd25519KeyPair()
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
				channelAliasPrompt := promptui.Prompt{
					Label: "Channel Alias",
					Validate: func(input string) error {
						// TODO :: Add public key format check
						return nil
					},
					Default: "my-channel",
				}
				channelAlias, err := channelAliasPrompt.Run()
				if err != nil {
					return err
				}
				channelIDPrompt := promptui.Prompt{
					Label: "Channel Id",
					Validate: func(input string) error {
						// TODO :: Add public key format check
						return nil
					},
					Default: "00000000000000000000000000000000000",
				}
				channelID, err := channelIDPrompt.Run()
				if err != nil {
					return err
				}

				channelURLPrompt := promptui.Prompt{
					Label: "Channel Url",
					Validate: func(input string) error {
						// TODO :: Add public key format check
						return nil
					},
					Default: "http://localhost:8080",
				}
				channelUrl, err := channelURLPrompt.Run()
				if err != nil {
					return err
				}

				channelCfg := cfg.ChannelCfg{
					Channel: cfg.Channel{
						ChannelURL:   channelUrl,
						ChannelID:    channelID,
						ChannelAlias: channelAlias,
					},
				}

				channels := make([]cfg.ChannelCfg, 0, 0)
				channels = append(channels, channelCfg)
				channels = append(channels, channelCfg)
				cliCfg.Channels = channels
				cliCfg.User.ActiveChannel = channelAlias
			}

			cliCfg.Version = version

			d, err := yaml.Marshal(&cliCfg)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			if err := ioutil.WriteFile(directory+configFileName, d, 0644); err != nil {
				return err
			}

			return nil
		},
	}

	cfgChannel := &cobra.Command{
		Use:   "channel",
		Short: "Configure a channel connection",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cfgRootCmd.AddCommand(cfgInitCmd, cfgChannel)
	return cfgRootCmd
}
