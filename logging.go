package main

import (
	"fmt"
	"time"
	"unsafe"
	"syscall"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/svc/eventlog"
)

var (
	advapi32                = syscall.NewLazyDLL("advapi32.dll")
	procCheckTokenMembership = advapi32.NewProc("CheckTokenMembership")
	procCreateWellKnownSid   = advapi32.NewProc("CreateWellKnownSid")
)

const (
	WinBuiltinAdministratorsSid = 26
)

func isAdmin() bool {
	var sid [1024]byte
	sidSize := uint32(len(sid))

	ret, _, _ := procCreateWellKnownSid.Call(
		uintptr(WinBuiltinAdministratorsSid),
		0,
		uintptr(unsafe.Pointer(&sid[0])),
		uintptr(unsafe.Pointer(&sidSize)),
	)

	if ret == 0 {
		return false
	}

	var isMember int32
	ret, _, _ = procCheckTokenMembership.Call(
		0,
		uintptr(unsafe.Pointer(&sid[0])),
		uintptr(unsafe.Pointer(&isMember)),
	)

	return ret != 0 && isMember != 0
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
