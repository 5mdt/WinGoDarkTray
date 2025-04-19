package main

import "github.com/getlantern/systray"

func handleMenuItemClicks(toggleSystemItem, toggleAppItem, toggleWindowsItem, autorunItem, quitItem *systray.MenuItem) {
	for {
		select {
		case <-toggleSystemItem.ClickedCh:
			toggleSystemMode()
		case <-toggleAppItem.ClickedCh:
			toggleAppMode()
		case <-toggleWindowsItem.ClickedCh:
			toggleWindowsMode()
		case <-autorunItem.ClickedCh:
			toggleAutorun(autorunItem)
		case <-quitItem.ClickedCh:
			systray.Quit()
		}
	}
}
