package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"tax-calculator/internal/ui/styles"
)

func (m *AppModel) renderInputForm() string {
	header := styles.HeaderStyle.Render("German Tax Calculator")

	taxClassTitle := styles.SubtitleStyle.Render("Tax Class:")

	var taxClassOptions strings.Builder
	for i, option := range m.taxClassOptions {
		classNum := option.Class
		style := styles.UnselectedItemStyle
		indicator := "  "

		if classNum == m.selectedTaxClass {
			style = styles.SelectedItemStyle
			indicator = "▶ "
		}

		if m.focusField == TaxClassField {
			indicator = "  "
			if classNum == m.selectedTaxClass {
				indicator = "▶ "
			}
		}

		fmt.Fprintf(&taxClassOptions, "%s%s\n",
			indicator,
			style.Render(fmt.Sprintf("Class %d: %s", classNum, option.Desc)))

		if i >= 2 && classNum != m.selectedTaxClass && classNum != m.selectedTaxClass+1 {
			continue
		}
	}

	incomeTitle := styles.SubtitleStyle.Render("Annual Income:")
	incomeField := styles.InputFieldStyle.Render(m.incomeInput.View())
	if m.focusField == IncomeField {
		incomeField = styles.ActiveInputStyle.Render(m.incomeInput.View())
	}

	yearTitle := styles.SubtitleStyle.Render("Tax Year:")
	yearField := styles.InputFieldStyle.Render(m.yearInput.View())
	if m.focusField == YearField {
		yearField = styles.ActiveInputStyle.Render(m.yearInput.View())
	}

	calculateButton := styles.ButtonStyle.Render(" Calculate ")
	if m.focusField == CalculateButtonField {
		calculateButton = styles.SelectedButtonStyle.Render(" Calculate ")
	}

	advancedButton := styles.ButtonStyle.Render(" Advanced Options ")
	if m.focusField == AdvancedButtonField {
		advancedButton = styles.SelectedButtonStyle.Render(" Advanced Options ")
	}

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		calculateButton,
		"     ",
		advancedButton,
	)

	errorMsg := ""
	if m.resultsError != "" {
		errorMsg = lipgloss.NewStyle().
			Foreground(styles.DangerColor).
			Render(m.resultsError)
	}

	localCalcText := ""
	if m.useLocalCalc {
		localCalcText = lipgloss.NewStyle().Foreground(styles.SuccessColor).Render(" • Local Calculation: ON")
	}

	helpText := styles.HelpStyle.Render(fmt.Sprintf("Tab: Next Field • Enter: Select/Calculate • ↑/↓: Navigate Tax Class • l: Toggle Local Calculation%s • Esc: Quit", localCalcText))

	content := lipgloss.JoinVertical(
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
		buttons,
		"",
		errorMsg,
	)

	formStyle := lipgloss.NewStyle().
		Border(styles.MinimalBorder).
		BorderForeground(styles.PrimaryColor).
		Padding(1, 3)

	formWidth := m.windowSize.Width
	if formWidth == 0 {
		formWidth = 100
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		lipgloss.NewStyle().
			Width(formWidth).
			Align(lipgloss.Center).
			Render(formStyle.Render(content)),
		"",
		lipgloss.NewStyle().
			Width(formWidth).
			Align(lipgloss.Center).
			Render(helpText),
	)
}

func (m *AppModel) renderResults() string {
	header := styles.HeaderStyle.Render("German Tax Calculator - Results")

	if m.resultsLoading {
		spinner := m.spinner.View()

		width := m.windowSize.Width
		if width == 0 {
			width = 100
		}

		loadingContent := lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			lipgloss.NewStyle().
				Foreground(styles.PrimaryColor).
				Render("Calculating tax..."),
			spinner,
		)

		return lipgloss.JoinVertical(
			lipgloss.Center,
			header,
			lipgloss.NewStyle().
				Width(width).
				Align(lipgloss.Center).
				Render(loadingContent),
		)
	}

	if m.resultsError != "" {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			header,
			"",
			lipgloss.NewStyle().
				Foreground(styles.DangerColor).
				Render("Error:"),
			lipgloss.NewStyle().
				Foreground(styles.DangerColor).
				Render(m.resultsError),
			"",
			styles.HelpStyle.Render("Press 'esc' to go back"),
		)
	}

	localCalcText := ""
	if m.useLocalCalc {
		localCalcText = lipgloss.NewStyle().Foreground(styles.SuccessColor).Render(" • Local Calculation: ON")
	}

	helpText := styles.HelpStyle.Render(fmt.Sprintf("d: Toggle Details • c: Compare Tax Rates • l: Toggle Local Calculation%s • ↑/↓: Scroll • b: Back • Esc: Quit", localCalcText))

	width := m.windowSize.Width
	if width == 0 {
		width = 100
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(m.resultsViewport.View()),
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(helpText),
	)
}

func (m *AppModel) renderComparison() string {
	header := styles.HeaderStyle.Render("German Tax Calculator - Income Comparison")

	if m.comparisonLoading {
		width := m.windowSize.Width
		if width == 0 {
			width = 100
		}

		loadingContent := lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			lipgloss.NewStyle().
				Foreground(styles.PrimaryColor).
				Bold(true).
				Render("Loading income comparison..."),
		)

		return lipgloss.JoinVertical(
			lipgloss.Center,
			header,
			lipgloss.NewStyle().
				Width(width).
				Align(lipgloss.Center).
				Render(loadingContent),
		)
	}

	if m.comparisonError != "" {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			header,
			"",
			lipgloss.NewStyle().
				Foreground(styles.DangerColor).
				Render("Error:"),
			lipgloss.NewStyle().
				Foreground(styles.DangerColor).
				Render(m.comparisonError),
			"",
			styles.HelpStyle.Render("Press 'esc' to go back"),
		)
	}

	helpText := styles.HelpStyle.Render("↑/↓: Scroll • b: Back to Results • Esc: Quit")

	width := m.windowSize.Width
	if width == 0 {
		width = 100
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(m.comparisonViewport.View()),
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(helpText),
	)
}

