package main

import (
	_ "embed"
)

var (
	//go:embed icon.ico
	icon []byte
)

const (
	autorunRegistryKey = `Software\Microsoft\Windows\CurrentVersion\Run`
	appName            = "WinGoDarkTray"
	projectLink        = "https://github.com/5mdt/WinGoDarkTray"
)

var messages = struct {
	ToggleTooltip   string
	AutorunEnabled  string
	AutorunDisabled string
	Error           string
	ModeSwitched    string
}{
	ToggleTooltip:   "Toggle between themes for system-wide, app, and Windows",
	AutorunEnabled:  "Autorun enabled!",
	AutorunDisabled: "Autorun disabled!",
	Error:           "Error: ",
	ModeSwitched:    "Mode switched successfully!",
}

var menuTitles = struct {
	AppName                string
	ToggleSystemMode       string
	ToggleAppMode          string
	ToggleWindowsMode      string
	EnableAutorun          string
	Quit                   string
	EnableAutorunChecked   string
	EnableAutorunUnchecked string
}{
	AppName:                "WinGoDarkTray 🔗",
	ToggleSystemMode:       "Toggle System-Wide theme",
	ToggleAppMode:          "Toggle Apps theme",
	ToggleWindowsMode:      "Toggle Windows theme",
	EnableAutorun:          "Enable Autorun",
	Quit:                   "Quit",
	EnableAutorunChecked:   "Enable Autorun (✔)",
	EnableAutorunUnchecked: "Enable Autorun (❌)",
}
