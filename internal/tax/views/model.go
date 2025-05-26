package views

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"tax-calculator/internal/tax/bmf"
	"tax-calculator/internal/tax/models"
	"tax-calculator/internal/tax/views/styles"
)

type Screen int

const (
	MainScreen Screen = iota
	ResultsScreen
	ComparisonScreen
	AdvancedScreen
)

type Tab int

const (
	BasicTab Tab = iota
	DetailsTab
	AboutTab
)

type Field int

const (
	TaxClassField Field = iota
	IncomeField
	YearField
	CalculateButtonField
	AdvancedButtonField

	// Advanced input fields
	AJAHR_Field
	ALTER1_Field
	KRV_Field
	KVZ_Field
	PVS_Field
	PVZ_Field
	R_Field
	ZKF_Field
	VBEZ_Field
	VJAHR_Field
	PKPV_Field
	PKV_Field
	PVA_Field
	BackButtonField
)

type TaxClassOption struct {
	Class int
	Desc  string
	Icon  string
}

type AdvancedField struct {
	Label       string
	Description string
	Model       textinput.Model
	Field       Field
}

type RetroApp struct {
	screen     Screen
	activeTab  Tab
	focusField Field

	taxClassOptions  []TaxClassOption
	selectedTaxClass int
	incomeInput      textinput.Model
	yearInput        textinput.Model
	useLocalCalc     bool

	// Advanced input fields in a more organized structure
	advancedFields []AdvancedField

	// Viewports for scrollable content
	mainViewport       viewport.Model
	resultsViewport    viewport.Model
	advancedViewport   viewport.Model
	comparisonViewport viewport.Model

	resultsLoading bool
	resultsError   string
	result         *bmf.TaxCalculationResponse
	showDetails    bool

	comparisonLoading     bool
	comparisonError       string
	comparisonResults     []models.TaxResult
	selectedComparisonIdx int
	showBreakdown         bool
	completedCalls        int
	totalCalls            int

	spinner    spinner.Model
	windowSize tea.WindowSizeMsg
}

func NewRetroApp() *RetroApp {
	taxClassOptions := []TaxClassOption{
		{Class: 1, Desc: "Single, separated", Icon: "👤"},
		{Class: 2, Desc: "Single with child", Icon: "👨‍👦"},
		{Class: 3, Desc: "Married (higher income)", Icon: "💍"},
		{Class: 4, Desc: "Married (equal income)", Icon: "👫"},
		{Class: 5, Desc: "Married (lower income)", Icon: "👪"},
		{Class: 6, Desc: "Multiple employments", Icon: "🏢"},
	}

	// Custom text input style
	incomeInput := textinput.New()
	incomeInput.Placeholder = "50000"
	incomeInput.Width = 20
	incomeInput.CharLimit = 10
	incomeInput.Prompt = "€ "
	incomeInput.TextStyle = lipgloss.NewStyle().Foreground(styles.FgColor)
	incomeInput.PromptStyle = lipgloss.NewStyle().Foreground(styles.SuccessColor)

	// Current year with properly styled input
	currentYear := time.Now().Year()
	yearInput := textinput.New()
	yearInput.Placeholder = fmt.Sprintf("%d", currentYear)
	yearInput.Width = 6
	yearInput.CharLimit = 4
	yearInput.SetValue(fmt.Sprintf("%d", currentYear))
	yearInput.TextStyle = lipgloss.NewStyle().Foreground(styles.FgColor)
	yearInput.PromptStyle = lipgloss.NewStyle().Foreground(styles.AccentColor)

	// Initialize advanced input fields with organized structure

	// Build advanced fields data structure
	advancedFields := []AdvancedField{
		createAdvancedField(
			"Year after 64th birthday",
			"Enter year after taxpayer's 64th birthday (0 if not applicable)",
			"0", 10, 4, AJAHR_Field),

		createAdvancedField(
			"Completed 64 years",
			"Enter 1 if taxpayer was 64+ at start of calendar year, 0 otherwise",
			"0", 5, 1, ALTER1_Field),

		createAdvancedField(
			"Social insurance type",
			"0: Normal statutory pension, 1: No compulsory insurance, 2: Reduced rate",
			"0", 5, 1, KRV_Field),

		createAdvancedField(
			"Additional health rate %",
			"Additional health insurance percentage (standard is 1.3%)",
			"1.3", 10, 5, KVZ_Field),

		createAdvancedField(
			"Employer in Saxony",
			"1 if employer is in Saxony (different nursing care), 0 otherwise",
			"0", 5, 1, PVS_Field),

		createAdvancedField(
			"Childless surcharge",
			"1 if employee (23+) pays childless surcharge for nursing care",
			"0", 5, 1, PVZ_Field),

		createAdvancedField(
			"Religion code",
			"0: No church tax, 1: Catholic church, 2: Protestant church",
			"0", 5, 1, R_Field),

		createAdvancedField(
			"Child allowance",
			"Number of children for tax allowance (can be decimal, e.g. 0.5)",
			"0.0", 10, 5, ZKF_Field),

		createAdvancedField(
			"Pension payments €",
			"Annual pension income in euros (0 if none)",
			"0", 10, 10, VBEZ_Field),

		createAdvancedField(
			"First pension year",
			"Year when taxpayer first received pension (0 if not applicable)",
			"0", 10, 4, VJAHR_Field),

		createAdvancedField(
			"Private insurance payment €",
			"Monthly private health insurance premium in euros",
			"0", 10, 10, PKPV_Field),

		createAdvancedField(
			"Health insurance type",
			"0: Statutory, 1: Private without employer subsidy, 2: Private with subsidy",
			"0", 5, 1, PKV_Field),

		createAdvancedField(
			"Children for care insurance",
			"Number of children for reduced nursing care insurance (0-4)",
			"0", 5, 1, PVA_Field),
	}

	// Create fancy spinner
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = styles.SpinnerStyle

	// Initialize viewports with proper styling
	mainVp := viewport.New(100, 40)
	mainVp.Style = styles.BaseStyle

	resultsVp := viewport.New(100, 40)
	resultsVp.Style = styles.BaseStyle

	compVp := viewport.New(100, 40)
	compVp.Style = styles.BaseStyle

	advVp := viewport.New(100, 40)
	advVp.Style = styles.BaseStyle
	advVp.MouseWheelEnabled = true

	return &RetroApp{
		screen:           MainScreen,
		activeTab:        BasicTab,
		focusField:       TaxClassField,
		taxClassOptions:  taxClassOptions,
		selectedTaxClass: 1, // Default selection
		incomeInput:      incomeInput,
		yearInput:        yearInput,
		useLocalCalc:     false,

		// Advanced fields
		advancedFields: advancedFields,

		// Viewports
		mainViewport:       mainVp,
		resultsViewport:    resultsVp,
		comparisonViewport: compVp,
		advancedViewport:   advVp,

		spinner:               s,
		selectedComparisonIdx: 0,
		showBreakdown:         false,
		completedCalls:        0,
		totalCalls:            0,
	}
}

