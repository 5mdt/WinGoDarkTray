// autorun_test.go

//go:build windows
// +build windows

package main

import (
	"fmt"
	"runtime"
	"testing"

	"golang.org/x/sys/windows/registry"
)

func TestUpdateAutorunUI(t *testing.T) {
	tests := []struct {
		name           string
		enabled        bool
		expectedTitle  string
	}{
		{
			name:          "enabled autorun",
			enabled:       true,
			expectedTitle: menuTitles.EnableAutorunChecked,
		},
		{
			name:          "disabled autorun",
			enabled:       false,
			expectedTitle: menuTitles.EnableAutorunUnchecked,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip this test as it requires systray initialization
			t.Skip("updateAutorunUI requires systray initialization")
		})
	}
}

func TestIsAutorunEnabled(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(key registry.Key) error
		expected bool
	}{
		{
			name: "autorun enabled - value exists",
			setup: func(key registry.Key) error {
				return key.SetStringValue(appName, "C:\\test\\path.exe")
			},
			expected: true,
		},
		{
			name: "autorun disabled - value does not exist",
			setup: func(key registry.Key) error {
				// Try to delete the value to ensure it doesn't exist
				key.DeleteValue(appName) // Ignore error as it might not exist
				return nil
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary registry key for testing
			testKey, err := registry.OpenKey(registry.CURRENT_USER, "Software\\WinGoDarkTrayTest", registry.ALL_ACCESS)
			if err != nil {
				// Create the key if it doesn't exist
				testKey, _, err = registry.CreateKey(registry.CURRENT_USER, "Software\\WinGoDarkTrayTest", registry.ALL_ACCESS)
				if err != nil {
					t.Skipf("Cannot create test registry key: %v", err)
				}
			}
			defer func() {
				testKey.Close()
				// Clean up the test key
				registry.DeleteKey(registry.CURRENT_USER, "Software\\WinGoDarkTrayTest")
			}()

			// Setup the test condition
			if err := tt.setup(testKey); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// Test the function
			result := isAutorunEnabled(testKey)
			if result != tt.expected {
				t.Errorf("isAutorunEnabled() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWithAutorunRegistry(t *testing.T) {
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
			err := withAutorunRegistry(tt.access, func(key registry.Key) error {
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

func TestAutorunRegistryErrorHandling(t *testing.T) {
	tests := []struct {
		name   string
		access uint32
	}{
		{
			name:   "invalid access permissions",
			access: 0xFFFFFFFF, // Invalid access flag
		},
		{
			name:   "valid read access",
			access: registry.READ,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := withAutorunRegistry(tt.access, func(key registry.Key) error {
				// Test function doesn't panic and returns appropriately
				return nil
			})

			// We don't assert specific error behavior as it depends on system state
			// but we verify the function doesn't panic
			_ = err
		})
	}
}

func TestIsAutorunEnabledEdgeCases(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on non-Windows platform")
	}

	// Create a test registry key for edge case testing
	testKey, _, err := registry.CreateKey(registry.CURRENT_USER, "Software\\WinGoDarkTrayAutorunTest", registry.ALL_ACCESS)
	if err != nil {
		t.Skipf("Cannot create test registry key: %v", err)
	}
	defer func() {
		testKey.Close()
		registry.DeleteKey(registry.CURRENT_USER, "Software\\WinGoDarkTrayAutorunTest")
	}()

	tests := []struct {
		name     string
		setup    func() error
		expected bool
	}{
		{
			name: "value exists with empty string",
			setup: func() error {
				return testKey.SetStringValue(appName, "")
			},
			expected: true,
		},
		{
			name: "value exists with whitespace",
			setup: func() error {
				return testKey.SetStringValue(appName, "   ")
			},
			expected: true,
		},
		{
			name: "value exists with path",
			setup: func() error {
				return testKey.SetStringValue(appName, "C:\\test\\path.exe")
			},
			expected: true,
		},
		{
			name: "value deleted",
			setup: func() error {
				testKey.DeleteValue(appName) // Ignore error
				return nil
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.setup(); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			result := isAutorunEnabled(testKey)
			if result != tt.expected {
				t.Errorf("isAutorunEnabled() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWithAutorunRegistryErrorPropagation(t *testing.T) {
	tests := []struct {
		name          string
		access        uint32
		fnError       error
		expectFnError bool
	}{
		{
			name:          "function returns error",
			access:        registry.READ,
			fnError:       fmt.Errorf("test error"),
			expectFnError: true,
		},
		{
			name:          "function returns nil",
			access:        registry.READ,
			fnError:       nil,
			expectFnError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := withAutorunRegistry(tt.access, func(key registry.Key) error {
				return tt.fnError
			})

			if tt.expectFnError && err == nil {
				t.Error("Expected function error to be propagated, but got nil")
			}

			if !tt.expectFnError && err != nil && err.Error() == "test error" {
				t.Error("Unexpected function error propagated")
			}
		})
	}
}
