package components

import (
	"fmt"
	"strings"

	"tax-calculator/internal/tax/views/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TaxClassInfo struct {
	Class       int
	Name        string
	Description string
	Icon        string
	Details     string
	CommonUse   string
}

type TaxClassSelector struct {
	Options     []TaxClassInfo
	Selected    int
	Focused     bool
	ShowDetails bool
}

func NewTaxClassSelector() TaxClassSelector {
	options := []TaxClassInfo{
		{
			Class:       1,
			Name:        "Single Person",
			Description: "Unmarried, divorced, or widowed",
			Icon:        "👤",
			Details:     "Standard tax class for single individuals",
			CommonUse:   "Most common for unmarried people",
		},
		{
			Class:       2,
			Name:        "Single Parent",
			Description: "Single with dependent children",
			Icon:        "👨‍👦",
			Details:     "Includes child allowance benefits",
			CommonUse:   "Single parents eligible for tax relief",
		},
		{
			Class:       3,
			Name:        "Married (Higher Earner)",
			Description: "Married, higher income spouse",
			Icon:        "💍",
			Details:     "Lower tax rate for main earner",
			CommonUse:   "When one spouse earns significantly more",
		},
		{
			Class:       4,
			Name:        "Married (Equal Income)",
			Description: "Married, both spouses work equally",
			Icon:        "👫",
			Details:     "Both spouses taxed individually",
			CommonUse:   "When both spouses have similar incomes",
		},
		{
			Class:       5,
			Name:        "Married (Lower Earner)",
			Description: "Married, lower income spouse",
			Icon:        "👪",
			Details:     "Higher tax rate, paired with Class 3",
			CommonUse:   "Lower-earning spouse in unequal income marriage",
		},
		{
			Class:       6,
			Name:        "Secondary Employment",
			Description: "Additional jobs or secondary income",
			Icon:        "🏢",
			Details:     "Higher tax rate for second job",
			CommonUse:   "Multiple employment relationships",
		},
	}

	return TaxClassSelector{
		Options:     options,
		Selected:    1,
		Focused:     false,
		ShowDetails: false,
	}
}

func (tcs *TaxClassSelector) Update(msg tea.Msg) (TaxClassSelector, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			tcs.Selected--
			if tcs.Selected < 1 {
				tcs.Selected = 6
			}
		case "down", "j":
			tcs.Selected++
			if tcs.Selected > 6 {
				tcs.Selected = 1
			}
		case "h", "?":
			tcs.ShowDetails = !tcs.ShowDetails
		}
	}
	return *tcs, nil
}

func (tcs *TaxClassSelector) View() string {
	var builder strings.Builder

	titleStyle := styles.SubtitleStyle
	builder.WriteString(titleStyle.Render("Tax Class"))
	builder.WriteString("\n")

	helpStyle := lipgloss.NewStyle().
		Foreground(styles.NeutralColor).
		Italic(true)
	builder.WriteString(helpStyle.Render("Choose your German tax classification"))
	builder.WriteString("\n\n")

	for _, option := range tcs.Options {
		isSelected := option.Class == tcs.Selected

		var style lipgloss.Style
		var indicator string

		if isSelected {
			style = styles.SelectedItemStyle
			indicator = "▶ "
		} else {
			style = styles.UnselectedItemStyle
			indicator = "  "
		}

		mainLine := fmt.Sprintf("%s%s %s",
			indicator,
			option.Icon,
			option.Name)

		builder.WriteString(style.Render(mainLine))
		builder.WriteString("\n")

		descStyle := lipgloss.NewStyle().
			Foreground(styles.NeutralColor).
			MarginLeft(4)
		builder.WriteString(descStyle.Render(option.Description))
		builder.WriteString("\n")

		if isSelected && tcs.ShowDetails {
			detailStyle := lipgloss.NewStyle().
				Foreground(styles.AccentColor).
				MarginLeft(4).
				Italic(true)
			builder.WriteString(detailStyle.Render("• " + option.Details))
			builder.WriteString("\n")
			builder.WriteString(detailStyle.Render("• " + option.CommonUse))
			builder.WriteString("\n")
		}

		builder.WriteString("\n")
	}

	if tcs.Focused {
		helpText := lipgloss.NewStyle().
			Foreground(styles.NeutralColor).
			Italic(true).
			Render("Press 'h' for more details • ↑/↓ to navigate")
		builder.WriteString(helpText)
	}

	return builder.String()
}

func (tcs *TaxClassSelector) GetSelected() TaxClassInfo {
	for _, option := range tcs.Options {
		if option.Class == tcs.Selected {
			return option
		}
	}
	return tcs.Options[0]
}

func (tcs *TaxClassSelector) SetSelected(class int) {
	if class >= 1 && class <= 6 {
		tcs.Selected = class
	}
}

func (tcs *TaxClassSelector) Focus() {
	tcs.Focused = true
}

func (tcs *TaxClassSelector) Blur() {
	tcs.Focused = false
}
