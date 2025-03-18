package main

import (
        _ "embed"
        "fmt"
        "log"
        "os"
        "os/exec"
        "runtime"
        "time"

        "github.com/getlantern/systray"
        "golang.org/x/sys/windows/registry"
)

//go:embed icon.ico
var icon []byte

var Version = "dev-build"

// Constants for registry keys and application name
const (
        autorunRegistryKey = `Software\Microsoft\Windows\CurrentVersion\Run`
        appName            = "WinGoDarkTray"
        projectLink        = "https://github.com/5mdt/WinGoDarkTray"
)

// Centralized message and title storage
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
        // Initialize the system tray application
        systray.Run(onReady, onExit)
}

func onReady() {
        // Set the tray icon and tooltip
        systray.SetIcon(icon)
        systray.SetTooltip(messages.ToggleTooltip)

        // Add inactive AppName that opens ProjectLink in the browser
        appNameItem := systray.AddMenuItem(menuTitles.AppName, "") // Just displaying as an info item, no action

        // Open the ProjectLink URL when clicked
        go func() {
                <-appNameItem.ClickedCh
                openBrowser(projectLink)
        }()

        autorunItem := systray.AddMenuItem(menuTitles.EnableAutorun, "Enable/Disable autorun at startup")
        systray.AddSeparator()

        // Add menu items for toggling system, app, and windows themes, enabling autorun, and quitting the app
        toggleSystemItem := systray.AddMenuItem(menuTitles.ToggleSystemMode, "Toggle system-wide theme (light/dark)")
        systray.AddSeparator()
        toggleAppItem := systray.AddMenuItem(menuTitles.ToggleAppMode, "Toggle app theme (light/dark)")
        toggleWindowsItem := systray.AddMenuItem(menuTitles.ToggleWindowsMode, "Toggle Windows theme (light/dark)")
        systray.AddSeparator()
        quitItem := systray.AddMenuItem(menuTitles.Quit, "Quit the app")

        // Update the autorun status based on current registry settings
        updateAutorunStatus(autorunItem)

        // Event loop to handle user interactions with the tray menu
        go handleMenuItemClicks(toggleSystemItem, toggleAppItem, toggleWindowsItem, autorunItem, quitItem)
}

func onExit() {
        // Perform any necessary cleanup when the app exits
}

func handleMenuItemClicks(toggleSystemItem, toggleAppItem, toggleWindowsItem, autorunItem, quitItem *systray.MenuItem) {
        // Listen for menu item click events and trigger corresponding actions
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
        var err error
        // Open the URL in the default browser without cmd window
        switch runtime.GOOS {
        case "windows":
                err = exec.Command("rundll32", "url.dll,FileProtocolHandler", urlStr).Start()
        case "darwin":
                err = exec.Command("open", urlStr).Start()
        case "linux":
                err = exec.Command("xdg-open", urlStr).Start()
        default:
                err = fmt.Errorf("unsupported platform")
        }
        if err != nil {
                showError("Failed to open browser: " + err.Error())
        }
}

