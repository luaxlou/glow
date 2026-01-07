## 1. Documentation
- [ ] 1.1 Update `app-management` spec.
- [ ] 1.2 Update `manifest-application` spec.

## 2. Implementation
- [ ] 2.1 Rename `cmd/glow/cmd/deploy.go` to `cmd/glow/cmd/app.go` and update commands to `app`.
- [ ] 2.2 Update `pkg/api/types.go` (Rename `Deployment` to `App`, verify `App` struct).
- [ ] 2.3 Update `pkg/manifest/parser.go` to prefer `kind: App`.
- [ ] 2.4 Verify `glow get app` output.
