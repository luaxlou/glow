# Proposal: Unify Apply Command & Remove Legacy Create Commands

## Summary
Refactor the CLI to use `glow apply -f` as the single source of truth for applying configurations, including `Config` and `Ingress` resources. This aligns with Kubernetes patterns and simplifies the CLI interface by removing ad-hoc creation commands like `glow create ingress` and `glow config apply`.

## Motivation
Currently, resource management is fragmented:
- `App` and `Host` use `glow apply`.
- `Config` uses `glow config apply`.
- `Ingress` uses `glow create ingress`.

This inconsistency increases cognitive load and makes automation harder (e.g., applying a folder of manifests). Unifying these under `glow apply` with standard YAML manifests provides a consistent, declarative workflow.

## Proposed Changes
1.  **Manifest Support**: Extend `glow apply` to support `Config` and `Ingress` Kinds.
2.  **Config Management**: Deprecate/Remove `glow config apply` in favor of `glow apply -f config.yaml`.
3.  **Ingress Automation**: Remove `glow create ingress` in favor of `glow apply -f ingress.yaml`.
4.  **API Types**: Ensure `Config` is a first-class citizen in the manifest parser and API types.

## Design Details
### Config Resource
```yaml
apiVersion: v1
kind: Config
metadata:
  name: my-app
data:
  DB_HOST: "localhost"
  DB_PORT: 3306
```

### Ingress Resource
```yaml
apiVersion: v1
kind: Ingress
metadata:
  name: my-ingress
spec:
  domain: example.com
  service: my-app
  port: 8080
```

## Impact
- **CLI**: `glow create ingress` removed. `glow config apply` removed/deprecated.
- **Docs**: CLI manual and specs updated to reflect declarative workflows.
