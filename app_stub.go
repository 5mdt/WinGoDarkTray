// app_stub.go

//go:build !windows
// +build !windows

package main

import "fmt"

// App represents the application state (stub for non-Windows)
type App struct{}

// NewApp creates a new App instance (stub for non-Windows)
func NewApp() *App {
	return &App{}
}

// Run starts the application (stub for non-Windows)
func (a *App) Run() error {
	return fmt.Errorf("App.Run not supported on non-Windows platforms")
}
