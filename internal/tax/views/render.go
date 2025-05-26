package views

import (
	"fmt"
	"strings"

	"tax-calculator/internal/tax/bmf"
	"tax-calculator/internal/tax/views/styles"

	"github.com/charmbracelet/lipgloss"
)

// Main View - determines which screen to render
func (m *RetroApp) View() string {
	switch m.screen {
	case MainScreen:
		return m.renderMainScreen()
	case AdvancedScreen:
		return m.renderAdvancedScreen()
	case ResultsScreen:
		return m.renderResultsScreen()
	case ComparisonScreen:
		return m.renderComparisonScreen()
	default:
		return m.renderMainScreen()
	}
}

// Main screen with minimal, elegant UI
func (m *RetroApp) renderMainScreen() string {
	// Create clean, elegant header
	logoText := "STEUER-GO"
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.PrimaryColor).
		Padding(1, 2).
		Render(logoText)

	subtitle := lipgloss.NewStyle().
		Foreground(styles.NeutralColor).
		Italic(true).
		Render("German Tax Calculator")

	// Create the top section with clean spacing
	topSection := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		subtitle,
	)

	// Tax class selector with icons - cleaner layout
	var taxClassOptions strings.Builder

	// Build styled tax class options with minimal layout
	taxClassTitle := styles.SubtitleStyle.Render("Tax Class")

	for _, option := range m.taxClassOptions {
		classNum := option.Class
		style := styles.UnselectedItemStyle
		indicator := "  "

		if classNum == m.selectedTaxClass {
			style = styles.SelectedItemStyle
			indicator = "• "
		}

		fmt.Fprintf(&taxClassOptions, "%s%s%s %s\n",
			indicator,
			option.Icon,
			style.Render(fmt.Sprintf(" Class %d:", classNum)),
			style.Render(option.Desc))
	}

	// Income input field - clean styling
	incomeTitle := styles.SubtitleStyle.Render("Annual Income")
	incomeField := styles.InputFieldStyle.Render(m.incomeInput.View())
	if m.focusField == IncomeField {
		incomeField = styles.ActiveInputStyle.Render(m.incomeInput.View())
	}

	// Year input field - clean styling
	yearTitle := styles.SubtitleStyle.Render("Tax Year")
	yearField := styles.InputFieldStyle.Render(m.yearInput.View())
	if m.focusField == YearField {
		yearField = styles.ActiveInputStyle.Render(m.yearInput.View())
	}

	// Minimal button styling
	calculateButton := styles.ButtonStyle.Render(" Calculate Tax ")
	if m.focusField == CalculateButtonField {
		calculateButton = styles.SelectedButtonStyle.Render(" Calculate Tax ")
	}

	advancedButton := styles.ButtonStyle.Render(" Advanced Options ")
	if m.focusField == AdvancedButtonField {
		advancedButton = styles.SelectedButtonStyle.Render(" Advanced Options ")
	}

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		calculateButton,
		"  ",
		advancedButton,
	)

	// Subtle mode indicator
	modeText := "Remote Calculation"
	if m.useLocalCalc {
		modeText = "Local Calculation"
	}

	modeIndicator := lipgloss.NewStyle().
		Foreground(styles.NeutralColor).
		Bold(true).
		Italic(true).
		Render(modeText)

	// Clean, minimal help text
	helpText := lipgloss.JoinHorizontal(
		lipgloss.Center,
		formatKeyHint("↑/↓", "Change Tax Class"),
		"  ",
		formatKeyHint("Tab", "Next Field"),
		"  ",
		formatKeyHint("Enter", "Select"),
		"  ",
		formatKeyHint("L", "Toggle Local Mode"),
	)

	// Main content with proper spacing
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		taxClassTitle,
		taxClassOptions.String(),
		"",
		incomeTitle,
		incomeField,
		"",
		yearTitle,
		yearField,
		"",
		lipgloss.JoinHorizontal(lipgloss.Center, buttons),
		"",
	)

	// Clean container styling
	formStyle := styles.ContainerStyle

	// Width handling
	width := m.windowSize.Width
	if width == 0 {
		width = 100
	}

	// Footer with minimal styling
	footer := lipgloss.JoinVertical(
		lipgloss.Center,
		modeIndicator,
		"",
		helpText,
	)

	// Combine all elements with clean spacing
	return lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		topSection,
		"",
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(formStyle.Render(mainContent)),
		"",
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(footer),
	)
}

