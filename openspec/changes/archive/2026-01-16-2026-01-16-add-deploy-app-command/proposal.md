# Change: Add Deploy App Command

## Why
Users need a way to deploy and update applications safely and efficiently. The `apply` command has been removed, and a dedicated `deploy app` command is required. This command should optimize bandwidth and time by checking if the binary has changed (via hash comparison) before uploading.

## What Changes
- **CLI**: Add `glow deploy` command (replaces `glow deploy app`).
- **Server**:
    - Add `POST /apps/upload` endpoint to receive binary files.
    - Update `StartApp` (or `DeployApp` logic) to calculate and store SHA256 hash of the binary.
    - Ensure `AppInfo` returned by `/apps/list` or `/apps/describe` includes the `BinaryHash`.
- **Logic**:
    - Client calculates local binary SHA256.
    - Client fetches remote app info.
    - If hashes match, skip upload and notify user.
    - If mismatch, upload binary and trigger app update/start.

## Impact
- **Specs**: `app-management`.
- **Code**: `cmd/glow/cmd/deploy.go` (new), `cmd/glow/cmd/verbs.go`, `internal/apiserver/server.go`, `internal/manager/manager.go`.
