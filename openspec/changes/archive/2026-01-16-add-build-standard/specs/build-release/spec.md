## ADDED Requirements
### Requirement: Binary Output Location
All compiled binaries MUST be output to the `./bin` directory in the project root.

#### Scenario: Build a component
- **WHEN** a developer or agent compiles a component (e.g., `glow-server`)
- **THEN** the resulting binary file is placed in `./bin/`
- **AND** the project root remains free of binary artifacts
