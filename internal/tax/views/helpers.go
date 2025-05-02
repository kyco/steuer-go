package views

import (
	"fmt"
	"strconv"
	"strings"

	"tax-calculator/internal/tax/models"
	"tax-calculator/internal/tax/views/styles"

	"github.com/charmbracelet/lipgloss"
)

// Format a float number as a euro amount
func formatEuro(amount float64) string {
	return fmt.Sprintf("€ %.2f", amount)
}

// Format a percentage
func formatPercent(percent float64) string {
	return fmt.Sprintf("%.2f%%", percent)
}

// Create a minimal progress bar
func createProgressBar(percent float64, width int, highlight bool, label string, withBorder bool) string {
	if width < 3 {
		width = 3
	}

	filled := int((percent / 100) * float64(width))
	if filled > width {
		filled = width
	}

	empty := width - filled

	filledStr := styles.ProgressBarFilledStyle.Render(strings.Repeat("━", filled))
	emptyStr := styles.ProgressBarEmptyStyle.Render(strings.Repeat("─", empty))

	bar := filledStr + emptyStr

	if withBorder {
		barStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder())
		if highlight {
			barStyle = barStyle.BorderForeground(styles.AccentColor)
		} else {
			barStyle = barStyle.BorderForeground(styles.NeutralColor)
		}
		return barStyle.Render(bar + " " + label)
	}
	return bar
}

// Format a title with proper styling
func formatTitle(text string) string {
	return styles.TitleStyle.Render(text)
}

// Format a subtitle with proper styling
func formatSubTitle(text string) string {
	return styles.SubtitleStyle.Render(text)
}

// Create a key hint for help text
func formatKeyHint(key, description string) string {
	keyStyle := lipgloss.NewStyle().
		Foreground(styles.PrimaryColor).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(styles.NeutralColor)

	return fmt.Sprintf("%s %s", keyStyle.Render(key), descStyle.Render(description))
}

// Join key hints with dot separators
func formatKeyHints(hints ...string) string {
	return strings.Join(hints, "  · ")
}

// Create a data row with label and value
func formatTableRow(label string, value string, highlight bool) string {
	var labelStyle, valueStyle lipgloss.Style

	if highlight {
		labelStyle = styles.HighlightStyle
		valueStyle = styles.HighlightStyle
	} else {
		labelStyle = styles.BaseStyle
		valueStyle = styles.BaseStyle
	}

	// Always align colon at the same position, and right-align value
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		labelStyle.Width(20).Render(label),
		valueStyle.Width(14).Align(lipgloss.Right).Render(value),
	)
}

// Create simple horizontal tabs
func createTabs(tabs []string, activeIdx int, width int) string {
	var renderedTabs []string

	for i, tab := range tabs {
		var tabStyle lipgloss.Style

		if i == activeIdx {
			tabStyle = styles.ActiveTabStyle
		} else {
			tabStyle = styles.TabStyle
		}

		renderedTabs = append(renderedTabs, tabStyle.Render(tab))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, renderedTabs...)
}

// Parse a float with a default fallback value
func parseFloatWithDefault(s string, defaultVal float64) (float64, error) {
	if s == "" {
		return defaultVal, nil
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultVal, err
	}

	return val, nil
}

// Parse an int with a default fallback value
func parseIntWithDefault(s string, defaultVal int) (int, error) {
	if s == "" {
		return defaultVal, nil
	}

	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal, err
	}

	return val, nil
}

