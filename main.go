package main

import (
	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/svc/eventlog"
)

func main() {
	_ = installEventLogSource()
	systray.Run(onReady, onExit)
}

func onExit() {}

func onReady() {
	systray.SetIcon(icon)
	systray.SetTooltip(tooltips.Default)
	logEvent(eventlog.Info, "WinGoDarkTray started and running")
	err := installEventLogSource()
	if err != nil {
		showError("Failed to install event log source: " + err.Error())
		logEvent(eventlog.Error, "Failed to install event log source: "+err.Error())
		return
	}

	appNameItem := systray.AddMenuItem(menuTitles.AppName, "")
	go func() {
		<-appNameItem.ClickedCh
		openBrowser(projectLink)
	}()

	autorunItem := systray.AddMenuItem(menuTitles.EnableAutorun, "")
	systray.AddSeparator()

	toggleSystemItem = systray.AddMenuItem("", "")
	systray.AddSeparator()

	toggleAppItem = systray.AddMenuItem("", "")
	toggleWindowsItem = systray.AddMenuItem("", "")
	systray.AddSeparator()

	quitItem := systray.AddMenuItem(menuTitles.Quit, "Exit the application")

	updateAutorunStatus(autorunItem)
	updateThemeToggleTitles(toggleSystemItem, toggleAppItem, toggleWindowsItem)

	go handleMenuItemClicks(toggleSystemItem, toggleAppItem, toggleWindowsItem, autorunItem, quitItem)
}
