# Story 1.3: Refine exit codes

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **developer**,
I want **the CLI to use consistent exit codes (0=success, 1=usage, 2=not found, 3=I/O or network)**,
so that **scripts and callers can reliably detect success and failure (FR16, FR17)**.

## Acceptance Criteria

1. **Given** the CLI with root, `ls`, `get`, and `write` commands  
   **When** a command completes successfully (e.g. `--help`), it exits with **0**  
   **When** a command fails due to usage (e.g. invalid flags or malformed args where defined), it exits with **1**  
   **When** a command fails due to not found (e.g. unknown license ID—stub behavior is acceptable for now), it exits with **2**  
   **When** a command fails due to I/O or network (stub behavior acceptable for now), it exits with **3**  
   **Then** all process exit is via `os.Exit` or equivalent from `main`/root; no prompts or interactive input are required (FR17)  
   **And** the exit code semantics are documented in code or a short comment for future implementation

## Tasks / Subtasks

- [x] **Task 1: Define exit-code semantics and error sentinels** (AC: #1)
  - [x] Add a short comment or block in `cmd/` (e.g. in `root.go` or a small `exit.go`) documenting: **0** = success, **1** = usage, **2** = not found, **3** = I/O or network; all `os.Exit` only from `main`/root.
  - [x] Define error sentinels or types in `cmd` (e.g. `ErrUsage`, `ErrNotFound`, `ErrIOOrNetwork`) so that `root` can classify returned errors. Use `errors.Is` for classification. Any error that is not `ErrNotFound` or `ErrIOOrNetwork` is treated as usage → exit **1**.
- [x] **Task 2: Centralize exit logic in root** (AC: #1)
  - [x] In `cmd.Execute()` (or the root’s execution path): run `rootCmd.Execute()`. If it returns `nil` → `os.Exit(0)`.
  - [x] If it returns an error: print the error message to **stderr** (project-context: all errors to stderr), then:  
     - `errors.Is(err, ErrNotFound)` → `os.Exit(2)`  
     - `errors.Is(err, ErrIOOrNetwork)` → `os.Exit(3)`  
     - else (usage or any other) → `os.Exit(1)`
  - [x] Ensure **no** `os.Exit` in `RunE` of `ls`/`get`/`write`; subcommands **return** errors and let root exit. When `internal/` is added later, it must never call `os.Exit`; only `main`/root may.
- [x] **Task 3: Wire `get` and `write` for usage and not-found** (AC: #1)
  - [x] **get:** Require one positional arg (license ID). Use Cobra `Args: cobra.ExactArgs(1)` (or equivalent) so missing or wrong count is validated by Cobra → returned error → root maps to **1**.
  - [x] **get RunE (stub):** For any provided ID, return `ErrNotFound` (and write a short message to stderr, e.g. `license not found: <id>`). This satisfies “unknown license ID—stub behavior is acceptable for now” and proves the **2** path.
  - [x] **write:** Require one or two args: `<id>` and optional `[path]`. Use `cobra.MinimumNArgs(1)` and `cobra.MaximumNArgs(2)` so malformed args → **1**.
  - [x] **write RunE (stub):** For any provided ID, return `ErrNotFound` and write to stderr. Proves the **2** path.
- [x] **Task 4: Wire I/O-or-network (3) and optional stub** (AC: #1)
  - [x] Implement the **3** path in root (as in Task 2). Document in the exit semantics comment when **3** is used: I/O or network failures (e.g. SPDX fetch failure, file write failure) when those are implemented in later stories.
  - [x] **Stub for 3 (optional):** Add a hidden or dev-only flag (e.g. `--simulate-io-error` on `get` or `write`) that makes `RunE` return `ErrIOOrNetwork` to prove the **3** path. If omitted, documenting the **3** semantics and the root mapping is sufficient; AC allows “stub behavior acceptable for now.”
- [x] **Task 5: Ensure `ls` and usage errors exit 1** (AC: #1)
  - [x] **ls:** No required positional args for MVP. Invalid flags or malformed args (if any are defined) → Cobra returns an error → root maps to **1**. `ls` RunE can remain a no-op (success, **0**).
  - [x] Verify: `ligma --help`, `ligma ls --help`, etc. still exit **0**; `ligma get` (missing ID), `ligma get foo` (stub not-found), `ligma write` (missing ID), `ligma write bar` (stub not-found) produce **1** and **2** as above.
- [x] **Task 6: Verify all four exit codes** (AC: #1)
  - [x] **0:** `./ligma --help`, `./ligma ls`, `./ligma ls --help`, `./ligma get --help`, `./ligma write --help`.
  - [x] **1:** `./ligma get` (missing arg), `./ligma write` (missing arg), `./ligma get --unknown-flag x`, or similar usage errors.
  - [x] **2:** `./ligma get MIT`, `./ligma write MIT` (stub: any ID is “not found”).
  - [x] **3:** If `--simulate-io-error` was added: `./ligma get --simulate-io-error x` (or equivalent). Otherwise, confirm the comment and root logic for **3** are in place.

## Dev Notes

- **Pre-requisite:** Stories 1.1 and 1.2 done: `main.go` → `cmd.Execute()`, `cmd/root.go`, `cmd/ls.go`, `cmd/get.go`, `cmd/write.go`. No `internal/` yet.
- **Exit only from main/root:** project-context and implementation-patterns: `os.Exit` only from `main`/root. Subcommands and, later, `internal/` must **return** errors; they must **not** call `os.Exit`. This keeps a single place to map errors → 1/2/3 and avoids surprises when `internal/spdx` is added.
- **Errors to stderr:** All error messages must go to stderr. Before `os.Exit(1|2|3)`, ensure the error is printed to stderr (root can do this when handling the error from `Execute`). Subcommands may also `fmt.Fprintln(os.Stderr, ...)` before returning the error; avoid double-printing if root already does.
- **Cobra and usage:** Cobra validates `Args` and flags. When validation fails, `Execute()` returns an error. We do **not** need to type-assert Cobra’s error; any error that is not `ErrNotFound` or `ErrIOOrNetwork` is treated as **1** (usage).
- **Stub not-found:** For `get` and `write`, the stub has no SPDX client, so every ID is “not found.” Returning `ErrNotFound` for any ID demonstrates the **2** path. In later stories, SPDX will determine not-found for real.
- **No `internal/` in 1.3:** Error sentinels live in `cmd`. When `internal/spdx` is added, it will return plain `error`s; `cmd` will wrap or translate them to `ErrNotFound` or `ErrIOOrNetwork` as appropriate.

### Project Structure Notes

- **This story:** Only `cmd/` and possibly `main.go` (if `Execute` lives in `cmd`, `main` may stay as-is). Optional: `cmd/exit.go` for sentinels and a helper (e.g. `exitCodeFrom(err) int`) if it keeps `root.go` cleaner. No `internal/`.
- **Future:** `internal/` packages must never call `os.Exit`; they return errors. `cmd` and root are responsible for mapping to 1/2/3 and exiting.

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-1-cli-foundation-help.md#story-13-refine-exit-codes]
- [Source: bmad_docs/planning-artifacts/architecture/core-architectural-decisions.md] — exit 0/1/2/3, errors to stderr
- [Source: bmad_docs/planning-artifacts/architecture/implementation-patterns-consistency-rules.md] — Exit codes, Stderr
- [Source: bmad_docs/project-context.md] — Exit codes, Errors to stderr, “Exit from main/root only; no os.Exit inside internal/”

---

## Developer Context

### Technical Requirements

- **Exit codes:** **0** = success; **1** = usage; **2** = not found; **3** = I/O or network. No other codes without updating project-context and architecture.
- **os.Exit:** Only from `main` or the root command’s `Execute` path. Never from `RunE` of subcommands or from `internal/`.
- **Stderr:** All error output to stderr. Stdout only for `ls` and `get` (and `--json` in Growth); `write` does not write to stdout.
- **FR16:** CLI exits 0 on success, non-zero on error. **FR17:** Non-interactive, no prompts; scriptable.

### Architecture Compliance

- **core-architectural-decisions:** “Error handling: Stderr for errors; exit 0 = success, 1 = usage, 2 = not found, 3 = I/O or network.”
- **implementation-patterns:** “Exit codes: 0 success, 1 usage, 2 not found, 3 I/O or network; always to `os.Exit` (or equivalent) from `main`/root.”
- **project-context:** “Exit from `main`/root only; no `os.Exit` inside `internal/` libraries.”

### Library / Framework Requirements

- **Cobra:** Use `RunE` to return `error`; do not call `os.Exit` inside `RunE`. Use `Args` (e.g. `cobra.ExactArgs(1)`, `cobra.MinimumNArgs(1)`, `cobra.MaximumNArgs(2)`) so Cobra enforces usage before `RunE`. Cobra’s own validation errors will be returned from `rootCmd.Execute()` and handled by root as **1**.
- **stdlib `errors`:** Use `errors.Is` to detect `ErrNotFound` and `ErrIOOrNetwork`. Sentinels can be `var ErrNotFound = errors.New("...")` or custom types implementing `error` and used with `errors.Is`.

### File Structure Requirements

| Path            | Purpose                                                    | This story   |
|-----------------|------------------------------------------------------------|-------------|
| `cmd/root.go`   | `Execute()`: run `rootCmd.Execute()`, classify error, `os.Exit(0\|1\|2\|3)`; exit semantics comment | **Modify**  |
| `cmd/exit.go`   | Optional: `ErrUsage`, `ErrNotFound`, `ErrIOOrNetwork`; optional `exitCodeFrom(err)` | **Create** (optional) |
| `cmd/get.go`    | `Args: cobra.ExactArgs(1)`; RunE stub returns `ErrNotFound`; optional `--simulate-io-error` | **Modify**  |
| `cmd/write.go`  | `Args: cobra.MinimumNArgs(1), MaximumNArgs(2)`; RunE stub returns `ErrNotFound`; optional `--simulate-io-error` | **Modify**  |
| `cmd/ls.go`     | No required args; invalid flags → Cobra error → 1; RunE no-op → 0 | **Modify** (if needed) |
| `main.go`       | Unchanged if `cmd.Execute()` owns all exit logic           | **No change** or **Modify** only if exit moves here |
| `internal/`     | Not used in 1.3                                           | **Do not add** |

### Testing Requirements

- **This story:** Manually verify 0/1/2/3 as in Task 6. Optionally add `cmd/root_test.go` or `cmd/exit_test.go` to test `exitCodeFrom(err)` or the mapping logic. Keep tests light; project-context allows lighter coverage for command wiring.
- **Later:** When `internal/spdx` exists, integration tests will exercise real not-found and I/O/network; 1.3 only establishes the stub and the centralized exit behavior.

### Previous Story Intelligence (1.1, 1.2)

- **1.1:** `go.mod`, `main.go`, `.gitignore`. `main` originally printed “OK” and exited 0; 1.2 replaced that with `cmd.Execute()`.
- **1.2:** `main.go` calls `cmd.Execute()`. `cmd/root.go` has `Execute()` and `rootCmd`; `cmd/ls.go`, `cmd/get.go`, `cmd/write.go` are stubs (no-ops or minimal). Cobra’s default `Execute()` may currently call `os.Exit(1)` on any error, or the generator’s `main` might. For 1.3 we **replace** that with our own 1/2/3 logic in `cmd.Execute()` (or equivalent). Ensure the generator’s `os.Exit(1)` or similar is removed in favor of our centralized behavior.
- **get/write in 1.2:** They may not have had `Args` or `RunE` defined. 1.3 adds `Args` and `RunE` that return `ErrNotFound` in stubs. If 1.2 used `Run` instead of `RunE`, switch to `RunE` so we can return errors.

### Latest Tech Information

- **Cobra `RunE`:** Returning a non-nil `error` from `RunE` causes `Execute()` to return that error. Cobra does not call `os.Exit` itself; the caller (our `Execute` or `main`) does. Our `Execute` must therefore run `rootCmd.Execute()`, inspect the error, and call `os.Exit` with the correct code.
- **Cobra `Args`:** `cobra.ExactArgs(1)` fails when the number of positional args is not 1; `cobra.MinimumNArgs(1)` and `cobra.MaximumNArgs(2)` enforce 1–2 args. On failure, Cobra returns an error before `RunE` runs; we map that to **1**.

### Project Context Reference

- **bmad_docs/project-context.md:**
  - **Exit codes:** 0 = success, 1 = usage, 2 = not found, 3 = I/O or network. Exit from `main`/root only; no `os.Exit` inside `internal/`.
  - **Errors:** Always to stderr; stdout only for `ls` and `get` (and `--json` when used).
  - **Critical:** Do not introduce extra exit codes without updating project-context and architecture.

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- **1-3 (2026-01-20):** AC#1: `cmd/exit.go` with 0/1/2/3 semantics, `ErrNotFound`, `ErrIOOrNetwork`; `cmd.Execute()` classifies via `errors.Is`, prints to stderr, `os.Exit(0|1|2|3)`. `get`: `ExactArgs(1)`, RunE→`ErrNotFound` (wrap), hidden `--simulate-io-error`→`ErrIOOrNetwork`. `write`: `RangeArgs(1,2)`, RunE→`ErrNotFound`. `ls`: RunE no-op. `SilenceUsage`/`SilenceErrors` on get/write to avoid double-print. Verified 0/1/2/3 per Task 6.

### File List

- `cmd/exit.go` (created)
- `cmd/root.go` (modified)
- `cmd/get.go` (modified)
- `cmd/write.go` (modified)
- `cmd/ls.go` (modified)

## Change Log

- 2026-01-20: Story 1.3 implemented — exit 0/1/2/3, `cmd/exit.go` sentinels, centralized root `Execute()`; get/write Args+RunE stubs; `--simulate-io-error` on get; ls RunE no-op. Verified all four exit codes.
