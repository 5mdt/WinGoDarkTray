#!/usr/bin/env bash
set -euo pipefail

# Configuration
readonly APP_NAME="WinGoDarkTray"
readonly ICON_FILE="icon.ico"
readonly BUILD_DIR="./build"
readonly RSRC_VERSION="v0.10.2"
readonly LDFLAGS="-s -w -H=windowsgui"  # -s: strip symbol table, -w: strip debug info, -H: hide console window
readonly BUILD_OPTS="-trimpath -buildvcs=false"  # -trimpath: remove file paths, -buildvcs: disable VCS stamping

# Architecture targets: GOARCH:suffix
readonly TARGETS=(
  "amd64:x64"
  "386:x32"
  "arm64:arm64"
)

# Colors (only if terminal supports them)
if [[ -t 1 ]]; then
  readonly RED='\033[0;31m'
  readonly GREEN='\033[0;32m'
  readonly BLUE='\033[0;34m'
  readonly YELLOW='\033[1;33m'
  readonly CYAN='\033[0;36m'
  readonly BOLD='\033[1m'
  readonly RESET='\033[0m'
else
  readonly RED='' GREEN='' BLUE='' YELLOW='' CYAN='' BOLD='' RESET=''
fi

# Logging functions
log_info() { echo -e "${CYAN}→${RESET} $*"; }
log_success() { echo -e "${GREEN}✓${RESET} $*"; }
log_error() { echo -e "${RED}✗ Error:${RESET} $*" >&2; }
log_build() { echo -e "${YELLOW}🔨${RESET} $*"; }
log_header() { echo -e "${BOLD}${BLUE}$*${RESET}"; }

# Validation
validate_requirements() {
  log_info "🔍 Validating requirements..."

  [[ -f "$ICON_FILE" ]] || { log_error "Icon file '$ICON_FILE' not found"; exit 1; }

  command -v go >/dev/null || { log_error "Go not installed"; exit 1; }

  if ! command -v rsrc >/dev/null; then
    log_info "📦 Installing rsrc tool (${RSRC_VERSION})..."
    go install "github.com/akavel/rsrc@${RSRC_VERSION}" || { log_error "Failed to install rsrc"; exit 1; }
  fi

  log_success "All requirements validated"
}

# Clean and prepare
prepare_build() {
  log_info "🧹 Preparing build environment..."

  # Remove old build artifacts
  find . -name "*.exe" -o -name "*.syso" | xargs -r rm -f

  # Create build directory
  mkdir -p "$BUILD_DIR"

  # Format code
  go fmt ./...

  log_success "Build environment ready"
}

# Build for all targets
build_targets() {
  for target in "${TARGETS[@]}"; do
    local arch="${target%:*}"
    local suffix="${target#*:}"
    local output="$BUILD_DIR/$APP_NAME-$suffix.exe"

    log_build "Building for $arch ($suffix)..."

    # Generate arch-specific .syso file
    log_info "🖼️  Generating $arch-specific icon resource..."
    GOARCH="$arch" GOOS="windows" rsrc -ico "$ICON_FILE" \
      || { log_error "Failed to embed icon for $arch"; exit 1; }

    # Build for target architecture
    GOARCH="$arch" GOOS="windows" go build $BUILD_OPTS -ldflags="$LDFLAGS" -o "$output" \
      || { log_error "Build failed for $arch"; exit 1; }

    # Clean up arch-specific .syso immediately
    find . -name "*.syso" -delete

    log_success "Built: $(basename "$output")"
  done
}

# Main execution
main() {
  log_header "🚀 Starting build process for $APP_NAME"

  validate_requirements
  prepare_build
  build_targets

  log_header "🎉 Build completed successfully!"
  log_info "📦 Artifacts created in: $BUILD_DIR/"
  if ls "$BUILD_DIR"/*.exe >/dev/null 2>&1; then
    echo -e "${BOLD}${GREEN}$(ls -la "$BUILD_DIR"/*.exe)${RESET}"
  fi
}

# Execute if run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  main "$@"
fi
