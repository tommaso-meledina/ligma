# Epic 8: Rebrand to ligma

Replace the name `licensegen` with `ligma` (License Generation Minimal Assistant) across the codebase, config, and all documentation. Rename files where the old name appears in the filename.
**FRs covered:** N/A (branding only)

## Story 8.1: Rebrand Go module, code, and build

As a **developer**,
I want **the Go module, all source code, and build-related files to use the name `ligma`**,
So that **the binary, config paths, and programmatic identifiers reflect the new brand**.

**Acceptance Criteria:**

**Given** the current codebase under `github.com/tom/licensegen`
**When** the rebrand is applied to Go and build artifacts
**Then** `go.mod` declares `module github.com/tom/ligma`
**And** every import path `github.com/tom/licensegen/...` is updated to `github.com/tom/ligma/...` in: `main.go`, `cmd/*.go`, `cmd/*_test.go`, `internal/config/config.go`, `internal/config/config_test.go`, `internal/cache/cache.go`, `internal/cache/cache_test.go`
**And** `cmd/root.go` uses `Use: "ligma"` (replacing `"licensegen"`)
**And** `internal/config`: the function `LicensegenDir` is renamed to `LigmaDir`; all references in `internal/config/config_test.go`, `cmd/get.go`, and `cmd/ls.go` are updated; the tests `TestLicensegenDir_Override` and `TestLicensegenDir_UserHomeDir` are renamed to `TestLigmaDir_Override` and `TestLigmaDir_UserHomeDir`
**And** all user-facing or comment references to `~/.licensegen` are changed to `~/.ligma` in: `internal/config/config.go`, `internal/config/config_test.go`, `internal/cache/cache.go`, `cmd/write.go` (error message: `~/.ligma/config.json`)
**And** `.gitignore` lists the binary as `ligma` instead of `licensegen`
**And** `go build -o ligma` succeeds and tests pass

**Files to edit (15):** `go.mod`, `main.go`, `cmd/root.go`, `cmd/root_test.go`, `cmd/ls.go`, `cmd/ls_test.go`, `cmd/get.go`, `cmd/get_test.go`, `cmd/write.go`, `cmd/write_test.go`, `internal/config/config.go`, `internal/config/config_test.go`, `internal/cache/cache.go`, `internal/cache/cache_test.go`, `.gitignore`

---

## Story 8.2: Rebrand config and project docs

As a **developer**,
I want **BMad config and root-level project documentation to use the name `ligma`**,
So that **tooling and project brief align with the new brand**.

**Acceptance Criteria:**

**Given** `_bmad/bmm/config.yaml`, `BRIEF.md`, and `bmad_docs/project-context.md`
**When** the rebrand is applied
**Then** `_bmad/bmm/config.yaml` has `project_name: ligma`
**And** `BRIEF.md` refers to the CLI as `ligma` and uses `ligma ls`, `ligma get`, `ligma write` in examples
**And** `bmad_docs/project-context.md` uses `ligma` and `~/.ligma` where the tool or config paths are described

**Files to edit (3):** `_bmad/bmm/config.yaml`, `BRIEF.md`, `bmad_docs/project-context.md`

---

## Story 8.3: Rebrand planning and implementation docs; rename artifact file

As a **developer**,
I want **all planning-artifacts and implementation-artifacts markdown to use `ligma` and `~/.ligma`**,
So that **documentation is consistent and no filename contains the old brand**.

**Acceptance Criteria:**

**Given** all `.md` files under `bmad_docs/planning-artifacts/` and `bmad_docs/implementation-artifacts/`
**When** the rebrand is applied
**Then** every occurrence of `licensegen` (CLI name in examples, e.g. `licensegen ls`, `licensegen get`, `licensegen write`) is replaced with `ligma`
**And** every occurrence of `~/.licensegen` or `.licensegen` (config/cache paths) is replaced with `~/.ligma` or `.ligma`
**And** the file `bmad_docs/implementation-artifacts/5-1-licensegen-dir-config-json-viper.md` is **renamed** to `bmad_docs/implementation-artifacts/5-1-ligma-dir-config-json-viper.md`; the title and content inside are updated to `~/.ligma` and `ligma` as above
**And** any references to “licensegen” or “LicensegenDir” in these docs (e.g. in 5.1, 5.3, architecture) are updated to `ligma` / `LigmaDir` where applicable

**Files to edit (38):**  
- **planning-artifacts:** `implementation-readiness-report-2026-01-20.md`, `project-context.md`, `epics/index.md`, `epics/overview.md`, `epics/requirements-inventory.md`, `epics/epic-1-cli-foundation-help.md` through `epic-7-view-write-enhancements-growth.md`, `prd/index.md`, `prd/product-scope.md`, `prd/user-journeys.md`, `prd/functional-requirements.md`, `prd/non-functional-requirements.md`, `prd/project-scoping-phased-development.md`, `prd/cli-tool-specific-requirements.md`, `prd/success-criteria.md`, `architecture/index.md`, `architecture/project-context-analysis.md`, `architecture/core-architectural-decisions.md`, `architecture/project-structure-boundaries.md`, `architecture/starter-template-evaluation.md`  
- **implementation-artifacts:** `1-1-go-module-project-structure-dummy-main.md`, `1-2-add-cobra-ls-get-write-subcommands.md`, `1-3-refine-exit-codes.md`, `2-2-ls-command-list-all-spdx-licenses.md`, `3-2-get-command-output-license-text-by-spdx-id.md`, `4-1-write-command-write-license-to-file-by-spdx-id.md`, `5-1-licensegen-dir-config-json-viper.md` (edit then **rename** to `5-1-ligma-dir-config-json-viper.md`), `5-2-spdx-client-uses-config-urls.md`, `5-3-cache-layer-ls-get.md`, `6-1-ls-json-output.md`, `6-2-ls-filter-popular-set.md`, `7-1-alias-resolution-get-write.md`, `7-2-get-json-output.md`, `7-3-write-no-args-favorite-or-error.md`

**File to rename (1):**  
`5-1-licensegen-dir-config-json-viper.md` → `5-1-ligma-dir-config-json-viper.md`

---

