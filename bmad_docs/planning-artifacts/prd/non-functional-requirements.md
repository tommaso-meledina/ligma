# Non-Functional Requirements

## Performance

- **NFR-P1:** Under typical network conditions, `ls` and `get` complete within **15 seconds** (network and SPDX response time dominate).
- **NFR-P2:** CLI startup adds negligible delay before performing the requested command (no heavy init before the first network call).

## Security

- **NFR-S1:** The CLI does not collect, store, or transmit user secrets or personal data. It only reads `~/.ligma/config.json` (Growth) and writes license files to user-specified paths.

## Integration (SPDX)

- **NFR-I1:** When the SPDX source is reachable and returns expected formats, the CLI successfully fetches and parses the license list and individual license data.
- **NFR-I2:** When the SPDX source is unreachable (e.g. network failure, 4xx/5xx), the CLI fails with a clear error within a **30 second** timeout (or an implementation-defined, documented limit).
