# Story 1.1: Go module, project structure, and dummy main

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **developer**,
I want **to run `go mod init` and lay out the basic project structure with a dummy `main.go` that prints "OK"**,
so that **the project builds and running the binary produces a known output**.

## Acceptance Criteria

1. **Given** an empty project directory  
   **When** I run `go mod init <module-path>` and create a minimal layout (`go.mod`, `main.go`, `.gitignore` as needed) with a `main.go` whose `main` function prints `OK` to stdout  
   **Then** `go build -o licensegen` succeeds and `./licensegen` prints `OK` to stdout and exits with 0  
   **And** the layout is consistent with the eventual Architecture (e.g. `main.go` at repo root; `cmd/` and `internal/` may be added in later stories)

## Tasks / Subtasks

- [x] **Task 1: Initialize Go module** (AC: #1)
  - [x] Run `go mod init <module-path>` (e.g. `github.com/yourusername/licensegen`). Use the actual repo/module path for this project.
  - [x] Ensure `go.mod` is at repo root.
- [x] **Task 2: Add main.go** (AC: #1)
  - [x] Create `main.go` at repo root.
  - [x] Implement `main` to print `OK` to **stdout** (e.g. `fmt.Println("OK")` or equivalent) and exit with 0.
  - [x] Do **not** add `cmd/` or `internal/` in this story; those come in Story 1.2.
- [x] **Task 3: Add .gitignore** (AC: #1)
  - [x] Create `.gitignore` with entries typical for a Go project (e.g. `licensegen` binary, `vendor/` if used, IDE/OS junk). Align with project conventions.
- [x] **Task 4: Verify build and run** (AC: #1)
  - [x] Run `go build -o licensegen` and confirm it succeeds.
  - [x] Run `./licensegen` and confirm it prints `OK` to stdout and exits with 0.

## Dev Notes

- **Scope:** This story is strictly **foundation only**. Do **not** add Cobra, `cmd/`, or `internal/` here. Story 1.2 will add the Cobra scaffold and `ls`/`get`/`write` stubs.
- **Stdout vs stderr:** For this story, the only programmatic output is `OK` to **stdout**. Error handling and stderr come in later stories; keep `main` minimal.
- **Exit code:** Success must exit with **0** only from `main`; no `os.Exit` inside `internal/` (there is no `internal/` yet).
- **Architecture alignment:** `main.go` at repo root matches the final layout in `project-structure-boundaries.md`. `cmd/` and `internal/` are intentionally deferred to 1.2.

### Project Structure Notes

- **This story only:** repo root should contain at least: `go.mod`, `main.go`, `.gitignore`. Optional: `README.md` if it already exists; do not introduce new docs.
- **Future layout (for context only, do NOT implement now):** `cmd/root.go`, `cmd/ls.go`, `cmd/get.go`, `cmd/write.go`, `internal/spdx`, `internal/config` [Growth], `internal/cache` [Growth]. See [Source: bmad_docs/planning-artifacts/architecture/project-structure-boundaries.md].

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-1-cli-foundation-help.md#story-11-go-module-project-structure-and-dummy-main]
- [Source: bmad_docs/planning-artifacts/architecture/project-structure-boundaries.md] — Project Directory Structure, Boundaries
- [Source: bmad_docs/planning-artifacts/architecture/starter-template-evaluation.md] — Cobra init comes in 1.2; here only `go mod init` and `main.go`
- [Source: bmad_docs/project-context.md] — Technology Stack, Language Rules, Dev Workflow (build: `go build -o licensegen`)

---

## Developer Context

### Technical Requirements

- **Go:** Use only the standard library for this story. No Cobra, Viper, or other deps. `go mod init` creates `go.mod`; after that, `main.go` may use only `fmt` (or equivalent) for printing.
- **Build:** `go build -o licensegen` must succeed. `go install` is optional for this story.
- **Binary name:** Output binary must be `licensegen` (or as built) so that `./licensegen` runs the dummy main.
- **Exit code:** `main` must cause the process to exit with **0** on success. Do not call `os.Exit` with any other code in this story.

### Architecture Compliance

- **Layout:** `main.go` at **repo root**. No `cmd/` or `internal/` in this story. The architecture’s `cmd/` and `internal/` are introduced in Story 1.2.
- **Decision Impact (from core-architectural-decisions):** “Cobra scaffold: `go mod init`, `cobra-cli init`, `cobra-cli add ls|get|write`.” For 1.1, only `go mod init` and a minimal `main.go`; Cobra and subcommands are 1.2.
- **Boundaries:** No SPDX, config, or cache in this story. Those live in `internal/` in later stories.

### Library / Framework Requirements

- **This story:** None. Stdlib only. No `cobra-cli`, no Viper, no `net/http`, no `encoding/json`.
- **Module path:** Choose a valid `go mod init` path (e.g. `github.com/yourusername/licensegen`). The project name in `go.mod` should align with the repo.

### File Structure Requirements

| Path           | Purpose                                      | This story |
|----------------|----------------------------------------------|------------|
| `go.mod`       | Go module definition at repo root            | **Create** |
| `main.go`      | Entrypoint; prints `OK` to stdout, exits 0   | **Create** |
| `.gitignore`   | Ignore binary, `vendor/`, common junk        | **Create** |
| `cmd/`         | Cobra root + subcommands                     | **Do not add** (1.2) |
| `internal/`    | SPDX, config [Growth], cache [Growth]        | **Do not add** (later) |

- Do **not** create `cmd/`, `internal/`, or any `*_test.go` in this story unless you are only adding a trivial `main` test; the architecture expects `*_test.go` next to packages, and there are no packages yet.

### Testing Requirements

- **This story:** No mandatory `*_test.go`. The AC is verified by: `go build -o licensegen` and `./licensegen` → `OK` and exit 0. If you add a trivial test (e.g. for a small helper), place it next to `main.go`; no `internal/` tests yet.
- **Later stories:** `*_test.go` alongside packages; `go test ./...`; coverage target 90%. Not in scope for 1.1.

### Previous Story Intelligence

- **N/A.** This is the first story. There is no prior implementation or learnings to reuse.

### Latest Tech Information

- **Go / go mod:** Use a current stable Go toolchain. `go mod init` and `go build` are stable; no special version or API requirements for this minimal story.
- **No third‑party deps:** No web research or version checks needed for 1.1.

### Project Context Reference

- **bmad_docs/project-context.md** is the source of truth for:
  - Technology: Go, stdlib for MVP; `go mod` for deps.
  - Build: `go build -o licensegen`; `go install` for a local binary.
  - First implementation: Cobra scaffold in 1.2; here only `go mod init` and minimal `main.go`.
  - Errors to stderr / exit codes: defined from Story 1.3 onward; this story only requires exit 0 and stdout for `OK`.

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- **1-1 (2026-01-20):** Implemented AC#1: `go mod init github.com/tom/licensegen`, `main.go` (prints `OK` to stdout, exit 0), `.gitignore` (licensegen, vendor/, IDE/OS). Verified: `go build -o licensegen` and `./licensegen` → `OK` and exit 0. No `*_test.go` per story (AC verified by build/run).

### File List

- `go.mod` (created)
- `main.go` (created)
- `.gitignore` (created)

## Change Log

- 2026-01-20: Story 1.1 implemented — go.mod, main.go, .gitignore; `go build -o licensegen` and `./licensegen` verified (OK, exit 0).
