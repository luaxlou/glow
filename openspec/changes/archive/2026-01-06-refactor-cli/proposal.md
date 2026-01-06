# Change: Refactor CLI

## Why
The current CLI mixes imperative commands (`start`, `stop`) with declarative `apply -f` patterns. To simplify usage and maintenance, we are removing the `apply -f <file.yaml>` mode and standardizing on a purely imperative CLI structure.

## What Changes
- Remove `apply` command.
- Standardize `glow <resource> <action>` pattern (e.g., `glow app start`, `glow app list`).
- Deprecate `glow start`, `glow stop` aliases in favor of structured subcommands.

## Impact
- Affected specs: `app-management`
- Affected code: `cmd/glow`, `pkg/api`
