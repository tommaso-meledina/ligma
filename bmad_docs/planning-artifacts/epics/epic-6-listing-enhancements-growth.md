# Epic 6: Listing Enhancements [Growth]

User can filter, restrict to a popular set, and get the license list in JSON for scripting.
**FRs covered:** FR2, FR3, FR4

## Story 6.1: ls — JSON output

As a **user**,
I want **to run `ligma ls --json` and get the license list in JSON**,
So that **I can script and integrate with other tools (FR4)**.

**Acceptance Criteria:**

**Given** the `ls` command
**When** I run `ligma ls --json`
**Then** the CLI prints the license list as JSON to stdout (e.g. an array of objects with at least SPDX ID; exact field names TBD in implementation) and exits with 0
**And** without `--json`, `ls` keeps the existing human-readable format

## Story 6.2: ls — filter and popular set

As a **user**,
I want **to run `ligma ls` with an optional filter or `--popular`**,
So that **I can narrow the list by search term or a static popular set (FR2, FR3)**.

**Acceptance Criteria:**

**Given** the `ls` command
**When** I run `ligma ls --filter <term>` (or equivalent flag/arg), the list is restricted to licenses whose SPDX ID (or other defined fields) match the term **case-insensitively** (FR2)
**When** I run `ligma ls --popular`, the list is restricted to a static, hardcoded "popular" set of SPDX IDs (e.g. MIT, Apache-2.0, GPL-2.0, etc.; exact set TBD in implementation) (FR3)
**And** `--filter` and `--popular` can be combined or used independently; without them, the full list is shown

---
