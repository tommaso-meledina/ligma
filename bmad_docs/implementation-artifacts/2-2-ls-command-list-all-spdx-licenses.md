# Story 2.2: ls command — list all SPDX licenses

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **user**,
I want **to run `ligma ls` and see all available SPDX licenses in the terminal**,
so that **I can choose which license to use (FR1, FR25)**.

## Acceptance Criteria

1. **Given** the SPDX client that can fetch the license list  
   **When** I run `ligma ls` and the SPDX source is reachable  
   **Then** the CLI prints the list of license identifiers (e.g. SPDX IDs, one per line or a clear, consistent format) to stdout and exits with 0  
   **When** the SPDX source is unreachable (network failure, 4xx/5xx, timeout), the CLI prints a clear, actionable error to stderr and exits with 3 (FR15, NFR-I1, NFR-I2)  
   **And** under typical network conditions, the command completes within 15 seconds (NFR-P1)

## Tasks / Subtasks

- [x] **Task 1: Wire ls to internal/spdx** (AC: #1)
  - [x] In `cmd/ls.go`, call the SPDX client (from 2.1) to fetch the license list. Use the hardcoded list URL (same as 2.1). Pass a context or rely on the client’s 30s timeout.
  - [x] On success: format the list (e.g. one SPDX ID per line, or ID and name; keep format consistent and human-readable) and print to **stdout only**. Exit 0 via root (return nil from `RunE`).
- [x] **Task 2: Map SPDX client errors to exit 3** (AC: #1)
  - [x] If the client returns an error (network failure, 4xx/5xx, timeout, invalid JSON), `ls` must return an error that root maps to **exit 3** (I/O or network). Use the existing `ErrIOOrNetwork` or equivalent from 1.3. Do **not** call `os.Exit` in `cmd/ls`; return the error.
  - [x] Print a clear, actionable error message to **stderr** before returning (or ensure root does; avoid double-print). E.g. “failed to fetch license list: …” (FR15).
- [x] **Task 3: Remove stub behavior for ls** (AC: #1)
  - [x] Replace any 1.3 stub (no-op or synthetic not-found) with the real SPDX fetch and output. `ls` takes no required args for MVP; invalid flags still → exit 1 via Cobra/root.
- [x] **Task 4: Verify NFR-P1 and help** (AC: #1)
  - [x] Under typical conditions, `ligma ls` should complete within 15 seconds (NFR-P1). The client’s 30s timeout is an upper bound; no additional delay. `ligma ls --help` continues to exit 0.

## Dev Notes

- **Pre-requisite:** Story 2.1 done: `internal/spdx` fetches and parses `licenses.json`; returns slice or error.
- **Output:** Human-oriented list to stdout (e.g. one ID per line). No `--json` in this story (Epic 6.1). Stderr only for errors.
- **Exit codes:** 0 = success; 3 = SPDX unreachable. Usage errors (invalid flags) → 1 via existing root logic. `ls` does not use exit 2 (not-found) for the list command.
- **URL:** Hardcoded; config override comes in 5.2.

### Project Structure Notes

- **This story:** Modify `cmd/ls.go` only. `internal/spdx` is unchanged. No `internal/cache` or `internal/config` yet.

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-2-list-licenses.md#story-22-ls-command-list-all-spdx-licenses]
- [Source: bmad_docs/planning-artifacts/architecture/project-structure-boundaries.md] — `cmd/ls` uses `internal/spdx`
- [Source: bmad_docs/project-context.md] — stdout for `ls`, stderr for errors, exit 3 for I/O/network

---

## Developer Context

### Technical Requirements

- **cmd/ls:** Import `internal/spdx`; call `FetchLicenseList` (or equivalent). Format and `fmt.Println`/`fmt.Fprintf` to stdout. On error: write to stderr, return `ErrIOOrNetwork` (or type root treats as 3).
- **Root:** Already maps `ErrIOOrNetwork` → `os.Exit(3)`. No change to root needed if the error type is consistent with 1.3.

### Architecture Compliance

- **Boundaries:** `cmd/ls` handles flags and orchestration; `internal/spdx` does HTTP and JSON. `ls` must not perform HTTP or JSON parse itself.

### File Structure Requirements

| Path           | Purpose                                      | This story   |
|----------------|----------------------------------------------|-------------|
| `cmd/ls.go`    | Call SPDX client, format list to stdout, map errors → 3 | **Modify**  |

### Testing Requirements

- Manual: `ligma ls` prints list and exits 0; with network down or invalid URL, exits 3 and error on stderr. Optional: `cmd/ls` unit test with mocked SPDX client.

### Previous Story Intelligence

- **2.1:** `internal/spdx` provides `FetchLicenseList` (or similar) and returns `([]License, error)`. 2.2 consumes it.
- **1.3:** `ls` had no required args; could be no-op. 2.2 replaces with real implementation. Root and `ErrIOOrNetwork` already in place.

### Project Context Reference

- **bmad_docs/project-context.md:** Stdout for `ls`; stderr for errors; exit 3 = I/O or network.

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- **2-2 (2026-01-20):** AC#1: `cmd/ls.go` calls `spdx.FetchLicenseList(ctx, spdx.DefaultListURL)` (or `lsListURLOverride` in tests). One SPDX ID per line to stdout; on fetch error returns `fmt.Errorf("%w: failed to fetch license list: %v", ErrIOOrNetwork, err)`. Replaced 1.3 stub with real fetch. `ls --help` and `ls` verified. Tests: success (httptest), fetch error → ErrIOOrNetwork.

### File List

- `cmd/ls.go` (modified)
- `cmd/ls_test.go` (modified)

## Change Log

- 2026-01-20: Story 2.2 implemented — `ls` wires to `spdx.FetchLicenseList`, one ID per line to stdout; fetch errors → `ErrIOOrNetwork` (exit 3); `lsListURLOverride` for tests; `TestLsRunE_Success`, `TestLsRunE_FetchErrorMapsToIOOrNetwork`.
