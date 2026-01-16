# Tasks: Unify Apply Command

- [x] Define `Config` and `Ingress` types in `pkg/api/types.go` <!-- id: 0 -->
- [x] Update `pkg/manifest/parser.go` to support `Config` and `Ingress` <!-- id: 1 -->
- [x] Refactor `cmd/glow/cmd/apply.go` to handle `Config` and `Ingress` resources <!-- id: 2 -->
- [x] Remove `glow create ingress` from `cmd/glow/cmd/ingress.go` <!-- id: 3 -->
- [x] Remove/Deprecate `glow config apply` from `cmd/glow/cmd/config.go` <!-- id: 4 -->
- [x] Verify `glow apply -f` works for all resource types <!-- id: 5 -->
