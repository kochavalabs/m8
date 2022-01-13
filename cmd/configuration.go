package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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
			// this is here to overwrite the root persistentPreRunE
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCfg := &cfg.Configuration{
				User: &cfg.UserCfg{},
			}

			// default config location
			dirname, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			// configuration directory prompt
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
				fmt.Println("here")
				return err
			}

			// check if existing configuration exists
			if _, err := os.Stat(directory + cfgDir + cfgName); !errors.Is(err, os.ErrNotExist) {
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
				Label:     "Generate Key Pair",
				Default:   "y",
				IsConfirm: true,
			}
			genKey, err := keyGenPrompt.Run()
			if err != nil {
				// prompt will return an error when confirmation is No
				// log error in case there are other critical errors during prompt execution
			}

			switch strings.ToLower(genKey) {
			case "n":
				// prompt to set keys
				privKeyPrompt := promptui.Prompt{
					Label: "Add Private Key",
					Validate: func(input string) error {
						// TODO :: add private key format check
						if len(input) < 128 {
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
				cliCfg.Channels = channels
				cliCfg.User.ActiveChannel = channelAlias
			}

			cliCfg.Version = version

			d, err := yaml.Marshal(&cliCfg)
			if err != nil {
				return err
			}

			if _, err := os.Stat(directory + cfgDir); errors.Is(err, os.ErrNotExist) {
				if err := os.Mkdir(directory+cfgDir, 0755); err != nil {
					return err
				}
			}

			if err := ioutil.WriteFile(directory+cfgDir+cfgName, d, 0644); err != nil {
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
			cfgYaml, err := yaml.Marshal(cfg)
			if err != nil {
				return err
			}
			fmt.Println(string(cfgYaml))
			return nil
		},
	}
	cfgShowCmd.Flags().String(cfgPath, "$HOME/.m8/cfg.yaml", "path to m8 cfg file")

	cfgRootCmd.AddCommand(cfgInitCmd, cfgShowCmd)
	return cfgRootCmd
}
