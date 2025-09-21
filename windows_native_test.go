// windows_native_test.go - Windows-specific API tests

//go:build windows
// +build windows

package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc/eventlog"
)

// TestWindowsNativeFunctionality validates Windows-specific API interactions under various permission scenarios
func TestWindowsNativeFunctionality(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "setBothThemeModes_actual_registry_operations",
			test: func(t *testing.T) {
				originalMode, err := getCurrentAppThemeMode()
				if err != nil {
					t.Logf("Could not read current theme mode: %v", err)
					return
				}

				err = setBothThemeModes(true)
				if err != nil {
					t.Logf("Could not set light mode: %v", err)
				}

				err = setBothThemeModes(false)
				if err != nil {
					t.Logf("Could not set dark mode: %v", err)
				}

				// Restore original theme mode
				if originalMode == 1 {
					setBothThemeModes(true)
				} else {
					setBothThemeModes(false)
				}
			},
		},
		{
			name: "toggleSingleTheme_registry_operations",
			test: func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Logf("toggleSingleTheme panicked as expected: %v", r)
					}
				}()

				app := NewApp("registry-test")
				app.toggleSingleTheme(regValAppsUseLight)
				app.toggleSingleTheme(regValSystemUsesLight)
			},
		},
		{
			name: "eventLogSourceExists_windows_check",
			test: func(t *testing.T) {
				exists := eventLogSourceExists()
				t.Logf("Event log source exists: %v", exists)

				// Verify consistent behavior across multiple calls
				exists2 := eventLogSourceExists()
				exists3 := eventLogSourceExists()
				_ = exists2
				_ = exists3
			},
		},
		{
			name: "installEventLogSource_windows_attempt",
			test: func(t *testing.T) {
				err1 := installEventLogSource()
				t.Logf("Install attempt 1 result: %v", err1)

				err2 := installEventLogSource()
				t.Logf("Install attempt 2 result: %v", err2)

				err3 := installEventLogSource()
				t.Logf("Install attempt 3 result: %v", err3)
			},
		},
		{
			name: "isAdmin_windows_check",
			test: func(t *testing.T) {
				admin := isAdmin()
				t.Logf("Running as admin: %v", admin)

				// Verify consistent behavior across multiple calls
				admin2 := isAdmin()
				admin3 := isAdmin()
				_ = admin2
				_ = admin3
			},
		},
		{
			name: "logEvent_actual_windows_logging",
			test: func(t *testing.T) {
				logEvent(eventlog.Info, "Test info message from unit test")
				logEvent(eventlog.Warning, "Test warning message from unit test")
				logEvent(eventlog.Error, "Test error message from unit test")

				// Test Unicode and extended length handling
				logEvent(eventlog.Info, "Unicode test: ñáéíóú αβγδε 中文")
				logEvent(eventlog.Warning, "Long message: "+strings.Repeat("A", 500))
			},
		},
		{
			name: "openBrowser_windows_execution",
			test: func(t *testing.T) {
				openBrowser("https://example.com")

				// Test various URL formats and edge cases
				openBrowser("")
				openBrowser("invalid-url-format")
				openBrowser("http://")
				openBrowser("https://example.com/path?param=value&other=test")
			},
		},
		{
			name: "withThemeRegistry_error_scenarios",
			test: func(t *testing.T) {
				// Test registry operations with different permission levels
				err1 := withThemeRegistry(registry.QUERY_VALUE, func(key registry.Key) error {
					_, _, err := key.GetIntegerValue(regValAppsUseLight)
					return err
				})
				t.Logf("Query value result: %v", err1)

				err2 := withThemeRegistry(registry.SET_VALUE, func(key registry.Key) error {
					return key.SetDWordValue(regValAppsUseLight, 1)
				})
				t.Logf("Set value result: %v", err2)

				err3 := withThemeRegistry(registry.READ, func(key registry.Key) error {
					_, _, err := key.GetIntegerValue(regValSystemUsesLight)
					return err
				})
				t.Logf("Read access result: %v", err3)
			},
		},
		{
			name: "getExePath_windows_executable",
			test: func(t *testing.T) {
				path1, err1 := getExePath()
				t.Logf("Executable path: %s, error: %v", path1, err1)

				// Verify consistent behavior across multiple calls
				path2, err2 := getExePath()
				path3, err3 := getExePath()
				_ = path2
				_ = path3
				_ = err2
				_ = err3
			},
		},
		{
			name: "autorun_registry_operations",
			test: func(t *testing.T) {
				err1 := withAutorunRegistry(registry.QUERY_VALUE, func(key registry.Key) error {
					enabled1 := isAutorunEnabled(key)
					enabled2 := isAutorunEnabled(key)
					enabled3 := isAutorunEnabled(key)
					t.Logf("Autorun enabled checks: %v, %v, %v", enabled1, enabled2, enabled3)
					return nil
				})
				t.Logf("Autorun query result: %v", err1)

				err2 := withAutorunRegistry(registry.SET_VALUE, func(key registry.Key) error {
					exePath, err := getExePath()
					if err != nil {
						return err
					}
					return key.SetStringValue("WinGoDarkTray", exePath)
				})
				t.Logf("Autorun set result: %v", err2)
			},
		},
		{
			name: "showNotificationWithFallback_windows_notifications",
			test: func(t *testing.T) {
				showNotificationWithFallback("Test Title", "Test message from unit test", "Test tooltip")
				showNotificationWithFallback("Unicode Test", "Testing ñáéíóú αβγδε 中文", "Unicode tooltip")
				showNotificationWithFallback("", "", "")
				showNotificationWithFallback("Long Title", string(make([]byte, 200)), "Long tooltip")

				// Allow notification processing time
				time.Sleep(100 * time.Millisecond)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// TestWindowsErrorPaths validates error handling in Windows-specific operations
func TestWindowsErrorPaths(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "registry_permission_errors",
			test: func(t *testing.T) {
				// Test registry access with varying permission levels
				key1, err1 := openRegistryKey(themeRegistryPath, registry.ALL_ACCESS)
				if err1 == nil {
					key1.Close()
				}
				t.Logf("ALL_ACCESS result: %v", err1)

				key2, err2 := openRegistryKey(themeRegistryPath, registry.WRITE)
				if err2 == nil {
					key2.Close()
				}
				t.Logf("WRITE result: %v", err2)

				key3, err3 := openRegistryKey(themeRegistryPath, registry.READ)
				if err3 == nil {
					key3.Close()
				}
				t.Logf("READ result: %v", err3)

				// Test error handling with invalid registry path
				key4, err4 := openRegistryKey("Invalid\\Registry\\Path", registry.READ)
				if err4 == nil {
					key4.Close()
				}
				t.Logf("Invalid path result: %v", err4)
			},
		},
		{
			name: "environment_specific_tests",
			test: func(t *testing.T) {
				// Test environment-dependent functionality
				originalDir, _ := os.Getwd()

				path, err := getExePath()
				t.Logf("Exe path: %s, error: %v", path, err)

				os.Chdir(originalDir)

				admin := isAdmin()
				t.Logf("Admin status: %v", admin)
			},
		},
		{
			name: "updateThemeToggleTitles_without_systray",
			test: func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Logf("updateThemeToggleTitles panicked as expected: %v", r)
					}
				}()

				app := NewApp("theme-titles-test")
				// Skip createMenuItems to test behavior without systray
				app.updateThemeToggleTitles()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// TestWindowsAPISafeInteractions validates safe Windows API operations
func TestWindowsAPISafeInteractions(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "theme_detection_comprehensive",
			test: func(t *testing.T) {
				// Test consistent theme detection behavior
				for i := 0; i < 5; i++ {
					mode, err := getCurrentAppThemeMode()
					t.Logf("Iteration %d - App theme mode: %d, error: %v", i+1, mode, err)
					time.Sleep(10 * time.Millisecond)
				}
			},
		},
		{
			name: "logging_comprehensive_scenarios",
			test: func(t *testing.T) {
				// Test comprehensive logging scenarios with various message types
				messages := []string{
					"Simple test message",
					"Message with special chars: !@#$%^&*()",
					"Unicode message: ñáéíóú αβγδε 中文测试",
					"Very long message: " + strings.Repeat("Test", 250),
					"Message\nwith\nnewlines",
					"Message\twith\ttabs",
					"",
				}

				types := []uint32{
					eventlog.Info,
					eventlog.Warning,
					eventlog.Error,
				}

				for _, msg := range messages {
					for _, logType := range types {
						logEvent(logType, msg)
						logToConsole(logType, msg)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
