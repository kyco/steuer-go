package views

import (
	"tax-calculator/internal/tax/models"
	"testing"
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
	cmd := PerformCalculationCmd(1, 50000, "2025")
	msg := cmd()
	
	_, ok := msg.(CalculationStartedMsg)
	if !ok {
		t.Fatalf("Expected CalculationStartedMsg, got %T", msg)
	}
}

func TestProgressUpdateCmd(t *testing.T) {
	cmd := ProgressUpdateCmd(5, 10)
	msg := cmd()
	
	progressMsg, ok := msg.(ComparisonProgressMsg)
	if !ok {
		t.Fatalf("Expected ComparisonProgressMsg, got %T", msg)
	}
	
	if progressMsg.CompletedCalls != 5 || progressMsg.TotalCalls != 10 {
		t.Errorf("Expected progress 5/10, got %d/%d", 
			progressMsg.CompletedCalls, progressMsg.TotalCalls)
	}
}

func TestCompletedResultsCmd(t *testing.T) {
	results := []models.TaxResult{
		{Income: 50000.0, IncomeTax: 8000.0},
	}
	
	cmd := CompletedResultsCmd(results)
	msg := cmd()
	
	compMsg, ok := msg.(ComparisonMsg)
	if !ok {
		t.Fatalf("Expected ComparisonMsg, got %T", msg)
	}
	
	if len(compMsg.Results) != 1 || compMsg.Results[0].Income != 50000.0 {
		t.Errorf("Expected 1 result with income 50000.0, got %d results with first income %f", 
			len(compMsg.Results), compMsg.Results[0].Income)
	}
}

func TestPerformComparisonCmd(t *testing.T) {
	cmd := PerformComparisonCmd()
	msg := cmd()
	
	_, ok := msg.(ComparisonStartedMsg)
	if !ok {
		t.Fatalf("Expected ComparisonStartedMsg, got %T", msg)
	}
}