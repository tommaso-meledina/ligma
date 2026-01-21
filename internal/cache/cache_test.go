package cache

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tom/ligma/internal/config"
	"github.com/tom/ligma/internal/spdx"
)

func TestFetchList_Miss(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"licenses":[{"licenseId":"MIT","name":"MIT"}]}`))
	}))
	defer srv.Close()

	cacheDir := filepath.Join(dir, "_cache")
	list, err := FetchList(context.Background(), cacheDir, 3600, srv.URL)
	if err != nil {
		t.Fatalf("FetchList: %v", err)
	}
	if len(list) != 1 || list[0].LicenseID != "MIT" {
		t.Errorf("list = %+v", list)
	}
	// cache file should exist
	if _, err := os.Stat(filepath.Join(cacheDir, "list.json")); os.IsNotExist(err) {
		t.Error("list.json should have been written")
	}
}

func TestFetchList_Hit(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	cacheDir := filepath.Join(dir, "_cache")
	_ = os.MkdirAll(cacheDir, 0755)
	listPath := filepath.Join(cacheDir, "list.json")
	_ = os.WriteFile(listPath, []byte(`{"licenses":[{"licenseId":"X","name":"X"}]}`), 0644)

	// ensure mtime is recent
	_ = os.Chtimes(listPath, time.Now(), time.Now())

	// use a fetcher that would fail if called (no server)
	FetchListFn = func(ctx context.Context, url string) ([]spdx.License, error) {
		t.Fatal("fetcher should not be called on cache hit")
		return nil, nil
	}
	defer func() { FetchListFn = spdx.FetchLicenseList }()

	list, err := FetchList(context.Background(), cacheDir, 3600, "http://unused")
	if err != nil {
		t.Fatalf("FetchList: %v", err)
	}
	if len(list) != 1 || list[0].LicenseID != "X" {
		t.Errorf("list = %+v", list)
	}
}

func TestFetchList_Stale(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	cacheDir := filepath.Join(dir, "_cache")
	_ = os.MkdirAll(cacheDir, 0755)
	listPath := filepath.Join(cacheDir, "list.json")
	_ = os.WriteFile(listPath, []byte(`{"licenses":[{"licenseId":"old","name":"Old"}]}`), 0644)
	// set mtime to 2 hours ago; ttl 3600
	old := time.Now().Add(-2 * time.Hour)
	_ = os.Chtimes(listPath, old, old)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"licenses":[{"licenseId":"NEW","name":"New"}]}`))
	}))
	defer srv.Close()

	list, err := FetchList(context.Background(), cacheDir, 3600, srv.URL)
	if err != nil {
		t.Fatalf("FetchList: %v", err)
	}
	if len(list) != 1 || list[0].LicenseID != "NEW" {
		t.Errorf("list = %+v, want NEW", list)
	}
}

func TestFetchList_TTLZero(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	cacheDir := filepath.Join(dir, "_cache")
	_ = os.MkdirAll(cacheDir, 0755)
	_ = os.WriteFile(filepath.Join(cacheDir, "list.json"), []byte(`{"licenses":[{"licenseId":"cached","name":"Cached"}]}`), 0644)
	_ = os.Chtimes(filepath.Join(cacheDir, "list.json"), time.Now(), time.Now())

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"licenses":[{"licenseId":"fresh","name":"Fresh"}]}`))
	}))
	defer srv.Close()

	list, err := FetchList(context.Background(), cacheDir, 0, srv.URL)
	if err != nil {
		t.Fatalf("FetchList: %v", err)
	}
	if len(list) != 1 || list[0].LicenseID != "fresh" {
		t.Errorf("list = %+v, want fresh (ttl 0 must bypass cache)", list)
	}
}

func TestFetchList_WriteFailureStillReturns(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"licenses":[{"licenseId":"OK","name":"OK"}]}`))
	}))
	defer srv.Close()

	// cacheDir is a file, not a dir, so MkdirAll/WriteFile will fail
	cacheDir := filepath.Join(dir, "obstacle")
	_ = os.WriteFile(cacheDir, []byte("x"), 0644)

	list, err := FetchList(context.Background(), cacheDir, 3600, srv.URL)
	if err != nil {
		t.Fatalf("FetchList: %v (should return data despite write failure)", err)
	}
	if len(list) != 1 || list[0].LicenseID != "OK" {
		t.Errorf("list = %+v", list)
	}
}

func TestFetchDetails_Miss(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	save := FetchDetailsFn
	FetchDetailsFn = func(ctx context.Context, template, id string) (string, error) {
		return "license text", nil
	}
	defer func() { FetchDetailsFn = save }()

	cacheDir := filepath.Join(dir, "_cache")
	text, err := FetchDetails(context.Background(), cacheDir, 3600, "http://x/{id}.json", "MIT")
	if err != nil {
		t.Fatalf("FetchDetails: %v", err)
	}
	if text != "license text" {
		t.Errorf("text = %q", text)
	}
	p := filepath.Join(cacheDir, "details", "MIT.json")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		t.Error("details/MIT.json should have been written")
	}
}

func TestFetchDetails_Hit(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	cacheDir := filepath.Join(dir, "_cache")
	detailsDir := filepath.Join(cacheDir, "details")
	_ = os.MkdirAll(detailsDir, 0755)
	path := filepath.Join(detailsDir, "MIT.json")
	b, _ := json.Marshal(struct{ LicenseText string }{"cached text"})
	_ = os.WriteFile(path, b, 0644)
	_ = os.Chtimes(path, time.Now(), time.Now())

	save := FetchDetailsFn
	FetchDetailsFn = func(ctx context.Context, template, id string) (string, error) {
		t.Fatal("fetcher should not be called on cache hit")
		return "", nil
	}
	defer func() { FetchDetailsFn = save }()

	text, err := FetchDetails(context.Background(), cacheDir, 3600, "http://unused/{id}.json", "MIT")
	if err != nil {
		t.Fatalf("FetchDetails: %v", err)
	}
	if text != "cached text" {
		t.Errorf("text = %q, want cached text", text)
	}
}

