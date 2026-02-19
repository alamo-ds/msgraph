package graph

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/time/rate"
)

func TestClientGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/test", `{"key":"value"}`)
	defer server.Close()

	client := newClient(server)

	resp, err := client.get(context.Background(), server.URL+"/test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestClientPost(t *testing.T) {
	server := newTestServer(t, http.MethodPost, "/test", `{"key":"value"}`)
	defer server.Close()

	client := newClient(server)

	resp, err := client.post(context.Background(), server.URL+"/test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}
}

func TestClientPatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Mock refreshETag
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"@odata.etag":"W/\"test-etag\""}`))
			return
		}
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		if r.URL.Path != "/test" {
			t.Errorf("Expected path /test, got %s", r.URL.Path)
		}
		if r.Header.Get("If-Match") != "W/\"test-etag\"" {
			t.Errorf("Expected If-Match header W/\"test-etag\", got %s", r.Header.Get("If-Match"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"key":"value"}`))
	}))
	defer server.Close()

	client := &Client{
		BaseURL:   server.URL,
		c:         server.Client(),
		limiter:   rate.NewLimiter(rate.Limit(100), 200),
		eTagCache: make(map[string]string),
	}

	resp, err := client.patch(context.Background(), server.URL+"/test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestRefreshETag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"@odata.etag":"W/\"new-etag\""}`))
	}))
	defer server.Close()

	client := &Client{
		BaseURL:   server.URL,
		c:         server.Client(),
		limiter:   rate.NewLimiter(rate.Limit(100), 200),
		eTagCache: make(map[string]string),
	}

	etag, err := client.refreshETag(context.Background(), server.URL+"/test")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if etag != "W/\"new-etag\"" {
		t.Errorf("Expected etag W/\"new-etag\", got %s", etag)
	}

	cachedEtag, ok := client.getETag(server.URL + "/test")
	if !ok || cachedEtag != "W/\"new-etag\"" {
		t.Errorf("Expected cached etag W/\"new-etag\", got %s", cachedEtag)
	}
}

func TestNewClientAndClose(t *testing.T) {
	// Test with provided config
	cfg := AzureADConfig{
		TenantID:     "test-tenant",
		ClientID:     "test-client",
		ClientSecret: "test-secret",
	}
	client := NewClient(context.Background(), cfg)
	if client.TenantID != "test-tenant" {
		t.Errorf("Expected tenant ID test-tenant, got %s", client.TenantID)
	}
	if client.ClientID != "test-client" {
		t.Errorf("Expected client ID test-client, got %s", client.ClientID)
	}
	if client.BaseURL != DefaultBaseURL {
		t.Errorf("Expected base URL %s, got %s", DefaultBaseURL, client.BaseURL)
	}

	client.MaxRequestsPerSecond(50)
	if client.limiter.Limit() != rate.Limit(50) {
		t.Errorf("Expected limit 50, got %v", client.limiter.Limit())
	}
	if client.limiter.Burst() != 100 {
		t.Errorf("Expected burst 100, got %v", client.limiter.Burst())
	}

	// Test Close (should write cache file)
	client.Close()
}
