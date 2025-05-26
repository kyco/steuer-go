package views

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"tax-calculator/internal/tax/calculation"
	"tax-calculator/internal/tax/models"
)

func (m *RetroApp) Init() tea.Cmd {
	// Start the spinner animation
	return tea.Batch(
		spinner.Tick,
		tea.EnterAltScreen,
	)
}

func (m *RetroApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Check if any input is currently focused
	var inputFocused bool
	
	// Main screen input fields
	if m.screen == MainScreen {
		if m.incomeInput.Focused() || m.yearInput.Focused() {
			inputFocused = true
		}
	}
	
	// Advanced screen input fields
	if m.screen == AdvancedScreen {
		for _, field := range m.advancedFields {
			if field.Model.Focused() {
				inputFocused = true
				break
			}
		}
	}
	
	// Handle special key events when an input is focused
	if inputFocused {
		// Special case for Escape key to blur input
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
			// Blur any focused input fields
			if m.screen == MainScreen {
				if m.incomeInput.Focused() {
					m.incomeInput.Blur()
				}
				if m.yearInput.Focused() {
					m.yearInput.Blur()
				}
			} else if m.screen == AdvancedScreen {
				for i, field := range m.advancedFields {
					if field.Model.Focused() {
						newModel := field.Model
						newModel.Blur()
						m.advancedFields[i].Model = newModel
					}
				}
			}
			// Continue processing the update
		} else if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
			// Blur any focused input fields on Enter
			if m.screen == MainScreen {
				if m.incomeInput.Focused() {
					m.incomeInput.Blur()
				}
				if m.yearInput.Focused() {
					m.yearInput.Blur()
				}
			} else if m.screen == AdvancedScreen {
				for i, field := range m.advancedFields {
					if field.Model.Focused() {
						newModel := field.Model
						newModel.Blur()
						m.advancedFields[i].Model = newModel
					}
				}
			}
			// Continue processing the update
		} else {
			// Regular input handling when focused
			// Will be handled by the input update section below
		}
	} else {
		// Regular key handling when no input is focused
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
	
			case "esc":
				switch m.screen {
				case ResultsScreen, ComparisonScreen, AdvancedScreen:
					m.screen = MainScreen
				default:
					return m, tea.Quit
				}
	
			case "tab", "shift+tab":
				// Handle tab navigation across fields
				m.handleTabNavigation(keyMsg.String() == "shift+tab")
	
			case "up", "down":
				// Handle up/down navigation
				m.handleUpDownNavigation(keyMsg.String() == "up")
	
			case "left", "right":
				// Handle left/right navigation
				m.handleLeftRightNavigation(keyMsg.String() == "left")
	
			case "enter":
				// Handle enter key selection
				cmd := m.handleEnterSelection()
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
	
			case "l":
				// Toggle local calculation
				m.useLocalCalc = !m.useLocalCalc
	
			case "c":
				// Start comparison mode from results screen
				if m.screen == ResultsScreen {
					m.screen = ComparisonScreen
					m.comparisonLoading = true
					cmds = append(cmds, m.startComparisonCmd())
				}
	
			case "b":
				// Go back from comparison to results
				if m.screen == ComparisonScreen {
					m.screen = ResultsScreen
				} else if m.screen == ResultsScreen {
					m.screen = MainScreen
				}
	
			case "d":
				// Toggle details in results
				if m.screen == ResultsScreen {
					m.showDetails = !m.showDetails
				}
			}
		}
	}

	// Handle other message types
	switch msgType := msg.(type) {
	case tea.WindowSizeMsg:
		// Store window size for responsive layouts
		m.windowSize = msgType
		m.updateViewportDimensions(msgType)

	case spinner.TickMsg:
		// Update spinner animation
		newSpinner, cmd := m.spinner.Update(msg)
		m.spinner = newSpinner
		cmds = append(cmds, cmd)

	case CalculationStartedMsg:
		// When calculation starts
		m.resultsLoading = true
		m.screen = ResultsScreen
		m.useLocalCalc = msgType.UseLocalCalculator

	case CalculationMsg:
		// When calculation completes
		m.resultsLoading = false
		
		if msgType.Error != nil {
			m.resultsError = msgType.Error.Error()
		} else {
			m.result = msgType.Result
			m.resultsError = ""
		}

	case ComparisonStartedMsg:
		// When comparison starts
		m.comparisonLoading = true
		m.completedCalls = 0
		m.totalCalls = 0

	case ComparisonProgressMsg:
		// Update comparison progress
		m.completedCalls = msgType.CompletedCalls
		m.totalCalls = msgType.TotalCalls

	case ComparisonMsg:
		// When comparison completes
		m.comparisonLoading = false
		
		if msgType.Error != nil {
			m.comparisonError = msgType.Error.Error()
		} else {
			m.comparisonResults = msgType.Results
			m.comparisonError = ""
		}

	case DebugLogMsg:
		// Skip debug messages in this UI
	}

	// Handle viewport updates
	switch m.screen {
	case ResultsScreen:
		newViewport, cmd := m.resultsViewport.Update(msg)
		m.resultsViewport = newViewport
		cmds = append(cmds, cmd)
	
	case ComparisonScreen:
		newViewport, cmd := m.comparisonViewport.Update(msg)
		m.comparisonViewport = newViewport
		cmds = append(cmds, cmd)
	
	case AdvancedScreen:
		newViewport, cmd := m.advancedViewport.Update(msg)
		m.advancedViewport = newViewport
		cmds = append(cmds, cmd)
	}

	// Always update input fields regardless of focus state
	if m.screen == MainScreen {
		// Always update income field
		newIncomeInput, incomeCmd := m.incomeInput.Update(msg)
		m.incomeInput = newIncomeInput
		if incomeCmd != nil {
			cmds = append(cmds, incomeCmd)
		}
		
		// Always update year field
		newYearInput, yearCmd := m.yearInput.Update(msg)
		m.yearInput = newYearInput
		if yearCmd != nil {
			cmds = append(cmds, yearCmd)
		}
	}
	
	// Always update all advanced fields
	if m.screen == AdvancedScreen {
		for i, field := range m.advancedFields {
			newModel, cmd := field.Model.Update(msg)
			m.advancedFields[i].Model = newModel
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// Handle tab navigation across input fields
func (m *RetroApp) handleTabNavigation(isBackward bool) {
	switch m.screen {
	case MainScreen:
		fields := []Field{TaxClassField, IncomeField, YearField, CalculateButtonField, AdvancedButtonField}
		m.navigateFields(fields, isBackward)

	case AdvancedScreen:
		var fields []Field
		for _, field := range m.advancedFields {
			fields = append(fields, field.Field)
		}
		fields = append(fields, BackButtonField, CalculateButtonField)
		m.navigateFields(fields, isBackward)
	}
}

// Helper for field navigation
func (m *RetroApp) navigateFields(fields []Field, isBackward bool) {
	currentIdx := -1
	for i, field := range fields {
		if field == m.focusField {
			currentIdx = i
			break
		}
	}

	if currentIdx == -1 {
		m.focusField = fields[0]
		return
	}

	if isBackward {
		currentIdx--
		if currentIdx < 0 {
			currentIdx = len(fields) - 1
		}
	} else {
		currentIdx++
		if currentIdx >= len(fields) {
			currentIdx = 0
		}
	}

	m.focusField = fields[currentIdx]
}

// Handle up/down navigation
func (m *RetroApp) handleUpDownNavigation(isUp bool) {
	switch m.screen {
	case MainScreen:
		if m.focusField == TaxClassField {
			// Tax class selection
			if isUp {
				m.selectedTaxClass--
				if m.selectedTaxClass < 1 {
					m.selectedTaxClass = 6
				}
			} else {
				m.selectedTaxClass++
				if m.selectedTaxClass > 6 {
					m.selectedTaxClass = 1
				}
			}
		}

	case ResultsScreen:
		if isUp {
			m.resultsViewport.LineUp(1)
		} else {
			m.resultsViewport.LineDown(1)
		}

	case ComparisonScreen:
		if isUp {
			m.comparisonViewport.LineUp(1)
		} else {
			m.comparisonViewport.LineDown(1)
		}

	case AdvancedScreen:
		if isUp {
			m.advancedViewport.LineUp(1)
		} else {
			m.advancedViewport.LineDown(1)
		}
	}
}

// Handle left/right navigation
func (m *RetroApp) handleLeftRightNavigation(isLeft bool) {
	if m.screen == ResultsScreen {
		if isLeft {
			m.activeTab--
			if m.activeTab < 0 {
				m.activeTab = AboutTab
			}
		} else {
			m.activeTab++
			if m.activeTab > AboutTab {
				m.activeTab = BasicTab
			}
		}
	}
}

// Handle enter key selection
func (m *RetroApp) handleEnterSelection() tea.Cmd {
	switch m.screen {
	case MainScreen:
		switch m.focusField {
		case IncomeField:
			// Focus the income input field
			m.incomeInput.Focus()
			return textinput.Blink
		case YearField:
			// Focus the year input field
			m.yearInput.Focus()
			return textinput.Blink
		case CalculateButtonField:
			return m.startCalculationCmd()
		case AdvancedButtonField:
			m.screen = AdvancedScreen
			m.focusField = AJAHR_Field
		}

	case AdvancedScreen:
		// Handle advanced field selection
		for i, field := range m.advancedFields {
			if field.Field == m.focusField {
				// Focus this input field
				newModel := field.Model
				newModel.Focus()
				m.advancedFields[i].Model = newModel
				return textinput.Blink
			}
		}
		
		switch m.focusField {
		case BackButtonField:
			m.screen = MainScreen
			m.focusField = TaxClassField
		case CalculateButtonField:
			return m.startAdvancedCalculationCmd()
		}
	}

	return nil
}

// Update viewport dimensions when window size changes
func (m *RetroApp) updateViewportDimensions(msg tea.WindowSizeMsg) {
	width := msg.Width - 20 // margin for borders
	height := msg.Height - 20 // margin for header/footer
	
	if width < 20 {
		width = 20
	}
	
	if height < 10 {
		height = 10
	}
	
	m.resultsViewport.Width = width
	m.resultsViewport.Height = height
	
	m.comparisonViewport.Width = width
	m.comparisonViewport.Height = height
	
	m.advancedViewport.Width = width
	m.advancedViewport.Height = height
}

// Start tax calculation command
func (m *RetroApp) startCalculationCmd() tea.Cmd {
	income, err := parseFloatWithDefault(m.incomeInput.Value(), 0)
	if err != nil || income <= 0 {
		return func() tea.Msg {
			return CalculationMsg{
				Error: err,
			}
		}
	}

	year := m.yearInput.Value()
	if strings.TrimSpace(year) == "" {
		year = time.Now().Format("2006")
	}

	// Use the proper calculation mode based on local flag
	return PerformCalculationWithAdvancedOptionsCmd(
		m.selectedTaxClass,
		income,
		year,
		models.TaxRequest{}, // Empty tax request with default values
		m.useLocalCalc,      // Pass the local calculator flag
	)
}

// Start advanced tax calculation command
func (m *RetroApp) startAdvancedCalculationCmd() tea.Cmd {
	income, err := parseFloatWithDefault(m.incomeInput.Value(), 0)
	if err != nil || income <= 0 {
		return func() tea.Msg {
			return CalculationMsg{
				Error: err,
			}
		}
	}

	year := m.yearInput.Value()
	if strings.TrimSpace(year) == "" {
		year = time.Now().Format("2006")
	}

	// Build the tax request from advanced fields
	taxRequest := m.buildTaxRequest()

	return PerformCalculationWithAdvancedOptionsCmd(
		m.selectedTaxClass,
		income,
		year,
		taxRequest,
		m.useLocalCalc,
	)
}

// Start comparison command
func (m *RetroApp) startComparisonCmd() tea.Cmd {
	income, err := parseFloatWithDefault(m.incomeInput.Value(), 0)
	if err != nil || income <= 0 {
		return func() tea.Msg {
			return ComparisonMsg{
				Error: err,
			}
		}
	}

	// Create a command for comparison calculations
	return tea.Sequence(
		func() tea.Msg {
			return ComparisonStartedMsg{}
		},
		func() tea.Msg {
			// Create the tax service
			taxService := calculation.NewTaxService()
			if m.useLocalCalc {
				taxService.EnableLocalCalculator()
			}

			// Run the calculation with progress updates
			return FetchComparisonCmd(m.selectedTaxClass, income)()
		},
	)
}