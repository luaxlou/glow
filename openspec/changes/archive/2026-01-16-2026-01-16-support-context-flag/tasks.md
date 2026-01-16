## 1. CLI Implementation
- [ ] 1.1 Add persistent `--context` flag to `rootCmd` in `cmd/glow/cmd/root.go`.
- [ ] 1.2 Update `ensureConfig` (or equivalent logic) in `cmd/glow/cmd/root.go` to use the flag value if present.
- [ ] 1.3 Verify invalid context names return an appropriate error.

## 2. Verification
- [ ] 2.1 Verify `glow get app --context=dev` uses the 'dev' context settings.
- [ ] 2.2 Verify `glow get app` (without flag) continues to use the active context.
