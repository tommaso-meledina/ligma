# Story 6.1: ls — JSON output [Growth]

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **user**,
I want **to run `ligma ls --json` and get the license list in JSON**,
so that **I can script and integrate with other tools (FR4)**.

## Acceptance Criteria

1. **Given** the `ls` command  
   **When** I run `ligma ls --json`  
   **Then** the CLI prints the license list as JSON to stdout (e.g. an array of objects with at least SPDX ID; exact field names TBD in implementation) and exits with 0  
   **And** without `--json`, `ls` keeps the existing human-readable format

## Tasks / Subtasks

- [x] **Task 1: Add --json flag to ls** (AC: #1)
  - [x] In `cmd/ls.go`, add a `--json` flag (boolean). Add the flag in the `ls` command (project-context: flags in the command that uses them). Default `false`.
  - [x] When `--json` is true, output JSON; when false, use the existing human-readable format (e.g. one ID per line).
- [x] **Task 2: JSON output format** (AC: #1)
  - [x] Print a JSON array to stdout. Each element is an object with at least the SPDX ID (e.g. `{"id":"MIT"}` or `{"licenseId":"MIT"}`; exact keys TBD). You may include `name` or other fields from the SPDX list if the `ls` data source provides them. Ensure valid JSON (e.g. `json.Marshal` or `json.Encoder`).
  - [x] Only stdout for the list; errors to stderr. With `--json`, stdout is machine-readable; do not mix in human-only text.
- [x] **Task 3: Preserve non-JSON behavior** (AC: #1)
  - [x] When `--json` is false or omitted, `ls` behaves exactly as before: human-oriented list (e.g. one ID per line). No change to exit codes: 0 on success, 3 on SPDX/cache failure.
- [x] **Task 4: Verify** (AC: #1)
  - [x] `ligma ls --json` → valid JSON array to stdout, 0. `ligma ls` → same human format as pre–6.1, 0. `ligma ls --help` shows `--json`.

## Dev Notes

- **Pre-requisite:** `ls` already fetches the list (from SPDX or cache). 6.1 only changes the output format when `--json` is set. If `ls` is wired to cache (5.3), the cache returns the same list shape; use it for both human and JSON.
- **Field names:** “Exact field names TBD in implementation.” At minimum include SPDX ID; `licenseId` or `id` are both acceptable. Prefer consistency with SPDX `license-list-data` if the list struct has `licenseId`/`name`.

### Project Structure Notes

- **This story:** Modify `cmd/ls.go` only. No new packages. `internal/spdx` and `internal/cache` unchanged.

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-6-listing-enhancements-growth.md#story-61-ls-json-output]
- [Source: bmad_docs/project-context.md] — `--json` on `ls` and `get` [Growth]; stdout for `ls` and `get` (and `--json` when used)
- [Source: bmad_docs/planning-artifacts/prd/cli-tool-specific-requirements.md] — `ls --json` for machine-readable list

---

## Developer Context

### Technical Requirements

- **encoding/json:** `json.Marshal` or `json.NewEncoder(os.Stdout).Encode(...)`. Ensure no trailing newline issues if piping; one JSON value (the array) on stdout.
- **Flag:** `lsCmd.Flags().BoolP("json", "j", false, "output as JSON")` or similar. Use `lsCmd.Flags().GetBool("json")` in `RunE`.

### Architecture Compliance

- **project-context:** “Flags: `--popular`, `--json` on `ls` and `get` [Growth].” 6.1 adds `--json` to `ls`.

### File Structure Requirements

| Path           | Purpose                                      | This story   |
|----------------|----------------------------------------------|-------------|
| `cmd/ls.go`    | Add `--json`, branch output: JSON vs human   | **Modify**  |

### Testing Requirements

- Manual: `ligma ls --json` is valid JSON and parseable. Optional: `jq` or `json.Marshal` round-trip. Without `--json`, output unchanged.

### Previous Story Intelligence

- **2.2, 5.3:** `ls` fetches list and prints human format. 6.1 adds a second branch for `--json`; same data source, different formatting.

### Project Context Reference

- **bmad_docs/project-context.md:** `--json` on `ls` and `get`; stdout for `ls` and `get` (and `--json` when set).

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- ls: `--json`/`-j` in init; runLs: GetBool("json") → json.NewEncoder(os.Stdout).Encode(list) (spdx.License has licenseId/name); else existing one-ID-per-line. ls_test: TestLsRunE_JSON (capture stdout, parse JSON, assert licenseId/name); TestLsRunE_NoJSONHumanFormat (assert not JSON, one line per ID).

### File List

- cmd/ls.go (modified)
- cmd/ls_test.go (modified)
