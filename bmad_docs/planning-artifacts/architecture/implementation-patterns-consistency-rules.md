# Implementation Patterns & Consistency Rules

## Pattern Categories Defined

**Critical conflict points:** Naming (Go, files, SPDX IDs); layout (`cmd/`, `internal/`); errors and exit codes; test placement.

## Naming Patterns

- **Go:** Exported: `PascalCase`; unexported: `mixedCaps`. Files: `lowercase.go` or `snake_case.go` to match main symbol. Packages: short, lowercase, no `_`.
- **SPDX IDs:** Use as-is (e.g. `MIT`, `Apache-2.0`); no normalization for storage or cache keys.
- **Commands/flags:** `ls`, `get`, `write`, `help`; flags `--popular`, `--json` (Growth).

## Structure Patterns

- **Commands:** `cmd/root.go`, `cmd/ls.go`, `cmd/get.go`, `cmd/write.go`; one file per subcommand.
- **Libraries:** `internal/spdx` (fetch, parse), `internal/config` [Growth], `internal/cache` [Growth]. `internal/` not importable by external modules.
- **Tests:** `*_test.go` alongside package (e.g. `internal/spdx/client_test.go`); `go test ./...`.

## Format Patterns

- **Exit codes:** 0 success, 1 usage, 2 not found, 3 I/O or network; always to `os.Exit` (or equivalent) from `main`/root.
- **Stderr:** All errors and diagnostics to stderr; stdout only for `ls` and `get` (and `--json` when used).
- **JSON (consumed):** SPDX as-is. **JSON (produced, Growth):** `--json` for `ls`/`get`; field naming TBD in implementation (e.g. Go struct tags).

## Process Patterns

- **Error handling:** Fail fast; message to stderr; exit non-zero. No interactive recovery.
- **Cache (Growth):** Wrapper around SPDX HTTP; on write failure, return fetched data, do not fail the command.

## Enforcement

**All agents MUST:** use `cmd/` and `internal/` as above; follow exit code semantics; write errors to stderr; use mtime and `_cache` layout for cache [Growth]; no `--no-cache`, use `cache_ttl: 0`.

---
