# Functional Requirements

*(MVP unless marked [Growth].)*

## License Listing

- **FR1:** User can list all available SPDX licenses in the terminal.
- **FR2:** User can filter the license list by a case-insensitive search term. [Growth]
- **FR3:** User can restrict the list to a static "popular" set of licenses. [Growth]
- **FR4:** User can obtain the license list in JSON format for scripting. [Growth]

## License Viewing

- **FR5:** User can output the full text of a specific license by SPDX ID to stdout.
- **FR6:** User can resolve a user-defined alias to an SPDX ID when viewing a license. [Growth]
- **FR7:** User can obtain license content in JSON format. [Growth]

## License Writing

- **FR8:** User can write the full text of a specific license to a file by SPDX ID.
- **FR9:** User can specify the target path for the written license file; if omitted, the file is written as `LICENSE` in the current directory.
- **FR10:** User can resolve a user-defined alias to an SPDX ID when writing. [Growth]
- **FR11:** User can write their configured favorite license by invoking `write` with no arguments. [Growth]
- **FR12:** User receives an error when invoking `write` with no arguments and no favorite is configured. [Growth]

## Help & Error Handling

- **FR13:** User can display help and usage information.
- **FR14:** User receives clear, actionable error messages for invalid or unknown license IDs.
- **FR15:** User receives a clear error when SPDX data cannot be fetched (e.g. network failure).
- **FR16:** CLI exits with zero on success and non-zero on error.
- **FR17:** User can run all commands in a non-interactive, scriptable manner (no prompts).

## Configuration

- **FR18:** User has a `~/.ligma/` directory and `~/.ligma/config.json` created when absent, at the start of each run. [Growth]
- **FR19:** User can set a favorite license in config for `write` with no arguments. [Growth]
- **FR20:** User can define aliases that map custom names to SPDX IDs. [Growth]
- **FR21:** User can override SPDX list and per-license URLs via config. [Growth]

## Local Cache

- **FR22:** User can have `ls` and `get` results served from a local cache under `~/.ligma/` when the cache is valid (within `cache_ttl`). [Growth]
- **FR23:** User can bypass the cache for `ls` and `get` by setting `cache_ttl` to `0` in config. [Growth]
- **FR24:** User can set cache TTL via the `cache_ttl` config property (how it is enforced TBD in implementation, e.g. stamp files). [Growth]

## SPDX Data Access

- **FR25:** CLI obtains the license list and individual license details from the official SPDX source (hardcoded URLs in MVP).
- **FR26:** CLI reads SPDX source URLs from the user's config when config is used. [Growth]
