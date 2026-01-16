# Manifest Application Deltas

## MODIFIED Requirements

### Requirement: 声明式资源 (Declarative Resources)
The system MUST support K8s-style YAML resource definitions.

#### Scenario: App Resource Definition
- `App`: Defines an application, its configuration, and execution parameters.
- **`spec.binaryPath` (optional)**: A new field specifying the path to the application's binary on the *client's* local filesystem. If present, `glow apply` will hash and upload this binary if it differs from the server's version.
