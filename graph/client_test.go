package graph

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestClientGet(t *testing.T) {
	server := newTestServer(t, http.MethodGet, "/test", `{"key":"value"}`)
	defer server.Close()

	client := newClient(server)

	resp, err := client.get(context.Background(), server.URL+"/test", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestClientPost(t *testing.T) {
	server := newTestServer(t, http.MethodPost, "/test", `{"key":"value"}`)
	defer server.Close()

	client := newClient(server)

	resp, err := client.post(context.Background(), server.URL+"/test", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestClientPatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Mock refreshETag
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"@odata.etag":"W/\"test-etag\""}`))
			return
		}
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, "/test", r.URL.Path)
		require.Equal(t, "W/\"test-etag\"", r.Header.Get("If-Match"))

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
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
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
	require.NoError(t, err)
	require.Equal(t, "W/\"new-etag\"", etag)

	cachedEtag, ok := client.getETag(server.URL + "/test")
	require.True(t, ok)
	require.True(t, cachedEtag == "W/\"new-etag\"")
}

func TestNewClientAndClose(t *testing.T) {
	// Test with provided config
	cfg := AzureADConfig{
		TenantID:     "test-tenant",
		ClientID:     "test-client",
		ClientSecret: "test-secret",
	}
	client := NewClient(context.Background(), cfg)
	require.Equal(t, "test-tenant", client.TenantID)
	require.Equal(t, "test-client", client.ClientID)
	require.Equal(t, DefaultBaseURL, client.BaseURL)

	client.MaxRequestsPerSecond(50)
	require.Equal(t, rate.Limit(50), client.limiter.Limit())
	require.Equal(t, 100, client.limiter.Burst())

	// Test Close (should write cache file)
	client.Close()
}

func TestMarshalAADConfig(t *testing.T) {
	// NOTE: this tests to make sure the secret isn't serialized (see CWE-499)
	cfg := AzureADConfig{
		TenantID:     "test-tenant",
		ClientID:     "test-client",
		ClientSecret: "supersecret",
	}

	b, _ := json.Marshal(cfg)
	ret := AzureADConfig{}
	require.NoError(t, json.Unmarshal(b, &ret))
	require.Equal(t, ret.ClientSecret, "")
}
