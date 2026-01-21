package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/tom/ligma/internal/cache"
	"github.com/tom/ligma/internal/config"
	"github.com/tom/ligma/internal/spdx"
)

func TestGetRunE_NotFound(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	getSimulateIO = false
	save := cache.FetchDetailsFn
	cache.FetchDetailsFn = func(ctx context.Context, template, id string) (string, error) {
		return "", spdx.ErrNotFound
	}
	defer func() { getSimulateIO = false; cache.FetchDetailsFn = save }()

	err := getCmd.RunE(getCmd, []string{"x"})
	if err == nil {
		t.Fatal("RunE: expected error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("RunE: expected ErrNotFound, got %v", err)
	}
}

func TestGetRunE_SimulateIO(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	getSimulateIO = true
	defer func() { getSimulateIO = false }()

	err := getCmd.RunE(getCmd, []string{"x"})
	if err == nil {
		t.Fatal("RunE: expected error, got nil")
	}
	if !errors.Is(err, ErrIOOrNetwork) {
		t.Errorf("RunE: expected ErrIOOrNetwork, got %v", err)
	}
}

func TestGetRunE_Success(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	getSimulateIO = false
	save := cache.FetchDetailsFn
	cache.FetchDetailsFn = func(ctx context.Context, template, id string) (string, error) {
		return "license text", nil
	}
	defer func() { cache.FetchDetailsFn = save }()

	err := getCmd.RunE(getCmd, []string{"MIT"})
	if err != nil {
		t.Errorf("RunE: expected nil, got %v", err)
	}
}

func TestGetRunE_AliasResolved(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{"aliases":{"mit":"MIT"}}`), 0644); err != nil {
		t.Fatal(err)
	}

	getSimulateIO = false
	save := cache.FetchDetailsFn
	cache.FetchDetailsFn = func(ctx context.Context, template, id string) (string, error) {
		if id != "MIT" {
			return "", spdx.ErrNotFound
		}
		return "from MIT", nil
	}
	defer func() { cache.FetchDetailsFn = save }()

	// get mit -> resolve to MIT -> fetch MIT -> "from MIT"
	err := getCmd.RunE(getCmd, []string{"mit"})
	if err != nil {
		t.Errorf("RunE(get mit): %v", err)
	}
	// stdout is not captured in RunE; the mock returns "from MIT" only when id is MIT, so no error means resolve worked
}

func TestGetRunE_UnknownAliasPassesThrough(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{"aliases":{"mit":"MIT"}}`), 0644); err != nil {
		t.Fatal(err)
	}

	save := cache.FetchDetailsFn
	cache.FetchDetailsFn = func(ctx context.Context, template, id string) (string, error) {
		return "", spdx.ErrNotFound
	}
	defer func() { cache.FetchDetailsFn = save }()

	err := getCmd.RunE(getCmd, []string{"notanalias"})
	if err == nil {
		t.Fatal("RunE: expected error for unknown alias (pass-through to SPDX, 404)")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("RunE: expected ErrNotFound, got %v", err)
	}
}

func TestGetRunE_JSON(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	save := cache.FetchDetailsFn
	cache.FetchDetailsFn = func(ctx context.Context, template, id string) (string, error) {
		return "license text", nil
	}
	defer func() { cache.FetchDetailsFn = save }()

	_ = getCmd.Flags().Set("json", "true")
	defer func() { _ = getCmd.Flags().Set("json", "false") }()

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	err := getCmd.RunE(getCmd, []string{"MIT"})
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}
	out, _ := io.ReadAll(r)
	var v struct {
		ID          string `json:"id"`
		LicenseText string `json:"licenseText"`
	}
	if err := json.Unmarshal(out, &v); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\nraw: %s", err, out)
	}
	if v.ID != "MIT" || v.LicenseText != "license text" {
		t.Errorf("JSON: id=%q licenseText=%q, want id=MIT licenseText=license text", v.ID, v.LicenseText)
	}
}

func TestGetRunE_NoJSONPlainText(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	save := cache.FetchDetailsFn
	cache.FetchDetailsFn = func(ctx context.Context, template, id string) (string, error) {
		return "license text", nil
	}
	defer func() { cache.FetchDetailsFn = save }()

	_ = getCmd.Flags().Set("json", "false")

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	err := getCmd.RunE(getCmd, []string{"MIT"})
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}
	out, _ := io.ReadAll(r)
	if string(out) != "license text" {
		t.Errorf("plain output = %q, want license text", out)
	}
	if len(out) > 0 && out[0] == '{' {
		t.Errorf("plain output should not be JSON; got: %s", out)
	}
}

func TestGetRunE_UsesConfigGetURLTemplate(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{"spdx_get_url_template":"https://example.com/details/{id}.json"}`), 0644); err != nil {
		t.Fatal(err)
	}

	save := cache.FetchDetailsFn
	cache.FetchDetailsFn = func(ctx context.Context, template, id string) (string, error) {
		if template != "https://example.com/details/{id}.json" {
			return "", spdx.ErrNotFound
		}
		return "ok", nil
	}
	defer func() { cache.FetchDetailsFn = save }()

	err := getCmd.RunE(getCmd, []string{"MIT"})
	if err != nil {
		t.Errorf("RunE: %v", err)
	}
}

func TestGetRunE_FetchErrorMapsToIOOrNetwork(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	save := cache.FetchDetailsFn
	cache.FetchDetailsFn = func(ctx context.Context, template, id string) (string, error) {
		return "", fmt.Errorf("connection refused")
	}
	defer func() { cache.FetchDetailsFn = save }()

	err := getCmd.RunE(getCmd, []string{"MIT"})
	if err == nil {
		t.Fatal("RunE: expected error")
	}
	if !errors.Is(err, ErrIOOrNetwork) {
		t.Errorf("RunE: expected ErrIOOrNetwork, got %v", err)
	}
}
