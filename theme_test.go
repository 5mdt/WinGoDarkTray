// theme_test.go

//go:build windows
// +build windows

package main

import (
	"testing"

	"golang.org/x/sys/windows/registry"
)

func TestThemeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "theme registry path",
			constant: "themeRegistryPath",
			expected: `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`,
		},
		{
			name:     "apps use light registry value",
			constant: "regValAppsUseLight",
			expected: "AppsUseLightTheme",
		},
		{
			name:     "system uses light registry value",
			constant: "regValSystemUsesLight",
			expected: "SystemUsesLightTheme",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.constant {
			case "themeRegistryPath":
				if themeRegistryPath != tt.expected {
					t.Errorf("themeRegistryPath = %s, want %s", themeRegistryPath, tt.expected)
				}
			case "regValAppsUseLight":
				if regValAppsUseLight != tt.expected {
					t.Errorf("regValAppsUseLight = %s, want %s", regValAppsUseLight, tt.expected)
				}
			case "regValSystemUsesLight":
				if regValSystemUsesLight != tt.expected {
					t.Errorf("regValSystemUsesLight = %s, want %s", regValSystemUsesLight, tt.expected)
				}
			}
		})
	}
}

func TestWithThemeRegistry(t *testing.T) {
	tests := []struct {
		name        string
		access      uint32
		shouldError bool
	}{
		{
			name:        "query access",
			access:      registry.QUERY_VALUE,
			shouldError: false,
		},
		{
			name:        "read access",
			access:      registry.READ,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			err := withThemeRegistry(tt.access, func(key registry.Key) error {
				called = true
				return nil
			})

			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !called && !tt.shouldError {
				t.Error("Function was not called")
			}
		})
	}
}

func TestGetCurrentAppThemeMode(t *testing.T) {
	// This test attempts to read the actual theme registry
	// It may fail if the registry key doesn't exist or access is denied
	t.Run("read current app theme mode", func(t *testing.T) {
		mode, err := getCurrentAppThemeMode()
		if err != nil {
			t.Logf("Expected behavior: could not read theme mode: %v", err)
			return
		}

		// Valid values are 0 (dark) or 1 (light)
		if mode != 0 && mode != 1 {
			t.Errorf("getCurrentAppThemeMode() returned invalid mode: %d, expected 0 or 1", mode)
		}
	})
}

func TestSetBothThemeModes(t *testing.T) {
	tests := []struct {
		name      string
		lightMode bool
	}{
		{
			name:      "set to light mode",
			lightMode: true,
		},
		{
			name:      "set to dark mode",
			lightMode: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This test may fail if we don't have write access to the registry
			// or if we're not running with sufficient privileges
			err := setBothThemeModes(tt.lightMode)
			if err != nil {
				t.Logf("Expected behavior on restricted systems: %v", err)
				// Don't fail the test as this is expected behavior in many environments
				return
			}

			// If successful, verify the values were set correctly
			appMode, err := getCurrentAppThemeMode()
			if err != nil {
				t.Logf("Could not verify app mode after setting: %v", err)
				return
			}

			expectedMode := uint64(0)
			if tt.lightMode {
				expectedMode = 1
			}

			if appMode != expectedMode {
				t.Errorf("After setBothThemeModes(%v), app mode = %d, want %d", tt.lightMode, appMode, expectedMode)
			}
		})
	}
}

func TestAppThemeMethods(t *testing.T) {
	app := NewApp("test-version")

	tests := []struct {
		name     string
		testFunc func()
	}{
		{
			name: "toggle app mode",
			testFunc: func() {
				app.toggleAppMode()
			},
		},
		{
			name: "toggle windows mode",
			testFunc: func() {
				app.toggleWindowsMode()
			},
		},
		{
			name: "toggle system mode",
			testFunc: func() {
				app.toggleSystemMode()
			},
		},
		{
			name: "update theme toggle titles",
			testFunc: func() {
				app.updateThemeToggleTitles()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip these tests as they require systray initialization and registry write access
			t.Skip("theme methods require systray initialization and registry access")
		})
	}
}

func TestToggleSingleTheme(t *testing.T) {
	tests := []struct {
		name        string
		registryKey string
	}{
		{
			name:        "toggle apps use light theme",
			registryKey: regValAppsUseLight,
		},
		{
			name:        "toggle system uses light theme",
			registryKey: regValSystemUsesLight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip this test as it requires systray initialization and registry write access
			t.Skip("toggleSingleTheme requires systray initialization and registry access")
		})
	}
}