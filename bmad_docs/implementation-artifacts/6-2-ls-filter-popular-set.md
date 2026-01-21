# Story 6.2: ls — filter and popular set [Growth]

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **user**,
I want **to run `ligma ls` with an optional filter or `--popular`**,
so that **I can narrow the list by search term or a static popular set (FR2, FR3)**.

## Acceptance Criteria

1. **Given** the `ls` command  
   **When** I run `ligma ls --filter <term>` (or equivalent flag/arg), the list is restricted to licenses whose SPDX ID (or other defined fields) match the term **case-insensitively** (FR2)  
   **When** I run `ligma ls --popular`, the list is restricted to a static, hardcoded "popular" set of SPDX IDs (e.g. MIT, Apache-2.0, GPL-2.0, etc.; exact set TBD in implementation) (FR3)  
   **And** `--filter` and `--popular` can be combined or used independently; without them, the full list is shown

## Tasks / Subtasks

- [x] **Task 1: Add --filter and --popular to ls** (AC: #1)
  - [x] In `cmd/ls.go`, add `--filter <term>` (string) and `--popular` (bool). When `--filter` is non-empty, filter the list before output. When `--popular` is true, restrict to the popular set. Both can be set: first restrict to popular (if `--popular`), then apply `--filter` on that subset; or apply both in an order that matches “combined” (e.g. popular ∩ filter, or filter on popular set). “Without them, the full list” → no filter and no popular means show all.
- [x] **Task 2: Implement --filter (case-insensitive)** (AC: #1)
  - [x] Filter licenses where SPDX ID (or other defined fields, e.g. `name`) contains the term **case-insensitively** (FR2). Use `strings.Contains(strings.ToLower(id), strings.ToLower(term))` or equivalent. If the list includes `name`, matching either `licenseId` or `name` is acceptable; define and document. Substring match is typical; exact match is stricter—prefer substring for “filter” semantics.
- [x] **Task 3: Implement --popular (static set)** (AC: #1)
  - [x] Define a hardcoded slice of SPDX IDs, e.g. `[]string{"MIT", "Apache-2.0", "GPL-2.0", "BSD-3-Clause", "ISC"}` or similar (exact set TBD). When `--popular` is true, restrict the list to those IDs that are in both the fetched list and the popular set. If a popular ID is not in the SPDX list, it can be omitted from the output.
- [x] **Task 4: Combine --filter and --popular** (AC: #1)
  - [x] If both are set: e.g. first apply `--popular` (restrict to popular set), then apply `--filter` on that subset. Or: apply `--filter` on the full list, then intersect with popular. Choose one and document. The AC says “can be combined”; either order is acceptable as long as it’s consistent.
- [x] **Task 5: Apply to both human and JSON output** (AC: #1)
  - [x] Filtering and popular apply **before** formatting. So both `ligma ls --filter mit` and `ligma ls --json --filter mit` return the filtered list; same for `--popular`. 6.1’s `--json` is unchanged; 6.2 only narrows the list earlier in the pipeline.
- [x] **Task 6: Verify and help** (AC: #1)
  - [x] `ligma ls --help` documents `--filter` and `--popular`. `ligma ls --filter apache` and `ligma ls --popular` and `ligma ls --popular --filter 2` behave as specified.

## Dev Notes

- **Pre-requisite:** `ls` has full list (from SPDX or cache) and supports `--json` (6.1). 6.2 adds filtering and popular **before** output formatting.
- **Popular set:** “Exact set TBD.” A reasonable default: MIT, Apache-2.0, GPL-2.0, BSD-3-Clause, ISC. Adjust if the product brief or PRD specifies. Keep it small (e.g. 5–10) for “popular.”

### Project Structure Notes

- **This story:** Modify `cmd/ls.go` only. Filtering logic can live in `ls` or a small helper in `cmd`; no new packages.

### References

- [Source: bmad_docs/planning-artifacts/epics/epic-6-listing-enhancements-growth.md#story-62-ls-filter-and-popular-set]
- [Source: bmad_docs/project-context.md] — `--popular` on `ls` and `get` [Growth]; actually `--popular` is for `ls` per PRD
- [Source: bmad_docs/planning-artifacts/prd/functional-requirements.md] — FR2 (filter), FR3 (popular)

---

## Developer Context

### Technical Requirements

- **Case-insensitive:** `strings.EqualFold` or `strings.Contains(strings.ToLower(a), strings.ToLower(b))`. Apply to `licenseId` and optionally `name`.
- **Order of operations:** e.g. (1) fetch full list, (2) if `--popular` then restrict to popular set, (3) if `--filter` then filter by term, (4) output (human or JSON).

### Architecture Compliance

- **project-context:** “Commands/flags: … `--popular` … (Growth).” 6.2 adds `--filter` and `--popular` to `ls`.

### File Structure Requirements

| Path           | Purpose                                      | This story   |
|----------------|----------------------------------------------|-------------|
| `cmd/ls.go`    | Add `--filter`, `--popular`; filter and popular logic before output | **Modify**  |

### Testing Requirements

- Manual: `--filter` case-insensitive; `--popular` returns subset; both combined. Optional: unit test for filter and popular logic with a fixed list.

### Previous Story Intelligence

- **6.1:** `--json` in `ls`. 6.2 filters before format; both human and `--json` see the filtered/popular list.

### Project Context Reference

- **bmad_docs/project-context.md:** `--popular` for `ls`; flags in the command that uses them.

### Story Completion Status

- **Status:** ready-for-dev  
- **Completion note:** Ultimate context engine analysis completed — comprehensive developer guide created.

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

- ls: --filter (string), --popular (bool). popularIDs = MIT, Apache-2.0, GPL-2.0, BSD-3-Clause, ISC. Order: fetch → --popular (intersect with set) → --filter (case-insensitive substring on licenseId or name) → output (human or --json). Filter: strings.Contains(ToLower(id/name), ToLower(term)). ls_test: Filter (apache→Apache-2.0), FilterCaseInsensitive (mit→MIT), Popular (5, no X), PopularAndFilter (--popular --filter 2 → Apache-2.0, GPL-2.0).

### File List

- cmd/ls.go (modified)
- cmd/ls_test.go (modified)
