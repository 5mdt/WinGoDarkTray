# 📝 Functionality Requirements Document

## ℹ️ Project Information

- **Name:** WinGoDarkTray
- **Platform:** Windows-only
- **Tested On:** Windows 10 and above
- **Goal:** A lightweight system tray utility to toggle Windows light/dark themes and autorun setting.
- **Philosophy:** Simple, dependency-light, and efficient (KISS). Code reuse preferred where applicable (DRY).

-

## ✅ Functional Requirements

| ID  | Feature  | Description  | User Story | Expected Behavior  |
|-|-|-|-|-|
| **F-001** | App Theme Toggle | Toggles light/dark mode for Windows apps via tray. Displays the target state (e.g., "Toggle app theme to Dark" or "Toggle app theme to Light") on the menu. | As a user, I want to easily see and instantly toggle the app theme between light and dark modes with a single click. | Modifies the `AppsUseLightTheme` registry key. The tray button displays the current target state, e.g., "Toggle app theme to Dark" if in Light mode, or "Toggle app theme to Light" if in Dark mode. A tooltip briefly shows "App theme changed to Dark" or "App theme changed to Light" for immediate feedback. If the registry key modification fails, an error icon will appear in the system tray with a tooltip: "Failed to toggle app theme." |
| **F-002** | Windows Theme Toggle | Toggles light/dark mode for system UI (taskbar, Start menu). Displays the target state (e.g., "Toggle Windows theme to Dark" or "Toggle Windows theme to Light") on the menu. | As a user, I want to easily see and instantly toggle the system theme between light and dark modes with a single click. | Modifies the `SystemUsesLightTheme` registry key. The tray button displays the current target state, e.g., "Toggle Windows theme to Dark" if in Light mode, or "Toggle Windows theme to Light" if in Dark mode. A tooltip briefly shows "System theme changed to Dark" or "System theme changed to Light" for immediate feedback. If the registry key modification fails, an error icon will appear in the system tray with a tooltip: "Failed to toggle system theme." |
| **F-003** | System Theme Toggle (App + Windows) | Toggles both themes (app and system) together. The menu item is labeled as "Toggle both to Light/Dark" and updates based on the next target state. If the app theme and Windows theme differ when the button is pressed, both will be synchronized to the toggled app theme. Displays the target state (e.g., "Toggle both to Dark" or "Toggle both to Light") on the menu. | As a user, I want to change both app and system themes in one click and have both themes synchronized to the same state. | If app and system themes differ, toggles both to the new app theme. The tray item shows the current target state for both themes. A tooltip shows "System theme changed to Light/Dark" after toggle, e.g., "Toggle both to Dark" or "Toggle both to Light." |
| **F-004** | Autorun Control | Enables/disables app auto-launch on login. | As a user, I want the app to start automatically if I choose. | Adds/removes autorun key for current executable. Menu shows ✅/❌. |
| **F-005** | Tray UI | Displays a system tray icon and menu. | As a user, I want easy access without launching a full app window. | Menu includes: `Toggle System`, `App`, `Windows`, `Autorun`, `GitHub`, `Quit`. |
| **F-006** | Open GitHub | Opens the project’s GitHub repo from the tray. | As a curious user or dev, I want to see the source or contribute. | Clicking “WinGoDarkTray 🔗” opens the GitHub link in the default browser. |
| **F-007** | Error Handling & Windows Event Logging | Uses Windows Event Log to record operational or error messages. | As a user, I want to know when something fails. As a developer, I want standard logs for troubleshooting. | Errors and important events are logged using the built-in Windows Event Log system under a custom source (`WinGoDarkTray`). Tooltip shows a temporary friendly error message. Errors and important events are logged using the built-in Windows Event Log system under a custom source (WinGoDarkTray). The log source is silently installed on first run.|
| **F-008** | Version Update Checker | On app start, checks for newer version from GitHub. Notifies user via Windows notification. | As a user, I want to know if there's a newer version available. | Calls GitHub API, compares current version. If newer version found, uses beeep.Notify() to alert. Menu item allows manual check. |
-

## 🔧 Technical Constraints

- **Language:** Go (Golang)
- **Dependencies:**
- [`github.com/getlantern/systray`](https://github.com/getlantern/systray) for tray icon/menu.
- [`golang.org/x/sys/windows/registry`](https://pkg.go.dev/golang.org/x/sys/windows/registry) for accessing Windows registry.
- **No GUI windows**, config files, or background services.
- All settings are persisted via **Windows Registry** only.
- App is designed for **simplicity, minimalism, and zero config**.
- **Performance Considerations:** App should start quickly with minimal CPU/memory usage to ensure a seamless experience for the user.
