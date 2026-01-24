#!/bin/bash
set -e

# Glow Server Installation Script
# Installs both glow-server and glow CLI from GitHub Releases

# Configuration
REPO="luaxlou/glow"
INSTALL_DIR="/usr/local/bin"
DATA_DIR="/var/lib/glow-server"
SERVICE_NAME="glow-server"

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

# Backup existing installation
backup_existing_installation() {
    local backup_dir=$1

    log_info "Backing up existing installation to ${backup_dir}..."

    mkdir -p "${backup_dir}"

    # Backup existing binaries
    if [ -f "${INSTALL_DIR}/glow-server" ]; then
        cp "${INSTALL_DIR}/glow-server" "${backup_dir}/"
        log_info "Backed up existing glow-server"
    fi

    if [ -f "${INSTALL_DIR}/glow" ]; then
        cp "${INSTALL_DIR}/glow" "${backup_dir}/"
        log_info "Backed up existing glow"
    fi

    # Backup service files
    if [ "$OS" = "linux" ] && [ -f "/etc/systemd/system/${SERVICE_NAME}.service" ]; then
        sudo cp "/etc/systemd/system/${SERVICE_NAME}.service" "${backup_dir}/"
        log_info "Backed up systemd service file"
    elif [ "$OS" = "darwin" ] && [ -f "/Library/LaunchDaemons/com.glowserver.${SERVICE_NAME}.plist" ]; then
        sudo cp "/Library/LaunchDaemons/com.glowserver.${SERVICE_NAME}.plist" "${backup_dir}/"
        log_info "Backed up launchd plist file"
    fi

    # Backup config and database if they exist
    if [ -d "${DATA_DIR}/config" ]; then
        sudo cp -r "${DATA_DIR}/config" "${backup_dir}/"
        log_info "Backed up configuration"
    fi

    if [ -d "${DATA_DIR}/db" ]; then
        sudo cp -r "${DATA_DIR}/db" "${backup_dir}/"
        log_info "Backed up database"
    fi
}

# Install binaries
install_binaries() {
    log_info "Installing binaries to ${INSTALL_DIR}..."

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    export TMP_DIR

    # Download SHA256SUMS.txt
    log_info "Downloading SHA256SUMS.txt..."
    local checksum_url="https://github.com/${REPO}/releases/download/v${VERSION}/SHA256SUMS.txt"
    if ! curl -fSL -o "${TMP_DIR}/SHA256SUMS.txt" "${checksum_url}"; then
        log_warn "Failed to download SHA256SUMS.txt, skipping verification"
    fi

    # Check if we have write permission
    if [ ! -w "${INSTALL_DIR}" ]; then
        SUDO="sudo"
    else
        SUDO=""
    fi

    # Check if binaries already exist and create backup
    if [ -f "${INSTALL_DIR}/glow-server" ] || [ -f "${INSTALL_DIR}/glow" ]; then
        BACKUP_DIR="/tmp/glow-server-backup-$(date +%Y%m%d-%H%M%S)"
        backup_existing_installation "${BACKUP_DIR}"
        log_info "Backup created at ${BACKUP_DIR}"
    fi

    # Download and install glow-server
    glow_server_tmp="${TMP_DIR}/glow-server-${VERSION}"
    download_binary "glow-server" "${glow_server_tmp}"
    $SUDO mv "${glow_server_tmp}" "${INSTALL_DIR}/glow-server"
    log_info "Installed glow-server"

    # Download and install glow
    glow_tmp="${TMP_DIR}/glow-${VERSION}"
    download_binary "glow" "${glow_tmp}"
    $SUDO mv "${glow_tmp}" "${INSTALL_DIR}/glow"
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
    log_info "Setting up data directory at ${DATA_DIR}..."

    if [ ! -w "/var/lib" ]; then
        SUDO="sudo"
    else
        SUDO=""
    fi

    # Check if database already exists
    if check_existing_database; then
        log_info "Skipping database creation - will reuse existing one"
        log_info "To perform a clean install, manually remove:"
        log_info "  - ${DATA_DIR}/db/"
        log_info "  - ${DATA_DIR}/config/"
    else
        $SUDO mkdir -p "${DATA_DIR}/db"
        $SUDO mkdir -p "${DATA_DIR}/config"
        log_info "Created new database and config directories"
    fi

    # Always create logs and apps directories (safe to recreate)
    $SUDO mkdir -p "${DATA_DIR}/logs"
    $SUDO mkdir -p "${DATA_DIR}/apps"

    # Set permissions
    $SUDO chown -R $(whoami):$(whoami) "${DATA_DIR}" 2>/dev/null || true

    log_info "Data directory setup completed"
}

# Generate API key and configure glow
setup_api_key() {
    log_info "Generating API key..."

    # Run glow-server keygen and capture output
    KEYGEN_OUTPUT=$(glow-server keygen 2>&1)
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
    log_info "Configuring glow CLI..."

    glow context add default \
        --url="http://localhost:32102" \
        --key="${API_KEY}"

    log_info "glow CLI configured with default context"
}

# Install and start service
install_service() {
    log_info "Installing system service..."

    if [ "$OS" = "linux" ]; then
        # Systemd
        if command -v systemctl &> /dev/null; then
            # Install service file
            sudo glow-server service install

            # Start and enable service
            sudo systemctl daemon-reload
            sudo systemctl enable ${SERVICE_NAME}
            sudo systemctl start ${SERVICE_NAME}
            log_info "Service installed and started"
        else
            log_warn "systemd not found, service not installed"
            log_info "You can start glow-server manually with: glow-server serve"
        fi
    elif [ "$OS" = "darwin" ]; then
        # Launchd
        if command -v launchctl &> /dev/null; then
            # Install service file
            sudo glow-server service install

            # Load and start service
            sudo launchctl load -w /Library/LaunchDaemons/com.glowserver.${SERVICE_NAME}.plist
            log_info "Service installed and started"
        else
            log_warn "launchctl not found, service not installed"
            log_info "You can start glow-server manually with: glow-server serve"
        fi
    fi
}

# Print summary
print_summary() {
    echo ""
    echo "=========================================="
    log_info "Installation completed successfully!"
    echo "=========================================="
    echo ""
    echo "Installed binaries:"
    echo "  - glow-server: $(which glow-server)"
    echo "  - glow:       $(which glow)"
    echo ""
    echo "Data directory: ${DATA_DIR}"
    echo ""
    echo "Quick start:"
    echo "  - Check status:  glow-server info"
    echo "  - View logs:    tail -f ${DATA_DIR}/logs/server.log"
    echo "  - Stop service: glow-server service stop"
    echo "  - Start service: glow-server service start"
    echo ""
    echo "For more information, see: https://github.com/${REPO}"
    echo ""
}

# Main installation flow
main() {
    echo ""
    echo "Glow Server Installation Script"
    echo "==============================="
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
    install_service
    print_summary
}

# Run main function
main
