package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/elewis787/boa"
	"github.com/kochavalabs/m8/cmd/channel"
	"github.com/kochavalabs/m8/cmd/config"
	"github.com/kochavalabs/m8/internal/cfg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func Execute() error {
	// root command entry to application
	rootCmd := &cobra.Command{
		Use:     "m8",
		Version: "v0.0.1",
		Aliases: []string{"mazzaroth"},
		Short:   "mazzaroth command line interface",
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

	rootCmd.AddCommand(
		initialize(),
		show(),
		pause(),
		delete(),
		channel.ChannelCmdChain(),
		config.ConfigurationCmdChain())

	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	styles := boa.DefaultStyles()
	styles.Title.BorderForeground(lipgloss.AdaptiveColor{Light: `#E3BD2D`, Dark: `#E3BD2D`})
	styles.Border.BorderForeground(lipgloss.AdaptiveColor{Light: `#E3BD2D`, Dark: `#E3BD2D`})
	styles.SelectedItem.Foreground(lipgloss.AdaptiveColor{Light: `#353C3B`, Dark: `#353C3B`}).
		Background(lipgloss.AdaptiveColor{Light: `#E3BD2D`, Dark: `#E3BD2D`})
	b := boa.New(boa.WithAltScreen(true), boa.WithStyles(styles))

	rootCmd.SetHelpFunc(b.HelpFunc)
	rootCmd.SetUsageFunc(b.UsageFunc)
	rootCmd.PersistentFlags().String(cfgPath, dir+cfgDir+cfgName, "location of the mazzaroth config file")
	rootCmd.PersistentFlags().String(channelId, "", "defaults to the active channel id in the cfg")
	rootCmd.PersistentFlags().String(channelAddress, "", "defaults to active channel address in the cfg")

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
