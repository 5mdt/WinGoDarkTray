// tray_handlers.go

//go:build windows
// +build windows

package main

import (
	"context"

	"github.com/getlantern/systray"
)

func (a *App) handleMenuItemClicks(autorunItem, quitItem *systray.MenuItem) {
	a.handleMenuItemClicksWithContext(context.Background(), autorunItem, quitItem)
}

func (a *App) handleMenuItemClicksWithContext(ctx context.Context, autorunItem, quitItem *systray.MenuItem) {
	for {
		select {
		case <-ctx.Done():
			return
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
