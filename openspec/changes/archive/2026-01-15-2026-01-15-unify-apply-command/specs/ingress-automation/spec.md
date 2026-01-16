# Ingress Automation Deltas

## MODIFIED Requirements

### Requirement: Ingress 管理 (Ingress Management)
The system MUST support declarative ingress management via `glow apply`.

#### Scenario: Create Ingress (Manifest)
- **WHEN** user executes `glow apply -f ingress.yaml`
- **THEN** CLI creates the ingress rule
- **NOTE** `glow create ingress` is removed in favor of this workflow.
