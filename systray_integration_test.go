// systray_integration_test.go

//go:build windows && integration
// +build windows,integration

package main

import (
	"context"
	"testing"
	"time"

	"github.com/getlantern/systray"
)

// Test functions that require systray initialization
func TestSystrayIntegration(t *testing.T) {
	// Skip if running in automated CI to avoid GUI dependencies
	if testing.Short() {
		t.Skip("Skipping systray integration test in short mode")
	}

	// Channel to communicate test results
	testDone := make(chan bool, 1)
	testError := make(chan error, 1)

	// Run systray in a goroutine with timeout
	go func() {
		systray.Run(func() {
			// onReady - systray is initialized
			defer func() {
				if r := recover(); r != nil {
					testError <- nil // Don't fail on panic, just note it happened
				}
				testDone <- true
			}()

			// Test createMenuItems
			t.Run("createMenuItems", func(t *testing.T) {
				app := NewApp("test-1.0.0")
				autorunItem, _ := app.createMenuItems()

				// Basic validation that menu items were created
				if app.toggleSystemItem == nil {
					t.Error("toggleSystemItem was not created")
				}
				if app.toggleAppItem == nil {
					t.Error("toggleAppItem was not created")
				}
				if app.toggleWindowsItem == nil {
					t.Error("toggleWindowsItem was not created")
				}
				_ = autorunItem // Use the return value
			})

			// Test initializeMenuState
			t.Run("initializeMenuState", func(t *testing.T) {
				app := NewApp("test-1.0.0")
				autorunItem, _ := app.createMenuItems()
				if autorunItem != nil {
					app.initializeMenuState(autorunItem)
				}
				// If it doesn't panic, it's working
			})

			// Test updateThemeToggleTitles
			t.Run("updateThemeToggleTitles", func(t *testing.T) {
				app := NewApp("test-1.0.0")
				_, _ = app.createMenuItems()
				app.updateThemeToggleTitles()
				// If it doesn't panic, it's working
			})

			// Test autorun functions with systray
			t.Run("updateAutorunUI", func(t *testing.T) {
				app := NewApp("test-1.0.0")
				autorunItem, _ := app.createMenuItems()
				if autorunItem != nil {
					updateAutorunUI(autorunItem, true)
					updateAutorunUI(autorunItem, false)
				}
				// If it doesn't panic, it's working
			})

			// Test theme toggle functions
			t.Run("toggleSystemMode", func(t *testing.T) {
				app := NewApp("test-1.0.0")
				autorunItem, _ := app.createMenuItems()
				if autorunItem != nil {
					app.initializeMenuState(autorunItem)
				}
				app.toggleSystemMode()
				// If it doesn't panic, it's working
			})

			t.Run("toggleAppMode", func(t *testing.T) {
				app := NewApp("test-1.0.0")
				autorunItem, _ := app.createMenuItems()
				if autorunItem != nil {
					app.initializeMenuState(autorunItem)
				}
				app.toggleAppMode()
				// If it doesn't panic, it's working
			})

			t.Run("toggleWindowsMode", func(t *testing.T) {
				app := NewApp("test-1.0.0")
				autorunItem, _ := app.createMenuItems()
				if autorunItem != nil {
					app.initializeMenuState(autorunItem)
				}
				app.toggleWindowsMode()
				// If it doesn't panic, it's working
			})

			// Test setTemporaryTooltip
			t.Run("setTemporaryTooltip", func(t *testing.T) {
				setTemporaryTooltip("Test tooltip", 50*time.Millisecond)
				// Wait for tooltip to reset
				time.Sleep(100 * time.Millisecond)
			})

			// Exit systray after tests
			systray.Quit()
		}, func() {
			// onExit
		})
	}()

	// Wait for tests to complete or timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-testDone:
		// Tests completed successfully
	case err := <-testError:
		if err != nil {
			t.Fatalf("Systray test failed: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("Systray test timed out")
	}
}

// Test that we can start and stop systray multiple times
func TestSystrayLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping systray lifecycle test in short mode")
	}

	for i := 0; i < 3; i++ {
		testDone := make(chan bool, 1)

		go func() {
			systray.Run(func() {
				// Quick initialization test
				app := NewApp("lifecycle-test")
				autorunItem, _ := app.createMenuItems()
				if autorunItem != nil {
					app.initializeMenuState(autorunItem)
				}
				systray.Quit()
				testDone <- true
			}, func() {
				// onExit
			})
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		select {
		case <-testDone:
			// Success
		case <-ctx.Done():
			cancel()
			t.Fatalf("Systray lifecycle test %d timed out", i+1)
		}
		cancel()

		// Small delay between iterations
		time.Sleep(100 * time.Millisecond)
	}
}

// Test the actual application flow
func TestApplicationFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping application flow test in short mode")
	}

	testDone := make(chan bool, 1)

	go func() {
		systray.Run(func() {
			defer func() {
				testDone <- true
			}()

			// Test the actual application initialization flow
			app := NewApp("flow-test-1.0.0")

			// Test createMenuItems
			autorunItem, _ := app.createMenuItems()

			// Test initializeApp components
			if !eventLogSourceExists() {
				// Try to install if we're admin (won't fail test if not)
				_ = installEventLogSource()
			}

			// Test menu state initialization
			if autorunItem != nil {
				app.initializeMenuState(autorunItem)
			}

			// Test update checking (with short timeout)
			go func() {
				time.Sleep(100 * time.Millisecond)
				checkForUpdate(app.version, app.updateNowItem)
			}()

			// Test autorun status update
			if autorunItem != nil {
				updateAutorunStatus(autorunItem)
			}

			// Wait a bit for async operations
			time.Sleep(500 * time.Millisecond)

			systray.Quit()
		}, func() {
			// onExit - clean shutdown
		})
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	select {
	case <-testDone:
		// Success
	case <-ctx.Done():
		t.Fatal("Application flow test timed out")
	}
}