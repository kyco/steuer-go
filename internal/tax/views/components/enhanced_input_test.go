package components

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewEnhancedInput(t *testing.T) {
	label := "Test Label"
	description := "Test Description"
	placeholder := "Test Placeholder"

	ei := NewEnhancedInput(label, description, placeholder)

	if ei.Label != label {
		t.Errorf("Expected label %q, got %q", label, ei.Label)
	}

	if ei.Description != description {
		t.Errorf("Expected description %q, got %q", description, ei.Description)
	}

	if ei.Placeholder != placeholder {
		t.Errorf("Expected placeholder %q, got %q", placeholder, ei.Placeholder)
	}

	if !ei.IsValid {
		t.Error("Expected IsValid to be true by default")
	}

	if ei.Required {
		t.Error("Expected Required to be false by default")
	}

	if ei.ErrorMessage != "" {
		t.Errorf("Expected empty ErrorMessage, got %q", ei.ErrorMessage)
	}

	if len(ei.ValidationRules) != 0 {
		t.Errorf("Expected empty ValidationRules, got %d rules", len(ei.ValidationRules))
	}
}

func TestEnhancedInputAddValidation(t *testing.T) {
	ei := NewEnhancedInput("Test", "Test", "Test")

	rule1 := func(value string) error {
		if value == "invalid" {
			return fmt.Errorf("invalid value")
		}
		return nil
	}

	rule2 := func(value string) error {
		if len(value) > 10 {
			return fmt.Errorf("too long")
		}
		return nil
	}

	ei.AddValidation(rule1)
	if len(ei.ValidationRules) != 1 {
		t.Errorf("Expected 1 validation rule, got %d", len(ei.ValidationRules))
	}

	ei.AddValidation(rule2)
	if len(ei.ValidationRules) != 2 {
		t.Errorf("Expected 2 validation rules, got %d", len(ei.ValidationRules))
	}
}

func TestEnhancedInputValidate(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		required    bool
		rules       []ValidationRule
		expectValid bool
		expectError string
	}{
		{
			name:        "empty value not required",
			value:       "",
			required:    false,
			rules:       nil,
			expectValid: true,
			expectError: "",
		},
		{
			name:        "empty value required",
			value:       "",
			required:    true,
			rules:       nil,
			expectValid: false,
			expectError: "Test is required",
		},
		{
			name:     "valid value with rules",
			value:    "valid",
			required: false,
			rules: []ValidationRule{
				func(value string) error {
					if value == "invalid" {
						return fmt.Errorf("invalid value")
					}
					return nil
				},
			},
			expectValid: true,
			expectError: "",
		},
		{
			name:     "invalid value with rules",
			value:    "invalid",
			required: false,
			rules: []ValidationRule{
				func(value string) error {
					if value == "invalid" {
						return fmt.Errorf("invalid value")
					}
					return nil
				},
			},
			expectValid: false,
			expectError: "invalid value",
		},
		{
			name:     "multiple rules first fails",
			value:    "test",
			required: false,
			rules: []ValidationRule{
				func(value string) error {
					return fmt.Errorf("first rule fails")
				},
				func(value string) error {
					return fmt.Errorf("second rule fails")
				},
			},
			expectValid: false,
			expectError: "first rule fails",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ei := NewEnhancedInput("Test", "Test", "Test")
			ei.Required = tt.required
			ei.SetValue(tt.value)

			for _, rule := range tt.rules {
				ei.AddValidation(rule)
			}

			ei.Validate()

			if ei.IsValid != tt.expectValid {
				t.Errorf("Expected IsValid %v, got %v", tt.expectValid, ei.IsValid)
			}

			if ei.ErrorMessage != tt.expectError {
				t.Errorf("Expected ErrorMessage %q, got %q", tt.expectError, ei.ErrorMessage)
			}
		})
	}
}

