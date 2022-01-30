package verbs

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/kochavalabs/m8/internal/cfg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func ShowCfg() *cobra.Command {
	showCfg := &cobra.Command{
		Use:   "show",
		Short: "show the current m8 configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := viper.Get("cfg").(*cfg.Configuration)
			if cfg == nil {
				return errors.New("missing configuration")
			}
			cfgYaml, err := yaml.Marshal(cfg)
			if err != nil {
				return err
			}

			barStyle := lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}).
				Background(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#353C3B"}).
				Padding(0, 1, 0, 1).Align(lipgloss.Center)
			m8Text := barStyle.Copy().
				Bold(true).
				Foreground(lipgloss.Color(darkGrey)).
				Background(lipgloss.Color(gold)).MarginLeft(1).Render("m8")
			fileType := barStyle.Copy().Bold(true).
				Background(lipgloss.Color(teal)).Render("yaml")
			cfgPathVal := barStyle.Copy().
				Bold(true).
				Width(101 - lipgloss.Width(m8Text) - lipgloss.Width(fileType)).
				Render(viper.GetString("cfg-path"))

			barText := lipgloss.JoinHorizontal(lipgloss.Top,
				m8Text,
				cfgPathVal,
				fileType,
			)

			yamlText := lipgloss.NewStyle().
				Bold(true).
				Width(100).
				Foreground(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#FFFFFF"}).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.AdaptiveColor{Light: "#01A299", Dark: "#01A299"}).
				Padding(1, 1, 1, 1).Render(string(cfgYaml))

			output := lipgloss.JoinVertical(lipgloss.Top, barText, yamlText)

			fmt.Println(output)
			return nil
		},
	}
	return showCfg
}
