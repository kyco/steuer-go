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

type Step int

const (
	InputStep Step = iota
	AdvancedInputStep
	ResultsStep
	ComparisonStep
)

type Field int

const (
	TaxClassField Field = iota
	IncomeField
	YearField
	AdvancedButtonField
	CalculateButtonField
	
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
}

type AppModel struct {
	step       Step
	focusField Field

	taxClassOptions  []TaxClassOption
	selectedTaxClass int
	incomeInput      textinput.Model
	yearInput        textinput.Model
	useLocalCalc     bool

	// Advanced input fields
	ajahr  textinput.Model
	alter1 textinput.Model
	krv    textinput.Model
	kvz    textinput.Model
	pvs    textinput.Model
	pvz    textinput.Model
	r      textinput.Model
	zkf    textinput.Model
	vbez   textinput.Model
	vjahr  textinput.Model
	pkpv   textinput.Model
	pkv    textinput.Model
	pva    textinput.Model
	
	// Viewports for scrollable content
	advancedViewport viewport.Model
	resultsViewport  viewport.Model
	comparisonViewport viewport.Model
	
	resultsLoading  bool
	resultsError    string
	result          *bmf.TaxCalculationResponse
	showDetails     bool
	debugMessages   []string

	comparisonLoading  bool
	comparisonError    string
	comparisonResults  []models.TaxResult
	completedCalls     int
	totalCalls         int

	spinner    spinner.Model
	windowSize tea.WindowSizeMsg
}

func NewAppModel() *AppModel {
	taxClassOptions := []TaxClassOption{
		{Class: 1, Desc: "Single or permanently separated persons"},
		{Class: 2, Desc: "Single or permanently separated persons with child"},
		{Class: 3, Desc: "Married person (higher income)"},
		{Class: 4, Desc: "Married person (equal income)"},
		{Class: 5, Desc: "Married person (lower income)"},
		{Class: 6, Desc: "Person with multiple employments"},
	}

	incomeInput := textinput.New()
	incomeInput.Placeholder = "Enter income (e.g. 50000)"
	incomeInput.Width = 20
	incomeInput.CharLimit = 10
	incomeInput.Prompt = "€ "
	incomeInput.TextStyle = lipgloss.NewStyle().Foreground(styles.FgColor)

	currentYear := time.Now().Year()
	yearInput := textinput.New()
	yearInput.Placeholder = fmt.Sprintf("%d", currentYear)
	yearInput.Width = 6
	yearInput.CharLimit = 4
	yearInput.SetValue(fmt.Sprintf("%d", currentYear))
	yearInput.TextStyle = lipgloss.NewStyle().Foreground(styles.FgColor)
	
	// Create advanced input fields with defaults
	textStyle := lipgloss.NewStyle().Foreground(styles.FgColor)
	
	// Initialize all advanced inputs
	ajahr := textinput.New()
	ajahr.Placeholder = "0"
	ajahr.Width = 10
	ajahr.CharLimit = 4
	ajahr.TextStyle = textStyle

	alter1 := textinput.New()
	alter1.Placeholder = "0"
	alter1.Width = 5
	alter1.CharLimit = 1
	alter1.TextStyle = textStyle
	alter1.SetValue("0")

	krv := textinput.New()
	krv.Placeholder = "0"
	krv.Width = 5
	krv.CharLimit = 1
	krv.TextStyle = textStyle
	krv.SetValue("0")

	kvz := textinput.New()
	kvz.Placeholder = "1.3"
	kvz.Width = 10
	kvz.CharLimit = 5
	kvz.TextStyle = textStyle
	kvz.SetValue("1.3")

	pvs := textinput.New()
	pvs.Placeholder = "0"
	pvs.Width = 5
	pvs.CharLimit = 1
	pvs.TextStyle = textStyle
	pvs.SetValue("0")

	pvz := textinput.New()
	pvz.Placeholder = "0"
	pvz.Width = 5
	pvz.CharLimit = 1
	pvz.TextStyle = textStyle
	pvz.SetValue("0")

	r := textinput.New()
	r.Placeholder = "0"
	r.Width = 5
	r.CharLimit = 1
	r.TextStyle = textStyle
	r.SetValue("0")

	zkf := textinput.New()
	zkf.Placeholder = "0.0"
	zkf.Width = 10
	zkf.CharLimit = 5
	zkf.TextStyle = textStyle
	zkf.SetValue("0.0")

	vbez := textinput.New()
	vbez.Placeholder = "0"
	vbez.Width = 10
	vbez.CharLimit = 10
	vbez.TextStyle = textStyle
	vbez.SetValue("0")

	vjahr := textinput.New()
	vjahr.Placeholder = "0"
	vjahr.Width = 10
	vjahr.CharLimit = 4
	vjahr.TextStyle = textStyle
	vjahr.SetValue("0")

	pkpv := textinput.New()
	pkpv.Placeholder = "0"
	pkpv.Width = 10
	pkpv.CharLimit = 10
	pkpv.TextStyle = textStyle
	pkpv.SetValue("0")

	pkv := textinput.New()
	pkv.Placeholder = "0"
	pkv.Width = 5
	pkv.CharLimit = 1
	pkv.TextStyle = textStyle
	pkv.SetValue("0")

	pva := textinput.New()
	pva.Placeholder = "0"
	pva.Width = 5
	pva.CharLimit = 1
	pva.TextStyle = textStyle
	pva.SetValue("0")

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.PrimaryColor)

	// Initialize viewports
	resultsVp := viewport.New(100, 40)
	resultsVp.Style = styles.ResultsBoxStyle

	compVp := viewport.New(100, 40)
	compVp.Style = styles.ResultsBoxStyle
	
	advVp := viewport.New(100, 40)
	advVp.Style = styles.ResultsBoxStyle
	advVp.MouseWheelEnabled = true

	return &AppModel{
		step:               InputStep,
		focusField:         TaxClassField,
		taxClassOptions:    taxClassOptions,
		selectedTaxClass:   1, // Default selection
		incomeInput:        incomeInput,
		yearInput:          yearInput,
		useLocalCalc:       false,
		
		// Advanced inputs
		ajahr:              ajahr,
		alter1:             alter1,
		krv:                krv,
		kvz:                kvz,
		pvs:                pvs,
		pvz:                pvz,
		r:                  r,
		zkf:                zkf,
		vbez:               vbez,
		vjahr:              vjahr,
		pkpv:               pkpv,
		pkv:                pkv,
		pva:                pva,
		
		// Viewports
		resultsViewport:    resultsVp,
		comparisonViewport: compVp,
		advancedViewport:   advVp,
		
		spinner:            s,
		completedCalls:     0,
		totalCalls:         0,
		debugMessages:      []string{},
	}
}
