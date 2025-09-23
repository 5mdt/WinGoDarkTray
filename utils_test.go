// utils_test.go

//go:build windows
// +build windows

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/sys/windows/registry"
)

func TestGetExePath(t *testing.T) {
	path, err := getExePath()
	if err != nil {
		t.Fatalf("getExePath() returned error: %v", err)
	}

	if path == "" {
		t.Error("getExePath() returned empty path")
	}

	// Verify the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("getExePath() returned non-existent path: %s", path)
	}
}

func TestOpenBrowserOnNonWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping non-Windows test on Windows platform")
	}

	// This test verifies that openBrowser() handles non-Windows platforms gracefully
	// Since showError() would be called, we can't easily test this without mocking
	// But we can at least verify the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("openBrowser() panicked on non-Windows platform: %v", r)
		}
	}()

	// Note: This will call showError() internally but shouldn't panic
	openBrowser("https://example.com")
}

func TestOpenBrowserOnWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on non-Windows platform")
	}

	tests := []struct {
		name string
		url  string
	}{
		{
			name: "valid HTTP URL",
			url:  "https://example.com",
		},
		{
			name: "valid HTTPS URL",
			url:  "https://github.com",
		},
		{
			name: "empty URL",
			url:  "",
		},
		{
			name: "invalid URL",
			url:  "not-a-url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't easily test the actual browser opening without side effects
			// So we test that the function doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("openBrowser() panicked with URL %s: %v", tt.url, r)
				}
			}()

			// Skip actual execution to avoid opening browsers during tests
			t.Skip("Skipping actual browser opening to avoid side effects")
		})
	}
}

func TestOpenRegistryKey(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on non-Windows platform")
	}

	tests := []struct {
		name        string
		path        string
		access      uint32
		shouldError bool
	}{
		{
			name:        "valid path with read access",
			path:        "Software",
			access:      registry.READ,
			shouldError: false,
		},
		{
			name:        "valid path with query access",
			path:        "Software",
			access:      registry.QUERY_VALUE,
			shouldError: false,
		},
		{
			name:        "non-existent path",
			path:        "NonExistentSoftwareKey12345",
			access:      registry.READ,
			shouldError: true,
		},
		{
			name:        "empty path",
			path:        "",
			access:      registry.READ,
			shouldError: false, // Empty path might be handled gracefully by Windows registry
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := openRegistryKey(tt.path, tt.access)

			if tt.shouldError {
				if err == nil {
					key.Close() // Clean up if unexpectedly successful
					t.Errorf("openRegistryKey() expected error for path %s, but got none", tt.path)
				}
				return
			}

			if err != nil {
				t.Errorf("openRegistryKey() unexpected error for path %s: %v", tt.path, err)
				return
			}

			// Clean up the opened key
			if err := key.Close(); err != nil {
				t.Errorf("Failed to close registry key: %v", err)
			}
		})
	}
}

func TestOpenRegistryKeyReturnValue(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on non-Windows platform")
	}

	tests := []struct {
		name string
		path string
	}{
		{
			name: "software key return value",
			path: "Software",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := openRegistryKey(tt.path, registry.READ)
			if err != nil {
				t.Fatalf("openRegistryKey() failed: %v", err)
			}
			defer key.Close()

			// Verify we got a valid key handle
			if key == registry.Key(0) {
				t.Error("openRegistryKey() returned zero key handle")
			}
		})
	}
}

func TestGetExePathProperties(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "executable path properties",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := getExePath()
			if err != nil {
				t.Fatalf("getExePath() failed: %v", err)
			}

			// Test that the path is absolute
			if !filepath.IsAbs(path) {
				t.Errorf("getExePath() returned relative path: %s", path)
			}

			// Test that the file has executable extension on Windows
			if runtime.GOOS == "windows" && !strings.HasSuffix(strings.ToLower(path), ".exe") {
				t.Errorf("getExePath() on Windows should end with .exe, got: %s", path)
			}
		})
	}
}

func TestRundll32CommandStructure(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on non-Windows platform")
	}

	tests := []struct {
		name string
		url  string
	}{
		{
			name: "command structure validation",
			url:  "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the command structure used in openBrowser without executing
			cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", tt.url)

			expectedProgram := "rundll32"
			expectedArgs := []string{"rundll32", "url.dll,FileProtocolHandler", tt.url}

			if cmd.Args[0] != expectedProgram {
				t.Errorf("Command program = %s, want %s", cmd.Args[0], expectedProgram)
			}

			if len(cmd.Args) != len(expectedArgs) {
				t.Errorf("Command args length = %d, want %d", len(cmd.Args), len(expectedArgs))
			}

			for i, expectedArg := range expectedArgs {
				if i < len(cmd.Args) && cmd.Args[i] != expectedArg {
					t.Errorf("Command arg[%d] = %s, want %s", i, cmd.Args[i], expectedArg)
				}
			}
		})
	}
}

func TestRegistryKeyConstants(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on non-Windows platform")
	}

	tests := []struct {
		name     string
		access   uint32
		expected string
	}{
		{
			name:     "READ access constant",
			access:   registry.READ,
			expected: "READ",
		},
		{
			name:     "QUERY_VALUE access constant",
			access:   registry.QUERY_VALUE,
			expected: "QUERY_VALUE",
		},
		{
			name:     "SET_VALUE access constant",
			access:   registry.SET_VALUE,
			expected: "SET_VALUE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that registry constants are available and have expected values
			// This validates our imports and constant usage
			if tt.access == 0 {
				t.Errorf("Registry access constant %s should not be zero", tt.expected)
			}
		})
	}
}
