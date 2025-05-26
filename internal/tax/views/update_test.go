package views

import (
	"tax-calculator/internal/tax/models"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestRetroAppHandleTabNavigation(t *testing.T) {
	app := NewRetroApp()

	// Test main screen navigation
	app.screen = MainScreen
	app.focusField = TaxClassField

	// Forward navigation
	app.handleTabNavigation(false)
	if app.focusField != IncomeField {
		t.Errorf("Expected IncomeField, got %v", app.focusField)
	}

	app.handleTabNavigation(false)
	if app.focusField != YearField {
		t.Errorf("Expected YearField, got %v", app.focusField)
	}

	app.handleTabNavigation(false)
	if app.focusField != CalculateButtonField {
		t.Errorf("Expected CalculateButtonField, got %v", app.focusField)
	}

	app.handleTabNavigation(false)
	if app.focusField != AdvancedButtonField {
		t.Errorf("Expected AdvancedButtonField, got %v", app.focusField)
	}

	// Wrap around
	app.handleTabNavigation(false)
	if app.focusField != TaxClassField {
		t.Errorf("Expected TaxClassField after wrap around, got %v", app.focusField)
	}

	// Backward navigation
	app.handleTabNavigation(true)
	if app.focusField != AdvancedButtonField {
		t.Errorf("Expected AdvancedButtonField, got %v", app.focusField)
	}
}

func TestRetroAppHandleUpDownNavigation(t *testing.T) {
	app := NewRetroApp()

	// Test tax class selection
	app.screen = MainScreen
	app.focusField = TaxClassField
	app.selectedTaxClass = 3

	// Navigate up
	app.handleUpDownNavigation(true)
	if app.selectedTaxClass != 2 {
		t.Errorf("Expected tax class 2, got %d", app.selectedTaxClass)
	}

	// Navigate down
	app.handleUpDownNavigation(false)
	if app.selectedTaxClass != 3 {
		t.Errorf("Expected tax class 3, got %d", app.selectedTaxClass)
	}

	// Test wrapping at boundaries
	app.selectedTaxClass = 1
	app.handleUpDownNavigation(true)
	if app.selectedTaxClass != 6 {
		t.Errorf("Expected tax class 6 after wrap, got %d", app.selectedTaxClass)
	}

	app.selectedTaxClass = 6
	app.handleUpDownNavigation(false)
	if app.selectedTaxClass != 1 {
		t.Errorf("Expected tax class 1 after wrap, got %d", app.selectedTaxClass)
	}
}

func TestRetroAppHandleLeftRightNavigation(t *testing.T) {
	app := NewRetroApp()
	app.screen = ResultsScreen
	app.activeTab = BasicTab

	// Navigate right
	app.handleLeftRightNavigation(false)
	if app.activeTab != DetailsTab {
		t.Errorf("Expected DetailsTab, got %v", app.activeTab)
	}

	app.handleLeftRightNavigation(false)
	if app.activeTab != AboutTab {
		t.Errorf("Expected AboutTab, got %v", app.activeTab)
	}

	// Test wrapping
	app.handleLeftRightNavigation(false)
	if app.activeTab != BasicTab {
		t.Errorf("Expected BasicTab after wrap, got %v", app.activeTab)
	}

	// Navigate left
	app.handleLeftRightNavigation(true)
	if app.activeTab != AboutTab {
		t.Errorf("Expected AboutTab, got %v", app.activeTab)
	}
}

func TestRetroAppHandleEnterSelection(t *testing.T) {
	app := NewRetroApp()

	// Test main screen enter handling
	app.screen = MainScreen
	app.focusField = IncomeField

	cmd := app.handleEnterSelection()
	if cmd != nil {
		t.Error("Expected no command for income field enter")
	}

	if app.focusField != YearField {
		t.Errorf("Expected focus to move to YearField, got %v", app.focusField)
	}

	// Test year field
	app.focusField = YearField
	cmd = app.handleEnterSelection()
	if cmd != nil {
		t.Error("Expected no command for year field enter")
	}

	if app.focusField != CalculateButtonField {
		t.Errorf("Expected focus to move to CalculateButtonField, got %v", app.focusField)
	}

	// Test advanced button
	app.focusField = AdvancedButtonField
	cmd = app.handleEnterSelection()
	if cmd != nil {
		t.Error("Expected no command for advanced button")
	}

	if app.screen != AdvancedScreen {
		t.Errorf("Expected screen to change to AdvancedScreen, got %v", app.screen)
	}

	if app.focusField != AJAHR_Field {
		t.Errorf("Expected focus to move to AJAHR_Field, got %v", app.focusField)
	}
}

func TestRetroAppBlurAllInputs(t *testing.T) {
	app := NewRetroApp()

	// Test main screen blur
	app.screen = MainScreen
	app.incomeInput.Focus()
	app.yearInput.Focus()

	app.blurAllInputs()

	if app.incomeInput.Focused() {
		t.Error("Expected income input to be blurred")
	}

	if app.yearInput.Focused() {
		t.Error("Expected year input to be blurred")
	}
}

func TestRetroAppAutoFocusInputField(t *testing.T) {
	app := NewRetroApp()

	// Test main screen auto focus
	app.screen = MainScreen
	app.focusField = IncomeField

	app.autoFocusInputField()

	if !app.incomeInput.Focused() {
		t.Error("Expected income input to be focused")
	}

	if app.yearInput.Focused() {
		t.Error("Expected year input to remain blurred")
	}

	// Test year field focus
	app.focusField = YearField
	app.autoFocusInputField()

	if app.incomeInput.Focused() {
		t.Error("Expected income input to be blurred")
	}

	if !app.yearInput.Focused() {
		t.Error("Expected year input to be focused")
	}
}

func TestRetroAppNavigateFields(t *testing.T) {
	app := NewRetroApp()

	fields := []Field{TaxClassField, IncomeField, YearField}
	app.focusField = TaxClassField

	// Forward navigation
	app.navigateFields(fields, false)
	if app.focusField != IncomeField {
		t.Errorf("Expected IncomeField, got %v", app.focusField)
	}

	// Backward navigation
	app.navigateFields(fields, true)
	if app.focusField != TaxClassField {
		t.Errorf("Expected TaxClassField, got %v", app.focusField)
	}

	// Test wrap around backward
	app.navigateFields(fields, true)
	if app.focusField != YearField {
		t.Errorf("Expected YearField after backward wrap, got %v", app.focusField)
	}

	// Test wrap around forward
	app.navigateFields(fields, false)
	if app.focusField != TaxClassField {
		t.Errorf("Expected TaxClassField after forward wrap, got %v", app.focusField)
	}
}

func TestRetroAppUpdateViewportDimensions(t *testing.T) {
	app := NewRetroApp()

	msg := tea.WindowSizeMsg{
		Width:  100,
		Height: 50,
	}

	app.updateViewportDimensions(msg)

	expectedWidth := 80  // 100 - 20
	expectedHeight := 30 // 50 - 20

	if app.resultsViewport.Width != expectedWidth {
		t.Errorf("Expected results viewport width %d, got %d", expectedWidth, app.resultsViewport.Width)
	}

	if app.resultsViewport.Height != expectedHeight {
		t.Errorf("Expected results viewport height %d, got %d", expectedHeight, app.resultsViewport.Height)
	}

	if app.comparisonViewport.Width != expectedWidth {
		t.Errorf("Expected comparison viewport width %d, got %d", expectedWidth, app.comparisonViewport.Width)
	}

	if app.comparisonViewport.Height != expectedHeight {
		t.Errorf("Expected comparison viewport height %d, got %d", expectedHeight, app.comparisonViewport.Height)
	}

	if app.advancedViewport.Width != expectedWidth {
		t.Errorf("Expected advanced viewport width %d, got %d", expectedWidth, app.advancedViewport.Width)
	}

	if app.advancedViewport.Height != expectedHeight {
		t.Errorf("Expected advanced viewport height %d, got %d", expectedHeight, app.advancedViewport.Height)
	}
}

func TestRetroAppUpdateViewportDimensionsMinimum(t *testing.T) {
	app := NewRetroApp()

	// Test minimum dimensions
	msg := tea.WindowSizeMsg{
		Width:  10,
		Height: 5,
	}

	app.updateViewportDimensions(msg)

	if app.resultsViewport.Width != 20 {
		t.Errorf("Expected minimum width 20, got %d", app.resultsViewport.Width)
	}

	if app.resultsViewport.Height != 10 {
		t.Errorf("Expected minimum height 10, got %d", app.resultsViewport.Height)
	}
}

func TestRetroAppComparisonNavigation(t *testing.T) {
	app := NewRetroApp()
	app.screen = ComparisonScreen

	// Create some mock comparison results
	app.comparisonResults = []models.TaxResult{
		{Income: 30000},
		{Income: 50000},
		{Income: 70000},
	}

	app.selectedComparisonIdx = 0

	// Navigate down
	app.handleUpDownNavigation(false)
	if app.selectedComparisonIdx != 1 {
		t.Errorf("Expected selectedComparisonIdx 1, got %d", app.selectedComparisonIdx)
	}

	// Navigate up
	app.handleUpDownNavigation(true)
	if app.selectedComparisonIdx != 0 {
		t.Errorf("Expected selectedComparisonIdx 0, got %d", app.selectedComparisonIdx)
	}

	// Test wrap around up
	app.handleUpDownNavigation(true)
	if app.selectedComparisonIdx != 2 {
		t.Errorf("Expected selectedComparisonIdx 2 after wrap up, got %d", app.selectedComparisonIdx)
	}

	// Test wrap around down
	app.handleUpDownNavigation(false)
	if app.selectedComparisonIdx != 0 {
		t.Errorf("Expected selectedComparisonIdx 0 after wrap down, got %d", app.selectedComparisonIdx)
	}
}

func TestRetroAppComparisonEnterHandling(t *testing.T) {
	app := NewRetroApp()
	app.screen = ComparisonScreen

	// Create some mock comparison results
	app.comparisonResults = []models.TaxResult{
		{Income: 30000},
		{Income: 50000},
	}

	app.showBreakdown = false

	// Test enter toggles breakdown
	cmd := app.handleEnterSelection()
	if cmd != nil {
		t.Error("Expected no command from comparison enter")
	}

	if !app.showBreakdown {
		t.Error("Expected showBreakdown to be true after enter")
	}

	// Test enter toggles back
	cmd = app.handleEnterSelection()
	if cmd != nil {
		t.Error("Expected no command from comparison enter")
	}

	if app.showBreakdown {
		t.Error("Expected showBreakdown to be false after second enter")
	}
}

func TestRetroAppComparisonWithEmptyResults(t *testing.T) {
	app := NewRetroApp()
	app.screen = ComparisonScreen
	app.comparisonResults = []models.TaxResult{}

	// Navigation should not crash with empty results
	app.handleUpDownNavigation(true)
	app.handleUpDownNavigation(false)

	// Enter should not crash with empty results
	cmd := app.handleEnterSelection()
	if cmd != nil {
		t.Error("Expected no command with empty comparison results")
	}
}

func TestRetroAppAdvancedScreenNavigation(t *testing.T) {
	app := NewRetroApp()
	app.screen = AdvancedScreen

	// Test back button
	app.focusField = BackButtonField
	cmd := app.handleEnterSelection()
	if cmd != nil {
		t.Error("Expected no command for back button")
	}

	if app.screen != MainScreen {
		t.Errorf("Expected screen to change to MainScreen, got %v", app.screen)
	}

	if app.focusField != TaxClassField {
		t.Errorf("Expected focus to move to TaxClassField, got %v", app.focusField)
	}
}

func TestFieldConstants(t *testing.T) {
	// Test that field constants are defined correctly
	fields := []Field{
		TaxClassField,
		IncomeField,
		YearField,
		CalculateButtonField,
		AdvancedButtonField,
		BackButtonField,
		AJAHR_Field,
		ALTER1_Field,
		KRV_Field,
		KVZ_Field,
		PVS_Field,
		PVZ_Field,
		R_Field,
		ZKF_Field,
		VBEZ_Field,
		VJAHR_Field,
		PKPV_Field,
		PKV_Field,
		PVA_Field,
	}

	for i, field := range fields {
		if field == Field(0) && i > 0 {
			t.Errorf("Field at index %d is zero value", i)
		}
	}
}

func TestTabConstants(t *testing.T) {
	// Test that tab constants are defined correctly
	tabs := []Tab{
		BasicTab,
		DetailsTab,
		AboutTab,
	}

	for i, tab := range tabs {
		if tab == Tab(0) && i > 0 {
			t.Errorf("Tab at index %d is zero value", i)
		}
	}
}

func TestScreenConstants(t *testing.T) {
	// Test that screen constants are defined correctly
	screens := []Screen{
		MainScreen,
		ResultsScreen,
		ComparisonScreen,
		AdvancedScreen,
	}

	for i, screen := range screens {
		if screen == Screen(0) && i > 0 {
			t.Errorf("Screen at index %d is zero value", i)
		}
	}
}
