package main

import (
	"os"
	"runtime"
	"testing"
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
