## 1. Preparation & Planning
- [x] 1.1 Review current `glow init` implementation (cmd/glow/cmd/init.go)
- [x] 1.2 Review project-initialization spec requirements
- [x] 1.3 Review system-initialization spec for online script patterns
- [x] 1.4 Create proposal document

## 2. Create Online Initialization Script
- [x] 2.1 Create `scripts/init-project.sh` script with the following features:
  - [x] 2.1.1 Detect current directory and project structure
  - [x] 2.1.2 Create standard directories (cmd/, bin/, scripts/)
  - [x] 2.1.3 Generate deploy.sh script
  - [x] 2.1.4 Support AI tool integration (Claude Code skills)
  - [x] 2.1.5 Handle idempotency (skip existing files unless --force)
  - [x] 2.1.6 Support command-line arguments (--skip-ai, --force, --help)
  - [x] 2.1.7 Add colored output and progress indicators
  - [x] 2.1.8 Error handling and validation

## 3. Update CLI Code
- [x] 3.1 Remove `cmd/glow/cmd/init.go` file
- [x] 3.2 Update `cmd/glow/cmd/root.go` to remove initCmd registration
- [x] 3.3 Add deprecation notice for `glow init` (optional transition period)
  - [x] 3.3.1 Print migration message when old `glow init` is invoked
  - [x] 3.3.2 Exit with link to online script
- [x] 3.4 Test CLI build to ensure no compilation errors

## 4. Update Documentation
- [x] 4.1 Update README.md to use online script initialization
- [x] 4.2 Update getting started guides
- [x] 4.3 Update CLI manual (docs/cli_manual.md)
- [x] 4.4 Add migration notes for existing users
- [x] 4.5 Document init script usage and options

## 5. Publish & Deploy
- [x] 5.1 Upload init-project.sh to GitHub repository (scripts/)
- [x] 5.2 Add installation instructions in README
- [x] 5.3 Create raw.githubusercontent.com permalink or CDN URL
- [x] 5.4 Test script execution: `curl -fsSL <url> | bash`
- [x] 5.5 Test with various options: `curl ... | bash -s -- --skip-ai --force`

## 6. Update OpenSpec Artifacts
- [ ] 6.1 Create spec delta for project-initialization
- [ ] 6.2 Mark old requirements as REMOVED
- [ ] 6.3 Add new requirements for online script initialization
- [ ] 6.4 Validate proposal with `openspec validate --strict`

## 7. Testing & Validation
- [x] 7.1 Test online script on fresh project directory
- [x] 7.2 Test idempotency (run script twice on same project)
- [x] 7.3 Test --force flag to overwrite existing files
- [x] 7.4 Test --skip-ai flag
- [x] 7.5 Test AI tool integration (Claude Code skills)
- [x] 7.6 Test on different platforms (Linux, macOS)
- [x] 7.7 Verify generated deploy.sh works correctly

## 8. Cleanup & Archive
- [ ] 8.1 Verify all tasks completed
- [ ] 8.2 Run final validation: `openspec validate remove-glow-init-curl-init --strict`
- [ ] 8.3 Update CHANGELOG.md (after implementation)
- [ ] 8.4 Create PR for code review
- [ ] 8.5 Archive change after deployment
