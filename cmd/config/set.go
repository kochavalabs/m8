package config

import (
	"errors"

	"github.com/kochavalabs/m8/internal/cfg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	channelAlias = `channel-alias`
	cfgPath      = `cfg-path`
)

func set() *cobra.Command {
	set := &cobra.Command{
		Use:   "set",
		Short: "set values in the mazzaroth config",
	}
	set.AddCommand(setActiveChannel())
	return set
}

func setActiveChannel() *cobra.Command {
	setActiveChannel := &cobra.Command{
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
	setActiveChannel.Flags().String(channelAlias, "", "channel alias for the channel to set as active")
	setActiveChannel.MarkFlagRequired(channelAlias)
	return setActiveChannel
}
