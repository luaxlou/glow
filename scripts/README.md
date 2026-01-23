# Glow Scripts

This directory contains utility scripts for Glow development and release management.

## Quick Start

### Complete Release Workflow (One Command)

```bash
# Everything in one command!
./scripts/release.sh --version v1.0.0-beta.6
```

That's it! This single command will:
1. ✅ Build binaries for all platforms
2. ✅ Commit and push changes (if needed)
3. ✅ Create GitHub release
4. ✅ Upload all assets

### Advanced Usage

```bash
# Release with custom notes
./scripts/release.sh \
  --version v1.0.0 \
  --title "Version 1.0.0" \
  --notes /tmp/release_notes.md

# Pre-release
./scripts/release.sh \
  --version v1.0.0-beta.7 \
  --pre-release

# Skip commit (if already pushed)
./scripts/release.sh \
  --version v1.0.0-beta.7 \
  --skip-commit

# Use existing binaries (already built)
./scripts/release.sh \
  --version v1.0.0 \
  --skip-build
```

---

## Release Script

**`release.sh`** - Complete release automation

### What It Does

The `release.sh` script handles the entire release process:

1. **Commit & Push** (optional)
   - Checks for uncommitted changes
   - Prompts for commit message if needed
   - Pushes to GitHub
   - Creates version tag

2. **Build Binaries** (optional)
   - Compiles for 4 platforms: darwin/amd64, darwin/arm64, linux/amd64, linux/arm64
   - Builds both glow-server and glow CLI
   - Generates SHA256 checksums

3. **Create Release**
   - Creates GitHub release with tag
   - Generates release notes from git log (or custom notes)
   - Supports pre-release and draft releases

4. **Upload Assets**
   - Uploads all compiled binaries
   - Includes SHA256SUMS.txt

### Usage

```bash
./scripts/release.sh [OPTIONS]
```

### Options

| Option | Description |
|--------|-------------|
| `-v, --version VERSION` | Version tag (required, e.g., v1.0.0-beta.6) |
| `-t, --title TITLE` | Release title (optional) |
| `-n, --notes FILE` | Release notes file in markdown (optional) |
| `-p, --pre-release` | Mark as pre-release |
| `-d, --draft` | Create as draft |
| `--skip-build` | Skip building binaries (use existing ./dist) |
| `--skip-commit` | Skip commit/push (assume already done) |
| `-h, --help` | Show help message |

### Examples

#### Standard Release

```bash
./scripts/release.sh --version v1.0.0-beta.6
```

**Output**:
```
[INFO] Glow Complete Release
[INFO] Version:     v1.0.0-beta.6
[INFO] Pre-release: false

[STEP] Step 1: Commit and Push
[INFO] Pushing to GitHub...

[STEP] Step 2: Build Binaries
[INFO] Building glow-server-darwin-arm64...
[INFO] Building glow-darwin-arm64...
...
[INFO] Build completed!

[STEP] Step 3: Create GitHub Release
[INFO] Creating GitHub release...

[STEP] Step 4: Upload Release Assets
[INFO] Files to upload: 8
[INFO] Uploading glow-server-darwin-arm64...
✓ Uploaded: glow-server-darwin-arm64
...

[INFO] Release Completed Successfully!
[INFO] URL: https://github.com/luaxlou/glow/releases/tag/v1.0.0-beta.6
```

#### Pre-release with Custom Notes

```bash
# Prepare release notes
cat > /tmp/v1.0.0-beta.7.md << EOF
# v1.0.0-beta.7

## New Features
- Feature 1
- Feature 2

## Bug Fixes
- Fix 1
EOF

# Create release
./scripts/release.sh \
  --version v1.0.0-beta.7 \
  --title "Beta 7 - New Features" \
  --notes /tmp/v1.0.0-beta.7.md \
  --pre-release
```

#### Incremental Release (Already Committed)

```bash
# If you've already pushed changes and just need to build/upload
./scripts/release.sh \
  --version v1.0.0-beta.7 \
  --skip-commit
```

#### Use Existing Binaries

```bash
# If you've already built binaries in ./dist
./scripts/release.sh \
  --version v1.0.0 \
  --skip-build
```

### Workflow Options

The script provides flexibility for different workflows:

1. **Complete Release** (default)
   ```
   Commit → Push → Build → Create Release → Upload
   ```

2. **Skip Commit** (use `--skip-commit`)
   ```
   Build → Create Release → Upload
   ```

