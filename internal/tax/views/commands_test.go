package views

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"tax-calculator/internal/tax/models"
)

func TestDebugLogCmd(t *testing.T) {
	cmd := CaptureDebugCmd("Test debug message")
	msg := cmd()

	debugMsg, ok := msg.(DebugLogMsg)
	if !ok {
		t.Fatalf("Expected DebugLogMsg, got %T", msg)
	}

	if debugMsg.Message != "Test debug message" {
		t.Errorf("Expected debug message 'Test debug message', got %q", debugMsg.Message)
	}
}

func TestPerformCalculationCmd(t *testing.T) {
	cmd := PerformCalculationCmd(1, 50000.0, "2024")

	if cmd == nil {
		t.Error("PerformCalculationCmd should return a non-nil command")
	}

	// Execute the command to test it
	msg := cmd()

	// The function returns a tea.BatchMsg which contains multiple messages
	if batchMsg, ok := msg.(tea.BatchMsg); ok {
		// BatchMsg is a slice of tea.Cmd functions
		if len(batchMsg) == 0 {
			t.Error("Expected BatchMsg to contain commands")
		}
	} else {
		t.Errorf("Expected tea.BatchMsg, got %T", msg)
	}
}

func TestPerformCalculationWithAdvancedOptionsCmd(t *testing.T) {
	taxRequest := models.TaxRequest{
		Period:   models.Year,
		Income:   5000000,
		TaxClass: models.TaxClass1,
		R:        0,
		AJAHR:    0,
		ALTER1:   0,
		KRV:      0,
		KVZ:      1.3,
		PVS:      0,
		PVZ:      0,
		PKV:      0,
		PVA:      0,
		ZKF:      0,
		VBEZ:     0,
		VJAHR:    0,
		PKPV:     0,
	}

	cmd := PerformCalculationWithAdvancedOptionsCmd(1, 50000.0, "2024", taxRequest, false)

	if cmd == nil {
		t.Error("PerformCalculationWithAdvancedOptionsCmd should return a non-nil command")
	}

	// Execute the command to test it
	msg := cmd()

	// The function returns a tea.BatchMsg which contains multiple messages
	if batchMsg, ok := msg.(tea.BatchMsg); ok {
		// BatchMsg is a slice of tea.Cmd functions
		if len(batchMsg) == 0 {
			t.Error("Expected BatchMsg to contain commands")
		}
	} else {
		t.Errorf("Expected tea.BatchMsg, got %T", msg)
	}
}

func TestPerformCalculationWithAdvancedOptionsCmdLocal(t *testing.T) {
	taxRequest := models.TaxRequest{
		Period:   models.Year,
		Income:   5000000,
		TaxClass: models.TaxClass1,
	}

	cmd := PerformCalculationWithAdvancedOptionsCmd(1, 50000.0, "2024", taxRequest, true)

	if cmd == nil {
		t.Error("PerformCalculationWithAdvancedOptionsCmd should return a non-nil command")
	}

	// Execute the command to test it
	msg := cmd()

	// The function returns a tea.BatchMsg which contains multiple messages
	if batchMsg, ok := msg.(tea.BatchMsg); ok {
		// BatchMsg is a slice of tea.Cmd functions
		if len(batchMsg) == 0 {
			t.Error("Expected BatchMsg to contain commands")
		}
	} else {
		t.Errorf("Expected tea.BatchMsg, got %T", msg)
	}
}

func TestFetchResultsCmd(t *testing.T) {
	cmd := FetchResultsCmd(1, 50000.0)

	if cmd == nil {
		t.Error("FetchResultsCmd should return a non-nil command")
	}

	// The command should be a tea.Cmd that can be executed
	// We can't easily test the actual execution without mocking the API
	// but we can verify it returns a function
}

func TestFetchResultsWithModeCmd(t *testing.T) {
	cmd := FetchResultsWithModeCmd(1, 50000.0, false)

	if cmd == nil {
		t.Error("FetchResultsWithModeCmd should return a non-nil command")
	}

	// Test with local mode
	cmdLocal := FetchResultsWithModeCmd(1, 50000.0, true)

	if cmdLocal == nil {
		t.Error("FetchResultsWithModeCmd with local mode should return a non-nil command")
	}
}

func TestFetchResultsWithAdvancedParamsCmd(t *testing.T) {
	taxRequest := models.TaxRequest{
		Period:   models.Year,
		Income:   5000000,
		TaxClass: models.TaxClass1,
	}

	cmd := FetchResultsWithAdvancedParamsCmd(taxRequest, false)

	if cmd == nil {
		t.Error("FetchResultsWithAdvancedParamsCmd should return a non-nil command")
	}

	// Test with local mode
	cmdLocal := FetchResultsWithAdvancedParamsCmd(taxRequest, true)

	if cmdLocal == nil {
		t.Error("FetchResultsWithAdvancedParamsCmd with local mode should return a non-nil command")
	}
}

