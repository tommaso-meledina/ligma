# Story 8.3: Rebrand planning and implementation docs; rename artifact file

Status: review

## Story

As a **developer**,
I want **all planning-artifacts and implementation-artifacts markdown to use `ligma` and `~/.ligma`**,
So that **documentation is consistent and no filename contains the old brand**.

## Acceptance Criteria

**Given** all `.md` files under `bmad_docs/planning-artifacts/` and `bmad_docs/implementation-artifacts/`
**When** the rebrand is applied
**Then** every occurrence of `licensegen` (CLI name in examples, e.g. `licensegen ls`, `licensegen get`, `licensegen write`) is replaced with `ligma`
**And** every occurrence of `~/.licensegen` or `.licensegen` (config/cache paths) is replaced with `~/.ligma` or `.ligma`
**And** the file `5-1-licensegen-dir-config-json-viper.md` is **renamed** to `5-1-ligma-dir-config-json-viper.md`; the title and content inside are updated to `~/.ligma` and `ligma` as above
**And** any references to "licensegen" or "LicensegenDir" in these docs are updated to `ligma` / `LigmaDir` where applicable

## Tasks / Subtasks

- [x] **Task 1: Replace ~/.licensegen and .licensegen in all planning and implementation artifacts**
- [x] **Task 2: Replace licensegen (CLI examples), LicensegenDir, go build -o licensegen, ./licensegen, product name**
- [x] **Task 3: Edit 5-1-licensegen-dir, then rename to 5-1-ligma-dir; update links (epic-5, epics/index, 5-1 refs)**
- [x] **Task 4: Exclude 8-1, 8-2 (historical AC); leave epic-8 rebrand-spec phrasing as-is**

## Dev Notes

- **Source:** epic-8-rebrand-to-ligma.md Story 8.3. Scope: `bmad_docs/planning-artifacts/`, `bmad_docs/implementation-artifacts/`.
- **Rename:** `5-1-licensegen-dir-config-json-viper.md` → `5-1-ligma-dir-config-json-viper.md`. Update epic-5 heading and `#story-51-licensegen-and-configjson-with-viper` → `#story-51-ligma-and-configjson-with-viper`.

## Dev Agent Record

### File List

- **Renamed:** `5-1-licensegen-dir-config-json-viper.md` → `5-1-ligma-dir-config-json-viper.md` (created 5-1-ligma, deleted 5-1-licensegen)
- **Planning-artifacts (24):** implementation-readiness-report, epics/index, overview, requirements-inventory, epic-1–7, prd/index, product-scope, user-journeys, functional, non-functional, project-scoping, cli-tool-specific, success-criteria, architecture/index, project-context-analysis, core-architectural-decisions, project-structure-boundaries, starter-template-evaluation
- **Implementation-artifacts (14):** 1-3, 2-2, 3-2, 4-1, 5-2, 5-3, 6-1, 6-2, 7-1, 7-2, 7-3; 5-1 (edit+rename). Excluded: 8-1, 8-2 (historical); epic-8 (rebrand spec)

## Change Log

- 2026-01-20: Story 8.3 implemented — rebrand planning/impl docs: ~/.licensegen→~/.ligma, licensegen→ligma (CLI examples), LicensegenDir→LigmaDir, go build -o ligma, product names; 5-1 renamed to 5-1-ligma-dir-config-json-viper.md; epic-5 heading and epics/index anchor #story-51-ligma-and-configjson-with-viper.
