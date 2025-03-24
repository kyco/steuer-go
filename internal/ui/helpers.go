package ui

import (
	"fmt"
	"strconv"
	"strings"
	
	"github.com/charmbracelet/lipgloss"
	
	"tax-calculator/internal/adapters/api"
	"tax-calculator/internal/ui/styles"
)

func (m *AppModel) updateFocus() {
	m.incomeInput.Blur()
	m.yearInput.Blur()

	switch m.focusField {
	case IncomeField:
		m.incomeInput.Focus()
	case YearField:
		m.yearInput.Focus()
	}
}

func (m *AppModel) validateAndCalculate() (bool, string) {
	if strings.TrimSpace(m.incomeInput.Value()) == "" {
		return false, "Income cannot be empty"
	}
	
	income, err := strconv.ParseFloat(strings.TrimSpace(m.incomeInput.Value()), 64)
	if err != nil || income <= 0 {
		return false, "Income must be a positive number"
	}
	
	year := strings.TrimSpace(m.yearInput.Value())
	if year == "" {
		m.yearInput.SetValue("2025")
	} else {
		yearVal, err := strconv.Atoi(year)
		if err != nil || yearVal < 2024 || yearVal > 2030 {
			return false, "Year must be between 2024 and 2030"
		}
	}
	
	return true, ""
}

