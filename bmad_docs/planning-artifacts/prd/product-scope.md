# Product Scope

## MVP - Minimum Viable Product

- **`ligma ls`** — list all licenses (no `[filter]`, no `--popular`).
- **`ligma get <id>`** — output license text by SPDX ID only.
- **`ligma write <id> [path]`** — write `LICENSE` (or `[path]`) by SPDX ID only.
- **`ligma help`** — help.
- **Data:** SPDX license-list-data (JSON).

No config file, no aliases, no favorite, no `ls [filter]`, no `ls --popular`.

## Growth (Post-MVP)

- **`~/.ligma/` directory** (not just a config file): created when absent; holds `config.json` and local cache. Config at `~/.ligma/config.json`.
- **Local cache:** `ls` and `get` results are cached under `~/.ligma/`; subsequent invocations use cache when valid. User bypasses cache by setting `cache_ttl` to `0` in config. `cache_ttl` in config enforces TTL (how enforced TBD, e.g. stamp files).
- **Config:** `favorite`, `aliases`, `spdx_list_url`, `spdx_get_url_template`, `cache_ttl`; user-configurable.
- **Aliases** in config; `get` and `write` accept aliases.
- **`favorite`** in config: `ligma write` with no args uses favorite; if no args and no `favorite` → error.
- **`ligma ls [filter]`** — optional, case-insensitive filter; works with `--popular`.
- **`ligma ls --popular [filter]`** — top 5 (static for now; generic so we can plug in better filtering later).

## Vision (Future)

- **Conversational / "help me choose" mode** — out of scope for now; possible later.
