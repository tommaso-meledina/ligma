# Story 5.3: Cache layer for ls and get [Growth]

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **user**,
I want **`ls` and `get` to use a local cache under `~/.ligma/_cache/` when the cache is valid**,
so that **repeated runs are fast and work offline within the TTL (FR22, FR23, FR24)**.

## Acceptance Criteria

1. **Given** `internal/cache` that wraps SPDX list and details fetches  
   **When** `ls` or `get` runs, the cache checks `_cache/list.json` or `_cache/details/<id>.json`; if the file exists and its mtime is within `cache_ttl` (from config), use cached data and do not perform HTTP  
   **When** the file is missing or stale, fetch from SPDX, write to the cache (best-effort: on write failure, still return fetched data), then return  
   **When** `cache_ttl` is `0` in config, always fetch and do not use cache (FR23); when `cache_ttl` is null/absent, use an implementation default (e.g. 24h)  
   **Then** `_cache` and `_cache/details` are created on first use; `write` does not use the cache  
   **And** there is no `--no-cache` flag; bypass is via `cache_ttl: 0` in config (Architecture)

## Tasks / Subtasks

- [x] **Task 1: Create internal/cache package** (AC: #1)
  - [x] Create `internal/cache/cache.go`. The cache wraps calls to `internal/spdx`: for the list, it checks `~/.ligma/_cache/list.json`; for details by ID, `~/.ligma/_cache/details/<id>.json`. SPDX IDs as-is for the path (e.g. `details/MIT.json`, `details/Apache-2.0.json`).
  - [x] **TTL:** Use file **mtime** (modification time), not birth/creation time (project-context). If `now - mtime < cache_ttl` (in seconds), treat as valid and read from file; otherwise fetch. `cache_ttl` 0: always fetch, do not read from cache. `cache_ttl` null/absent: use default (e.g. 24h = 86400). `cache_ttl` in config is in seconds (PRD).
- [x] **Task 2: List cache: _cache/list.json** (AC: #1)
  - [x] For `ls`: cache gets list URL from config (via caller). Check `_cache/list.json`. If exists and mtime within TTL: decode JSON and return the list (same shape as SPDX client), no HTTP. If missing or stale: call SPDX list fetch, then write result to `_cache/list.json` (best-effort; on write failure, still return the fetched data). Create `_cache/` if needed (`os.MkdirAll`). The cache layer lives under `~/.ligma/`; resolve with same logic as config (e.g. `UserHomeDir`).
- [x] **Task 3: Details cache: _cache/details/<id>.json** (AC: #1)
  - [x] For `get`: cache gets details URL template and ID. Check `_cache/details/<id>.json` (ID as-is; sanitize only if necessary for filesystem, e.g. path traversal; SPDX IDs are safe). If exists and mtime within TTL: decode and return `licenseText` (or the string). If missing or stale: call SPDX details fetch, write to `_cache/details/<id>.json` (best-effort; on write failure, return fetched data). Create `_cache` and `_cache/details` on first use.
- [x] **Task 4: Wire ls and get to cache; write unchanged** (AC: #1)
  - [x] **cmd/ls:** Instead of calling `internal/spdx` directly, call `internal/cache` (e.g. `cache.FetchList(ctx, cfg, listURL)`). The cache uses `config.Load()` for `cache_ttl` and for `_cache` base path; it calls SPDX when needed. `ls` continues to print to stdout and map errors to 3.
  - [x] **cmd/get:** Call `cache.FetchDetails(ctx, cfg, template, id)` (or equivalent). Same rules: cache uses config for TTL and path; on cache miss or stale, calls SPDX. `get` still maps not-found→2, I/O→3.
  - [x] **cmd/write:** Do **not** use the cache. `write` continues to call `internal/spdx` (or the same SPDX API the cache uses internally) directly. Only `ls` and `get` go through the cache (Architecture, PRD).
- [x] **Task 5: No --no-cache; config only** (AC: #1)
  - [x] Do not add a `--no-cache` flag. Bypass only via `cache_ttl: 0` in config (project-context, Architecture).
- [x] **Task 6: Tests** (AC: #1)
  - [x] Add `internal/cache/cache_test.go`. Test: hit (valid mtime), miss (no file), stale (old mtime), `cache_ttl` 0 (always fetch), write failure (still return fetched data), and `_cache`/`_cache/details` creation. Use a temp dir for `~/.ligma`.

## Dev Notes

- **Pre-requisite:** 5.1 (config, `cache_ttl`) and 5.2 (config URLs). 5.3 adds the cache; `ls` and `get` must already be using config for URLs so the cache can pass them through to SPDX on miss.
- **Boundary:** `internal/cache` wraps SPDX; it uses `internal/spdx` for HTTP. Only `internal/spdx` does the actual HTTP and JSON for SPDX. Cache does filesystem and JSON decode of cached files; that is not “SPDX HTTP” so it can live in `internal/cache`.
- **Resilience:** On cache file write failure (permission, disk), still return the fetched data; do not fail the command (project-context, core-architectural-decisions).

### Project Structure Notes

- **This story:** Create `internal/cache/cache.go`, `cache_test.go`; modify `cmd/ls.go` and `cmd/get.go` to use cache. `cmd/write.go` stays on SPDX only.

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-5-configuration-local-cache-growth.md#story-53-cache-layer-for-ls-and-get]
- [Source: bmad_docs/planning-artifacts/architecture/core-architectural-decisions.md] — Cache: `_cache/list.json`, `_cache/details/<id>.json`, mtime TTL, `cache_ttl` 0, resilience, no `--no-cache`
- [Source: bmad_docs/project-context.md] — Cache: mtime TTL; `cache_ttl` 0 = always fetch; on write failure still return data; create `_cache` and `_cache/details` on first use; no `--no-cache`

---

## Developer Context

### Technical Requirements

- **mtime:** `os.Stat(path)` then `ModTime()`. `time.Since(mtime) < time.Duration(cache_ttl)*time.Second` → valid. `cache_ttl == 0` → skip cache read, always fetch.
- **Paths:** `~/.ligma/_cache/list.json`, `~/.ligma/_cache/details/<id>.json`. Reuse home resolution from config. For `<id>`, use the ID as the filename (e.g. `MIT.json`, `Apache-2.0.json`); avoid path traversal (IDs are alphanumeric and `-`.`.).
- **Best-effort write:** After SPDX fetch, attempt `os.MkdirAll` and `os.WriteFile`. On any error, log or ignore and return the fetched data. Do not return an error to the caller for cache write failure.

### Architecture Compliance

- **project-structure-boundaries:** “Cache [Growth]: `internal/cache` wraps SPDX fetches; used by `ls` and `get` only.”
- **core-architectural-decisions:** Cache layout, mtime, TTL, `cache_ttl` 0, resilience, no `--no-cache`.

### File Structure Requirements

| Path                          | Purpose                                      | This story   |
|-------------------------------|----------------------------------------------|-------------|
| `internal/cache/cache.go`     | Wrap list and details; mtime TTL; best-effort write | **Create**  |
| `internal/cache/cache_test.go`| Tests for hit, miss, stale, TTL 0, write failure | **Create**  |
| `cmd/ls.go`                   | Use cache for list instead of SPDX directly  | **Modify**  |
| `cmd/get.go`                  | Use cache for details instead of SPDX directly | **Modify**  |

### Testing Requirements

- Prefer coverage for cache (project-context). Test mtime logic, TTL 0, and “write failure, still return data.”

### Previous Story Intelligence

- **5.1, 5.2:** Config with `cache_ttl`, `spdx_list_url`, `spdx_get_url_template`. Cache will call `config.Load()` (or receive config from `cmd`) to get TTL and to resolve `~/.ligma`. On miss, it needs the list URL and details template to call SPDX—same as 5.2.
- **2.1, 3.1:** SPDX client. Cache calls it when cache miss or stale. SPDX stays the single place for HTTP.

### Project Context Reference

- **bmad_docs/project-context.md:** Cache: `_cache/list.json`, `_cache/details/<id>.json`; mtime TTL; `cache_ttl` null = default, 0 = always fetch; on write failure still return data; create dirs on first use; no `--no-cache`.

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- internal/cache: FetchList, FetchDetails with mtime TTL; TTL(cfg) 0=always fetch, nil=86400. _cache/list.json, _cache/details/<id>.json; MkdirAll on first use; tryWrite best-effort. FetchListFn/FetchDetailsFn for tests. config.LigmaDir() for ~/.ligma. ls/get use cache; write unchanged (spdx only). get: removed getFetchDetails, use cache.FetchDetailsFn in tests. cache_test: Miss, Hit, Stale, TTLZero, WriteFailureStillReturns (list+details), TTL. config: LigmaDir(), TestLigmaDir_Override. No --no-cache.

### File List

- internal/cache/cache.go (created)
- internal/cache/cache_test.go (created)
- internal/config/config.go (modified; LigmaDir)
- internal/config/config_test.go (modified; TestLigmaDir_Override)
- cmd/ls.go (modified; cache.FetchList)
- cmd/get.go (modified; cache.FetchDetails, removed getFetchDetails)
- cmd/get_test.go (modified; cache.FetchDetailsFn)
- cmd/root_test.go (modified; cache.FetchDetailsFn)
