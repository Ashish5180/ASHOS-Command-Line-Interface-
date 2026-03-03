package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	PrimaryColor   = lipgloss.Color("#7D56F4") // Purple
	SecondaryColor = lipgloss.Color("#04B575") // Green
	AccentColor    = lipgloss.Color("#EE6FF8") // Pink
	GrayColor      = lipgloss.Color("#626262")
	WhiteColor     = lipgloss.Color("#EEEEEE")

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(WhiteColor).
			Background(PrimaryColor).
			Padding(0, 2).
			MarginBottom(1)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(1).
			MarginRight(2)

	StatsBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SecondaryColor).
			Padding(1)

	FocusBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(AccentColor).
			Padding(1).
			MarginTop(1)

	StatusStyle = lipgloss.NewStyle().
			Foreground(GrayColor).
			Italic(true)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	SubTitleStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true).
			MarginBottom(1)

	GrayStyle = lipgloss.NewStyle().
			Foreground(GrayColor)
)

func MakeBox(title string, content string, style lipgloss.Style) string {
	t := SubTitleStyle.Render(strings.ToUpper(title))
	return style.Render(t + "\n" + content)
}
