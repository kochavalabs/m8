package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kochavalabs/m8/internal/cfg"
	"gopkg.in/yaml.v2"
)

var _ tea.Model = &CfgModel{}

type CfgCmd func() tea.Msg

type cfgMsg struct {
	path string
	body []byte
}

type CfgModel struct {
	cmd tea.Cmd

	cfg *cfgMsg
	err error
}

func NewCfgModel(cfgCmd CfgCmd) *CfgModel {
	return &CfgModel{
		cmd: tea.Cmd(cfgCmd),
		cfg: &cfgMsg{},
	}
}

func (c CfgModel) Init() tea.Cmd {
	return c.cmd
}

func (c CfgModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return c, tea.Quit
		default:
			return c, nil
		}
	case *cfgMsg:
		c.cfg = &cfgMsg{
			body: []byte(msg.body),
			path: msg.path,
		}
		return c, tea.Quit
	case error:
		c.err = error(msg)
		return c, tea.Quit
	}
	return c, nil
}

func (c CfgModel) View() string {
	output := ""
	bodyText := ""

	m8Text := barStyle.Copy().
		Foreground(lipgloss.Color(darkGrey)).
		Background(lipgloss.Color(gold)).MarginLeft(1).Render("m8")

	fileType := barStyle.Copy().
		Background(lipgloss.Color(teal)).Render("yaml")

	cfgPathVal := barStyle.Copy().
		Width(101 - lipgloss.Width(m8Text) - lipgloss.Width(fileType)).
		Render(c.cfg.path)

	barText := lipgloss.JoinHorizontal(lipgloss.Top, m8Text, cfgPathVal, fileType)

	if c.cfg != nil {
		if c.cfg.body != nil {
			bodyText = string(c.cfg.body)
		}

		yamlText := lipgloss.NewStyle().
			Bold(true).
			Width(100).
			Foreground(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#FFFFFF"}).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#01A299", Dark: "#01A299"}).
			Padding(1, 1, 1, 1).Render(string(bodyText))

		output = lipgloss.JoinVertical(lipgloss.Top, barText, yamlText)
	}

	if c.err != nil {
		errText := lipgloss.NewStyle().
			Bold(true).
			Width(100).
			Foreground(lipgloss.AdaptiveColor{Light: red, Dark: red}).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#01A299", Dark: "#01A299"}).
			Padding(1, 1, 1, 1).Render("error: " + c.err.Error())

		output = lipgloss.JoinVertical(lipgloss.Top, barText, errText)
	}

	return output + "\n"
}

func CfgShow(cfg *cfg.Configuration, path string) CfgCmd {
	return func() tea.Msg {
		cfgYaml, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}
		return &cfgMsg{
			path: path,
			body: cfgYaml,
		}
	}
}
