package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	PrimaryColor   = lipgloss.Color("#82AAFF") 
	SecondaryColor = lipgloss.Color("#C3E88D") 
	AccentColor    = lipgloss.Color("#89DDFF") 
	SuccessColor   = lipgloss.Color("#C3E88D") 
	DangerColor    = lipgloss.Color("#FF5370") 
	WarningColor   = lipgloss.Color("#FFCB6B") 
	NeutralColor   = lipgloss.Color("#A9B8E8") 
	BgColor        = lipgloss.Color("#0F111A") 
	FgColor        = lipgloss.Color("#FFFFFF") 
	
	MinimalBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
	}

	BaseStyle = lipgloss.NewStyle().
			Foreground(FgColor).
			Background(BgColor).
			PaddingLeft(1).
			MarginBottom(0)

	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Background(BgColor).
			MarginLeft(1).
			MarginTop(1).
			MarginBottom(1).
			PaddingLeft(1).
			PaddingRight(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Background(BgColor).
			MarginLeft(1).
			MarginBottom(1).
			Bold(true)

	InputFieldStyle = lipgloss.NewStyle().
			Foreground(FgColor).
			Background(BgColor).
			Padding(0, 1).
			MarginLeft(1).
			MarginRight(1).
			Border(MinimalBorder).
			BorderForeground(AccentColor)

	ActiveInputStyle = InputFieldStyle.Copy().
				BorderForeground(PrimaryColor).
				Bold(true)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(BgColor).
				Bold(true).
				Background(PrimaryColor).
				PaddingLeft(1).
				PaddingRight(1)

	UnselectedItemStyle = lipgloss.NewStyle().
				Foreground(FgColor).
				Bold(false)

	ButtonStyle = lipgloss.NewStyle().
			Foreground(BgColor).
			Background(SecondaryColor).
			Padding(0, 2).
			MarginTop(0).
			MarginRight(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SecondaryColor)

	SelectedButtonStyle = ButtonStyle.Copy().
				Background(PrimaryColor).
				BorderForeground(PrimaryColor).
				Foreground(BgColor).
				Bold(true)

	ResultsBoxStyle = lipgloss.NewStyle().
				Border(MinimalBorder).
				BorderForeground(AccentColor).
				Padding(1, 1).
				MarginTop(1).
				MarginBottom(1).
				MarginLeft(1).
				MarginRight(1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Italic(true).
			MarginLeft(1)
			
	HeaderStyle = lipgloss.NewStyle().
			Foreground(BgColor).
			Background(PrimaryColor).
			Bold(true).
			Width(100).
			Align(lipgloss.Center).
			Padding(0, 0).
			MarginBottom(1)
			
	FooterStyle = lipgloss.NewStyle().
			Foreground(NeutralColor).
			Background(BgColor).
			Width(100).
			Align(lipgloss.Center).
			Padding(0, 0).
			MarginTop(1).
			Border(lipgloss.Border{
				Top: "─",
			}).
			BorderTop(true).
			BorderForeground(SecondaryColor)
)