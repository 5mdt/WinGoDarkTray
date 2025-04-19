package main

import (
	"time"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
)

func toggleAutorun(autorunItem *systray.MenuItem) {
	key, err := openRegistryKey(autorunRegistryKey, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		showError("Failed to open autorun registry key: " + err.Error())
		return
	}
	defer key.Close()

	if isAutorunEnabled(key) {
		removeAutorun(key, autorunItem)
	} else {
		addAutorun(key, autorunItem)
	}

	time.Sleep(2 * time.Second)
	systray.SetTooltip(messages.ToggleTooltip)
}

func removeAutorun(key registry.Key, autorunItem *systray.MenuItem) error {
	if err := key.DeleteValue(appName); err != nil {
		showError("Failed to remove autorun: " + err.Error())
		return err
	}
	autorunItem.SetTitle(menuTitles.EnableAutorunUnchecked)
	systray.SetTooltip(messages.AutorunDisabled)
	return nil
}

func addAutorun(key registry.Key, autorunItem *systray.MenuItem) error {
	exePath, err := getExePath()
	if err != nil {
		showError("Failed to find executable path: " + err.Error())
		return err
	}
	if err := key.SetStringValue(appName, exePath); err != nil {
		showError("Failed to set autorun: " + err.Error())
		return err
	}
	autorunItem.SetTitle(menuTitles.EnableAutorunChecked)
	systray.SetTooltip(messages.AutorunEnabled)
	return nil
}

func updateAutorunStatus(autorunItem *systray.MenuItem) {
	key, err := openRegistryKey(autorunRegistryKey, registry.QUERY_VALUE)
	if err != nil {
		showError("Failed to open autorun registry key: " + err.Error())
		return
	}
	defer key.Close()

	if isAutorunEnabled(key) {
		autorunItem.SetTitle(menuTitles.EnableAutorunChecked)
	} else {
		autorunItem.SetTitle(menuTitles.EnableAutorunUnchecked)
	}
}

func isAutorunEnabled(key registry.Key) bool {
	_, _, err := key.GetStringValue(appName)
	return err == nil
}
