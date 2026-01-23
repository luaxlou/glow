#!/bin/bash
# Glow Release Script
# Automates the process of committing changes, pushing to GitHub, and creating a release

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    print_error "GitHub CLI (gh) is not installed"
    print_info "Install it from: https://cli.github.com/"
    exit 1
fi

# Check if user is authenticated
if ! gh auth status &> /dev/null; then
    print_error "Not authenticated with GitHub CLI"
    print_info "Run: gh auth login"
    exit 1
fi

# Parse arguments
VERSION=""
RELEASE_TITLE=""
RELEASE_NOTES=""
PRE_RELEASE=false
DRAFT=false

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
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -v, --version VERSION       Version tag (e.g., v1.0.0-beta.6)"
            echo "  -t, --title TITLE          Release title"
            echo "  -n, --notes FILE           Release notes file (markdown)"
            echo "  -p, --pre-release          Mark as pre-release"
            echo "  -d, --draft                Create as draft"
            echo "  -h, --help                 Show this help message"
            echo ""
            echo "Example:"
            echo "  $0 --version v1.0.0-beta.6 --title 'Beta 6' --notes release.md"
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

# Check if there are uncommitted changes
if [ -n "$(git status --porcelain)" ]; then
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
fi

# Check if tag already exists
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    print_error "Tag $VERSION already exists"
    print_info "Delete it first with: git tag -d $VERSION && git push origin :refs/tags/$VERSION"
    exit 1
fi

print_info "=========================================="
print_info "Glow Release Process"
print_info "=========================================="
echo ""
print_info "Version:     $VERSION"
print_info "Title:       ${RELEASE_TITLE:-$VERSION}"
print_info "Pre-release: $PRE_RELEASE"
print_info "Draft:       $DRAFT"
echo ""

# Confirm before proceeding
read -p "Continue with release? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_info "Release cancelled"
    exit 0
fi

# Push to GitHub
print_info "Pushing to GitHub..."
git push origin main

if [ $? -ne 0 ]; then
    print_error "Failed to push to GitHub"
    exit 1
fi

print_info "Push successful"

# Create tag
print_info "Creating tag: $VERSION"
git tag -a "$VERSION" -m "Release $VERSION"

# Push tag
print_info "Pushing tag to GitHub..."
git push origin "$VERSION"

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
    if [ -n "$LAST_TAG" ]; then
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
        NOTES_FLAG="--notes 'Release $VERSION'"
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

if [ $? -eq 0 ]; then
    print_info "=========================================="
    print_info "Release created successfully!"
    print_info "=========================================="
    print_info "Version: $VERSION"
    print_info "URL: $(gh release view $VERSION --json url -q .url)"
    print_info "=========================================="

    # Clean up temp file if created
    if [ -f "$TEMP_NOTES" ]; then
        rm -f "$TEMP_NOTES"
    fi
else
    print_error "Failed to create release"
    exit 1
fi
