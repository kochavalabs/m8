package tui

import (
	"context"
	"encoding/json"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kochavalabs/mazzaroth-go"
	"github.com/kochavalabs/mazzaroth-xdr/go-xdr/xdr"
)

var _ tea.Model = &TxModel{}

type TxModel struct {
	stopwatch stopwatch.Model

	cmd  tea.Cmd
	tx   *xdr.Transaction
	rcpt *xdr.Receipt
	id   xdr.ID
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
	case xdr.ID:
		t.id = msg
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
	if t.tx != nil {
		v, err := json.MarshalIndent(t.tx, "", "\t")
		if err != nil {
			return err.Error()
		}
		return string(v)
	}
	output := ""
	if t.err != nil {
		output = t.err.Error() + "\n"
	}
	// Note: you could further customize the time output by getting the
	// duration from m.stopwatch.Elapsed(), which returns a time.Duration, and
	// skip m.stopwatch.View() altogether.

	s := t.stopwatch.View()
	s = "Elapsed: " + s + "\n"

	return output + s
}

func TxLookup(ctx context.Context, client mazzaroth.Client, channelId string, transactionId string) TxCmd {
	return func() tea.Msg {
		time.Sleep(time.Second * 2)
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
