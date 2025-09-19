package main

import (
	"strings"
	"testing"
)

func TestMenuTitlesStructure(t *testing.T) {
	// Test that all menu titles are non-empty
	if menuTitles.AppName == "" {
		t.Error("AppName should not be empty")
	}
	if menuTitles.EnableAutorun == "" {
		t.Error("EnableAutorun should not be empty")
	}
	if menuTitles.Quit == "" {
		t.Error("Quit should not be empty")
	}
	if menuTitles.EnableAutorunChecked == "" {
		t.Error("EnableAutorunChecked should not be empty")
	}
	if menuTitles.EnableAutorunUnchecked == "" {
		t.Error("EnableAutorunUnchecked should not be empty")
	}
	if menuTitles.ToggleAppToDark == "" {
		t.Error("ToggleAppToDark should not be empty")
	}
	if menuTitles.ToggleAppToLight == "" {
		t.Error("ToggleAppToLight should not be empty")
	}
	if menuTitles.ToggleWinToDark == "" {
		t.Error("ToggleWinToDark should not be empty")
	}
	if menuTitles.ToggleWinToLight == "" {
		t.Error("ToggleWinToLight should not be empty")
	}
	if menuTitles.ToggleBothToDark == "" {
		t.Error("ToggleBothToDark should not be empty")
	}
	if menuTitles.ToggleBothToLight == "" {
		t.Error("ToggleBothToLight should not be empty")
	}
	if menuTitles.UpdateNow == "" {
		t.Error("UpdateNow should not be empty")
	}
}

func TestMenuTitlesContent(t *testing.T) {
	// Test specific content expectations
	if !strings.Contains(menuTitles.AppName, "WinGoDarkTray") {
		t.Error("AppName should contain 'WinGoDarkTray'")
	}

	// Test autorun titles contain appropriate indicators
	if !strings.Contains(menuTitles.EnableAutorunChecked, "✔") {
		t.Error("EnableAutorunChecked should contain checkmark")
	}
	if !strings.Contains(menuTitles.EnableAutorunUnchecked, "✗") {
		t.Error("EnableAutorunUnchecked should contain X mark")
	}

	// Test theme toggle titles contain appropriate direction
	if !strings.Contains(menuTitles.ToggleAppToDark, "Dark") {
		t.Error("ToggleAppToDark should contain 'Dark'")
	}
	if !strings.Contains(menuTitles.ToggleAppToLight, "Light") {
		t.Error("ToggleAppToLight should contain 'Light'")
	}
	if !strings.Contains(menuTitles.ToggleWinToDark, "Windows") {
		t.Error("ToggleWinToDark should contain 'Windows'")
	}
	if !strings.Contains(menuTitles.ToggleWinToLight, "Windows") {
		t.Error("ToggleWinToLight should contain 'Windows'")
	}
	if !strings.Contains(menuTitles.ToggleBothToDark, "both") {
		t.Error("ToggleBothToDark should contain 'both'")
	}
	if !strings.Contains(menuTitles.ToggleBothToLight, "both") {
		t.Error("ToggleBothToLight should contain 'both'")
	}
}

func TestTooltipsStructure(t *testing.T) {
	// Test that all tooltips are non-empty
	if tooltips.Default == "" {
		t.Error("Default tooltip should not be empty")
	}
	if tooltips.Error == "" {
		t.Error("Error tooltip should not be empty")
	}

	// Test content expectations
	if !strings.Contains(tooltips.Default, "windows") || !strings.Contains(tooltips.Default, "toggle") {
		t.Error("Default tooltip should describe the app's purpose")
	}
}

func TestNotificationTextsStructure(t *testing.T) {
	// Test that all notification texts are non-empty
	if notificationTexts.Error == "" {
		t.Error("Error notification text should not be empty")
	}
	if notificationTexts.UpdateAvailableTitle == "" {
		t.Error("UpdateAvailableTitle should not be empty")
	}
	if notificationTexts.UpdateAvailableMessage == "" {
		t.Error("UpdateAvailableMessage should not be empty")
	}

	// Test content expectations
	if !strings.Contains(notificationTexts.Error, "WinGoDarkTray") {
		t.Error("Error notification should contain app name")
	}
	if !strings.Contains(notificationTexts.UpdateAvailableTitle, "Update") {
		t.Error("UpdateAvailableTitle should contain 'Update'")
	}
}
