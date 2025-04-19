package main

import (
	"github.com/getlantern/systray"
)

func main() {
	_ = installEventLogSource()
	systray.Run(onReady, onExit)
}

func onExit() {}

func onReady() {
	systray.SetIcon(icon)
	systray.SetTooltip(tooltips.Default)

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
