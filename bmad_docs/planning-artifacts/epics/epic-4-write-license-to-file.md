# Epic 4: Write License to File

User can write a license to a file by SPDX ID, with configurable path or default `LICENSE`.
**FRs covered:** FR8, FR9, FR14

## Story 4.1: write command — write license to file by SPDX ID

As a **user**,
I want **to run `ligma write <id>` or `ligma write <id> <path>` and have the license written to a file**,
So that **I can add a license to my project (FR8, FR9)**.

**Acceptance Criteria:**

**Given** the SPDX client that can fetch license details by ID
**When** I run `ligma write <SPDX-ID>` and the SPDX source is reachable
**Then** the CLI writes the full license text to `LICENSE` in the current directory and exits with 0
**When** I run `ligma write <SPDX-ID> <path>`
**Then** the CLI writes the full license text to `<path>` and exits with 0; if the file exists, it is overwritten
**When** the ID is unknown or invalid, the CLI prints a clear, actionable error to stderr and exits with 2 (FR14)
**When** the SPDX source is unreachable or file write fails (e.g. permission, disk), the CLI prints a clear error to stderr and exits with 3
**And** `write` requires at least one argument (the SPDX ID); otherwise the CLI prints usage to stderr and exits with 1

---
