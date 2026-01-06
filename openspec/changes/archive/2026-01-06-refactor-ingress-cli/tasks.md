## 1. Implementation
- [x] 1.1 Refactor `internal/manager/nginx.go` to support independent listing and status checks.
- [x] 1.2 Implement API handler `handleIngressUpdate` (POST `/ingress/update`) in `internal/apiserver`.
- [x] 1.3 Implement API handler `handleIngressDelete` (POST `/ingress/delete`) in `internal/apiserver`.
- [x] 1.4 Implement API handler `handleIngressList` (GET `/ingress/list`) in `internal/apiserver`.
- [x] 1.5 Implement `glow ingress apply` client command.
- [x] 1.6 Implement `glow ingress delete` client command.
- [x] 1.7 Implement `glow ingress list` client command.
- [x] 1.8 Update `StartApp` to use the shared Ingress logic (maintain backward compatibility).
- [x] 1.9 Add integration tests for Ingress API.
