// icon_test.go

package main

import (
	"testing"
)

func TestIconConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "autorun registry key",
			constant: "autorunRegistryKey",
			expected: `Software\Microsoft\Windows\CurrentVersion\Run`,
		},
		{
			name:     "app name",
			constant: "appName",
			expected: "WinGoDarkTray",
		},
		{
			name:     "project link",
			constant: "projectLink",
			expected: "https://github.com/5mdt/WinGoDarkTray",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.constant {
			case "autorunRegistryKey":
				if autorunRegistryKey != tt.expected {
					t.Errorf("autorunRegistryKey = %s, want %s", autorunRegistryKey, tt.expected)
				}
			case "appName":
				if appName != tt.expected {
					t.Errorf("appName = %s, want %s", appName, tt.expected)
				}
			case "projectLink":
				if projectLink != tt.expected {
					t.Errorf("projectLink = %s, want %s", projectLink, tt.expected)
				}
			}
		})
	}
}

func TestIconEmbedded(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "icon data is embedded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the icon byte slice exists and is not empty
			if icon == nil {
				t.Error("icon byte slice is nil")
			}
			if len(icon) == 0 {
				t.Error("icon byte slice is empty")
			}

			// Basic validation that it looks like an ICO file
			// ICO files start with 0x00 0x00 0x01 0x00
			if len(icon) >= 4 {
				if icon[0] != 0x00 || icon[1] != 0x00 || icon[2] != 0x01 || icon[3] != 0x00 {
					t.Logf("Warning: icon data may not be a valid ICO file (header: %02x %02x %02x %02x)",
						icon[0], icon[1], icon[2], icon[3])
				}
			} else {
				t.Error("icon data is too short to be a valid ICO file")
			}
		})
	}
}
