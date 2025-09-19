# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

WinGoDarkTray is a Go-based Windows system tray application that allows users to toggle between light and dark themes for Windows UI, apps, and system-wide settings. The application is built for Windows platforms using cross-compilation.

**Important**: All features and requirements are defined in the [Functional Requirements Document (FRD.md)](./FRD.md). Refer to this document when making changes to ensure compliance with project specifications.

## Development Commands

### Building
- **Build for all platforms**: `./build.sh`
- **Build for development**: `go build -o WinGoDarkTray.exe`
- **Format code**: `go fmt ./...`
- **Install dependencies**: `go mod tidy`

### Testing
- **Run all tests**: `go test -v -cover ./...`
- **Run with coverage**: `go test -v -race -coverprofile=coverage.out -covermode=atomic ./...`
- **Generate coverage report**: `go tool cover -html=coverage.out -o coverage.html`
- **Run specific test**: `go test -run TestSemver`
- **Run linting**: `golangci-lint run`

**Note**: Full test suite requires Windows environment due to Windows API dependencies. See `TESTING.md` for comprehensive testing strategy.\n\n### CI/CD Workflows\nThe project includes comprehensive GitHub Actions workflows:\n- **Test Suite**: Windows integration testing and cross-platform validation\n- **Build and Publish**: Multi-architecture builds with security scanning\n- **Security**: Automated vulnerability and static analysis scanning\n- **Dependency Review**: License and security review for PRs

### Pre-commit hooks
The project uses pre-commit hooks configured in `.pre-commit-config.yaml`. Install with:
```bash
pre-commit install
```

## Architecture Overview

### Core Components

- **main.go**: Application entry point and systray initialization using `github.com/getlantern/systray`
- **theme.go**: Windows registry manipulation for theme switching via `golang.org/x/sys/windows/registry`
- **tray_handlers.go**: Event handlers for systray menu interactions
- **autorun.go**: Windows autorun functionality management
- **updater.go**: Self-update mechanism checking GitHub releases
- **logging.go**: Windows Event Log integration
- **semver.go**: Semantic version parsing and comparison

### Key Dependencies

- `github.com/getlantern/systray`: System tray interface
- `golang.org/x/sys/windows`: Windows API access
- `github.com/gen2brain/beeep`: Cross-platform notifications

### Build System

The `build.sh` script handles cross-compilation for Windows platforms:
- Targets: x64 (amd64), x32 (386), ARM64
- Uses `rsrc` tool to embed icon resources
- Applies build flags: `-trimpath`, `-buildvcs=false`, `-ldflags` for version injection
- Outputs to `./build/` directory

### Registry Integration

Theme switching works through Windows registry manipulation:
- Path: `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`
- Keys: `AppsUseLightTheme`, `SystemUsesLightTheme`
- Values: 0 (dark), 1 (light)

### Version Management

Version is injected at build time via ldflags (`-X main.buildVersion=${VERSION}`). The build script extracts version from git tags using `git describe --tags --always --dirty`.

### Event Logging

The application integrates with Windows Event Log for monitoring and debugging. Log events are written to the Windows Event Log with different severity levels.

## Build Constraints

The project uses Go build constraints for cross-platform compatibility:
- **Windows-specific files**: Tagged with `//go:build windows` (main.go, theme.go, etc.)
- **Cross-platform stub**: `main_stub.go` tagged with `//go:build !windows`
- **Tests**: Most test files work cross-platform, but full functionality requires Windows

## File Structure Notes

- All Go source files are in the root directory (flat structure)
- `main.go`: Windows-only main application entry point
- `main_stub.go`: Cross-platform stub for non-Windows validation
- `icon.ico`: Application icon embedded during build
- `build/`: Build artifacts directory
- `_assets/`: Documentation assets (screenshots)
