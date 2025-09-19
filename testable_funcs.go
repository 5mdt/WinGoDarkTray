//go:build !windows
// +build !windows

package main

// Stub implementations for testing on non-Windows platforms
// These allow tests to compile and run basic logic tests

import (
	"errors"
	"fmt"
)

// Mock Windows registry key for testing
type mockRegistryKey struct{}

func (m mockRegistryKey) Close() error { return nil }

// Stub implementations that don't use Windows APIs
func showError(message string) {
	fmt.Printf("Error: %s\n", message)
}

func logEvent(level interface{}, message string) {
	fmt.Printf("Log: %s\n", message)
}

func setTemporaryTooltip(message string, duration interface{}) {
	fmt.Printf("Tooltip: %s\n", message)
}

func openRegistryKey(path string, access uint32) (interface{}, error) {
	return mockRegistryKey{}, errors.New("registry not available on non-Windows")
}

// Stub systray functions for compilation
type mockMenuItem struct{}

func (m *mockMenuItem) SetTitle(title string) {}

type mockApp struct {
	version           string
	toggleSystemItem  *mockMenuItem
	toggleAppItem     *mockMenuItem
	toggleWindowsItem *mockMenuItem
	updateNowItem     *mockMenuItem
}

func NewAppForTesting(version string) *mockApp {
	if version == "" {
		version = "v0.0.0"
	}
	return &mockApp{
		version: version,
	}
}
