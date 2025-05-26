package components

import (
	"fmt"
	"strings"

	"tax-calculator/internal/tax/views/styles"

	"github.com/charmbracelet/lipgloss"
)

type StatusBar struct {
	CurrentScreen string
	CurrentMode   string
	AvailableKeys []KeyHint
	Width         int
}

type KeyHint struct {
	Key         string
	Description string
	Important   bool
}

func NewStatusBar(width int) StatusBar {
	return StatusBar{
		Width: width,
	}
}

func (sb *StatusBar) SetScreen(screen string) {
	sb.CurrentScreen = screen
}

func (sb *StatusBar) SetMode(mode string) {
	sb.CurrentMode = mode
}

func (sb *StatusBar) SetKeys(keys []KeyHint) {
	sb.AvailableKeys = keys
}

func (sb *StatusBar) View() string {
	leftSection := sb.renderLeftSection()
	rightSection := sb.renderRightSection()

	// Calculate available space for center section
	leftWidth := lipgloss.Width(leftSection)
	rightWidth := lipgloss.Width(rightSection)
	centerWidth := sb.Width - leftWidth - rightWidth - 4 // padding

	centerSection := sb.renderCenterSection(centerWidth)

	statusStyle := lipgloss.NewStyle().
		Background(styles.SecondaryColor).
		Foreground(styles.FgColor).
		Padding(0, 1).
		Width(sb.Width)

	content := lipgloss.JoinHorizontal(
		lipgloss.Left,
		leftSection,
		centerSection,
		rightSection,
	)

	return statusStyle.Render(content)
}

func (sb *StatusBar) renderLeftSection() string {
	screenStyle := lipgloss.NewStyle().
		Foreground(styles.PrimaryColor).
		Bold(true)

	modeStyle := lipgloss.NewStyle().
		Foreground(styles.AccentColor)

	var parts []string
	parts = append(parts, screenStyle.Render(sb.CurrentScreen))

	if sb.CurrentMode != "" {
		parts = append(parts, modeStyle.Render(fmt.Sprintf("(%s)", sb.CurrentMode)))
	}

	return strings.Join(parts, " ")
}

func (sb *StatusBar) renderCenterSection(width int) string {
	if len(sb.AvailableKeys) == 0 {
		return strings.Repeat(" ", width)
	}

	var keyHints []string
	for _, hint := range sb.AvailableKeys {
		keyStyle := lipgloss.NewStyle().
			Foreground(styles.WarningColor).
			Bold(true)

		descStyle := lipgloss.NewStyle().
			Foreground(styles.NeutralColor)

		if hint.Important {
			keyStyle = keyStyle.Background(styles.WarningColor).Foreground(styles.BgColor)
		}

		keyHints = append(keyHints, fmt.Sprintf("%s %s",
			keyStyle.Render(hint.Key),
			descStyle.Render(hint.Description)))
	}

	hintsText := strings.Join(keyHints, " • ")

	// Center the hints in the available space
	if width > len(hintsText) {
		padding := (width - len(hintsText)) / 2
		hintsText = strings.Repeat(" ", padding) + hintsText
	}

	return lipgloss.NewStyle().Width(width).Render(hintsText)
}

func (sb *StatusBar) renderRightSection() string {
	timeStyle := lipgloss.NewStyle().
		Foreground(styles.NeutralColor)

	return timeStyle.Render("🧮 SteuerGo")
}

// Predefined key hints for different screens
func GetMainScreenKeys() []KeyHint {
	return []KeyHint{
		{Key: "Tab", Description: "Next Field", Important: false},
		{Key: "↑/↓", Description: "Tax Class", Important: false},
		{Key: "Enter", Description: "Calculate", Important: true},
		{Key: "L", Description: "Toggle Mode", Important: false},
		{Key: "Q", Description: "Quit", Important: false},
	}
}

func GetResultsScreenKeys() []KeyHint {
	return []KeyHint{
		{Key: "←/→", Description: "Switch Tab", Important: false},
		{Key: "C", Description: "Compare", Important: true},
		{Key: "B", Description: "Back", Important: false},
		{Key: "↑/↓", Description: "Scroll", Important: false},
		{Key: "Esc", Description: "Main", Important: false},
	}
}

func GetAdvancedScreenKeys() []KeyHint {
	return []KeyHint{
		{Key: "Tab", Description: "Next Field", Important: false},
		{Key: "Enter", Description: "Calculate", Important: true},
		{Key: "Esc", Description: "Back", Important: false},
		{Key: "↑/↓", Description: "Scroll", Important: false},
	}
}

func GetComparisonScreenKeys() []KeyHint {
	return []KeyHint{
		{Key: "↑/↓", Description: "Select Item", Important: false},
		{Key: "Enter", Description: "Show Breakdown", Important: true},
		{Key: "B", Description: "Back to Results", Important: false},
		{Key: "Esc", Description: "Main", Important: false},
	}
}

func GetComparisonBreakdownKeys() []KeyHint {
	return []KeyHint{
		{Key: "Enter", Description: "Back to List", Important: true},
		{Key: "B", Description: "Back to Results", Important: false},
		{Key: "Esc", Description: "Main", Important: false},
	}
}