func TestFetchDetails_Stale(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	cacheDir := filepath.Join(dir, "_cache")
	_ = os.MkdirAll(filepath.Join(cacheDir, "details"), 0755)
	path := filepath.Join(cacheDir, "details", "MIT.json")
	b, _ := json.Marshal(struct{ LicenseText string }{"old"})
	_ = os.WriteFile(path, b, 0644)
	old := time.Now().Add(-2 * time.Hour)
	_ = os.Chtimes(path, old, old)

	save := FetchDetailsFn
	FetchDetailsFn = func(ctx context.Context, template, id string) (string, error) {
		return "fresh", nil
	}
	defer func() { FetchDetailsFn = save }()

	text, err := FetchDetails(context.Background(), cacheDir, 3600, "http://unused/{id}.json", "MIT")
	if err != nil {
		t.Fatalf("FetchDetails: %v", err)
	}
	if text != "fresh" {
		t.Errorf("text = %q, want fresh (stale should refetch)", text)
	}
}

func TestFetchDetails_TTLZero(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	cacheDir := filepath.Join(dir, "_cache")
	_ = os.MkdirAll(filepath.Join(cacheDir, "details"), 0755)
	path := filepath.Join(cacheDir, "details", "MIT.json")
	b, _ := json.Marshal(struct{ LicenseText string }{"cached"})
	_ = os.WriteFile(path, b, 0644)
	_ = os.Chtimes(path, time.Now(), time.Now())

	save := FetchDetailsFn
	FetchDetailsFn = func(ctx context.Context, template, id string) (string, error) {
		return "fresh", nil
	}
	defer func() { FetchDetailsFn = save }()

	text, err := FetchDetails(context.Background(), cacheDir, 0, "http://unused/{id}.json", "MIT")
	if err != nil {
		t.Fatalf("FetchDetails: %v", err)
	}
	if text != "fresh" {
		t.Errorf("text = %q, want fresh (ttl 0 must bypass cache)", text)
	}
}

func TestFetchDetails_pathTraversalSkipsCache(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	cacheDir := filepath.Join(dir, "_cache")
	_ = os.MkdirAll(filepath.Join(cacheDir, "details"), 0755)

	save := FetchDetailsFn
	FetchDetailsFn = func(ctx context.Context, template, id string) (string, error) {
		return "from-fetcher", nil
	}
	defer func() { FetchDetailsFn = save }()

	text, err := FetchDetails(context.Background(), cacheDir, 3600, "http://x/{id}.json", "a/b")
	if err != nil {
		t.Fatalf("FetchDetails: %v", err)
	}
	if text != "from-fetcher" {
		t.Errorf("text = %q, want from-fetcher (id with / skips cache read, calls fetcher)", text)
	}
}

func TestFetchList_CorruptJSONReturnsError(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	cacheDir := filepath.Join(dir, "_cache")
	_ = os.MkdirAll(cacheDir, 0755)
	_ = os.WriteFile(filepath.Join(cacheDir, "list.json"), []byte(`{invalid`), 0644)
	_ = os.Chtimes(filepath.Join(cacheDir, "list.json"), time.Now(), time.Now())

	_, err := FetchList(context.Background(), cacheDir, 3600, "http://unused")
	if err == nil {
		t.Fatal("FetchList: expected error for corrupt cache JSON")
	}
}

func TestFetchDetails_CorruptJSONReturnsError(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	cacheDir := filepath.Join(dir, "_cache")
	_ = os.MkdirAll(filepath.Join(cacheDir, "details"), 0755)
	_ = os.WriteFile(filepath.Join(cacheDir, "details", "MIT.json"), []byte(`{invalid`), 0644)
	_ = os.Chtimes(filepath.Join(cacheDir, "details", "MIT.json"), time.Now(), time.Now())

	_, err := FetchDetails(context.Background(), cacheDir, 3600, "http://unused/{id}.json", "MIT")
	if err == nil {
		t.Fatal("FetchDetails: expected error for corrupt cache JSON")
	}
}

func TestFetchDetails_WriteFailureStillReturns(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDirOverride(dir)
	defer config.SetConfigDirOverride("")

	save := FetchDetailsFn
	FetchDetailsFn = func(ctx context.Context, template, id string) (string, error) {
		return "fetched", nil
	}
	defer func() { FetchDetailsFn = save }()

	// _cache/details as a file so writing MIT.json under it fails
	cacheDir := filepath.Join(dir, "_cache")
	_ = os.MkdirAll(cacheDir, 0755)
	_ = os.WriteFile(filepath.Join(cacheDir, "details"), []byte("x"), 0644)

	text, err := FetchDetails(context.Background(), cacheDir, 3600, "http://x/{id}.json", "MIT")
	if err != nil {
		t.Fatalf("FetchDetails: %v (should return data despite write failure)", err)
	}
	if text != "fetched" {
		t.Errorf("text = %q", text)
	}
}

func TestTTL(t *testing.T) {
	if TTL(nil) != defaultTTL {
		t.Errorf("TTL(nil) = %d, want %d", TTL(nil), defaultTTL)
	}
	z := 0
	if TTL(&z) != 0 {
		t.Errorf("TTL(0) = %d, want 0", TTL(&z))
	}
	n := 100
	if TTL(&n) != 100 {
		t.Errorf("TTL(100) = %d, want 100", TTL(&n))
	}
}
