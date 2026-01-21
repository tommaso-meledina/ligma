# Project Scoping & Phased Development

## MVP Strategy & Philosophy

**MVP Approach:** Problem-solving — replace copy-paste with a single command to get a correct `LICENSE` from SPDX.

**Resource:** Solo (hobby). MVP kept minimal so it can ship and be used quickly.

**Config in MVP:** No config file. SPDX URLs are hardcoded. Growth introduces `~/.ligma/` directory, `~/.ligma/config.json`, and create-if-absent.

---

## MVP Feature Set (Phase 1)

**Core journeys:** Primary success (ls → get/write), primary edge (errors + `help`).

**Must-have capabilities:**

- `ligma ls` — full list, human-oriented stdout
- `ligma get <id>` — license text to stdout, by SPDX ID only
- `ligma write <id> [path]` — write `LICENSE` (or `[path]`), by SPDX ID only
- `ligma help` — usage and help
- SPDX fetch from hardcoded URLs (e.g. `licenses.json` and `details/{id}.json`)
- Non-interactive, scriptable; clear errors to stderr; defined exit codes
- No config file, no aliases, no favorite, no `ls [filter]`, no `--popular`

---

## Post-MVP Features

**Phase 2 (Growth):**

- `~/.ligma/` directory and `config.json` (create-if-absent; schema and defaults: see **Config Schema**). Local cache for `ls` and `get`; `cache_ttl` (0 = bypass cache; enforcement TBD).
- `favorite` and `aliases`; `ligma write` with no args when favorite is set (error when not)
- `ligma ls [filter]` (case-insensitive); `ligma ls --popular [filter]`
- `--json` for `ls` and `get`
- SPDX URLs always read from config (no hardcoded URLs in Growth)

**Phase 3 (Expansion / later Growth):**

- Shell completion (bash, zsh, fish)
- Optional: conversational "help me choose" (Vision)

---

## Risk Mitigation Strategy

| Risk | Mitigation |
|------|------------|
| **Technical:** SPDX schema or URLs change | Hardcoded URLs in MVP; Growth uses config so users can override. Clear errors when fetch fails. |
| **Technical:** Cache/offline | Growth: local cache for `ls` and `get` with `cache_ttl` (0 = bypass cache); offline when cache is valid. TTL enforcement (e.g. stamp files) TBD. |
| **Market** | None; hobby. Ship MVP, use it, iterate. |
| **Resource** | Lean MVP; Growth only when useful. |
