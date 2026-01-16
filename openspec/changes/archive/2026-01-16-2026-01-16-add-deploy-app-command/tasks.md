## 1. Server Implementation
- [ ] 1.1 Update `internal/manager` to support calculating SHA256 of binaries and storing in `AppInfo`.
- [ ] 1.2 Implement `POST /apps/upload` in `internal/apiserver` to handle file uploads (multipart).
- [ ] 1.3 Ensure `StartApp` uses the uploaded binary and updates the hash.

## 2. CLI Implementation
- [ ] 2.1 Add `deploy` verb to `cmd/glow/cmd/verbs.go`.
- [ ] 2.2 Create `cmd/glow/cmd/deploy.go` with `deploy app` subcommand.
- [ ] 2.3 Implement hash calculation (SHA256) on client.
- [ ] 2.4 Implement `getApp` logic to fetch remote hash.
- [ ] 2.5 Implement `uploadBinary` logic (multipart POST).
- [ ] 2.6 Implement the full control flow: Check Hash -> (Skip OR Upload) -> Start.

## 3. Verification
- [ ] 3.1 Verify `glow deploy app` with new binary (should upload).
- [ ] 3.2 Verify `glow deploy app` with same binary (should skip).
