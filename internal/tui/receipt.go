package tui

import (
	"context"
	"encoding/json"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kochavalabs/mazzaroth-go"
	"github.com/kochavalabs/mazzaroth-xdr/go-xdr/xdr"
)

var _ tea.Model = &RcptModel{}

type RcptModel struct {
	stopwatch stopwatch.Model
	cmd       tea.Cmd

	rcpt *xdr.Receipt
	err  error
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
		r.rcpt = msg
		r.quit = true
		return r, tea.Quit
	case error:
		r.err = error(msg)
		r.quit = true
		return r, tea.Quit
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
	barStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}).
		Background(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#353C3B"}).
		Padding(0, 1, 0, 1).Align(lipgloss.Center)
	m8Text := barStyle.Copy().
		Bold(true).
		Foreground(lipgloss.Color("#353C3B")).
		Background(lipgloss.Color("#E3BD2D")).MarginLeft(1).Render("m8")
	fileType := barStyle.Copy().Bold(true).
		Background(lipgloss.Color("#01A299")).Render("json")
	cfgPathVal := barStyle.Copy().
		Bold(true).
		Width(101 - lipgloss.Width(m8Text) - lipgloss.Width(fileType)).
		Render("receipt")

	barText := lipgloss.JoinHorizontal(lipgloss.Top,
		m8Text,
		cfgPathVal,
		fileType,
	)

	output := ""
	if r.rcpt != nil {
		v, err := json.MarshalIndent(r.rcpt, "", " ")
		if err != nil {
			r.err = err
		} else {
			output = string(v)
		}
	} else if r.err != nil {
		errText := lipgloss.NewStyle().
			Bold(true).
			Width(100).
			Foreground(lipgloss.AdaptiveColor{Light: red, Dark: red}).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#01A299", Dark: "#01A299"}).
			Padding(1, 1, 1, 1).Render("error: " + r.err.Error())
		output = errText
	}

	sw := r.stopwatch.View()
	sw = "Elapsed: " + sw + "\n"

	return lipgloss.JoinVertical(lipgloss.Top, barText, output, sw)
}

func RcptLookup(ctx context.Context, client mazzaroth.Client, channelId string, transactionId string) RcptCmd {
	return func() tea.Msg {
		receipt, err := client.ReceiptLookup(ctx, channelId, transactionId)
		if err != nil {
			return err
		}
		return receipt
	}
}
