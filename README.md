# Glow Framework

`glow` is the application framework repository.

## Scope

This repo only contains framework-facing capabilities:

- `starter/glowconfig`
- `starter/glowhttp`
- `starter/glowmysql`
- `starter/glowredis`
- `starter/glowsqlite`
- `starter/glowwebsocket`
- SDK examples and SDK documentation

## Non-Goals

This repo does not contain operations lifecycle orchestration:

- no `glow-server`
- no `glow-cli`
- no process supervision / deployment / rollback orchestration
- no server-side control-plane APIs

Operations runtime has been split into a separate repository: `glow-ops`.

## Development

```bash
go test ./...
```

See [docs/sdk_manual.md](docs/sdk_manual.md) for starter usage.
