// autorun.go

//go:build windows
// +build windows

package main

import (
	"time"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc/eventlog"
)

// withAutorunRegistry executes a function with an open autorun registry key
func withAutorunRegistry(access uint32, fn func(registry.Key) error) error {
	key, err := openRegistryKey(autorunRegistryKey, access)
	if err != nil {
		return err
	}
	defer key.Close()
	return fn(key)
}

// updateAutorunUI updates the menu item title based on current state
func updateAutorunUI(autorunItem *systray.MenuItem, enabled bool) {
	if enabled {
		autorunItem.SetTitle(menuTitles.EnableAutorunChecked)
	} else {
		autorunItem.SetTitle(menuTitles.EnableAutorunUnchecked)
	}
}

func toggleAutorun(autorunItem *systray.MenuItem) {
	var currentlyEnabled bool

	err := withAutorunRegistry(registry.QUERY_VALUE|registry.SET_VALUE, func(key registry.Key) error {
		currentlyEnabled = isAutorunEnabled(key)

		if currentlyEnabled {
			if err := key.DeleteValue(appName); err != nil {
				return err
			}
			logEvent(eventlog.Info, "Autorun disabled")
		} else {
			exePath, err := getExePath()
			if err != nil {
				return err
			}
			if err := key.SetStringValue(appName, exePath); err != nil {
				return err
			}
			logEvent(eventlog.Info, "Autorun enabled")
		}
		return nil
	})

	if err != nil {
		showError("Failed to update autorun: " + err.Error())
		return
	}

	newState := !currentlyEnabled
	updateAutorunUI(autorunItem, newState)

	message := "Autorun disabled"
	if newState {
		message = "Autorun enabled"
	}
	// Use the setTemporaryTooltip from logging.go instead of duplicate function
	go setTemporaryTooltip(message, 2*time.Second)
}

func updateAutorunStatus(autorunItem *systray.MenuItem) {
	var enabled bool

	err := withAutorunRegistry(registry.QUERY_VALUE, func(key registry.Key) error {
		enabled = isAutorunEnabled(key)
		return nil
	})

	if err != nil {
		showError("Failed to open autorun registry key: " + err.Error())
		return
	}

	updateAutorunUI(autorunItem, enabled)
}

func isAutorunEnabled(key registry.Key) bool {
	_, _, err := key.GetStringValue(appName)
	return err == nil
}
