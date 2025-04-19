package main

import (
	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onExit() {}

func onReady() {
	systray.SetIcon(icon)
	systray.SetTooltip(messages.ToggleTooltip)

	appNameItem := systray.AddMenuItem(menuTitles.AppName, "")
	go func() {
		<-appNameItem.ClickedCh
		openBrowser(projectLink)
	}()

	autorunItem := systray.AddMenuItem(menuTitles.EnableAutorun, "Enable/Disable autorun at startup")
	systray.AddSeparator()

	toggleSystemItem = systray.AddMenuItem("", "Toggle system-wide theme (light/dark)")
	systray.AddSeparator()
	toggleAppItem = systray.AddMenuItem("", "Toggle app theme (light/dark)")
	toggleWindowsItem = systray.AddMenuItem(menuTitles.ToggleWindowsMode, "Toggle Windows theme (light/dark)")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem(menuTitles.Quit, "Quit the app")

	updateAutorunStatus(autorunItem)
	updateThemeToggleTitles(toggleSystemItem, toggleAppItem)

	go handleMenuItemClicks(toggleSystemItem, toggleAppItem, toggleWindowsItem, autorunItem, quitItem)
}