func (m *AppModel) renderAdvancedInputForm() string {
	header := styles.HeaderStyle.Render("German Tax Calculator - Advanced Parameters")

	// Helper function to render input fields with descriptions
	renderInputWithDesc := func(title string, description string, input textinput.Model, isFocused bool) string {
		titleText := styles.SubtitleStyle.Render(title)
		descText := lipgloss.NewStyle().
			Foreground(styles.NeutralColor).
			Italic(true).
			Width(60).
			Render(description)

		style := styles.InputFieldStyle
		if isFocused {
			style = styles.ActiveInputStyle
		}
		inputField := style.Render(input.View())

		return lipgloss.JoinVertical(
			lipgloss.Left,
			titleText,
			descText,
			inputField,
		)
	}

	// Create a single column for all inputs
	singleColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		renderInputWithDesc("Year following 64th year:",
			"Calendar year after taxpayer's 64th birthday. Enter 0 if not applicable.",
			m.ajahr, m.focusField == AJAHR_Field),
		"",
		renderInputWithDesc("Completed 64 years (0/1):",
			"1 if taxpayer was 64+ at the start of calendar year, otherwise 0.",
			m.alter1, m.focusField == ALTER1_Field),
		"",
		renderInputWithDesc("Social insurance (0-2):",
			"0: Normal statutory pension, 1: No compulsory insurance, 2: Reduced rate for miners/seamen",
			m.krv, m.focusField == KRV_Field),
		"",
		renderInputWithDesc("Additional health rate (%):",
			"Additional health insurance rate (usually 0.3-2.2%). Standard is 1.3%.",
			m.kvz, m.focusField == KVZ_Field),
		"",
		renderInputWithDesc("Employer in Saxony (0/1):",
			"1 if employer is based in Saxony (different nursing care contributions), otherwise 0.",
			m.pvs, m.focusField == PVS_Field),
		"",
		renderInputWithDesc("Childless surcharge (0/1):",
			"1 if employee (aged 23+) pays childless surcharge for nursing care insurance, otherwise 0.",
			m.pvz, m.focusField == PVZ_Field),
		"",
		renderInputWithDesc("Religion code (0-2):",
			"0: No church tax, 1: Catholic church, 2: Protestant church",
			m.r, m.focusField == R_Field),
		"",
		renderInputWithDesc("Child allowance:",
			"Number of children for tax allowance. Can be decimal (0.5 for shared custody).",
			m.zkf, m.focusField == ZKF_Field),
		"",
		renderInputWithDesc("Pension payments (euros):",
			"Annual pension income in euros. Enter 0 if no pension income.",
			m.vbez, m.focusField == VBEZ_Field),
		"",
		renderInputWithDesc("First pension year:",
			"Year when taxpayer first started receiving pension. Enter 0 if not applicable.",
			m.vjahr, m.focusField == VJAHR_Field),
		"",
		renderInputWithDesc("Private insurance payment (euros):",
			"Monthly private health insurance premium in euros. Only for private insurance.",
			m.pkpv, m.focusField == PKPV_Field),
		"",
		renderInputWithDesc("Health insurance type (0-2):",
			"0: Statutory health insurance, 1: Private without employer subsidy, 2: Private with subsidy",
			m.pkv, m.focusField == PKV_Field),
		"",
		renderInputWithDesc("Children for care insurance (0-4):",
			"Number of children for reduced nursing care insurance contributions.",
			m.pva, m.focusField == PVA_Field),
	)

	backButton := styles.ButtonStyle.Render(" Back ")
	if m.focusField == BackButtonField {
		backButton = styles.SelectedButtonStyle.Render(" Back ")
	}

	calculateButton := styles.ButtonStyle.Render(" Calculate with Advanced Options ")
	if m.focusField == CalculateButtonField {
		calculateButton = styles.SelectedButtonStyle.Render(" Calculate with Advanced Options ")
	}

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		backButton,
		"     ",
		calculateButton,
	)

	localCalcText := ""
	if m.useLocalCalc {
		localCalcText = lipgloss.NewStyle().Foreground(styles.SuccessColor).Render(" • Local Calculation: ON")
	}

	// Combine all content
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		singleColumn,
		"",
		lipgloss.NewStyle().Align(lipgloss.Center).Render(buttons),
	)

	// Set content to viewport for scrolling - using pointer receiver
	// Add extra padding at the bottom to ensure there's room to scroll
	paddedContent := content + "" // Extra padding lines
	m.advancedViewport.SetContent(paddedContent)

	formWidth := m.windowSize.Width
	if formWidth == 0 {
		formWidth = 100
	}

	helpText := styles.HelpStyle.Render(fmt.Sprintf("↑/↓: Scroll • Tab: Next Field • Enter: Select • l: Toggle Local Calculation%s • Esc: Quit", localCalcText))

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		lipgloss.NewStyle().
			Width(formWidth).
			Align(lipgloss.Center).
			Render(m.advancedViewport.View()),
		"",
		lipgloss.NewStyle().
			Width(formWidth).
			Align(lipgloss.Center).
			Render(helpText),
	)
}

