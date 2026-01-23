# Glow Scripts

This directory contains utility scripts for Glow development and release management.

## Release Script

**`release.sh`** - Automated release script for creating GitHub releases

### Features

- Validates Git status (checks for uncommitted changes)
- Creates Git tag
- Pushes changes and tags to GitHub
- Creates GitHub release with notes
- Supports pre-release and draft releases

### Usage

```bash
./scripts/release.sh [OPTIONS]
```

### Options

- `-v, --version VERSION` - Version tag (required, e.g., v1.0.0-beta.6)
- `-t, --title TITLE` - Release title (optional)
- `-n, --notes FILE` - Release notes file in markdown (optional)
- `-p, --pre-release` - Mark as pre-release
- `-d, --draft` - Create as draft
- `-h, --help` - Show help message

### Examples

#### Basic Release

```bash
./scripts/release.sh --version v1.0.0
```

#### Pre-release with Notes

```bash
./scripts/release.sh \
  --version v1.0.0-beta.7 \
  --title "Beta 7 - Bug Fixes" \
  --notes /tmp/release_notes.md \
  --pre-release
```

#### Draft Release

```bash
./scripts/release.sh \
  --version v1.0.0 \
  --title "Version 1.0.0" \
  --notes release-notes.md \
  --draft
```

### Workflow

1. **Prepare release notes** (optional):
   ```bash
   # Create release notes file
   cat > /tmp/v1.0.0-notes.md << EOF
   # v1.0.0

   ## New Features
   - Feature 1
   - Feature 2

   ## Bug Fixes
   - Fix 1
   EOF
   ```

2. **Run release script**:
   ```bash
   ./scripts/release.sh \
     --version v1.0.0 \
     --title "Version 1.0.0" \
     --notes /tmp/v1.0.0-notes.md
   ```

3. **Script will**:
   - Check for uncommitted changes
   - Prompt to commit if needed
   - Push changes to GitHub
   - Create and push Git tag
   - Create GitHub release

### Automatic Release Notes

If you don't provide a `--notes` file, the script will automatically generate release notes from the Git commit log since the last tag.

### Requirements

- Git
- GitHub CLI (`gh`)
- Authenticated GitHub session (`gh auth login`)

### Installation

```bash
# Install GitHub CLI (if not installed)
# On macOS:
brew install gh

# On Linux:
# See https://cli.github.com/

# Authenticate
gh auth login
```

### Color Output

The script uses colored output:
- 🟢 Green - Info messages
- 🟡 Yellow - Warnings
- 🔴 Red - Errors

### Error Handling

The script will:
- Exit on any error (`set -e`)
- Validate version format
- Check for existing tags
- Verify GitHub authentication
- Confirm before creating release

### Version Format

Versions must follow semantic versioning:
- `v1.0.0` - Standard release
- `v1.0.0-beta.1` - Pre-release
- `v1.0.0-rc.1` - Release candidate
- `v1.0.0-alpha.1` - Alpha release

### Related Commands

```bash
# List all releases
gh release list

# View specific release
gh release view v1.0.0

# Delete a release (and tag)
gh release delete v1.0.0 --yes
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0

# Edit release notes
gh release edit v1.0.0 --notes-file new-notes.md
```

### Troubleshooting

**Tag already exists**:
```bash
# Delete local and remote tag
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0
```

**Not authenticated**:
```bash
# Re-authenticate
gh auth logout
gh auth login
```

**Push failed**:
```bash
# Check remote
git remote -v

# Ensure you have access to the repository
gh repo view
```

## Future Scripts

Additional utility scripts may be added in the future:

- `build.sh` - Build all Glow components
- `test.sh` - Run comprehensive tests
- `install.sh` - Install Glow locally
- `clean.sh` - Clean build artifacts

## Contributing

When adding new scripts:

1. Make them executable: `chmod +x script.sh`
2. Add usage documentation
3. Include error handling
4. Follow existing code style
5. Update this README
