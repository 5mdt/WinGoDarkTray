package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc/eventlog"
)

//go:embed icon.ico
var icon []byte

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

func main() {
	// Optional: register the app as an event source (run once with admin)
	// _ = installEventLogSource()

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon)
	systray.SetTooltip(messages.ToggleTooltip)

	appNameItem := systray.AddMenuItem(menuTitles.AppName, "")
	go func() {
		<-appNameItem.ClickedCh
		openBrowser(projectLink)
	}()

	autorunItem := systray.AddMenuItem(menuTitles.EnableAutorun, "Enable/Disable autorun at startup")
	systray.AddSeparator()

	toggleSystemItem := systray.AddMenuItem(menuTitles.ToggleSystemMode, "Toggle system-wide theme (light/dark)")
	systray.AddSeparator()
	toggleAppItem := systray.AddMenuItem(menuTitles.ToggleAppMode, "Toggle app theme (light/dark)")
	toggleWindowsItem := systray.AddMenuItem(menuTitles.ToggleWindowsMode, "Toggle Windows theme (light/dark)")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem(menuTitles.Quit, "Quit the app")

	updateAutorunStatus(autorunItem)

	go handleMenuItemClicks(toggleSystemItem, toggleAppItem, toggleWindowsItem, autorunItem, quitItem)
}

func onExit() {}

func handleMenuItemClicks(toggleSystemItem, toggleAppItem, toggleWindowsItem, autorunItem, quitItem *systray.MenuItem) {
	for {
		select {
		case <-toggleSystemItem.ClickedCh:
			toggleSystemMode()
		case <-toggleAppItem.ClickedCh:
			toggleAppMode()
		case <-toggleWindowsItem.ClickedCh:
			toggleWindowsMode()
		case <-autorunItem.ClickedCh:
			toggleAutorun(autorunItem)
		case <-quitItem.ClickedCh:
			systray.Quit()
		}
	}
}

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


func toggleSystemMode() {
	key, err := openRegistryKey(`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		showError("Failed to open registry key: " + err.Error())
		return
	}
	defer key.Close()

	// Step 1: Read current App theme
	appMode, _, err := key.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		showError("Failed to read AppsUseLightTheme: " + err.Error())
		return
	}

	// Step 2: Toggle the theme (invert)
	var newMode uint32
	if appMode == 1 {
		newMode = 0 // Switch to dark mode
		fmt.Println("Switching to dark mode (App + Windows)...")
	} else {
		newMode = 1 // Switch to light mode
		fmt.Println("Switching to light mode (App + Windows)...")
	}

	// Step 3: Apply to both AppsUseLightTheme and SystemUsesLightTheme
	if err := key.SetDWordValue("AppsUseLightTheme", newMode); err != nil {
		showError("Failed to set AppsUseLightTheme: " + err.Error())
		return
	}
	if err := key.SetDWordValue("SystemUsesLightTheme", newMode); err != nil {
		showError("Failed to set SystemUsesLightTheme: " + err.Error())
		return
	}

	systray.SetTooltip(messages.ModeSwitched)
	time.Sleep(2 * time.Second)
	systray.SetTooltip(messages.ToggleTooltip)
}

func toggleAppMode() {
	key, err := openRegistryKey(`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		showError("Failed to open registry key for app theme: " + err.Error())
		return
	}
	defer key.Close()

	current, _, err := key.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		showError("Failed to read AppsUseLightTheme: " + err.Error())
		return
	}

	newMode := uint32(1)
	if current == 1 {
		newMode = 0
		fmt.Println("Switching to dark mode...")
	} else {
		fmt.Println("Switching to light mode...")
	}

	if err := key.SetDWordValue("AppsUseLightTheme", newMode); err != nil {
		showError("Failed to set AppsUseLightTheme: " + err.Error())
		return
	}

	systray.SetTooltip(messages.ModeSwitched)
	time.Sleep(2 * time.Second)
	systray.SetTooltip(messages.ToggleTooltip)
}

func toggleWindowsMode() {
	toggleSystemMode()
}

func toggleAutorun(autorunItem *systray.MenuItem) {
	key, err := openRegistryKey(autorunRegistryKey, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		showError("Failed to open autorun registry key: " + err.Error())
		return
	}
	defer key.Close()

	if isAutorunEnabled(key) {
		if err := removeAutorun(key, autorunItem); err != nil {
			return
		}
	} else {
		if err := addAutorun(key, autorunItem); err != nil {
			return
		}
	}

	time.Sleep(2 * time.Second)
	systray.SetTooltip(messages.ToggleTooltip)
}

func removeAutorun(key registry.Key, autorunItem *systray.MenuItem) error {
	if err := key.DeleteValue(appName); err != nil {
		showError("Failed to remove autorun: " + err.Error())
		return err
	}
	autorunItem.SetTitle(menuTitles.EnableAutorunUnchecked)
	systray.SetTooltip(messages.AutorunDisabled)
	return nil
}

func addAutorun(key registry.Key, autorunItem *systray.MenuItem) error {
	exePath, err := getExePath()
	if err != nil {
		showError("Failed to find executable path: " + err.Error())
		return err
	}
	if err := key.SetStringValue(appName, exePath); err != nil {
		showError("Failed to set autorun: " + err.Error())
		return err
	}
	autorunItem.SetTitle(menuTitles.EnableAutorunChecked)
	systray.SetTooltip(messages.AutorunEnabled)
	return nil
}

func getExePath() (string, error) {
	return os.Executable()
}

func updateAutorunStatus(autorunItem *systray.MenuItem) {
	key, err := openRegistryKey(autorunRegistryKey, registry.QUERY_VALUE)
	if err != nil {
		showError("Failed to open autorun registry key: " + err.Error())
		return
	}
	defer key.Close()

	if isAutorunEnabled(key) {
		autorunItem.SetTitle(menuTitles.EnableAutorunChecked)
	} else {
		autorunItem.SetTitle(menuTitles.EnableAutorunUnchecked)
	}
}

func isAutorunEnabled(key registry.Key) bool {
	_, _, err := key.GetStringValue(appName)
	return err == nil
}

func showError(message string) {
	logError(message)
	systray.SetTooltip(messages.Error + message)
	time.Sleep(3 * time.Second)
	systray.SetTooltip(messages.ToggleTooltip)
}

// ✅ NEW: Logs error to Windows Event Log
func logError(message string) {
	elog, err := eventlog.Open(appName)
	if err != nil {
		fmt.Println("EventLog open error:", err)
		return
	}
	defer elog.Close()

	elog.Error(1, message)
}

// ✅ OPTIONAL: Call this once with admin to register event source
func installEventLogSource() error {
	return eventlog.InstallAsEventCreate(appName, eventlog.Error|eventlog.Warning|eventlog.Info)
}

// Opens registry with error-safe fallback
func openRegistryKey(path string, access uint32) (registry.Key, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, path, access)
	if err != nil {
		return registry.Key(0), fmt.Errorf("failed to open registry key: %v", err)
	}
	return key, nil
}
