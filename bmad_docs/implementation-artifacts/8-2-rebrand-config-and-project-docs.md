# Story 8.2: Rebrand config and project docs

Status: review

## Story

As a **developer**,
I want **BMad config and root-level project documentation to use the name `ligma`**,
So that **tooling and project brief align with the new brand**.

## Acceptance Criteria

**Given** `_bmad/bmm/config.yaml`, `BRIEF.md`, and `bmad_docs/project-context.md`
**When** the rebrand is applied
**Then** `_bmad/bmm/config.yaml` has `project_name: ligma`
**And** `BRIEF.md` refers to the CLI as `ligma` and uses `ligma ls`, `ligma get`, `ligma write` in examples
**And** `bmad_docs/project-context.md` uses `ligma` and `~/.ligma` where the tool or config paths are described

## Tasks / Subtasks

- [x] **Task 1: _bmad/bmm/config.yaml** (AC: project_name)
  - [x] `project_name: licensegen` → `project_name: ligma`
- [x] **Task 2: BRIEF.md** (AC: CLI name and examples)
  - [x] CLI named `licensegen` → `ligma`
  - [x] Examples: `licensegen ls`, `licensegen get <license id>`, `licensegen write <license id> [location]` → `ligma ls`, `ligma get`, `ligma write`
- [x] **Task 3: bmad_docs/project-context.md** (AC: ligma and ~/.ligma)
  - [x] `project_name` frontmatter, intro "licensegen" → "ligma"
  - [x] All `~/.licensegen` and `.licensegen` paths → `~/.ligma` / `.ligma`
  - [x] Build example `go build -o licensegen` → `go build -o ligma`

## Dev Notes

- **Source:** epic-8-rebrand-to-ligma.md Story 8.2
- **Files to edit (3):** `_bmad/bmm/config.yaml`, `BRIEF.md`, `bmad_docs/project-context.md`

## Dev Agent Record

### File List

- `_bmad/bmm/config.yaml`
- `BRIEF.md`
- `bmad_docs/project-context.md`

## Change Log

- 2026-01-20: Story 8.2 implemented — rebrand config and project docs: config.yaml `project_name: ligma`; BRIEF.md CLI `ligma` and examples `ligma ls`/`get`/`write`; project-context.md `ligma`, `~/.ligma`, `go build -o ligma`.
