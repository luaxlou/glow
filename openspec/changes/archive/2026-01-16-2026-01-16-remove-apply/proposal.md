# Proposal: Remove Declarative Apply Functionality

## Summary
Remove the `glow apply` command and all associated server-side logic (`/apply/app` endpoint, manifest parsing). The system will no longer support declarative resource management via YAML/JSON manifests.

## Motivation
The user has explicitly requested to remove the entire `apply` functionality and stop supporting it. This simplifies the CLI and reduces maintenance burden for features that may be deemed unnecessary or redundant.

## Proposed Changes
1.  **CLI**: Remove `glow apply` command (`cmd/glow/cmd/apply.go`).
2.  **Server**: Remove `POST /apply/app` endpoint and handler (`internal/apiserver/server.go`).
3.  **Manager**: Remove `ApplyApp` logic (`internal/manager/manager.go`).
4.  **Package**: Remove `pkg/manifest` package as it is no longer used.
5.  **Documentation**: Update `cli_manual.md` and `server_manual.md` to remove references to `apply`.
6.  **Specs**: Update OpenSpec documents to reflect the removal of declarative requirements.

## Impact
- Users will no longer be able to use `glow apply -f ...` to deploy or update resources.
- Application management must be done via imperative commands (`start`, `stop`, `config edit`, etc.) or other means.
