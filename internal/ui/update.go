package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *AppModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case DebugLogMsg:
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
		} else if m.step == AdvancedInputStep {
			m.advancedViewport.Width = min(msg.Width-4, 100)
			m.advancedViewport.Height = msg.Height - 8
		}
		
	case tea.KeyMsg:
		switch m.step {
		case InputStep:
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
				
			case "tab", "shift+tab":
				if msg.String() == "tab" {
					m.focusField = (m.focusField + 1) % 5 // Updated to include AdvancedButtonField
				} else {
					m.focusField = (m.focusField - 1 + 5) % 5
				}
				m.updateFocus()
				
			case "up", "down":
				if m.focusField == TaxClassField {
					if msg.String() == "up" {
						m.selectedTaxClass = max(1, m.selectedTaxClass-1)
					} else {
						m.selectedTaxClass = min(6, m.selectedTaxClass+1)
					}
				}
				
			case "l":
				// Toggle local calculation mode
				m.useLocalCalc = !m.useLocalCalc
				m.debugMessages = append(m.debugMessages, fmt.Sprintf("Local calculation mode: %v", m.useLocalCalc))
				
			case "enter":
				if m.focusField == CalculateButtonField {
					valid, errMsg := m.validateAndCalculate()
					if valid {
						m.step = ResultsStep
						m.resultsLoading = true
						
						income := parseIncome(m.incomeInput.Value())
						cmds = append(cmds, m.spinner.Tick)
						cmds = append(cmds, PerformCalculationCmd(m.selectedTaxClass, income, m.yearInput.Value()))
					} else {
						m.resultsError = errMsg
					}
				} else if m.focusField == AdvancedButtonField {
					// Switch to advanced options screen
					m.step = AdvancedInputStep
					m.focusField = AJAHR_Field
					m.updateFocus()
					// Reset viewport to top
					m.advancedViewport.SetYOffset(0)
				} else {
					m.focusField = (m.focusField + 1) % 5
					m.updateFocus()
				}
			}
			
		case AdvancedInputStep:
			switch msg.String() {
			case "ctrl+c", "esc":
				m.step = InputStep
				return m, nil
				
			case "tab", "shift+tab":
				// Store current focused field
				// Previous field no longer needed

				if msg.String() == "tab" {
					// Forward tabbing
					newField := int(m.focusField) + 1
					if newField > int(BackButtonField) {
						newField = int(AJAHR_Field)
					}
					m.focusField = Field(newField)
				} else {
					// Backward tabbing
					newField := int(m.focusField) - 1
					if newField < int(AJAHR_Field) {
						newField = int(BackButtonField)
					}
					m.focusField = Field(newField)
				}
				m.updateFocus()
				
				// Use direct field scrolling for more precise positioning
				m.scrollToField(m.focusField)
				
			case "enter":
				if m.focusField == BackButtonField {
					m.step = InputStep
					m.focusField = TaxClassField
					m.updateFocus()
				} else if m.focusField == CalculateButtonField {
					valid, errMsg := m.validateAndCalculate()
					if valid {
						m.step = ResultsStep
						m.resultsLoading = true
						
						income := parseIncome(m.incomeInput.Value())
						cmds = append(cmds, m.spinner.Tick)
						
						// Extract advanced parameters using our helper function
						advancedParams := GetAdvancedParametersFromModel(m)
						
						cmds = append(cmds, PerformCalculationWithAdvancedOptionsCmd(
							m.selectedTaxClass, income, m.yearInput.Value(), 
							advancedParams,
							m.useLocalCalc,
						))
					} else {
						m.resultsError = errMsg
					}
				} else {
					// Store current focused field
					// Previous field no longer needed
					
					// Move to the next field when pressing enter
					newField := int(m.focusField) + 1
					if newField > int(BackButtonField) {
						newField = int(AJAHR_Field)
					}
					m.focusField = Field(newField)
					m.updateFocus()
					
					// First update focus, then explicitly scroll to ensure visibility
					m.scrollToField(m.focusField)
				}
			
			// Handle up/down, pageup/pagedown in advanced view
			case "up", "down", "pageup", "pagedown", "home", "end":
				// Let the viewport handle these directly
				var cmd tea.Cmd
				m.advancedViewport, cmd = m.advancedViewport.Update(msg)
				cmds = append(cmds, cmd)
			}
		
		case ResultsStep:
			switch msg.String() {
			case "esc", "q", "b":
				m.step = InputStep
				m.resultsLoading = false
				m.resultsError = ""
				return m, nil
				
			case "d":
				m.showDetails = !m.showDetails
				m.updateResultsContent()
			
			case "l":
				m.useLocalCalc = !m.useLocalCalc
				return m, CaptureDebugCmd(fmt.Sprintf("Local calculation mode: %v", m.useLocalCalc))
			
			case "c":
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
				m.step = ResultsStep
				return m, nil
				
			case "ctrl+c":
				return m, tea.Quit
			}
		}
		
	case CalculationStartedMsg:
		cmds = append(cmds, FetchResultsWithModeCmd(m.selectedTaxClass, parseIncome(m.incomeInput.Value()), m.useLocalCalc))
		
	case CalculationMsg:
		m.resultsLoading = false
		if msg.Error != nil {
			m.resultsError = msg.Error.Error()
		} else {
			m.result = msg.Result
			m.updateResultsContent()
		}
		
	case ComparisonStartedMsg:
		m.comparisonLoading = true
		m.totalCalls = 0
		m.completedCalls = 0
		m.comparisonResults = nil
		m.debugMessages = []string{}
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
	} else if m.step == AdvancedInputStep {
		// Handle input for advanced fields based on focus
		var cmd tea.Cmd
		switch m.focusField {
		case AJAHR_Field:
			m.ajahr, cmd = m.ajahr.Update(msg)
			cmds = append(cmds, cmd)
		case ALTER1_Field:
			m.alter1, cmd = m.alter1.Update(msg)
			cmds = append(cmds, cmd)
		case KRV_Field:
			m.krv, cmd = m.krv.Update(msg)
			cmds = append(cmds, cmd)
		case KVZ_Field:
			m.kvz, cmd = m.kvz.Update(msg)
			cmds = append(cmds, cmd)
		case PVS_Field:
			m.pvs, cmd = m.pvs.Update(msg)
			cmds = append(cmds, cmd)
		case PVZ_Field:
			m.pvz, cmd = m.pvz.Update(msg)
			cmds = append(cmds, cmd)
		case R_Field:
			m.r, cmd = m.r.Update(msg)
			cmds = append(cmds, cmd)
		case ZKF_Field:
			m.zkf, cmd = m.zkf.Update(msg)
			cmds = append(cmds, cmd)
		case VBEZ_Field:
			m.vbez, cmd = m.vbez.Update(msg)
			cmds = append(cmds, cmd)
		case VJAHR_Field:
			m.vjahr, cmd = m.vjahr.Update(msg)
			cmds = append(cmds, cmd)
		case PKPV_Field:
			m.pkpv, cmd = m.pkpv.Update(msg)
			cmds = append(cmds, cmd)
		case PKV_Field:
			m.pkv, cmd = m.pkv.Update(msg)
			cmds = append(cmds, cmd)
		case PVA_Field:
			m.pva, cmd = m.pva.Update(msg)
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

// scrollToAdvancedField adjusts the viewport to show the focused field
func (m *AppModel) scrollToAdvancedField(prevField, newField Field) {
	// Map fields to their precise vertical positions in the rendered content
	// These values are carefully calibrated to match the actual rendered positions
	fieldPositions := []int{
		15,     // AJAHR_Field
		30,     // ALTER1_Field
		45,     // KRV_Field
		60,     // KVZ_Field
		75,     // PVS_Field
		90,     // PVZ_Field
		105,    // R_Field
		120,    // ZKF_Field
		135,    // VBEZ_Field
		150,    // VJAHR_Field
		165,    // PKPV_Field
		180,    // PKV_Field
		195,    // PVA_Field
		215,    // BackButtonField
		220,    // CalculateButtonField - slightly below BackButtonField
	}
	
	// Get field index for position lookup
	newIndex := int(newField) - int(AJAHR_Field)
	
	// Check bounds
	if newIndex < 0 || newIndex >= len(fieldPositions) {
		return
	}
	
	// Get the position of the new field
	newPos := fieldPositions[newIndex]
	
	// Calculate viewport parameters
	viewportHeight := m.advancedViewport.Height
	currentOffset := m.advancedViewport.YOffset
	
	// Calculate the visible range
	visibleStart := currentOffset
	visibleEnd := currentOffset + viewportHeight - 10 // Margin for visibility
	
	// If the field is not visible, adjust scrolling
	if newPos < visibleStart+5 || newPos > visibleEnd {
		// Calculate ideal scroll position to center field in viewport
		idealScrollPos := max(0, newPos - (viewportHeight / 2))
		
		// For button fields (near bottom), make sure they appear in the lower part of the viewport
		if newField == BackButtonField || newField == CalculateButtonField {
			// Center the buttons in the bottom half of the screen instead of scrolling all the way
			idealScrollPos = max(0, newPos - (viewportHeight * 3/4))
		}
		
		// Apply the scroll position
		m.advancedViewport.SetYOffset(idealScrollPos)
	}
}

// scrollToField ensures a specific field is visible in the viewport
func (m *AppModel) scrollToField(field Field) {
	// This is used for direct field scrolling
	if m.step != AdvancedInputStep {
		return
	}
	
	// Map of exact field positions
	fieldPositions := map[Field]int{
		AJAHR_Field:           15,
		ALTER1_Field:          30,
		KRV_Field:             45,
		KVZ_Field:             60,
		PVS_Field:             75,
		PVZ_Field:             90,
		R_Field:               105,
		ZKF_Field:             120,
		VBEZ_Field:            135,
		VJAHR_Field:           150,
		PKPV_Field:            165,
		PKV_Field:             180,
		PVA_Field:             195,
		BackButtonField:       215,
		CalculateButtonField:  220,
	}
	
	// Get position for requested field
	position, ok := fieldPositions[field]
	if !ok {
		return
	}
	
	// Calculate viewport dimensions
	viewportHeight := m.advancedViewport.Height
	
	// Center the field in the viewport
	idealOffset := max(0, position - (viewportHeight / 2))
	
	// For buttons at bottom, position them in lower part of viewport
	if field == BackButtonField || field == CalculateButtonField {
		idealOffset = max(0, position - (viewportHeight * 3/4))
	}
	
	// Apply the scroll position
	m.advancedViewport.SetYOffset(idealOffset)
}

func (m *AppModel) View() string {
	switch m.step {
	case InputStep:
		return m.renderInputForm()
	case AdvancedInputStep:
		return m.renderAdvancedInputForm()
	case ResultsStep:
		return m.renderResults()
	case ComparisonStep:
		return m.renderComparison()
	default:
		return "Loading..."
	}
}