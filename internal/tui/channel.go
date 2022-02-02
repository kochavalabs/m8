package tui

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kochavalabs/mazzaroth-go"
	"github.com/kochavalabs/mazzaroth-xdr/go-xdr/xdr"
)

var _ tea.Model = &ChannelModel{}

type ChannelModel struct {
	stopwatch stopwatch.Model
	cmd       tea.Cmd
	id        *xdr.ID
	abi       *xdr.Abi
	err       error

	quit bool
}

type ChannelCmd func() tea.Msg

func NewChannelModel(channelCmd ChannelCmd) *ChannelModel {
	return &ChannelModel{
		cmd:       tea.Cmd(channelCmd),
		stopwatch: stopwatch.NewWithInterval(time.Millisecond),
	}
}

func (c ChannelModel) Init() tea.Cmd {
	return tea.Batch(c.stopwatch.Init(), c.cmd)
}

func (c ChannelModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return c, tea.Quit
		default:
			return c, nil
		}
	case *xdr.ID:
		c.id = msg
		c.quit = true
		return c, tea.Quit
	case *xdr.Abi:
		c.abi = msg
		c.quit = true
		return c, tea.Quit
	case error:
		c.err = error(msg)
		c.quit = true
		return c, tea.Quit
	default:
		if !c.quit {
			var cmd tea.Cmd
			c.stopwatch, cmd = c.stopwatch.Update(msg)
			return c, cmd
		}
		return c, tea.Quit
	}
}

func (c ChannelModel) View() string {

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
		Render("channel")

	barText := lipgloss.JoinHorizontal(lipgloss.Top,
		m8Text,
		cfgPathVal,
		fileType,
	)

	output := ""
	if c.id != nil {
		output = hex.EncodeToString(c.id[:])
	} else if c.abi != nil {
		v, err := json.MarshalIndent(c.abi, "", " ")
		if err != nil {
			c.err = err
		} else {
			output = string(v)
		}
	} else if c.err != nil {
		errText := lipgloss.NewStyle().
			Bold(true).
			Width(100).
			Foreground(lipgloss.AdaptiveColor{Light: red, Dark: red}).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#01A299", Dark: "#01A299"}).
			Padding(1, 1, 1, 1).Render("error: " + c.err.Error())
		output = errText
	}

	sw := c.stopwatch.View()
	sw = "Elapsed: " + sw + "\n"
	return lipgloss.JoinVertical(lipgloss.Top, barText, output, sw)
}

func ChannelDelete(ctx context.Context, client mazzaroth.Client, tx *xdr.Transaction) ChannelCmd {
	return func() tea.Msg {
		if tx.Data.Category.Type == xdr.CategoryTypeDELETE {
			id, _, err := client.TransactionSubmit(ctx, tx)
			if err != nil {
				return err
			}
			return id
		}
		return errors.New("invalid trnasaction type supplied to delete cmd")
	}
}

func ChannelPause(ctx context.Context, client mazzaroth.Client, tx *xdr.Transaction) ChannelCmd {
	return func() tea.Msg {
		if tx.Data.Category.Type == xdr.CategoryTypePAUSE {
			id, _, err := client.TransactionSubmit(ctx, tx)
			if err != nil {
				return err
			}
			return id
		}
		return errors.New("invalid transaction type supplied to pause cmd")
	}
}

func ChannelAbiLookup(ctx context.Context, client mazzaroth.Client, channelId string) ChannelCmd {
	return func() tea.Msg {
		abi, err := client.ChannelAbi(ctx, channelId)
		if err != nil {
			return err
		}
		return abi
	}
}
