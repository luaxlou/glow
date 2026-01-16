# App Management Deltas

## ADDED Requirements

### Requirement: Application State Validation
The system MUST provide an endpoint to query the configuration and binary state of a deployed application.

#### Scenario: Get App State
- **GIVEN** an application "my-app" is deployed
- **WHEN** the client sends `GET /apps/my-app/state`
- **THEN** the server MUST return a JSON object containing the current `configHash` and `binaryHash`.

### Requirement: Application Binary Upload
The system MUST provide an endpoint for clients to upload application binaries.

#### Scenario: Upload a New Binary
- **WHEN** the client sends a `POST` request to `/apps/my-app/binary` with a multipart/form-data payload containing the binary file
- **THEN** the server MUST securely save the binary to a managed location.
- **AND** the server MUST update the application's stored `binaryHash`.

### Requirement: Idempotent Application Apply
The `glow apply` command MUST be idempotent for applications. If the configuration and binary are already up-to-date, no update action should be performed.

#### Scenario: Apply an Unchanged Application
- **GIVEN** the client's local manifest and binary file have not changed since the last apply
- **WHEN** the user executes `glow apply -f app.yaml`
- **THEN** the client MUST validate hashes with the server.
- **AND** the client MUST determine that no changes are needed.
- **AND** the client MUST NOT send any update requests to the server.
- **AND** the CLI SHOULD print a message like "Application 'my-app' is already up-to-date."