func (m *AppModel) updateResultsContent() {
	if m.result == nil {
		return
	}

	var content strings.Builder

	var incomeTax, solidarityTax string
	for _, output := range m.result.Outputs.Output {
		if output.Name == "LSTLZZ" {
			incomeTax = output.Value
		} else if output.Name == "SOLZLZZ" {
			solidarityTax = output.Value
		}
	}

	incomeTaxEuros := float64(api.MustParseInt(incomeTax)) / 100
	solidarityTaxEuros := float64(api.MustParseInt(solidarityTax)) / 100
	totalTax := incomeTaxEuros + solidarityTaxEuros
	
	income, _ := strconv.ParseFloat(m.incomeInput.Value(), 64)
	netIncome := income - totalTax
	taxPercentage := (totalTax / income) * 100
	
	title := lipgloss.NewStyle().
		Foreground(styles.PrimaryColor).
		Bold(true).
		Render("Tax Calculation Results")
	fmt.Fprintf(&content, "%s\n", title)
	
	summaryStyle := lipgloss.NewStyle().
		Border(styles.MinimalBorder).
		BorderForeground(styles.SecondaryColor).
		Padding(0, 1)
	
	rightAlignedValue := func(label string, value string) string {
		return fmt.Sprintf("%-20s %s", label, lipgloss.NewStyle().Align(lipgloss.Right).Width(15).Render(value))
	}
	
	var inputSummary strings.Builder
	fmt.Fprintf(&inputSummary, "%s\n", rightAlignedValue("Income:", fmt.Sprintf("€%.2f", income)))
	fmt.Fprintf(&inputSummary, "%s\n", rightAlignedValue("Tax Class:", fmt.Sprintf("%d", m.selectedTaxClass)))
	fmt.Fprintf(&inputSummary, "%s", rightAlignedValue("Year:", m.yearInput.Value()))
	
	fmt.Fprintf(&content, "%s\n", summaryStyle.Render(inputSummary.String()))
	
	var resultsSummary strings.Builder
	fmt.Fprintf(&resultsSummary, "%s\n", rightAlignedValue("Income Tax:", fmt.Sprintf("€%.2f", incomeTaxEuros)))
	fmt.Fprintf(&resultsSummary, "%s\n", rightAlignedValue("Solidarity Tax:", fmt.Sprintf("€%.2f", solidarityTaxEuros)))
	fmt.Fprintf(&resultsSummary, "%s\n", rightAlignedValue("Total Tax:", fmt.Sprintf("€%.2f", totalTax)))
	fmt.Fprintf(&resultsSummary, "%s\n", rightAlignedValue("Net Income:", fmt.Sprintf("€%.2f", netIncome)))
	fmt.Fprintf(&resultsSummary, "%s", rightAlignedValue("Tax Rate:", fmt.Sprintf("%.2f%%", taxPercentage)))
	
	fmt.Fprintf(&content, "%s\n", lipgloss.NewStyle().
		Border(styles.MinimalBorder).
		BorderForeground(styles.PrimaryColor).
		Padding(0, 1).
		Render(resultsSummary.String()))
	
	monthlyIncome := income / 12
	monthlyTax := totalTax / 12
	monthlyNet := monthlyIncome - monthlyTax
	
	var monthlySummary strings.Builder
	fmt.Fprintf(&monthlySummary, "%s\n", rightAlignedValue("Monthly Income:", fmt.Sprintf("€%.2f", monthlyIncome)))
	fmt.Fprintf(&monthlySummary, "%s\n", rightAlignedValue("Monthly Tax:", fmt.Sprintf("€%.2f", monthlyTax)))
	fmt.Fprintf(&monthlySummary, "%s", rightAlignedValue("Monthly Net Income:", fmt.Sprintf("€%.2f", monthlyNet)))
	
	fmt.Fprintf(&content, "%s\n", summaryStyle.Render(monthlySummary.String()))
	
	barWidth := 40
	netRatio := netIncome / income
	netChars := int(netRatio * float64(barWidth))
	taxChars := barWidth - netChars
	
	netBar := strings.Repeat("█", netChars)
	taxBar := strings.Repeat("▒", taxChars)
	
	barChart := lipgloss.NewStyle().Foreground(styles.SuccessColor).Render(netBar) + 
		lipgloss.NewStyle().Foreground(styles.DangerColor).Render(taxBar)
	
	fmt.Fprintf(&content, "%s\n", barChart)
	fmt.Fprintf(&content, "%s %s\n", 
		lipgloss.NewStyle().Foreground(styles.SuccessColor).Render("█ Net Income"),
		lipgloss.NewStyle().Foreground(styles.DangerColor).Render("▒ Tax"))
	
	compareText := lipgloss.NewStyle().
		Foreground(styles.AccentColor).
		Bold(true).
		Render("Press 'c' to compare with other income levels")
	fmt.Fprintf(&content, "\n%s\n", compareText)
	
	if m.showDetails {
		detailsTitle := lipgloss.NewStyle().
			Foreground(styles.PrimaryColor).
			Render("Detailed Tax Information")
			fmt.Fprintf(&content, "%s\n", detailsTitle)

		fmt.Fprintf(&content, "%-10s %-14s %-10s\n", 
			lipgloss.NewStyle().Foreground(styles.AccentColor).Bold(true).Render("CODE"),
			lipgloss.NewStyle().Foreground(styles.AccentColor).Bold(true).Render("VALUE (EUR)"),
			lipgloss.NewStyle().Foreground(styles.AccentColor).Bold(true).Render("TYPE"))
		fmt.Fprintf(&content, "%s\n", 
			lipgloss.NewStyle().Foreground(styles.SecondaryColor).Render(strings.Repeat("─", 40)))

		for _, output := range m.result.Outputs.Output {
			valueInEuros := float64(api.MustParseInt(output.Value)) / 100
			fmt.Fprintf(&content, "%-10s %-14.2f %-10s\n",
				output.Name,
				valueInEuros,
				output.Type)
		}
	} else {
		helpText := lipgloss.NewStyle().
			Foreground(styles.NeutralColor).
			Italic(true).
			Render("Press 'd' to toggle detailed tax information")
		fmt.Fprintf(&content, "%s\n", helpText)
	}

	m.resultsViewport.SetContent(content.String())
	m.resultsViewport.GotoTop()
}

