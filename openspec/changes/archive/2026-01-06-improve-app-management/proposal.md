# Change: Improve App Management

## Why
Current app management implementation has several limitations:
1. Start operation is not idempotent (errors if already running).
2. Forced port allocation wastes resources for non-web apps.
3. Lack of explicit restart and delete capabilities.
4. Ambiguous status handling (STOPPED vs EXITED).

## What Changes
- Update `StartApp` to be idempotent.
- Make port allocation optional (read from config).
- Add `DeleteApp` capability.
- Refine status management to distinguish manual stop from abnormal exit.

## Impact
- Affected specs: `app-management`
- Affected code: `internal/manager`, `internal/apiserver`, `pkg/api`
