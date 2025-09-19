// tray_handlers.go

//go:build windows
// +build windows

package main

import "github.com/getlantern/systray"

func (a *App) handleMenuItemClicks(autorunItem, quitItem *systray.MenuItem) {
	for {
		select {
		case <-a.toggleSystemItem.ClickedCh:
			a.toggleSystemMode()
		case <-a.toggleAppItem.ClickedCh:
			a.toggleAppMode()
		case <-a.toggleWindowsItem.ClickedCh:
			a.toggleWindowsMode()
		case <-autorunItem.ClickedCh:
			toggleAutorun(autorunItem)
		case <-quitItem.ClickedCh:
			systray.Quit()
		}
	}
}
