#!/bin/bash
# Glow Build Script
# Compiles glow-server and glow CLI for multiple platforms

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Supported platforms
PLATFORMS=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
)

# Output directory
DIST_DIR="./dist"
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

print_info "=========================================="
print_info "Glow Build Script"
print_info "=========================================="
echo ""

# Get version from git tag or use default
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0-dev")
VERSION=${VERSION#v}  # Remove 'v' prefix if present

print_info "Version: $VERSION"
print_info "Output directory: $DIST_DIR"
echo ""

# Build function
build_binary() {
    local goos=$1
    local goarch=$2
    local binary_name=$3

    local output_name="${binary_name}-${goos}-${goarch}"
    local output_path="${DIST_DIR}/${output_name}"

    print_step "Building ${output_name}..."

    GOOS=$goos GOARCH=$goarch go build \
        -ldflags="-X 'main.Version=$VERSION' -X 'main.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" \
        -o "$output_path" \
        ./"$(get_binary_path $binary_name)"

    print_info "✓ Built: ${output_name}"

    # Generate checksum
    if command -v sha256sum &> /dev/null; then
        checksum=$(sha256sum "$output_path" | awk '{print $1}')
    elif command -v shasum &> /dev/null; then
        checksum=$(shasum -a 256 "$output_path" | awk '{print $1}')
    else
        print_error "No checksum tool available"
        exit 1
    fi

    echo "${checksum}  ${output_name}" >> "${DIST_DIR}/SHA256SUMS.txt"
    print_info "  Checksum: ${checksum}"
}

get_binary_path() {
    case $1 in
        "glow-server")
            echo "cmd/glow-server"
            ;;
        "glow")
            echo "cmd/glow"
            ;;
        *)
            print_error "Unknown binary: $1"
            exit 1
            ;;
    esac
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed"
    exit 1
fi

print_info "Go version: $(go version)"
echo ""

# Build for all platforms
for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r goos goarch <<< "$platform"

    print_step "Building for ${goos}/${goarch}..."

    # Build glow-server
    build_binary "$goos" "$goarch" "glow-server"

    # Build glow
    build_binary "$goos" "$goarch" "glow"

    echo ""
done

# Sort checksums file
if [ -f "${DIST_DIR}/SHA256SUMS.txt" ]; then
    sort "${DIST_DIR}/SHA256SUMS.txt" -o "${DIST_DIR}/SHA256SUMS.txt"
fi

print_info "=========================================="
print_info "Build completed!"
print_info "=========================================="
echo ""
print_info "Built binaries:"
ls -lh "$DIST_DIR"
echo ""
print_info "Checksums file: ${DIST_DIR}/SHA256SUMS.txt"
echo ""
print_info "Total files: $(ls -1 "$DIST_DIR" | wc -l)"
echo ""
