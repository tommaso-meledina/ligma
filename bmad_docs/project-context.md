---
project_name: 'ligma'
user_name: 'Tom'
date: '2026-01-20'
sections_completed: ['technology_stack', 'language_rules', 'framework_rules', 'testing_rules', 'quality_rules', 'workflow_rules', 'anti_patterns']
status: 'complete'
optimized_for_llm: true
---

# Project Context for AI Agents

_Critical rules and patterns for implementing code in ligma. Focus on unobvious details agents might otherwise miss._

---

## Technology Stack & Versions

- **Go** — stdlib only for MVP; `go mod` for deps.
- **Cobra** — `github.com/spf13/cobra`; CLI via `cobra-cli init` and `cobra-cli add <cmd>`.
- **cobra-cli** — `github.com/spf13/cobra-cli@latest`; scaffold `cmd/root.go`, `cmd/ls.go`, `cmd/get.go`, `cmd/write.go`.
- **SPDX:** `net/http` + `encoding/json`; 30s timeout; `licenses.json` + `details/{id}.json`.
- **[Growth]** **Viper** — `~/.ligma/config.json`; `--viper` on `cobra-cli init` when adding config.
- **[Growth]** **Cache** — file-based under `~/.ligma/_cache/`; mtime-based TTL; no `--no-cache` (use `cache_ttl: 0`).

---

## Critical Implementation Rules

### Language-Specific Rules (Go)

- **Naming:** Exported `PascalCase`; unexported `mixedCaps`; packages short, lowercase, no `_`.
- **Errors:** Always to **stderr**; never to stdout. Stdout only for `ls` and `get` (and `--json` when used).
- **Exit codes:** `0` = success, `1` = usage, `2` = not found, `3` = I/O or network. Exit from `main`/root only; no `os.Exit` inside `internal/` libraries.
- **SPDX IDs:** Use as-is (e.g. `MIT`, `Apache-2.0`); no case or format normalization for cache keys or storage.

### Framework-Specific Rules (Cobra)

- **Layout:** `cmd/root.go`, `cmd/ls.go`, `cmd/get.go`, `cmd/write.go`; one file per subcommand. `main.go` calls `cmd.Execute()`.
- **Flags:** `--popular`, `--json` on `ls` and `get` [Growth]; no `--no-cache`. Add flags in the command that uses them.
- **Help:** Use Cobra’s default `help`; do not replace or duplicate.

### Testing Rules

- **Placement:** `*_test.go` next to the package (e.g. `internal/spdx/client_test.go`). `go test ./...`.
- **Coverage:** Strive for 100%; project-wide threshold is **90%** to allow tactical omissions (e.g. tests that prove too hard to write). Prioritise SPDX client, cache [Growth], and config [Growth]; command wiring can be lighter.
- **Style:** Keep tests as simple and lean as possible; avoid external test libraries if the stdlib (`testing`, etc.) is sufficient.

### Code Quality & Style Rules

- **Layout:** `cmd/` for commands; `internal/spdx`, `internal/config`, `internal/cache`. No SPDX HTTP or JSON parse outside `internal/spdx`.
- **Formatting:** `gofmt`; follow standard Go style.
- **Cache [Growth]:** `~/.ligma/_cache/list.json`, `_cache/details/<id>.json`; use file **mtime** for TTL (not birth/creation time). `cache_ttl` null = default (e.g. 24h); `0` = always fetch. On cache write failure, still return fetched data.

### Development Workflow Rules

- **Build:** `go build -o ligma`; `go install` for a local binary.
- **First implementation:** Cobra scaffold (`go mod init`, `cobra-cli init`, `cobra-cli add ls|get|write`) per architecture.

### Critical Don’t-Miss Rules

- **No `--no-cache` flag.** Bypass cache only via `cache_ttl: 0` in config.
- **Exit code semantics:** 0/1/2/3 as above; do not introduce extra codes without updating this and the architecture.
- **Stderr for all errors and diagnostics.** Stdout only for `ls` and `get` output (and `--json` when set).
- **Cache [Growth]:** Create `_cache` and `_cache/details` on first use; do not fail the command if cache write fails.
- **SPDX:** Only `internal/spdx` performs HTTP and JSON for SPDX; [Growth] cache wraps those calls, config supplies URLs.
- **Config [Growth]:** `~/.ligma/` and `config.json` create-if-absent at startup (FR18); schema as PRD.

---

## Usage Guidelines

**For AI Agents:**

- Read this file before implementing any code.
- Follow all rules as documented; when in doubt, prefer the more restrictive option.
- Update this file if new patterns emerge during implementation.

**For Humans:**

- Keep this file lean and focused on agent needs.
- Update when the technology stack or architecture changes.
- Review periodically; remove rules that become obvious over time.

Last Updated: 2026-01-20
