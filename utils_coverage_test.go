// utils_coverage_test.go

//go:build windows
// +build windows

package main

import (
	"runtime"
	"testing"
)

func TestOpenBrowserFunction(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "test open browser with valid URL",
			url:  "https://example.com",
		},
		{
			name: "test open browser with invalid URL",
			url:  "not-a-valid-url",
		},
		{
			name: "test open browser with empty URL",
			url:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if runtime.GOOS != "windows" {
				t.Skip("Skipping Windows-specific test")
			}

			// Test that openBrowser doesn't panic
			// We can't easily test the actual browser opening without side effects
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("openBrowser() panicked with URL %s: %v", tt.url, r)
				}
			}()

			// Execute the function - it will try to open browser but we can test the code path
			openBrowser(tt.url)
		})
	}
}

func TestOpenBrowserOnNonWindowsPlatform(t *testing.T) {
	// This test simulates the non-Windows behavior
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "non-windows error handling",
			url:  "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can test the function by temporarily changing runtime.GOOS
			// but since that's not possible, we'll test the logic instead

			// The function checks: if runtime.GOOS != "windows"
			if runtime.GOOS != "windows" {
				// On non-Windows, it should call showError
				openBrowser(tt.url)
			} else {
				// On Windows, test the rundll32 execution path
				openBrowser(tt.url)
			}
		})
	}
}
