package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_CreateIfAbsent(t *testing.T) {
	dir := t.TempDir()
	SetConfigDirOverride(dir)
	defer SetConfigDirOverride("")

	path := filepath.Join(dir, "config.json")
	if _, err := os.Stat(path); err == nil {
		t.Fatal("config.json should not exist yet")
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load: expected config, got nil")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("Load: config.json should have been created")
	}
	// Second Load reads the file
	cfg2, err := Load()
	if err != nil {
		t.Fatalf("Load second: %v", err)
	}
	if cfg2.SPDXListURL != cfg.SPDXListURL {
		t.Errorf("second Load SPDXListURL = %q, want %q", cfg2.SPDXListURL, cfg.SPDXListURL)
	}
}

func TestLoad_EmptyFile_Defaults(t *testing.T) {
	dir := t.TempDir()
	SetConfigDirOverride(dir)
	defer SetConfigDirOverride("")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	// Default URLs from Viper
	if cfg.SPDXListURL != defaultListURL {
		t.Errorf("SPDXListURL = %q, want %q", cfg.SPDXListURL, defaultListURL)
	}
	if cfg.SPDXGetURLTemplate != defaultDetailsURLTmpl {
		t.Errorf("SPDXGetURLTemplate = %q, want %q", cfg.SPDXGetURLTemplate, defaultDetailsURLTmpl)
	}
	if cfg.Favorite != nil {
		t.Errorf("Favorite = %v, want nil", cfg.Favorite)
	}
	if cfg.CacheTTL != nil {
		t.Errorf("CacheTTL = %v, want nil", cfg.CacheTTL)
	}
	if cfg.Aliases == nil || len(cfg.Aliases) != 0 {
		t.Errorf("Aliases = %v, want empty map", cfg.Aliases)
	}
}

func TestLoad_WithContent(t *testing.T) {
	dir := t.TempDir()
	SetConfigDirOverride(dir)
	defer SetConfigDirOverride("")

	path := filepath.Join(dir, "config.json")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	// Pre-create with content
	body := `{"favorite":"MIT","aliases":{"apache":"Apache-2.0"},"cache_ttl":0}`
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Favorite == nil || *cfg.Favorite != "MIT" {
		t.Errorf("Favorite = %v, want *MIT", cfg.Favorite)
	}
	if cfg.Aliases["apache"] != "Apache-2.0" {
		t.Errorf("Aliases[apache] = %q, want Apache-2.0", cfg.Aliases["apache"])
	}
	if cfg.CacheTTL == nil || *cfg.CacheTTL != 0 {
		t.Errorf("CacheTTL = %v, want *0", cfg.CacheTTL)
	}
}

func TestLoad_FavoriteEmptyStringYieldsNil(t *testing.T) {
	dir := t.TempDir()
	SetConfigDirOverride(dir)
	defer SetConfigDirOverride("")

	path := filepath.Join(dir, "config.json")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"favorite":""}`), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Favorite != nil {
		t.Errorf("Favorite with empty string should be nil, got %v", cfg.Favorite)
	}
}

func TestLigmaDir_Override(t *testing.T) {
	dir := t.TempDir()
	SetConfigDirOverride(dir)
	defer SetConfigDirOverride("")

	got, err := LigmaDir()
	if err != nil {
		t.Fatalf("LigmaDir: %v", err)
	}
	if got != dir {
		t.Errorf("LigmaDir() = %q, want %q", got, dir)
	}
}

func TestLigmaDir_UserHomeDir(t *testing.T) {
	SetConfigDirOverride("")
	defer SetConfigDirOverride("")

	got, err := LigmaDir()
	if err != nil {
		t.Fatalf("LigmaDir: %v", err)
	}
	if got == "" {
		t.Error("LigmaDir() = empty")
	}
	if !filepath.IsAbs(got) && got != "" {
		t.Errorf("LigmaDir() = %q, expect absolute or non-empty", got)
	}
	// Should end with .ligma when using UserHomeDir
	if len(got) > 0 && filepath.Base(got) != ".ligma" {
		t.Errorf("LigmaDir() = %q, expect path ending in .ligma", got)
	}
}

func TestLoad_InvalidJSONReturnsError(t *testing.T) {
	dir := t.TempDir()
	SetConfigDirOverride(dir)
	defer SetConfigDirOverride("")

	path := filepath.Join(dir, "config.json")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{invalid json}`), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load()
	if err == nil {
		t.Fatal("Load: expected error for invalid config JSON")
	}
}

func TestResolve_AliasKeyReturnsMappedID(t *testing.T) {
	cfg := &Config{Aliases: map[string]string{"x": "MIT", "apache": "Apache-2.0"}}
	if got := cfg.Resolve("x"); got != "MIT" {
		t.Errorf("Resolve(%q) = %q, want MIT", "x", got)
	}
	if got := cfg.Resolve("apache"); got != "Apache-2.0" {
		t.Errorf("Resolve(%q) = %q, want Apache-2.0", "apache", got)
	}
}

func TestResolve_NonKeyReturnsAsIs(t *testing.T) {
	cfg := &Config{Aliases: map[string]string{"x": "MIT"}}
	if got := cfg.Resolve("MIT"); got != "MIT" {
		t.Errorf("Resolve(%q) = %q, want MIT (pass-through)", "MIT", got)
	}
	if got := cfg.Resolve("unknown"); got != "unknown" {
		t.Errorf("Resolve(%q) = %q, want unknown", "unknown", got)
	}
}

func TestResolve_NilAliasesReturnsAsIs(t *testing.T) {
	cfg := &Config{Aliases: nil}
	if got := cfg.Resolve("MIT"); got != "MIT" {
		t.Errorf("Resolve with nil Aliases: got %q, want MIT", got)
	}
}

func TestLoad_WithoutOverrideUsesUserHomeDir(t *testing.T) {
	SetConfigDirOverride("")
	defer SetConfigDirOverride("")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load: expected config")
	}
	// Load used UserHomeDir; we can't assert path without duplicating logic
	dir, err := LigmaDir()
	if err != nil {
		t.Fatalf("LigmaDir: %v", err)
	}
	if dir == "" {
		t.Error("LigmaDir empty")
	}
}
