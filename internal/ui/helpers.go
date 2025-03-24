package ui

import (
	"fmt"
	"strconv"
	"strings"
	
	"github.com/charmbracelet/lipgloss"
	
	"tax-calculator/internal/adapters/api"
	"tax-calculator/internal/ui/styles"
)

// updateFocus updates the focus state of form fields
func (m *AppModel) updateFocus() {
	// First blur everything
	m.incomeInput.Blur()
	m.yearInput.Blur()

	// Then focus the active field
	switch m.focusField {
	case IncomeField:
		m.incomeInput.Focus()
	case YearField:
		m.yearInput.Focus()
	}
}

// validateAndCalculate validates the form inputs and initiates calculation
func (m *AppModel) validateAndCalculate() (bool, string) {
	// Validate income
	if strings.TrimSpace(m.incomeInput.Value()) == "" {
		return false, "Income cannot be empty"
	}
	
	income, err := strconv.ParseFloat(strings.TrimSpace(m.incomeInput.Value()), 64)
	if err != nil || income <= 0 {
		return false, "Income must be a positive number"
	}
	
	// Validate year
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

// updateResultsContent updates the viewport content based on calculation results
func (m *AppModel) updateResultsContent() {
	if m.result == nil {
		return
	}

	var content strings.Builder

	// Extract tax values
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
	
	// Get input values
	income, _ := strconv.ParseFloat(m.incomeInput.Value(), 64)
	netIncome := income - totalTax
	taxPercentage := (totalTax / income) * 100
	
	// Debug logs section if we have any
	if len(m.debugMessages) > 0 {
		debugTitle := lipgloss.NewStyle().
			Foreground(styles.WarningColor).
			Bold(true).
			Render("Debug Messages:")
		fmt.Fprintf(&content, "%s\n", debugTitle)
		
		// Display up to the last 20 debug messages to avoid overwhelming the display
		// and to ensure the most recent messages are visible
		startIdx := 0
		if len(m.debugMessages) > 20 {
			startIdx = len(m.debugMessages) - 20
		}
		
		// Show message count if some are hidden
		if startIdx > 0 {
			hiddenMsg := fmt.Sprintf("(Showing last 20 of %d messages)", len(m.debugMessages))
			fmt.Fprintf(&content, "%s\n", 
				lipgloss.NewStyle().
					Foreground(styles.NeutralColor).
					Italic(true).
					Render(hiddenMsg))
		}
		
		// Process each debug message with appropriate styling
		for i := startIdx; i < len(m.debugMessages); i++ {
			msg := m.debugMessages[i]
			
			// Style messages differently based on content
			style := lipgloss.NewStyle().Foreground(styles.WarningColor)
			
			// Success messages get green color
			if strings.Contains(msg, "✓ Success") {
				style = lipgloss.NewStyle().Foreground(styles.SuccessColor)
			}
			
			// Failure messages get red color
			if strings.Contains(msg, "✗ Failed") {
				style = lipgloss.NewStyle().Foreground(styles.DangerColor)
			}
			
			// Progress messages get neutral color
			if strings.Contains(msg, "Completed") {
				style = lipgloss.NewStyle().Foreground(styles.NeutralColor)
			}
			
			fmt.Fprintf(&content, "%s\n", style.Render(msg))
		}
		
		// Add divider
		fmt.Fprintf(&content, "%s\n\n", 
			lipgloss.NewStyle().
				Foreground(styles.WarningColor).
				Render(strings.Repeat("─", 60)))
	}
	
	// Create summary header
	title := lipgloss.NewStyle().
		Foreground(styles.PrimaryColor).
		Bold(true).
		Render("Tax Calculation Results")
	fmt.Fprintf(&content, "%s\n", title)
	
	// Create compact summary tables with right-aligned values
	summaryStyle := lipgloss.NewStyle().
		Border(styles.MinimalBorder).
		BorderForeground(styles.SecondaryColor).
		Padding(0, 1)
	
	// Right-align formatting helper
	rightAlignedValue := func(label string, value string) string {
		return fmt.Sprintf("%-20s %s", label, lipgloss.NewStyle().Align(lipgloss.Right).Width(15).Render(value))
	}
	
	// Input summary with right-aligned values
	var inputSummary strings.Builder
	fmt.Fprintf(&inputSummary, "%s\n", rightAlignedValue("Income:", fmt.Sprintf("€%.2f", income)))
	fmt.Fprintf(&inputSummary, "%s\n", rightAlignedValue("Tax Class:", fmt.Sprintf("%d", m.selectedTaxClass)))
	fmt.Fprintf(&inputSummary, "%s", rightAlignedValue("Year:", m.yearInput.Value()))
	
	fmt.Fprintf(&content, "%s\n", summaryStyle.Render(inputSummary.String()))
	
	// Results summary with right-aligned values
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
	
	// Monthly breakdown with right-aligned values
	monthlyIncome := income / 12
	monthlyTax := totalTax / 12
	monthlyNet := monthlyIncome - monthlyTax
	
	var monthlySummary strings.Builder
	fmt.Fprintf(&monthlySummary, "%s\n", rightAlignedValue("Monthly Income:", fmt.Sprintf("€%.2f", monthlyIncome)))
	fmt.Fprintf(&monthlySummary, "%s\n", rightAlignedValue("Monthly Tax:", fmt.Sprintf("€%.2f", monthlyTax)))
	fmt.Fprintf(&monthlySummary, "%s", rightAlignedValue("Monthly Net Income:", fmt.Sprintf("€%.2f", monthlyNet)))
	
	fmt.Fprintf(&content, "%s\n", summaryStyle.Render(monthlySummary.String()))
	
	// Display compact visualization
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
	
	// Add a "Compare" option
	compareText := lipgloss.NewStyle().
		Foreground(styles.AccentColor).
		Bold(true).
		Render("Press 'c' to compare with other income levels")
	fmt.Fprintf(&content, "\n%s\n", compareText)
	
	// Detailed breakdown if requested
	if m.showDetails {
		detailsTitle := lipgloss.NewStyle().
			Foreground(styles.PrimaryColor).
			Render("Detailed Tax Information")
			fmt.Fprintf(&content, "%s\n", detailsTitle)

		// Format output in columns
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
	m.resultsViewport.GotoTop() // Ensure we start at the top
}

// updateComparisonContent updates the comparison viewport content
func (m *AppModel) updateComparisonContent() {
	// Safety check in case the comparison step is triggered but no results are available
	if m.comparisonResults == nil || len(m.comparisonResults) == 0 {
		return
	}

	var content strings.Builder

	// Create comparison header
	title := lipgloss.NewStyle().
		Foreground(styles.PrimaryColor).
		Bold(true).
		Render("Income and Tax Comparison")
	fmt.Fprintf(&content, "%s\n\n", title)
	
	// Debug logs section if we have any
	if len(m.debugMessages) > 0 {
		debugTitle := lipgloss.NewStyle().
			Foreground(styles.WarningColor).
			Bold(true).
			Render("Debug Messages:")
		fmt.Fprintf(&content, "%s\n", debugTitle)
		
		// Display up to the last 20 debug messages to avoid overwhelming the display
		// and to ensure the most recent messages are visible
		startIdx := 0
		if len(m.debugMessages) > 20 {
			startIdx = len(m.debugMessages) - 20
		}
		
		// Show message count if some are hidden
		if startIdx > 0 {
			hiddenMsg := fmt.Sprintf("(Showing last 20 of %d messages)", len(m.debugMessages))
			fmt.Fprintf(&content, "%s\n", 
				lipgloss.NewStyle().
					Foreground(styles.NeutralColor).
					Italic(true).
					Render(hiddenMsg))
		}
		
		// Process each debug message with appropriate styling
		for i := startIdx; i < len(m.debugMessages); i++ {
			msg := m.debugMessages[i]
			
			// Style messages differently based on content
			style := lipgloss.NewStyle().Foreground(styles.WarningColor)
			
			// Success messages get green color
			if strings.Contains(msg, "✓ Success") {
				style = lipgloss.NewStyle().Foreground(styles.SuccessColor)
			}
			
			// Failure messages get red color
			if strings.Contains(msg, "✗ Failed") {
				style = lipgloss.NewStyle().Foreground(styles.DangerColor)
			}
			
			// Progress messages get neutral color
			if strings.Contains(msg, "Completed") {
				style = lipgloss.NewStyle().Foreground(styles.NeutralColor)
			}
			
			fmt.Fprintf(&content, "%s\n", style.Render(msg))
		}
		
		// Add divider
		fmt.Fprintf(&content, "%s\n\n", 
			lipgloss.NewStyle().
				Foreground(styles.WarningColor).
				Render(strings.Repeat("─", 60)))
	}
	
	// Get original input values
	originalIncome, _ := strconv.ParseFloat(m.incomeInput.Value(), 64)
	
	// Create description text
	description := fmt.Sprintf("Showing tax calculations from half (€%.2f) to double (€%.2f) your income (€%.2f)",
		originalIncome/2, originalIncome*2, originalIncome)
	
	descriptionStyle := lipgloss.NewStyle().
		Foreground(styles.SecondaryColor).
		Italic(true)
	
	fmt.Fprintf(&content, "%s\n\n", descriptionStyle.Render(description))
	
	// Create table headers
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.AccentColor).
		Bold(true)
	
	// Table column widths
	incomeWidth := 16
	taxWidth := 16
	netWidth := 16
	rateWidth := 10
	
	// Table header
	fmt.Fprintf(&content, "%s %s %s %s\n",
		headerStyle.Width(incomeWidth).Render("Income"),
		headerStyle.Width(taxWidth).Render("Total Tax"),
		headerStyle.Width(netWidth).Render("Net Income"),
		headerStyle.Width(rateWidth).Render("Tax Rate"))
	
	// Table divider
	fmt.Fprintf(&content, "%s\n", 
		lipgloss.NewStyle().
			Foreground(styles.SecondaryColor).
			Render(strings.Repeat("─", incomeWidth+taxWidth+netWidth+rateWidth+3)))
	
	// Table data with highlighting for original value
	originalValueStyle := lipgloss.NewStyle().
		Foreground(styles.BgColor).
		Background(styles.PrimaryColor).
		Bold(true)
	
	normalValueStyle := lipgloss.NewStyle().
		Foreground(styles.FgColor)
		
	// Print each row, highlighting the row that matches the original income
	for _, result := range m.comparisonResults {
		// Skip rows with errors
		if result.Error != nil {
			continue
		}
		
		// Format values
		incomeStr := fmt.Sprintf("€%.2f", result.Income)
		taxStr := fmt.Sprintf("€%.2f", result.TotalTax)
		netStr := fmt.Sprintf("€%.2f", result.NetIncome)
		rateStr := fmt.Sprintf("%.2f%%", result.TaxRate)
		
		// Determine style based on whether this is close to the original income
		style := normalValueStyle
		if result.Income >= originalIncome*0.99 && result.Income <= originalIncome*1.01 {
			style = originalValueStyle
		}
		
		// Print row
		fmt.Fprintf(&content, "%s %s %s %s\n",
			style.Width(incomeWidth).Align(lipgloss.Right).Render(incomeStr),
			style.Width(taxWidth).Align(lipgloss.Right).Render(taxStr),
			style.Width(netWidth).Align(lipgloss.Right).Render(netStr),
			style.Width(rateWidth).Align(lipgloss.Right).Render(rateStr))
	}
	
	// Visual tax rate progression
	fmt.Fprintf(&content, "\n%s\n", headerStyle.Render("Tax Rate Progression"))
	fmt.Fprintf(&content, "%s\n", 
		lipgloss.NewStyle().
			Foreground(styles.SecondaryColor).
			Render(strings.Repeat("─", 60)))
	
	// Create a visual representation of tax rates
	barWidth := 40
	for _, result := range m.comparisonResults {
		// Skip rows with errors
		if result.Error != nil {
			continue
		}
		
		// Create a progress bar for the tax rate
		taxRateChars := int((result.TaxRate / 60) * float64(barWidth)) // Assume max 60% tax rate
		if taxRateChars > barWidth {
			taxRateChars = barWidth
		}
		
		taxBar := strings.Repeat("█", taxRateChars)
		
		// Determine style based on whether this is close to the original income
		style := lipgloss.NewStyle().Foreground(styles.DangerColor)
		labelStyle := lipgloss.NewStyle().Foreground(styles.FgColor)
		
		if result.Income >= originalIncome*0.99 && result.Income <= originalIncome*1.01 {
			style = lipgloss.NewStyle().Foreground(styles.PrimaryColor)
			labelStyle = lipgloss.NewStyle().Foreground(styles.PrimaryColor).Bold(true)
		}
		
		// Format income for label
		incomeLabel := fmt.Sprintf("€%.0f", result.Income)
		
		// Print bar with label and percentage
		fmt.Fprintf(&content, "%s %s %s\n",
			labelStyle.Width(10).Align(lipgloss.Right).Render(incomeLabel),
			style.Render(taxBar),
			labelStyle.Render(fmt.Sprintf(" %.2f%%", result.TaxRate)))
	}
	
	// Help text
	helpText := lipgloss.NewStyle().
		Foreground(styles.NeutralColor).
		Italic(true).
		Render("Press 'b' to go back to regular results")
	fmt.Fprintf(&content, "\n%s\n", helpText)

	m.comparisonViewport.SetContent(content.String())
	m.comparisonViewport.GotoTop() // Ensure we start at the top
}

// Helper functions
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