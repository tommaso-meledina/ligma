package cmd

import (
	"context"
	"testing"

	"github.com/tom/ligma/internal/cache"
	"github.com/tom/ligma/internal/config"
	"github.com/tom/ligma/internal/spdx"
)

func TestExecute_Success(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	rootCmd.SetArgs([]string{"ls"})
	defer rootCmd.SetArgs(nil)

	got := Execute()
	if got != 0 {
		t.Errorf("Execute() = %d, want 0", got)
	}
}

func TestExecute_UsageError(t *testing.T) {
	rootCmd.SetArgs([]string{"get"}) // missing required arg
	defer rootCmd.SetArgs(nil)

	got := Execute()
	if got != 1 {
		t.Errorf("Execute() = %d, want 1 (usage)", got)
	}
}

func TestExecute_WriteNoArgs(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	rootCmd.SetArgs([]string{"write"})
	defer rootCmd.SetArgs(nil)

	got := Execute()
	if got != 1 {
		t.Errorf("Execute() = %d, want 1 (no favorite set)", got)
	}
}

func TestExecute_WriteThreeArgs(t *testing.T) {
	rootCmd.SetArgs([]string{"write", "MIT", "a", "b"})
	defer rootCmd.SetArgs(nil)

	got := Execute()
	if got != 1 {
		t.Errorf("Execute() = %d, want 1 (usage)", got)
	}
}

func TestExecute_NotFound(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	save := cache.FetchDetailsFn
	cache.FetchDetailsFn = func(ctx context.Context, template, id string) (string, error) {
		return "", spdx.ErrNotFound
	}
	defer func() { cache.FetchDetailsFn = save }()

	rootCmd.SetArgs([]string{"get", "x"})
	defer rootCmd.SetArgs(nil)

	got := Execute()
	if got != 2 {
		t.Errorf("Execute() = %d, want 2 (not found)", got)
	}
}

func TestExecute_IOOrNetwork(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	rootCmd.SetArgs([]string{"get", "--simulate-io-error", "x"})
	defer rootCmd.SetArgs(nil)

	got := Execute()
	if got != 3 {
		t.Errorf("Execute() = %d, want 3 (I/O or network)", got)
	}
}
