package main

import (
	"fmt"
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
	logError(message)
	systray.SetTooltip(messages.Error + message)
	time.Sleep(3 * time.Second)
	systray.SetTooltip(messages.ToggleTooltip)
}

func logError(message string) {
	elog, err := eventlog.Open(appName)
	if err != nil {
		fmt.Println("EventLog open error:", err)
		return
	}
	defer elog.Close()

	elog.Error(1, message)
}

func installEventLogSource() error {
	return eventlog.InstallAsEventCreate(appName, eventlog.Error|eventlog.Warning|eventlog.Info)
}

func openRegistryKey(path string, access uint32) (registry.Key, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, path, access)
	if err != nil {
		return registry.Key(0), fmt.Errorf("failed to open registry key: %v", err)
	}
	return key, nil
}

func getExePath() (string, error) {
	return os.Executable()
}
