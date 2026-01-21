# Story 7.2: get — JSON output [Growth]

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **user**,
I want **to run `ligma get <id> --json` and receive the license in JSON**,
so that **I can script and pipe structured data (FR7)**.

## Acceptance Criteria

1. **Given** the `get` command  
   **When** I run `ligma get <SPDX-ID> --json`  
   **Then** the CLI prints the license (e.g. `licenseText` and optionally ID) as JSON to stdout and exits with 0  
   **And** without `--json`, `get` keeps the existing plain-text output

## Tasks / Subtasks

- [x] **Task 1: Add --json flag to get** (AC: #1)
  - [x] In `cmd/get.go`, add a `--json` flag (boolean). Add it to the `get` command (project-context: flags in the command that uses them). Default `false`.
  - [x] When `--json` is true, output JSON; when false, keep the existing behavior (plain license text to stdout).
- [x] **Task 2: JSON output format** (AC: #1)
  - [x] Print a single JSON object to stdout. At least `licenseText` (the full license text) and optionally `id` or `licenseId` (the SPDX ID). Example: `{"id":"MIT","licenseText":"..."}`. Use `encoding/json`. Ensure valid JSON; one object. No extra human-only text when `--json`.
- [x] **Task 3: Preserve non-JSON behavior** (AC: #1)
  - [x] When `--json` is false or omitted, `get` prints only the plain license text to stdout, as before. Exit codes unchanged: 0 success, 1 usage, 2 not found, 3 I/O.
- [x] **Task 4: Alias and cache** (AC: #1)
  - [x] If 7.1 (alias resolution) is implemented, resolve the ID/alias before fetch; the JSON can include the **resolved** SPDX ID in the `id` field. Cache (5.3) is transparent; same output whether from cache or network.
- [x] **Task 5: Verify** (AC: #1)
  - [x] `ligma get MIT --json` → valid JSON object with `licenseText` (and optionally `id`) to stdout, 0. `ligma get MIT` → plain text only, 0. `ligma get --help` shows `--json`.

## Dev Notes

- **Pre-requisite:** `get` already fetches license text (from SPDX or cache). 7.2 adds a `--json` branch: same data, different format. Works with or without 7.1 (alias).
- **Field names:** “e.g. `licenseText` and optionally ID” — at minimum `licenseText`; `id` or `licenseId` recommended for scripting.

### Project Structure Notes

- **This story:** Modify `cmd/get.go` only. No `internal/` changes.

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-7-view-write-enhancements-growth.md#story-72-get-json-output]
- [Source: bmad_docs/project-context.md] — `--json` on `ls` and `get` [Growth]; stdout for `get` (and `--json` when set)
- [Source: bmad_docs/planning-artifacts/prd/cli-tool-specific-requirements.md] — `get --json` for structured output

---

## Developer Context

### Technical Requirements

- **encoding/json:** `json.Marshal(struct{ID string `json:"id"`; LicenseText string `json:"licenseText"`}{...})` or similar. One JSON object to stdout.
- **Flag:** `getCmd.Flags().BoolP("json","j", false, "output as JSON")` or similar.

### Architecture Compliance

- **project-context:** “Flags: `--popular`, `--json` on `ls` and `get` [Growth].” 7.2 adds `--json` to `get`.

### File Structure Requirements

| Path           | Purpose                                      | This story   |
|----------------|----------------------------------------------|-------------|
| `cmd/get.go`   | Add `--json`; branch: JSON object vs plain text | **Modify**  |

### Testing Requirements

- Manual: `ligma get MIT --json` is valid JSON; `jq .licenseText` or similar works. Without `--json`, plain text only.

### Previous Story Intelligence

- **3.2, 5.3, 7.1:** `get` fetches text and prints. 7.2 adds formatting branch; same fetch, different output.

### Project Context Reference

- **bmad_docs/project-context.md:** `--json` on `get`; stdout for `get` (and `--json` when set).

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- get: `--json` / `-j` (BoolP, default false). If true: `json.Marshal(struct{ID, LicenseText})` with resolved `id`, write to stdout. If false: plain `text` only. Format: `{"id":"<resolved>","licenseText":"..."}`. Tests: TestGetRunE_JSON (capture stdout, unmarshal, assert id+licenseText), TestGetRunE_NoJSONPlainText (plain only, not JSON).

### File List

- cmd/get.go (modified)
- cmd/get_test.go (modified)
