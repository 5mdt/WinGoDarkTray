package main

import (
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/eventlog"
)

const (
	adminRoleID = "S-1-5-32-544"
)

func isAdmin() bool {
	token, err := windows.OpenCurrentProcessToken()
	if err != nil {
		return false
	}
	defer token.Close()

	adminSid, err := windows.StringToSid(adminRoleID)
	if err != nil {
		return false
	}

	isMember, err := token.IsMember(adminSid)
	return err == nil && isMember
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
	// If already exists, no need for admin privileges
	if eventLogSourceExists() {
		return nil
	}

	if !isAdmin() {
		return fmt.Errorf("administrative privileges required for event log registration")
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
		fmt.Printf("[%s] %s\n", getEventTypeName(eventType), message)
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

func getEventTypeName(eventType uint32) string {
	switch eventType {
	case eventlog.Info:
		return "INFO"
	case eventlog.Warning:
		return "WARNING"
	case eventlog.Error:
		return "ERROR"
	default:
		return "UNKNOWN"
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
