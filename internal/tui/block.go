package tui

import tea "github.com/charmbracelet/bubbletea"

var _ tea.Model = &BlockModel{}

type BlockModel struct {
}

func (b BlockModel) Init() tea.Cmd {
	return nil
}

func (b BlockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (b BlockModel) View() string {
	return ""
}
