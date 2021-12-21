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
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const (
	configFileName = "/.mazzaroth-cli"
)

func configure() *cobra.Command {
	config := &cobra.Command{
		Use:   "cfg",
		Short: "Setup of the mazzaroth cli configurations and preferences",
	}

	initCfg := &cobra.Command{
		Use:   "init",
		Short: "Initalize the mazzaroth cli configuration and preferences",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCfg := &cfg.Configuration{}
			// Pull configured logger from viper registery
			zlogger := viper.Get("logger").(*zap.Logger)
			if zlogger == nil {
				// if not logger was found, create a no-op logger
				zlogger = zap.NewNop()
			}

			// Default config location $HOME
			dirname, err := os.UserHomeDir()
			if err != nil {
				zlogger.Debug(err.Error())
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
			zlogger.Debug(directory)

			// Check if existing configuration exists
			if _, err := os.Stat(directory + configFileName); !errors.Is(err, os.ErrNotExist) {
				zlogger.Debug("file exists")
				overwriteExistingPrompt := promptui.Prompt{
					Label:     "Overwrite existing config at " + directory,
					Default:   "n",
					IsConfirm: true,
				}
				if v, err := overwriteExistingPrompt.Run(); err != nil || strings.ToLower(v) == "n" {
					zlogger.Debug("exit to avoid config overwriting file")
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
				zlogger.Debug(err.Error())
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
				cliCfg.PrivateKey = priv
				cliCfg.PublicKey = pub
			default: // default case is y
				// Generate Key
				priv, pub, err := crypto.GenerateEd25519KeyPair()
				if err != nil {
					return err
				}
				cliCfg.PrivateKey = crypto.ToHex(priv)
				cliCfg.PublicKey = crypto.ToHex(pub)
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
				zlogger.Debug("error" + err.Error())
			}
			switch strings.ToLower(addChannel) {
			case "n":
				zlogger.Debug("no channel to configure")
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
					Default: "0x00000000000000000000000000000000000",
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

				channelCfg := cfg.ChannelConfiguration{
					ChannelURL:   channelUrl,
					ChannelID:    channelID,
					ChannelAlias: channelAlias,
				}
				channels := make([]cfg.ChannelConfiguration, 0, 0)
				channels = append(channels, channelCfg)
				cliCfg.Channels = channels
				cliCfg.ActiveChannel = channelAlias
			}

			d, err := yaml.Marshal(&cliCfg)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			if err := ioutil.WriteFile(directory+configFileName, d, 0644); err != nil {
				return err
			}

			zlogger.Debug("done")
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

	config.AddCommand(initCfg, cfgChannel)
	return config
}
