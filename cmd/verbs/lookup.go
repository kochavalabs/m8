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

// lookup
// blockheight
// block
// tx
// receipt

const (
	gold     = `#E3BD2D`
	darkGrey = `#353C3B`
	teal     = `#01A299`
	white    = `#FFFFFF`
)

func Lookup(resource string) *cobra.Command {
	lookup := &cobra.Command{
		Use:   "lookup",
		Short: "look up items on a mazzaroth node",
	}
	// sub command chain by resource type
	switch resource {
	case "channel":
		lookup.AddCommand(lookupAbi(), lookupBlock(), lookupTx(), lookupReceipt())
	case "cfg":
		lookup.AddCommand(lookupCfg())
	}
	return lookup
}

func lookupCfg() *cobra.Command {
	cfg := &cobra.Command{
		Use:   "cfg",
		Short: "look up items on a mazzaroth node",
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
	return cfg
}

func lookupAbi() *cobra.Command {
	abi := &cobra.Command{
		Use:   "abi",
		Short: "look up items on a mazzaroth node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return abi
}

func lookupBlock() *cobra.Command {
	block := &cobra.Command{
		Use:   "block",
		Short: "look up items on a mazzaroth node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return block
}

func lookupTx() *cobra.Command {
	tx := &cobra.Command{
		Use:   "tx",
		Short: "look up items on a mazzaroth node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return tx
}

func lookupReceipt() *cobra.Command {
	rcpt := &cobra.Command{
		Use:   "rcpt",
		Short: "look up items on a mazzaroth node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return rcpt
}
