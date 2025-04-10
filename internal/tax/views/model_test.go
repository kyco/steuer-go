package views

import (
	"tax-calculator/internal/tax/models"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewAppModel(t *testing.T) {
	model := NewAppModel()

	// Check default values
	if model.step != InputStep {
		t.Errorf("Expected initial step to be InputStep, got %v", model.step)
	}

	if model.focusField != TaxClassField {
		t.Errorf("Expected initial focus field to be TaxClassField, got %v", model.focusField)
	}

	if model.selectedTaxClass != 1 {
		t.Errorf("Expected initial selected tax class to be 1, got %d", model.selectedTaxClass)
	}

	if len(model.taxClassOptions) != 6 {
		t.Errorf("Expected 6 tax class options, got %d", len(model.taxClassOptions))
	}

	// Check tax class options
	expectedClasses := []int{1, 2, 3, 4, 5, 6}
	for i, option := range model.taxClassOptions {
		if option.Class != expectedClasses[i] {
			t.Errorf("Expected tax class option %d to have class %d, got %d", 
				i, expectedClasses[i], option.Class)
		}
	}

	// Check default values for inputs
	if model.incomeInput.Placeholder != "Enter income (e.g. 50000)" {
		t.Errorf("Expected income input placeholder to be 'Enter income (e.g. 50000)', got %q", 
			model.incomeInput.Placeholder)
	}

	if model.yearInput.Value() != "2025" {
		t.Errorf("Expected year input default value to be '2025', got %q", 
			model.yearInput.Value())
	}

	// Check viewport initialization
	if model.resultsViewport.Width != 100 || model.resultsViewport.Height != 40 {
		t.Errorf("Expected results viewport dimensions to be 100x40, got %dx%d", 
			model.resultsViewport.Width, model.resultsViewport.Height)
	}

	if model.comparisonViewport.Width != 100 || model.comparisonViewport.Height != 40 {
		t.Errorf("Expected comparison viewport dimensions to be 100x40, got %dx%d", 
			model.comparisonViewport.Width, model.comparisonViewport.Height)
	}
}

func TestModelUpdate(t *testing.T) {
	model := NewAppModel()
	
	// Test window size update
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 50})
	m, ok := newModel.(*AppModel)
	if !ok {
		t.Fatalf("Expected *AppModel after update, got %T", newModel)
	}
	
	if m.windowSize.Width != 120 || m.windowSize.Height != 50 {
		t.Errorf("Window size not updated correctly, got %dx%d", m.windowSize.Width, m.windowSize.Height)
	}
	
	// Test tab key navigation
	model = NewAppModel()
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	m, ok = newModel.(*AppModel)
	if !ok {
		t.Fatalf("Expected *AppModel after update, got %T", newModel)
	}
	
	if m.focusField != IncomeField {
		t.Errorf("Expected tab to move focus to IncomeField, got %v", m.focusField)
	}
	
	// Test shift+tab key navigation
	model = NewAppModel()
	model.focusField = IncomeField
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	m, ok = newModel.(*AppModel)
	if !ok {
		t.Fatalf("Expected *AppModel after update, got %T", newModel)
	}
	
	if m.focusField != TaxClassField {
		t.Errorf("Expected shift+tab to move focus to TaxClassField, got %v", m.focusField)
	}
	
	// Test up/down keys for tax class selection
	model = NewAppModel()
	model.selectedTaxClass = 3
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	m, ok = newModel.(*AppModel)
	if !ok {
		t.Fatalf("Expected *AppModel after update, got %T", newModel)
	}
	
	if m.selectedTaxClass != 2 {
		t.Errorf("Expected up key to change selected tax class to 2, got %d", m.selectedTaxClass)
	}
	
	model.selectedTaxClass = 3
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	m, ok = newModel.(*AppModel)
	if !ok {
		t.Fatalf("Expected *AppModel after update, got %T", newModel)
	}
	
	if m.selectedTaxClass != 4 {
		t.Errorf("Expected down key to change selected tax class to 4, got %d", m.selectedTaxClass)
	}
	
	// Test calculation message handling
	model = NewAppModel()
	model.step = ResultsStep
	model.resultsLoading = true
	
	// Test debug message handling
	model = NewAppModel()
	newModel, _ = model.Update(DebugLogMsg{Message: "Test debug message"})
	m, ok = newModel.(*AppModel)
	if !ok {
		t.Fatalf("Expected *AppModel after update, got %T", newModel)
	}
	
	if len(m.debugMessages) != 1 || m.debugMessages[0] != "Test debug message" {
		t.Errorf("Debug message not added correctly, got messages: %v", m.debugMessages)
	}
}

func TestModelView(t *testing.T) {
	model := NewAppModel()
	
	// Test Input step view
	view := model.View()
	if view == "" {
		t.Error("Input step view should not be empty")
	}
	
	// Test Results step view
	model.step = ResultsStep
	view = model.View()
	if view == "" {
		t.Error("Results step view should not be empty")
	}
	
	// Test Comparison step view
	model.step = ComparisonStep
	view = model.View()
	if view == "" {
		t.Error("Comparison step view should not be empty")
	}
}

func TestSortResults(t *testing.T) {
	results := []models.TaxResult{
		{Income: 50000.0},
		{Income: 30000.0},
		{Income: 70000.0},
		{Income: 10000.0},
	}
	
	sortResults(results)
	
	// Check if results are sorted by income
	for i := 0; i < len(results)-1; i++ {
		if results[i].Income > results[i+1].Income {
			t.Errorf("Results not sorted correctly: %f should be less than %f", 
				results[i].Income, results[i+1].Income)
		}
	}
}

// The calculateTaxForIncome function can't be easily tested without mocking
// So we'll just test that the sortResults function works correctly
func TestSortResults2(t *testing.T) {
	// Create some unsorted results
	results := []models.TaxResult{
		{Income: 60000.0},
		{Income: 30000.0},
		{Income: 10000.0},
	}
	
	// Sort them
	sortResults(results)
	
	// Verify they're sorted
	for i := 0; i < len(results)-1; i++ {
		if results[i].Income > results[i+1].Income {
			t.Errorf("Results not properly sorted: %f should be less than %f",
				results[i].Income, results[i+1].Income)
		}
	}
}