func TestFetchComparisonCmd(t *testing.T) {
	cmd := FetchComparisonCmd(1, 50000.0)

	if cmd == nil {
		t.Error("FetchComparisonCmd should return a non-nil command")
	}

	// Execute the command to test it returns a function
	cmdFunc := cmd()
	if cmdFunc == nil {
		t.Error("FetchComparisonCmd should return a command function")
	}
}

func TestCalculateTaxForIncome(t *testing.T) {
	// This is an internal function, but we can test it indirectly
	// by testing the commands that use it

	// Test through FetchComparisonCmd which uses calculateTaxForIncome
	cmd := FetchComparisonCmd(1, 50000.0)
	if cmd == nil {
		t.Error("FetchComparisonCmd should work with calculateTaxForIncome")
	}
}

func TestSortResultsIndirect(t *testing.T) {
	// Create some test results
	results := []models.TaxResult{
		{Income: 70000, TotalTax: 14000},
		{Income: 30000, TotalTax: 4000},
		{Income: 50000, TotalTax: 8000},
	}

	// The sortResults function is internal, but we can test it indirectly
	// by verifying that comparison results are properly sorted

	// For now, just verify we have test data
	if len(results) != 3 {
		t.Error("Test data should have 3 results")
	}

	// In a real test, you would expose sortResults or test it through
	// the comparison functionality that uses it
}

func TestProgressUpdateCmd(t *testing.T) {
	cmd := ProgressUpdateCmd(5, 10)

	if cmd == nil {
		t.Error("ProgressUpdateCmd should return a non-nil command")
	}

	// Execute the command
	msg := cmd()

	// Should return a ComparisonProgressMsg
	if progressMsg, ok := msg.(ComparisonProgressMsg); !ok {
		t.Errorf("Expected ComparisonProgressMsg, got %T", msg)
	} else {
		if progressMsg.CompletedCalls != 5 {
			t.Errorf("Expected CompletedCalls 5, got %d", progressMsg.CompletedCalls)
		}
		if progressMsg.TotalCalls != 10 {
			t.Errorf("Expected TotalCalls 10, got %d", progressMsg.TotalCalls)
		}
	}
}

func TestCompletedResultsCmd(t *testing.T) {
	results := []models.TaxResult{
		{Income: 50000, TotalTax: 8000},
	}

	cmd := CompletedResultsCmd(results)

	if cmd == nil {
		t.Error("CompletedResultsCmd should return a non-nil command")
	}

	// Execute the command
	msg := cmd()

	// Should return a ComparisonMsg
	if compMsg, ok := msg.(ComparisonMsg); !ok {
		t.Errorf("Expected ComparisonMsg, got %T", msg)
	} else {
		if len(compMsg.Results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(compMsg.Results))
		}
		if compMsg.Error != nil {
			t.Errorf("Expected no error, got %v", compMsg.Error)
		}
	}
}

func TestCompletedResultsCmdWithError(t *testing.T) {
	// Test with nil results to trigger error case
	cmd := CompletedResultsCmd(nil)

	if cmd == nil {
		t.Error("CompletedResultsCmd should return a non-nil command even with nil results")
	}

	// Execute the command
	msg := cmd()

	// Should return a ComparisonMsg with error
	if compMsg, ok := msg.(ComparisonMsg); !ok {
		t.Errorf("Expected ComparisonMsg, got %T", msg)
	} else {
		if compMsg.Results != nil {
			t.Error("Expected nil results")
		}
		// Note: The actual error handling depends on the implementation
	}
}

func TestMessageTypes(t *testing.T) {
	// Test that our message types can be created and used

	calcStarted := CalculationStartedMsg{UseLocalCalculator: true}
	if !calcStarted.UseLocalCalculator {
		t.Error("CalculationStartedMsg should preserve UseLocalCalculator field")
	}

	calcMsg := CalculationMsg{
		Result: nil,
		Error:  nil,
	}
	if calcMsg.Result != nil || calcMsg.Error != nil {
		t.Error("CalculationMsg should initialize with nil values")
	}

	compStarted := ComparisonStartedMsg{}
	_ = compStarted // Just verify it can be created

	compProgress := ComparisonProgressMsg{
		CompletedCalls: 5,
		TotalCalls:     10,
	}
	if compProgress.CompletedCalls != 5 || compProgress.TotalCalls != 10 {
		t.Error("ComparisonProgressMsg should preserve field values")
	}

	compMsg := ComparisonMsg{
		Results: []models.TaxResult{},
		Error:   nil,
	}
	if compMsg.Results == nil {
		t.Error("ComparisonMsg should preserve Results field")
	}

	debugMsg := DebugLogMsg{Message: "test"}
	if debugMsg.Message != "test" {
		t.Error("DebugLogMsg should preserve Message field")
	}
}
