package main

import (
	"time"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc/eventlog"
)

const themeRegistryPath = `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`

// withThemeRegistry executes a function with an open theme registry key
func withThemeRegistry(access uint32, fn func(registry.Key) error) error {
	key, err := openRegistryKey(themeRegistryPath, access)
	if err != nil {
		return err
	}
	defer key.Close()
	return fn(key)
}

func (a *App) toggleSystemMode() {
	var appMode uint64

	err := withThemeRegistry(registry.QUERY_VALUE|registry.SET_VALUE, func(key registry.Key) error {
		var err error
		appMode, _, err = key.GetIntegerValue("AppsUseLightTheme")
		if err != nil {
			return err
		}

		_, _, err = key.GetIntegerValue("SystemUsesLightTheme")
		if err != nil {
			return err
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
		return nil
	})

	if err != nil {
		showError("Failed to toggle system mode: " + err.Error())
		logEvent(eventlog.Error, "Failed to toggle system mode: "+err.Error())
		return
	}

	systray.SetTooltip("Both app and system theme switched")
	time.Sleep(2 * time.Second)
	systray.SetTooltip(tooltips.Default)

	a.updateThemeToggleTitles()
}

func (a *App) toggleTheme(appKey, sysKey string) {
	err := withThemeRegistry(registry.QUERY_VALUE|registry.SET_VALUE, func(key registry.Key) error {
		current, _, err := key.GetIntegerValue(appKey)
		if err != nil {
			return err
		}

		var newMode uint32
		if current == 1 {
			newMode = 0
		} else {
			newMode = 1
		}

		return key.SetDWordValue(sysKey, newMode)
	})

	if err != nil {
		showError("Failed to toggle theme: " + err.Error())
		return
	}

	systray.SetTooltip(tooltips.Default)
	time.Sleep(2 * time.Second)
	systray.SetTooltip(tooltips.Default)

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
	a.toggleTheme("AppsUseLightTheme", "AppsUseLightTheme")
}

func (a *App) toggleWindowsMode() {
	a.toggleTheme("SystemUsesLightTheme", "SystemUsesLightTheme")
}
