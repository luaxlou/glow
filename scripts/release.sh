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
BUMP_TYPE=""

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
        --bump)
            BUMP_TYPE="$2"
            shift 2
            ;;
        --help|-h)
            echo "Glow Complete Release Script"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "This script automates the entire release process:"
            echo "  1. Auto-detects current version and calculates new version (if not provided)"
            echo "  2. Updates VERSION file and commits changes (unless --skip-commit)"
            echo "  3. Creates git tag and pushes to GitHub (unless --skip-commit)"
            echo "  4. Builds binaries for all platforms (unless --skip-build)"
            echo "  5. Creates GitHub release"
            echo "  6. Uploads binaries as release assets"
            echo ""
            echo "Options:"
            echo "  -v, --version VERSION       Version tag (optional, auto-detected if not provided)"
            echo "  --bump TYPE                Bump version type: major, minor, patch, prerelease (default: auto)"
            echo "  -t, --title TITLE          Release title (optional)"
            echo "  -n, --notes FILE           Release notes file (markdown)"
            echo "  -p, --pre-release          Mark as pre-release"
            echo "  -d, --draft                Create as draft"
            echo "  --skip-build               Skip building binaries (use existing ./dist)"
            echo "  --skip-commit              Skip commit/push (assume already done)"
            echo "  -h, --help                 Show this help message"
            echo ""
            echo "Examples:"
            echo "  # Auto-detect and bump version (default: patch or prerelease)"
            echo "  $0"
            echo ""
            echo "  # Auto-bump with specific type"
            echo "  $0 --bump minor"
            echo ""
            echo "  # Complete release workflow with explicit version"
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

# Function to get latest version from git tags or VERSION file
get_latest_version() {
    # Try to get latest tag from git (sorted by version)
    local latest_tag=""
    
    # Get all version tags and sort them properly
    if git rev-parse --git-dir > /dev/null 2>&1; then
        # Get all tags that match version pattern and sort by version
        latest_tag=$(git tag -l "v[0-9]*" | sort -V | tail -1)
        
        # If no v-prefixed tags, try without prefix
        if [ -z "$latest_tag" ]; then
            latest_tag=$(git tag -l "[0-9]*" | sort -V | tail -1)
        fi
    fi
    
    if [ -n "$latest_tag" ]; then
        # Remove 'v' prefix if present
        echo "${latest_tag#v}"
    elif [ -f "VERSION" ]; then
        # Fallback to VERSION file
        local version=$(cat VERSION 2>/dev/null | tr -d '[:space:]')
        if [ -n "$version" ]; then
            echo "$version"
        else
            echo "0.0.0"
        fi
    else
        echo "0.0.0"
    fi
}

# Function to bump version
bump_version() {
    local current_version=$1
    local bump_type=${2:-"auto"}
    
    # Remove 'v' prefix if present
    current_version=${current_version#v}
    
    # Parse version components
    local major minor patch prerelease=""
    
    if [[ $current_version =~ ^([0-9]+)\.([0-9]+)\.([0-9]+)(-.*)?$ ]]; then
        major=${BASH_REMATCH[1]}
        minor=${BASH_REMATCH[2]}
        patch=${BASH_REMATCH[3]}
        prerelease=${BASH_REMATCH[4]}
        prerelease=${prerelease#-}  # Remove leading dash
    else
        print_error "Invalid version format: $current_version"
        exit 1
    fi
    
    # Auto-detect bump type if not specified
    if [ "$bump_type" = "auto" ] || [ -z "$bump_type" ]; then
        if [ -n "$prerelease" ]; then
            # If current version is prerelease, bump prerelease number
            bump_type="prerelease"
        else
            # Otherwise, bump patch version
            bump_type="patch"
        fi
    fi
    
    # Perform bump
    case $bump_type in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            prerelease=""
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            prerelease=""
            ;;
        patch)
            patch=$((patch + 1))
            prerelease=""
            ;;
        prerelease)
            if [ -n "$prerelease" ]; then
                # Extract prerelease type and number (e.g., "beta.10" -> "beta" and "10")
                if [[ $prerelease =~ ^([a-zA-Z]+)\.([0-9]+)$ ]]; then
                    local pre_type=${BASH_REMATCH[1]}
                    local pre_num=${BASH_REMATCH[2]}
                    pre_num=$((pre_num + 1))
                    prerelease="${pre_type}.${pre_num}"
                elif [[ $prerelease =~ ^([a-zA-Z]+)$ ]]; then
                    # If no number, add .1
                    prerelease="${prerelease}.1"
                else
                    # Unknown format, append .1
                    prerelease="${prerelease}.1"
                fi
            else
                # No prerelease, create beta.1
                prerelease="beta.1"
            fi
            ;;
        *)
            print_error "Invalid bump type: $bump_type"
            print_info "Valid types: major, minor, patch, prerelease"
            exit 1
            ;;
    esac
    
    # Construct new version
    local new_version="${major}.${minor}.${patch}"
    if [ -n "$prerelease" ]; then
        new_version="${new_version}-${prerelease}"
    fi
    
    echo "$new_version"
}