func TestEnhancedInputUpdate(t *testing.T) {
	ei := NewEnhancedInput("Test", "Test", "Test")
	ei.AddValidation(func(value string) error {
		if value == "invalid" {
			return fmt.Errorf("invalid value")
		}
		return nil
	})

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("test")}
	updatedEi, cmd := ei.Update(keyMsg)

	// Command may or may not be returned depending on the underlying textinput implementation
	_ = cmd

	if !updatedEi.IsValid {
		t.Error("Expected input to be valid after update with valid value")
	}
}

func TestEnhancedInputView(t *testing.T) {
	ei := NewEnhancedInput("Test Label", "Test Description", "Test Placeholder")

	view := ei.View()
	if view == "" {
		t.Error("View should not be empty")
	}

	if !strings.Contains(view, "Test Label") {
		t.Error("View should contain label")
	}

	if !strings.Contains(view, "Test Description") {
		t.Error("View should contain description")
	}
}

func TestEnhancedInputViewWithError(t *testing.T) {
	ei := NewEnhancedInput("Test", "Test", "Test")
	ei.Required = true
	ei.SetValue("")
	ei.Validate()

	view := ei.View()
	if !strings.Contains(view, "⚠") {
		t.Error("View should contain error indicator when invalid")
	}

	if !strings.Contains(view, "required") {
		t.Error("View should contain error message")
	}
}

func TestEnhancedInputViewWithoutDescription(t *testing.T) {
	ei := NewEnhancedInput("Test Label", "", "Test Placeholder")

	view := ei.View()
	if !strings.Contains(view, "Test Label") {
		t.Error("View should contain label")
	}

	lines := strings.Split(view, "\n")
	if len(lines) < 2 {
		t.Error("View should have at least label and input lines")
	}
}

func TestValidateIncome(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
		errorText   string
	}{
		{
			name:        "empty value",
			value:       "",
			expectError: false,
		},
		{
			name:        "valid income",
			value:       "50000",
			expectError: false,
		},
		{
			name:        "valid decimal income",
			value:       "50000.50",
			expectError: false,
		},
		{
			name:        "zero income",
			value:       "0",
			expectError: false,
		},
		{
			name:        "negative income",
			value:       "-1000",
			expectError: true,
			errorText:   "cannot be negative",
		},
		{
			name:        "invalid number",
			value:       "not-a-number",
			expectError: true,
			errorText:   "valid number",
		},
		{
			name:        "unrealistically high income",
			value:       "20000000",
			expectError: true,
			errorText:   "unrealistically high",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIncome(tt.value)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorText) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorText, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateYear(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
		errorText   string
	}{
		{
			name:        "empty value",
			value:       "",
			expectError: false,
		},
		{
			name:        "valid year",
			value:       "2024",
			expectError: false,
		},
		{
			name:        "minimum valid year",
			value:       "2020",
			expectError: false,
		},
		{
			name:        "maximum valid year",
			value:       "2030",
			expectError: false,
		},
		{
			name:        "year too low",
			value:       "2019",
			expectError: true,
			errorText:   "between 2020 and 2030",
		},
		{
			name:        "year too high",
			value:       "2031",
			expectError: true,
			errorText:   "between 2020 and 2030",
		},
		{
			name:        "invalid year format",
			value:       "not-a-year",
			expectError: true,
			errorText:   "valid year",
		},
		{
			name:        "decimal year",
			value:       "2024.5",
			expectError: true,
			errorText:   "valid year",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateYear(tt.value)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorText) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorText, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestEnhancedInputFocusHandling(t *testing.T) {
	ei := NewEnhancedInput("Test", "Test", "Test")

	ei.Focus()
	if !ei.Focused() {
		t.Error("Expected input to be focused after Focus()")
	}

	ei.Blur()
	if ei.Focused() {
		t.Error("Expected input to be blurred after Blur()")
	}
}

func TestEnhancedInputValueHandling(t *testing.T) {
	ei := NewEnhancedInput("Test", "Test", "Test")

	testValue := "test value"
	ei.SetValue(testValue)

	if ei.Value() != testValue {
		t.Errorf("Expected value %q, got %q", testValue, ei.Value())
	}
}
