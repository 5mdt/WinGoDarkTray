package main

import (
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/sys/windows/registry"
)

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

func openRegistryKey(path string, access uint32) (registry.Key, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, path, access)
	if err != nil {
		return registry.Key(0), err
	}
	return key, nil
}

func getExePath() (string, error) {
	return os.Executable()
}
