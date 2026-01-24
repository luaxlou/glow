#!/bin/bash
set -e

# Deploy script for Glow
# Usage: ./scripts/deploy.sh [app_name]
#
# If app_name is not provided, will scan cmd/ directory and prompt for selection

APP_NAME="$1"

# If no argument provided, scan cmd/ directory
if [ -z "$APP_NAME" ]; then
    if [ -d "cmd" ]; then
        echo "🔍 Scanning cmd/ directory for applications..."

        # Find all subdirectories in cmd/
        apps=($(find cmd -maxdepth 1 -type d ! -name "cmd" -exec basename {} \;))

        if [ ${#apps[@]} -eq 0 ]; then
            echo "❌ Error: No applications found in cmd/ directory"
            echo "   Usage: $0 <app_name>"
            exit 1
        fi

        # If only one app found, use it
        if [ ${#apps[@]} -eq 1 ]; then
            APP_NAME="${apps[0]}"
            echo "✓ Found application: $APP_NAME"
        else
            # Multiple apps found, prompt for selection
            echo ""
            echo "Found ${#apps[@]} application(s):"
            echo ""

            # Display menu
            for i in "${!apps[@]}"; do
                echo "$((i+1))) ${apps[$i]}"
            done
            echo ""

            while true; do
                read -p "Select application to deploy (default=all): " choice
                echo ""

                # Empty input means deploy all
                if [ -z "$choice" ]; then
                    echo "🚀 Deploying all applications..."
                    for app_to_deploy in "${apps[@]}"; do
                        echo ""
                        echo "➤ Deploying $app_to_deploy..."
                        BINARY_PATH="bin/${app_to_deploy}"

                        # Build if needed (with cross-compilation)
                        if [ ! -f "${BINARY_PATH}" ]; then
                            CURRENT_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
                            TARGET_OS="${GOOS:-linux}"
                            TARGET_ARCH="${GOARCH:-amd64}"
                            
                            # Default to linux/amd64 for server deployment
                            if [ -z "$GOOS" ] && [ -z "$GOARCH" ]; then
                                TARGET_OS="linux"
                                TARGET_ARCH="amd64"
                            fi
                            
                            # Always set target platform for consistent builds
                            export GOOS="$TARGET_OS"
                            export GOARCH="$TARGET_ARCH"
                            
                            if go build -o "${BINARY_PATH}" "./cmd/${app_to_deploy}"; then
                                echo "  ✓ Built: ${BINARY_PATH} (${TARGET_OS}/${TARGET_ARCH})"
                            else
                                echo "  ❌ Build failed for $app_to_deploy"
                                continue
                            fi
                        fi

                        # Deploy
                        if command -v glow &> /dev/null; then
                            glow deploy "${BINARY_PATH}" --name "${app_to_deploy}"
                            echo "  ✓ Deployed: $app_to_deploy"
                        else
                            echo "  ❌ Error: glow CLI not found"
                            exit 1
                        fi
                    done
                    echo ""
                    echo "✅ All applications deployed!"
                    exit 0
                fi

                # Check if choice is a valid app number
                if [ "$choice" -ge 1 ] && [ "$choice" -le ${#apps[@]} ]; then
                    APP_NAME="${apps[$((choice-1))]}"
                    break
                fi

                echo "❌ Invalid selection. Please try again."
                echo ""
            done
        fi
    else
        echo "❌ Error: cmd/ directory not found"
        echo "   Usage: $0 <app_name>"
        exit 1
    fi
fi

BINARY_PATH="bin/${APP_NAME}"

echo "🚀 Deploying ${APP_NAME} to Glow..."

# Detect current platform
CURRENT_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
CURRENT_ARCH=$(uname -m)

# Determine target platform for cross-compilation
# Glow server typically runs on Linux amd64, so cross-compile if we're on macOS
# Allow override via environment variables
TARGET_OS="${GOOS:-linux}"
TARGET_ARCH="${GOARCH:-amd64}"

# If not explicitly set, default to linux/amd64 for server deployment
if [ -z "$GOOS" ] && [ -z "$GOARCH" ]; then
    TARGET_OS="linux"
    TARGET_ARCH="amd64"
fi

# Build the application if binary doesn't exist or is outdated
if [ ! -f "${BINARY_PATH}" ] || [ "$CURRENT_OS" != "linux" ]; then
    echo "📦 Building ${APP_NAME}..."

    # Create bin directory if it doesn't exist
    mkdir -p bin

    # Cross-compile if needed
    if [ "$CURRENT_OS" != "linux" ] || [ "$CURRENT_ARCH" != "$TARGET_ARCH" ]; then
        echo "   Cross-compiling for ${TARGET_OS}/${TARGET_ARCH}..."
    fi
    export GOOS="$TARGET_OS"
    export GOARCH="$TARGET_ARCH"

    # Build the application
    if go build -o "${BINARY_PATH}" "./cmd/${APP_NAME}"; then
        echo "✓ Build completed: ${BINARY_PATH} (${TARGET_OS}/${TARGET_ARCH})"
    else
        echo "❌ Error: Build failed"
        echo "   Please check your code and try again"
        exit 1
    fi
else
    echo "✓ Using existing binary: ${BINARY_PATH}"
fi

# Deploy using glow CLI
if command -v glow &> /dev/null; then
    glow deploy "${BINARY_PATH}" --name "${APP_NAME}"
else
    echo "❌ Error: glow CLI not found"
    echo "   Install glow from: https://github.com/luaxlou/glow"
    exit 1
fi

echo "✅ Deployment complete!"
