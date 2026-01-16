# Design: App Validation and Update Flow

This document outlines the technical design for implementing client-side validation and conditional updates for `glow apply app`.

## 1. API Changes

### `GET /apps/{name}/state`
A new endpoint to retrieve the current state hashes of a deployed application.

**Request:**
`GET /apps/my-cool-app/state`

**Response (Success):**
```json
{
  "success": true,
  "data": {
    "configHash": "sha256:abc...",
    "binaryHash": "sha256:def..."
  }
}
```
If the app does not exist, the hashes will be empty or null, signaling a new deployment.

### `POST /apps/{name}/binary`
A new endpoint to upload a new binary for an application.

**Request:**
`POST /apps/my-cool-app/binary`
- **Content-Type**: `multipart/form-data`
- **Body**: A form field named `binary` containing the file data.

**Server-Side Logic:**
- The server will save the uploaded file to a managed directory, e.g., `/var/glow/binaries/{appName}/{sha256_of_file}`.
- After saving, it will update the application's stored `binaryHash` and potentially update the `command` path in the app's configuration to point to this new binary.

## 2. Client-Side `apply` Flow

The `cmd/glow/cmd/apply.go` logic for `api.App` will be modified as follows:

1.  **Parse Manifest**: Read the `app.yaml` file.
2.  **Calculate Hashes**:
    *   **Config Hash**: Serialize the `spec` field of the manifest to canonical JSON and compute its SHA256 hash.
    *   **Binary Hash**:
        *   Look for a new `spec.binaryPath` field in the manifest.
        *   If it exists and points to a valid local file, compute the file's SHA256 hash.
        *   If not, the binary hash is considered empty.
3.  **Pre-flight Check**:
    *   Execute `GET /apps/{name}/state` to fetch the server's current hashes.
4.  **Compare and Update**:
    *   `if localConfigHash != serverConfigHash`:
        *   `POST /apply/app` with the new manifest `spec`.
        *   Log "Configuration updated."
    *   `if localBinaryHash != serverBinaryHash`:
        *   `POST /apps/{name}/binary` with the local binary file from `binaryPath`.
        *   Log "Binary updated."
    *   `if both hashes match`:
        *   Log "Application is already up-to-date."
        *   Do nothing.

## 3. Manifest Changes

A new, optional field `binaryPath` will be added to the `AppSpecOld` struct in `pkg/api/types.go` (and the corresponding spec).

**Example `app.yaml`:**
```yaml
apiVersion: v1
kind: App
metadata:
  name: my-app
spec:
  # The path to the binary on the CLIENT machine.
  binaryPath: ./bin/my-app-executable
  # The command that will be run on the SERVER.
  # The server will replace this with the path to the uploaded binary.
  command: ./my-app-executable
  args:
    - --config
    - /etc/app/config.json
```
This design ensures that changes are only pushed when necessary, reducing network traffic and avoiding unnecessary restarts, while also guaranteeing binary consistency.
