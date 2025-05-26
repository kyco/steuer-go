package components

import (
	"fmt"
	"strings"

	"tax-calculator/internal/tax/views/styles"

	"github.com/charmbracelet/lipgloss"
)

type TaxResult struct {
	Income           float64
	IncomeTax        float64
	SolidarityTax    float64
	TotalTax         float64
	NetIncome        float64
	EffectiveTaxRate float64
	TaxClass         int
	Year             string
}

type ResultsDashboard struct {
	Result         TaxResult
	ShowComparison bool
	ComparisonData []TaxResult
	Width          int
}

func NewResultsDashboard(result TaxResult) ResultsDashboard {
	return ResultsDashboard{
		Result:         result,
		ShowComparison: false,
		Width:          80,
	}
}

func (rd *ResultsDashboard) View() string {
	var builder strings.Builder

	builder.WriteString(rd.renderHeader())
	builder.WriteString("\n\n")
	builder.WriteString(rd.renderKeyMetrics())
	builder.WriteString("\n\n")
	builder.WriteString(rd.renderVisualBreakdown())
	builder.WriteString("\n\n")
	builder.WriteString(rd.renderMonthlyBreakdown())

	if rd.ShowComparison && len(rd.ComparisonData) > 0 {
		builder.WriteString("\n\n")
		builder.WriteString(rd.renderComparison())
	}

	return builder.String()
}

func (rd *ResultsDashboard) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.PrimaryColor).
		Bold(true).
		Align(lipgloss.Center).
		Width(rd.Width)

	return headerStyle.Render(fmt.Sprintf("🧮 Tax Calculation Results for %s", rd.Result.Year))
}

func (rd *ResultsDashboard) renderKeyMetrics() string {
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.PrimaryColor).
		Padding(1, 2).
		Width(rd.Width/3 - 2)

	grossCard := cardStyle.Render(fmt.Sprintf(
		"%s\n%s\n%s",
		lipgloss.NewStyle().Foreground(styles.NeutralColor).Render("Gross Income"),
		lipgloss.NewStyle().Foreground(styles.FgColor).Bold(true).Render(formatEuro(rd.Result.Income)),
		lipgloss.NewStyle().Foreground(styles.SuccessColor).Render("Annual"),
	))

	taxCard := cardStyle.Render(fmt.Sprintf(
		"%s\n%s\n%s",
		lipgloss.NewStyle().Foreground(styles.NeutralColor).Render("Total Tax"),
		lipgloss.NewStyle().Foreground(styles.DangerColor).Bold(true).Render(formatEuro(rd.Result.TotalTax)),
		lipgloss.NewStyle().Foreground(styles.DangerColor).Render(formatPercent(rd.Result.EffectiveTaxRate)),
	))

	netCard := cardStyle.Render(fmt.Sprintf(
		"%s\n%s\n%s",
		lipgloss.NewStyle().Foreground(styles.NeutralColor).Render("Net Income"),
		lipgloss.NewStyle().Foreground(styles.SuccessColor).Bold(true).Render(formatEuro(rd.Result.NetIncome)),
		lipgloss.NewStyle().Foreground(styles.SuccessColor).Render("Take Home"),
	))

	return lipgloss.JoinHorizontal(lipgloss.Top, grossCard, "  ", taxCard, "  ", netCard)
}

func (rd *ResultsDashboard) renderVisualBreakdown() string {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.AccentColor).
		Bold(true)
	builder.WriteString(titleStyle.Render("💰 Tax Breakdown"))
	builder.WriteString("\n\n")

	barWidth := rd.Width - 30

	incomeTaxPercent := (rd.Result.IncomeTax / rd.Result.Income) * 100
	solidarityPercent := (rd.Result.SolidarityTax / rd.Result.Income) * 100
	netPercent := (rd.Result.NetIncome / rd.Result.Income) * 100

	builder.WriteString(rd.renderBreakdownRow("Income Tax", rd.Result.IncomeTax, incomeTaxPercent, barWidth, styles.DangerColor))
	builder.WriteString("\n")
	builder.WriteString(rd.renderBreakdownRow("Solidarity Tax", rd.Result.SolidarityTax, solidarityPercent, barWidth, styles.WarningColor))
	builder.WriteString("\n")
	builder.WriteString(rd.renderBreakdownRow("Net Income", rd.Result.NetIncome, netPercent, barWidth, styles.SuccessColor))

	return builder.String()
}

