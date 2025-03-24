package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"tax-calculator/internal/adapters/api"
	"tax-calculator/internal/domain/models"
	"tax-calculator/internal/ui/styles"
)

// Step represents the current step in the application flow
type Step int

const (
	InputStep Step = iota
	ResultsStep
	ComparisonStep
)

// Field type represents form field focus
type Field int

const (
	TaxClassField Field = iota
	IncomeField
	YearField
	CalculateButtonField
)

// TaxClassOption represents a tax class option
type TaxClassOption struct {
	Class int
	Desc  string
}

// AppModel represents the unified main application model
type AppModel struct {
	step       Step
	focusField Field
	
	// Form fields
	taxClassOptions []TaxClassOption
	selectedTaxClass int
	incomeInput     textinput.Model
	yearInput       textinput.Model
	
	// Results data
	resultsLoading bool
	resultsError   string
	result         *api.TaxCalculationResponse
	resultsViewport viewport.Model
	showDetails     bool
	debugMessages   []string
	
	// Comparison data
	comparisonLoading bool
	comparisonError   string
	comparisonResults []models.TaxResult
	comparisonViewport viewport.Model
	completedCalls     int
	totalCalls         int
	
	// UI elements
	spinner     spinner.Model
	windowSize  tea.WindowSizeMsg
}

// NewAppModel creates a new unified application model
func NewAppModel() AppModel {
	// Set up tax class options
	taxClassOptions := []TaxClassOption{
		{Class: 1, Desc: "Single or permanently separated persons"},
		{Class: 2, Desc: "Single or permanently separated persons with child"},
		{Class: 3, Desc: "Married person (higher income)"},
		{Class: 4, Desc: "Married person (equal income)"},
		{Class: 5, Desc: "Married person (lower income)"},
		{Class: 6, Desc: "Person with multiple employments"},
	}

	// Set up income input
	incomeInput := textinput.New()
	incomeInput.Placeholder = "Enter income (e.g. 50000)"
	incomeInput.Width = 20
	incomeInput.CharLimit = 10
	incomeInput.Prompt = "€ "
	incomeInput.TextStyle = lipgloss.NewStyle().Foreground(styles.FgColor)
	
	// Set up year input
	yearInput := textinput.New()
	yearInput.Placeholder = "2025"
	yearInput.Width = 6
	yearInput.CharLimit = 4
	yearInput.SetValue("2025") // Default value
	yearInput.TextStyle = lipgloss.NewStyle().Foreground(styles.FgColor)

	// Set up spinner and viewport
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.PrimaryColor)

	// Use full width and height viewport for results
	vp := viewport.New(100, 40)
	vp.Style = styles.ResultsBoxStyle
	
	// Use full width and height viewport for comparison table
	compVp := viewport.New(100, 40)
	compVp.Style = styles.ResultsBoxStyle

	return AppModel{
		step:               InputStep,
		focusField:         TaxClassField,
		taxClassOptions:    taxClassOptions,
		selectedTaxClass:   1, // Default selection
		incomeInput:        incomeInput,
		yearInput:          yearInput,
		spinner:            s,
		resultsViewport:    vp,
		comparisonViewport: compVp,
		completedCalls:     0,
		totalCalls:         0,
		debugMessages:      []string{}, // Initialize empty debug messages
	}
}