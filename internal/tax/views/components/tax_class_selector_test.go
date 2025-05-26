package components

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewTaxClassSelector(t *testing.T) {
	tcs := NewTaxClassSelector()

	if len(tcs.Options) != 6 {
		t.Errorf("Expected 6 tax class options, got %d", len(tcs.Options))
	}

	if tcs.Selected != 1 {
		t.Errorf("Expected default selected class 1, got %d", tcs.Selected)
	}

	if tcs.Focused {
		t.Error("Expected Focused to be false by default")
	}

	if tcs.ShowDetails {
		t.Error("Expected ShowDetails to be false by default")
	}

	expectedClasses := []int{1, 2, 3, 4, 5, 6}
	for i, option := range tcs.Options {
		if option.Class != expectedClasses[i] {
			t.Errorf("Expected class %d at index %d, got %d", expectedClasses[i], i, option.Class)
		}

		if option.Name == "" {
			t.Errorf("Expected non-empty name for class %d", option.Class)
		}

		if option.Description == "" {
			t.Errorf("Expected non-empty description for class %d", option.Class)
		}

		if option.Icon == "" {
			t.Errorf("Expected non-empty icon for class %d", option.Class)
		}

		if option.Details == "" {
			t.Errorf("Expected non-empty details for class %d", option.Class)
		}

		if option.CommonUse == "" {
			t.Errorf("Expected non-empty common use for class %d", option.Class)
		}
	}
}

func TestTaxClassSelectorUpdate(t *testing.T) {
	tests := []struct {
		name            string
		initialClass    int
		keyMsg          tea.KeyMsg
		expectedClass   int
		expectedDetails bool
	}{
		{
			name:          "down arrow from class 1",
			initialClass:  1,
			keyMsg:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("down")},
			expectedClass: 2,
		},
		{
			name:          "up arrow from class 1 wraps to 6",
			initialClass:  1,
			keyMsg:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("up")},
			expectedClass: 6,
		},
		{
			name:          "down arrow from class 6 wraps to 1",
			initialClass:  6,
			keyMsg:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("down")},
			expectedClass: 1,
		},
		{
			name:          "up arrow from class 6",
			initialClass:  6,
			keyMsg:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("up")},
			expectedClass: 5,
		},
		{
			name:          "j key moves down",
			initialClass:  3,
			keyMsg:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")},
			expectedClass: 4,
		},
		{
			name:          "k key moves up",
			initialClass:  3,
			keyMsg:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")},
			expectedClass: 2,
		},
		{
			name:            "h key toggles details",
			initialClass:    3,
			keyMsg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")},
			expectedClass:   3,
			expectedDetails: true,
		},
		{
			name:            "? key toggles details",
			initialClass:    3,
			keyMsg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")},
			expectedClass:   3,
			expectedDetails: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcs := NewTaxClassSelector()
			tcs.Selected = tt.initialClass

			updatedTcs, cmd := tcs.Update(tt.keyMsg)

			if cmd != nil {
				t.Error("Expected no command from Update")
			}

			if updatedTcs.Selected != tt.expectedClass {
				t.Errorf("Expected selected class %d, got %d", tt.expectedClass, updatedTcs.Selected)
			}

			if updatedTcs.ShowDetails != tt.expectedDetails {
				t.Errorf("Expected ShowDetails %v, got %v", tt.expectedDetails, updatedTcs.ShowDetails)
			}
		})
	}
}

func TestTaxClassSelectorView(t *testing.T) {
	tcs := NewTaxClassSelector()

	view := tcs.View()
	if view == "" {
		t.Error("View should not be empty")
	}

	if !strings.Contains(view, "Tax Class") {
		t.Error("View should contain title")
	}

	if !strings.Contains(view, "German tax classification") {
		t.Error("View should contain help text")
	}

	for _, option := range tcs.Options {
		if !strings.Contains(view, option.Name) {
			t.Errorf("View should contain option name %q", option.Name)
		}

		if !strings.Contains(view, option.Description) {
			t.Errorf("View should contain option description %q", option.Description)
		}

		if !strings.Contains(view, option.Icon) {
			t.Errorf("View should contain option icon %q", option.Icon)
		}
	}
}

func TestTaxClassSelectorViewWithSelection(t *testing.T) {
	tcs := NewTaxClassSelector()
	tcs.Selected = 3

	view := tcs.View()
	if !strings.Contains(view, "▶") {
		t.Error("View should contain selection indicator")
	}
}

func TestTaxClassSelectorViewWithDetails(t *testing.T) {
	tcs := NewTaxClassSelector()
	tcs.Selected = 2
	tcs.ShowDetails = true

	view := tcs.View()

	selectedOption := tcs.GetSelected()
	if !strings.Contains(view, selectedOption.Details) {
		t.Error("View should contain details when ShowDetails is true")
	}

	if !strings.Contains(view, selectedOption.CommonUse) {
		t.Error("View should contain common use when ShowDetails is true")
	}
}

