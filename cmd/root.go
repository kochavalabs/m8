package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/kochavalabs/mazzaroth-cli/internal/cfg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func Execute() error {
	// root command entry to application
	rootCmd := &cobra.Command{
		Use:   "m8",
		Short: "mazzaroth command line interface",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Bind Cobra flags with viper
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			// Environment variables are expected to be ALL CAPS
			viper.AutomaticEnv()
			viper.SetEnvPrefix("m8")

			if _, err := os.Stat(viper.GetString(cfgPath)); errors.Is(err, os.ErrNotExist) {
				return errors.New(err.Error())
			}

			/*
				viper.SetConfigType("yaml")
				viper.SetConfigFile(viper.GetString(cfgPath))
				err := viper.ReadInConfig()
				if err != nil {
					return err
				}
				fmt.Println(viper.AllKeys())
			*/

			config, err := cfg.FromFile(viper.GetString(cfgPath))
			if err != nil {
				return err
			}

			viper.Set("cfg", config)

			// Set flag values from cfg if not set
			if !viper.IsSet(channelId) {
				channel, err := config.ActiveChannel()
				if err != nil {
					return errors.New("missing required channel ID")
				}
				viper.Set(channelId, channel.ChannelID)
			}

			if !viper.IsSet(channelAddress) {
				channel, err := config.ActiveChannel()
				if err != nil {
					return err
				}
				viper.Set(channelAddress, channel.ChannelAddress)
			}

			if !viper.IsSet(privateKey) {
				viper.Set(privateKey, config.User.PrivateKey)
			}

			if !viper.IsSet(publicKey) {
				viper.Set(publicKey, config.User.PublicKey)
			}

			return nil
		},
	}

	rootCmd.AddCommand(channelCmdChain())
	rootCmd.AddCommand(configurationCmdChain())

	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	rootCmd.PersistentFlags().String(cfgPath, dir+cfgDir+cfgName, "location of the mazzaroth config file")

	ctx, cancel := context.WithCancel(context.Background())
	errGrp, errctx := errgroup.WithContext(ctx)
	errGrp.Go(func() error {
		defer cancel()
		if err := rootCmd.ExecuteContext(errctx); err != nil {
			return err
		}
		return nil
	})
	return errGrp.Wait()
}
