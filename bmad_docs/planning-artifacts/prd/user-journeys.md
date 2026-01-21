# User Journeys

## 1. Primary – Success (Alex)

**Opening scene:** Alex has just created a new side project and needs an MIT license. They’re used to copy-pasting from a random site or another repo and want to be done in seconds.

**Rising action:** They run `ligma ls`, spot MIT in the list, run `ligma get MIT` to confirm the text, then `ligma write MIT`. A correct `LICENSE` appears in the project root.

**Climax:** They’re done in three commands (or two if they skip `get`). No copy-paste, no hunting for an authoritative source.

**Resolution:** They commit and push. The license step is off their mind.

---

## 2. Primary – Edge / Errors (Riley)

**Opening scene:** Riley runs `ligma write` with no arguments, or types a wrong ID (e.g. `MIT-0` instead of `MIT`), or the SPDX source is unreachable.

**Rising action:** The CLI exits with an error. Riley runs `ligma help` to see usage. The error message explains what went wrong (e.g. unknown license ID, or “no favorite set; pass an ID or set favorite in ~/.ligma/config.json”) and how to fix it.

**Climax:** They correct the ID or set a favorite, run `ligma write` (or `ligma write <id>`), and get a valid `LICENSE`.

**Resolution:** Clear errors plus `help` are enough to recover without hunting through docs or the web.

---

## 3. Configuring User (Growth)

**Opening scene:** A user runs ligma for the first time. The tool creates the `~/.ligma/` directory and `~/.ligma/config.json`. They want a default license (favorite) and maybe an alias so they don’t have to remember SPDX IDs.

**Rising action:** They open the config, set `favorite` to e.g. `MIT` and optionally an alias like `"apache": "Apache-2.0"`. They run `ligma write` with no args and get an error until favorite is set; after setting it, `ligma write` works. `ls` and `get` use the local cache when available, so repeat runs are faster.

**Climax:** From then on, `ligma write` with no args writes their default license. For other licenses they use `ligma write Apache-2.0` or, with aliases, `ligma write apache`.

**Resolution:** One-time setup; daily use is one command for the common case and simple commands for the rest.

---

## 4. Power User / Scripter (Sam)

**Opening scene:** Sam maintains many OSS repos and often adds or refreshes licenses. They want to script it and avoid any interactive prompts.

**Rising action:** They use `ligma ls --popular` or `ligma ls MIT` to find IDs, then `ligma write MIT` (or `ligma write MIT /path/to/repo`) in a loop or one-liner. Every call is non-interactive with deterministic exit codes. The local cache makes repeat `ls` and `get` calls faster; setting `cache_ttl` to `0` in config forces a fresh fetch when needed.

**Climax:** A small script or `find` plus `ligma write` updates many repos. No manual steps, no prompts.

**Resolution:** License management scales with the number of repos; the CLI stays a reliable building block for automation.

---

## Journey Requirements Summary

| Journey | Capabilities |
|---------|--------------|
| **Primary – Success** | `ls` (full list), `get <id>`, `write <id> [path]`; correct SPDX text; non-interactive. |
| **Primary – Edge** | `help`; clear, actionable errors (unknown ID, no favorite, invalid args); stable exit codes. |
| **Configuring** | Create `~/.ligma/` directory and `config.json` on first run; `favorite`; aliases; `write` with no args when favorite is set; local cache for `ls` and `get`. |
| **Power user** | `ls [filter]`, `ls --popular [filter]`; `cache_ttl: 0` to bypass cache; `write <id> [path]` scriptable; no prompts; deterministic behavior; cache speeds repeat `ls`/`get`. |
