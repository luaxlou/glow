# Glow Framework Examples

This directory contains framework-only examples for `starter/*` usage.

## Scope

Examples here demonstrate:

- configuration loading via `starter/glowconfig`
- HTTP/router initialization via `starter/glowhttp`
- storage/client initialization via starter packages

They intentionally do **not** include operations lifecycle workflows.

## simple-app

`examples/simple-app` is a minimal SDK sample that reads values from `config.json`.

### Run

```bash
cd examples/simple-app
cat > config.json <<'JSON'
{
  "log_level": "debug",
  "max_connections": 100,
  "mysql_dsn": "root:password@tcp(localhost:3306)/demo",
  "redis_addr": "localhost:6379",
  "port": 8080
}
JSON

go run ./cmd/simple-app
```

### Expected output

The app prints configured values and exits.

## Notes

- Operations runtime (`glow-server`, deployment orchestration) is in the `glow-ops` repository.
- Keep examples focused on SDK behavior only.
