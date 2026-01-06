# Change: Refactor Config CLI

## Why
Currently, the `glow config` command is limited to setting client-side connection details (URL and Key). To better manage application configurations (JSON configs stored in SQLite), we need a full CRUD CLI interface.

## What Changes
- Expand `glow config` to support subcommands: `set`, `get`.
- Note: Current `glow config --url ...` usage for client setup should be moved to `glow init` or `glow client config`. For now, we will repurpose `glow config` for *Application Configuration* as per the `config-management` spec.
- Wait, `glow config` currently sets *client* config. The user request is "config 增删改查" (CRUD). This likely refers to the *App Configs* managed by `configmanager`.
- We should distinguish:
    - `glow client-config` (or keep as `glow config` for client?) -> Maybe rename current `config` to `setup` or `login`? Or just `glow config set-client`.
    - `glow config` -> App Config Management.

**Decision**:
- `glow config set <app> <key> <value>`
- `glow config get <app> [key]`
- `glow config list <app>` (shows all keys)
- `glow config delete <app> <key>` (if supported by backend)

But first, let's check what `handleConfig` currently does. It saves client config.
We will rename the current `handleConfig` to `handleClientConfig` and expose it as `glow client setup` (or similar), and free up `glow config` for App Configs.
Or, to be less disruptive, maybe `glow app config ...`?
The user said "glow config 增删改查". So `glow config` should probably be the command.

## Impact
- Affected specs: `config-management`
- Affected code: `cmd/glow`
