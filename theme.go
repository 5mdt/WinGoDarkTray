package main

import (
	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc/eventlog"
	"time"
)

func toggleSystemMode() {
	key, err := openRegistryKey(`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		showError("Failed to open registry key: " + err.Error())
		return
	}
	defer key.Close()

	appMode, _, err := key.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		showError("Failed to read AppsUseLightTheme: " + err.Error())
		return
	}

	var newMode uint32
	if appMode == 1 {
		newMode = 0
		logEvent(eventlog.Info, "Switching to dark mode (App + Windows)...")
	} else {
		newMode = 1
		logEvent(eventlog.Info, "Switching to light mode (App + Windows)...")
	}

	key.SetDWordValue("AppsUseLightTheme", newMode)
	key.SetDWordValue("SystemUsesLightTheme", newMode)

	systray.SetTooltip(messages.ModeSwitched)
	time.Sleep(2 * time.Second)
	systray.SetTooltip(messages.ToggleTooltip)
}

func toggleAppMode() {
	key, err := openRegistryKey(`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		showError("Failed to open registry key: " + err.Error())
		return
	}
	defer key.Close()

	current, _, err := key.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		showError("Failed to read AppsUseLightTheme: " + err.Error())
		return
	}

	newMode := uint32(1)
	if current == 1 {
		newMode = 0
		logEvent(eventlog.Info, "Switching to dark mode...")
	} else {
		logEvent(eventlog.Info, "Switching to light mode...")
	}

	key.SetDWordValue("AppsUseLightTheme", newMode)
	systray.SetTooltip(messages.ModeSwitched)
	time.Sleep(2 * time.Second)
	systray.SetTooltip(messages.ToggleTooltip)
}

func toggleWindowsMode() {
	toggleSystemMode()
}
