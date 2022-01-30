package tui

import tea "github.com/charmbracelet/bubbletea"

var _ tea.Model = &RcptModel{}

type RcptModel struct {
}

func (r RcptModel) Init() tea.Cmd {
	return nil
}

func (r RcptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (r RcptModel) View() string {
	return ""
}