// Helper function to create advanced input fields with consistent styling
func createAdvancedField(label, description, defaultValue string, width, charLimit int, fieldType Field) AdvancedField {
	input := textinput.New()
	input.Placeholder = defaultValue
	input.Width = width
	input.CharLimit = charLimit
	input.TextStyle = lipgloss.NewStyle().Foreground(styles.FgColor)
	input.PromptStyle = lipgloss.NewStyle().Foreground(styles.AccentColor)
	input.SetValue(defaultValue)

	return AdvancedField{
		Label:       label,
		Description: description,
		Model:       input,
		Field:       fieldType,
	}
}

// Helper method to get a specific advanced field by its field type
func (m *RetroApp) getAdvancedField(field Field) *AdvancedField {
	for i := range m.advancedFields {
		if m.advancedFields[i].Field == field {
			return &m.advancedFields[i]
		}
	}
	return nil
}

// Helper method to build a tax request with all parameters
func (m *RetroApp) buildTaxRequest() models.TaxRequest {
	income, _ := parseFloatWithDefault(m.incomeInput.Value(), 0)

	// Basic request
	request := models.TaxRequest{
		Period:   models.Year,
		Income:   int(income * 100),
		TaxClass: models.TaxClass(m.selectedTaxClass),
	}

	// Add advanced parameters
	if field := m.getAdvancedField(AJAHR_Field); field != nil {
		request.AJAHR, _ = parseIntWithDefault(field.Model.Value(), 0)
	}

	if field := m.getAdvancedField(ALTER1_Field); field != nil {
		request.ALTER1, _ = parseIntWithDefault(field.Model.Value(), 0)
	}

	if field := m.getAdvancedField(KRV_Field); field != nil {
		request.KRV, _ = parseIntWithDefault(field.Model.Value(), 0)
	}

	if field := m.getAdvancedField(KVZ_Field); field != nil {
		request.KVZ, _ = parseFloatWithDefault(field.Model.Value(), 1.3)
	}

	if field := m.getAdvancedField(PVS_Field); field != nil {
		request.PVS, _ = parseIntWithDefault(field.Model.Value(), 0)
	}

	if field := m.getAdvancedField(PVZ_Field); field != nil {
		request.PVZ, _ = parseIntWithDefault(field.Model.Value(), 0)
	}

	if field := m.getAdvancedField(R_Field); field != nil {
		request.R, _ = parseIntWithDefault(field.Model.Value(), 0)
	}

	if field := m.getAdvancedField(ZKF_Field); field != nil {
		request.ZKF, _ = parseFloatWithDefault(field.Model.Value(), 0)
	}

	if field := m.getAdvancedField(VBEZ_Field); field != nil {
		request.VBEZ, _ = parseIntWithDefault(field.Model.Value(), 0)
	}

	if field := m.getAdvancedField(VJAHR_Field); field != nil {
		request.VJAHR, _ = parseIntWithDefault(field.Model.Value(), 0)
	}

	if field := m.getAdvancedField(PKPV_Field); field != nil {
		request.PKPV, _ = parseIntWithDefault(field.Model.Value(), 0)
	}

	if field := m.getAdvancedField(PKV_Field); field != nil {
		request.PKV, _ = parseIntWithDefault(field.Model.Value(), 0)
	}

	if field := m.getAdvancedField(PVA_Field); field != nil {
		request.PVA, _ = parseIntWithDefault(field.Model.Value(), 0)
	}

	return request
}
