# Tasks: App Validation and Update

### Server-Side
- [x] **API**: Implement `GET /apps/{name}/state` endpoint. <!-- id: 0 -->
- [x] **API**: Implement `POST /apps/{name}/binary` endpoint for file uploads. <!-- id: 1 -->
- [x] **State Management**: Extend the application's data model to store `configHash` and `binaryHash`. <!-- id: 2 -->
- [x] **File Management**: Implement logic to save and manage uploaded binaries on the server's filesystem. <!-- id: 3 -->

### Client-Side
- [x] **API Types**: Add `binaryPath` field to the `AppSpecOld` struct. <!-- id: 4 -->
- [x] **Hashing**: Implement SHA256 hashing for JSON config and binary files. <!-- id: 5 -->
- [x] **`apply` Command**: Refactor `apply` logic for `App` resources to include the pre-flight check and conditional update flow. <!-- id: 6 -->
- [x] **`apply` Command**: Implement binary upload logic using a multipart/form-data request. <!-- id: 7 -->

### Documentation
- [x] Update `cli_manual.md` to document the new `binaryPath` field and the validation behavior. <!-- id: 8 -->
- [x] Update `app-management` and `manifest-application` specs with the new requirements. <!-- id: 9 -->
