# Change: Add Build Binary Standard

## Why
Currently, there is no strict enforcement of where compiled binaries should be placed. Standardizing this to `./bin` ensures a clean project root and predictable tooling behavior.

## What Changes
- Introduce `build-release` capability.
- Require all build artifacts to be output to `./bin`.

## Impact
- Affected specs: `build-release` (new)
- Affected code: Build scripts, documentation, and Agent instructions.