func TestTaxClassSelectorViewWithFocus(t *testing.T) {
	tcs := NewTaxClassSelector()
	tcs.Focus()

	view := tcs.View()
	if !strings.Contains(view, "Press 'h' for more details") {
		t.Error("View should contain help text when focused")
	}

	if !strings.Contains(view, "↑/↓ to navigate") {
		t.Error("View should contain navigation help when focused")
	}
}

func TestTaxClassSelectorGetSelected(t *testing.T) {
	tcs := NewTaxClassSelector()

	for class := 1; class <= 6; class++ {
		tcs.Selected = class
		selected := tcs.GetSelected()

		if selected.Class != class {
			t.Errorf("Expected selected class %d, got %d", class, selected.Class)
		}

		if selected.Name == "" {
			t.Errorf("Expected non-empty name for selected class %d", class)
		}
	}
}

func TestTaxClassSelectorGetSelectedInvalid(t *testing.T) {
	tcs := NewTaxClassSelector()
	tcs.Selected = 999 // Invalid class

	selected := tcs.GetSelected()
	if selected.Class != 1 {
		t.Errorf("Expected fallback to class 1 for invalid selection, got %d", selected.Class)
	}
}

func TestTaxClassSelectorSetSelected(t *testing.T) {
	tcs := NewTaxClassSelector()

	validClasses := []int{1, 2, 3, 4, 5, 6}
	for _, class := range validClasses {
		tcs.SetSelected(class)
		if tcs.Selected != class {
			t.Errorf("Expected selected class %d, got %d", class, tcs.Selected)
		}
	}

	invalidClasses := []int{0, 7, -1, 100}
	for _, class := range invalidClasses {
		originalSelected := tcs.Selected
		tcs.SetSelected(class)
		if tcs.Selected != originalSelected {
			t.Errorf("Expected selected class to remain %d for invalid class %d, got %d", originalSelected, class, tcs.Selected)
		}
	}
}

func TestTaxClassSelectorFocusBlur(t *testing.T) {
	tcs := NewTaxClassSelector()

	if tcs.Focused {
		t.Error("Expected Focused to be false initially")
	}

	tcs.Focus()
	if !tcs.Focused {
		t.Error("Expected Focused to be true after Focus()")
	}

	tcs.Blur()
	if tcs.Focused {
		t.Error("Expected Focused to be false after Blur()")
	}
}

func TestTaxClassSelectorToggleDetails(t *testing.T) {
	tcs := NewTaxClassSelector()

	if tcs.ShowDetails {
		t.Error("Expected ShowDetails to be false initially")
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")}
	updatedTcs, _ := tcs.Update(keyMsg)

	if !updatedTcs.ShowDetails {
		t.Error("Expected ShowDetails to be true after pressing 'h'")
	}

	updatedTcs2, _ := updatedTcs.Update(keyMsg)

	if updatedTcs2.ShowDetails {
		t.Error("Expected ShowDetails to be false after pressing 'h' again")
	}
}

func TestTaxClassInfo(t *testing.T) {
	info := TaxClassInfo{
		Class:       1,
		Name:        "Test Class",
		Description: "Test Description",
		Icon:        "🧪",
		Details:     "Test Details",
		CommonUse:   "Test Common Use",
	}

	if info.Class != 1 {
		t.Errorf("Expected Class 1, got %d", info.Class)
	}

	if info.Name != "Test Class" {
		t.Errorf("Expected Name 'Test Class', got %q", info.Name)
	}

	if info.Description != "Test Description" {
		t.Errorf("Expected Description 'Test Description', got %q", info.Description)
	}

	if info.Icon != "🧪" {
		t.Errorf("Expected Icon '🧪', got %q", info.Icon)
	}

	if info.Details != "Test Details" {
		t.Errorf("Expected Details 'Test Details', got %q", info.Details)
	}

	if info.CommonUse != "Test Common Use" {
		t.Errorf("Expected CommonUse 'Test Common Use', got %q", info.CommonUse)
	}
}

func TestTaxClassSelectorNavigationBoundaries(t *testing.T) {
	tcs := NewTaxClassSelector()

	tcs.Selected = 1
	upMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("up")}
	updatedTcs, _ := tcs.Update(upMsg)
	if updatedTcs.Selected != 6 {
		t.Errorf("Expected wrapping from 1 to 6, got %d", updatedTcs.Selected)
	}

	tcs.Selected = 6
	downMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("down")}
	updatedTcs, _ = tcs.Update(downMsg)
	if updatedTcs.Selected != 1 {
		t.Errorf("Expected wrapping from 6 to 1, got %d", updatedTcs.Selected)
	}
}
