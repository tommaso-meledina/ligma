# Story 3.2: get command — output license text by SPDX ID

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **user**,
I want **to run `ligma get <id>` and see the full license text in the terminal**,
so that **I can read or pipe a license before writing it (FR5)**.

## Acceptance Criteria

1. **Given** the SPDX client that can fetch license details by ID  
   **When** I run `ligma get <SPDX-ID>` with a valid ID and the SPDX source is reachable  
   **Then** the CLI prints the full license text to stdout and exits with 0  
   **When** the ID is unknown or invalid (e.g. 404 from SPDX, or ID not in the list), the CLI prints a clear, actionable error to stderr and exits with 2 (FR14)  
   **When** the SPDX source is unreachable (network failure, 4xx/5xx, timeout), the CLI prints a clear error to stderr and exits with 3  
   **And** `get` requires exactly one argument (the SPDX ID); otherwise the CLI prints usage to stderr and exits with 1

## Tasks / Subtasks

- [x] **Task 1: Wire get to internal/spdx FetchLicenseDetails** (AC: #1)
  - [x] In `cmd/get.go`, call the SPDX client's details fetch (from 3.1) with the single positional arg (SPDX ID). Use the hardcoded details URL template. `get` already has `Args: cobra.ExactArgs(1)` from 1.3.
  - [x] On success: print the license text to **stdout** only. Return nil from `RunE` → exit 0.
- [x] **Task 2: Map client errors to exit 2 and 3** (AC: #1)
  - [x] If the client returns a **not-found** error (404, or ID not in list): print a clear, actionable error to **stderr** (FR14), return `ErrNotFound` (or equivalent) so root exits **2**.
  - [x] If the client returns an **I/O or network** error (timeout, 5xx, network failure): print a clear error to stderr, return `ErrIOOrNetwork` (or equivalent) so root exits **3**.
  - [x] Do **not** call `os.Exit` in `cmd/get`; only return errors. Root (from 1.3) performs the mapping and `os.Exit`.
- [x] **Task 3: Remove stub not-found for get** (AC: #1)
  - [x] Replace the 1.3 stub that always returned `ErrNotFound` with the real SPDX fetch. Not-found is now determined by the client (404 or equivalent); all other failures → 3.
- [x] **Task 4: Verify usage and help** (AC: #1)
  - [x] `ligma get` (missing arg) and `ligma get a b` (too many args) must exit **1** via Cobra `ExactArgs(1)`. `ligma get --help` exits 0. `ligma get MIT` prints text and exits 0; `ligma get UnknownID` exits 2 with error on stderr.

## Dev Notes

- **Pre-requisite:** Story 3.1 done: `internal/spdx` provides `FetchLicenseDetails(id)` (or similar) and returns `(string, error)`; 404 → not-found error; others → I/O/network error.
- **Output:** Plain license text to stdout. No `--json` in this story (Epic 7.2). Stderr only for errors.
- **Exit codes:** 0 = success; 1 = usage (ExactArgs(1)); 2 = not found; 3 = I/O or network. Root already maps these.

### Project Structure Notes

- **This story:** Modify `cmd/get.go` only. `internal/spdx` and root are unchanged.

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-3-view-license.md#story-32-get-command-output-license-text-by-spdx-id]
- [Source: bmad_docs/project-context.md] — stdout for `get`, stderr for errors; exit 2 = not found, 3 = I/O/network

---

## Developer Context

### Technical Requirements

- **cmd/get:** Import `internal/spdx`; call `FetchLicenseDetails` with the first (and only) arg. Print result to stdout. On not-found error → return `ErrNotFound`; on other → return `ErrIOOrNetwork` (or types root maps to 2 and 3). Write error message to stderr before returning.
- **Root:** No change if 1.3 already maps `ErrNotFound`→2 and `ErrIOOrNetwork`→3.

### Architecture Compliance

- **Boundaries:** `cmd/get` orchestrates; `internal/spdx` does HTTP and JSON. `get` must not do SPDX HTTP/JSON itself.

### File Structure Requirements

| Path           | Purpose                                      | This story   |
|----------------|----------------------------------------------|-------------|
| `cmd/get.go`   | Call SPDX details, print text to stdout; map not-found→2, I/O→3 | **Modify**  |

### Testing Requirements

- Manual: `ligma get MIT` → text, 0; `ligma get BadId` → 2, stderr; with network down → 3, stderr; `ligma get` → 1.

### Previous Story Intelligence

- **3.1:** Client returns `(string, error)`; 404 or “not in list” → not-found; timeout/5xx/network → I/O. 3.2 consumes and maps to 2/3.
- **1.3:** `get` has `ExactArgs(1)` and stub `ErrNotFound`. 3.2 replaces stub with real fetch and correct 2/3 mapping.

### Project Context Reference

- **bmad_docs/project-context.md:** Stdout for `get`; stderr for errors; 2 = not found, 3 = I/O or network.

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- runGet: getFetchDetails(cmd.Context(), DefaultDetailsURLTemplate, id); spdx.ErrNotFound→ErrNotFound (exit 2), else wrap ErrIOOrNetwork (exit 3); success→os.Stdout.WriteString(text). getFetchDetails var for tests. get_test and root_test mock getFetchDetails for NotFound/Success; --simulate-io-error unchanged. Root prints err to stderr.

### File List

- cmd/get.go (modified)
- cmd/get_test.go (modified)
- cmd/root_test.go (modified)
