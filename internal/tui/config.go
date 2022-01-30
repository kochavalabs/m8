package tui

import tea "github.com/charmbracelet/bubbletea"

var _ tea.Model = &CfgModel{}

type CfgModel struct {
}

func (c CfgModel) Init() tea.Cmd {
	return nil
}

func (c CfgModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (c CfgModel) View() string {
	return ""
}
