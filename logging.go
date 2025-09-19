// logging.go

//go:build windows
// +build windows

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

	// Event log message IDs
	infoEventID    = 1
	warningEventID = 2
	errorEventID   = 3
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

// logToConsole outputs a message to console with event type prefix
func logToConsole(eventType uint32, message string) {
	fmt.Printf("[%s] %s\n", getEventTypeName(eventType), message)
}

// logEvent attempts to log to Windows Event Log, falls back to console on failure
func logEvent(eventType uint32, message string) error {
	elog, err := eventlog.Open(appName)
	if err != nil {
		logToConsole(eventType, message)
		return fmt.Errorf("failed to open event log: %v", err)
	}
	defer elog.Close()

	switch eventType {
	case eventlog.Info:
		return elog.Info(infoEventID, message)
	case eventlog.Warning:
		return elog.Warning(warningEventID, message)
	case eventlog.Error:
		return elog.Error(errorEventID, message)
	default:
		return fmt.Errorf("unknown event type: %d", eventType)
	}
}

// logEventSafe logs an event and ignores errors (for fire-and-forget logging)
func logEventSafe(eventType uint32, message string) {
	if err := logEvent(eventType, message); err != nil {
		logToConsole(eventType, message)
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
	logEventSafe(eventlog.Error, message)
	showNotificationWithFallback(notificationTexts.Error, message, tooltips.Error+message)
}

// showNotificationWithFallback attempts desktop notification, falls back to tooltip
func showNotificationWithFallback(title, message, fallbackTooltip string) {
	if err := beeep.Notify(title, message, ""); err != nil {
		setTemporaryTooltip(fallbackTooltip, 3*time.Second)
	}
}

// setTemporaryTooltip sets a temporary tooltip that reverts after duration
func setTemporaryTooltip(message string, duration time.Duration) {
	systray.SetTooltip(message)
	go func() {
		time.Sleep(duration)
		systray.SetTooltip(tooltips.Default)
	}()
}
