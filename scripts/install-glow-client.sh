#!/bin/bash
set -e

# Glow Client-Only Installation Script
# Installs only the glow CLI from GitHub Releases
# Does NOT install or modify glow-server

# Configuration
REPO="luaxlou/glow"
INSTALL_DIR="${HOME}/.local/bin"
DATA_DIR=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    OS="$(uname -s)"
    ARCH="$(uname -m)"

    case "$OS" in
        Linux)
            OS="linux"
            ;;
        Darwin)
            OS="darwin"
            ;;
        *)
            log_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac

    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac

    log_info "Detected platform: $OS-$ARCH"
}

# Get latest release version
get_latest_version() {
    log_info "Fetching latest release version..."
    VERSION=$(curl -s https://api.github.com/repos/${REPO}/releases/latest | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/' | sed 's/^v//')

    if [ -z "$VERSION" ]; then
        log_error "Failed to fetch latest version"
        exit 1
    fi

    log_info "Latest version: $VERSION"
}

# Download and verify binary
download_binary() {
    local binary_name=$1
    local output_path=$2

    local filename="${binary_name}-${OS}-${ARCH}"
    local download_url="https://github.com/${REPO}/releases/download/v${VERSION}/${filename}"
    local checksum_url="${download_url}.sha256"

    log_info "Downloading ${binary_name} from ${download_url}..."

    # Download binary
    if ! curl -fSL -o "${output_path}" "${download_url}"; then
        log_error "Failed to download ${binary_name}"
        exit 1
    fi

    # Download and verify checksum
    log_info "Verifying checksum..."
    if curl -fSL -o "${output_path}.sha256" "${checksum_url}"; then
        if command -v sha256sum &> /dev/null; then
            downloaded_checksum=$(sha256sum "${output_path}" | awk '{print $1}')
        elif command -v shasum &> /dev/null; then
            downloaded_checksum=$(shasum -a 256 "${output_path}" | awk '{print $1}')
        else
            log_warn "No checksum tool available, skipping verification"
            rm -f "${output_path}.sha256"
            return
        fi

        expected_checksum=$(cat "${output_path}.sha256" | awk '{print $1}')

        if [ "$downloaded_checksum" != "$expected_checksum" ]; then
            log_error "Checksum verification failed!"
            log_error "Expected: $expected_checksum"
            log_error "Got: $downloaded_checksum"
            rm -f "${output_path}" "${output_path}.sha256"
            exit 1
        fi

        log_info "Checksum verified successfully"
        rm -f "${output_path}.sha256"
    else
        log_warn "Failed to download checksum, skipping verification"
    fi

    # Make executable
    chmod +x "${output_path}"
}

# Install binary
install_binary() {
    log_info "Installing glow CLI to ${INSTALL_DIR}..."

    # Create install directory if it doesn't exist
    mkdir -p "${INSTALL_DIR}"

    # Check if binary already exists and create backup
    if [ -f "${INSTALL_DIR}/glow" ]; then
        BACKUP_FILE="${INSTALL_DIR}/glow.backup.$(date +%Y%m%d-%H%M%S)"
        cp "${INSTALL_DIR}/glow" "${BACKUP_FILE}"
        log_info "Backed up existing glow to ${BACKUP_FILE}"
    fi

    # Download and install glow
    glow_tmp="/tmp/glow-${VERSION}"
    download_binary "glow" "${glow_tmp}"
    mv "${glow_tmp}" "${INSTALL_DIR}/glow"
    log_info "Installed glow"
}

# Check PATH configuration
check_path() {
    if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
        echo ""
        log_warn "WARNING: ${INSTALL_DIR} is not in your PATH"
        echo ""
        echo "To add glow to your PATH, run one of the following commands:"
        echo ""
        echo "  For Bash:"
        echo "    echo 'export PATH=\"\$PATH:${INSTALL_DIR}\"' >> ~/.bashrc"
        echo "    source ~/.bashrc"
        echo ""
        echo "  For Zsh:"
        echo "    echo 'export PATH=\"\$PATH:${INSTALL_DIR}\"' >> ~/.zshrc"
        echo "    source ~/.zshrc"
        echo ""
        echo "  For temporary (current session only):"
        echo "    export PATH=\"\$PATH:${INSTALL_DIR}\""
        echo ""
    fi
}

# Print summary
print_summary() {
    echo ""
    echo "=========================================="
    log_info "Installation completed successfully!"
    echo "=========================================="
    echo ""
    echo "Installed binary:"
    echo "  - glow: ${INSTALL_DIR}/glow"
    echo ""
    echo "Quick start:"
    echo "  1. Add glow to your PATH (see warning above if needed)"
    echo "  2. Configure a connection context:"
    echo "     glow config add myserver --server-url=http://your-server:32102 --api-key=your-key"
    echo "  3. Set as default context:"
    echo "     glow config use myserver"
    echo "  4. Test connection:"
    echo "     glow get apps"
    echo ""
    echo "For more information, see: https://github.com/${REPO}"
    echo ""
}

# Main installation flow
main() {
    echo ""
    echo "Glow Client-Only Installation Script"
    echo "===================================="
    echo ""
    echo "This script installs ONLY the glow CLI."
    echo "It does NOT install or modify glow-server."
    echo ""

    # Check prerequisites
    if ! command -v curl &> /dev/null; then
        log_error "curl is required but not installed"
        exit 1
    fi

    detect_platform
    get_latest_version
    install_binary
    check_path
    print_summary
}

# Run main function
main
