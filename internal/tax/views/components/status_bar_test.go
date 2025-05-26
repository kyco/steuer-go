package components

import (
	"strings"
	"testing"
)

func TestNewStatusBar(t *testing.T) {
	width := 100
	sb := NewStatusBar(width)

	if sb.Width != width {
		t.Errorf("Expected width %d, got %d", width, sb.Width)
	}

	if sb.CurrentScreen != "" {
		t.Errorf("Expected empty CurrentScreen, got %q", sb.CurrentScreen)
	}

	if sb.CurrentMode != "" {
		t.Errorf("Expected empty CurrentMode, got %q", sb.CurrentMode)
	}

	if len(sb.AvailableKeys) != 0 {
		t.Errorf("Expected empty AvailableKeys, got %d items", len(sb.AvailableKeys))
	}
}

func TestStatusBarSetters(t *testing.T) {
	sb := NewStatusBar(100)

	screen := "Main Screen"
	sb.SetScreen(screen)
	if sb.CurrentScreen != screen {
		t.Errorf("Expected CurrentScreen %q, got %q", screen, sb.CurrentScreen)
	}

	mode := "Local Mode"
	sb.SetMode(mode)
	if sb.CurrentMode != mode {
		t.Errorf("Expected CurrentMode %q, got %q", mode, sb.CurrentMode)
	}

	keys := []KeyHint{
		{Key: "Tab", Description: "Next", Important: false},
		{Key: "Enter", Description: "Select", Important: true},
	}
	sb.SetKeys(keys)
	if len(sb.AvailableKeys) != len(keys) {
		t.Errorf("Expected %d keys, got %d", len(keys), len(sb.AvailableKeys))
	}

	for i, key := range keys {
		if sb.AvailableKeys[i].Key != key.Key {
			t.Errorf("Expected key %q, got %q", key.Key, sb.AvailableKeys[i].Key)
		}
		if sb.AvailableKeys[i].Description != key.Description {
			t.Errorf("Expected description %q, got %q", key.Description, sb.AvailableKeys[i].Description)
		}
		if sb.AvailableKeys[i].Important != key.Important {
			t.Errorf("Expected important %v, got %v", key.Important, sb.AvailableKeys[i].Important)
		}
	}
}

func TestStatusBarView(t *testing.T) {
	sb := NewStatusBar(100)
	sb.SetScreen("Test Screen")
	sb.SetMode("Test Mode")
	sb.SetKeys([]KeyHint{
		{Key: "Tab", Description: "Next", Important: false},
	})

	view := sb.View()
	if view == "" {
		t.Error("View should not be empty")
	}

	if !strings.Contains(view, "Test Screen") {
		t.Error("View should contain screen name")
	}
}

func TestStatusBarRenderLeftSection(t *testing.T) {
	sb := NewStatusBar(100)

	sb.SetScreen("Main")
	left := sb.renderLeftSection()
	if !strings.Contains(left, "Main") {
		t.Error("Left section should contain screen name")
	}

	sb.SetMode("Local")
	left = sb.renderLeftSection()
	if !strings.Contains(left, "Main") {
		t.Error("Left section should contain screen name")
	}
	if !strings.Contains(left, "Local") {
		t.Error("Left section should contain mode when set")
	}
}

func TestStatusBarRenderCenterSection(t *testing.T) {
	sb := NewStatusBar(100)

	center := sb.renderCenterSection(50)
	if len(center) == 0 {
		t.Error("Center section should not be empty even with no keys")
	}

	sb.SetKeys([]KeyHint{
		{Key: "Tab", Description: "Next", Important: false},
		{Key: "Enter", Description: "Select", Important: true},
	})

	center = sb.renderCenterSection(50)
	if !strings.Contains(center, "Tab") {
		t.Error("Center section should contain key hints")
	}
	if !strings.Contains(center, "Next") {
		t.Error("Center section should contain key descriptions")
	}
}

