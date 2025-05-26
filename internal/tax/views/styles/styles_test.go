package styles

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestColors(t *testing.T) {
	// Verify color definitions
	colorTests := []struct {
		name     string
		actual   lipgloss.Color
		expected string
	}{
		{"PrimaryColor", PrimaryColor, "#7b9cd9"},
		{"SecondaryColor", SecondaryColor, "#323232"},
		{"AccentColor", AccentColor, "#db7093"},
		{"SuccessColor", SuccessColor, "#76b876"},
		{"DangerColor", DangerColor, "#e06c75"},
		{"WarningColor", WarningColor, "#e5c07b"},
		{"NeutralColor", NeutralColor, "#abb2bf"},
		{"BgColor", BgColor, "#282c34"},
		{"FgColor", FgColor, "#f8f8f2"},
	}

	for _, tc := range colorTests {
		if string(tc.actual) != tc.expected {
			t.Errorf("Expected %s to be %q, got %q", tc.name, tc.expected, tc.actual)
		}
	}
}

func TestBorders(t *testing.T) {
	// Verify the SimpleBorder definition
	if SimpleBorder.Top != "─" {
		t.Errorf("Expected SimpleBorder.Top to be '─', got %q", SimpleBorder.Top)
	}
	if SimpleBorder.Bottom != "─" {
		t.Errorf("Expected SimpleBorder.Bottom to be '─', got %q", SimpleBorder.Bottom)
	}
	if SimpleBorder.Left != "│" {
		t.Errorf("Expected SimpleBorder.Left to be '│', got %q", SimpleBorder.Left)
	}
	if SimpleBorder.Right != "│" {
		t.Errorf("Expected SimpleBorder.Right to be '│', got %q", SimpleBorder.Right)
	}
	if SimpleBorder.TopLeft != "┌" {
		t.Errorf("Expected SimpleBorder.TopLeft to be '┌', got %q", SimpleBorder.TopLeft)
	}
	if SimpleBorder.TopRight != "┐" {
		t.Errorf("Expected SimpleBorder.TopRight to be '┐', got %q", SimpleBorder.TopRight)
	}
	if SimpleBorder.BottomLeft != "└" {
		t.Errorf("Expected SimpleBorder.BottomLeft to be '└', got %q", SimpleBorder.BottomLeft)
	}
	if SimpleBorder.BottomRight != "┘" {
		t.Errorf("Expected SimpleBorder.BottomRight to be '┘', got %q", SimpleBorder.BottomRight)
	}
}

func TestStyles(t *testing.T) {
	// Test BaseStyle properties
	if BaseStyle.GetForeground() != FgColor {
		t.Errorf("Expected BaseStyle foreground to be %s, got %s", FgColor, BaseStyle.GetForeground())
	}
	if BaseStyle.GetBackground() != BgColor {
		t.Errorf("Expected BaseStyle background to be %s, got %s", BgColor, BaseStyle.GetBackground())
	}

	// Test TitleStyle properties
	if TitleStyle.GetForeground() != PrimaryColor {
		t.Errorf("Expected TitleStyle foreground to be %s, got %s", PrimaryColor, TitleStyle.GetForeground())
	}
	if !TitleStyle.GetBold() {
		t.Error("Expected TitleStyle to be bold")
	}

	// Test SubtitleStyle properties
	if SubtitleStyle.GetForeground() != AccentColor {
		t.Errorf("Expected SubtitleStyle foreground to be %s, got %s", AccentColor, SubtitleStyle.GetForeground())
	}

	// Test ButtonStyle properties
	if ButtonStyle.String() == "" {
		t.Errorf("Expected ButtonStyle to be properly defined")
	}

	// Test SelectedButtonStyle properties
	if !SelectedButtonStyle.GetBold() {
		t.Error("Expected SelectedButtonStyle to be bold")
	}

	// Test HelpStyle properties
	if HelpStyle.GetForeground() != NeutralColor {
		t.Errorf("Expected HelpStyle foreground to be %s, got %s", NeutralColor, HelpStyle.GetForeground())
	}
	if !HelpStyle.GetItalic() {
		t.Error("Expected HelpStyle to be italic")
	}
}
