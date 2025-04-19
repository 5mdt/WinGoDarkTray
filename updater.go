package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

const githubReleasesAPI = "https://api.github.com/repos/5mdt/WinGoDarkTray/releases/latest"

func checkForUpdate(currentVersion string) {

	resp, err := http.Get(githubReleasesAPI)
	if err != nil {
		logEvent(3, "Failed to fetch latest release: "+err.Error())
		return
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
	}

	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		logEvent(3, "Failed to parse GitHub release response: "+err.Error())
		return
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	current := strings.TrimPrefix(currentVersion, "v")

	if isVersionNewer(latest, current) {
		logEvent(1, fmt.Sprintf("New version available: %s (current: %s)", latest, current))
		message := fmt.Sprintf("New version %s is available!.", release.TagName)

		err := beeep.Notify(notificationTexts.UpdateAvailableTitle, message, "")
		if err != nil {
			systray.SetTooltip(notificationTexts.UpdateAvailableMessage + release.TagName)
		}

		updateNowItem.Show()

		go func() {
			for {
				select {
				case <-updateNowItem.ClickedCh:
					runUpdateCommand()
				}
			}
		}()
	}
}

func runUpdateCommand() {

	cmd := exec.Command("cmd", "/C", "start", "cmd", "/K", "echo Trying to update WinGoDarkTray, using winget app && winget install 5mdt.WinGoDarkTray && echo If nothing installed, please try again later. Winget repository moderators must approve new package first.")
	err := cmd.Start()
	if err != nil {
		showError("Failed to run update command: " + err.Error())
	}
}
