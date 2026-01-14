package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georgejdanforth/rckv/kv"
)

func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	kvStore := kv.NewMemoryStore()
	server := NewServer(kvStore)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /get", server.HandleGet)
	mux.HandleFunc("PUT /set", server.HandleSet)
	return httptest.NewServer(mux)
}

func TestServer(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	t.Run("set and get value", func(t *testing.T) {
		req, err := http.NewRequest("PUT", ts.URL+"/set?foo=bar", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to set: %v", err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}

		resp, err = http.Get(ts.URL + "/get?key=foo")
		if err != nil {
			t.Fatalf("failed to get: %v", err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if string(body) != "bar" {
			t.Fatalf("expected 'bar', got '%s'", body)
		}
	})

	t.Run("get key not found", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/get?key=nonexistent")
		if err != nil {
			t.Fatalf("failed to get: %v", err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", resp.StatusCode)
		}
	})

	t.Run("get missing key param", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/get")
		if err != nil {
			t.Fatalf("failed to get: %v", err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("get multiple key values", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/get?key=a&key=b")
		if err != nil {
			t.Fatalf("failed to get: %v", err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("set multiple keys", func(t *testing.T) {
		req, err := http.NewRequest("PUT", ts.URL+"/set?foo=bar&baz=qux", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to set: %v", err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("set multiple values for key", func(t *testing.T) {
		req, err := http.NewRequest("PUT", ts.URL+"/set?foo=bar&foo=baz", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to set: %v", err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})
}
