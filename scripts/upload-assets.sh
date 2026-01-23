#!/bin/bash
# Upload Release Assets Script
# Uploads compiled binaries to a GitHub release

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

# Check arguments
if [ $# -lt 1 ]; then
    echo "Usage: $0 <version> [dist-directory]"
    echo ""
    echo "Example:"
    echo "  $0 v1.0.0-beta.6"
    echo "  $0 v1.0.0-beta.6 ./dist"
    exit 1
fi

VERSION="$1"
DIST_DIR="${2:-./dist}"

# Remove 'v' prefix if user forgot it
if [[ ! $VERSION =~ ^v ]]; then
    VERSION="v${VERSION}"
fi

if [ ! -d "$DIST_DIR" ]; then
    print_error "Distribution directory not found: $DIST_DIR"
    exit 1
fi

print_info "=========================================="
print_info "Upload Release Assets"
print_info "=========================================="
echo ""
print_info "Version: $VERSION"
print_info "Source: $DIST_DIR"
echo ""

# List files to upload
FILES=($(find "$DIST_DIR" -type f -not -name "*.txt"))

if [ ${#FILES[@]} -eq 0 ]; then
    print_error "No files found in $DIST_DIR"
    print_info "Run ./scripts/build.sh first"
    exit 1
fi

print_info "Files to upload:"
for file in "${FILES[@]}"; do
    print_info "  - $(basename $file) ($(du -h $file | cut -f1))"
done
echo ""

# Confirm
read -p "Continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_info "Cancelled"
    exit 0
fi

# Upload each file
print_step "Uploading files to GitHub release..."

for file in "${FILES[@]}"; do
    filename=$(basename "$file")

    print_step "Uploading ${filename}..."

    if gh release upload "$VERSION" "$file" --clobber; then
        print_info "✓ Uploaded: ${filename}"
    else
        print_error "✗ Failed: ${filename}"
    fi
done

# Upload checksums file separately
if [ -f "${DIST_DIR}/SHA256SUMS.txt" ]; then
    print_step "Uploading SHA256SUMS.txt..."
    if gh release upload "$VERSION" "${DIST_DIR}/SHA256SUMS.txt" --clobber; then
        print_info "✓ Uploaded: SHA256SUMS.txt"
    else
        print_error "✗ Failed: SHA256SUMS.txt"
    fi
fi

echo ""
print_info "=========================================="
print_info "Upload completed!"
print_info "=========================================="
echo ""

# Show release assets
print_info "Release assets:"
gh release view "$VERSION" --json assets --jq '.assets[].name'

echo ""
print_info "Release URL:"
gh release view "$VERSION" --json url -q .url

echo ""
