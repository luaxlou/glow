#!/bin/bash
set -e

# Glow Local Installation Script
# Installs both glow-server and glow CLI for local use
# Does NOT register or start any system service/daemon
# Works on both macOS and Linux

# Configuration
REPO="luaxlou/glow"
INSTALL_DIR="${HOME}/.local/bin"

# Platform-specific data directory
if [[ "$OSTYPE" == "darwin"* ]]; then
    DATA_DIR="${HOME}/Library/Application Support/glow-server"
else
    DATA_DIR="${HOME}/.glow-server"
fi

SERVER_URL="http://localhost:32102"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    OS="$(uname -s)"
    ARCH="$(uname -m)"

    case "$OS" in
        Linux)
            OS="linux"
            # Adjust data dir for Linux
            DATA_DIR="${HOME}/.local/share/glow-server"
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

# Get latest release version (including pre-releases)
get_latest_version() {
    log_info "Fetching latest release version..."

    # Method 1: Get all releases (including pre-releases) and find the latest by version number
    # Get first page of releases (up to 30, should be enough)
    local releases_json=$(curl -s "https://api.github.com/repos/${REPO}/releases?per_page=30" 2>/dev/null)
    
    if [ -n "$releases_json" ]; then
        # Extract all tag names, remove 'v' prefix, and sort by version number
        # This ensures we get the truly latest version, including pre-releases
        local temp_file=$(mktemp)
        echo "$releases_json" | grep -oE '"tag_name":\s*"v[0-9][^"]*"' | sed -E 's/.*"v([^"]+)".*/\1/' > "$temp_file"
        
        if [ -s "$temp_file" ]; then
            # Sort by version and get the latest
            VERSION=$(sort -V "$temp_file" | tail -1)
        fi
        rm -f "$temp_file"
    fi

    # Method 2: Fallback to /releases/latest (non-prerelease only)
    if [ -z "$VERSION" ]; then
        log_warn "Failed to get all releases, trying /releases/latest (non-prerelease only)..."
        VERSION=$(curl -s https://api.github.com/repos/${REPO}/releases/latest 2>/dev/null | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/' | sed 's/^v//')
    fi

    # Method 3: Fallback to VERSION file from main branch
    if [ -z "$VERSION" ]; then
        log_warn "GitHub API unavailable, trying VERSION file from main branch..."
        VERSION=$(curl -s https://raw.githubusercontent.com/${REPO}/main/VERSION 2>/dev/null | tr -d '[:space:]' || echo "")
    fi

    if [ -z "$VERSION" ]; then
        log_error "Failed to fetch latest version"
        log_error "Please specify version manually or check your internet connection"
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

    log_info "Downloading ${binary_name} from ${download_url}..."

    # Download binary
    if ! curl -fSL -o "${output_path}" "${download_url}"; then
        log_error "Failed to download ${binary_name}"
        exit 1
    fi

    # Verify checksum from SHA256SUMS.txt
    log_info "Verifying checksum..."
    if [ -f "${TMP_DIR:-/tmp}/SHA256SUMS.txt" ]; then
        if command -v sha256sum &> /dev/null; then
            downloaded_checksum=$(sha256sum "${output_path}" | awk '{print $1}')
        elif command -v shasum &> /dev/null; then
            downloaded_checksum=$(shasum -a 256 "${output_path}" | awk '{print $1}')
        else
            log_warn "No checksum tool available, skipping verification"
            return
        fi

        expected_checksum=$(grep "${filename}" "${TMP_DIR:-/tmp}/SHA256SUMS.txt" | awk '{print $1}')

        if [ "$downloaded_checksum" != "$expected_checksum" ]; then
            log_error "Checksum verification failed!"
            log_error "Expected: $expected_checksum"
            log_error "Got: $downloaded_checksum"
            rm -f "${output_path}"
            exit 1
        fi

        log_info "Checksum verified successfully"
    else
        log_warn "SHA256SUMS.txt not found, skipping verification"
    fi

    # Make executable
    chmod +x "${output_path}"
}

# Install binaries
install_binaries() {
    log_step "Installing binaries to ${INSTALL_DIR}..."

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    export TMP_DIR

    # Download SHA256SUMS.txt
    log_info "Downloading SHA256SUMS.txt..."
    local checksum_url="https://github.com/${REPO}/releases/download/v${VERSION}/SHA256SUMS.txt"
    if ! curl -fSL -o "${TMP_DIR}/SHA256SUMS.txt" "${checksum_url}"; then
        log_warn "Failed to download SHA256SUMS.txt, skipping verification"
    fi

    # Create install directory if it doesn't exist
    mkdir -p "${INSTALL_DIR}"

    # Check if binaries already exist and create backup
    if [ -f "${INSTALL_DIR}/glow-server" ] || [ -f "${INSTALL_DIR}/glow" ]; then
        BACKUP_DIR="${HOME}/.glow-backup-$(date +%Y%m%d-%H%M%S)"
        mkdir -p "${BACKUP_DIR}"
        if [ -f "${INSTALL_DIR}/glow-server" ]; then
            cp "${INSTALL_DIR}/glow-server" "${BACKUP_DIR}/"
            log_info "Backed up existing glow-server to ${BACKUP_DIR}"
        fi
        if [ -f "${INSTALL_DIR}/glow" ]; then
            cp "${INSTALL_DIR}/glow" "${BACKUP_DIR}/"
            log_info "Backed up existing glow to ${BACKUP_DIR}"
        fi
    fi

    # Download and install glow-server
    glow_server_tmp="${TMP_DIR}/glow-server-${VERSION}"
    download_binary "glow-server" "${glow_server_tmp}"
    mv "${glow_server_tmp}" "${INSTALL_DIR}/glow-server"
    log_info "Installed glow-server"

    # Download and install glow
    glow_tmp="${TMP_DIR}/glow-${VERSION}"
    download_binary "glow" "${glow_tmp}"
    mv "${glow_tmp}" "${INSTALL_DIR}/glow"
    log_info "Installed glow"

    # Cleanup
    rm -rf "${TMP_DIR}"
}

# Check if database exists
check_existing_database() {
    local db_file="${DATA_DIR}/db/glow.db"
    local config_dir="${DATA_DIR}/config"

    if [ -f "$db_file" ]; then
        log_warn "Detected existing database at ${db_file}"
        log_warn "Reusing existing database and configuration"
        return 0
    fi

    return 1
}

# Create data directory
create_data_dir() {
    log_step "Setting up data directory at ${DATA_DIR}..."

    # Check if database already exists
    if check_existing_database; then
        log_info "Skipping database creation - will reuse existing one"
        log_info "To perform a clean install, manually remove:"
        log_info "  - ${DATA_DIR}/db/"
        log_info "  - ${DATA_DIR}/config/"
    else
        mkdir -p "${DATA_DIR}/db"
        mkdir -p "${DATA_DIR}/config"
        log_info "Created new database and config directories"
    fi

    # Always create logs and apps directories (safe to recreate)
    mkdir -p "${DATA_DIR}/logs"
    mkdir -p "${DATA_DIR}/apps"

    log_info "Data directory setup completed"
}

# Generate API key and configure glow
setup_api_key() {
    log_step "Generating API key..."

    # Ensure glow-server is in PATH for this session
    export PATH="${PATH}:${INSTALL_DIR}"

    # Run glow-server keygen with GLOW_DATA_DIR environment variable
    # Capture the output to get the API key
    KEYGEN_OUTPUT=$(GLOW_DATA_DIR="${DATA_DIR}" glow-server keygen 2>&1)
    if [ $? -ne 0 ]; then
        log_error "Failed to generate API key"
        log_error "${KEYGEN_OUTPUT}"
        exit 1
    fi

    # Extract the API key from the output
    # Output format: "Generated New API Key: <key>" or "Existing API Key: <key>"
    API_KEY=$(echo "${KEYGEN_OUTPUT}" | grep -oE 'Generated New API Key: [0-9a-f]+|Existing API Key: [0-9a-f]+' | sed 's/.*: //')

    if [ -z "$API_KEY" ]; then
        log_error "Failed to extract API key from keygen output"
        log_error "Output was: ${KEYGEN_OUTPUT}"
        exit 1
    fi

    log_info "API Key generated: ${API_KEY}"

    # Configure glow CLI with default context
    log_step "Configuring glow CLI..."

    if ! glow context add default \
        --url="${SERVER_URL}" \
        --key="${API_KEY}"; then
        log_warn "Failed to configure glow CLI automatically"
        log_warn "You can configure it manually with:"
        echo ""
        echo "  glow context add default --url=${SERVER_URL} --key=${API_KEY}"
        echo ""
    else
        log_info "glow CLI configured with default context"
    fi
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
    log_info "Local installation completed!"
    echo "=========================================="
    echo ""
    echo "Installed binaries:"
    echo "  - glow-server: ${INSTALL_DIR}/glow-server"
    echo "  - glow:       ${INSTALL_DIR}/glow"
    echo ""
    echo "Data directory: ${DATA_DIR}"
    echo ""
    echo "NOTE: This is a LOCAL installation."
    echo "No system service or daemon has been registered."
    echo ""
    echo "To start using Glow:"
    echo ""
    echo "  1. Make sure ${INSTALL_DIR} is in your PATH (see warning above)"
    echo ""
    echo "  2. Start the server in foreground:"
    echo "     glow-server serve --data-dir=\"${DATA_DIR}\""
    echo ""
    echo "  3. In another terminal, test the connection:"
    echo "     glow get apps"
    echo ""
    echo "  4. To stop the server, press Ctrl+C in the server terminal"
    echo ""
    echo "For more information, see: https://github.com/${REPO}"
    echo ""
}

# Main installation flow
main() {
    echo ""
    echo "Glow Local Installation Script"
    echo "=========================================="
    echo ""
    echo "This script installs glow-server and glow for LOCAL USE."
    echo "It will NOT register any system service or daemon."
    echo "You will start the server manually in foreground mode."
    echo ""

    # Check prerequisites
    if ! command -v curl &> /dev/null; then
        log_error "curl is required but not installed"
        exit 1
    fi

    detect_platform
    get_latest_version
    install_binaries
    create_data_dir
    setup_api_key
    check_path
    print_summary
}

# Run main function
main
