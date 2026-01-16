# Manifest Application Deltas

## ADDED Requirements

### Requirement: Config Resource Support
The system MUST support the `Config` kind in manifests to manage application configuration.

#### Scenario: Apply Config
- **WHEN** user executes `glow apply -f config.yaml` with `kind: Config`
- **THEN** CLI updates the application configuration via the API

### Requirement: Ingress Resource Support
The system MUST support the `Ingress` kind in manifests to manage ingress rules.

#### Scenario: Apply Ingress
- **WHEN** user executes `glow apply -f ingress.yaml` with `kind: Ingress`
- **THEN** CLI creates or updates the ingress rule via the API
