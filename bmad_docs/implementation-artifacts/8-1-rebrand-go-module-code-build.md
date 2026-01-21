# Story 8.1: Rebrand Go module, code, and build

Status: review

## Story

As a **developer**,
I want **the Go module, all source code, and build-related files to use the name `ligma`**,
So that **the binary, config paths, and programmatic identifiers reflect the new brand**.

## Acceptance Criteria

**Given** the current codebase under `github.com/tom/licensegen`
**When** the rebrand is applied to Go and build artifacts
**Then** `go.mod` declares `module github.com/tom/ligma`
**And** every import path `github.com/tom/licensegen/...` is updated to `github.com/tom/ligma/...` in: `main.go`, `cmd/*.go`, `cmd/*_test.go`, `internal/config/config.go`, `internal/config/config_test.go`, `internal/cache/cache.go`, `internal/cache/cache_test.go`
**And** `cmd/root.go` uses `Use: "ligma"` (replacing `"licensegen"`)
**And** `internal/config`: the function `LicensegenDir` is renamed to `LigmaDir`; all references in `internal/config/config_test.go`, `cmd/get.go`, and `cmd/ls.go` are updated; the tests `TestLicensegenDir_Override` and `TestLicensegenDir_UserHomeDir` are renamed to `TestLigmaDir_Override` and `TestLigmaDir_UserHomeDir`
**And** all user-facing or comment references to `~/.licensegen` are changed to `~/.ligma` in: `internal/config/config.go`, `internal/config/config_test.go`, `internal/cache/cache.go`, `cmd/write.go` (error message: `~/.ligma/config.json`)
**And** `.gitignore` lists the binary as `ligma` instead of `licensegen`
**And** `go build -o ligma` succeeds and tests pass

## Tasks / Subtasks

- [x] **Task 1: go.mod and import paths** (AC: module, imports)
  - [x] `go.mod`: `module github.com/tom/ligma`
  - [x] Update `github.com/tom/licensegen/...` → `github.com/tom/ligma/...` in: `main.go`, `cmd/root.go`, `cmd/root_test.go`, `cmd/ls.go`, `cmd/ls_test.go`, `cmd/get.go`, `cmd/get_test.go`, `cmd/write.go`, `cmd/write_test.go`, `internal/cache/cache.go`, `internal/cache/cache_test.go` (internal/config has no self-import)
- [x] **Task 2: root Use and config renames** (AC: Use, LicensegenDir→LigmaDir)
  - [x] `cmd/root.go`: `Use: "ligma"`
  - [x] `internal/config/config.go`: `LicensegenDir` → `LigmaDir`
  - [x] `internal/config/config_test.go`: `LicensegenDir` → `LigmaDir`; `TestLicensegenDir_Override` → `TestLigmaDir_Override`; `TestLicensegenDir_UserHomeDir` → `TestLigmaDir_UserHomeDir`
  - [x] `cmd/get.go`, `cmd/ls.go`: `config.LicensegenDir` → `config.LigmaDir`
- [x] **Task 3: ~/.licensegen → ~/.ligma** (AC: user-facing paths)
  - [x] `internal/config/config.go`: comments and `filepath.Join(home, ".licensegen")` → `.ligma`
  - [x] `internal/config/config_test.go`: `filepath.Base(got) != ".licensegen"` → `.ligma`
  - [x] `internal/cache/cache.go`: comment `~/.licensegen/_cache` → `~/.ligma/_cache`
  - [x] `cmd/write.go`: error `~/.licensegen/config.json` → `~/.ligma/config.json`
- [x] **Task 4: .gitignore and build** (AC: .gitignore, build)
  - [x] `.gitignore`: `licensegen` → `ligma`
  - [x] Run `go build -o ligma` and `go test ./...`

## Dev Notes

- **Source:** epic-8-rebrand-to-ligma.md Story 8.1
- **Files to edit (15):** `go.mod`, `main.go`, `cmd/root.go`, `cmd/root_test.go`, `cmd/ls.go`, `cmd/ls_test.go`, `cmd/get.go`, `cmd/get_test.go`, `cmd/write.go`, `cmd/write_test.go`, `internal/config/config.go`, `internal/config/config_test.go`, `internal/cache/cache.go`, `internal/cache/cache_test.go`, `.gitignore`
- **internal/spdx:** No `github.com/tom/licensegen` import; no change. internal/config has no self-import; only renames and `~/.licensegen`→`~/.ligma`.

## Dev Agent Record

### Implementation Plan

- Rebrand: `go.mod` → `github.com/tom/ligma`; all import paths in main, cmd, internal/cache (internal/config has no self-import; internal/spdx unchanged).
- `cmd/root.go` Use: `"ligma"`. `LicensegenDir` → `LigmaDir` in config; `config.LicensegenDir` → `config.LigmaDir` in cmd/get, cmd/ls. Test renames in config_test.
- All `~/.licensegen` and `.licensegen` → `~/.ligma` / `.ligma` in config, config_test, cache (comments), write (error msg).
- `.gitignore`: binary `licensegen` → `ligma`. Verified `go build -o ligma` and `go test ./...` pass.

### File List

- `go.mod`
- `main.go`
- `cmd/root.go`
- `cmd/root_test.go`
- `cmd/ls.go`
- `cmd/ls_test.go`
- `cmd/get.go`
- `cmd/get_test.go`
- `cmd/write.go`
- `cmd/write_test.go`
- `internal/config/config.go`
- `internal/config/config_test.go`
- `internal/cache/cache.go`
- `internal/cache/cache_test.go`
- `.gitignore`

## Change Log

- 2026-01-20: Story 8.1 implemented — rebrand to ligma: go.mod, imports, root Use, LicensegenDir→LigmaDir, ~/.licensegen→~/.ligma, .gitignore; `go build -o ligma` and tests pass.
