## 1. Implementation
- [ ] 1.1 Update `StartApp` to check status and return success if already RUNNING (Idempotency).
- [ ] 1.2 Refactor port allocation: use config-provided port if available, otherwise allocate.
- [ ] 1.3 Implement `RestartApp` logic (Stop then Start).
- [ ] 1.4 Implement `DeleteApp` logic (Stop process, remove config/data).
- [ ] 1.5 Update status handling: `STOPPED` (manual), `EXITED` (code 0), `ERROR` (non-0).
- [ ] 1.6 Update API handlers in `internal/apiserver`.
