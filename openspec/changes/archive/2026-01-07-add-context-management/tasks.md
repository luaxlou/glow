## 1. Documentation
- [ ] 1.1 Update `authentication` spec to include Context Management requirements.

## 2. Implementation
- [ ] 2.1 Update `cmd/glow/cmd/root.go` to support new config schema (Contexts). Migration logic for existing config?
- [ ] 2.2 Implement `cmd/glow/cmd/context.go` with `list`, `use`, `add`, `delete`.
- [ ] 2.3 Update `cmd/glow/cmd/auth.go` to use the current context from the new schema.
- [ ] 2.4 Add integration tests for context switching.