// Advanced options screen with minimal styling
func (m *RetroApp) renderAdvancedScreen() string {
	// Clean header
	header := formatTitle("Advanced Tax Parameters")
	subtitle := lipgloss.NewStyle().
		Foreground(styles.NeutralColor).
		Italic(true).
		Render("Adjust parameters for accurate calculations")

	// Clean form layout
	var formContent strings.Builder

	for i, field := range m.advancedFields {
		isFocused := m.focusField == field.Field

		// Simplified styling
		labelStyle := styles.SubtitleStyle
		descStyle := lipgloss.NewStyle().
			Foreground(styles.NeutralColor).
			Italic(true)

		inputStyle := styles.InputFieldStyle
		if isFocused {
			inputStyle = styles.ActiveInputStyle
		}

		// Render with clean spacing
		formContent.WriteString(labelStyle.Render(field.Label))
		formContent.WriteString("\n")
		formContent.WriteString(descStyle.Render(field.Description))
		formContent.WriteString("\n")
		formContent.WriteString(inputStyle.Render(field.Model.View()))

		// Spacing between fields
		if i < len(m.advancedFields)-1 {
			formContent.WriteString("\n\n")
		}
	}

	// Minimal button styling
	backButton := styles.ButtonStyle.Render(" Back ")
	if m.focusField == BackButtonField {
		backButton = styles.SelectedButtonStyle.Render(" Back ")
	}

	calculateButton := styles.ButtonStyle.Render(" Calculate with Parameters ")
	if m.focusField == CalculateButtonField {
		calculateButton = styles.SelectedButtonStyle.Render(" Calculate with Parameters ")
	}

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		backButton,
		"  ",
		calculateButton,
	)

	// Set content with clean layout
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		formContent.String(),
		"",
		lipgloss.NewStyle().Align(lipgloss.Center).Render(buttons),
		"",
	)

	m.advancedViewport.SetContent(content)

	// Clean help text
	helpText := lipgloss.JoinHorizontal(
		lipgloss.Center,
		formatKeyHint("↑/↓", "Scroll"),
		"  ",
		formatKeyHint("Tab", "Next Field"),
		"  ",
		formatKeyHint("Enter", "Select"),
	)

	// Width handling
	width := m.windowSize.Width
	if width == 0 {
		width = 100
	}

	// Clean layout with proper spacing
	return lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(lipgloss.JoinVertical(
				lipgloss.Center,
				header,
				subtitle,
			)),
		"",
		lipgloss.NewStyle().
			Width(width-10).
			Align(lipgloss.Center).
			Render(styles.ContainerStyle.Render(m.advancedViewport.View())),
		"",
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(helpText),
	)
}

