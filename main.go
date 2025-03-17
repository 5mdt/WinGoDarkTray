package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
)

//go:embed icon.ico
var icon []byte

// Constants for registry keys and application name
const (
	autorunRegistryKey = `Software\Microsoft\Windows\CurrentVersion\Run`
	appName            = "WinGoDarkTray"
)

func main() {
	// Initialize the system tray application
	systray.Run(onReady, onExit)
}

func onReady() {
	// Set the tray icon and tooltip
	systray.SetIcon(icon)
	systray.SetTooltip("Toggle between light and dark app mode")

	// Add menu items for toggling mode, enabling autorun, and quitting the app
	toggleItem := systray.AddMenuItem("Toggle app mode", "Toggle between light and dark app mode")
	autorunItem := systray.AddMenuItem("Enable Autorun", "Enable/Disable autorun at startup")
	quitItem := systray.AddMenuItem("Quit", "Quit the app")

	// Update the autorun status based on current registry settings
	updateAutorunStatus(autorunItem)

	// Event loop to handle user interactions with the tray menu
	go handleMenuItemClicks(toggleItem, autorunItem, quitItem)
}

func onExit() {
	// Perform any necessary cleanup when the app exits
}

func handleMenuItemClicks(toggleItem, autorunItem, quitItem *systray.MenuItem) {
	// Listen for menu item click events and trigger corresponding actions
	for {
		select {
		case <-toggleItem.ClickedCh:
			toggleMode()
		case <-autorunItem.ClickedCh:
			toggleAutorun(autorunItem)
		case <-quitItem.ClickedCh:
			systray.Quit()
		}
	}
}

func toggleMode() {
	// Open the registry key for Personalize settings
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		showError("Failed to open registry key: " + err.Error())
		return
	}
	defer key.Close()

	// Read the current app mode and system mode from the registry
	currentMode, _, err := key.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		showError("Failed to read AppsUseLightTheme: " + err.Error())
		return
	}
	systemMode, _, err := key.GetIntegerValue("SystemUsesLightTheme")
	if err != nil {
		showError("Failed to read SystemUsesLightTheme: " + err.Error())
		return
	}

	// Toggle the app mode based on the current settings
	var newMode uint32
	if currentMode == 0 && systemMode == 0 {
		// Switch to light mode
		newMode = 1
		fmt.Println("Switching to light app mode and light system mode...")
	} else {
		// Switch to dark mode
		newMode = 0
		fmt.Println("Switching to dark app mode and dark system mode...")
	}

	// Update the app mode and system mode in the registry
	if err := updateRegistryMode(key, newMode); err != nil {
		return
	}

	// Provide feedback on the mode change
	systray.SetTooltip("Mode switched successfully!")
	time.Sleep(2 * time.Second)
	systray.SetTooltip("Toggle between light and dark app mode") // Reset tooltip
}

func updateRegistryMode(key registry.Key, newMode uint32) error {
	// Update both app mode and system mode in the registry
	if err := key.SetDWordValue("AppsUseLightTheme", newMode); err != nil {
		showError("Failed to set AppsUseLightTheme: " + err.Error())
		return err
	}
	if err := key.SetDWordValue("SystemUsesLightTheme", newMode); err != nil {
		showError("Failed to set SystemUsesLightTheme: " + err.Error())
		return err
	}
	return nil
}

func showError(message string) {
	// Display an error message in the system tray tooltip
	systray.SetTooltip("Error: " + message)
	time.Sleep(3 * time.Second)
	systray.SetTooltip("Toggle between light and dark app mode") // Reset tooltip
}

func toggleAutorun(autorunItem *systray.MenuItem) {
	// Open the registry key for Run settings
	key, err := registry.OpenKey(registry.CURRENT_USER, autorunRegistryKey, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		showError("Failed to open autorun registry key: " + err.Error())
		return
	}
	defer key.Close()

	// Check if the app is set to run at startup
	_, _, err = key.GetStringValue(appName)
	if err == nil {
		// If the app is set to autorun, remove it
		if err := removeAutorun(key, autorunItem); err != nil {
			return
		}
	} else {
		// If the app is not set to autorun, add it
		if err := addAutorun(key, autorunItem); err != nil {
			return
		}
	}

	time.Sleep(2 * time.Second)
	systray.SetTooltip("Toggle between light and dark app mode") // Reset tooltip
}

func removeAutorun(key registry.Key, autorunItem *systray.MenuItem) error {
	// Remove the app from the autorun registry
	if err := key.DeleteValue(appName); err != nil {
		showError("Failed to remove autorun: " + err.Error())
		return err
	}
	// Update the menu item label to indicate autorun is disabled
	autorunItem.SetTitle("Enable Autorun (❌)")
	systray.SetTooltip("Autorun disabled!")
	return nil
}

func addAutorun(key registry.Key, autorunItem *systray.MenuItem) error {
	// Get the executable path for the app
	exePath, err := getExePath()
	if err != nil {
		showError("Failed to find executable path: " + err.Error())
		return err
	}
	// Set the app to autorun at startup
	if err := key.SetStringValue(appName, exePath); err != nil {
		showError("Failed to set autorun: " + err.Error())
		return err
	}
	// Update the menu item label to indicate autorun is enabled
	autorunItem.SetTitle("Enable Autorun (✔)")
	systray.SetTooltip("Autorun enabled!")
	return nil
}

func getExePath() (string, error) {
	// Return the path of the currently running executable
	return exec.LookPath(os.Args[0])
}

func updateAutorunStatus(autorunItem *systray.MenuItem) {
	// Open the registry key for Run settings
	key, err := registry.OpenKey(registry.CURRENT_USER, autorunRegistryKey, registry.QUERY_VALUE)
	if err != nil {
		showError("Failed to open autorun registry key: " + err.Error())
		return
	}
	defer key.Close()

	// Check if the app is already set to autorun and update the menu item label
	_, _, err = key.GetStringValue(appName)
	if err == nil {
		// App is set to autorun, update the menu item to reflect that
		autorunItem.SetTitle("Enable Autorun (✔)")
	} else {
		// App is not set to autorun, update the menu item to reflect that
		autorunItem.SetTitle("Enable Autorun (❌)")
	}
}
