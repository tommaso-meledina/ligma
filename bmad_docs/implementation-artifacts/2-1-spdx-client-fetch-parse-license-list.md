# Story 2.1: SPDX client — fetch and parse license list

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **developer**,
I want **`internal/spdx` to fetch `licenses.json` from the official SPDX URL and parse it**,
so that **the `ls` command can display the license list (FR25)**.

## Acceptance Criteria

1. **Given** the official SPDX `licenses.json` URL (hardcoded for MVP)  
   **When** the client fetches the license list with a 30s timeout (NFR-I2)  
   **Then** it returns the parsed list (structs matching `license-list-data` layout) or an error if the request fails, times out, or the response is not valid JSON  
   **And** it uses `net/http` and `encoding/json` only; no third-party HTTP client (Architecture)  
   **And** the package lives under `internal/spdx` (e.g. `client.go`)

## Tasks / Subtasks

- [x] **Task 1: Create internal/spdx package** (AC: #1)
  - [x] Create `internal/spdx/` with `client.go`. Define a struct for the license list matching SPDX `license-list-data` JSON (e.g. `licenses.json` root structure: array of objects with at least `licenseId`; include fields needed for `ls`).
  - [x] Define a function to fetch and parse the list (e.g. `FetchLicenseList(ctx, listURL string) ([]License, error)` or `Client.FetchList(ctx) ([]License, error)`). The function must **return** errors; it must **not** call `os.Exit` (project-context: no `os.Exit` inside `internal/`).
- [x] **Task 2: Implement HTTP fetch with 30s timeout** (AC: #1)
  - [x] Use `net/http` only. Create `http.Client` with `Timeout: 30 * time.Second` (NFR-I2).
  - [x] GET the hardcoded URL: `https://raw.githubusercontent.com/spdx/license-list-data/main/json/licenses.json` (match PRD/cli-tool-specific-requirements default; this is MVP hardcoded).
  - [x] On non-2xx status, network error, or timeout: return a descriptive error. Do not parse body on 4xx/5xx as JSON.
- [x] **Task 3: Parse JSON with encoding/json** (AC: #1)
  - [x] Use `encoding/json` to decode the response. SPDX `licenses.json` has a structure like `{"licenses": [{"licenseId":"MIT", ...}, ...]}` or similar; follow the actual `license-list-data` schema. Return an error on invalid JSON.
  - [x] Return the parsed list (or a slice usable by `ls`). The client does not print to stdout/stderr; it only returns data or error.
- [x] **Task 4: Add tests** (AC: #1)
  - [x] Add `internal/spdx/client_test.go`. Test: successful parse (e.g. with a small fixture or mocked HTTP), timeout behavior, non-2xx handling, invalid JSON. `go test ./internal/spdx/...`. Prioritise coverage for the SPDX client (project-context).

## Dev Notes

- **Boundary:** `internal/spdx` is the **only** place that does HTTP and JSON for SPDX (project-structure-boundaries). No HTTP or JSON parsing for SPDX in `cmd/` or anywhere else.
- **URL:** Hardcoded for MVP. Epic 5.2 will introduce config override; for 2.1, use the official URL above.
- **SPDX IDs:** Use as-is in structs (e.g. `Apache-2.0`); no normalization (project-context).
- **No `os.Exit` in internal/spdx:** All errors are returned. `cmd/ls` (Story 2.2) will map errors to exit 2/3 via root.

### Project Structure Notes

- **This story:** Create `internal/spdx/client.go` and `internal/spdx/client_test.go`. No changes to `cmd/` or `main.go` in 2.1.

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-2-list-licenses.md#story-21-spdx-client-fetch-and-parse-license-list]
- [Source: bmad_docs/planning-artifacts/architecture/core-architectural-decisions.md] — SPDX: `net/http`, `encoding/json`, 30s timeout, `licenses.json`
- [Source: bmad_docs/planning-artifacts/architecture/project-structure-boundaries.md] — `internal/spdx`; only place for SPDX HTTP/JSON
- [Source: bmad_docs/project-context.md] — SPDX, no `os.Exit` in `internal/`, `*_test.go` next to package

---

## Developer Context

### Technical Requirements

- **net/http:** `http.Get` or `http.Client{Timeout: 30*time.Second}.Get`. No third-party HTTP client.
- **encoding/json:** `json.Decoder` or `json.Unmarshal` on response body. Struct tags to match SPDX `license-list-data` field names.
- **Context:** Accept `context.Context` (e.g. for timeout/retry) or at least honor a 30s client timeout. The 30s is mandatory (NFR-I2).

### Architecture Compliance

- **core-architectural-decisions:** “SPDX: `licenses.json` for list; … `net/http` GET; 30s timeout; `encoding/json`; structs for list.”
- **project-structure-boundaries:** “SPDX: `internal/spdx` is the only place that does HTTP and JSON parse for SPDX.”

### Library / Framework Requirements

- **Stdlib only for SPDX:** `net/http`, `encoding/json`, `context`, `time`. No `goresty`, `go-retryablehttp`, etc.

### File Structure Requirements

| Path                         | Purpose                                      | This story   |
|------------------------------|----------------------------------------------|-------------|
| `internal/spdx/client.go`    | Fetch `licenses.json`, parse, return list or error | **Create**  |
| `internal/spdx/client_test.go` | Tests for list fetch, timeout, 4xx/5xx, invalid JSON | **Create**  |

### Testing Requirements

- **Placement:** `client_test.go` next to `client.go`. `go test ./internal/spdx/...`.
- **Coverage:** Prioritise SPDX client; project threshold 90%. Use httptest or fixtures for HTTP.

### Previous Story Intelligence (Epic 1)

- **1.1–1.3:** `cmd/root`, `cmd/ls`, `cmd/get`, `cmd/write` exist; exit 0/1/2/3 from root. `internal/` did not exist. 2.1 introduces `internal/spdx`; `cmd/ls` will call it in 2.2.

### Latest Tech Information

- **SPDX license-list-data:** `https://raw.githubusercontent.com/spdx/license-list-data/main/json/licenses.json`. Schema: see `license-list-data` repository; typically `{"licenses": [{"licenseId": "...", "name": "..."}, ...]}`. Use structs that match the fields `ls` needs (at least `licenseId`).

### Project Context Reference

- **bmad_docs/project-context.md:** SPDX: `net/http` + `encoding/json`; 30s timeout; `licenses.json`; only `internal/spdx` does SPDX HTTP/JSON. No `os.Exit` in `internal/`.

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- **2-1 (2026-01-20):** AC#1: `internal/spdx/client.go` with `License` (licenseId, name), `listResponse`, `DefaultListURL`, `FetchLicenseList(ctx, listURL) ([]License, error)`. `http.Client` 30s timeout; GET listURL; non-2xx → drain body, return error; `encoding/json` decode. `client_test.go`: success (httptest+fixture), 404, 5xx, invalid JSON, context canceled. Coverage 93.3%. No `cmd/` or `main.go` changes.

### File List

- `internal/spdx/client.go` (created)
- `internal/spdx/client_test.go` (created)

## Change Log

- 2026-01-20: Story 2.1 implemented — `internal/spdx` with FetchLicenseList (30s timeout, net/http, encoding/json), License struct, DefaultListURL; tests for success, 4xx/5xx, invalid JSON, context canceled.
