# Project Structure & Boundaries

## Project Directory Structure

```
ligma/
├── README.md
├── go.mod
├── go.sum
├── .gitignore
├── main.go
├── cmd/
│   ├── root.go
│   ├── ls.go
│   ├── get.go
│   └── write.go
└── internal/
    ├── spdx/           # HTTP client, fetch list & details, parse JSON
    │   ├── client.go
    │   └── client_test.go
    ├── config/         # [Growth] load ~/.ligma/config.json (Viper)
    │   └── config.go
    └── cache/          # [Growth] file-based cache wrapper
        ├── cache.go
        └── cache_test.go
```

Runtime (not in repo): `~/.ligma/` (config.json, `_cache/`, `_cache/list.json`, `_cache/details/<id>.json`).

## Boundaries

- **SPDX:** `internal/spdx` is the only place that does HTTP and JSON parse for SPDX. All `ls`/`get`/`write` logic calls into it (or via cache in Growth).
- **Cache [Growth]:** `internal/cache` wraps SPDX fetches; used by `ls` and `get` only.
- **Config [Growth]:** `internal/config` loads and parses `~/.ligma/config.json`; called at startup by commands that need it.
- **Commands:** `cmd/*` handle flags, args, and orchestration; business logic in `internal/`.

## Requirements → Structure

- **FR1–FR4, FR25, FR26:** `internal/spdx`, `cmd/ls`.
- **FR5–FR7:** `internal/spdx`, `cmd/get`; [Growth] `internal/cache`, `internal/config` (aliases).
- **FR8–FR12:** `internal/spdx`, `cmd/write`; [Growth] `internal/config` (favorite, aliases).
- **FR13–FR17:** Cobra `help`; `cmd/*` for errors, stderr, exit codes.
- **FR18–FR24:** `internal/config`, `internal/cache`; `~/.ligma/` and `_cache/` layout.

---
