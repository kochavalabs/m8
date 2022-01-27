package tui

import (
	"context"
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kochavalabs/mazzaroth-go"
	"github.com/kochavalabs/mazzaroth-xdr/go-xdr/xdr"
)

var _ tea.Model = &TxModel{}

type TxModel struct {
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
	return t.cmd
}

func (t TxModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (t TxModel) View() string {
	return ""
}

func txLookup(ctx context.Context, client mazzaroth.Client, channelId string, transactionId string) tea.Cmd {
	return func() tea.Msg {
		tx, err := client.TransactionLookup(ctx, channelId, transactionId)
		if err != nil {
			return err
		}

		v, err := json.MarshalIndent(tx, "", "\t")
		if err != nil {
			return err
		}

		return v
	}
}

func txCall(ctx context.Context, client mazzaroth.Client, channelId string, sender string, function string, args []xdr.Argument)
