package tui

import (
	"context"
	"encoding/json"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kochavalabs/mazzaroth-go"
	"github.com/kochavalabs/mazzaroth-xdr/go-xdr/xdr"
)

var _ tea.Model = &TxModel{}

type rcptMsg *xdr.Receipt

type txIdMsg xdr.ID

type txMsg *xdr.Transaction

type TxModel struct {
	stopwatch stopwatch.Model

	cmd  tea.Cmd
	tx   *xdr.Transaction
	rcpt *xdr.Receipt
	id   xdr.ID
}

func NewTxModel(cmd tea.Cmd) *TxModel {
	return &TxModel{
		cmd: cmd,
	}
}

func (t TxModel) Init() tea.Cmd {
	return tea.Batch(t.stopwatch.Init(), t.cmd)
}

func (t TxModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case *xdr.Transaction:
		return t, tea.Quit
	case xdr.ID:
		return t, tea.Quit
	case error:
		return t, tea.Quit
	default:
		var cmd tea.Cmd
		t.stopwatch, cmd = t.stopwatch.Update(msg)
		return t, cmd
	}
}

func (t TxModel) View() string {
	if t.tx != nil {
		v, err := json.MarshalIndent(t.tx, "", "\t")
		if err != nil {
			return err.Error()
		}
		return string(v)
	}
	// Note: you could further customize the time output by getting the
	// duration from m.stopwatch.Elapsed(), which returns a time.Duration, and
	// skip m.stopwatch.View() altogether.
	s := t.stopwatch.View() + "\n"
	s = "Elapsed: " + s

	return s
}

func TxLookup(ctx context.Context, client mazzaroth.Client, channelId string, transactionId string) tea.Cmd {
	return func() tea.Msg {
		tx, err := client.TransactionLookup(ctx, channelId, transactionId)
		if err != nil {
			return err
		}
		return tx
	}
}

func TxCall(ctx context.Context, client mazzaroth.Client, tx *xdr.Transaction) tea.Cmd {
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
