package spdx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// DefaultListURL is the official SPDX license list (MVP hardcoded). Epic 5.2 adds config override.
const DefaultListURL = "https://raw.githubusercontent.com/spdx/license-list-data/main/json/licenses.json"

// DefaultDetailsURLTemplate is the SPDX details URL with {id} placeholder. Epic 5.2 adds config override.
const DefaultDetailsURLTemplate = "https://raw.githubusercontent.com/spdx/license-list-data/main/json/details/{id}.json"

// ErrNotFound is returned when the license details URL returns HTTP 404. get/write map this to exit 2.
var ErrNotFound = errors.New("spdx: not found")

// License holds minimal fields from SPDX licenses.json for the list. Use as-is; no normalization (project-context).
type License struct {
	LicenseID string `json:"licenseId"`
	Name      string `json:"name"`
}

type listResponse struct {
	Licenses []License `json:"licenses"`
}

// FetchLicenseList GETs listURL, parses JSON, and returns the licenses. Returns errors only; no os.Exit (internal/).
// Uses 30s timeout (NFR-I2). On non-2xx, network error, or timeout: returns a descriptive error.
// On 4xx/5xx the body is not parsed as JSON.
func FetchLicenseList(ctx context.Context, listURL string) ([]License, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, listURL, nil)
	if err != nil {
		return nil, fmt.Errorf("spdx: new request: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("spdx: fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, fmt.Errorf("spdx: list fetch failed: HTTP %s", resp.Status)
	}

	var list listResponse
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("spdx: invalid JSON: %w", err)
	}
	return list.Licenses, nil
}

type detailsResponse struct {
	LicenseText string `json:"licenseText"`
}

// FetchLicenseDetails GETs the details URL (template with {id} replaced by id as-is), parses JSON,
// and returns the licenseText. Uses 30s timeout (NFR-I2). On 404 returns ErrNotFound (get→exit 2);
// on other 4xx/5xx, network, timeout, invalid JSON, or missing licenseText returns an error (get→exit 3).
// No os.Exit in internal/.
func FetchLicenseDetails(ctx context.Context, detailsURLTemplate, id string) (string, error) {
	url := strings.ReplaceAll(detailsURLTemplate, "{id}", id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("spdx: new request: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("spdx: fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		_, _ = io.Copy(io.Discard, resp.Body)
		return "", ErrNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, resp.Body)
		return "", fmt.Errorf("spdx: details fetch failed: HTTP %s", resp.Status)
	}

	var d detailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return "", fmt.Errorf("spdx: invalid JSON: %w", err)
	}
	if d.LicenseText == "" {
		return "", fmt.Errorf("spdx: missing licenseText")
	}
	return d.LicenseText, nil
}
