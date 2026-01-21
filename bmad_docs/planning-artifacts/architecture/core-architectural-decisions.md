# Core Architectural Decisions

## Decision Priority Analysis

**Critical (block implementation):** HTTP client stdlib `net/http`; JSON `encoding/json`; SPDX `licenses.json` + `details/{id}.json`; 30s timeout; exit codes 0/1/2/3, errors to stderr; no auth, NFR-S1; build via `go build` / `go install`.

**Important (shape architecture, incl. Growth):** Cache: file-based under `~/.ligma/_cache/` with mtime-based TTL; config (Growth): Viper, `~/.ligma/config.json`, schema as PRD.

**Deferred:** Release packaging (GoReleaser, multi-arch, CI); finer exit-code semantics.

## Data Architecture

- **SPDX:** `licenses.json` for list; `details/{id}.json` for `licenseText`; `encoding/json`; structs for list and details. Follow `license-list-data` layout.
- **Cache [Growth]:** File-based under `~/.ligma/_cache/`. **Layout:** `_cache/list.json` (license list), `_cache/details/<id>.json` (per-license). **Behaviour:** SPDX HTTP calls are wrapped by a cache layer. The layer checks for an existing cache file; if present, compares file **mtime** to `cache_ttl`. If mtime is within TTL, use cached data (no HTTP); otherwise fetch, save/overwrite, return. `cache_ttl` null: implementation default (e.g. 24h). `cache_ttl` 0: always bypass (always fetch). **Bypass:** No `--no-cache` flag; use `cache_ttl: 0` in config. **Resilience:** On cache write failure, still return fetched data; caching is best-effort. **Creation:** `_cache` and `_cache/details` created on first use.
- **Config [Growth]:** Viper; `~/.ligma/config.json`; schema as PRD (favorite, aliases, spdx_list_url, spdx_get_url_template, cache_ttl).

## Authentication & Security

- **Authentication:** None.
- **Security:** NFR-S1 only: no secrets/PII; read `~/.ligma/config.json` (Growth), write license files only to user-specified or default path.

## API & Communication Patterns

- **SPDX (external):** `net/http` GET; 30s timeout (NFR-I2); `encoding/json`. No third-party HTTP client.
- **Error handling:** Stderr for errors; exit 0 = success, 1 = usage, 2 = not found, 3 = I/O or network; other semantics TBD in implementation.
- **Output:** Stdout for `ls` and `get` (and `--json` in Growth); `write` only to files.

## Frontend Architecture

- Not applicable — CLI only. Output formats (human vs `--json`) per PRD (FR4, FR7).

## Infrastructure & Deployment

- **Build:** `go build -o ligma`; `go install` for local use.
- **Distribution:** Deferred; MVP: `go build` / `go install` only.

## Decision Impact / Implementation Order

1. Cobra scaffold: `go mod init`, `cobra-cli init`, `cobra-cli add ls|get|write`.
2. SPDX client: `net/http` + `encoding/json`, 30s timeout; list + details.
3. Commands: `ls` (list), `get` (license text to stdout), `write` (to file); errors and exit codes; `help` (Cobra).
4. Growth: `~/.ligma/` and config (Viper); cache layer (`_cache`, mtime, TTL); aliases, favorite, `--popular`, `--json`.

---
