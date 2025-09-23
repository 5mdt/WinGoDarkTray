// utils_stub.go

//go:build !windows
// +build !windows

package main

import "fmt"

// Stub functions for non-Windows platforms to enable testing

func getExePath() (string, error) {
	return "", fmt.Errorf("getExePath not supported on non-Windows platforms")
}

func openBrowser(url string) error {
	return fmt.Errorf("openBrowser not supported on non-Windows platforms")
}
