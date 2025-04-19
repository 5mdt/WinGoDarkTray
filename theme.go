package main

import (
	"time"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc/eventlog"
)

func toggleSystemMode() {
	key, err := openRegistryKey(`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		showError("Failed to open registry key: " + err.Error())
		logEvent(eventlog.Error, "Failed to open registry key for system theme: "+err.Error())
		return
	}
	defer key.Close()

	appMode, _, err := key.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		showError("Failed to read AppsUseLightTheme: " + err.Error())
		logEvent(eventlog.Error, "Failed to read AppsUseLightTheme: "+err.Error())
		return
	}
	_, _, err = key.GetIntegerValue("SystemUsesLightTheme")
	if err != nil {
		showError("Failed to read SystemUsesLightTheme: " + err.Error())
		logEvent(eventlog.Error, "Failed to read SystemUsesLightTheme: "+err.Error())
		return
	}

	var newMode uint32
	if appMode == 1 {
		newMode = 0
		logEvent(eventlog.Info, "Switching both to dark mode...")
	} else {
		newMode = 1
		logEvent(eventlog.Info, "Switching both to light mode...")
	}

	key.SetDWordValue("AppsUseLightTheme", newMode)
	key.SetDWordValue("SystemUsesLightTheme", newMode)

	systray.SetTooltip("Both app and system theme switched")
	time.Sleep(2 * time.Second)
	systray.SetTooltip(tooltips.Default)

	updateThemeToggleTitles(toggleSystemItem, toggleAppItem, toggleWindowsItem)
}

func toggleTheme(appKey, sysKey string) {
	key, err := openRegistryKey(`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		showError("Failed to open registry key: " + err.Error())
		return
	}
	defer key.Close()

	current, _, err := key.GetIntegerValue(appKey)
	if err != nil {
		showError("Failed to read registry value: " + err.Error())
		return
	}

	var newMode uint32
	if current == 1 {
		newMode = 0
	} else {
		newMode = 1
	}

	key.SetDWordValue(sysKey, newMode)

	systray.SetTooltip(tooltips.Default)
	time.Sleep(2 * time.Second)
	systray.SetTooltip(tooltips.Default)

	updateThemeToggleTitles(toggleSystemItem, toggleAppItem, toggleWindowsItem)
}

func updateThemeToggleTitles(bothItem, appItem, windowsItem *systray.MenuItem) {
	key, err := openRegistryKey(`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE)
	if err != nil {
		showError("Failed to read current theme: " + err.Error())
		return
	}
	defer key.Close()

	appMode, _, _ := key.GetIntegerValue("AppsUseLightTheme")
	systemMode, _, _ := key.GetIntegerValue("SystemUsesLightTheme")

	if appMode == 1 {
		appItem.SetTitle(trayTitles.ToggleAppToDark)
		bothItem.SetTitle(trayTitles.ToggleBothToDark)
	} else {
		appItem.SetTitle(trayTitles.ToggleAppToLight)
		bothItem.SetTitle(trayTitles.ToggleBothToLight)
	}

	if systemMode == 1 {
		windowsItem.SetTitle(trayTitles.ToggleWinToDark)
	} else {
		windowsItem.SetTitle(trayTitles.ToggleWinToLight)
	}
}

func toggleAppMode() {
	toggleTheme("AppsUseLightTheme", "AppsUseLightTheme")
}

func toggleWindowsMode() {
	toggleTheme("SystemUsesLightTheme", "SystemUsesLightTheme")
}
