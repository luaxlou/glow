# Proposal: Add Validation and Conditional Update to `apply app`

## Summary
This proposal enhances the `glow apply` command for applications by introducing a client-side validation step. Before applying changes, the client will compare the local configuration and binary file hash against the state on the server. Updates will only be performed if inconsistencies are detected, making the apply process more efficient and idempotent.

## Motivation
Currently, `glow apply` unconditionally pushes the manifest configuration to the server. This can lead to unnecessary application restarts even when no actual changes have been made. Furthermore, there is no mechanism to ensure that the binary file being executed by the server matches the user's local version, which can lead to configuration drift and hard-to-debug issues.

## Proposed Changes
1.  **Client-Side Validation**: Before applying, the `glow` client will perform a pre-flight check with the server.
2.  **Configuration Hashing**: The client will compute a hash of the `spec` section of the application manifest.
3.  **Binary Hashing & Upload**:
    *   A new `binaryPath` field will be added to the `App` manifest `spec`.
    *   If `binaryPath` is present, the client will compute a SHA256 hash of the local binary file.
    *   If the binary hash on the server is different, the client will upload the new binary.
4.  **Conditional Updates**: The client will only send update requests (for config or binary) if the corresponding hashes do not match the server's state.

This change introduces a more robust, "Git-like" workflow for managing applications, where the client intelligently syncs local changes with the remote environment.
