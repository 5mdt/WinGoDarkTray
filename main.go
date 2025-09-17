package main

import (
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

var version string
var app *App

func main() {
	app = NewApp(version)
	_ = installEventLogSource()
	systray.Run(app.onReady, onExit)
}

func onExit() {}

func (a *App) onReady() {
	a.initializeApp()
	if err := a.setupEventLog(); err != nil {
		return
	}

	autorunItem, quitItem := a.createMenuItems()
	a.initializeMenuState(autorunItem)
	a.startEventHandlers(autorunItem, quitItem)
}

func (a *App) initializeApp() {
	systray.SetIcon(icon)
	systray.SetTooltip(tooltips.Default)
	logEvent(eventlog.Info, fmt.Sprintf("WinGoDarkTray started and running, Version: %s", a.version))
}

func (a *App) setupEventLog() error {
	if err := installEventLogSource(); err != nil {
		showError("Failed to install event log source: " + err.Error())
		logEvent(eventlog.Error, fmt.Sprintf("Failed to install event log source: %s, Version: %s", err.Error(), a.version))
		return err
	}
	return nil
}

func (a *App) createMenuItems() (*systray.MenuItem, *systray.MenuItem) {
	appNameItem := systray.AddMenuItem(menuTitles.AppName, "")
	go func() {
		<-appNameItem.ClickedCh
		openBrowser(projectLink)
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
	go a.handleMenuItemClicks(autorunItem, quitItem)
	go checkForUpdate(a.version, a.updateNowItem)
}
