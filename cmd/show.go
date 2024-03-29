package cmd

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kochavalabs/m8/internal/cfg"
	"github.com/kochavalabs/m8/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func show() *cobra.Command {
	show := &cobra.Command{
		Use:   "show",
		Short: "show resources",
	}
	show.AddCommand(showCfg())
	return show
}

func showCfg() *cobra.Command {
	showCfg := &cobra.Command{
		Use:   "cfg",
		Short: "show the current m8 configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := viper.Get("cfg").(*cfg.Configuration)
			if cfg == nil {
				return errors.New("missing configuration")
			}
			cfgCmd := tui.CfgShow(cfg, viper.GetString(cfgPath))
			cfgModel := tui.NewCfgModel(cfgCmd)

			if err := tea.NewProgram(cfgModel).Start(); err != nil {
				return err
			}
			return nil
		},
	}
	return showCfg
}
