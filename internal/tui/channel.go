package tui

import tea "github.com/charmbracelet/bubbletea"

var _ tea.Model = &ChannelModel{}

type ChannelModel struct {
}

func (c ChannelModel) Init() tea.Cmd {
	return nil
}

func (c ChannelModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (c ChannelModel) View() string {
	return ""
}
