#!/bin/bash
set -e

# Glow Project Initialization Script
# Usage: curl -fsSL https://raw.githubusercontent.com/{owner}/{repo}/main/scripts/init-project.sh | bash
# Or with options: curl ... | bash -s -- --skip-ai --force

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
FORCE=false
SKIP_AI=false
SHOW_HELP=false

# Parse command-line arguments
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --force)
                FORCE=true
                shift
                ;;
            --skip-ai)
                SKIP_AI=true
                shift
                ;;
            --help|-h)
                SHOW_HELP=true
                shift
                ;;
            *)
                # Unknown option
                echo -e "${RED}Unknown option: $1${NC}"
                echo "Use --help to see available options"
                exit 1
                ;;
        esac
    done
}

# Show help message
show_help() {
    cat << EOF
Glow Project Initialization Script

USAGE:
    curl -fsSL <url> | bash
    curl -fsSL <url> | bash -s -- [options]

OPTIONS:
    --force       Skip confirmation prompts and overwrite existing files
    --skip-ai     Skip AI tool integration (Claude Code skills)
    --help, -h    Show this help message

DESCRIPTION:
    This script initializes a project for Glow governance by creating:
    - Standard directory structure (cmd/, bin/, scripts/)
    - deploy.sh script for building and deploying applications
    - AI tool integration (optional)

EXAMPLES:
    # Basic initialization
    curl -fsSL https://raw.githubusercontent.com/luaxlou/glow/main/scripts/init-project.sh | bash

    # Skip AI tool integration
    curl -fsSL https://raw.githubusercontent.com/luaxlou/glow/main/scripts/init-project.sh | bash -s -- --skip-ai

    # Force overwrite existing files
    curl -fsSL https://raw.githubusercontent.com/luaxlou/glow/main/scripts/init-project.sh | bash -s -- --force

EOF
    exit 0
}

# Detect platform
detect_platform() {
    OS=$(uname -s)
    case "$OS" in
        Linux*)
            PLATFORM="linux"
            ;;
        Darwin*)
            PLATFORM="macos"
            ;;
        *)
            echo -e "${RED}Unsupported platform: $OS${NC}"
            exit 1
            ;;
    esac
}

# Analyze project structure
analyze_project() {
    HAS_CMD_DIR=false
    HAS_BIN_DIR=false
    HAS_SCRIPTS_DIR=false
    HAS_DEPLOY_SH=false
    HAS_CLAUDE_DIR=false
    HAS_APP_YAML=false

    [ -d "cmd" ] && HAS_CMD_DIR=true || true
    [ -d "bin" ] && HAS_BIN_DIR=true || true
    [ -d "scripts" ] && HAS_SCRIPTS_DIR=true || true
    [ -f "scripts/deploy.sh" ] && HAS_DEPLOY_SH=true || true
    [ -d ".claude" ] && HAS_CLAUDE_DIR=true || true
    [ -f "app.yaml" ] && HAS_APP_YAML=true || true
}

# Print project analysis
print_analysis() {
    echo -e "${BLUE}📊 Project Structure Analysis:${NC}"
    echo ""

    local items=(
        "cmd/ $HAS_CMD_DIR"
        "bin/ $HAS_BIN_DIR"
        "scripts/ $HAS_SCRIPTS_DIR"
        "scripts/deploy.sh $HAS_DEPLOY_SH"
        "app.yaml $HAS_APP_YAML"
        ".claude/ $HAS_CLAUDE_DIR"
    )

    for item in "${items[@]}"; do
        local name=$(echo "$item" | cut -d' ' -f1) || continue
        local exists=$(echo "$item" | cut -d' ' -f2) || continue

        if [ "$exists" = "true" ]; then
            echo -e "   ${GREEN}✅${NC} $name"
        else
            echo -e "   ${RED}❌${NC} $name"
        fi
    done
    echo ""
}

# Confirm before proceeding
confirm_proceed() {
    if [ "$FORCE" = true ]; then
        return 0
    fi

    # Check if running in non-interactive mode (piped from curl)
    # If stdin is not a terminal, skip confirmation and proceed
    if [ ! -t 0 ]; then
        echo "   ℹ Non-interactive mode detected, proceeding automatically..."
        return 0
    fi

    # Interactive mode: ask for confirmation
    read -p "Continue with initialization? [y/N] " input
    input=$(echo "$input" | tr '[:upper:]' '[:lower:]')

    if [ "$input" = "y" ] || [ "$input" = "yes" ]; then
        return 0
    else
        echo -e "${RED}❌ Initialization cancelled.${NC}"
        exit 0
    fi
}

