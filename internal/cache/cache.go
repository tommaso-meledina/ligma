package cache

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tom/ligma/internal/spdx"
)

// FetchListFn is the list fetcher; defaults to spdx.FetchLicenseList. Override in tests.
var FetchListFn = spdx.FetchLicenseList

// FetchDetailsFn is the details fetcher; defaults to spdx.FetchLicenseDetails. Override in tests.
var FetchDetailsFn = spdx.FetchLicenseDetails

const defaultTTL = 86400 // 24h in seconds

// TTL returns the effective cache TTL in seconds. If cfg is nil, use default. If *cfg is 0, always fetch.
func TTL(cfg *int) int {
	if cfg == nil {
		return defaultTTL
	}
	return *cfg
}

// FetchList returns the SPDX license list, from cache if valid (mtime within ttl) or via spdx.FetchLicenseList.
// cacheDir is ~/.ligma/_cache. ttl 0: always fetch. On cache write failure, still returns fetched data.
func FetchList(ctx context.Context, cacheDir string, ttl int, listURL string) ([]spdx.License, error) {
	listPath := filepath.Join(cacheDir, "list.json")

	if ttl == 0 {
		list, err := FetchListFn(ctx, listURL)
		tryWriteList(cacheDir, listPath, list)
		return list, err
	}

	fi, err := os.Stat(listPath)
	if err == nil && time.Since(fi.ModTime()) < time.Duration(ttl)*time.Second {
		return readListFile(listPath)
	}

	list, err := FetchListFn(ctx, listURL)
	tryWriteList(cacheDir, listPath, list)
	return list, err
}

type listFile struct {
	Licenses []spdx.License `json:"licenses"`
}

func readListFile(path string) ([]spdx.License, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var f listFile
	if err := json.Unmarshal(b, &f); err != nil {
		return nil, err
	}
	return f.Licenses, nil
}

func tryWriteList(cacheDir, listPath string, list []spdx.License) {
	_ = os.MkdirAll(cacheDir, 0755)
	b, err := json.Marshal(listFile{Licenses: list})
	if err != nil {
		return
	}
	_ = os.WriteFile(listPath, b, 0644)
}

// FetchDetails returns the license text for id, from cache if valid or via spdx.FetchLicenseDetails.
// cacheDir is ~/.ligma/_cache. ttlSec 0: always fetch. On cache write failure, still returns fetched data.
// ID is used as-is for the path (e.g. details/MIT.json); if id contains ".." or path separators, cache is skipped.
func FetchDetails(ctx context.Context, cacheDir string, ttl int, template, id string) (string, error) {
	if strings.Contains(id, "..") || strings.ContainsAny(id, `/\`) {
		return FetchDetailsFn(ctx, template, id)
	}

	detailsPath := filepath.Join(cacheDir, "details", id+".json")

	if ttl == 0 {
		text, err := FetchDetailsFn(ctx, template, id)
		tryWriteDetails(cacheDir, detailsPath, text)
		return text, err
	}

	fi, err := os.Stat(detailsPath)
	if err == nil && time.Since(fi.ModTime()) < time.Duration(ttl)*time.Second {
		return readDetailsFile(detailsPath)
	}

	text, err := FetchDetailsFn(ctx, template, id)
	tryWriteDetails(cacheDir, detailsPath, text)
	return text, err
}

type detailsFile struct {
	LicenseText string `json:"licenseText"`
}

func readDetailsFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var f detailsFile
	if err := json.Unmarshal(b, &f); err != nil {
		return "", err
	}
	return f.LicenseText, nil
}

func tryWriteDetails(cacheDir, detailsPath, text string) {
	dir := filepath.Dir(detailsPath)
	_ = os.MkdirAll(dir, 0755)
	b, err := json.Marshal(detailsFile{LicenseText: text})
	if err != nil {
		return
	}
	_ = os.WriteFile(detailsPath, b, 0644)
}