func (rd *ResultsDashboard) renderBreakdownRow(label string, amount, percent float64, barWidth int, color lipgloss.Color) string {
	labelStyle := lipgloss.NewStyle().Width(15)
	amountStyle := lipgloss.NewStyle().Width(12).Align(lipgloss.Right).Foreground(color).Bold(true)
	percentStyle := lipgloss.NewStyle().Width(8).Align(lipgloss.Right).Foreground(color)

	filled := int((percent / 100) * float64(barWidth))
	bar := strings.Repeat("━", filled) + strings.Repeat("─", barWidth-filled)
	barStyle := lipgloss.NewStyle().Foreground(color)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		labelStyle.Render(label),
		amountStyle.Render(formatEuro(amount)),
		percentStyle.Render(formatPercent(percent)),
		"  ",
		barStyle.Render(bar),
	)
}

func (rd *ResultsDashboard) renderMonthlyBreakdown() string {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.AccentColor).
		Bold(true)
	builder.WriteString(titleStyle.Render("📅 Monthly Breakdown"))
	builder.WriteString("\n\n")

	monthlyIncome := rd.Result.Income / 12
	monthlyTax := rd.Result.TotalTax / 12
	monthlyNet := rd.Result.NetIncome / 12

	tableStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.NeutralColor).
		Padding(1)

	table := fmt.Sprintf(
		"%-15s %12s\n"+
			"%-15s %12s\n"+
			"%-15s %12s\n"+
			"%s\n"+
			"%-15s %12s",
		"Gross Monthly:", formatEuro(monthlyIncome),
		"Tax Monthly:", formatEuro(monthlyTax),
		"Net Monthly:", formatEuro(monthlyNet),
		strings.Repeat("─", 28),
		"Daily Net:", formatEuro(monthlyNet/30),
	)

	builder.WriteString(tableStyle.Render(table))

	return builder.String()
}

func (rd *ResultsDashboard) renderComparison() string {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.AccentColor).
		Bold(true)
	builder.WriteString(titleStyle.Render("📊 Tax Rate Comparison"))
	builder.WriteString("\n\n")

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.PrimaryColor).
		Bold(true)

	builder.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Left,
		headerStyle.Width(15).Render("Income"),
		headerStyle.Width(12).Render("Tax Rate"),
		headerStyle.Width(40).Render("Comparison"),
	))
	builder.WriteString("\n\n")

	for _, comparison := range rd.ComparisonData {
		isCurrentIncome := comparison.Income == rd.Result.Income

		incomeStr := formatEuro(comparison.Income)
		rateStr := formatPercent(comparison.EffectiveTaxRate)

		style := styles.BaseStyle
		if isCurrentIncome {
			style = styles.HighlightStyle
		}

		barWidth := 35
		filled := int((comparison.EffectiveTaxRate / 50) * float64(barWidth))
		bar := strings.Repeat("━", filled) + strings.Repeat("─", barWidth-filled)

		line := lipgloss.JoinHorizontal(
			lipgloss.Left,
			style.Width(15).Render(incomeStr),
			style.Width(12).Render(rateStr),
			style.Render(bar),
		)

		if isCurrentIncome {
			line += " ← Your Income"
		}

		builder.WriteString(line)
		builder.WriteString("\n")
	}

	return builder.String()
}

func (rd *ResultsDashboard) SetComparison(data []TaxResult) {
	rd.ComparisonData = data
	rd.ShowComparison = true
}

func (rd *ResultsDashboard) ToggleComparison() {
	rd.ShowComparison = !rd.ShowComparison
}

func formatEuro(amount float64) string {
	return fmt.Sprintf("€ %.2f", amount)
}

func formatPercent(percent float64) string {
	return fmt.Sprintf("%.1f%%", percent)
}
