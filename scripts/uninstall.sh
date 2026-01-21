#!/bin/bash
set -e

# Glow Server Uninstallation Script
# Removes binaries and service files, but preserves configuration and database

# Configuration
INSTALL_DIR="/usr/local/bin"
DATA_DIR="/var/lib/glow-server"
SERVICE_NAME="glow-server"
BACKUP_DIR="/tmp/glow-server-backup-$(date +%Y%m%d-%H%M%S)"

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

# Detect OS
OS="$(uname -s)"
case "$OS" in
    Linux) OS="linux" ;;
    Darwin) OS="darwin" ;;
    *)
        log_error "Unsupported operating system: $OS"
        exit 1
        ;;
esac

# Print warning and get confirmation
print_warning() {
    echo ""
    echo "=========================================="
    echo "Glow Server Uninstallation"
    echo "=========================================="
    echo ""
    log_warn "This will uninstall glow-server and glow from your system."
    echo ""
    echo "The following will be REMOVED:"
    echo "  - Binaries: ${INSTALL_DIR}/glow-server, ${INSTALL_DIR}/glow"
    echo "  - System service files (systemd/launchd)"
    echo ""
    echo "The following will be PRESERVED:"
    echo "  - Configuration: ${DATA_DIR}/config/"
    echo "  - Database: ${DATA_DIR}/db/"
    echo "  - Logs: ${DATA_DIR}/logs/"
    echo "  - Applications: ${DATA_DIR}/apps/"
    echo ""
    echo "To completely remove glow-server, manually delete:"
    echo "  sudo rm -rf ${DATA_DIR}"
    echo ""
}

# Backup binaries and service files before uninstallation
backup_before_uninstall() {
    log_step "Creating backup before uninstallation..."

    mkdir -p "${BACKUP_DIR}"

    # Backup binaries if they exist
    if [ -f "${INSTALL_DIR}/glow-server" ]; then
        cp "${INSTALL_DIR}/glow-server" "${BACKUP_DIR}/"
        log_info "Backed up glow-server binary"
    fi

    if [ -f "${INSTALL_DIR}/glow" ]; then
        cp "${INSTALL_DIR}/glow" "${BACKUP_DIR}/"
        log_info "Backed up glow binary"
    fi

    # Backup service files
    if [ "$OS" = "linux" ]; then
        if [ -f "/etc/systemd/system/${SERVICE_NAME}.service" ]; then
            sudo cp "/etc/systemd/system/${SERVICE_NAME}.service" "${BACKUP_DIR}/"
            log_info "Backed up systemd service file"
        fi
    elif [ "$OS" = "darwin" ]; then
        if [ -f "/Library/LaunchDaemons/com.glowserver.${SERVICE_NAME}.plist" ]; then
            sudo cp "/Library/LaunchDaemons/com.glowserver.${SERVICE_NAME}.plist" "${BACKUP_DIR}/"
            log_info "Backed up launchd plist file"
        fi
    fi

    log_info "Backup created at: ${BACKUP_DIR}"
}

# Stop and disable service
stop_service() {
    log_step "Stopping and disabling service..."

    if [ "$OS" = "linux" ]; then
        if systemctl is-active --quiet ${SERVICE_NAME} 2>/dev/null; then
            sudo systemctl stop ${SERVICE_NAME}
            log_info "Service stopped"
        fi

        if systemctl is-enabled --quiet ${SERVICE_NAME} 2>/dev/null; then
            sudo systemctl disable ${SERVICE_NAME}
            log_info "Service disabled"
        fi
    elif [ "$OS" = "darwin" ]; then
        if sudo launchctl list | grep -q "com.glowserver.${SERVICE_NAME}"; then
            sudo launchctl unload -w /Library/LaunchDaemons/com.glowserver.${SERVICE_NAME}.plist 2>/dev/null || true
            log_info "Service unloaded"
        fi
    fi
}

# Remove service files
remove_service_files() {
    log_step "Removing service files..."

    if [ "$OS" = "linux" ]; then
        if [ -f "/etc/systemd/system/${SERVICE_NAME}.service" ]; then
            sudo rm -f "/etc/systemd/system/${SERVICE_NAME}.service"
            sudo systemctl daemon-reload
            log_info "Systemd service file removed"
        fi
    elif [ "$OS" = "darwin" ]; then
        if [ -f "/Library/LaunchDaemons/com.glowserver.${SERVICE_NAME}.plist" ]; then
            sudo rm -f "/Library/LaunchDaemons/com.glowserver.${SERVICE_NAME}.plist"
            log_info "Launchd plist file removed"
        fi
    fi
}

# Remove binaries
remove_binaries() {
    log_step "Removing binaries..."

    if [ -w "${INSTALL_DIR}" ]; then
        SUDO=""
    else
        SUDO="sudo"
    fi

    if [ -f "${INSTALL_DIR}/glow-server" ]; then
        $SUDO rm -f "${INSTALL_DIR}/glow-server"
        log_info "Removed glow-server binary"
    fi

    if [ -f "${INSTALL_DIR}/glow" ]; then
        $SUDO rm -f "${INSTALL_DIR}/glow"
        log_info "Removed glow binary"
    fi
}

# Print summary
print_summary() {
    echo ""
    echo "=========================================="
    log_info "Uninstallation completed successfully!"
    echo "=========================================="
    echo ""
    echo "Removed:"
    echo "  - Binaries from ${INSTALL_DIR}"
    echo "  - Service files"
    echo ""
    echo "Preserved:"
    echo "  - Configuration and data in ${DATA_DIR}"
    echo ""
    echo "Backup location:"
    echo "  - ${BACKUP_DIR}"
    echo ""
    echo "To completely remove glow-server, run:"
    echo "  sudo rm -rf ${DATA_DIR}"
    echo "  sudo rm -rf ${BACKUP_DIR}"
    echo ""
}

# Main uninstallation flow
main() {
    print_warning

    # Ask for confirmation
    read -p "Do you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Uninstallation cancelled"
        exit 0
    fi

    backup_before_uninstall
    stop_service
    remove_service_files
    remove_binaries
    print_summary
}

# Run main function
main
