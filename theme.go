package main

import (
	"time"

	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc/eventlog"
)

const (
	themeRegistryPath     = `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`
	regValAppsUseLight    = "AppsUseLightTheme"
	regValSystemUsesLight = "SystemUsesLightTheme"
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
		appMode, _, err = key.GetIntegerValue(regValAppsUseLight)
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
		if err := key.SetDWordValue(regValAppsUseLight, newMode); err != nil {
			return err
		}
		return key.SetDWordValue(regValSystemUsesLight, newMode)
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

	target := "Light"
	if !switchingToLight {
		target = "Dark"
	}
	setTemporaryTooltip("Switched both to "+target+" mode", 2*time.Second)
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
		showError("Failed to toggle theme (" + registryKey + "): " + err.Error())
		logEvent(eventlog.Error, "Failed to toggle theme ("+registryKey+"): "+err.Error())
		return
	}

	setTemporaryTooltip("Theme toggled", 2*time.Second)
	a.updateThemeToggleTitles()
}

func (a *App) updateThemeToggleTitles() {
	var appMode, systemMode uint64

	err := withThemeRegistry(registry.QUERY_VALUE, func(key registry.Key) error {
		var err error
		if appMode, _, err = key.GetIntegerValue(regValAppsUseLight); err != nil {
			return err
		}
		if systemMode, _, err = key.GetIntegerValue(regValSystemUsesLight); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		showError("Failed to read current theme: " + err.Error())
		return
	}

	bothLight := (appMode == 1 && systemMode == 1)
	if appMode == 1 {
		a.toggleAppItem.SetTitle(menuTitles.ToggleAppToDark)
	} else {
		a.toggleAppItem.SetTitle(menuTitles.ToggleAppToLight)
	}
	if bothLight {
		a.toggleSystemItem.SetTitle(menuTitles.ToggleBothToDark)
	} else {
		a.toggleSystemItem.SetTitle(menuTitles.ToggleBothToLight)
	}

	if systemMode == 1 {
		a.toggleWindowsItem.SetTitle(menuTitles.ToggleWinToDark)
	} else {
		a.toggleWindowsItem.SetTitle(menuTitles.ToggleWinToLight)
	}
}

func (a *App) toggleAppMode() {
	a.toggleSingleTheme(regValAppsUseLight)
}

func (a *App) toggleWindowsMode() {
	a.toggleSingleTheme(regValSystemUsesLight)
}
