# Change: Support Context Flag

## Why
Currently, switching contexts requires using `glow context use <name>`, which changes the global state. This makes scripting difficult as running commands against different environments in parallel or in sequence requires constantly changing the global configuration. Users need a way to specify the context for a single command execution without affecting the global state.

## What Changes
- **CLI**: Add a global `--context` flag to `glow`.
- **Logic**: Update command execution flow to prioritize the `--context` flag over the currently active context in the configuration file.

## Impact
- **Specs**: `authentication` (as it handles connection context).
- **Code**: `cmd/glow/cmd/root.go`.
