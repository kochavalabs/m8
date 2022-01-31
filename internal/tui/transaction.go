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

var _ tea.Model = &TxModel{}

type TxModel struct {
	stopwatch stopwatch.Model

	cmd  tea.Cmd
	tx   *xdr.Transaction
	rcpt *xdr.Receipt
	id   *xdr.ID
	err  error

	quit bool
}

// TODO explain that this is to protect the TxModel from other CMD types
type TxCmd func() tea.Msg

func NewTxModel(txCmd TxCmd) *TxModel {
	return &TxModel{
		cmd:       tea.Cmd(txCmd),
		stopwatch: stopwatch.NewWithInterval(time.Millisecond),
	}
}

func (t TxModel) Init() tea.Cmd {
	return tea.Batch(t.stopwatch.Init(), t.cmd)
}

func (t TxModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return t, tea.Quit
		default:
			return t, nil
		}
	case *xdr.Transaction:
		t.tx = msg
		return t, tea.Quit
	case *xdr.ID:
		t.id = msg
		t.quit = true
		return t, tea.Quit
	case *xdr.Receipt:
		t.rcpt = msg
		t.quit = true
		return t, tea.Quit
	case error:
		t.err = error(msg)
		t.quit = true
		return t, nil
	default:
		if !t.quit {
			var cmd tea.Cmd
			t.stopwatch, cmd = t.stopwatch.Update(msg)
			return t, cmd
		}
		return t, tea.Quit
	}
}

func (t TxModel) View() string {

	barStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}).
		Background(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#353C3B"}).
		Padding(0, 1, 0, 1).Align(lipgloss.Center)
	m8Text := barStyle.Copy().
		Bold(true).
		Foreground(lipgloss.Color("#353C3B")).
		Background(lipgloss.Color("#E3BD2D")).MarginLeft(1).Render("m8")
	fileType := barStyle.Copy().Bold(true).
		Background(lipgloss.Color("#01A299")).Render("exec")
	cfgPathVal := barStyle.Copy().
		Bold(true).
		Width(101 - lipgloss.Width(m8Text) - lipgloss.Width(fileType)).
		Render("")

	barText := lipgloss.JoinHorizontal(lipgloss.Top,
		m8Text,
		cfgPathVal,
		fileType,
	)

	output := ""
	if t.tx != nil {
		v, err := json.MarshalIndent(t.tx, "", "\t")
		if err != nil {
			return err.Error()
		}
		output = string(v)
	}

	if t.rcpt != nil {
		v, err := json.MarshalIndent(t.rcpt, "", "\t")
		if err != nil {
			return err.Error()
		}
		output = string(v)
	}

	if t.err != nil {
		errText := lipgloss.NewStyle().
			Bold(true).
			Width(100).
			Foreground(lipgloss.AdaptiveColor{Light: red, Dark: red}).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#01A299", Dark: "#01A299"}).
			Padding(1, 1, 1, 1).Render("error: " + t.err.Error())
		output = errText
	}

	sw := t.stopwatch.View()
	sw = "Elapsed: " + sw + "\n"

	return lipgloss.JoinVertical(lipgloss.Top, barText, output, sw)
}

func TxLookup(ctx context.Context, client mazzaroth.Client, channelId string, transactionId string) TxCmd {
	return func() tea.Msg {
		tx, err := client.TransactionLookup(ctx, channelId, transactionId)
		if err != nil {
			return err
		}
		return tx
	}
}

func TxCall(ctx context.Context, client mazzaroth.Client, tx *xdr.Transaction) TxCmd {
	return func() tea.Msg {
		id, rcpt, err := client.TransactionSubmit(ctx, tx)
		if err != nil {
			return err
		}
		if rcpt != nil {
			return rcpt
		}
		return id
	}
}
