package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/kochavalabs/crypto"
	"github.com/kochavalabs/m8/internal/cfg"
	"github.com/kochavalabs/m8/internal/tui"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initialize() *cobra.Command {
	init := &cobra.Command{
		Use:   "init",
		Short: "initialize resources",
	}
	init.AddCommand(initCfg())
	return init
}

func initCfg() *cobra.Command {
	initCfg := &cobra.Command{
		Use:   "cfg",
		Short: "initialize the mazzaroth cli configuration and preferences",
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
				channelCfg, err := tui.ChannelPrompt()
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
	return initCfg
}

func initChannel() *cobra.Command {
	initChannel := &cobra.Command{
		Use:   "channel",
		Short: "initialize a mazzaroth channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Generate Key
			pub, priv, err := crypto.GenerateEd25519KeyPair()
			if err != nil {
				return err
			}
			fmt.Println("channel address:", crypto.ToHex(pub))
			fmt.Println("channel private key:", crypto.ToHex(priv))
			fmt.Println("channel Id:", crypto.ToHex(pub))
			// TODO
			// self signed cert for channel
			// mazzaroth.io cert generate for channel
			return nil
		},
	}
	return initChannel
}
