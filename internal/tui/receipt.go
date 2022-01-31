package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kochavalabs/mazzaroth-xdr/go-xdr/xdr"
)

var _ tea.Model = &RcptModel{}

type RcptModel struct {
	stopwatch stopwatch.Model
	cmd       tea.Cmd

	rcpt *xdr.Receipt
	quit bool
}

type RcptCmd func() tea.Msg

func NewRcptModel(rcptCmd RcptCmd) *RcptModel {
	return &RcptModel{
		cmd:       tea.Cmd(rcptCmd),
		stopwatch: stopwatch.NewWithInterval(time.Millisecond),
	}
}

func (r RcptModel) Init() tea.Cmd {
	return tea.Batch(r.stopwatch.Init(), r.cmd)
}

func (r RcptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return r, tea.Quit
		default:
			return r, nil
		}
	case *xdr.Receipt:
		return r, nil
	case error:
		return r, nil
	default:
		if !r.quit {
			var cmd tea.Cmd
			r.stopwatch, cmd = r.stopwatch.Update(msg)
			return r, cmd
		}
		return r, tea.Quit
	}
}

func (r RcptModel) View() string {
	return ""
}
