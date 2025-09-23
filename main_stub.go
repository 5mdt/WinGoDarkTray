// main_stub.go

//go:build !windows
// +build !windows

package main

import (
	"fmt"
	"os"
)

// Stub main function for non-Windows platforms
// This allows basic compilation validation on non-Windows CI runners
func main() {
	fmt.Fprintln(os.Stderr, "WinGoDarkTray is only supported on Windows platforms")
	os.Exit(1)
}
