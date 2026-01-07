# Change: Rename Deployment to App

## Why
The user requested to change the naming convention from "Deployment" back to "App". This simplifies the terminology, aligning with the project name "Glow App" and the underlying `AppInfo` structure, while maintaining the K8s-style CLI structure (`glow get app`).

## What Changes
- **Terminology**: Rename "Deployment" -> "App" in specs and CLI.
- **CLI Commands**:
    - `glow get deploy` -> `glow get app`
    - `glow describe deploy` -> `glow describe app`
    - `glow start/stop/... deploy` -> `glow start/stop/... app`
- **Manifest**: `kind: Deployment` -> `kind: App`.

## Impact
- **Specs**: `app-management`, `manifest-application`.
- **Code**: `cmd/glow/cmd/deploy.go` (rename to `app.go`?), `pkg/api/types.go` (rename `Deployment` struct to `App` or alias), `manifest/parser.go`.
