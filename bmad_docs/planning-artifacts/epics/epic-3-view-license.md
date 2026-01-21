# Epic 3: View License

User can output the full text of a license by SPDX ID to stdout.
**FRs covered:** FR5, FR14

## Story 3.1: SPDX client — fetch license details by ID

As a **developer**,
I want **`internal/spdx` to fetch `details/{id}.json` for a given SPDX ID and return the license text**,
So that **the `get` command can output the full license (FR25)**.

**Acceptance Criteria:**

**Given** an SPDX license ID (e.g. `MIT`, `Apache-2.0`) and the official SPDX `details/{id}.json` URL template (hardcoded for MVP)
**When** the client fetches the details with a 30s timeout (NFR-I2)
**Then** it returns the license text (e.g. from the `licenseText` field) or an error if the request fails, times out, returns 404, or the response is not valid JSON
**And** it uses the same `net/http` and `encoding/json` approach as the list client; no third-party HTTP client
**And** SPDX IDs are used as-is (e.g. `Apache-2.0`); no normalization for the URL path

## Story 3.2: get command — output license text by SPDX ID

As a **user**,
I want **to run `ligma get <id>` and see the full license text in the terminal**,
So that **I can read or pipe a license before writing it (FR5)**.

**Acceptance Criteria:**

**Given** the SPDX client that can fetch license details by ID
**When** I run `ligma get <SPDX-ID>` with a valid ID and the SPDX source is reachable
**Then** the CLI prints the full license text to stdout and exits with 0
**When** the ID is unknown or invalid (e.g. 404 from SPDX, or ID not in the list), the CLI prints a clear, actionable error to stderr and exits with 2 (FR14)
**When** the SPDX source is unreachable (network failure, 4xx/5xx, timeout), the CLI prints a clear error to stderr and exits with 3
**And** `get` requires exactly one argument (the SPDX ID); otherwise the CLI prints usage to stderr and exits with 1

---