func TestStatusBarRenderRightSection(t *testing.T) {
	sb := NewStatusBar(100)
	right := sb.renderRightSection()

	if !strings.Contains(right, "SteuerGo") {
		t.Error("Right section should contain app name")
	}
}

func TestGetMainScreenKeys(t *testing.T) {
	keys := GetMainScreenKeys()

	if len(keys) == 0 {
		t.Error("Main screen keys should not be empty")
	}

	expectedKeys := []string{"Tab", "↑/↓", "Enter", "L", "Q"}
	for _, expectedKey := range expectedKeys {
		found := false
		for _, key := range keys {
			if strings.Contains(key.Key, expectedKey) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find key containing %q", expectedKey)
		}
	}

	importantFound := false
	for _, key := range keys {
		if key.Important {
			importantFound = true
			break
		}
	}
	if !importantFound {
		t.Error("Expected at least one important key in main screen")
	}
}

func TestGetResultsScreenKeys(t *testing.T) {
	keys := GetResultsScreenKeys()

	if len(keys) == 0 {
		t.Error("Results screen keys should not be empty")
	}

	expectedKeys := []string{"←/→", "C", "B", "↑/↓", "Esc"}
	for _, expectedKey := range expectedKeys {
		found := false
		for _, key := range keys {
			if strings.Contains(key.Key, expectedKey) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find key containing %q", expectedKey)
		}
	}
}

func TestGetAdvancedScreenKeys(t *testing.T) {
	keys := GetAdvancedScreenKeys()

	if len(keys) == 0 {
		t.Error("Advanced screen keys should not be empty")
	}

	expectedKeys := []string{"Tab", "Enter", "Esc", "↑/↓"}
	for _, expectedKey := range expectedKeys {
		found := false
		for _, key := range keys {
			if strings.Contains(key.Key, expectedKey) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find key containing %q", expectedKey)
		}
	}
}

func TestGetComparisonScreenKeys(t *testing.T) {
	keys := GetComparisonScreenKeys()

	if len(keys) == 0 {
		t.Error("Comparison screen keys should not be empty")
	}

	expectedKeys := []string{"↑/↓", "Enter", "B", "Esc"}
	for _, expectedKey := range expectedKeys {
		found := false
		for _, key := range keys {
			if strings.Contains(key.Key, expectedKey) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find key containing %q", expectedKey)
		}
	}
}

func TestGetComparisonBreakdownKeys(t *testing.T) {
	keys := GetComparisonBreakdownKeys()

	if len(keys) == 0 {
		t.Error("Comparison breakdown keys should not be empty")
	}

	expectedKeys := []string{"Enter", "B", "Esc"}
	for _, expectedKey := range expectedKeys {
		found := false
		for _, key := range keys {
			if strings.Contains(key.Key, expectedKey) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find key containing %q", expectedKey)
		}
	}
}

func TestKeyHintStruct(t *testing.T) {
	hint := KeyHint{
		Key:         "Test",
		Description: "Test Description",
		Important:   true,
	}

	if hint.Key != "Test" {
		t.Errorf("Expected Key %q, got %q", "Test", hint.Key)
	}
	if hint.Description != "Test Description" {
		t.Errorf("Expected Description %q, got %q", "Test Description", hint.Description)
	}
	if !hint.Important {
		t.Error("Expected Important to be true")
	}
}

func TestStatusBarWithEmptyKeys(t *testing.T) {
	sb := NewStatusBar(100)
	sb.SetKeys([]KeyHint{})

	center := sb.renderCenterSection(50)
	if len(center) == 0 {
		t.Error("Center section should handle empty keys gracefully")
	}
}

func TestStatusBarWithLongContent(t *testing.T) {
	sb := NewStatusBar(20) // Small width
	sb.SetScreen("Very Long Screen Name")
	sb.SetMode("Very Long Mode Name")
	sb.SetKeys([]KeyHint{
		{Key: "VeryLongKey", Description: "Very Long Description", Important: false},
	})

	view := sb.View()
	if view == "" {
		t.Error("View should handle long content gracefully")
	}
}