# Create standard directories
create_directories() {
    echo -e "${BLUE}📁 Creating standard directories...${NC}"

    local dirs=(
        "cmd:$HAS_CMD_DIR"
        "bin:$HAS_BIN_DIR"
        "scripts:$HAS_SCRIPTS_DIR"
    )

    for item in "${dirs[@]}"; do
        local dir=$(echo "$item" | cut -d':' -f1)
        local exists=$(echo "$item" | cut -d':' -f2)

        if [ "$exists" = "true" ]; then
            echo "   ⊙ $dir/ already exists, skipping"
        else
            mkdir -p "$dir"
            echo -e "   ${GREEN}✓${NC} Created $dir/"
        fi
    done
}

# Create deploy script
create_deploy_script() {
    echo -e "${BLUE}📜 Creating deploy script...${NC}"

    if [ "$HAS_DEPLOY_SH" = "true" ] && [ "$FORCE" != "true" ]; then
        echo "   ⊙ scripts/deploy.sh already exists, skipping"
        return 0
    fi

    cat > scripts/deploy.sh << 'DEPLOY_EOF'
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
DEPLOY_EOF

    chmod +x scripts/deploy.sh
    echo -e "   ${GREEN}✓${NC} Created scripts/deploy.sh"
}

# Create app.yaml template
create_app_yaml() {
    echo -e "${BLUE}📄 Creating app.yaml template...${NC}"

    if [ "$HAS_APP_YAML" = "true" ] && [ "$FORCE" != "true" ]; then
        echo "   ⊙ app.yaml already exists, skipping"
        return 0
    fi

    # Detect application name
    APP_NAME=""
    if [ -d "cmd" ]; then
        # Find all subdirectories in cmd/
        apps=($(find cmd -maxdepth 1 -type d ! -name "cmd" -exec basename {} \;))
        
        if [ ${#apps[@]} -eq 1 ]; then
            # Single app, use its name
            APP_NAME="${apps[0]}"
        elif [ ${#apps[@]} -gt 1 ]; then
            # Multiple apps, use first one or project directory name
            APP_NAME="${apps[0]}"
            echo "   ⚠ Multiple apps found, using first app: $APP_NAME"
        fi
    fi

    # Fallback to project directory name if no app found
    if [ -z "$APP_NAME" ]; then
        APP_NAME=$(basename "$(pwd)")
        echo "   ⚠ No app found in cmd/, using directory name: $APP_NAME"
    fi

    # Generate a random port in the range 33000-33999
    RANDOM_PORT=$((33000 + RANDOM % 1000))

    cat > app.yaml << APP_YAML_EOF
# Glow App 配置文件
#
# 使用方法：
#   1. 构建应用: ./scripts/deploy.sh          # 自动构建并部署
#   2. 应用配置: glow apply -f app.yaml       # 应用配置并绑定资源
#   3. 启动应用: glow start app ${APP_NAME}  # 启动应用
#   4. 查看状态: glow get app ${APP_NAME}     # 查看应用状态
#
# 说明：
# - 应用使用 MySQL 和 HTTP 服务
# - Ingress（域名绑定）通过 spec.domain 声明
# - 资源绑定通过 spec.resources 声明
# - binary 和 workingDir 采用约定大于配置原则（可省略）

apiVersion: v1
kind: App
metadata:
  name: ${APP_NAME}
spec:
  # 应用执行配置（约定大于配置）
  # binary: 默认为 <data-dir>/apps/<app-name>/<app-name>
  # workingDir: 默认为 <data-dir>/apps/<app-name>
  # 因此可以省略这两个字段

  args: []

  # 环境变量（可选）
  env:
    - name: ENV
      value: development
    - name: LOG_LEVEL
      value: debug

  # 应用配置（可选）⭐
  # 这些配置会写入到 <data-dir>/apps/<app-name>/<app-name>_local_config.json
  # 应用代码可以使用 glowconfig 读取这些配置
  config:
    # 数据库配置（用户自行提供）
    mysql_dsn: "root:password@tcp(localhost:3306)/${APP_NAME}_db"

    # Redis 配置（用户自行提供）
    redis_addr: "localhost:6379"

    # 应用配置
    log_level: "debug"
    max_connections: 100

  # HTTP 服务端口
  # 应用会通过 OP_APP_PORT 环境变量收到端口 ${RANDOM_PORT}
  port: ${RANDOM_PORT}

  # Ingress 配置（可选）
  # 指定 domain 后，glow-server 会自动配置 Nginx 反向代理
  # domain: ${APP_NAME}.local
APP_YAML_EOF

    echo -e "   ${GREEN}✓${NC} Created app.yaml (app: ${APP_NAME}, port: ${RANDOM_PORT})"
}

# Setup AI tools
setup_ai_tools() {
    if [ "$SKIP_AI" = true ]; then
        echo -e "\n${YELLOW}🤖 AI Tool Integration skipped (--skip-ai flag used)${NC}"
        return 0
    fi

    echo -e "\n${BLUE}🤖 AI Tool Integration${NC}"

    # Check for glow CLI installation
    if ! command -v glow &> /dev/null; then
        echo -e "   ${YELLOW}⚠ glow CLI not found${NC}"
        echo "   ⊙ Skipping AI tool configuration"
        echo "   💡 Install glow CLI first to enable AI tool integration"
        return 0
    fi

    # Try to find glow installation path
    GLOW_PATH=$(find_glow_install_path)
    if [ -z "$GLOW_PATH" ]; then
        echo -e "   ${YELLOW}⚠ Could not find glow installation directory${NC}"
        echo "   ⊙ Skipping AI tool configuration"
        return 0
    fi

    # Setup Claude Code skills
    setup_claude_code
}

# Setup Claude Code integration
setup_claude_code() {
    echo "   📦 Setting up Claude Code integration..."

    local source_skills_dir="$GLOW_PATH/.claude/skills"
    local dest_skills_dir=".claude/skills"

    # Check if source skills directory exists
    if [ ! -d "$source_skills_dir" ]; then
        echo -e "   ${YELLOW}⊙ Glow skills not found at $source_skills_dir${NC}"
        echo "   ⊙ Skipping skill copy"
        return 0
    fi

    # Create destination directory
    mkdir -p "$dest_skills_dir"

    # Skills to copy
    local skills=("glow-sdk" "glow-deploy" "glow-debug")

    for skill in "${skills[@]}"; do
        local source_skill="$source_skills_dir/$skill"
        local dest_skill="$dest_skills_dir/$skill"

        if [ ! -d "$source_skill" ]; then
            echo "   ⊙ Skill $skill not found, skipping"
            continue
        fi

        # Copy the skill directory
        if cp -r "$source_skill" "$dest_skill" 2>/dev/null; then
            echo "   ✓ Copied skill: $skill"
        else
            echo -e "   ${YELLOW}⚠ Warning: Failed to copy skill $skill${NC}"
        fi
    done
}

# Find glow installation path
find_glow_install_path() {
    # Priority 1: Check if we're in the glow repo (development environment)
    local current_dir="$(pwd)"
    while [ "$current_dir" != "/" ]; do
        if [ -d "$current_dir/.git" ] && [ -d "$current_dir/cmd/glow" ]; then
            echo "$current_dir"
            return 0
        fi
        current_dir="$(dirname "$current_dir")"
    done

    # Priority 2: Try to find the path relative to the binary
    if [ -n "$(which glow)" ]; then
        local glow_bin="$(which glow)"
        local glow_dir="$(dirname "$glow_bin")"

        # Check current directory
        if [ -d "$glow_dir/.claude/skills" ]; then
            echo "$glow_dir"
            return 0
        fi

        # Check parent directory (in case binary is in bin/)
        local parent_dir="$(dirname "$glow_dir")"
        if [ -d "$parent_dir/.claude/skills" ]; then
            echo "$parent_dir"
            return 0
        fi
    fi

    return 1
}

# Main function
main() {
    parse_arguments "$@"

    if [ "$SHOW_HELP" = true ]; then
        show_help
    fi

    echo -e "${BLUE}🚀 Initializing project for Glow governance...${NC}"

    # Get current directory
    local pwd="$(pwd)"
    echo -e "${BLUE}📁 Working directory: ${pwd}${NC}"
    echo ""

    # Detect platform
    detect_platform

    # Analyze project structure
    analyze_project
    print_analysis

    # Ask for confirmation
    confirm_proceed
    echo ""

    # Create standard directories
    create_directories

    # Create deploy script
    create_deploy_script

    # Create app.yaml template
    create_app_yaml

    # Setup AI tools
    setup_ai_tools

    echo ""
    echo -e "${GREEN}✅ Project initialization complete!${NC}"
    echo -e "${BLUE}📝 Next steps:${NC}"
    echo "   - Review the generated files (app.yaml, scripts/deploy.sh)"
    echo "   - Edit app.yaml to configure your application"
    echo "   - Run './scripts/deploy.sh' to build and deploy your application"
    echo "   - Run 'glow apply -f app.yaml' to apply configuration"
}

# Run main only if script is executed directly (not sourced)
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    main "$@"
fi