// Format tax comparison results for display
func formatComparisonResults(results []models.TaxResult, currentIncome float64) string {
	var sb strings.Builder

	sb.WriteString("\n")

	// Create a header row
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.PrimaryColor).
		Bold(true)

	sb.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Left,
		headerStyle.Width(15).Render("Income"),
		headerStyle.Width(15).Render("Tax Amount"),
		headerStyle.Width(15).Render("Tax Rate"),
		headerStyle.Width(44).Render("Comparison"),
	))
	sb.WriteString("\n\n")

	barWidth := 36

	// Create rows for each result
	for _, result := range results {
		isCurrentIncome := result.Income == currentIncome

		incomeStr := formatEuro(result.Income)
		taxStr := formatEuro(result.TotalTax)
		rateStr := formatPercent(result.TaxRate)

		bar := createProgressBar(result.TaxRate, barWidth, isCurrentIncome, rateStr, false)

		rowStyle := styles.BaseStyle
		if isCurrentIncome {
			rowStyle = styles.HighlightStyle
		}

		sb.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Left,
			rowStyle.Width(15).Render(incomeStr),
			rowStyle.Width(15).Render(taxStr),
			rowStyle.Width(15).Render(rateStr),
			bar,
		))

		if isCurrentIncome {
			sb.WriteString(" ← Your income")
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

// Format tax results for display
func formatTaxResults(income, incomeTax, solidarityTax, totalTax, netIncome, taxRate float64) string {
	var sb strings.Builder

	// Calculate tax breakdown percentages
	incomeTaxPercent := (incomeTax / income) * 100
	solidarityTaxPercent := (solidarityTax / income) * 100
	netIncomePercent := (netIncome / income) * 100

	// Unified two-column layout for all rows
	sb.WriteString("\n")
	sb.WriteString(formatTableRow("Annual Income:", formatEuro(income), false))
	sb.WriteString("\n")
	sb.WriteString(formatTableRow("Income Tax:", formatEuro(incomeTax), false))
	sb.WriteString("\n")
	sb.WriteString(formatTableRow("Solidarity Tax:", formatEuro(solidarityTax), false))
	sb.WriteString("\n\n") // Add a blank line between sections
	sb.WriteString(lipgloss.NewStyle().
		Foreground(styles.NeutralColor).
		Render(strings.Repeat("─", 45)))
	sb.WriteString("\n")
	sb.WriteString(formatTableRow("Total Tax:", formatEuro(totalTax), true))
	sb.WriteString("\n")
	sb.WriteString(formatTableRow("Net Income:", formatEuro(netIncome), true))
	sb.WriteString("\n")
	sb.WriteString(formatTableRow("Effective Tax Rate:", formatPercent(taxRate), true))
	sb.WriteString("\n\n")

	// Visual breakdown - cleaner visualization
	sb.WriteString(formatSubTitle("Tax Breakdown"))
	sb.WriteString("\n\n")

	// Styled bars with consistent width
	barWidth := 40

	// Helper for aligned breakdown row
	breakdownRow := func(label, value string, percent float64, highlight bool) string {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.BaseStyle.Width(16).Render(label),
			styles.BaseStyle.Width(7).Align(lipgloss.Right).Render(value),
			styles.BaseStyle.Width(2).Render(""), // Spacer column
			createProgressBar(percent, barWidth, highlight, "", false),
		)
	}

	// Income tax row
	sb.WriteString(breakdownRow("Income Tax:", formatPercent(incomeTaxPercent), incomeTaxPercent, false))
	sb.WriteString("\n\n")
	// Solidarity row
	sb.WriteString(breakdownRow("Solidarity:", formatPercent(solidarityTaxPercent), solidarityTaxPercent, false))
	sb.WriteString("\n\n")
	// Net income row
	sb.WriteString(breakdownRow("Net Income:", formatPercent(netIncomePercent), netIncomePercent, true))
	sb.WriteString("\n\n")

	// Monthly breakdown - cleaner layout
	sb.WriteString(formatSubTitle("Monthly Breakdown"))
	sb.WriteString("\n\n")

	monthlyIncome := income / 12
	monthlyTax := totalTax / 12
	monthlyNet := netIncome / 12

	sb.WriteString(formatTableRow("Monthly Income:", formatEuro(monthlyIncome), false))
	sb.WriteString("\n")
	sb.WriteString(formatTableRow("Monthly Tax:", formatEuro(monthlyTax), false))
	sb.WriteString("\n")
	sb.WriteString(formatTableRow("Monthly Net:", formatEuro(monthlyNet), true))

	return sb.String()
}
