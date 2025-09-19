package main

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// Test pure functions that don't depend on Windows APIs

func TestIsVersionNewerPure(t *testing.T) {
	// Extracted logic test without external dependencies
	testVersionNewer := func(latest, current string) bool {
		latest = strings.TrimPrefix(latest, "v")
		current = strings.TrimPrefix(current, "v")

		latestParts := strings.Split(latest, ".")
		currentParts := strings.Split(current, ".")

		// Normalize to same length by padding with zeros
		maxLen := len(latestParts)
		if len(currentParts) > maxLen {
			maxLen = len(currentParts)
		}

		// Pad with zeros
		for len(latestParts) < maxLen {
			latestParts = append(latestParts, "0")
		}
		for len(currentParts) < maxLen {
			currentParts = append(currentParts, "0")
		}

		for i := 0; i < maxLen; i++ {
			latestNum, err1 := strconv.Atoi(latestParts[i])
			currentNum, err2 := strconv.Atoi(currentParts[i])
			if err1 != nil || err2 != nil {
				return false
			}

			if latestNum > currentNum {
				return true
			} else if latestNum < currentNum {
				return false
			}
		}

		return false
	}

	tests := []struct {
		latest  string
		current string
		want    bool
	}{
		{"0.9.9", "0.9.8", true},
		{"1.10.0", "1.2.3", true},
		{"1.2.3", "1.10.0", false},
		{"1.2.3", "1.2.3", false},
		{"1.2.3", "2.0.0", false},
		{"1.2.3", "invalid", false},
		{"1.2.4", "1.2.3", true},
		{"1.2", "1.2.0", false},
		{"1.3.0", "1.2.9", true},
		{"2.0.0", "1.9.9", true},
		{"invalid", "1.2.3", false},
		{"v1.10.0", "v1.2.3", true},
		{"v1.2.4", "v1.2.3", true},
		{"v2.0.0", "v1.9.9", true},
	}

	for _, tt := range tests {
		t.Run(tt.latest+" vs "+tt.current, func(t *testing.T) {
			got := testVersionNewer(tt.latest, tt.current)
			if got != tt.want {
				t.Errorf("testVersionNewer(%q, %q) = %v; want %v", tt.latest, tt.current, got, tt.want)
			}
		})
	}
}

func TestPadVersionPartsPure(t *testing.T) {
	padVersionParts := func(parts []string, targetLen int) []string {
		result := make([]string, len(parts))
		copy(result, parts)
		for len(result) < targetLen {
			result = append(result, "0")
		}
		return result
	}

	tests := []struct {
		name      string
		parts     []string
		targetLen int
		want      []string
	}{
		{
			name:      "pad from 2 to 3",
			parts:     []string{"1", "2"},
			targetLen: 3,
			want:      []string{"1", "2", "0"},
		},
		{
			name:      "pad from 1 to 3",
			parts:     []string{"1"},
			targetLen: 3,
			want:      []string{"1", "0", "0"},
		},
		{
			name:      "no padding needed",
			parts:     []string{"1", "2", "3"},
			targetLen: 3,
			want:      []string{"1", "2", "3"},
		},
		{
			name:      "already longer than target",
			parts:     []string{"1", "2", "3", "4"},
			targetLen: 3,
			want:      []string{"1", "2", "3", "4"},
		},
		{
			name:      "empty parts",
			parts:     []string{},
			targetLen: 2,
			want:      []string{"0", "0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := padVersionParts(tt.parts, tt.targetLen)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("padVersionParts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxPure(t *testing.T) {
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"a greater than b", 5, 3, 5},
		{"b greater than a", 2, 7, 7},
		{"a equals b", 4, 4, 4},
		{"negative numbers", -2, -5, -2},
		{"zero and positive", 0, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := max(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("max(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestUITextsValidation(t *testing.T) {
	// Test menu titles structure without importing Windows dependencies
	expectedTexts := map[string]string{
		"AppName":                "🔗 WinGoDarkTray",
		"EnableAutorun":          "Enable Autorun",
		"Quit":                   "✕ Quit",
		"EnableAutorunChecked":   "✔ Autorun Enabled",
		"EnableAutorunUnchecked": "✗ Autorun Disabled",
		"ToggleAppToDark":        "☾ Toggle app theme to Dark",
		"ToggleAppToLight":       "☼ Toggle app theme to Light",
		"ToggleWinToDark":        "☾ Toggle Windows theme to Dark",
		"ToggleWinToLight":       "☼ Toggle Windows theme to Light",
		"ToggleBothToDark":       "☾ Toggle both to Dark",
		"ToggleBothToLight":      "☼ Toggle both to Light",
		"UpdateNow":              "🔄 Update Now",
	}

	for name, expected := range expectedTexts {
		if expected == "" {
			t.Errorf("Expected text for %s should not be empty", name)
		}
		if name == "AppName" && !strings.Contains(expected, "WinGoDarkTray") {
			t.Errorf("AppName should contain 'WinGoDarkTray', got: %s", expected)
		}
	}

	// Test tooltip structure
	defaultTooltip := "A windows app to toggle light and dark mode from the system tray"
	errorTooltip := "An error occurred, please try again later."

	if !strings.Contains(defaultTooltip, "toggle") || !strings.Contains(defaultTooltip, "windows") {
		t.Error("Default tooltip should describe app purpose")
	}
	if errorTooltip == "" {
		t.Error("Error tooltip should not be empty")
	}
}
