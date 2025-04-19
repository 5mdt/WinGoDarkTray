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
