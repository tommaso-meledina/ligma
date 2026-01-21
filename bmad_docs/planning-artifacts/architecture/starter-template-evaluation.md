# Starter Template Evaluation

## Primary Technology Domain

**CLI** — from PRD and project context: Go, single binary, subcommands `ls`, `get`, `write`, `help`.

## Starter Options Considered

- **Cobra** (`github.com/spf13/cobra` + `github.com/spf13/cobra-cli`): multi-command structure, `cobra-cli init` / `cobra-cli add`, built-in help and flags; pairs with Viper for config (Growth). Heavier than urfave/cli.
- **urfave/cli v3** (`github.com/urfave/cli/v3`): lighter, no generator; you define `cli.App` and `Commands` in code. Suited to small, single-purpose CLIs.

**Choice:** Cobra — better fit for four subcommands, future config, flags, and shell completion.

## Selected Starter: Cobra + cobra-cli

**Rationale**

- Subcommand/flag model matches `ls`, `get`, `write`, `help` and Growth flags (`--popular`, `--json`).
- `cobra-cli` gives a repeatable layout and `cobra-cli add` for new commands.
- Viper can be added when implementing `~/.ligma/config.json` in Growth; MVP can omit it.
- Shell completion (Phase 3) is supported by Cobra.

**Initialization Command**

```bash
go mod init <module-path>   # e.g. github.com/yourusername/ligma
go install github.com/spf13/cobra-cli@latest
cobra-cli init
cobra-cli add ls
cobra-cli add get
cobra-cli add write
```

(`help` is provided by the root command. Add `--viper` to `cobra-cli init` when implementing config in Growth.)

## Architectural Decisions Provided by Starter

**Language & Runtime:** Go; standard `go mod` and `main.go` → `cmd.Execute()`.

**CLI Structure:** `main.go` calls `cmd.Execute()`; `cmd/root.go` for root command and global flags; `cmd/ls.go`, `cmd/get.go`, `cmd/write.go` for subcommands; flags in each command.

**Build & Layout:** `go build -o ligma`; `go install` for a global binary; `cmd/` for commands; `internal/` or `pkg/` for SPDX client, (later) config, cache.

**Development:** `go run main.go [ls|get|write]`; `go test ./...`.

**Note:** Project initialization using this sequence should be the first implementation story.