func toggleSystemMode() {
        // Open the registry key for system-wide theme settings
        key, err := openRegistryKey(`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE|registry.SET_VALUE)
        if err != nil {
                showError("Failed to open registry key for system theme: " + err.Error())
                return
        }
        defer key.Close()

        // Read the current system theme mode from the registry
        currentSystemMode, _, err := key.GetIntegerValue("SystemUsesLightTheme")
        if err != nil {
                showError("Failed to read SystemUsesLightTheme: " + err.Error())
                return
        }

        // Toggle the system-wide theme mode
        var newMode uint32
        if currentSystemMode == 0 {
                // Switch to light mode
                newMode = 1
                fmt.Println("Switching system-wide theme to light mode...")
        } else {
                // Switch to dark mode
                newMode = 0
                fmt.Println("Switching system-wide theme to dark mode...")
        }

        // Update the system theme mode in the registry
        if err := key.SetDWordValue("SystemUsesLightTheme", newMode); err != nil {
                showError("Failed to set SystemUsesLightTheme: " + err.Error())
                return
        }

        // Provide feedback on the mode change
        systray.SetTooltip(messages.ModeSwitched)
        time.Sleep(2 * time.Second)
        systray.SetTooltip(messages.ToggleTooltip) // Reset tooltip
}

func toggleAppMode() {
        // Open the registry key for app-specific theme settings
        key, err := openRegistryKey(`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE|registry.SET_VALUE)
        if err != nil {
                showError("Failed to open registry key for app theme: " + err.Error())
                return
        }
        defer key.Close()

        // Read the current app theme mode from the registry
        currentAppMode, _, err := key.GetIntegerValue("AppsUseLightTheme")
        if err != nil {
                showError("Failed to read AppsUseLightTheme: " + err.Error())
                return
        }

        // Toggle the app-specific theme mode
        var newMode uint32
        if currentAppMode == 0 {
                // Switch to light mode
                newMode = 1
                fmt.Println("Switching app theme to light mode...")
        } else {
                // Switch to dark mode
                newMode = 0
                fmt.Println("Switching app theme to dark mode...")
        }

        // Update the app theme mode in the registry
        if err := key.SetDWordValue("AppsUseLightTheme", newMode); err != nil {
                showError("Failed to set AppsUseLightTheme: " + err.Error())
                return
        }

        // Provide feedback on the mode change
        systray.SetTooltip(messages.ModeSwitched)
        time.Sleep(2 * time.Second)
        systray.SetTooltip(messages.ToggleTooltip) // Reset tooltip
}

func toggleWindowsMode() {
        // This can be implemented as the same logic as system mode toggle
        // But assuming that there's a separate registry or method to handle this, you could extend this.
        // For now, it's the same as system mode toggle
        toggleSystemMode() // Reusing system mode toggle for simplicity
}

func toggleAutorun(autorunItem *systray.MenuItem) {
        // Open the registry key for Run settings
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
        systray.SetTooltip(messages.ToggleTooltip) // Reset tooltip
}

func removeAutorun(key registry.Key, autorunItem *systray.MenuItem) error {
        // Remove the app from the autorun registry
        if err := key.DeleteValue(appName); err != nil {
                showError("Failed to remove autorun: " + err.Error())
                return err
        }
        // Update the menu item label to indicate autorun is disabled
        autorunItem.SetTitle(menuTitles.EnableAutorunUnchecked)
        systray.SetTooltip(messages.AutorunDisabled)
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
        autorunItem.SetTitle(menuTitles.EnableAutorunChecked)
        systray.SetTooltip(messages.AutorunEnabled)
        return nil
}

func getExePath() (string, error) {
        // Return the path of the currently running executable
        return os.Executable()
}

func updateAutorunStatus(autorunItem *systray.MenuItem) {
        // Open the registry key for Run settings
        key, err := openRegistryKey(autorunRegistryKey, registry.QUERY_VALUE)
        if err != nil {
                showError("Failed to open autorun registry key: " + err.Error())
                return
        }
        defer key.Close()

        // Check if the app is already set to autorun and update the menu item label
        if isAutorunEnabled(key) {
                // App is set to autorun, update the menu item to reflect that
                autorunItem.SetTitle(menuTitles.EnableAutorunChecked)
        } else {
                // App is not set to autorun, update the menu item to reflect that
                autorunItem.SetTitle(menuTitles.EnableAutorunUnchecked)
        }
}

func isAutorunEnabled(key registry.Key) bool {
        _, _, err := key.GetStringValue(appName)
        return err == nil
}

func showError(message string) {
        logError(message) // Log to file for persistence
        systray.SetTooltip(messages.Error + message)
        time.Sleep(3 * time.Second)
        systray.SetTooltip(messages.ToggleTooltip) // Reset tooltip
}

func logError(message string) {
        file, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
                fmt.Println("Failed to open error log:", err)
                return
        }
        defer file.Close()
        logger := log.New(file, "", log.LstdFlags)
        logger.Println(message)
}

func openRegistryKey(path string, access uint32) (registry.Key, error) {
        key, err := registry.OpenKey(registry.CURRENT_USER, path, access)
        if err != nil {
                return registry.Key(0), fmt.Errorf("failed to open registry key: %v", err) // Return an empty Key (0) instead of nil
        }
        return key, nil
}
