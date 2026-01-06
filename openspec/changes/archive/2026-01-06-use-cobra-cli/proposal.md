# Change: Use Cobra for CLI

## Why
The current CLI implementation in `cmd/glow/main.go` uses a manual `switch` statement for command dispatching. This is hard to maintain, lacks features like automatic help generation, flag parsing, and shell completion. We want to refactor the CLI to use the `spf13/cobra` library, which is the standard for Go CLIs.

## What Changes
- Introduce `cobra` and `viper` (optional, but good for config) dependencies.
- Structure `cmd/glow` into multiple files/packages:
    - `cmd/glow/root.go`: Root command.
    - `cmd/glow/app.go`: `app` subcommands.
    - `cmd/glow/config.go`: `config` subcommands.
    - `cmd/glow/ingress.go`: `ingress` subcommands.
    - `cmd/glow/setup.go`: `setup` command (renamed from old `config`).
- Replace manual flag parsing with Cobra flags.

## Impact
- Affected specs: No functional requirement changes, but implementation details change.
- Affected code: `cmd/glow/*.go`
