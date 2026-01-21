# Story 5.2: SPDX client uses config URLs when present [Growth]

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **user**,
I want **to override SPDX list and per-license URLs via `config.json`**,
so that **I can use a mirror or custom SPDX source (FR21, FR26)**.

## Acceptance Criteria

1. **Given** config with `spdx_list_url` and `spdx_get_url_template` set  
   **When** `ls` or `get` (or `write`) fetches from SPDX  
   **Then** the SPDX client uses those URLs instead of the hardcoded defaults  
   **When** those keys are absent or empty, the client uses the hardcoded URLs  
   **And** config is loaded at the start of the command (triggering create-if-absent from 5.1 when first used)

## Tasks / Subtasks

- [x] **Task 1: SPDX client accepts URL parameters** (AC: #1)
  - [x] In `internal/spdx`, ensure the list and details functions accept URL (or template) as an argument (e.g. `FetchLicenseList(ctx, listURL string)` and `FetchLicenseDetails(ctx, template, id string)`). 2.1 and 3.1 may already have these; if they used hardcoded URLs inside, refactor to take URLs from the caller.
  - [x] Hardcoded defaults: `https://raw.githubusercontent.com/spdx/license-list-data/main/json/licenses.json` and `https://raw.githubusercontent.com/spdx/license-list-data/main/json/details/{id}.json` (PRD defaults). These are used when config values are absent or empty.
- [x] **Task 2: Commands load config and pass URLs to SPDX** (AC: #1)
  - [x] **cmd/ls:** At start of `RunE`, call `config.Load()` (5.1). If `config.SPDXListURL()` or equivalent is non-empty, pass it to the SPDX list fetch; otherwise pass the hardcoded list URL. Loading config triggers create-if-absent on first run.
  - [x] **cmd/get:** At start of `RunE`, call `config.Load()`. If `config.SPDXGetURLTemplate()` is non-empty, pass it to the SPDX details fetch; otherwise pass the hardcoded details template. Pass the SPDX ID as-is.
  - [x] **cmd/write:** Same as `get`: `config.Load()` and pass `spdx_get_url_template` or hardcoded. `write` does not use the list URL.
- [x] **Task 3: Absent or empty means use hardcoded** (AC: #1)
  - [x] If `spdx_list_url` or `spdx_get_url_template` is missing, empty string, or null, use the hardcoded default. Do not treat empty as an error.
- [x] **Task 4: Verify and tests** (AC: #1)
  - [x] Manual: with custom `spdx_list_url` / `spdx_get_url_template` in `~/.ligma/config.json`, `ls` and `get`/`write` use them. With keys removed or empty, they use hardcoded. Unit tests: pass a mock config or URLs into the client to assert the requested URL is used.

## Dev Notes

- **Pre-requisite:** 5.1 done: `internal/config` with `Load()`, `spdx_list_url`, `spdx_get_url_template`. Create-if-absent on first `Load()`.
- **Boundary:** `internal/spdx` still does all HTTP and JSON; it does **not** read config. `cmd/*` load config and pass URLs. This keeps SPDX independent of Viper/config.
- **Template:** `spdx_get_url_template` has `{id}`; replace with the SPDX ID as-is (no normalization).

### Project Structure Notes

- **This story:** Modify `internal/spdx` (signatures to accept URLs if not already) and `cmd/ls.go`, `cmd/get.go`, `cmd/write.go` to call `config.Load()` and pass URLs.

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-5-configuration-local-cache-growth.md#story-52-spdx-client-uses-config-urls-when-present]
- [Source: bmad_docs/planning-artifacts/prd/cli-tool-specific-requirements.md] — Config Schema, `spdx_list_url`, `spdx_get_url_template`
- [Source: bmad_docs/project-context.md] — Config supplies URLs; [Growth] cache wraps SPDX, config supplies URLs

---

## Developer Context

### Technical Requirements

- **Config API:** `config.Load()` returns a struct or accessors: `SPDXListURL() string`, `SPDXGetURLTemplate() string`. Empty or unset returns `""` so callers can fall back to hardcoded.
- **Call order:** `config.Load()` first (which may create `~/.ligma/` and `config.json`), then call SPDX with the resolved URL or default.

### Architecture Compliance

- **project-structure-boundaries:** “Config [Growth]: … called at startup by commands that need it.” “SPDX: … [Growth] config supplies URLs.”

### File Structure Requirements

| Path               | Purpose                                      | This story   |
|--------------------|----------------------------------------------|-------------|
| `internal/spdx/client.go` | Accept listURL, detailsURLTemplate as args   | **Modify**  |
| `cmd/ls.go`        | config.Load(), pass list URL to SPDX         | **Modify**  |
| `cmd/get.go`       | config.Load(), pass details template to SPDX | **Modify**  |
| `cmd/write.go`     | config.Load(), pass details template to SPDX | **Modify**  |

### Testing Requirements

- Tests: SPDX client receives and uses the URL/template passed in. Integration: config with custom URLs and `ls`/`get` hit the expected endpoint (or use httptest).

### Previous Story Intelligence

- **5.1:** `config.Load()` and create-if-absent. 5.2 is the first consumer of `spdx_list_url` and `spdx_get_url_template`. **2.1, 3.1:** Refactor to accept URLs from caller if they are currently hardcoded inside.

### Project Context Reference

- **bmad_docs/project-context.md:** Config [Growth] supplies URLs; SPDX only in `internal/spdx`; config supplies URLs when cache/config are used.

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- ls: config.Load() first; url = lsListURLOverride else cfg.SPDXListURL else spdx.DefaultListURL. get: config.Load() first; template = cfg.SPDXGetURLTemplate else spdx.DefaultDetailsURLTemplate. write: config.Load() first; same template resolution. Absent/empty → hardcoded. All cmd tests that run ls/get/write set config.SetConfigDirOverride(t.TempDir()). TestLsRunE_UsesConfigListURL: config.json with spdx_list_url=httptest.URL, no lsListURLOverride, run ls→success. internal/spdx unchanged (already takes URL/template).

### File List

- cmd/ls.go (modified)
- cmd/ls_test.go (modified)
- cmd/get.go (modified)
- cmd/get_test.go (modified)
- cmd/write.go (modified)
- cmd/write_test.go (modified)
- cmd/root_test.go (modified)
