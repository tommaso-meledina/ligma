# Project Context Analysis

## Requirements Overview

**Functional Requirements**

- **License Listing:** FR1 (list all SPDX licenses); FR2 (case-insensitive filter [Growth]); FR3 (`--popular` [Growth]); FR4 (`--json` for `ls` [Growth]).
- **License Viewing:** FR5 (`get <id>` to stdout); FR6 (alias resolution for `get` [Growth]); FR7 (`--json` for `get` [Growth]).
- **License Writing:** FR8 (`write <id>` by SPDX ID); FR9 (optional `[path]`, default `LICENSE` in cwd); FR10 (alias resolution for `write` [Growth]); FR11 (`write` with no args using `favorite` [Growth]); FR12 (error when `write` no args and no `favorite` [Growth]).
- **Help & Error Handling:** FR13 (`help`); FR14 (clear errors for unknown/invalid IDs); FR15 (clear error when SPDX unreachable); FR16 (exit 0 on success, non-zero on error); FR17 (non-interactive, scriptable, no prompts).
- **Configuration [Growth]:** FR18 (`~/.ligma/` and `config.json` create-if-absent); FR19 (`favorite` for `write` with no args); FR20 (aliases); FR21 (override SPDX URLs via config).
- **Local Cache [Growth]:** FR22 (cache for `ls` and `get` under `~/.ligma/`, respect `cache_ttl`); FR23 (bypass cache via `cache_ttl: 0` in config; no `--no-cache` flag); FR24 (`cache_ttl` in config; enforcement TBD).
- **SPDX Data Access:** FR25 (list and details from SPDX; hardcoded URLs in MVP); FR26 (SPDX URLs from config [Growth]).

**Non-Functional Requirements**

- **Performance:** NFR-P1 (`ls`/`get` within ~15s under typical network); NFR-P2 (negligible startup before first network call).
- **Security:** NFR-S1 (no secrets or PII; only `~/.ligma/config.json` and license file writes).
- **Integration (SPDX):** NFR-I1 (correct fetch/parse when SPDX reachable and format expected); NFR-I2 (clear failure and ~30s timeout when unreachable).

**Scale & Complexity**

- Primary domain: CLI.
- Complexity: low.
- Estimated architectural components: CLI entry and subcommands; SPDX client; config loader (Growth); cache (Growth); file I/O for `write`.

## Technical Constraints & Dependencies

- **Go:** Single binary; no extra runtime.
- **SPDX:** `license-list-data` JSON (list + `details/{id}.json`). MVP: hardcoded URLs; Growth: from config.
- **No daemon:** On-demand fetch or cache; offline when cache valid (Growth).
- **Filesystem:** `write` to `LICENSE` or `[path]`; Growth: `~/.ligma/` for config and cache.

## Cross-Cutting Concerns

- **CC1 — Scriptability:** Exit codes, stderr vs stdout, no prompts; `--json` (Growth) for machine-readable `ls`/`get` (FR4, FR7, FR16, FR17).
- **CC2 — SPDX integration:** Single source for list and text; shared fetch/parse, error and timeout handling (FR25, FR26, NFR-I1, NFR-I2).
- **CC3 — Config and cache lifecycle [Growth]:** Create-if-absent for `~/.ligma/` and `config.json`; TTL (`cache_ttl: 0` = always fetch, bypassing cache); resolution order for ID vs alias and `favorite` (FR18–FR24, FR10, FR11, FR12).
