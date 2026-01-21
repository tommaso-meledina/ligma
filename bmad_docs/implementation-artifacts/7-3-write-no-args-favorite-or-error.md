# Story 7.3: write with no arguments — favorite or error [Growth]

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **user**,
I want **to run `ligma write` with no arguments and have it write my configured favorite license**,
so that **I can quickly add my usual license (FR11, FR12)**.

## Acceptance Criteria

1. **Given** config with `favorite` set to an SPDX ID (or an alias that resolves to one)  
   **When** I run `ligma write` with no arguments  
   **Then** the CLI writes the favorite license to `LICENSE` in the current directory and exits with 0 (FR11)  
   **Given** config with no `favorite` set (or favorite is empty)  
   **When** I run `ligma write` with no arguments  
   **Then** the CLI prints a clear error to stderr and exits with 1 (or an implementation-defined code for "config error") (FR12)  
   **And** `ligma write` with no arguments always writes to `LICENSE` in the current directory; a custom path for favorite-only is out of scope for this story

## Tasks / Subtasks

- [x] **Task 1: Allow zero args for write** (AC: #1)
  - [x] Today `write` has `RangeArgs(0,2)`; when args empty use favorite, else args[0] (and optional args[1]). (ID required). For 7.3, when args are **empty**, treat as “write favorite to LICENSE.” When args are non-empty, keep current behavior: first arg is ID (or alias), optional second is path. So: `cobra.MinimumNArgs(0)` and `MaximumNArgs(2)`. In `RunE`: if `len(args) == 0`, use favorite path; if `len(args) >= 1`, use `args[0]` as ID and `args[1]` as optional path.
- [x] **Task 2: Resolve favorite from config** (AC: #1)
  - [x] Call `config.Load()`. Read `favorite`. If `favorite` is non-empty (and not null), use it as the SPDX ID or alias. Resolve via `cfg.Resolve(*cfg.Favorite)`. Then fetch and write to `LICENSE` in the current directory.
  - [x] If `favorite` is absent, null, or empty string: return error that root maps to **exit 1** (usage or config error). The AC allows “1 (or an implementation-defined code for config error)”; use **1** to stay within 0/1/2/3. Do not introduce a new exit code without updating project-context.
- [x] **Task 3: Path for favorite-only** (AC: #1)
  - [x] When `ligma write` is called with no args, **always** write to `LICENSE` in the current directory. No optional path for this mode; “custom path for favorite-only is out of scope.”
- [x] **Task 4: Preserve write <id> [path] behavior** (AC: #1)
  - [x] `ligma write <id>` and `ligma write <id> <path>` behave as before (4.1, 7.1). Only the zero-arg case is new.
- [x] **Task 5: Not-found and I/O for favorite** (AC: #1)
  - [x] If the resolved favorite ID is not found (404), or SPDX/cache fails (I/O): same as 4.1 — exit 2 for not-found, exit 3 for I/O. File write failure to `LICENSE` → 3.
- [x] **Task 6: Verify and help** (AC: #1)
  - [x] With `favorite: "MIT"` (or `"mit"` and alias), `ligma write` → writes `LICENSE` in cwd, 0. With `favorite: null` or missing, `ligma write` → error on stderr, exit 1. Long text: "With no arguments, writes the configured favorite to LICENSE in the current directory."

## Dev Notes

- **Pre-requisite:** 5.1 (config, `favorite`), 7.1 (alias resolution) so that `favorite` can be an alias. 4.1 (write behavior). Config is loaded in `write` (5.2); 7.3 adds a branch when args are empty.
- **favorite in schema:** PRD: `favorite` (string | null). Default null. Empty string can be treated like null (no favorite).

### Project Structure Notes

- **This story:** Modify `cmd/write.go` only. `internal/config` already has `favorite`; ensure an accessor exists if needed.

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-7-view-write-enhancements-growth.md#story-73-write-with-no-arguments-favorite-or-error]
- [Source: bmad_docs/planning-artifacts/prd/cli-tool-specific-requirements.md] — `write` with no args uses `favorite`; error if unset; Config: `favorite`
- [Source: bmad_docs/project-context.md] — Config: `favorite`; FR11, FR12; do not add new exit codes without updating project-context

---

## Developer Context

### Technical Requirements

- **Args:** `cobra.MinimumNArgs(0)`, `MaximumNArgs(2)`. In `RunE`: `if len(args)==0 { use favorite, path=LICENSE } else { id=args[0], path=args[1] or LICENSE }`.
- **Exit 1 for no favorite:** Return `ErrUsage` or an error root maps to 1. Message e.g. “favorite license is not set; run with <id> or set favorite in ~/.ligma/config.json”.

### Architecture Compliance

- **project-structure-boundaries:** “FR8–FR12: … [Growth] `internal/config` (favorite, aliases).”
- **project-context:** Exit 1 for usage/config; do not introduce extra codes.

### File Structure Requirements

| Path           | Purpose                                      | This story   |
|----------------|----------------------------------------------|-------------|
| `cmd/write.go` | Allow 0 args; branch on args; favorite→LICENSE; no favorite→exit 1 | **Modify**  |

### Testing Requirements

- Manual: `favorite` set → `ligma write` creates `LICENSE`, 0. `favorite` null/empty → `ligma write` → stderr, 1. `ligma write MIT` and `ligma write MIT ./x` unchanged.

### Previous Story Intelligence

- **4.1:** `write` requires at least ID. 7.3 adds 0-arg mode. **7.1:** Alias resolution; `favorite` can be an alias, so resolve it the same way. **5.1:** `favorite` in config.

### Project Context Reference

- **bmad_docs/project-context.md:** Config `favorite`; “`write` with no args” uses favorite, errors if unset; exit codes 0/1/2/3 only.

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- write: `Args: cobra.RangeArgs(0, 2)`. If `len(args)==0`: cfg.Favorite nil or empty → `fmt.Errorf("favorite license is not set; run with <id> or set favorite in ~/.ligma/config.json")` (exit 1); else `id=cfg.Resolve(*cfg.Favorite)`, path=`cwd/LICENSE`. Else: id=Resolve(args[0]), path=args[1] or cwd/LICENSE. Fetch and write unchanged. Not-found→2, I/O→3. Long: "With no arguments, writes the configured favorite to LICENSE in the current directory."
- Tests: TestWriteRunE_ZeroArgsFavorite, TestWriteRunE_ZeroArgsNoFavorite, TestWriteRunE_ZeroArgsFavoriteResolvedViaAlias; TestExecute_WriteNoArgs uses config override (no favorite) → exit 1.

### File List

- cmd/write.go (modified)
- cmd/write_test.go (modified)
- cmd/root_test.go (modified: TestExecute_WriteNoArgs sets config override)
