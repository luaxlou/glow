# Config Management Deltas

## MODIFIED Requirements

### Requirement: 应用配置管理 (App Config Management)
The system MUST support declarative configuration management via `glow apply`.

#### Scenario: Update Application Config (Manifest)
- **WHEN** user executes `glow apply -f config.yaml`
- **THEN** CLI reads the `Config` resource and updates the application configuration
- **NOTE** `glow config apply` is replaced by this workflow.
