# License Generation Minimal Assistant (LiGMA)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## What's LiGMA?

...it's a small CLI to list, view, and write SPDX licenses. For OSS maintainers who need the right license in a new project without copy‑pasting from the web (data comes from the [official SPDX JSON](https://github.com/spdx/license-list-data/tree/main/json)).

```bash
ligma get MIT
ligma write Apache-2.0
ligma ls --popular
```

---

## Install

**From source** (requires [Go 1.25+](https://go.dev/dl/)):

```bash
git clone https://github.com/tommaso-meledina/ligma.git /tmp/ligma && cd /tmp/ligma && go build -o ligma . && sudo mv ligma /usr/local/bin/ligma
```

---

## Usage

- List licenses: `ligma ls` (use `--popular` or `--filter <term>` to narrow)
- View full text of a license: `ligma get <SPDX-ID>`
- Write license to a file: `ligma write <SPDX-ID>` (writes to `LICENSE` in the current directory) or `ligma write <SPDX-ID> <path>`. With no arguments, `write` uses the configured favorite and writes to `LICENSE`.

Run `ligma <cmd> --help` for all flags.

---

## Commands

| Command | Description | Flags / notes |
|---------|-------------|---------------|
| `ls` | List available SPDX license IDs (and names when using `--json`). | `--json`, `--filter <term>`, `--popular` |
| `get <id>` | Fetch and print the full license text for an SPDX ID. | `--json` |
| `write [id] [path]` | Fetch the license by ID and write it to a file. If no args are provided, uses the configured `favorite` ID; if one arg is provided, it is interpreted as the ID of the license; if two args are provided, the second arg overrides the output path. Overwrites if the file exists. | — |

---

## Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Usage error (invalid flags, malformed args) |
| `2` | Not found error (e.g. unknown license ID) |
| `3` | I/O or network error (e.g. SPDX fetch or file write failure) |

---

## Configuration (optional)

The program creates a `config.json` file in `~/.ligma/`; you can uodate it in order to set a default `favorite` license ID (for calling `ligma write` with no args), `cache_ttl`, SPDX list/details URLs, and aliases. List and details are cached under `~/.ligma/_cache/`. Run `ligma <cmd> --help` or see the repository for details.

---

## License

[MIT](LICENSE)