// Results screen with minimal, clean tab design
func (m *RetroApp) renderResultsScreen() string {
	// Clean loading indicator
	if m.resultsLoading {
		spinner := m.spinner.View()

		loadingContent := lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			lipgloss.NewStyle().
				Foreground(styles.PrimaryColor).
				Bold(true).
				Render("Computing tax breakdown..."),
			"",
			spinner,
		)

		width := m.windowSize.Width
		if width == 0 {
			width = 100
		}

		return lipgloss.JoinVertical(
			lipgloss.Center,
			formatTitle("Tax Calculator"),
			"",
			lipgloss.NewStyle().
				Width(width).
				Align(lipgloss.Center).
				Render(loadingContent),
		)
	}

	// Clean error display
	if m.resultsError != "" {
		errorContent := lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			lipgloss.NewStyle().
				Foreground(styles.DangerColor).
				Bold(true).
				Render("Calculation Error"),
			"",
			lipgloss.NewStyle().
				Foreground(styles.DangerColor).
				Render(m.resultsError),
			"",
			styles.ButtonStyle.Render(" Press Esc to go back "),
		)

		width := m.windowSize.Width
		if width == 0 {
			width = 100
		}

		return lipgloss.JoinVertical(
			lipgloss.Center,
			formatTitle("Tax Calculator"),
			"",
			lipgloss.NewStyle().
				Width(width).
				Align(lipgloss.Center).
				Render(errorContent),
		)
	}

	// Calculate tax values
	income, _ := parseFloatWithDefault(m.incomeInput.Value(), 0)

	var incomeTax, solidarityTax float64
	for _, output := range m.result.Outputs.Output {
		if output.Name == "LSTLZZ" {
			incomeTax = float64(bmf.MustParseInt(output.Value)) / 100
		} else if output.Name == "SOLZLZZ" {
			solidarityTax = float64(bmf.MustParseInt(output.Value)) / 100
		}
	}

	totalTax := incomeTax + solidarityTax
	netIncome := income - totalTax
	taxRate := 0.0
	if income > 0 {
		taxRate = (totalTax / income) * 100
	}

	// Set tab content with clean styling
	var tabContent string
	switch m.activeTab {
	case BasicTab:
		tabContent = formatTaxResults(income, incomeTax, solidarityTax, totalTax, netIncome, taxRate)

	case DetailsTab:
		// Clean detailed view
		var details strings.Builder

		// Input parameters
		details.WriteString(formatSubTitle("Input Parameters"))
		details.WriteString("\n\n")
		details.WriteString(formatTableRow("Tax Class:", fmt.Sprintf("%d", m.selectedTaxClass), false))
		details.WriteString("\n")
		details.WriteString(formatTableRow("Income:", formatEuro(income), false))
		details.WriteString("\n")
		details.WriteString(formatTableRow("Year:", m.yearInput.Value(), false))
		details.WriteString("\n\n")

		// Advanced parameters if used
		details.WriteString(formatSubTitle("Advanced Parameters"))
		details.WriteString("\n\n")

		for _, field := range m.advancedFields {
			details.WriteString(formatTableRow(field.Label+":", field.Model.Value(), false))
			details.WriteString("\n")
		}

		details.WriteString("\n")
		details.WriteString(formatSubTitle("Raw Calculation Values"))
		details.WriteString("\n\n")

		// Show raw output values from the BMF response
		for _, output := range m.result.Outputs.Output {
			details.WriteString(formatTableRow(output.Name+":", output.Value, false))
			details.WriteString("\n")
		}

		tabContent = details.String()

	case AboutTab:
		// Clean about tab
		var about strings.Builder
		about.WriteString(formatSubTitle("About Steuer-Go"))
		about.WriteString("\n\n")
		about.WriteString(styles.BaseStyle.Render("A modern German tax calculator built with Go"))
		about.WriteString("\n\n")
		about.WriteString(formatSubTitle("Features"))
		about.WriteString("\n\n")
		about.WriteString("• Calculates income tax for all tax classes\n")
		about.WriteString("• Supports basic and advanced tax parameters\n")
		about.WriteString("• Local calculation mode for offline use\n")
		about.WriteString("• Tax rate comparison across income levels\n")
		about.WriteString("• Visualizes tax breakdown\n")
		about.WriteString("\n")
		about.WriteString(formatSubTitle("Credits"))
		about.WriteString("\n\n")
		about.WriteString("• Interface built with Bubble Tea, Bubbles, and Lip Gloss\n")
		about.WriteString("• Tax formulas from the German Federal Ministry of Finance\n")

		tabContent = about.String()
	}

	// Set viewport content
	m.resultsViewport.SetContent(tabContent)

	// Minimal button styling
	compareButton := styles.ButtonStyle.Render(" Compare Rates ")
	backButton := styles.ButtonStyle.Render(" Back ")

	actions := lipgloss.JoinHorizontal(
		lipgloss.Center,
		compareButton,
		"  ",
		backButton,
	)

	// Clean help text
	helpText := lipgloss.JoinHorizontal(
		lipgloss.Center,
		formatKeyHint("←/→", "Change Tab"),
		"  ",
		formatKeyHint("↑/↓", "Scroll"),
		"  ",
		formatKeyHint("C", "Compare"),
		"  ",
		formatKeyHint("B", "Back"),
	)

	// Width handling
	width := m.windowSize.Width
	if width == 0 {
		width = 100
	}

	// Create simple, elegant tabs
	tabs := []string{"Basic Results", "Details", "About"}
	renderedTabs := createTabs(tabs, int(m.activeTab), width)

	// Clean container styling
	resultContainer := styles.ResultsContainerStyle.Render(m.resultsViewport.View())

	// Clean layout with proper spacing
	return lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		formatTitle("Tax Calculation Results"),
		"",
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(renderedTabs),
		"",
		lipgloss.NewStyle().
			Width(width-10).
			Align(lipgloss.Center).
			Render(resultContainer),
		"",
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(actions),
		"",
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(helpText),
	)
}

