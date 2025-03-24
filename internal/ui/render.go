package ui

import (
	"fmt"
	"strings"
	
	"github.com/charmbracelet/lipgloss"
	
	"tax-calculator/internal/ui/styles"
)

func (m AppModel) renderInputForm() string {
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
	
	errorMsg := ""
	if m.resultsError != "" {
		errorMsg = lipgloss.NewStyle().
			Foreground(styles.DangerColor).
			Render(m.resultsError)
	}
	
	helpText := styles.HelpStyle.Render("Tab: Next Field • Enter: Select/Calculate • ↑/↓: Navigate Tax Class • Esc: Quit")
	
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
		calculateButton,
		"",
		errorMsg,
	)
	
	formWidth := m.windowSize.Width
	if formWidth == 0 {
		formWidth = 100
	}
	
	centeredForm := lipgloss.NewStyle().
		Width(formWidth).
		Align(lipgloss.Center).
		Render(content)
	
	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		centeredForm,
		"",
		lipgloss.NewStyle().
			Width(formWidth).
			Align(lipgloss.Center).
			Render(helpText),
	)
}

func (m AppModel) renderResults() string {
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
	
	helpText := styles.HelpStyle.Render("d: Toggle Details • c: Compare Tax Rates • ↑/↓: Scroll • b: Back • Esc: Quit")
	
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

func (m AppModel) renderComparison() string {
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