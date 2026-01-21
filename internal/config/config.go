package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// configDirOverride, when non-empty, replaces ~/.ligma as the config directory for Load. Used by tests.
var configDirOverride string

// SetConfigDirOverride sets the config directory for Load (e.g. a temp dir in tests). Restore to "" when done.
func SetConfigDirOverride(dir string) {
	configDirOverride = dir
}

// Resolve resolves an alias to an SPDX ID. If idOrAlias is a key in Aliases, the mapped
// SPDX ID is returned. Otherwise idOrAlias is returned as-is (treated as an SPDX ID).
// Aliases take precedence over a literal SPDX ID (e.g. an alias "MIT" -> "Apache-2.0" would win).
func (c *Config) Resolve(idOrAlias string) string {
	if c.Aliases == nil {
		return idOrAlias
	}
	if v, ok := c.Aliases[idOrAlias]; ok {
		return v
	}
	return idOrAlias
}

// LigmaDir returns the ~/.ligma directory (or configDirOverride when set). Used by cache to resolve _cache.
func LigmaDir() (string, error) {
	if configDirOverride != "" {
		return configDirOverride, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("config: home dir: %w", err)
	}
	return filepath.Join(home, ".ligma"), nil
}

// Config holds the parsed config. Only favorite, aliases, spdx_list_url, spdx_get_url_template, cache_ttl (NFR-S1: no secrets, no PII).
type Config struct {
	Favorite           *string
	Aliases            map[string]string
	SPDXListURL        string
	SPDXGetURLTemplate string
	CacheTTL           *int
}

const (
	defaultListURL        = "https://raw.githubusercontent.com/spdx/license-list-data/main/json/licenses.json"
	defaultDetailsURLTmpl = "https://raw.githubusercontent.com/spdx/license-list-data/main/json/details/{id}.json"
)

// Load creates ~/.ligma/ and ~/.ligma/config.json if absent (with {}), then reads and parses config.
// Uses os.UserHomeDir() to resolve ~. When configDirOverride is set (e.g. in tests), uses that instead of ~/.ligma.
func Load() (*Config, error) {
	var dir string
	if configDirOverride != "" {
		dir = configDirOverride
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("config: home dir: %w", err)
		}
		dir = filepath.Join(home, ".ligma")
	}
	path := filepath.Join(dir, "config.json")

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("config: mkdir %s: %w", dir, err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte("{}"), 0644); err != nil {
			return nil, fmt.Errorf("config: create %s: %w", path, err)
		}
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetDefault("spdx_list_url", defaultListURL)
	v.SetDefault("spdx_get_url_template", defaultDetailsURLTmpl)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}

	cfg := &Config{
		SPDXListURL:        v.GetString("spdx_list_url"),
		SPDXGetURLTemplate: v.GetString("spdx_get_url_template"),
		Aliases:            v.GetStringMapString("aliases"),
	}
	if cfg.Aliases == nil {
		cfg.Aliases = make(map[string]string)
	}
	if v.IsSet("favorite") {
		if s := v.GetString("favorite"); s != "" {
			cfg.Favorite = &s
		}
	}
	if v.IsSet("cache_ttl") {
		n := v.GetInt("cache_ttl")
		cfg.CacheTTL = &n
	}
	return cfg, nil
}
