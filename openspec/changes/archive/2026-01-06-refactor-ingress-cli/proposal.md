# Change: Refactor Ingress Management to Client-Server Model

## Why
Currently, Ingress (Nginx) configuration is implicitly tied to `AppStart`. We want to decouple this and manage Ingress explicitly via the CLI. To align with Kubernetes and modern operations patterns, this management should be performed via the `glow` client (sending requests to the server) rather than running administrative commands directly on the server binary (`glow-server`).

## What Changes
- **New API Endpoints**: Add Ingress management endpoints to `glow-server` API (e.g., POST `/ingress/update`, POST `/ingress/delete`, GET `/ingress/list`).
- **Client Update**: Update `glow` CLI to support `glow ingress` subcommand group.
    - `glow ingress apply`: Create or update an ingress rule.
    - `glow ingress delete`: Remove an ingress rule.
    - `glow ingress list`: List active ingress rules.
- **Spec Update**: Redefine `ingress-automation` to be an API-driven capability exposed via the standard client.
- **Refactoring**: Isolate Nginx logic in `internal/manager/nginx.go` to be callable from the new API handlers.

## Impact
- **Affected Specs**: `ingress-automation`
- **Affected Code**:
    - `internal/apiserver/`: Add ingress handlers.
    - `cmd/glow/`: Add ingress subcommands.
    - `internal/manager/`: Refactor `nginx.go` for independent use.