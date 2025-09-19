// ui_texts.go

package main

var menuTitles = struct {
	AppName                string
	EnableAutorun          string
	Quit                   string
	EnableAutorunChecked   string
	EnableAutorunUnchecked string
	ToggleAppToDark        string
	ToggleAppToLight       string
	ToggleWinToDark        string
	ToggleWinToLight       string
	ToggleBothToDark       string
	ToggleBothToLight      string
	UpdateNow              string
}{
	AppName:                "🔗 WinGoDarkTray",
	EnableAutorun:          "Enable Autorun",
	Quit:                   "✕ Quit",
	EnableAutorunChecked:   "✔ Autorun Enabled",
	EnableAutorunUnchecked: "✗ Autorun Disabled",
	ToggleAppToDark:        "☾ Toggle app theme to Dark",
	ToggleAppToLight:       "☼ Toggle app theme to Light",
	ToggleWinToDark:        "☾ Toggle Windows theme to Dark",
	ToggleWinToLight:       "☼ Toggle Windows theme to Light",
	ToggleBothToDark:       "☾ Toggle both to Dark",
	ToggleBothToLight:      "☼ Toggle both to Light",
	UpdateNow:              "🔄 Update Now",
}

var tooltips = struct {
	Default string
	Error   string
}{
	Default: "A windows app to toggle light and dark mode from the system tray",
	Error:   "An error occurred, please try again later.",
}

var notificationTexts = struct {
	Error                  string
	UpdateAvailableTitle   string
	UpdateAvailableMessage string
}{
	Error:                  "WinGoDarkTray Error",
	UpdateAvailableTitle:   "WinGoDarkTray Update",
	UpdateAvailableMessage: "Update available: ",
}
