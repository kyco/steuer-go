package views

import (
	"testing"
)

func TestParseIncome(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"100", 100.0},
		{"50000.50", 50000.50},
		{"   45000   ", 45000.0},
		{"", 0.0},
		{"invalid", 0.0},
	}

	for _, tc := range tests {
		result := parseIncome(tc.input)
		if result != tc.expected {
			t.Errorf("parseIncome(%q): expected %f, got %f", tc.input, tc.expected, result)
		}
	}
}

func TestMinAndMax(t *testing.T) {
	minTests := []struct {
		a, b, expected int
	}{
		{5, 10, 5},
		{10, 5, 5},
		{0, 0, 0},
		{-5, -10, -10},
		{-10, -5, -10},
	}

	for _, tc := range minTests {
		result := min(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("min(%d, %d): expected %d, got %d", tc.a, tc.b, tc.expected, result)
		}
	}

	maxTests := []struct {
		a, b, expected int
	}{
		{5, 10, 10},
		{10, 5, 10},
		{0, 0, 0},
		{-5, -10, -5},
		{-10, -5, -5},
	}

	for _, tc := range maxTests {
		result := max(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("max(%d, %d): expected %d, got %d", tc.a, tc.b, tc.expected, result)
		}
	}
}

func TestValidateAndCalculate(t *testing.T) {
	tests := []struct {
		income        string
		year          string
		expectedValid bool
		expectedError string
	}{
		{"50000", "2025", true, ""},
		{"", "2025", false, "Income cannot be empty"},
		{"invalid", "2025", false, "Income must be a positive number"},
		{"-1000", "2025", false, "Income must be a positive number"},
		{"0", "2025", false, "Income must be a positive number"},
		{"50000", "", true, ""}, // Should default to 2025
		{"50000", "invalid", false, "Year must be between 2024 and 2030"},
		{"50000", "2023", false, "Year must be between 2024 and 2030"},
		{"50000", "2031", false, "Year must be between 2024 and 2030"},
	}

	for _, tc := range tests {
		model := NewAppModel()
		model.incomeInput.SetValue(tc.income)
		model.yearInput.SetValue(tc.year)

		valid, errMsg := model.validateAndCalculate()

		if valid != tc.expectedValid {
			t.Errorf("validateAndCalculate(%q, %q): expected valid=%v, got valid=%v", 
				tc.income, tc.year, tc.expectedValid, valid)
		}

		if errMsg != tc.expectedError {
			t.Errorf("validateAndCalculate(%q, %q): expected error=%q, got error=%q", 
				tc.income, tc.year, tc.expectedError, errMsg)
		}

		// Check if empty year is defaulted to 2025
		if tc.year == "" && valid {
			if model.yearInput.Value() != "2025" {
				t.Errorf("validateAndCalculate with empty year: expected year to default to 2025, got %q", 
					model.yearInput.Value())
			}
		}
	}
}

func TestAdvancedViewportScrolling(t *testing.T) {
	// Create a new app model
	model := NewAppModel()
	
	// Set window height smaller than the full content
	model.windowSize.Height = 30
	model.windowSize.Width = 100
	
	// Switch to advanced view
	model.step = AdvancedInputStep
	
	// Set the viewport height based on window size
	model.advancedViewport.Height = model.windowSize.Height - 8
	
	// Test different focus fields to ensure scrolling works
	testCases := []struct {
		name       string
		focusField Field
		wantScroll bool
	}{
		{"Focus on first field", AJAHR_Field, false},
		{"Focus on middle field", ZKF_Field, true},
		{"Focus on last field", PVA_Field, true},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set initial scroll position to top
			model.advancedViewport.SetYOffset(0)
			
			// Render the form with the model as pointer
			// This should update the viewport's content and scroll position
			content := model.renderAdvancedInputForm()
			
			// Focus on the test field and scroll to it
			prevField := model.focusField
			model.focusField = tc.focusField
			model.scrollToAdvancedField(prevField, tc.focusField)
			
			// Verify scrolling behavior
			if tc.wantScroll && model.advancedViewport.YOffset == 0 {
				t.Errorf("Expected viewport to scroll for field %v, but it remained at the top", tc.focusField)
			}
			
			if !tc.wantScroll && model.advancedViewport.YOffset > 0 {
				t.Errorf("Expected viewport to remain at the top for field %v, but it scrolled to %d", 
					tc.focusField, model.advancedViewport.YOffset)
			}
			
			// Make sure content is not empty
			if content == "" {
				t.Error("Rendered content is empty")
			}
		})
	}
}