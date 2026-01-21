# Story 7.1: Alias resolution for get and write [Growth]

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **user**,
I want **to use aliases from config in `get` and `write`**,
so that **I can type short names like `mit` instead of `MIT` (FR6, FR10)**.

## Acceptance Criteria

1. **Given** config with `aliases` mapping names to SPDX IDs (FR20)  
   **When** I run `ligma get <alias>` or `ligma write <alias>` (or `write <alias> <path>`), the CLI resolves the alias to the SPDX ID before fetching  
   **Then** behavior is the same as if the SPDX ID had been passed; if the alias is unknown, the CLI treats it as an unknown ID (clear error, exit 2)  
   **And** when both an alias and an SPDX ID could match, resolution rules are well-defined (e.g. aliases take precedence, or exact SPDX ID; TBD in implementation)

## Tasks / Subtasks

- [x] **Task 1: Config aliases and resolver** (AC: #1)
  - [x] `internal/config` (from 5.1) already has `aliases` (object/map: string → SPDX ID). Expose an accessor (e.g. `config.Aliases() map[string]string`) or a `Resolve(idOrAlias string) (spdxID string, ok bool)`. If the input is a key in `aliases`, return the mapped SPDX ID. If not, return the input as-is (treat as SPDX ID) and `ok` true, or a distinct “unresolved” for “use as SPDX ID.” Define clearly for 7.1.
  - [x] **Resolution rule (TBD in implementation):** Recommended: **aliases take precedence**. If the argument equals an alias key (exact match, case-sensitive for the alias key as in the JSON), use the mapped SPDX ID. Otherwise, treat the argument as an SPDX ID and pass through. So `get mit` with `"mit" -> "MIT"` → fetch MIT; `get MIT` with no alias `MIT` → fetch MIT as SPDX ID. If an alias shadows an SPDX ID (e.g. `"MIT"` → `"Apache-2.0"`), the alias wins. Document the rule in code or config.
- [x] **Task 2: get: resolve before fetch** (AC: #1)
  - [x] In `cmd/get`, after `config.Load()` and before calling SPDX/cache, resolve the first argument: `resolved := config.Resolve(args[0])` (or lookup in `config.Aliases()`). Use `resolved` as the SPDX ID for the fetch. If the alias is unknown, `Resolve` returns the arg as-is (so we still call SPDX); SPDX 404 or “not in list” will then produce exit 2. The AC says “if the alias is unknown, treat as unknown ID” → same as passing an unknown SPDX ID: clear error, exit 2. So no special case for “alias not in map”; we pass through and let SPDX/cache determine not-found.
- [x] **Task 3: write: resolve before fetch** (AC: #1)
  - [x] In `cmd/write`, same as `get`: resolve the first arg (the ID or alias) via config before calling SPDX. Use the resolved SPDX ID for the fetch. Path (second arg) is unchanged. Unknown alias → same as unknown ID → exit 2.
- [x] **Task 4: Aliases config shape** (AC: #1)
  - [x] PRD: `aliases` is an object, e.g. `{"apache":"Apache-2.0","mit":"MIT"}`. Config must parse this. If `aliases` is missing or null, treat as empty map (no aliases). If it’s `[]`, also empty.
- [x] **Task 5: Verify and document** (AC: #1)
  - [x] With `aliases: {"mit":"MIT"}`, `ligma get mit` and `ligma write mit` behave like `get MIT` and `write MIT`. `ligma get unknownalias` with no such alias → passes through to SPDX → 404/not-found → exit 2. Document resolution rule (aliases take precedence) in code or user-facing docs.

## Dev Notes

- **Pre-requisite:** 5.1 (config, `aliases`). 5.2 and 5.3 are in use for `get`; 4.1 for `write`. Both `get` and `write` already call `config.Load()` (5.2). 7.1 adds a resolve step after Load and before fetch.
- **get/write only:** `ls` does not take an ID or alias; 7.1 does not change `ls`.

### Project Structure Notes

- **This story:** Modify `internal/config` (add `Resolve` or alias lookup) and `cmd/get.go`, `cmd/write.go`. No new packages.

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-7-view-write-enhancements-growth.md#story-71-alias-resolution-for-get-and-write]
- [Source: bmad_docs/planning-artifacts/prd/cli-tool-specific-requirements.md] — Config: `aliases` (object)
- [Source: bmad_docs/project-context.md] — Config [Growth]; aliases; resolve before SPDX lookup

---

## Developer Context

### Technical Requirements

- **Alias lookup:** Exact key match in the `aliases` map. If `arg` is a key, return `aliases[arg]`. Else return `arg` (treat as SPDX ID). Case-sensitive for the key; SPDX IDs are case-sensitive.
- **Unknown alias:** No special handling. “Unknown alias” means it’s not in the map, so we use the arg as SPDX ID. If that ID doesn’t exist in SPDX, we get 404 → exit 2. So “unknown alias” and “unknown ID” both end as exit 2.

### Architecture Compliance

- **project-structure-boundaries:** “FR5–FR7: … [Growth] `internal/config` (aliases).” “FR8–FR12: … [Growth] `internal/config` (favorite, aliases).”

### File Structure Requirements

| Path                 | Purpose                                      | This story   |
|----------------------|----------------------------------------------|-------------|
| `internal/config/config.go` | Add `Resolve(idOrAlias string) string` or alias lookup | **Modify**  |
| `cmd/get.go`         | Resolve first arg before fetch               | **Modify**  |
| `cmd/write.go`       | Resolve first arg (ID) before fetch          | **Modify**  |

### Testing Requirements

- Config with `aliases: {"x":"MIT"}`: `get x` and `write x` use MIT. `get notanalias` passes through; if `notanalias` is not in SPDX, exit 2. Unit test for `Resolve`: alias key → mapped ID; non-key → pass-through.

### Previous Story Intelligence

- **5.1, 5.2:** Config has `aliases`. 5.2 passes URLs to SPDX; 7.1 adds a resolve step for the ID/alias argument. Order: Load config → Resolve ID/alias → Fetch (with resolved ID and URLs from config).

### Project Context Reference

- **bmad_docs/project-context.md:** “ID and alias resolution: … resolve … `aliases` before SPDX lookup; unknown ID/alias → clear error and non-zero exit.”

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- Config: `(c *Config) Resolve(idOrAlias string) string` — if `c.Aliases[idOrAlias]` exists return it, else return idOrAlias. Aliases take precedence. Nil Aliases: return idOrAlias. Load already normalizes nil aliases to empty map; GetStringMapString for missing/null/`[]` yields nil, then we make().
- get: `id := cfg.Resolve(args[0])` before cache.FetchDetails. write: `id := cfg.Resolve(args[0])` before writeFetchDetails. Unknown alias passes through as SPDX ID → 404 → exit 2.
- Tests: config TestResolve_AliasKeyReturnsMappedID, TestResolve_NonKeyReturnsAsIs, TestResolve_NilAliasesReturnsAsIs; get TestGetRunE_AliasResolved (aliases mit->MIT, get mit), TestGetRunE_UnknownAliasPassesThrough; write TestWriteRunE_AliasResolved.

### File List

- internal/config/config.go (modified)
- internal/config/config_test.go (modified)
- cmd/get.go (modified)
- cmd/get_test.go (modified)
- cmd/write.go (modified)
- cmd/write_test.go (modified)
