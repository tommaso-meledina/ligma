# Requirements Inventory

## Functional Requirements

FR1: User can list all available SPDX licenses in the terminal.
FR2: User can filter the license list by a case-insensitive search term. [Growth]
FR3: User can restrict the list to a static "popular" set of licenses. [Growth]
FR4: User can obtain the license list in JSON format for scripting. [Growth]
FR5: User can output the full text of a specific license by SPDX ID to stdout.
FR6: User can resolve a user-defined alias to an SPDX ID when viewing a license. [Growth]
FR7: User can obtain license content in JSON format. [Growth]
FR8: User can write the full text of a specific license to a file by SPDX ID.
FR9: User can specify the target path for the written license file; if omitted, the file is written as `LICENSE` in the current directory.
FR10: User can resolve a user-defined alias to an SPDX ID when writing. [Growth]
FR11: User can write their configured favorite license by invoking `write` with no arguments. [Growth]
FR12: User receives an error when invoking `write` with no arguments and no favorite is configured. [Growth]
FR13: User can display help and usage information.
FR14: User receives clear, actionable error messages for invalid or unknown license IDs.
FR15: User receives a clear error when SPDX data cannot be fetched (e.g. network failure).
FR16: CLI exits with zero on success and non-zero on error.
FR17: User can run all commands in a non-interactive, scriptable manner (no prompts).
FR18: User has a `~/.ligma/` directory and `~/.ligma/config.json` created when absent, at the start of each run. [Growth]
FR19: User can set a favorite license in config for `write` with no arguments. [Growth]
FR20: User can define aliases that map custom names to SPDX IDs. [Growth]
FR21: User can override SPDX list and per-license URLs via config. [Growth]
FR22: User can have `ls` and `get` results served from a local cache under `~/.ligma/` when the cache is valid (within `cache_ttl`). [Growth]
FR23: User can bypass the cache for `ls` and `get` by setting `cache_ttl` to `0` in config. [Growth]
FR24: User can set cache TTL via the `cache_ttl` config property (how it is enforced TBD in implementation, e.g. stamp files). [Growth]
FR25: CLI obtains the license list and individual license details from the official SPDX source (hardcoded URLs in MVP).
FR26: CLI reads SPDX source URLs from the user's config when config is used. [Growth]

## NonFunctional Requirements

NFR-P1: Under typical network conditions, `ls` and `get` complete within **15 seconds** (network and SPDX response time dominate).
NFR-P2: CLI startup adds negligible delay before performing the requested command (no heavy init before the first network call).
NFR-S1: The CLI does not collect, store, or transmit user secrets or personal data. It only reads `~/.ligma/config.json` (Growth) and writes license files to user-specified paths.
NFR-I1: When the SPDX source is reachable and returns expected formats, the CLI successfully fetches and parses the license list and individual license data.
NFR-I2: When the SPDX source is unreachable (e.g. network failure, 4xx/5xx), the CLI fails with a clear error within a **30 second** timeout (or an implementation-defined, documented limit).

## Additional Requirements

- **Starter template (Epic 1 Story 1):** Cobra + cobra-cli. Initialize with: `go mod init <module-path>`, `cobra-cli init`, `cobra-cli add ls`, `cobra-cli add get`, `cobra-cli add write`. Project initialization using this sequence must be the first implementation story.
- **Build & distribution:** `go build -o ligma`; `go install` for local use. Distribution (GoReleaser, multi-arch, CI) deferred for MVP.
- **SPDX integration:** `net/http` GET, `encoding/json`; 30s timeout (NFR-I2). `licenses.json` for list; `details/{id}.json` for license text. Structs follow `license-list-data` layout. No third-party HTTP client.
- **Exit codes and errors:** 0 = success, 1 = usage, 2 = not found, 3 = I/O or network; all errors to stderr; stdout only for `ls` and `get` (and `--json` in Growth); `write` only to files.
- **Project structure:** `cmd/` (root.go, ls.go, get.go, write.go); `internal/spdx` (HTTP + JSON for SPDX); [Growth] `internal/config`, `internal/cache`. `internal/` not importable by external modules.
- **Cache [Growth]:** File-based under `~/.ligma/_cache/`. Layout: `_cache/list.json`, `_cache/details/<id>.json`. mtime-based TTL; `cache_ttl` 0 = always bypass (no `--no-cache` flag). On cache write failure, return fetched data; caching is best-effort.
- **Config [Growth]:** Viper; `~/.ligma/config.json`; schema: favorite, aliases, spdx_list_url, spdx_get_url_template, cache_ttl.
- **Security:** No authentication; NFR-S1: no secrets/PII; read config (Growth), write license files only to user-specified or default path.
- **Naming & patterns:** SPDX IDs as-is; commands `ls`, `get`, `write`, `help`; flags `--popular`, `--json` (Growth). Tests: `*_test.go` alongside package; `go test ./...`.

## FR Coverage Map

FR1: Epic 2 - List all SPDX licenses in the terminal
FR2: Epic 6 - Filter list by search term [Growth]
FR3: Epic 6 - Restrict to popular set [Growth]
FR4: Epic 6 - List in JSON [Growth]
FR5: Epic 3 - Output full license text by SPDX ID to stdout
FR6: Epic 7 - Alias→SPDX when viewing [Growth]
FR7: Epic 7 - License content in JSON [Growth]
FR8: Epic 4 - Write license to file by SPDX ID
FR9: Epic 4 - Target path or default LICENSE
FR10: Epic 7 - Alias→SPDX when writing [Growth]
FR11: Epic 7 - Write favorite with no args [Growth]
FR12: Epic 7 - Error when write no-args and no favorite [Growth]
FR13: Epic 1 - Help and usage
FR14: Epic 3, 4 - Clear errors for invalid/unknown license ID
FR15: Epic 2 - Clear error when SPDX unreachable
FR16: Epic 1 - Exit 0 on success, non-zero on error
FR17: Epic 1 - Non-interactive, scriptable
FR18: Epic 5 - ~/.ligma/ and config.json when absent [Growth]
FR19: Epic 5 - Favorite in config [Growth]
FR20: Epic 5 - Aliases in config [Growth]
FR21: Epic 5 - Override SPDX URLs in config [Growth]
FR22: Epic 5 - ls/get from cache when valid [Growth]
FR23: Epic 5 - Bypass cache via cache_ttl: 0 [Growth]
FR24: Epic 5 - cache_ttl in config [Growth]
FR25: Epic 2 - License list and details from SPDX
FR26: Epic 5 - SPDX URLs from config [Growth]
