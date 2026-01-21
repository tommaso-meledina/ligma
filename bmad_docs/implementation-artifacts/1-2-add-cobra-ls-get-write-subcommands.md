# Story 1.2: Add Cobra and ls, get, write subcommands

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **developer**,
I want **to use cobra-cli to add the root command and `ls`, `get`, `write` subcommands**,
so that **`licensegen ls`, `licensegen get`, and `licensegen write` are callable from the CLI regardless of their output**.

## Acceptance Criteria

1. **Given** a Go project with a working `main.go`  
   **When** I run `cobra-cli init` (integrating its output with existing `main.go` as needed) and `cobra-cli add ls`, `cobra-cli add get`, `cobra-cli add write`  
   **Then** `licensegen ls`, `licensegen get`, `licensegen write` execute without error and are invokable from the command line  
   **And** the subcommands exist in the structure defined by Architecture (e.g. `cmd/root.go`, `cmd/ls.go`, `cmd/get.go`, `cmd/write.go`); the exact output or behavior of each subcommand is not yet specified  
   **And** `licensegen --help` and `licensegen <command> --help` display Cobra-generated help (FR13)

## Tasks / Subtasks

- [x] **Task 1: Install cobra-cli and add Cobra dependency** (AC: #1)
  - [x] Run `go install github.com/spf13/cobra-cli@latest` (ensure it is on PATH). Cobra itself will be added to `go.mod` when running `cobra-cli init`/`add`.
  - [x] Ensure the project root has `go.mod` and `main.go` from Story 1.1.
- [x] **Task 2: Run cobra-cli init and integrate with main.go** (AC: #1)
  - [x] Run `cobra-cli init` from the project root. Do **not** use `--viper` (reserved for Growth / config).
  - [x] If `cobra-cli init` creates or overwrites `main.go`, use the generated `main.go` that calls `cmd.Execute()` (or the root command’s `Execute`). If it does not create `main.go`, replace the 1.1 dummy body in `main.go` so it calls `cmd.Execute()` and remove the "OK" print.
  - [x] Ensure `main.go` remains at repo root and is the single entrypoint; `go build -o licensegen` must still succeed.
- [x] **Task 3: Add ls, get, write subcommands** (AC: #1)
  - [x] Run `cobra-cli add ls`, `cobra-cli add get`, `cobra-cli add write` from the project root.
  - [x] Confirm `cmd/ls.go`, `cmd/get.go`, `cmd/write.go` exist and each registers its command under the root (e.g. `rootCmd.AddCommand(lsCmd)` in `init()`). Do not use `--parent`; these are direct children of root.
- [x] **Task 4: Verify subcommands and help** (AC: #1)
  - [x] Run `./licensegen ls`, `./licensegen get`, `./licensegen write` — each must run without error. Output/behavior of the commands is unspecified for this story (stubs are fine).
  - [x] Run `./licensegen --help` and `./licensegen ls --help`, `./licensegen get --help`, `./licensegen write --help` — each must show Cobra-generated help (FR13).

## Dev Notes

- **Pre-requisite:** Story 1.1 must be done: `go.mod`, `main.go` (dummy "OK"), `.gitignore` at repo root. This story replaces the dummy `main` behavior with the Cobra entrypoint and adds `cmd/`.
- **No Viper:** Do **not** pass `--viper` to `cobra-cli init`. Viper is for Growth when implementing `~/.licensegen/config.json` (project-context, starter-template-evaluation).
- **Subcommand behavior:** AC says "the exact output or behavior of each subcommand is not yet specified." `ls`, `get`, `write` can be no-ops (e.g. `Run`/`RunE` that does nothing and returns `nil`). They must only be invokable and show `--help`.
- **Exit codes:** Story 1.3 will refine exit codes (0/1/2/3). Here, "execute without error" and help succeeding is enough; Cobra’s default behavior on `--help` is typically exit 0.
- **Help:** Use Cobra’s default help; do not replace or duplicate (project-context, framework rules).
- **No `internal/` in this story:** `internal/spdx`, `internal/config`, `internal/cache` are added in later stories. Only `cmd/` and `main.go` change here.

### Project Structure Notes

- **After this story:** repo root: `main.go` (calls `cmd.Execute()`), `go.mod`, `go.sum`, `.gitignore`. `cmd/`: `root.go`, `ls.go`, `get.go`, `write.go`. No `internal/` yet.
- **One file per subcommand:** `cmd/root.go`, `cmd/ls.go`, `cmd/get.go`, `cmd/write.go` (implementation-patterns, project-structure-boundaries).

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-1-cli-foundation-help.md#story-12-add-cobra-and-ls-get-write-subcommands]
- [Source: bmad_docs/planning-artifacts/architecture/project-structure-boundaries.md] — Project Directory Structure, Boundaries
- [Source: bmad_docs/planning-artifacts/architecture/starter-template-evaluation.md] — Cobra + cobra-cli, init/add sequence, no `--viper` for MVP
- [Source: bmad_docs/project-context.md] — Cobra layout, Help, Framework-Specific Rules

---

## Developer Context

### Technical Requirements

- **Cobra:** `github.com/spf13/cobra` (pulled in via `cobra-cli init`/`add`). `cobra-cli` is a generator: `go install github.com/spf13/cobra-cli@latest`; run `cobra-cli init` then `cobra-cli add <name>`.
- **main.go:** Must call `cmd.Execute()` (or equivalent) so the root Cobra command runs. No "OK" or other 1.1-only behavior.
- **Build:** `go build -o licensegen` must succeed. `go run main.go [ls|get|write]` should work.
- **Binary name:** The built binary remains `licensegen`; `licensegen ls`, `licensegen get`, `licensegen write` must be the invoked forms.

### Architecture Compliance

- **Layout:** `main.go` at repo root; `cmd/root.go`, `cmd/ls.go`, `cmd/get.go`, `cmd/write.go`. No `internal/` in 1.2 (core-architectural-decisions, project-structure-boundaries).
- **CLI structure:** `main.go` → `cmd.Execute()`; root and subcommands in `cmd/`; `help` provided by the root (starter-template-evaluation).
- **Boundaries:** Commands in `cmd/`; no SPDX, config, or cache. `internal/` appears in later stories.

### Library / Framework Requirements

- **Cobra:** Use `github.com/spf13/cobra` as brought in by the generator. No minimum version specified; `cobra-cli@latest` and `go get`/`go mod` will pull a compatible Cobra.
- **cobra-cli:** `github.com/spf13/cobra-cli@latest`; used only to scaffold. Do not add `cobra-cli` as a module dependency; it is a dev/CLI tool.
- **Viper:** Not used in 1.2. Do not add `--viper` to `cobra-cli init`.

### File Structure Requirements

| Path           | Purpose                                      | This story   |
|----------------|----------------------------------------------|-------------|
| `main.go`      | Entrypoint; calls `cmd.Execute()`            | **Modify**  |
| `cmd/root.go`  | Root command, global flags, `Execute()`      | **Create**  |
| `cmd/ls.go`    | `ls` subcommand, registered under root       | **Create**  |
| `cmd/get.go`   | `get` subcommand, registered under root      | **Create**  |
| `cmd/write.go` | `write` subcommand, registered under root    | **Create**  |
| `go.mod`       | Module; will gain `cobra`                    | **Modify**  |
| `go.sum`       | Sums for deps                               | **Modify**  |
| `internal/`    | SPDX, config [Growth], cache [Growth]        | **Do not add** |

- **Do not** create `internal/` or any `internal/*` files. No SPDX, config, or cache in 1.2.

### Testing Requirements

- **This story:** No mandatory `*_test.go`. AC is verified by: `./licensegen ls`, `./licensegen get`, `./licensegen write` run; `--help` works for root and each subcommand. If you add tests (e.g. that `Execute` runs or that commands are registered), place `*_test.go` next to the package (e.g. `cmd/root_test.go`). Keep them light; project-context says command wiring can have lighter coverage.

### Previous Story Intelligence (1.1)

- **1.1 produced:** `go.mod`, `main.go` (prints "OK" to stdout, exits 0), `.gitignore` at repo root. No `cmd/` or `internal/`.
- **Integrating with 1.2:** `main.go` must be changed from the 1.1 dummy to the Cobra entrypoint. Either (a) `cobra-cli init` overwrites `main.go` with a `cmd.Execute()` style, or (b) you manually replace the body to call `cmd.Execute()` and remove the "OK" logic. Do not keep both the "OK" print and Cobra; the 1.1 AC is superseded by 1.2’s CLI structure.
- **Build and binary:** 1.1 established `go build -o licensegen` and `./licensegen`. After 1.2, `./licensegen` with no args typically runs the root (often shows help); `./licensegen ls|get|write` run the subcommands. The 1.1 "OK" behavior is no longer required.

### Latest Tech Information

- **cobra-cli (spf13):** `cobra-cli init` creates `main.go` (or leaves it) and `cmd/root.go`; `cobra-cli add <name>` creates `cmd/<name>.go` and registers under `rootCmd` by default. Do not use `--viper` for MVP. Use `--viper` only when adding config in Growth.
- **Cobra:** Standard layout is `main` → `cmd.Execute()`; `rootCmd.Execute()` typically handles `os.Exit` on error. For 1.2, default Cobra behavior is sufficient; explicit exit code handling is 1.3.

### Project Context Reference

- **bmad_docs/project-context.md:**
  - **Cobra:** `github.com/spf13/cobra`; CLI via `cobra-cli init` and `cobra-cli add <cmd>`.
  - **Layout:** `cmd/root.go`, `cmd/ls.go`, `cmd/get.go`, `cmd/write.go`; one file per subcommand; `main.go` calls `cmd.Execute()`.
  - **Help:** Use Cobra’s default `help`; do not replace or duplicate.
  - **Viper:** `--viper` on `cobra-cli init` only when adding config [Growth]; omit in 1.2.

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- **1-2 (2026-01-20):** AC#1: `go install` cobra-cli; `cobra-cli init` (no `--viper`) overwrote `main.go` with `cmd.Execute()`; `cobra-cli add ls|get|write` created `cmd/ls.go`, `cmd/get.go`, `cmd/write.go` with `rootCmd.AddCommand`. Verified: `./licensegen ls|get|write` run; `--help` for root and subcommands. No `*_test.go` per story.

### File List

- `main.go` (modified)
- `cmd/root.go` (created)
- `cmd/ls.go` (created)
- `cmd/get.go` (created)
- `cmd/write.go` (created)
- `go.mod` (modified)
- `go.sum` (modified)
- `LICENSE` (created by cobra-cli init)

## Change Log

- 2026-01-20: Story 1.2 implemented — Cobra scaffold, `cmd/root.go`, `cmd/ls.go`, `cmd/get.go`, `cmd/write.go`; `main.go` → `cmd.Execute()`; subcommands and `--help` verified.
