package spdx

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const goodJSON = `{"licenseListVersion":"3.0","licenses":[{"licenseId":"MIT","name":"MIT License"},{"licenseId":"Apache-2.0","name":"Apache License 2.0"}]}`

func TestFetchLicenseList_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(goodJSON))
	}))
	defer srv.Close()

	ctx := context.Background()
	list, err := FetchLicenseList(ctx, srv.URL)
	if err != nil {
		t.Fatalf("FetchLicenseList: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("len(list) = %d, want 2", len(list))
	}
	if list[0].LicenseID != "MIT" || list[0].Name != "MIT License" {
		t.Errorf("list[0] = %+v", list[0])
	}
	if list[1].LicenseID != "Apache-2.0" || list[1].Name != "Apache License 2.0" {
		t.Errorf("list[1] = %+v", list[1])
	}
}

func TestFetchLicenseList_Non2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	ctx := context.Background()
	_, err := FetchLicenseList(ctx, srv.URL)
	if err == nil {
		t.Fatal("FetchLicenseList: expected error for 404")
	}
	if err.Error() != "spdx: list fetch failed: HTTP 404 Not Found" {
		t.Errorf("error = %v", err)
	}
}

func TestFetchLicenseList_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{invalid json`))
	}))
	defer srv.Close()

	ctx := context.Background()
	_, err := FetchLicenseList(ctx, srv.URL)
	if err == nil {
		t.Fatal("FetchLicenseList: expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "invalid JSON") {
		t.Errorf("error = %v", err)
	}
}

func TestFetchLicenseList_5xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	ctx := context.Background()
	_, err := FetchLicenseList(ctx, srv.URL)
	if err == nil {
		t.Fatal("FetchLicenseList: expected error for 500")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error = %v", err)
	}
}

func TestFetchLicenseList_ContextCanceled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := FetchLicenseList(ctx, srv.URL)
	if err == nil {
		t.Fatal("FetchLicenseList: expected error when context canceled")
	}
}

// --- FetchLicenseDetails ---

const goodDetailsJSON = `{"licenseId":"MIT","name":"MIT License","licenseText":"Copyright (c) <year> <copyright holders>\n\nPermission is hereby granted..."}`

func TestFetchLicenseDetails_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(goodDetailsJSON))
	}))
	defer srv.Close()

	template := srv.URL + "/{id}.json"
	ctx := context.Background()
	text, err := FetchLicenseDetails(ctx, template, "MIT")
	if err != nil {
		t.Fatalf("FetchLicenseDetails: %v", err)
	}
	if !strings.Contains(text, "Permission is hereby granted") {
		t.Errorf("licenseText = %q", text)
	}
}

func TestFetchLicenseDetails_404_ReturnsErrNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	template := srv.URL + "/{id}.json"
	ctx := context.Background()
	_, err := FetchLicenseDetails(ctx, template, "NoSuch")
	if err == nil {
		t.Fatal("FetchLicenseDetails: expected error for 404")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestFetchLicenseDetails_5xx_ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	template := srv.URL + "/{id}.json"
	ctx := context.Background()
	_, err := FetchLicenseDetails(ctx, template, "MIT")
	if err == nil {
		t.Fatal("FetchLicenseDetails: expected error for 5xx")
	}
	if errors.Is(err, ErrNotFound) {
		t.Errorf("5xx must not be ErrNotFound: %v", err)
	}
}

func TestFetchLicenseDetails_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{invalid`))
	}))
	defer srv.Close()

	template := srv.URL + "/{id}.json"
	ctx := context.Background()
	_, err := FetchLicenseDetails(ctx, template, "MIT")
	if err == nil {
		t.Fatal("FetchLicenseDetails: expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "invalid JSON") && !strings.Contains(err.Error(), "JSON") {
		t.Errorf("error = %v", err)
	}
}

func TestFetchLicenseDetails_MissingLicenseText(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"licenseId":"MIT","name":"MIT"}`))
	}))
	defer srv.Close()

	template := srv.URL + "/{id}.json"
	ctx := context.Background()
	_, err := FetchLicenseDetails(ctx, template, "MIT")
	if err == nil {
		t.Fatal("FetchLicenseDetails: expected error for missing licenseText")
	}
}

func TestFetchLicenseDetails_ContextCanceled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	template := srv.URL + "/{id}.json"
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := FetchLicenseDetails(ctx, template, "MIT")
	if err == nil {
		t.Fatal("FetchLicenseDetails: expected error when context canceled")
	}
}
