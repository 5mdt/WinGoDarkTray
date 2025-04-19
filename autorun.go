package main

import (
	"time"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc/eventlog"
)

func toggleAutorun(autorunItem *systray.MenuItem) {
	key, err := openRegistryKey(autorunRegistryKey, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		showError("Failed to open autorun registry key: " + err.Error())
		logEvent(eventlog.Error, "Failed to open autorun registry key: "+err.Error())
		return
	}
	defer key.Close()

	if isAutorunEnabled(key) {
		updateAutorun(key, autorunItem, false)
		logEvent(eventlog.Info, "Autorun disabled")
		autorunItem.SetTitle(menuTitles.EnableAutorunUnchecked)
		systray.SetTooltip("Autorun disabled")
	} else {
		updateAutorun(key, autorunItem, true)
		logEvent(eventlog.Info, "Autorun enabled")
		autorunItem.SetTitle(menuTitles.EnableAutorunChecked)
		systray.SetTooltip("Autorun enabled")
	}

	time.Sleep(2 * time.Second)
	systray.SetTooltip(tooltips.Default)
}

func updateAutorun(key registry.Key, autorunItem *systray.MenuItem, enable bool) {
	var err error
	if enable {
		err = addAutorun(key, autorunItem)
	} else {
		err = removeAutorun(key, autorunItem)
	}
	if err != nil {
		showError("Failed to update autorun: " + err.Error())
	}
}

func removeAutorun(key registry.Key, autorunItem *systray.MenuItem) error {
	if err := key.DeleteValue(appName); err != nil {
		return err
	}
	autorunItem.SetTitle(menuTitles.EnableAutorunUnchecked)
	systray.SetTooltip(tooltips.Default)
	return nil
}

func addAutorun(key registry.Key, autorunItem *systray.MenuItem) error {
	exePath, err := getExePath()
	if err != nil {
		return err
	}
	if err := key.SetStringValue(appName, exePath); err != nil {
		return err
	}
	autorunItem.SetTitle(menuTitles.EnableAutorunChecked)
	systray.SetTooltip(tooltips.Default)
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
