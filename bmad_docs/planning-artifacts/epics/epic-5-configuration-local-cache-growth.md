# Epic 5: Configuration & Local Cache [Growth]

User can configure favorite, aliases, SPDX URLs, cache TTL; ls and get use a local cache when valid.
**FRs covered:** FR18, FR19, FR20, FR21, FR22, FR23, FR24, FR26 | **NFRs:** NFR-S1

## Story 5.1: ~/.ligma/ and config.json with Viper

As a **user**,
I want **`~/.ligma/` and `~/.ligma/config.json` to be created when absent, at the start of a run that uses config**,
So that **I can later add favorite, aliases, URLs, and cache_ttl without manual setup (FR18)**.

**Acceptance Criteria:**

**Given** the internal/config package using Viper and a schema with `favorite`, `aliases`, `spdx_list_url`, `spdx_get_url_template`, `cache_ttl` (FR19, FR20, FR21, FR24 as config support)
**When** `config.Load()` runs and either `~/.ligma/` or `~/.ligma/config.json` is missing
**Then** the package creates `~/.ligma/` and `~/.ligma/config.json` (empty `{}` or minimal valid JSON) before reading; subsequent Load returns the parsed config
**And** the CLI only reads `~/.ligma/config.json` and writes license files to user-specified paths; it does not collect, store, or transmit secrets or PII (NFR-S1)

## Story 5.2: SPDX client uses config URLs when present

As a **user**,
I want **to override SPDX list and per-license URLs via `config.json`**,
So that **I can use a mirror or custom SPDX source (FR21, FR26)**.

**Acceptance Criteria:**

**Given** config with `spdx_list_url` and `spdx_get_url_template` set
**When** `ls` or `get` (or `write`) fetches from SPDX
**Then** the SPDX client uses those URLs instead of the hardcoded defaults
**When** those keys are absent or empty, the client uses the hardcoded URLs
**And** config is loaded at the start of the command (triggering create-if-absent from 5.1 when first used)

## Story 5.3: Cache layer for ls and get

As a **user**,
I want **`ls` and `get` to use a local cache under `~/.ligma/_cache/` when the cache is valid**,
So that **repeated runs are fast and work offline within the TTL (FR22, FR23, FR24)**.

**Acceptance Criteria:**

**Given** `internal/cache` that wraps SPDX list and details fetches
**When** `ls` or `get` runs, the cache checks `_cache/list.json` or `_cache/details/<id>.json`; if the file exists and its mtime is within `cache_ttl` (from config), use cached data and do not perform HTTP
**When** the file is missing or stale, fetch from SPDX, write to the cache (best-effort: on write failure, still return fetched data), then return
**When** `cache_ttl` is `0` in config, always fetch and do not use cache (FR23); when `cache_ttl` is null/absent, use an implementation default (e.g. 24h)
**Then** `_cache` and `_cache/details` are created on first use; `write` does not use the cache
**And** there is no `--no-cache` flag; bypass is via `cache_ttl: 0` in config (Architecture)

---
