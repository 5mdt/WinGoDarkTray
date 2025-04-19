package main

import (
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc/eventlog"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func openBrowser(urlStr string) {
	if runtime.GOOS != "windows" {
		showError("Opening URLs is only supported on Windows.")
		return
	}

	err := exec.Command("rundll32", "url.dll,FileProtocolHandler", urlStr).Start()
	if err != nil {
		showError("Failed to open browser: " + err.Error())
	}
}

func showError(message string) {
	logEvent(eventlog.Error, message)
	err := beeep.Notify("WinGoDarkTray Error", message, "")
	if err != nil {
		systray.SetTooltip(tooltips.Error + message)
		time.Sleep(3 * time.Second)
		systray.SetTooltip(tooltips.Default)
	}
}

func logEvent(eventType uint32, message string) {
	elog, err := eventlog.Open(appName)
	if err != nil {
		return
	}
	defer elog.Close()

	switch eventType {
	case eventlog.Info:
		elog.Info(1, message)
	case eventlog.Warning:
		elog.Warning(2, message)
	case eventlog.Error:
		elog.Error(3, message)
	}
}

func installEventLogSource() error {
	return eventlog.InstallAsEventCreate(appName, eventlog.Error|eventlog.Warning|eventlog.Info)
}

func openRegistryKey(path string, access uint32) (registry.Key, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, path, access)
	if err != nil {
		return registry.Key(0), err
	}
	return key, nil
}

func getExePath() (string, error) {
	return os.Executable()
}