func (m *AppModel) updateComparisonContent() {
	if m.comparisonResults == nil || len(m.comparisonResults) == 0 {
		return
	}

	var content strings.Builder

	title := lipgloss.NewStyle().
		Foreground(styles.PrimaryColor).
		Bold(true).
		Render("Income and Tax Comparison")
	fmt.Fprintf(&content, "%s\n\n", title)
	
	originalIncome, _ := strconv.ParseFloat(m.incomeInput.Value(), 64)
	
	description := fmt.Sprintf("Showing tax calculations from half (€%.2f) to double (€%.2f) your income (€%.2f)",
		originalIncome/2, originalIncome*2, originalIncome)
	
	descriptionStyle := lipgloss.NewStyle().
		Foreground(styles.SecondaryColor).
		Italic(true)
	
	fmt.Fprintf(&content, "%s\n\n", descriptionStyle.Render(description))
	
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.AccentColor).
		Bold(true)
	
	incomeWidth := 16
	taxWidth := 16
	netWidth := 16
	rateWidth := 10
	
	fmt.Fprintf(&content, "%s %s %s %s\n",
		headerStyle.Width(incomeWidth).Render("Income"),
		headerStyle.Width(taxWidth).Render("Total Tax"),
		headerStyle.Width(netWidth).Render("Net Income"),
		headerStyle.Width(rateWidth).Render("Tax Rate"))
	
	fmt.Fprintf(&content, "%s\n", 
		lipgloss.NewStyle().
			Foreground(styles.SecondaryColor).
			Render(strings.Repeat("─", incomeWidth+taxWidth+netWidth+rateWidth+3)))
	
	originalValueStyle := lipgloss.NewStyle().
		Foreground(styles.BgColor).
		Background(styles.PrimaryColor).
		Bold(true)
	
	normalValueStyle := lipgloss.NewStyle().
		Foreground(styles.FgColor)
		
	for _, result := range m.comparisonResults {
		if result.Error != nil {
			continue
		}
		
		incomeStr := fmt.Sprintf("€%.2f", result.Income)
		taxStr := fmt.Sprintf("€%.2f", result.TotalTax)
		netStr := fmt.Sprintf("€%.2f", result.NetIncome)
		rateStr := fmt.Sprintf("%.2f%%", result.TaxRate)
		
		style := normalValueStyle
		if result.Income >= originalIncome*0.99 && result.Income <= originalIncome*1.01 {
			style = originalValueStyle
		}
		
		fmt.Fprintf(&content, "%s %s %s %s\n",
			style.Width(incomeWidth).Align(lipgloss.Right).Render(incomeStr),
			style.Width(taxWidth).Align(lipgloss.Right).Render(taxStr),
			style.Width(netWidth).Align(lipgloss.Right).Render(netStr),
			style.Width(rateWidth).Align(lipgloss.Right).Render(rateStr))
	}
	
	fmt.Fprintf(&content, "\n%s\n", headerStyle.Render("Tax Rate Progression"))
	fmt.Fprintf(&content, "%s\n", 
		lipgloss.NewStyle().
			Foreground(styles.SecondaryColor).
			Render(strings.Repeat("─", 60)))
	
	barWidth := 40
	for _, result := range m.comparisonResults {
		if result.Error != nil {
			continue
		}
		
		taxRateChars := int((result.TaxRate / 60) * float64(barWidth))
		if taxRateChars > barWidth {
			taxRateChars = barWidth
		}
		
		taxBar := strings.Repeat("█", taxRateChars)
		
		style := lipgloss.NewStyle().Foreground(styles.DangerColor)
		labelStyle := lipgloss.NewStyle().Foreground(styles.FgColor)
		
		if result.Income >= originalIncome*0.99 && result.Income <= originalIncome*1.01 {
			style = lipgloss.NewStyle().Foreground(styles.PrimaryColor)
			labelStyle = lipgloss.NewStyle().Foreground(styles.PrimaryColor).Bold(true)
		}
		
		incomeLabel := fmt.Sprintf("€%.0f", result.Income)
		
		fmt.Fprintf(&content, "%s %s %s\n",
			labelStyle.Width(10).Align(lipgloss.Right).Render(incomeLabel),
			style.Render(taxBar),
			labelStyle.Render(fmt.Sprintf(" %.2f%%", result.TaxRate)))
	}
	
	helpText := lipgloss.NewStyle().
		Foreground(styles.NeutralColor).
		Italic(true).
		Render("Press 'b' to go back to regular results")
	fmt.Fprintf(&content, "\n%s\n", helpText)

	m.comparisonViewport.SetContent(content.String())
	m.comparisonViewport.GotoTop()
}

func parseIncome(s string) float64 {
	income, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0
	}
	return income
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}