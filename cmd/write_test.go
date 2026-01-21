package cmd

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/tom/ligma/internal/config"
	"github.com/tom/ligma/internal/spdx"
)

func TestWriteRunE_NotFound(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	save := writeFetchDetails
	writeFetchDetails = func(ctx context.Context, template, id string) (string, error) {
		return "", spdx.ErrNotFound
	}
	defer func() { writeFetchDetails = save }()

	err := writeCmd.RunE(writeCmd, []string{"x"})
	if err == nil {
		t.Fatal("RunE: expected error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("RunE: expected ErrNotFound, got %v", err)
	}
}

func TestWriteRunE_OneArg(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	save := writeFetchDetails
	writeFetchDetails = func(ctx context.Context, template, id string) (string, error) {
		return "license text", nil
	}
	defer func() { writeFetchDetails = save }()
	orig, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer func() { _ = os.Chdir(orig) }()

	err := writeCmd.RunE(writeCmd, []string{"MIT"})
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "LICENSE"))
	if err != nil {
		t.Fatalf("ReadFile LICENSE: %v", err)
	}
	if string(got) != "license text" {
		t.Errorf("LICENSE = %q, want %q", got, "license text")
	}
}

func TestWriteRunE_AliasResolved(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{"aliases":{"mit":"MIT"}}`), 0644); err != nil {
		t.Fatal(err)
	}

	save := writeFetchDetails
	writeFetchDetails = func(ctx context.Context, template, id string) (string, error) {
		if id != "MIT" {
			return "", spdx.ErrNotFound
		}
		return "license text", nil
	}
	defer func() { writeFetchDetails = save }()
	orig, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer func() { _ = os.Chdir(orig) }()

	err := writeCmd.RunE(writeCmd, []string{"mit"})
	if err != nil {
		t.Fatalf("RunE(write mit): %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "LICENSE"))
	if err != nil {
		t.Fatalf("ReadFile LICENSE: %v", err)
	}
	if string(got) != "license text" {
		t.Errorf("LICENSE = %q, want %q (alias mit must resolve to MIT)", got, "license text")
	}
}

func TestWriteRunE_TwoArgs(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	save := writeFetchDetails
	writeFetchDetails = func(ctx context.Context, template, id string) (string, error) {
		return "custom text", nil
	}
	defer func() { writeFetchDetails = save }()

	orig, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer func() { _ = os.Chdir(orig) }()

	err := writeCmd.RunE(writeCmd, []string{"GPL-3.0", "out.txt"})
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "out.txt"))
	if err != nil {
		t.Fatalf("ReadFile out.txt: %v", err)
	}
	if string(got) != "custom text" {
		t.Errorf("out.txt = %q, want %q", got, "custom text")
	}
}

func TestWriteRunE_WriteFails(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	save := writeFetchDetails
	writeFetchDetails = func(ctx context.Context, template, id string) (string, error) {
		return "x", nil
	}
	defer func() { writeFetchDetails = save }()

	dir := t.TempDir() // pass dir as path: WriteFile to a directory fails
	err := writeCmd.RunE(writeCmd, []string{"id", dir})
	if err == nil {
		t.Fatal("RunE: expected error (write to dir), got nil")
	}
	if !errors.Is(err, ErrIOOrNetwork) {
		t.Errorf("RunE: expected ErrIOOrNetwork, got %v", err)
	}
}

func TestWriteRunE_ZeroArgsFavorite(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{"favorite":"MIT"}`), 0644); err != nil {
		t.Fatal(err)
	}

	save := writeFetchDetails
	writeFetchDetails = func(ctx context.Context, template, id string) (string, error) {
		if id != "MIT" {
			return "", spdx.ErrNotFound
		}
		return "favorite text", nil
	}
	defer func() { writeFetchDetails = save }()
	orig, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer func() { _ = os.Chdir(orig) }()

	err := writeCmd.RunE(writeCmd, []string{})
	if err != nil {
		t.Fatalf("RunE(write with favorite): %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "LICENSE"))
	if err != nil {
		t.Fatalf("ReadFile LICENSE: %v", err)
	}
	if string(got) != "favorite text" {
		t.Errorf("LICENSE = %q, want %q", got, "favorite text")
	}
}

func TestWriteRunE_ZeroArgsNoFavorite(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	err := writeCmd.RunE(writeCmd, []string{})
	if err == nil {
		t.Fatal("RunE: expected error when favorite is not set")
	}
	if errors.Is(err, ErrNotFound) || errors.Is(err, ErrIOOrNetwork) {
		t.Errorf("RunE: expected usage/config error (exit 1), got %v", err)
	}
}

func TestWriteRunE_ZeroArgsFavoriteResolvedViaAlias(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{"favorite":"mit","aliases":{"mit":"MIT"}}`), 0644); err != nil {
		t.Fatal(err)
	}

	save := writeFetchDetails
	writeFetchDetails = func(ctx context.Context, template, id string) (string, error) {
		if id != "MIT" {
			return "", spdx.ErrNotFound
		}
		return "from alias", nil
	}
	defer func() { writeFetchDetails = save }()
	orig, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer func() { _ = os.Chdir(orig) }()

	err := writeCmd.RunE(writeCmd, []string{})
	if err != nil {
		t.Fatalf("RunE(write with favorite=mit alias): %v", err)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "LICENSE"))
	if string(got) != "from alias" {
		t.Errorf("LICENSE = %q, want from alias", got)
	}
}

func TestWriteRunE_UsesConfigGetURLTemplate(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{"spdx_get_url_template":"https://example.com/{id}.json"}`), 0644); err != nil {
		t.Fatal(err)
	}

	save := writeFetchDetails
	writeFetchDetails = func(ctx context.Context, template, id string) (string, error) {
		if template != "https://example.com/{id}.json" {
			return "", spdx.ErrNotFound
		}
		return "ok", nil
	}
	defer func() { writeFetchDetails = save }()
	orig, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer func() { _ = os.Chdir(orig) }()

	err := writeCmd.RunE(writeCmd, []string{"MIT"})
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "LICENSE"))
	if string(got) != "ok" {
		t.Errorf("LICENSE = %q, want ok", got)
	}
}
