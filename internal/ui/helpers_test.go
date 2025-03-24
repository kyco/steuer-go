package ui

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