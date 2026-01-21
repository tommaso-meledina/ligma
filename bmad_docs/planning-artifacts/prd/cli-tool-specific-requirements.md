# CLI Tool Specific Requirements

## Project-Type Overview

ligma is a **non-interactive, scriptable CLI**: all commands run without prompts. It targets terminal and script use. No visual UI, no touch; we skip visual_design, ux_principles, and touch_interactions.

---

## Technical Architecture Considerations

- **Single binary** (Go): one `ligma` executable; no extra runtime.
- **SPDX as source of truth**: fetch and parse SPDX `license-list-data` JSON (from the official GitHub repo).
- **Config and cache (Growth):** `~/.ligma/` directory with `config.json`; created when absent. Local cache for `ls` and `get` under `~/.ligma/`; `cache_ttl` (0 = bypass cache; see Config Schema). Schema extensible for future keys.
- **No daemon:** once license data is available (from network or cache), `get` and `write` can run; offline when cache is valid.

---

## Command Structure

| Command | MVP | Growth | Notes |
|--------|-----|--------|-------|
| `ligma ls` | ✓ | ✓ | Full list; human-only output. |
| `ligma ls [filter]` | — | ✓ | Case-insensitive filter. |
| `ligma ls --popular [filter]` | — | ✓ | Top 5 (static), works with `[filter]`. |
| `ligma ls --json` | — | ✓ | JSON output for scripting. |
| `ligma get <id>` | ✓ | ✓ | Stdout; by ID (Growth: by alias). |
| `ligma get <id> --json` | — | ✓ | JSON output. |
| `ligma write <id> [path]` | ✓ | ✓ | By ID (Growth: by alias); `[path]` optional. |
| `ligma write` | — | ✓ | Uses `favorite`; error if unset. |
| `ligma help` | ✓ | ✓ | Usage and help. |

- **Subcommands:** `ls`, `get`, `write`, `help` (and any future ones) as explicit verbs; no interactive wizards.
- **Flags:** `--popular`, `--json` (Growth, for `ls` and `get`); config and `favorite` via `~/.ligma/config.json`. Cache bypass via `cache_ttl: 0` in config.

---

## Output Formats

| Command | MVP | Growth |
|---------|-----|--------|
| **`ls`** | Human-oriented list to stdout (e.g. ID and/or name per line). | Same + `--json` for machine-readable list. |
| **`get`** | Plain license text to stdout. | Same + `--json` for structured (e.g. id + licenseText). |
| **`write`** | File only; no stdout. Writes `LICENSE` in cwd or `[path]`. | Same. |
| **`help`** | Usage/help text to stdout. | Same. |

- **Stderr:** Errors and diagnostics to stderr; stdout stays parseable when `--json` is used.
- **Exit codes:** 0 on success; non-zero on error (conventions TBD in implementation).

---

## Config Schema

- **Directory and file:** `~/.ligma/` is a directory in `$HOME`. Config lives at `~/.ligma/config.json`. The local cache for `ls` and `get` is stored under `~/.ligma/` (exact layout TBD in implementation).
- **When it's created:** On every run, before doing anything else, the CLI checks for the `~/.ligma/` directory. If it's missing, it creates it. Then it checks for `~/.ligma/config.json`; if missing, it creates it with the default structure below. So the directory and config exist when the process runs; there is no "run without config" in Growth.
- **Implication:** `spdx_list_url` and `spdx_get_url_template` are always read from config (including right after it's created with defaults). There are no URL fallbacks hardcoded in the binary. `cache_ttl` controls how long `ls` and `get` results are served from cache; how it's enforced (e.g. stamp files) is TBD in implementation.

**Default config (created when the file is absent):**

```json
{
  "favorite": null,
  "aliases": [],
  "spdx_list_url": "https://raw.githubusercontent.com/spdx/license-list-data/main/json/licenses.json",
  "spdx_get_url_template": "https://raw.githubusercontent.com/spdx/license-list-data/main/json/details/{id}.json",
  "cache_ttl": null
}
```

- **`favorite`** (string | null): SPDX ID or alias used by `ligma write` when no argument is given. Default `null`; if `null` and `write` is called with no args, the CLI errors.
- **`aliases`** (object): Map from user-defined names to SPDX IDs (e.g. `"apache": "Apache-2.0"`). Used by `get` and `write` to resolve identifiers.
- **`spdx_list_url`** (string): URL of the SPDX license list (e.g. `licenses.json`). Used by `ls` to load the full list.
- **`spdx_get_url_template`** (string): URL template for a single license's JSON. The substring `{id}` is replaced by the SPDX ID (or the alias resolved to an ID) to fetch e.g. `.../details/MIT.json`.
- **`cache_ttl`** (number | null): TTL in seconds for `ls` and `get` cache entries. When `null`, implementation uses a default. When `0`, cache is always bypassed (always fetch from network). How TTL is enforced (e.g. stamp files on cached data) is TBD in implementation.

**Extensibility:** The schema may gain more keys later (e.g. cache subdir, `ls` defaults). New keys will be optional so old config files keep working.

---

## Scripting Support

- **Non-interactive:** No prompts, no pager, no interactive choices; all behavior driven by args and config.
- **Exit codes:** 0 = success; non-zero = failure; semantics to be defined (e.g. 1 = usage, 2 = not found, 3 = I/O error).
- **Machine-readable output (Growth):** `--json` on `ls` and `get` for scripting and tooling.
- **Pipes and redirects:** `get` and `ls` write to stdout; `write` only to files. Safe to pipe `get` and `ls` in scripts.

---

## Implementation Considerations

- **SPDX fetch and cache (Growth):** `ls` and `get` results are cached under `~/.ligma/`. `cache_ttl` in config controls use of cache (0 = always fetch). How TTL is enforced (e.g. stamp files) TBD. Offline: when cache is valid, `ls` and `get` work without network.
- **Config load order:** Global config only (`~/.ligma/config.json`); no project-local config for now. If added later, document precedence.
- **ID and alias resolution:** In Growth, resolve `favorite` and `aliases` before SPDX lookup; unknown ID/alias → clear error and non-zero exit.
- **Growth – even later:** **Shell completion** (bash, zsh, fish) for subcommands and license IDs; explicitly out of scope for MVP and initial Growth, to be planned in a later phase.
