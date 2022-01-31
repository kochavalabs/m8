package tui

import tea "github.com/charmbracelet/bubbletea"

var _ tea.Model = &ManifestModel{}

type ManifestModel struct {
}

func (m ManifestModel) Init() tea.Cmd {
	return nil
}

func (m ManifestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m ManifestModel) View() string {
	return ""
}
