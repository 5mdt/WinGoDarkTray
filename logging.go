package main

import (
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc/eventlog"
)

func isAdmin() bool {

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, registry.QUERY_VALUE)
	if err != nil {

		return false
	}
	defer key.Close()
	return true
}

func eventLogSourceExists() bool {
	elog, err := eventlog.Open(appName)
	if err != nil {

		return false
	}
	defer elog.Close()
	return true
}

func installEventLogSource() error {

	if !isAdmin() {
		return fmt.Errorf("this action requires administrative privileges")
	}

	if eventLogSourceExists() {

		return nil
	}

	err := eventlog.InstallAsEventCreate(appName, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		return fmt.Errorf("failed to install event log source: %v", err)
	}

	fmt.Println("Event log source installed successfully")
	return nil
}

func logEvent(eventType uint32, message string) {
	elog, err := eventlog.Open(appName)
	if err != nil {

		fmt.Println("Failed to open event log:", err)
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

func showError(message string) {

	logEvent(eventlog.Error, message)

	err := beeep.Notify(notificationTexts.Error, message, "")
	if err != nil {

		systray.SetTooltip(tooltips.Error + message)
		time.Sleep(3 * time.Second)
		systray.SetTooltip(tooltips.Default)
	}
}