# Auto-detect version if not provided
AUTO_DETECTED=false
if [ -z "$VERSION" ]; then
    AUTO_DETECTED=true
    print_info "=========================================="
    print_info "Auto-detecting version..."
    print_info "=========================================="
    echo ""
    
    LATEST_VERSION=$(get_latest_version)
    
    if [ "$LATEST_VERSION" = "0.0.0" ]; then
        print_warn "No existing version found in git tags or VERSION file"
        print_info "Starting with initial version: v0.1.0-beta.1"
        VERSION="v0.1.0-beta.1"
        PRE_RELEASE=true
    else
        print_info "Current latest version: $LATEST_VERSION"
        
        # Determine bump type for display
        if [ -z "$BUMP_TYPE" ] || [ "$BUMP_TYPE" = "auto" ]; then
            if [[ $LATEST_VERSION =~ - ]]; then
                BUMP_TYPE="prerelease"
            else
                BUMP_TYPE="patch"
            fi
        fi
        
        print_info "Bump type: $BUMP_TYPE"
        NEW_VERSION=$(bump_version "$LATEST_VERSION" "$BUMP_TYPE")
        VERSION="v${NEW_VERSION}"
        
        # Auto-detect if it's a prerelease
        if [[ $NEW_VERSION =~ - ]]; then
            PRE_RELEASE=true
        fi
        
        print_info "New version: $VERSION"
    fi
    echo ""
fi

