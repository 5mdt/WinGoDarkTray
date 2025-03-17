package main

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
)

//go:embed icon.ico
var icon []byte

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon) // Set the tray icon
	systray.SetTooltip("Toggle between light and dark app mode")

	// Add context menu (right-click)
	toggleItem := systray.AddMenuItem("Toggle app mode", "Toggle between light and dark app mode")
	quitItem := systray.AddMenuItem("Quit", "Quit the app")

	// Event loop to handle menu item clicks
	go func() {
		for {
			select {
			case <-toggleItem.ClickedCh:
				toggleMode()
			case <-quitItem.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {}

func toggleMode() {
	// Open the registry key for Personalize
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		showError("Failed to open registry key: " + err.Error())
		return
	}
	defer key.Close()

	// Get the current app mode
	currentMode, _, err := key.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		showError("Failed to read AppsUseLightTheme: " + err.Error())
		return
	}

	// Get the current system mode
	systemMode, _, err := key.GetIntegerValue("SystemUsesLightTheme")
	if err != nil {
		showError("Failed to read SystemUsesLightTheme: " + err.Error())
		return
	}

	// Toggle the mode (0 = dark, 1 = light)
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

	// Set the new app mode in the registry
	err = key.SetDWordValue("AppsUseLightTheme", newMode)
	if err != nil {
		showError("Failed to set AppsUseLightTheme: " + err.Error())
		return
	}

	// Set the new system mode in the registry
	err = key.SetDWordValue("SystemUsesLightTheme", newMode)
	if err != nil {
		showError("Failed to set SystemUsesLightTheme: " + err.Error())
		return
	}

	// Provide feedback that the mode has been changed successfully
	systray.SetTooltip("Mode switched successfully!")
	time.Sleep(2 * time.Second)                                  // Show success tooltip for 2 seconds
	systray.SetTooltip("Toggle between light and dark app mode") // Reset tooltip
}

func showError(message string) {
	// Show an error message in the system tray
	systray.SetTooltip("Error: " + message)
	time.Sleep(3 * time.Second)                                  // Show the error message for 3 seconds
	systray.SetTooltip("Toggle between light and dark app mode") // Reset the tooltip to the original text
}
