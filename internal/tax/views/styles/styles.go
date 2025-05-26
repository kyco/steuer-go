package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Elegant minimal color palette
	PrimaryColor   = lipgloss.Color("#7b9cd9") // Soft blue
	SecondaryColor = lipgloss.Color("#323232") // Charcoal
	AccentColor    = lipgloss.Color("#db7093") // Pale violet red
	SuccessColor   = lipgloss.Color("#76b876") // Muted green
	DangerColor    = lipgloss.Color("#e06c75") // Soft red
	WarningColor   = lipgloss.Color("#e5c07b") // Soft amber
	NeutralColor   = lipgloss.Color("#abb2bf") // Silver
	BgColor        = lipgloss.Color("#282c34") // Dark background
	FgColor        = lipgloss.Color("#f8f8f2") // Off-white text

	// Simplified border styles
	SimpleBorder  = lipgloss.NormalBorder()
	ThickBorder   = lipgloss.ThickBorder()
	RoundedBorder = lipgloss.RoundedBorder()

	// Core styles
	BaseStyle = lipgloss.NewStyle().
			Foreground(FgColor).
			Background(BgColor).
			PaddingLeft(1)

	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Padding(1, 2)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			MarginLeft(1).
			MarginBottom(1).
			Padding(0, 1)

	// Input fields
	InputFieldStyle = lipgloss.NewStyle().
			Border(SimpleBorder).
			BorderForeground(NeutralColor).
			Padding(0, 1).
			MarginRight(1)

	ActiveInputStyle = InputFieldStyle.Copy().
				BorderForeground(PrimaryColor).
				Bold(true)

	// Selection styles
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(BgColor).
				Background(PrimaryColor).
				Bold(true).
				Padding(0, 1)

	UnselectedItemStyle = lipgloss.NewStyle().
				Foreground(FgColor)

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Border(SimpleBorder).
			BorderForeground(NeutralColor).
			Padding(0, 2).
			MarginRight(1)

	SelectedButtonStyle = ButtonStyle.Copy().
				BorderForeground(PrimaryColor).
				Bold(true)

	// Container styles
	ContainerStyle = lipgloss.NewStyle().
			Border(SimpleBorder).
			BorderForeground(NeutralColor).
			Padding(1)

	FocusedContainerStyle = ContainerStyle.Copy().
				BorderForeground(PrimaryColor)

	// Results styles
	ResultsContainerStyle = lipgloss.NewStyle().
				Border(SimpleBorder).
				BorderForeground(NeutralColor).
				Padding(1, 2)

	// Helper styles
	HelpStyle = lipgloss.NewStyle().
			Foreground(NeutralColor).
			Italic(true)

	KeyHintStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Background(SecondaryColor).
			Padding(0, 1).
			Bold(true)

	// Visual elements
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)

	// Tab styles - simplified
	TabStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(NeutralColor)

	ActiveTabStyle = TabStyle.Copy().
			Foreground(PrimaryColor).
			Underline(true).
			Bold(true)

	// Data visualization
	ProgressBarEmptyStyle = lipgloss.NewStyle().
				Foreground(NeutralColor).
				Background(lipgloss.Color("#23272e")). // darker background for empty
				Faint(true)

	ProgressBarFilledStyle = lipgloss.NewStyle().
				Foreground(FgColor).
				Background(PrimaryColor).
				Bold(true).
				Underline(true).
				Padding(0, 0)
)
