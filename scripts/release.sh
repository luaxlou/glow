#!/bin/bash
# Glow Complete Release Script
# Automates the entire release process: build, commit, push, create release, upload assets

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

# Check prerequisites
check_prerequisites() {
    local missing=0

    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        missing=1
    fi

    if ! command -v gh &> /dev/null; then
        print_error "GitHub CLI (gh) is not installed"
        print_info "Install it from: https://cli.github.com/"
        missing=1
    fi

    if ! gh auth status &> /dev/null; then
        print_error "Not authenticated with GitHub CLI"
        print_info "Run: gh auth login"
        missing=1
    fi

    if [ $missing -eq 1 ]; then
        exit 1
    fi
}

# Parse arguments
VERSION=""
RELEASE_TITLE=""
RELEASE_NOTES=""
PRE_RELEASE=false
DRAFT=false
SKIP_BUILD=false
SKIP_COMMIT=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --version|-v)
            VERSION="$2"
            shift 2
            ;;
        --title|-t)
            RELEASE_TITLE="$2"
            shift 2
            ;;
        --notes|-n)
            RELEASE_NOTES="$2"
            shift 2
            ;;
        --pre-release|-p)
            PRE_RELEASE=true
            shift
            ;;
        --draft|-d)
            DRAFT=true
            shift
            ;;
        --skip-build)
            SKIP_BUILD=true
            shift
            ;;
        --skip-commit)
            SKIP_COMMIT=true
            shift
            ;;
        --help|-h)
            echo "Glow Complete Release Script"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "This script automates the entire release process:"
            echo "  1. Builds binaries for all platforms (unless --skip-build)"
            echo "  2. Commits and pushes changes (unless --skip-commit)"
            echo "  3. Creates GitHub release"
            echo "  4. Uploads binaries as release assets"
            echo ""
            echo "Options:"
            echo "  -v, --version VERSION       Version tag (required, e.g., v1.0.0-beta.6)"
            echo "  -t, --title TITLE          Release title (optional)"
            echo "  -n, --notes FILE           Release notes file (markdown)"
            echo "  -p, --pre-release          Mark as pre-release"
            echo "  -d, --draft                Create as draft"
            echo "  --skip-build               Skip building binaries (use existing ./dist)"
            echo "  --skip-commit              Skip commit/push (assume already done)"
            echo "  -h, --help                 Show this help message"
            echo ""
            echo "Examples:"
            echo "  # Complete release workflow"
            echo "  $0 --version v1.0.0-beta.6"
            echo ""
            echo "  # Release with custom notes"
            echo "  $0 --version v1.0.0 --title 'Version 1.0.0' --notes release.md"
            echo ""
            echo "  # Pre-release (skip commit if already done)"
            echo "  $0 --version v1.0.0-beta.7 --pre-release --skip-commit"
            echo ""
            echo "  # Use existing binaries"
            echo "  $0 --version v1.0.0 --skip-build"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use --help to see available options"
            exit 1
            ;;
    esac
done

# Validate version
if [ -z "$VERSION" ]; then
    print_error "Version is required (use --version)"
    exit 1
fi

if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+ ]]; then
    print_error "Invalid version format: $VERSION"
    print_info "Expected format: v1.0.0, v1.0.0-beta.6, etc."
    exit 1
fi

