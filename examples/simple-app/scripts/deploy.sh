#!/bin/bash
set -e

# Deploy script for Glow
# Usage: ./scripts/deploy.sh [app_name]
#
# If app_name is not provided, will scan cmd/ directory and prompt for selection
#
# Features:
# - Automatic cross-compilation for Linux/amd64
# - Automatic fallback to curl upload if glow CLI times out (for large files)
# - Extended timeout (10 minutes) for file uploads via curl
#
# Environment variables:
# - GLOW_API_KEY: API key for Glow server (optional, will try to get from context)
# - GOOS, GOARCH: Override target platform (defaults to linux/amd64)

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

# Function to deploy using curl with extended timeout
deploy_with_curl() {
    local binary_path="$1"
    local app_name="$2"
    
    # Get current context info
    local context_info=$(glow context list 2>/dev/null | grep "^\*" || glow context list 2>/dev/null | head -2 | tail -1)
    if [ -z "$context_info" ]; then
        echo "❌ Error: Could not determine current glow context"
        return 1
    fi
    
    # Extract URL from context (3rd column)
    local server_url=$(echo "$context_info" | awk '{print $3}')
    if [ -z "$server_url" ]; then
        echo "❌ Error: Could not extract server URL from context"
        return 1
    fi
    
    # Get API key - try environment variable first, then try to extract from glow context
    local api_key="${GLOW_API_KEY:-}"
    if [ -z "$api_key" ]; then
        # Try to get from glow context show command (if supported)
        local context_show=$(glow context show 2>/dev/null)
        if [ -n "$context_show" ]; then
            api_key=$(echo "$context_show" | grep -iE "api.*key|key:" | awk -F: '{print $2}' | tr -d ' ' || echo "")
        fi
    fi
    
    if [ -z "$api_key" ]; then
        echo "❌ Error: API key not found"
        echo ""
        echo "   Please provide API key using one of these methods:"
        echo "   1. Set environment variable: export GLOW_API_KEY=your-api-key"
        echo "   2. Or ensure glow context is properly configured with: glow context show"
        echo ""
        echo "   To get your API key, run: glow-server keygen (on the server)"
        return 1
    fi
    
    echo "📤 Uploading ${binary_path} to ${server_url}..."
    echo "   Using curl with extended timeout (10 minutes)..."
    
    # Upload with curl using extended timeout (600 seconds = 10 minutes)
    local upload_url="${server_url}/apps/upload"
    local response=$(curl -s -w "\n%{http_code}" \
        --max-time 600 \
        --connect-timeout 30 \
        -H "Authorization: Bearer ${api_key}" \
        -F "file=@${binary_path}" \
        "${upload_url}" 2>&1)
    
    local http_code=$(echo "$response" | tail -1)
    local body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "200" ]; then
        echo "✅ File uploaded successfully"
        # Extract file path from response if available
        local uploaded_path=$(echo "$body" | grep -o '"data":"[^"]*"' | cut -d'"' -f4 || echo "")
        if [ -n "$uploaded_path" ]; then
            echo "   Uploaded to: ${uploaded_path}"
        fi
        return 0
    else
        echo "❌ Upload failed with HTTP code: ${http_code}"
        echo "   Response: ${body}"
        return 1
    fi
}

# Deploy using glow CLI with fallback to curl
if command -v glow &> /dev/null; then
    echo "🚀 Deploying using glow CLI..."
    local deploy_output
    deploy_output=$(glow deploy "${BINARY_PATH}" --name "${APP_NAME}" 2>&1)
    local deploy_exit_code=$?
    
    if [ $deploy_exit_code -eq 0 ]; then
        echo "✅ Deployment successful via glow CLI"
    else
        # Check if it's a timeout error
        if echo "$deploy_output" | grep -qi "timeout\|deadline exceeded"; then
            echo ""
            echo "⚠️  glow CLI upload timed out. Trying alternative method with curl..."
            echo ""
            if deploy_with_curl "${BINARY_PATH}" "${APP_NAME}"; then
                echo "✅ Deployment successful via curl"
            else
                echo "❌ Deployment failed with both methods"
                exit 1
            fi
        else
            echo "❌ Deployment failed:"
            echo "$deploy_output"
            exit 1
        fi
    fi
else
    echo "❌ Error: glow CLI not found"
    echo "   Install glow from: https://github.com/luaxlou/glow"
    exit 1
fi

echo "✅ Deployment complete!"
