# Story 3.1: SPDX client — fetch license details by ID

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **developer**,
I want **`internal/spdx` to fetch `details/{id}.json` for a given SPDX ID and return the license text**,
so that **the `get` command can output the full license (FR25)**.

## Acceptance Criteria

1. **Given** an SPDX license ID (e.g. `MIT`, `Apache-2.0`) and the official SPDX `details/{id}.json` URL template (hardcoded for MVP)  
   **When** the client fetches the details with a 30s timeout (NFR-I2)  
   **Then** it returns the license text (e.g. from the `licenseText` field) or an error if the request fails, times out, returns 404, or the response is not valid JSON  
   **And** it uses the same `net/http` and `encoding/json` approach as the list client; no third-party HTTP client  
   **And** SPDX IDs are used as-is (e.g. `Apache-2.0`); no normalization for the URL path

## Tasks / Subtasks

- [x] **Task 1: Add FetchLicenseDetails to internal/spdx** (AC: #1)
  - [x] In `internal/spdx/client.go`, add a function (e.g. `FetchLicenseDetails(ctx, detailsURLTemplate, id string) (string, error)` or `Client.FetchDetails(ctx, id) (string, error)`). The template has a `{id}` placeholder; replace with the given ID **as-is** (no case or format normalization; project-context).
  - [x] Hardcoded template for MVP: `https://raw.githubusercontent.com/spdx/license-list-data/main/json/details/{id}.json`.
  - [x] Return only the license text (e.g. `licenseText` from JSON). Return an error on HTTP 404, 4xx/5xx, network error, timeout, or invalid JSON. Do **not** call `os.Exit` (internal/ must only return errors).
- [x] **Task 2: Implement HTTP fetch with 30s timeout** (AC: #1)
  - [x] Use `net/http` and the same 30s timeout pattern as the list client. GET the resolved URL (template with `{id}` replaced by the given ID).
  - [x] On 404: return a distinct error so `get` can map to exit 2 (not found). On other 4xx/5xx, network, or timeout: return an error that `get` will map to exit 3 (I/O/network). The client does not need to distinguish 2 vs 3; it can return a generic error and let the command decide, or return a typified error (e.g. `ErrNotFound`) for 404. Architecture: 404 / “ID not in list” → 2; “unreachable, 5xx, timeout” → 3.
- [x] **Task 3: Parse JSON and extract licenseText** (AC: #1)
  - [x] Use `encoding/json` to decode. SPDX `details/{id}.json` includes a `licenseText` (or similar) field. Extract and return it as a string. Invalid JSON or missing `licenseText` → return error.
- [x] **Task 4: Tests** (AC: #1)
  - [x] In `internal/spdx/client_test.go`, add tests for: success and `licenseText` extraction, 404 → not-found-style error, 5xx/timeout/network → I/O-style error, invalid JSON. Reuse httptest or fixtures as in 2.1.

## Dev Notes

- **Boundary:** All SPDX HTTP and JSON stay in `internal/spdx`. `get` (3.2) and `write` (4.1) will call this; they must not perform HTTP or JSON for SPDX.
- **ID as-is:** Use the ID exactly for the URL path (e.g. `details/MIT.json`, `details/Apache-2.0.json`). No lowercasing or normalization (project-context).
- **404 vs 5xx/timeout:** 404 (or “not in list”) should be detectable so `get` can exit 2. Other failures → 3. The client can return `ErrNotFound` for 404 and a generic or `ErrIOOrNetwork`-like for the rest; or return a struct/error type that carries “not found” vs “I/O” so `get` can map correctly.

### Project Structure Notes

- **This story:** Extend `internal/spdx/client.go` and `client_test.go`. No `cmd/` changes.

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-3-view-license.md#story-31-spdx-client-fetch-license-details-by-id]
- [Source: bmad_docs/planning-artifacts/architecture/core-architectural-decisions.md] — `details/{id}.json`, `licenseText`, 30s, `net/http`, `encoding/json`
- [Source: bmad_docs/project-context.md] — SPDX IDs as-is; only `internal/spdx` does SPDX HTTP/JSON

---

## Developer Context

### Technical Requirements

- **URL template:** `.../details/{id}.json`. Replace `{id}` with the argument. IDs like `Apache-2.0` are valid in paths; no encoding beyond using the string as-is (path segment may need URL-encoding if the ID contained slashes; SPDX IDs in practice do not—use as-is unless the HTTP client requires encoding).
- **404:** `http.Res.StatusCode == 404` → not-found. Return an error that `get` can treat as `ErrNotFound` (exit 2).

### Architecture Compliance

- **core-architectural-decisions:** “`details/{id}.json` for `licenseText`; … 30s timeout; `encoding/json`.”
- **project-structure-boundaries:** SPDX HTTP/JSON only in `internal/spdx`.

### File Structure Requirements

| Path                         | Purpose                                      | This story   |
|------------------------------|----------------------------------------------|-------------|
| `internal/spdx/client.go`    | Add `FetchLicenseDetails` (or equivalent)    | **Modify**  |
| `internal/spdx/client_test.go` | Tests for details fetch, 404, 5xx, timeout, invalid JSON | **Modify**  |

### Testing Requirements

- Test 404 → not-found error; 5xx/timeout → I/O error; valid JSON → `licenseText` returned. Coverage focus on SPDX client.

### Previous Story Intelligence

- **2.1:** List client in same package; same `net/http` and 30s pattern. Reuse client or timeout config. 3.1 adds a second entry point for details.

### Project Context Reference

- **bmad_docs/project-context.md:** SPDX IDs as-is; `internal/spdx` only for SPDX HTTP/JSON; no `os.Exit` in internal.

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- FetchLicenseDetails(ctx, detailsURLTemplate, id) in internal/spdx/client.go; DefaultDetailsURLTemplate const; ErrNotFound for 404. 30s timeout, net/http, encoding/json. 404→ErrNotFound, other 4xx/5xx/network/timeout→wrapped error, invalid JSON or missing licenseText→error. Tests: success, 404→ErrNotFound, 5xx, invalid JSON, missing licenseText, context canceled.

### File List

- internal/spdx/client.go (modified)
- internal/spdx/client_test.go (modified)