# Main workflow
main() {
    print_info "=========================================="
    print_info "Glow Complete Release"
    print_info "=========================================="
    echo ""
    print_info "Version:     $VERSION"
    print_info "Title:       ${RELEASE_TITLE:-$VERSION}"
    print_info "Pre-release: $PRE_RELEASE"
    print_info "Draft:       $DRAFT"
    print_info "Skip Build:  $SKIP_BUILD"
    print_info "Skip Commit: $SKIP_COMMIT"
    echo ""

    check_prerequisites

    # Step 1: Update VERSION file and commit (if not skipped)
    if [ "$SKIP_COMMIT" = false ]; then
        print_step "Step 1: Update VERSION and Commit"
        echo ""

        # Update VERSION file
        VERSION_NO_V=${VERSION#v}
        print_info "Updating VERSION file to $VERSION_NO_V..."
        echo "$VERSION_NO_V" > VERSION

        # Check if there are other uncommitted changes
        if [ -n "$(git status --porcelain | grep -v VERSION)" ]; then
            print_warn "You have uncommitted changes"
            read -p "Do you want to commit them now? (y/n) " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                print_info "Staging all changes..."
                git add -A

                print_info "Please enter commit message:"
                read -e COMMIT_MSG

                if [ -z "$COMMIT_MSG" ]; then
                    print_error "Commit message cannot be empty"
                    exit 1
                fi

                print_info "Creating commit..."
                git commit -m "$COMMIT_MSG"
            else
                print_error "Cannot proceed with uncommitted changes"
                exit 1
            fi
        else
            # Only VERSION file changed
            print_info "Committing VERSION file..."
            git add VERSION
            git commit -m "chore: update VERSION to $VERSION_NO_V"
        fi

        # Check if tag already exists
        if git rev-parse "$VERSION" >/dev/null 2>&1; then
            print_error "Tag $VERSION already exists"
            print_info "Delete it first with: git tag -d $VERSION && git push origin :refs/tags/$VERSION"
            exit 1
        fi

        # Push to GitHub
        print_info "Pushing to GitHub..."
        git push origin main

        if [ $? -ne 0 ]; then
            print_error "Failed to push to GitHub"
            exit 1
        fi

        print_info "Push successful"
        echo ""
    else
        print_info "Skipping commit/push (--skip-commit)"
        echo ""
    fi

    # Step 2: Build binaries
    if [ "$SKIP_BUILD" = false ]; then
        print_step "Step 2: Build Binaries"
        echo ""

        DIST_DIR="./dist"
        rm -rf "$DIST_DIR"
        mkdir -p "$DIST_DIR"

        print_info "Building for all platforms..."

        # Supported platforms
        PLATFORMS=(
            "darwin/amd64"
            "darwin/arm64"
            "linux/amd64"
            "linux/arm64"
        )

        # Get version from git tag
        VERSION_NO_V=${VERSION#v}

        # Build function
        build_binary() {
            local goos=$1
            local goarch=$2
            local binary_name=$3

            local output_name="${binary_name}-${goos}-${goarch}"
            local output_path="${DIST_DIR}/${output_name}"

            print_info "Building ${output_name}..."

            GOOS=$goos GOARCH=$goarch go build \
                -ldflags="-X 'main.Version=$VERSION_NO_V' -X 'main.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" \
                -o "$output_path" \
                ./"$(get_binary_path $binary_name)"

            # Generate checksum
            if command -v sha256sum &> /dev/null; then
                checksum=$(sha256sum "$output_path" | awk '{print $1}')
            elif command -v shasum &> /dev/null; then
                checksum=$(shasum -a 256 "$output_path" | awk '{print $1}')
            fi

            echo "${checksum}  ${output_name}" >> "${DIST_DIR}/SHA256SUMS.txt"
        }

        get_binary_path() {
            case $1 in
                "glow-server")
                    echo "cmd/glow-server"
                    ;;
                "glow")
                    echo "cmd/glow"
                    ;;
            esac
        }

        # Build for all platforms
        for platform in "${PLATFORMS[@]}"; do
            IFS='/' read -r goos goarch <<< "$platform"

            build_binary "$goos" "$goarch" "glow-server"
            build_binary "$goos" "$goarch" "glow"
        done

        # Sort checksums file
        if [ -f "${DIST_DIR}/SHA256SUMS.txt" ]; then
            sort "${DIST_DIR}/SHA256SUMS.txt" -o "${DIST_DIR}/SHA256SUMS.txt"
        fi

        print_info "Build completed!"
        print_info "Binaries: $(ls -1 "$DIST_DIR" | wc -l) files"
        echo ""
    else
        print_info "Skipping build (--skip-build)"
        print_info "Using existing ./dist directory"

        if [ ! -d "./dist" ]; then
            print_error "dist directory not found"
            exit 1
        fi
        echo ""
    fi

    # Step 3: Create GitHub release
    print_step "Step 3: Create GitHub Release"
    echo ""

    # Prepare release notes
    if [ -n "$RELEASE_NOTES" ]; then
        if [ ! -f "$RELEASE_NOTES" ]; then
            print_error "Release notes file not found: $RELEASE_NOTES"
            exit 1
        fi
        NOTES_FLAG="--notes-file $RELEASE_NOTES"
    else
        # Generate basic release notes from git log
        LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
        if [ -n "$LAST_TAG" ] && [ "$LAST_TAG" != "$VERSION" ]; then
            print_info "Generating release notes from git log ($LAST_TAG..HEAD)"
            TEMP_NOTES=$(mktemp)
            cat > "$TEMP_NOTES" << EOF
# $VERSION

${RELEASE_TITLE:-"Release $VERSION"}

## Changes since $LAST_TAG

$(git log "$LAST_TAG..HEAD" --pretty=format:"- %s" --reverse)

EOF
            NOTES_FLAG="--notes-file $TEMP_NOTES"
        else
            TEMP_NOTES=$(mktemp)
            echo "# $VERSION\n\n${RELEASE_TITLE:-"Release $VERSION"}" > "$TEMP_NOTES"
            NOTES_FLAG="--notes-file $TEMP_NOTES"
        fi
    fi

    # Build release command
    RELEASE_CMD="gh release create $VERSION"
    RELEASE_CMD="$RELEASE_CMD --title '${RELEASE_TITLE:-$VERSION}'"

    if [ "$PRE_RELEASE" = true ]; then
        RELEASE_CMD="$RELEASE_CMD --prerelease"
    fi

    if [ "$DRAFT" = true ]; then
        RELEASE_CMD="$RELEASE_CMD --draft"
    fi

    if [ -n "$NOTES_FLAG" ]; then
        RELEASE_CMD="$RELEASE_CMD $NOTES_FLAG"
    fi

    print_info "Creating GitHub release..."
    eval $RELEASE_CMD

    if [ $? -ne 0 ]; then
        print_error "Failed to create release"
        exit 1
    fi

    print_info "Release created successfully"
    echo ""

    # Step 4: Upload assets
    print_step "Step 4: Upload Release Assets"
    echo ""

    DIST_DIR="${DIST_DIR:-./dist}"
    FILES=($(find "$DIST_DIR" -type f -not -name "*.txt"))

    if [ ${#FILES[@]} -eq 0 ]; then
        print_error "No files found in $DIST_DIR"
        exit 1
    fi

    print_info "Files to upload: ${#FILES[@]}"
    for file in "${FILES[@]}"; do
        print_info "  - $(basename $file) ($(du -h $file | cut -f1))"
    done
    echo ""

    # Upload each file
    for file in "${FILES[@]}"; do
        filename=$(basename "$file")
        print_info "Uploading ${filename}..."

        if gh release upload "$VERSION" "$file" --clobber; then
            print_info "✓ Uploaded: ${filename}"
        else
            print_error "✗ Failed: ${filename}"
        fi
    done

    # Upload checksums file
    if [ -f "${DIST_DIR}/SHA256SUMS.txt" ]; then
        print_info "Uploading SHA256SUMS.txt..."
        if gh release upload "$VERSION" "${DIST_DIR}/SHA256SUMS.txt" --clobber; then
            print_info "✓ Uploaded: SHA256SUMS.txt"
        fi
    fi

    echo ""

    # Cleanup temp file if created
    if [ -f "$TEMP_NOTES" ]; then
        rm -f "$TEMP_NOTES"
    fi

    # Done
    print_info "=========================================="
    print_info "Release Completed Successfully!"
    print_info "=========================================="
    echo ""
    print_info "Version: $VERSION"
    print_info "URL: $(gh release view $VERSION --json url -q .url)"
    echo ""
    print_info "Assets uploaded:"
    gh release view "$VERSION" --json assets --jq '.assets | length' | xargs echo "  Total:"
    echo ""
}

# Run main function
main
