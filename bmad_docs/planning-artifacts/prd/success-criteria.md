# Success Criteria

## User Success

- **Completion:** User has a correct `LICENSE` in the repo and can move on.
- **Relief:** One short command instead of copy-pasting; with **favorite** in config, `ligma write` (no args) covers the common case.
- **Friction:** One command to list, view, or write; optional favorite for zero-argument write.

## Business Success

- **Context:** Hobby / OSS project. No formal business metrics.
- **Indicators:** Personal use and others finding it useful.

## Technical Success

- Correct license text from the official SPDX source.
- Scriptable, non-interactive CLI (`get` / `write`; clear exit codes, no prompts).
- Data from SPDX `license-list-data` (JSON).

## Measurable Outcomes

- User can go from "I need a license" to a correct `LICENSE` with one or two commands.
- MVP: `get` and `write` by SPDX ID only. Growth: aliases, favorite, `ligma write` with no args when favorite is set.
