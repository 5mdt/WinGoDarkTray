package main

import (
	"time"

	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc/eventlog"
)

const (
	themeRegistryPath = `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`
)

// withThemeRegistry executes a function with an open theme registry key
func withThemeRegistry(access uint32, fn func(registry.Key) error) error {
	key, err := openRegistryKey(themeRegistryPath, access)
	if err != nil {
		return err
	}
	defer key.Close()
	return fn(key)
}

// getCurrentAppThemeMode reads the current app theme mode from registry
func getCurrentAppThemeMode() (uint64, error) {
	var appMode uint64
	err := withThemeRegistry(registry.QUERY_VALUE, func(key registry.Key) error {
		var err error
		appMode, _, err = key.GetIntegerValue("AppsUseLightTheme")
		return err
	})
	return appMode, err
}

// setBothThemeModes sets both app and system theme modes to the same value
func setBothThemeModes(lightMode bool) error {
	var newMode uint32
	if lightMode {
		newMode = 1
	} else {
		newMode = 0
	}

	return withThemeRegistry(registry.SET_VALUE, func(key registry.Key) error {
		if err := key.SetDWordValue("AppsUseLightTheme", newMode); err != nil {
			return err
		}
		return key.SetDWordValue("SystemUsesLightTheme", newMode)
	})
}

func (a *App) toggleSystemMode() {
	currentAppMode, err := getCurrentAppThemeMode()
	if err != nil {
		showError("Failed to read current theme: " + err.Error())
		return
	}

	switchingToLight := currentAppMode == 0

	if err := setBothThemeModes(switchingToLight); err != nil {
		showError("Failed to toggle system mode: " + err.Error())
		logEvent(eventlog.Error, "Failed to toggle system mode: "+err.Error())
		return
	}

	if switchingToLight {
		logEvent(eventlog.Info, "Switching both to light mode...")
	} else {
		logEvent(eventlog.Info, "Switching both to dark mode...")
	}

	setTemporaryTooltip("Both app and system theme switched", 2*time.Second)
	a.updateThemeToggleTitles()
}

// toggleSingleTheme toggles a specific theme setting by registry key
func (a *App) toggleSingleTheme(registryKey string) {
	err := withThemeRegistry(registry.QUERY_VALUE|registry.SET_VALUE, func(key registry.Key) error {
		current, _, err := key.GetIntegerValue(registryKey)
		if err != nil {
			return err
		}

		var newMode uint32
		if current == 1 {
			newMode = 0
		} else {
			newMode = 1
		}

		return key.SetDWordValue(registryKey, newMode)
	})

	if err != nil {
		showError("Failed to toggle theme: " + err.Error())
		return
	}

	a.updateThemeToggleTitles()
}

func (a *App) updateThemeToggleTitles() {
	var appMode, systemMode uint64

	err := withThemeRegistry(registry.QUERY_VALUE, func(key registry.Key) error {
		appMode, _, _ = key.GetIntegerValue("AppsUseLightTheme")
		systemMode, _, _ = key.GetIntegerValue("SystemUsesLightTheme")
		return nil
	})

	if err != nil {
		showError("Failed to read current theme: " + err.Error())
		return
	}

	if appMode == 1 {
		a.toggleAppItem.SetTitle(menuTitles.ToggleAppToDark)
		a.toggleSystemItem.SetTitle(menuTitles.ToggleBothToDark)
	} else {
		a.toggleAppItem.SetTitle(menuTitles.ToggleAppToLight)
		a.toggleSystemItem.SetTitle(menuTitles.ToggleBothToLight)
	}

	if systemMode == 1 {
		a.toggleWindowsItem.SetTitle(menuTitles.ToggleWinToDark)
	} else {
		a.toggleWindowsItem.SetTitle(menuTitles.ToggleWinToLight)
	}
}

func (a *App) toggleAppMode() {
	a.toggleSingleTheme("AppsUseLightTheme")
}

func (a *App) toggleWindowsMode() {
	a.toggleSingleTheme("SystemUsesLightTheme")
}
