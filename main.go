// main.go

//go:build windows
// +build windows

package main

import (
	"context"
	"fmt"
	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/svc/eventlog"
)

// App holds the application state and dependencies
type App struct {
	toggleSystemItem  *systray.MenuItem
	toggleAppItem     *systray.MenuItem
	toggleWindowsItem *systray.MenuItem
	updateNowItem     *systray.MenuItem
	version           string
}

// NewApp creates a new application instance
func NewApp(version string) *App {
	if version == "" {
		version = "v0.0.0"
	}
	return &App{
		version: version,
	}
}

// Build-time version injection - set via ldflags during build
var buildVersion string

func main() {
	app := NewApp(buildVersion)
	systray.Run(app.onReady, onExit)
}

func onExit() {}

func (a *App) onReady() {
	a.initializeApp()
	if err := a.setupEventLog(); err != nil {
		// proceed without event log
	}
	logEvent(eventlog.Info, fmt.Sprintf("WinGoDarkTray started and running, Version: %s", a.version))
	autorunItem, quitItem := a.createMenuItems()
	a.initializeMenuState(autorunItem)
	a.startEventHandlers(autorunItem, quitItem)
}

func (a *App) initializeApp() {
	systray.SetIcon(icon)
	systray.SetTooltip(tooltips.Default)
}

func (a *App) setupEventLog() error {
	if err := installEventLogSource(); err != nil {
		logEvent(eventlog.Warning, fmt.Sprintf("Event log install failed: %s; continuing without event log. Version: %s", err.Error(), a.version))
		return err
	}
	return nil
}

func (a *App) createMenuItems() (*systray.MenuItem, *systray.MenuItem) {
	appNameItem := systray.AddMenuItem(menuTitles.AppName, "")
	go func() {
		for range appNameItem.ClickedCh {
			openBrowser(projectLink)
		}
	}()

	autorunItem := systray.AddMenuItem(menuTitles.EnableAutorun, "")
	systray.AddSeparator()

	a.toggleSystemItem = systray.AddMenuItem("", "")
	systray.AddSeparator()

	a.toggleAppItem = systray.AddMenuItem("", "")
	a.toggleWindowsItem = systray.AddMenuItem("", "")
	systray.AddSeparator()

	a.updateNowItem = systray.AddMenuItem(menuTitles.UpdateNow, "Click to update the app")
	a.updateNowItem.Hide()

	quitItem := systray.AddMenuItem(menuTitles.Quit, "Exit the application")
	return autorunItem, quitItem
}

func (a *App) initializeMenuState(autorunItem *systray.MenuItem) {
	updateAutorunStatus(autorunItem)
	a.updateThemeToggleTitles()
}

func (a *App) startEventHandlers(autorunItem, quitItem *systray.MenuItem) {
	a.startEventHandlersWithContext(context.Background(), autorunItem, quitItem)
}

func (a *App) startEventHandlersWithContext(ctx context.Context, autorunItem, quitItem *systray.MenuItem) {
	quitCh := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			close(quitCh)
			return
		case <-quitItem.ClickedCh:
			close(quitCh)
			systray.Quit()
		}
	}()

	go a.handleMenuItemClicksWithContext(ctx, autorunItem, quitItem)
	go checkForUpdateWithContext(ctx, a.version, a.updateNowItem)
	go startUpdateClickHandlerWithContext(ctx, a.updateNowItem, quitCh)
}
