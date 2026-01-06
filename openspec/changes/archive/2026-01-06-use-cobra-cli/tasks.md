## 1. Implementation
- [ ] 1.1 Create `cmd/glow/cmd/root.go` and move `main` logic to `Execute()`.
- [ ] 1.2 Implement `setup` command (was `config` for client setup).
- [ ] 1.3 Implement `app` command group (`start`, `stop`, `restart`, `delete`, `list`, `logs`).
- [ ] 1.4 Implement `config` command group (`set`, `get`, `list`).
- [ ] 1.5 Implement `ingress` command group (`apply`, `delete`, `list`).
- [ ] 1.6 Implement `init` and `ping` commands.
- [ ] 1.7 Update `cmd/glow/main.go` to call `cmd.Execute()`.