# Validate version format
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
    if [ "$AUTO_DETECTED" = true ]; then
        print_info "Version:     $VERSION (auto-detected)"
    else
        print_info "Version:     $VERSION"
    fi
    print_info "Title:       ${RELEASE_TITLE:-$VERSION}"
    print_info "Pre-release: $PRE_RELEASE"
    print_info "Draft:       $DRAFT"
    print_info "Skip Build:  $SKIP_BUILD"
    print_info "Skip Commit: $SKIP_COMMIT"
    echo ""

    check_prerequisites

    # Step 1: Update VERSION file, commit, tag and push (if not skipped)
    if [ "$SKIP_COMMIT" = false ]; then
        print_step "Step 1: Update VERSION, Commit, Tag and Push"
        echo ""

        # Check if tag already exists
        if git rev-parse "$VERSION" >/dev/null 2>&1; then
            print_error "Tag $VERSION already exists"
            print_info "Delete it first with: git tag -d $VERSION && git push origin :refs/tags/$VERSION"
            exit 1
        fi

        # Check if we're on the correct branch
        CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
        if [ "$CURRENT_BRANCH" != "main" ] && [ "$CURRENT_BRANCH" != "master" ]; then
            print_warn "Current branch is '$CURRENT_BRANCH', not 'main' or 'master'"
            print_info "Continuing anyway..."
        fi

        # Update VERSION file
        VERSION_NO_V=${VERSION#v}
        print_info "Updating VERSION file to $VERSION_NO_V..."
        echo "$VERSION_NO_V" > VERSION

        # Check if there are other uncommitted changes
        UNCOMMITTED=$(git status --porcelain | grep -v "^ M VERSION$" | grep -v "^??" || true)
        if [ -n "$UNCOMMITTED" ]; then
            print_info "Staging all changes including VERSION file..."
            git add -A

            # Generate commit message from changes
            if [ "$AUTO_DETECTED" = true ]; then
                COMMIT_MSG="chore: release $VERSION"
            else
                COMMIT_MSG="chore: update VERSION to $VERSION_NO_V"
            fi

            print_info "Creating commit: $COMMIT_MSG"
            git commit -m "$COMMIT_MSG"
        else
            # Only VERSION file changed or no changes
            if git diff --quiet VERSION 2>/dev/null || [ -n "$(git status --porcelain VERSION)" ]; then
                print_info "Committing VERSION file..."
                git add VERSION
                git commit -m "chore: update VERSION to $VERSION_NO_V"
            else
                print_info "No changes to commit"
            fi
        fi

        # Create tag
        print_info "Creating tag $VERSION..."
        git tag -a "$VERSION" -m "Release $VERSION"

        if [ $? -ne 0 ]; then
            print_error "Failed to create tag"
            exit 1
        fi

        print_info "Tag created successfully"

        # Push commits and tags to GitHub
        print_info "Pushing commits to GitHub..."
        git push origin "$CURRENT_BRANCH"

        if [ $? -ne 0 ]; then
            print_error "Failed to push commits to GitHub"
            exit 1
        fi

        print_info "Pushing tag $VERSION to GitHub..."
        git push origin "$VERSION"

        if [ $? -ne 0 ]; then
            print_error "Failed to push tag to GitHub"
            exit 1
        fi

        print_info "Push successful"
        echo ""
    else
        print_info "Skipping commit/tag/push (--skip-commit)"
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
            "darwin/arm64"
            "linux/amd64"
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
        # Get the previous tag (before the current one)
        ALL_TAGS=($(git tag -l "v*" | sort -V))
        LAST_TAG=""
        
        # Find the tag before the current version
        for tag in "${ALL_TAGS[@]}"; do
            if [ "$tag" != "$VERSION" ]; then
                LAST_TAG="$tag"
            else
                break
            fi
        done
        
        # If no previous tag found, try without v prefix
        if [ -z "$LAST_TAG" ]; then
            ALL_TAGS=($(git tag -l | sort -V))
            for tag in "${ALL_TAGS[@]}"; do
                if [ "$tag" != "$VERSION" ]; then
                    LAST_TAG="$tag"
                else
                    break
                fi
            done
        fi
        
        TEMP_NOTES=$(mktemp)
        if [ -n "$LAST_TAG" ] && [ "$LAST_TAG" != "$VERSION" ]; then
            print_info "Generating release notes from git log ($LAST_TAG..$VERSION)"
            cat > "$TEMP_NOTES" << EOF
# $VERSION

${RELEASE_TITLE:-"Release $VERSION"}

## Changes since $LAST_TAG

$(git log "$LAST_TAG..$VERSION" --pretty=format:"- %s (%h)" --reverse 2>/dev/null || git log "$LAST_TAG..HEAD" --pretty=format:"- %s (%h)" --reverse)

EOF
        else
            print_info "Generating basic release notes (no previous tag found)"
            cat > "$TEMP_NOTES" << EOF
# $VERSION

${RELEASE_TITLE:-"Release $VERSION"}

## Changes

$(git log --pretty=format:"- %s (%h)" --reverse -20)

EOF
        fi
        NOTES_FLAG="--notes-file $TEMP_NOTES"
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
