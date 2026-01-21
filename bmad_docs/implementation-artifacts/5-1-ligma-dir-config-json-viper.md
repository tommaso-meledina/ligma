# Story 5.1: ~/.ligma/ and config.json with Viper [Growth]

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **user**,
I want **`~/.ligma/` and `~/.ligma/config.json` to be created when absent, at the start of a run that uses config**,
so that **I can later add favorite, aliases, URLs, and cache_ttl without manual setup (FR18)**.

## Acceptance Criteria

1. **Given** the `internal/config` package using Viper and a schema with `favorite`, `aliases`, `spdx_list_url`, `spdx_get_url_template`, `cache_ttl` (FR19, FR20, FR21, FR24 as config support)  
   **When** `config.Load()` runs and either `~/.ligma/` or `~/.ligma/config.json` is missing  
   **Then** the package creates `~/.ligma/` and `~/.ligma/config.json` (empty `{}` or minimal valid JSON) before reading; subsequent Load returns the parsed config  
   **And** the CLI only reads `~/.ligma/config.json` and writes license files to user-specified paths; it does not collect, store, or transmit secrets or PII (NFR-S1)

## Tasks / Subtasks

- [x] **Task 1: Create internal/config package with Viper** (AC: #1)
  - [x] Create `internal/config/config.go`. Use Viper (`github.com/spf13/viper`) to read `~/.ligma/config.json`. Define a config struct or accessors for: `favorite` (string|nil), `aliases` (map or slice as in schema), `spdx_list_url`, `spdx_get_url_template`, `cache_ttl` (int or nil). Schema per PRD/cli-tool-specific-requirements.
  - [x] Add Viper to `go.mod` (`go get github.com/spf13/viper` or equivalent). Cobra can already use Viper; for this story, `internal/config` is the owner of the config file.
- [x] **Task 2: Create-if-absent for ~/.ligma/ and config.json** (AC: #1)
  - [x] In `config.Load()` (or first-use init): if `~/.ligma/` does not exist, create it (`os.MkdirAll`, 0755). If `~/.ligma/config.json` does not exist, create it with minimal valid JSON: `{}` or the default structure `{"favorite":null,"aliases":{},"spdx_list_url":"...","spdx_get_url_template":"...","cache_ttl":null}`. PRD default URLs can be used as defaults when creating; Viper can also apply defaults when keys are absent.
  - [x] After create-if-absent, read the file with Viper and return the parsed config. If the file is new and empty `{}`, Viper returns empty values; treat `favorite` as nil, `aliases` as empty, `cache_ttl` as nil, and `spdx_list_url` / `spdx_get_url_template` as absent (to be overridden by 5.2 when using defaults).
- [x] **Task 3: Resolve ~ from $HOME** (AC: #1)
  - [x] Use `os.UserHomeDir()` to resolve `~`. Path: `filepath.Join(home, ".ligma", "config.json")`. Do not assume `$HOME` env on all platforms; `UserHomeDir` is preferred.
- [x] **Task 4: NFR-S1 and tests** (AC: #1)
  - [x] Document or ensure: config only holds `favorite`, `aliases`, `spdx_list_url`, `spdx_get_url_template`, `cache_ttl`. No secrets, no PII. License files are written only to user-specified or default `LICENSE` path. Add `internal/config/config_test.go` for: create-if-absent in a temp dir, Load returns parsable config, missing keys yield sensible defaults.

## Dev Notes

- **Growth:** This story and Epic 5 are [Growth]. Commands that need config will call `config.Load()` at start; 5.2 will wire `ls`/`get` to use URLs from config; 5.3 will add cache. For 5.1, `Load()` need not be invoked by `ls`/`get`/`write` yet; the story's scope is the package and create-if-absent. If you prefer to integrate one command (e.g. `ls`) to trigger `Load` and thus create-if-absent on first use, that is acceptable.
- **Default config:** When creating an absent `config.json`, use empty `{}` or the PRD default shape. viper's `SetDefault` can supply `spdx_list_url` and `spdx_get_url_template` when keys are missing; for 5.2, "absent or empty" will fall back to hardcoded. 5.1 only needs create-if-absent and a parsable struct.
- **Cobra —viper:** project-context says add `--viper` to `cobra-cli init` when adding config. That is optional for 5.1 if config is loaded explicitly in `internal/config`; `--viper` mainly ties Cobra flags to Viper. For 5.1, `internal/config` can use Viper independently.

### Project Structure Notes

- **This story:** Create `internal/config/config.go` and `config_test.go`. No `cmd/` changes required for 5.1 if we only implement the package; or minimal change to have one command call `config.Load()` to exercise create-if-absent.

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-5-configuration-local-cache-growth.md#story-51-ligma-and-configjson-with-viper]
- [Source: bmad_docs/planning-artifacts/prd/cli-tool-specific-requirements.md] — Config Schema, default config, when it's created
- [Source: bmad_docs/planning-artifacts/architecture/core-architectural-decisions.md] — Config: Viper, `~/.ligma/config.json`
- [Source: bmad_docs/project-context.md] — Config [Growth]; create-if-absent; schema as PRD

---

## Developer Context

### Technical Requirements

- **Viper:** `github.com/spf13/viper`. Set config file path to `~/.ligma/config.json`; `ReadInConfig` or `ReadConfig`. On first run with missing file, create dir and file first, then read. New file can be `{}`; Viper will return empty for missing keys.
- **Schema (PRD):** `favorite` (string|null), `aliases` (object), `spdx_list_url`, `spdx_get_url_template`, `cache_ttl` (number|null). Default URLs in PRD; use when creating or as Viper defaults.

### Architecture Compliance

- **project-structure-boundaries:** "Config [Growth]: `internal/config` loads and parses `~/.ligma/config.json`."
- **core-architectural-decisions:** "Config (Growth): Viper, `~/.ligma/config.json`, schema as PRD."

### Library / Framework Requirements

- **Viper:** Add `github.com/spf13/viper` to `go.mod`. No other new deps for 5.1.

### File Structure Requirements

| Path                          | Purpose                                      | This story   |
|-------------------------------|----------------------------------------------|-------------|
| `internal/config/config.go`   | Viper, Load(), create-if-absent, config struct | **Create**  |
| `internal/config/config_test.go` | Tests for create-if-absent, Load, defaults   | **Create**  |

### Testing Requirements

- Use a temp directory as `~/.ligma` (or override path in tests) to avoid touching real home. Test: first Load creates dir and file; second Load reads; empty/missing keys yield expected defaults.

### Previous Story Intelligence

- **2.1–4.1:** `internal/spdx` and `cmd/ls`, `get`, `write` use hardcoded URLs. 5.1 does not change them; 5.2 will switch to config URLs when present.

### Project Context Reference

- **bmad_docs/project-context.md:** Config [Growth]: `~/.ligma/config.json`; create-if-absent at startup (FR18); schema as PRD.

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- internal/config: Viper, Config struct (Favorite, Aliases, SPDXListURL, SPDXGetURLTemplate, CacheTTL). Load() creates ~/.ligma and config.json with {} if absent; UserHomeDir for ~; SetDefault for the two URLs. SetConfigDirOverride for tests. config_test: CreateIfAbsent (dir+file, second Load reads), EmptyFile_Defaults, WithContent (favorite, aliases, cache_ttl). go get github.com/spf13/viper. NFR-S1: Config doc states only those five keys; no secrets/PII.

### File List

- internal/config/config.go (created)
- internal/config/config_test.go (created)
- go.mod (modified; viper)
