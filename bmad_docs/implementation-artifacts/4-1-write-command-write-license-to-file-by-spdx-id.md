# Story 4.1: write command — write license to file by SPDX ID

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **user**,
I want **to run `ligma write <id>` or `ligma write <id> <path>` and have the license written to a file**,
so that **I can add a license to my project (FR8, FR9)**.

## Acceptance Criteria

1. **Given** the SPDX client that can fetch license details by ID  
   **When** I run `ligma write <SPDX-ID>` and the SPDX source is reachable  
   **Then** the CLI writes the full license text to `LICENSE` in the current directory and exits with 0  
   **When** I run `ligma write <SPDX-ID> <path>`  
   **Then** the CLI writes the full license text to `<path>` and exits with 0; if the file exists, it is overwritten  
   **When** the ID is unknown or invalid, the CLI prints a clear, actionable error to stderr and exits with 2 (FR14)  
   **When** the SPDX source is unreachable or file write fails (e.g. permission, disk), the CLI prints a clear error to stderr and exits with 3  
   **And** `write` requires at least one argument (the SPDX ID); otherwise the CLI prints usage to stderr and exits with 1

## Tasks / Subtasks

- [x] **Task 1: Wire write to internal/spdx and resolve path** (AC: #1)
  - [x] In `cmd/write.go`, call the SPDX client’s `FetchLicenseDetails` (from 3.1) with the first positional arg (SPDX ID). Use the hardcoded details URL template.
  - [x] **Path rules:** If only one arg: write to `LICENSE` in the current directory. If two args: write to the second arg as `<path>`. Use `cobra.MinimumNArgs(1)` and `cobra.MaximumNArgs(2)` (from 1.3; adjust if 1.3 used different bounds). Resolve `LICENSE` as `./LICENSE` or `filepath.Join(cwd, "LICENSE")`; resolve `<path>` as given (relative or absolute).
- [x] **Task 2: Write license text to file** (AC: #1)
  - [x] Use `os.WriteFile` or `ioutil.WriteFile` to write the license text to the chosen path. Overwrite if the file exists. Permissions: 0644 or equivalent. On write failure (permission, disk full, etc.): print error to stderr and return an error that root maps to **exit 3** (I/O or network). Do **not** call `os.Exit` in `cmd/write`.
- [x] **Task 3: Map client and write errors to exit 2 and 3** (AC: #1)
  - [x] **Not found (ID unknown/invalid):** client returns not-found → print clear error to stderr (FR14), return `ErrNotFound` → exit **2**.
  - [x] **SPDX unreachable or file write failure:** return `ErrIOOrNetwork` (or equivalent) → exit **3**. This includes: network failure, 5xx, timeout, and local `WriteFile` errors.
- [x] **Task 4: Remove stub and verify usage** (AC: #1)
  - [x] Replace the 1.3 stub that returned `ErrNotFound` for any ID with the real SPDX fetch and file write. `ligma write` (no args) → exit 1. `ligma write id path1 path2` (three args) → exit 1 if `MaximumNArgs(2)`; otherwise adjust bounds to match AC (“at least one” = ID; “optional path” = 0 or 1 extra). Ensure `ligma write --help` exits 0.

## Dev Notes

- **Pre-requisite:** Story 3.1 done: `internal/spdx.FetchLicenseDetails` returns `(string, error)`; not-found vs I/O as in 3.2.
- **Output:** `write` does **not** print to stdout (project-context, PRD). Only stderr for errors.
- **Overwrite:** If the target file exists, overwrite it. No confirmations (FR17: non-interactive).
- **Exit codes:** 0 = success; 1 = usage; 2 = not found; 3 = I/O or network (including file write failure).

### Project Structure Notes

- **This story:** Modify `cmd/write.go` only. No new packages. `internal/spdx` is reused; no `internal/cache` or `internal/config` (Epic 5).

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-4-write-license-to-file.md#story-41-write-command-write-license-to-file-by-spdx-id]
- [Source: bmad_docs/project-context.md] — `write` only to files, no stdout; stderr for errors; 2/3

---

## Developer Context

### Technical Requirements

- **Path:** One arg → `LICENSE` in cwd; two args → second is path. Use `os.Getwd()` for cwd when defaulting to `LICENSE`. `filepath` for joining; support relative and absolute paths.
- **WriteFile:** `os.WriteFile(path, []byte(text), 0644)`. On error, treat as I/O → 3.

### Architecture Compliance

- **Boundaries:** `cmd/write` orchestrates; fetches via `internal/spdx`; does file I/O in `cmd` (or a small helper). File write is not SPDX—it’s local I/O; it can live in `cmd`. No `internal/` for “write file” unless you introduce a tiny helper; keeping it in `cmd` is fine.

### File Structure Requirements

| Path           | Purpose                                      | This story   |
|----------------|----------------------------------------------|-------------|
| `cmd/write.go` | Call SPDX details, resolve path, `WriteFile`; map not-found→2, I/O→3 | **Modify**  |

### Testing Requirements

- Manual: `ligma write MIT` → `LICENSE` in cwd, 0; `ligma write MIT ./custom.txt` → `custom.txt`, 0; overwrite; `ligma write BadId` → 2; write to read-only or invalid path → 3; `ligma write` → 1.

### Previous Story Intelligence

- **3.1, 3.2:** `FetchLicenseDetails` and error semantics. 4.1 reuses the same client and error mapping. **1.3:** `write` had `MinimumNArgs(1)`, `MaximumNArgs(2)` and stub `ErrNotFound`; 4.1 replaces stub with fetch+write and adds file-I/O→3.

### Project Context Reference

- **bmad_docs/project-context.md:** `write` only to files; stdout only for `ls` and `get`; 2 = not found, 3 = I/O or network.

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- runWrite: writeFetchDetails (spdx.FetchLicenseDetails) with DefaultDetailsURLTemplate; 1 arg→filepath.Join(Getwd(),"LICENSE"), 2→args[1]; os.WriteFile 0644. spdx.ErrNotFound→ErrNotFound; fetch or WriteFile error→ErrIOOrNetwork. writeFetchDetails var for tests. write_test: NotFound, OneArg (temp dir + LICENSE), TwoArgs (out.txt), WriteFails (path=dir). root_test: WriteNoArgs, WriteThreeArgs → 1.

### File List

- cmd/write.go (modified)
- cmd/write_test.go (modified)
- cmd/root_test.go (modified)