// Comparison screen with minimal styling
func (m *RetroApp) renderComparisonScreen() string {
	// Clean loading display
	if m.comparisonLoading {
		spinner := m.spinner.View()

		loadingContent := lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			lipgloss.NewStyle().
				Foreground(styles.PrimaryColor).
				Bold(true).
				Render("Computing comparison data..."),
			"",
			spinner,
			"",
			lipgloss.NewStyle().
				Foreground(styles.NeutralColor).
				Render(fmt.Sprintf("Progress: %d/%d", m.completedCalls, m.totalCalls)),
		)

		width := m.windowSize.Width
		if width == 0 {
			width = 100
		}

		return lipgloss.JoinVertical(
			lipgloss.Center,
			formatTitle("Tax Rate Comparison"),
			"",
			lipgloss.NewStyle().
				Width(width).
				Align(lipgloss.Center).
				Render(loadingContent),
		)
	}

	// Clean error display
	if m.comparisonError != "" {
		errorContent := lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			lipgloss.NewStyle().
				Foreground(styles.DangerColor).
				Bold(true).
				Render("Comparison Error"),
			"",
			lipgloss.NewStyle().
				Foreground(styles.DangerColor).
				Render(m.comparisonError),
			"",
			styles.ButtonStyle.Render(" Press Esc to go back "),
		)

		width := m.windowSize.Width
		if width == 0 {
			width = 100
		}

		return lipgloss.JoinVertical(
			lipgloss.Center,
			formatTitle("Tax Rate Comparison"),
			"",
			lipgloss.NewStyle().
				Width(width).
				Align(lipgloss.Center).
				Render(errorContent),
		)
	}

	// Format comparison data
	income, _ := parseFloatWithDefault(m.incomeInput.Value(), 0)
	var comparisonContent string

	if m.showBreakdown && len(m.comparisonResults) > 0 && m.selectedComparisonIdx < len(m.comparisonResults) {
		// Show detailed breakdown for selected item
		selectedResult := m.comparisonResults[m.selectedComparisonIdx]
		comparisonContent = formatSelectedBreakdown(selectedResult)
	} else {
		// Show comparison list
		comparisonContent = formatComparisonResults(m.comparisonResults, income, m.selectedComparisonIdx)
	}

	// Set content in viewport
	m.comparisonViewport.SetContent(comparisonContent)

	// Minimal button styling
	backButton := styles.ButtonStyle.Render(" Back to Results ")

	// Clean help text
	var helpText string
	if m.showBreakdown {
		helpText = lipgloss.JoinHorizontal(
			lipgloss.Center,
			formatKeyHint("Enter", "Back to List"),
			"  ",
			formatKeyHint("B", "Back to Results"),
		)
	} else {
		helpText = lipgloss.JoinHorizontal(
			lipgloss.Center,
			formatKeyHint("↑/↓", "Select Item"),
			"  ",
			formatKeyHint("Enter", "Show Breakdown"),
			"  ",
			formatKeyHint("B", "Back to Results"),
		)
	}

	// Width handling
	width := m.windowSize.Width
	if width == 0 {
		width = 100
	}

	// Clean container styling
	comparisonContainer := styles.ResultsContainerStyle.Render(m.comparisonViewport.View())

	// Clean layout with proper spacing
	return lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		formatTitle("Tax Rate Comparison"),
		"",
		lipgloss.NewStyle().
			Width(width-10).
			Align(lipgloss.Center).
			Render(comparisonContainer),
		"",
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(backButton),
		"",
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(helpText),
	)
}
