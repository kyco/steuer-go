package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the model
func (m AppModel) Init() tea.Cmd {
	// We'll use a simpler approach - we'll modify the API function to directly
	// send messages to the program instead of trying to capture stdout
	
	return textinput.Blink
}

// Update handles events and updates the model
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case DebugLogMsg:
		// Add debug message to our list
		m.debugMessages = append(m.debugMessages, msg.Message)
		if m.step == ResultsStep {
			m.updateResultsContent()
		}
		
	case tea.WindowSizeMsg:
		m.windowSize = msg
		
		if m.step == ResultsStep {
			m.resultsViewport.Width = min(msg.Width-4, 100)
			m.resultsViewport.Height = msg.Height - 8
		} else if m.step == ComparisonStep {
			m.comparisonViewport.Width = min(msg.Width-4, 100)
			m.comparisonViewport.Height = msg.Height - 8
		}
		
	case tea.KeyMsg:
		switch m.step {
		case InputStep:
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
				
			case "tab", "shift+tab":
				// Cycle through input fields
				if msg.String() == "tab" {
					m.focusField = (m.focusField + 1) % 4
				} else {
					m.focusField = (m.focusField - 1 + 4) % 4
				}
				m.updateFocus()
				
			case "up", "down":
				if m.focusField == TaxClassField {
					// Navigate through tax class options
					if msg.String() == "up" {
						m.selectedTaxClass = max(1, m.selectedTaxClass-1)
					} else {
						m.selectedTaxClass = min(6, m.selectedTaxClass+1)
					}
				}
				
			case "enter":
				if m.focusField == CalculateButtonField {
					valid, errMsg := m.validateAndCalculate()
					if valid {
						// Start calculation
						m.step = ResultsStep
						m.resultsLoading = true
						
						income := parseIncome(m.incomeInput.Value())
						cmds = append(cmds, m.spinner.Tick)
						cmds = append(cmds, PerformCalculationCmd(m.selectedTaxClass, income, m.yearInput.Value()))
					} else {
						m.resultsError = errMsg
					}
				} else {
					// Move to next field
					m.focusField = (m.focusField + 1) % 4
					m.updateFocus()
				}
			}
		
		case ResultsStep:
			switch msg.String() {
			case "esc", "q", "b":
				// Return to input form
				m.step = InputStep
				m.resultsLoading = false
				m.resultsError = ""
				return m, nil
				
			case "d":
				// Toggle details
				m.showDetails = !m.showDetails
				m.updateResultsContent()
			
			case "c":
				// Switch to comparison view
				m.step = ComparisonStep
				m.comparisonLoading = true
				m.comparisonResults = nil
				m.completedCalls = 0
				m.totalCalls = 0
				
				cmds = append(cmds, m.spinner.Tick)
				cmds = append(cmds, PerformComparisonCmd())
				
			case "ctrl+c":
				return m, tea.Quit
			}
			
		case ComparisonStep:
			switch msg.String() {
			case "esc", "q", "b":
				// Return to regular results view
				m.step = ResultsStep
				return m, nil
				
			case "ctrl+c":
				return m, tea.Quit
			}
		}
		
	case CalculationStartedMsg:
		cmds = append(cmds, FetchResultsCmd(m.selectedTaxClass, parseIncome(m.incomeInput.Value())))
		
	case CalculationMsg:
		m.resultsLoading = false
		if msg.Error != nil {
			m.resultsError = msg.Error.Error()
		} else {
			m.result = msg.Result
			m.updateResultsContent()
		}
		
	case ComparisonStartedMsg:
		// Show loading state while we calculate
		m.comparisonLoading = true
		// Reset variables
		m.totalCalls = 0
		m.completedCalls = 0
		// Clear any old results
		m.comparisonResults = nil
		// Clear any old debug messages
		m.debugMessages = []string{}
		// Fetch tax comparison
		cmds = append(cmds, FetchComparisonCmd(m.selectedTaxClass, parseIncome(m.incomeInput.Value())))
		
	case ComparisonProgressMsg:
		m.completedCalls = msg.CompletedCalls
		m.totalCalls = msg.TotalCalls
		
	case ComparisonMsg:
		m.comparisonLoading = false
		if msg.Error != nil {
			m.comparisonError = msg.Error.Error()
		} else {
			m.comparisonResults = msg.Results
			// Only update content if we actually have results
			if len(msg.Results) > 0 {
				m.updateComparisonContent()
			}
		}
		
	case spinner.TickMsg:
		if m.resultsLoading || m.comparisonLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Handle input updates
	if m.step == InputStep {
		switch m.focusField {
		case IncomeField:
			newIncomeInput, cmd := m.incomeInput.Update(msg)
			m.incomeInput = newIncomeInput
			cmds = append(cmds, cmd)
		case YearField:
			newYearInput, cmd := m.yearInput.Update(msg)
			m.yearInput = newYearInput
			cmds = append(cmds, cmd)
		}
	} else if m.step == ResultsStep {
		var cmd tea.Cmd
		m.resultsViewport, cmd = m.resultsViewport.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.step == ComparisonStep {
		var cmd tea.Cmd
		m.comparisonViewport, cmd = m.comparisonViewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the model
func (m AppModel) View() string {
	switch m.step {
	case InputStep:
		return m.renderInputForm()
	case ResultsStep:
		return m.renderResults()
	case ComparisonStep:
		return m.renderComparison()
	default:
		return "Loading..."
	}
}