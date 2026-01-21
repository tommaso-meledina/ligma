# Epic 7: View & Write Enhancements [Growth]

User can use aliases and favorite for get/write, JSON for get, and write with no args (favorite or error).
**FRs covered:** FR6, FR7, FR10, FR11, FR12

## Story 7.1: Alias resolution for get and write

As a **user**,
I want **to use aliases from config in `get` and `write`**,
So that **I can type short names like `mit` instead of `MIT` (FR6, FR10)**.

**Acceptance Criteria:**

**Given** config with `aliases` mapping names to SPDX IDs (FR20)
**When** I run `ligma get <alias>` or `ligma write <alias>` (or `write <alias> <path>`), the CLI resolves the alias to the SPDX ID before fetching
**Then** behavior is the same as if the SPDX ID had been passed; if the alias is unknown, the CLI treats it as an unknown ID (clear error, exit 2)
**And** when both an alias and an SPDX ID could match, resolution rules are well-defined (e.g. aliases take precedence, or exact SPDX ID; TBD in implementation)

## Story 7.2: get — JSON output

As a **user**,
I want **to run `ligma get <id> --json` and receive the license in JSON**,
So that **I can script and pipe structured data (FR7)**.

**Acceptance Criteria:**

**Given** the `get` command
**When** I run `ligma get <SPDX-ID> --json`
**Then** the CLI prints the license (e.g. `licenseText` and optionally ID) as JSON to stdout and exits with 0
**And** without `--json`, `get` keeps the existing plain-text output

## Story 7.3: write with no arguments — favorite or error

As a **user**,
I want **to run `ligma write` with no arguments and have it write my configured favorite license**,
So that **I can quickly add my usual license (FR11, FR12)**.

**Acceptance Criteria:**

**Given** config with `favorite` set to an SPDX ID (or an alias that resolves to one)
**When** I run `ligma write` with no arguments
**Then** the CLI writes the favorite license to `LICENSE` in the current directory and exits with 0 (FR11)
**Given** config with no `favorite` set (or favorite is empty)
**When** I run `ligma write` with no arguments
**Then** the CLI prints a clear error to stderr and exits with 1 (or an implementation-defined code for "config error") (FR12)
**And** `ligma write` with no arguments always writes to `LICENSE` in the current directory; a custom path for favorite-only is out of scope for this story

