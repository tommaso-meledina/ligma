package cmd

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tom/ligma/internal/config"
	"github.com/tom/ligma/internal/spdx"
)

const lsGoodJSON = `{"licenseListVersion":"1.0","licenses":[{"licenseId":"MIT","name":"MIT License"}]}`

// lsMultiJSON has the 5 popular IDs plus X for --popular and --filter tests.
const lsMultiJSON = `{"licenses":[
  {"licenseId":"MIT","name":"MIT License"},
  {"licenseId":"Apache-2.0","name":"Apache License 2.0"},
  {"licenseId":"GPL-2.0","name":"GNU General Public License v2.0"},
  {"licenseId":"BSD-3-Clause","name":"BSD 3-Clause"},
  {"licenseId":"ISC","name":"ISC License"},
  {"licenseId":"X","name":"X License"}
]}`

func TestLsRunE_Success(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(lsGoodJSON))
	}))
	defer srv.Close()

	lsListURLOverride = srv.URL
	defer func() { lsListURLOverride = "" }()

	err := lsCmd.RunE(lsCmd, []string{})
	if err != nil {
		t.Errorf("RunE: expected nil, got %v", err)
	}
}

func TestLsRunE_FetchErrorMapsToIOOrNetwork(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	lsListURLOverride = srv.URL
	defer func() { lsListURLOverride = "" }()

	err := lsCmd.RunE(lsCmd, []string{})
	if err == nil {
		t.Fatal("RunE: expected error")
	}
	if !errors.Is(err, ErrIOOrNetwork) {
		t.Errorf("RunE: expected ErrIOOrNetwork, got %v", err)
	}
}

func TestLsRunE_UsesConfigListURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(lsGoodJSON))
	}))
	defer srv.Close()

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")
	if err := os.WriteFile(cfgPath, []byte(`{"spdx_list_url":"`+srv.URL+`"}`), 0644); err != nil {
		t.Fatal(err)
	}
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	// No lsListURLOverride: runLs must use config's spdx_list_url
	err := lsCmd.RunE(lsCmd, []string{})
	if err != nil {
		t.Errorf("RunE: %v", err)
	}
}

func TestLsRunE_JSON(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(lsGoodJSON))
	}))
	defer srv.Close()

	lsListURLOverride = srv.URL
	defer func() { lsListURLOverride = "" }()

	_ = lsCmd.Flags().Set("json", "true")
	defer func() { _ = lsCmd.Flags().Set("json", "false") }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = w

	err = lsCmd.RunE(lsCmd, []string{})
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}

	out, _ := io.ReadAll(r)
	var list []spdx.License
	if err := json.Unmarshal(out, &list); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\nraw: %s", err, out)
	}
	if len(list) != 1 || list[0].LicenseID != "MIT" || list[0].Name != "MIT License" {
		t.Errorf("json list = %+v", list)
	}
}

func TestLsRunE_NoJSONHumanFormat(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(lsGoodJSON))
	}))
	defer srv.Close()

	lsListURLOverride = srv.URL
	defer func() { lsListURLOverride = "" }()

	_ = lsCmd.Flags().Set("json", "false")

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	err := lsCmd.RunE(lsCmd, []string{})
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}

	out, _ := io.ReadAll(r)
	// Human format: one ID per line, no JSON
	if len(out) > 0 && out[0] == '[' {
		t.Errorf("human format should not be JSON array; got: %s", out)
	}
	if string(out) != "MIT\n" {
		t.Errorf("human output = %q, want MIT\\n", out)
	}
}

func TestLsRunE_Filter(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(lsMultiJSON))
	}))
	defer srv.Close()

	lsListURLOverride = srv.URL
	defer func() { lsListURLOverride = "" }()

	_ = lsCmd.Flags().Set("filter", "apache")
	defer func() { _ = lsCmd.Flags().Set("filter", "") }()

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	err := lsCmd.RunE(lsCmd, []string{})
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}
	out, _ := io.ReadAll(r)
	if string(out) != "Apache-2.0\n" {
		t.Errorf("--filter apache: got %q, want Apache-2.0\\n", out)
	}
}

func TestLsRunE_FilterCaseInsensitive(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(lsGoodJSON))
	}))
	defer srv.Close()

	lsListURLOverride = srv.URL
	defer func() { lsListURLOverride = "" }()

	_ = lsCmd.Flags().Set("filter", "mit")
	defer func() { _ = lsCmd.Flags().Set("filter", "") }()

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	err := lsCmd.RunE(lsCmd, []string{})
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}
	out, _ := io.ReadAll(r)
	if string(out) != "MIT\n" {
		t.Errorf("--filter mit (case-insensitive): got %q, want MIT\\n", out)
	}
}

func TestLsRunE_Popular(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(lsMultiJSON))
	}))
	defer srv.Close()

	lsListURLOverride = srv.URL
	defer func() { lsListURLOverride = "" }()

	_ = lsCmd.Flags().Set("popular", "true")
	defer func() { _ = lsCmd.Flags().Set("popular", "false") }()

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	err := lsCmd.RunE(lsCmd, []string{})
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}
	out, _ := io.ReadAll(r)
	lines := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	got := make(map[string]bool)
	for _, s := range lines {
		got[s] = true
	}
	want := map[string]bool{"MIT": true, "Apache-2.0": true, "GPL-2.0": true, "BSD-3-Clause": true, "ISC": true}
	if len(got) != 5 || got["X"] {
		t.Errorf("--popular: got %v, want exactly the 5 popular (no X)", got)
	}
	for id := range want {
		if !got[id] {
			t.Errorf("--popular: missing %s in %v", id, got)
		}
	}
}

func TestLsRunE_PopularAndFilter(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(lsMultiJSON))
	}))
	defer srv.Close()

	lsListURLOverride = srv.URL
	defer func() { lsListURLOverride = "" }()

	_ = lsCmd.Flags().Set("popular", "true")
	_ = lsCmd.Flags().Set("filter", "2")
	defer func() { _ = lsCmd.Flags().Set("popular", "false"); _ = lsCmd.Flags().Set("filter", "") }()

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	err := lsCmd.RunE(lsCmd, []string{})
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}
	out, _ := io.ReadAll(r)
	lines := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	got := make(map[string]bool)
	for _, s := range lines {
		got[s] = true
	}
	// Popular ∩ filter "2": Apache-2.0, GPL-2.0
	if len(got) != 2 || !got["Apache-2.0"] || !got["GPL-2.0"] {
		t.Errorf("--popular --filter 2: got %v, want {Apache-2.0, GPL-2.0}", got)
	}
}

func TestLsRunE_FilterMatchesName(t *testing.T) {
	config.SetConfigDirOverride(t.TempDir())
	defer config.SetConfigDirOverride("")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(lsMultiJSON))
	}))
	defer srv.Close()

	lsListURLOverride = srv.URL
	defer func() { lsListURLOverride = "" }()

	// "gnu" matches "GNU General Public License v2.0" (name) but not "GPL-2.0" (licenseId)
	_ = lsCmd.Flags().Set("filter", "gnu")
	defer func() { _ = lsCmd.Flags().Set("filter", "") }()

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	err := lsCmd.RunE(lsCmd, []string{})
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("RunE: %v", err)
	}
	out, _ := io.ReadAll(r)
	if string(out) != "GPL-2.0\n" {
		t.Errorf("--filter gnu (matches name only): got %q, want GPL-2.0\\n", out)
	}
}
