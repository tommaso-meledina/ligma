# Epic 1: CLI Foundation & Help

User can run the CLI, see help, and have ls/get/write stubs that exit correctly.
**FRs covered:** FR13, FR16, FR17

## Story 1.1: Go module, project structure, and dummy main

As a **developer**,
I want **to run `go mod init` and lay out the basic project structure with a dummy `main.go` that prints "OK"**,
So that **the project builds and running the binary produces a known output**.

**Acceptance Criteria:**

**Given** an empty project directory
**When** I run `go mod init <module-path>` and create a minimal layout (e.g. `go.mod`, `main.go`, `.gitignore` as needed) with a `main.go` whose `main` function prints `OK` to stdout
**Then** `go build -o ligma` succeeds and `./ligma` prints `OK` to stdout and exits with 0
**And** the layout is consistent with the eventual Architecture (e.g. `main.go` at repo root; `cmd/` and `internal/` may be added in later stories)

## Story 1.2: Add Cobra and ls, get, write subcommands

As a **developer**,
I want **to use cobra-cli to add the root command and `ls`, `get`, `write` subcommands**,
So that **`ligma ls`, `ligma get`, and `ligma write` are callable from the CLI regardless of their output**.

**Acceptance Criteria:**

**Given** a Go project with a working `main.go`
**When** I run `cobra-cli init` (integrating its output with existing `main.go` as needed) and `cobra-cli add ls`, `cobra-cli add get`, `cobra-cli add write`
**Then** `ligma ls`, `ligma get`, `ligma write` execute without error and are invokable from the command line
**And** the subcommands exist in the structure defined by Architecture (e.g. `cmd/root.go`, `cmd/ls.go`, `cmd/get.go`, `cmd/write.go`); the exact output or behavior of each subcommand is not yet specified
**And** `ligma --help` and `ligma <command> --help` display Cobra-generated help (FR13)

## Story 1.3: Refine exit codes

As a **developer**,
I want **the CLI to use consistent exit codes (0=success, 1=usage, 2=not found, 3=I/O or network)**,
So that **scripts and callers can reliably detect success and failure (FR16, FR17)**.

**Acceptance Criteria:**

**Given** the CLI with root, `ls`, `get`, and `write` commands
**When** a command completes successfully (e.g. `--help`), it exits with 0
**When** a command fails due to usage (e.g. invalid flags or malformed args where defined), it exits with 1
**When** a command fails due to not found (e.g. unknown license ID—stub behavior is acceptable for now), it exits with 2
**When** a command fails due to I/O or network (stub behavior acceptable for now), it exits with 3
**Then** all process exit is via `os.Exit` or equivalent from `main`/root; no prompts or interactive input are required (FR17)
**And** the exit code semantics are documented in code or a short comment for future implementation

---
