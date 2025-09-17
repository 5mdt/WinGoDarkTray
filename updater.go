package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/svc/eventlog"
)

const githubReleasesAPI = "https://api.github.com/repos/5mdt/WinGoDarkTray/releases/latest"

// startUpdateClickHandler starts a goroutine to handle update button clicks
func startUpdateClickHandler(updateNowItem *systray.MenuItem, quitCh chan struct{}) {
	go func() {
		for {
			select {
			case <-updateNowItem.ClickedCh:
				runUpdateCommand()
			case <-quitCh:
				return
			}
		}
	}()
}

func checkForUpdate(currentVersion string, updateNowItem *systray.MenuItem) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(githubReleasesAPI)
	if err != nil {
		logEvent(eventlog.Error, "Failed to fetch latest release: "+err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logEvent(eventlog.Error, fmt.Sprintf("GitHub releases API returned %d", resp.StatusCode))
		return
	}
	var release struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
	}

	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		logEvent(eventlog.Error, "Failed to parse GitHub release response: "+err.Error())
		return
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	current := strings.TrimPrefix(currentVersion, "v")

	if isVersionNewer(latest, current) {
		logEvent(eventlog.Info, fmt.Sprintf("New version available: %s (current: %s)", latest, current))
		showUpdateNotification(release.TagName)
		updateNowItem.Show()
	}
}

// showUpdateNotification displays update available notification
func showUpdateNotification(tagName string) {
	message := fmt.Sprintf("New version %s is available!.", tagName)

	err := beeep.Notify(notificationTexts.UpdateAvailableTitle, message, "")
	if err != nil {
		systray.SetTooltip(notificationTexts.UpdateAvailableMessage + tagName)
	}
}

func runUpdateCommand() {
	cmd := exec.Command("cmd", "/C", "start", "cmd", "/K", "echo Trying to update WinGoDarkTray, using winget app && winget install 5mdt.WinGoDarkTray && echo If nothing installed, please try again later. Winget repository moderators must approve new package first.")
	err := cmd.Start()
	if err != nil {
		showError("Failed to run update command: " + err.Error())
	}
}
