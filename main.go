package main

import (
	"fmt"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/svc/eventlog"
)

var version string

func main() {
	if version == "" {
		version = "v0.0.0"
	}
	_ = installEventLogSource()
	systray.Run(onReady, onExit)
}

func onExit() {}

func onReady() {

	systray.SetIcon(icon)
	systray.SetTooltip(tooltips.Default)

	logEvent(eventlog.Info, fmt.Sprintf("WinGoDarkTray started and running, Version: %s", version))

	err := installEventLogSource()
	if err != nil {
		showError("Failed to install event log source: " + err.Error())
		logEvent(eventlog.Error, fmt.Sprintf("Failed to install event log source: %s, Version: %s", err.Error(), version))
		return
	}

	appNameItem := systray.AddMenuItem(menuTitles.AppName, "")
	go func() {
		<-appNameItem.ClickedCh
		openBrowser(projectLink)
	}()

	autorunItem := systray.AddMenuItem(menuTitles.EnableAutorun, "")
	systray.AddSeparator()

	toggleSystemItem := systray.AddMenuItem("", "")
	systray.AddSeparator()

	toggleAppItem := systray.AddMenuItem("", "")
	toggleWindowsItem := systray.AddMenuItem("", "")
	systray.AddSeparator()

	updateNowItem = systray.AddMenuItem(menuTitles.UpdateNow, "Click to update the app")
	updateNowItem.Hide()

	quitItem := systray.AddMenuItem(menuTitles.Quit, "Exit the application")

	updateAutorunStatus(autorunItem)
	updateThemeToggleTitles(toggleSystemItem, toggleAppItem, toggleWindowsItem)

	go handleMenuItemClicks(toggleSystemItem, toggleAppItem, toggleWindowsItem, autorunItem, quitItem)

	go checkForUpdate(version)
}
