/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import "errors"

// Exit code semantics (all os.Exit only from main/root):
//
//   - 0 = success
//   - 1 = usage (invalid flags, malformed args, or any error not 2/3)
//   - 2 = not found (e.g. unknown license ID)
//   - 3 = I/O or network (e.g. SPDX fetch failure, file write failure when implemented in later stories)
var (
	ErrNotFound    = errors.New("not found")
	ErrIOOrNetwork = errors.New("I/O or network error")
)

// exitCodeFrom maps an error to exit code 0–3. Used by Execute; extracted for testing.
func exitCodeFrom(err error) int {
	if err == nil {
		return 0
	}
	if errors.Is(err, ErrNotFound) {
		return 2
	}
	if errors.Is(err, ErrIOOrNetwork) {
		return 3
	}
	return 1
}