3. **Skip Build** (use `--skip-build`)
   ```
   Create Release → Upload
   ```

4. **Skip Both** (use `--skip-commit --skip-build`)
   ```
   Create Release → Upload (fastest if everything is ready)
   ```

---

## Individual Build Tools

The following scripts are also available for individual use:

### `build.sh`

Separate build script if you only want to compile binaries:

```bash
./scripts/build.sh
```

**Output**: Creates binaries in `./dist/`

### `upload-assets.sh`

Separate upload script if you only want to upload:

```bash
./scripts/upload-assets.sh v1.0.0-beta.6 ./dist
```

**Note**: Most users should use `release.sh` instead, which combines all steps.

---

## Supported Platforms

Binaries are built for the following platforms:

| Platform | Architecture | Binary Name |
|----------|-------------|-------------|
| macOS | Intel (amd64) | `glow-server-darwin-amd64`, `glow-darwin-amd64` |
| macOS | Apple Silicon (arm64) | `glow-server-darwin-arm64`, `glow-darwin-arm64` |
| Linux | Intel (amd64) | `glow-server-linux-amd64`, `glow-linux-amd64` |
| Linux | ARM (arm64) | `glow-server-linux-arm64`, `glow-linux-arm64` |

---

## Requirements

### For Building

- **Go** 1.18+ (for compiling binaries)
- **Git** (for version control)

### For Releasing

- **GitHub CLI** (`gh`) - [Install](https://cli.github.com/)
- Authenticated GitHub session (`gh auth login`)

### Check Prerequisites

```bash
# Check Go
go version

# Check GitHub CLI
gh version

# Check authentication
gh auth status
```

---

## Complete Example

Here's a complete example of creating a release from scratch:

```bash
# 1. Make your changes
vim cmd/glow/cmd/apply.go

# 2. Test locally
go test ./...
go run cmd/glow/main.go version

# 3. Run the release script
./scripts/release.sh \
  --version v1.0.0-beta.7 \
  --title "Beta 7 - Bug Fixes" \
  --pre-release

# That's it! The script will:
# - Prompt for commit message if needed
# - Build all binaries
# - Create GitHub release
# - Upload everything
# - Show you the release URL
```

---

## Troubleshooting

### "Tag already exists"

```bash
# Delete local and remote tag
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0
```

### "Go not installed"

```bash
# Install Go
# On macOS:
brew install go

# On Linux:
# See https://golang.org/dl/
```

### "Not authenticated with gh"

```bash
# Authenticate
gh auth login
```

### Build failed

```bash
# Check Go version
go version

# Clean and retry
rm -rf ./dist
./scripts/release.sh --version v1.0.0
```

---

## Best Practices

### 1. Version Numbering

Follow semantic versioning:
- `v1.0.0` - Stable release
- `v1.0.0-beta.1` - Pre-release
- `v1.0.0-rc.1` - Release candidate
- `v1.0.0-alpha.1` - Alpha release

### 2. Release Notes

Always provide meaningful release notes:
```bash
./scripts/release.sh \
  --version v1.0.0 \
  --notes RELEASE_NOTES.md
```

### 3. Test Before Release

Always test before creating a release:
```bash
# Run tests
go test ./...

# Build locally
go build -o test-binary ./cmd/glow

# Test the binary
./test-binary version
```

### 4. Pre-release Testing

Use pre-release for beta testing:
```bash
./scripts/release.sh \
  --version v1.0.0-beta.7 \
  --pre-release
```

Then promote to stable:
```bash
./scripts/release.sh \
  --version v1.0.0 \
  --title "Stable Release 1.0.0"
```

---

## Installation Script

Users can install Glow using the installation script:

```bash
curl -fsSL "https://raw.githubusercontent.com/luaxlou/glow/main/scripts/install-local.sh" | bash
```

This script downloads the appropriate binaries for their platform from the latest GitHub release.

---

## Contributing

When adding new scripts:

1. Make them executable: `chmod +x script.sh`
2. Add usage documentation
3. Include error handling
4. Follow existing code style
5. Update this README

---

## Related Commands

```bash
# List all releases
gh release list

# View specific release
gh release view v1.0.0

# Download release assets
gh release download v1.0.0

# Delete a release
gh release delete v1.0.0 --yes
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0

# Edit release notes
gh release edit v1.0.0 --notes-file new-notes.md
```

---

## Summary

**For most users, just run:**

```bash
./scripts/release.sh --version v1.0.0
```

That one command handles everything!
