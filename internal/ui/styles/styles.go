package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// High contrast Catppuccino-inspired theme
	PrimaryColor   = lipgloss.Color("#82AAFF") // Vibrant blue - main color 
	SecondaryColor = lipgloss.Color("#C3E88D") // Bright green - secondary elements
	AccentColor    = lipgloss.Color("#89DDFF") // Bright cyan - highlights
	SuccessColor   = lipgloss.Color("#C3E88D") // Bright green - success indicators
	DangerColor    = lipgloss.Color("#FF5370") // Bright red - errors/warnings
	WarningColor   = lipgloss.Color("#FFCB6B") // Amber/yellow - warning/debug messages
	NeutralColor   = lipgloss.Color("#A9B8E8") // Light blue-gray - less important elements
	BgColor        = lipgloss.Color("#0F111A") // Deep navy - background 
	FgColor        = lipgloss.Color("#FFFFFF") // Pure white - main text
	
	// Simple, clean border for minimalist design
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

	// Base styles
	BaseStyle = lipgloss.NewStyle().
			Foreground(FgColor).
			Background(BgColor).
			PaddingLeft(1).
			MarginBottom(0)

	// Title style - high contrast
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Background(BgColor).
			MarginLeft(1).
			MarginTop(1).
			MarginBottom(1).
			PaddingLeft(1).
			PaddingRight(1)

	// Subtitle style with more contrast
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Background(BgColor).
			MarginLeft(1).
			MarginBottom(1).
			Bold(true)

	// Input field style - high contrast boxes
	InputFieldStyle = lipgloss.NewStyle().
			Foreground(FgColor).
			Background(BgColor).
			Padding(0, 1).
			MarginLeft(1).
			MarginRight(1).
			Border(MinimalBorder).
			BorderForeground(AccentColor)

	// Active input style
	ActiveInputStyle = InputFieldStyle.Copy().
				BorderForeground(PrimaryColor).
				Bold(true)

	// Selected option style
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(BgColor).
				Bold(true).
				Background(PrimaryColor).
				PaddingLeft(1).
				PaddingRight(1)

	// Unselected option style
	UnselectedItemStyle = lipgloss.NewStyle().
				Foreground(FgColor).
				Bold(false)

	// Button style - high contrast buttons
	ButtonStyle = lipgloss.NewStyle().
			Foreground(BgColor).
			Background(SecondaryColor).
			Padding(0, 2).
			MarginTop(0).
			MarginRight(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SecondaryColor)

	// Selected button style
	SelectedButtonStyle = ButtonStyle.Copy().
				Background(PrimaryColor).
				BorderForeground(PrimaryColor).
				Foreground(BgColor).
				Bold(true)

	// Results box style - clean minimalist panel
	ResultsBoxStyle = lipgloss.NewStyle().
				Border(MinimalBorder).
				BorderForeground(AccentColor).
				Padding(1, 1).
				MarginTop(1).
				MarginBottom(1).
				MarginLeft(1).
				MarginRight(1)

	// Help style - improved contrast help text
	HelpStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Italic(true).
			MarginLeft(1)
			
	// Header style - Catppuccino-themed header (full width)
	HeaderStyle = lipgloss.NewStyle().
			Foreground(BgColor).
			Background(PrimaryColor).
			Bold(true).
			Width(100).
			Align(lipgloss.Center).
			Padding(0, 0).
			MarginBottom(1)
			
	// Footer style - Catppuccino-themed footer (full width)
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