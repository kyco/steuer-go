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
		{"PrimaryColor", PrimaryColor, "#00AA00"},
		{"SecondaryColor", SecondaryColor, "#555555"},
		{"AccentColor", AccentColor, "#00FFFF"},
		{"SuccessColor", SuccessColor, "#00FF00"},
		{"DangerColor", DangerColor, "#FF0000"},
		{"WarningColor", WarningColor, "#FF8800"},
		{"NeutralColor", NeutralColor, "#AAAAAA"},
		{"BgColor", BgColor, "#000000"},
		{"FgColor", FgColor, "#FFFFFF"},
	}

	for _, tc := range colorTests {
		if string(tc.actual) != tc.expected {
			t.Errorf("Expected %s to be %q, got %q", tc.name, tc.expected, tc.actual)
		}
	}
}

func TestBorders(t *testing.T) {
	// Verify the minimal border definition
	if MinimalBorder.Top != "═" {
		t.Errorf("Expected MinimalBorder.Top to be '═', got %q", MinimalBorder.Top)
	}
	if MinimalBorder.Bottom != "═" {
		t.Errorf("Expected MinimalBorder.Bottom to be '═', got %q", MinimalBorder.Bottom)
	}
	if MinimalBorder.Left != "║" {
		t.Errorf("Expected MinimalBorder.Left to be '║', got %q", MinimalBorder.Left)
	}
	if MinimalBorder.Right != "║" {
		t.Errorf("Expected MinimalBorder.Right to be '║', got %q", MinimalBorder.Right)
	}
	if MinimalBorder.TopLeft != "╔" {
		t.Errorf("Expected MinimalBorder.TopLeft to be '╔', got %q", MinimalBorder.TopLeft)
	}
	if MinimalBorder.TopRight != "╗" {
		t.Errorf("Expected MinimalBorder.TopRight to be '╗', got %q", MinimalBorder.TopRight)
	}
	if MinimalBorder.BottomLeft != "╚" {
		t.Errorf("Expected MinimalBorder.BottomLeft to be '╚', got %q", MinimalBorder.BottomLeft)
	}
	if MinimalBorder.BottomRight != "╝" {
		t.Errorf("Expected MinimalBorder.BottomRight to be '╝', got %q", MinimalBorder.BottomRight)
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
	if SubtitleStyle.GetForeground() != SecondaryColor {
		t.Errorf("Expected SubtitleStyle foreground to be %s, got %s", SecondaryColor, SubtitleStyle.GetForeground())
	}
	if !SubtitleStyle.GetBold() {
		t.Error("Expected SubtitleStyle to be bold")
	}
	
	// Test ButtonStyle properties
	if ButtonStyle.GetForeground() != BgColor {
		t.Errorf("Expected ButtonStyle foreground to be %s, got %s", BgColor, ButtonStyle.GetForeground())
	}
	if ButtonStyle.GetBackground() != SecondaryColor {
		t.Errorf("Expected ButtonStyle background to be %s, got %s", SecondaryColor, ButtonStyle.GetBackground())
	}
	
	// Test SelectedButtonStyle properties
	if SelectedButtonStyle.GetBackground() != PrimaryColor {
		t.Errorf("Expected SelectedButtonStyle background to be %s, got %s", PrimaryColor, SelectedButtonStyle.GetBackground())
	}
	
	// Test HelpStyle properties
	if HelpStyle.GetForeground() != AccentColor {
		t.Errorf("Expected HelpStyle foreground to be %s, got %s", AccentColor, HelpStyle.GetForeground())
	}
	if !HelpStyle.GetItalic() {
		t.Error("Expected HelpStyle to be italic")
	}
	
	// Test HeaderStyle properties
	if HeaderStyle.GetForeground() != BgColor {
		t.Errorf("Expected HeaderStyle foreground to be %s, got %s", BgColor, HeaderStyle.GetForeground())
	}
	if HeaderStyle.GetBackground() != PrimaryColor {
		t.Errorf("Expected HeaderStyle background to be %s, got %s", PrimaryColor, HeaderStyle.GetBackground())
	}
	if HeaderStyle.GetAlign() != lipgloss.Center {
		t.Errorf("Expected HeaderStyle alignment to be Center, got %v", HeaderStyle.GetAlign())
	}
}