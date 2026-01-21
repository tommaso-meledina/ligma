# Epic 2: List Licenses

User can list all SPDX licenses in the terminal from the official source.
**FRs covered:** FR1, FR15, FR25 | **NFRs:** NFR-P1, NFR-I1, NFR-I2

## Story 2.1: SPDX client — fetch and parse license list

As a **developer**,
I want **`internal/spdx` to fetch `licenses.json` from the official SPDX URL and parse it**,
So that **the `ls` command can display the license list (FR25)**.

**Acceptance Criteria:**

**Given** the official SPDX `licenses.json` URL (hardcoded for MVP)
**When** the client fetches the license list with a 30s timeout (NFR-I2)
**Then** it returns the parsed list (structs matching `license-list-data` layout) or an error if the request fails, times out, or the response is not valid JSON
**And** it uses `net/http` and `encoding/json` only; no third-party HTTP client (Architecture)
**And** the package lives under `internal/spdx` (e.g. `client.go`)

## Story 2.2: ls command — list all SPDX licenses

As a **user**,
I want **to run `ligma ls` and see all available SPDX licenses in the terminal**,
So that **I can choose which license to use (FR1, FR25)**.

**Acceptance Criteria:**

**Given** the SPDX client that can fetch the license list
**When** I run `ligma ls` and the SPDX source is reachable
**Then** the CLI prints the list of license identifiers (e.g. SPDX IDs, one per line or a clear, consistent format) to stdout and exits with 0
**When** the SPDX source is unreachable (network failure, 4xx/5xx, timeout), the CLI prints a clear, actionable error to stderr and exits with 3 (FR15, NFR-I1, NFR-I2)
**And** under typical network conditions, the command completes within 15 seconds (NFR-P1)

---
