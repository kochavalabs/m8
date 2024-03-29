package tui

import "github.com/charmbracelet/lipgloss"

const (
	gold     = `#E3BD2D`
	darkGrey = `#353C3B`
	teal     = `#01A299`
	white    = `#FFFFFF`
	red      = `#F31849`
)

var (
	barStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.AdaptiveColor{Light: white, Dark: white}).
		Background(lipgloss.AdaptiveColor{Light: darkGrey, Dark: darkGrey}).
		Padding(0, 1, 0, 1).Align(lipgloss.Center)
)
