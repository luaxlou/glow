# Change: Add Context Management

## Why
Currently, the Glow client (`glow`) only supports a single server configuration (URL and API Key). Users managing multiple Glow environments (e.g., local, dev, prod) need a way to switch between them easily, similar to `kubectl config use-context`.

## What Changes
- **Configuration Schema**: Update `~/.glow.json` to store a list of contexts and a `current-context` pointer.
- **CLI Commands**: Introduce `glow context` command group.
    - `glow context list`: Show all contexts.
    - `glow context use <name>`: Switch current context.
    - `glow context add <name> --url <url> --key <key>`: Add a new context.
    - `glow context delete <name>`: Remove a context.
- **Auth Interaction**: `glow auth` commands will continue to work, operating on the *current* context.

## Impact
- **Specs**: `authentication`.
- **Code**: `cmd/glow/cmd/root.go` (Config loading logic), `cmd/glow/cmd/context.go` (New commands).
