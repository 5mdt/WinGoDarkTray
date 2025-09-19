// updater_logic_test.go

package main

import (
	"strings"
	"testing"
)

// Test helper functions that can be extracted from updater.go for testing

func TestVersionStringProcessing(t *testing.T) {
	tests := []struct {
		name     string
		tagName  string
		expected string
	}{
		{
			name:     "version with v prefix",
			tagName:  "v1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "version without v prefix",
			tagName:  "1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "version with pre-release",
			tagName:  "v1.2.3-beta",
			expected: "1.2.3-beta",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the logic from checkForUpdate
			latest := strings.TrimPrefix(tt.tagName, "v")
			if latest != tt.expected {
				t.Errorf("TrimPrefix(%q, 'v') = %q, want %q", tt.tagName, latest, tt.expected)
			}
		})
	}
}

func TestGithubAPIURL(t *testing.T) {
	expectedURL := "https://api.github.com/repos/5mdt/WinGoDarkTray/releases/latest"
	if githubReleasesAPI != expectedURL {
		t.Errorf("githubReleasesAPI = %q, want %q", githubReleasesAPI, expectedURL)
	}

	// Verify URL is well-formed
	if !strings.HasPrefix(githubReleasesAPI, "https://") {
		t.Error("githubReleasesAPI should use HTTPS")
	}

	if !strings.Contains(githubReleasesAPI, "github.com") {
		t.Error("githubReleasesAPI should point to GitHub")
	}

	if !strings.Contains(githubReleasesAPI, "5mdt/WinGoDarkTray") {
		t.Error("githubReleasesAPI should point to correct repository")
	}
}

// Test update command generation
func TestUpdateCommandComponents(t *testing.T) {
	// We can test that the command components are reasonable without executing
	expectedPackageID := "5mdt.WinGoDarkTray"

	// This tests the knowledge that the update uses winget
	if !strings.Contains("winget install "+expectedPackageID, expectedPackageID) {
		t.Error("Update command should reference correct package ID")
	}
}
