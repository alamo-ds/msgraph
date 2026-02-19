package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/alamo-ds/msgraph/env"
	"github.com/s-hammon/p"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/time/rate"
)

const (
	DefaultBaseURL = "https://graph.microsoft.com/v1.0"
	DefaultAuthURL = "https://login.microsoftonline.com/"
	DefaultScopes  = "https://graph.microsoft.com/.default"

	DefaultTimeoutSeconds         = 5
	DefaultRequestsPerSecondLimit = 100
	DefaultBurst                  = DefaultRequestsPerSecondLimit * 2
)

type AzureADConfig struct {
	TenantID     string   `json:"tenant_id"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`
}

type Client struct {
	BaseURL  string
	TenantID string
	ClientID string
	c        *http.Client
	limiter  *rate.Limiter

	eTagCache map[string]string
	// NOTE: if we expect the eTag to change for a resource, then this can become
	// just a sync.Mutex.
	eTagMu sync.RWMutex
}

func NewClient(ctx context.Context, azureADCfg ...AzureADConfig) *Client {
	var cfg AzureADConfig
	if len(azureADCfg) != 0 {
		cfg = azureADCfg[0]
	} else {
		data := env.LoadConfigFile()
		json.Unmarshal(data, &cfg)
	}
	if len(cfg.Scopes) == 0 {
		cfg.Scopes = append(cfg.Scopes, DefaultScopes)
	}

	eTagCache := env.LoadCacheFile().ETags
	if eTagCache == nil {
		eTagCache = make(map[string]string)
		env.WriteCacheFile("eTag", eTagCache)
	}

	adCfg := &clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     DefaultAuthURL + p.Format("%s/oauth2/v2.0/token", cfg.TenantID),
		Scopes:       cfg.Scopes,
	}

	client := &Client{
		BaseURL:   DefaultBaseURL,
		TenantID:  cfg.TenantID,
		ClientID:  cfg.ClientID,
		c:         adCfg.Client(ctx),
		limiter:   rate.NewLimiter(rate.Limit(DefaultRequestsPerSecondLimit), DefaultBurst),
		eTagCache: eTagCache,
	}

	return client
}

func (c *Client) Close() {
	env.WriteCacheFile("eTags", c.eTagCache)
}

// MaxRequestsPerSecond will also set the burst value to 2x the requests value
func (c *Client) MaxRequestsPerSecond(requests int) *Client {
	c.limiter.SetLimit(rate.Limit(requests))
	c.limiter.SetBurst(requests * 2)
	return c
}

func (c *Client) getETag(key string) (string, bool) {
	c.eTagMu.RLock()
	defer c.eTagMu.RUnlock()

	eTag, ok := c.eTagCache[key]
	return eTag, ok
}

func (c *Client) putETag(key, val string) {
	c.eTagMu.Lock()
	defer c.eTagMu.Unlock()

	c.eTagCache[key] = val
	log.Println("added eTag to cache w/ value:", val)
}

type refreshETagErr struct {
	err any
}

func (err refreshETagErr) Error() string {
	return p.Format("error fetching eTag for resource: %v", err.err)
}

func newRefreshETagErr(err any) refreshETagErr {
	return refreshETagErr{err: err}
}

func (c *Client) refreshETag(ctx context.Context, path string) (string, error) {
	var m map[string]any

	resp, err := c.get(ctx, path, nil)
	if err != nil {
		return "", newRefreshETagErr(err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return "", newRefreshETagErr(err)
	}

	eTag, ok := m["@odata.etag"]
	if !ok {
		return "", newRefreshETagErr("eTag not in response body")
	}

	switch s, ok := eTag.(string); ok {
	default:
		// NOTE: I don't think this should ever panic
		panic("type assertion for eTag failed!")
	case false:
		return "", fmt.Errorf("type for eTag: expected string, got %T", s)
	case true:
		c.putETag(path, s)
		return s, nil
	}
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	if err := c.limiter.Wait(req.Context()); err != nil {
		return nil, fmt.Errorf("limiter.Wait: %v", err)
	}
	return c.c.Do(req)
}

func (c *Client) get(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, requestErr(resp)
	}

	return resp, nil
}

func (c *Client) patch(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	eTag, ok := c.getETag(path)
	if !ok {
		var err error
		if eTag, err = c.refreshETag(ctx, path); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("If-Match", eTag)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, requestErr(resp)
	}

	return resp, nil
}

func (c *Client) post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %v", err)
	}
	if resp.StatusCode != 201 {
		return nil, requestErr(resp)
	}

	return resp, nil
}

// Use this only if JoinPath will not throw an error
func joinPath(base string, elem ...string) string {
	u, _ := url.JoinPath(base, elem...)
	return u
}

func readForError(r io.Reader) string {
	data, err := io.ReadAll(r)
	if err != nil {
		return "couldn't read body"
	}

	return string(data)
}

func toBody(v any) io.Reader {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(v)
	return &buf
}

func makeReqErr(v any) error {
	return fmt.Errorf("couldn't create request: %v", v)
}

func requestErr(resp *http.Response) error {
	if resp.Body != nil {
		defer resp.Body.Close()
		return fmt.Errorf("request returned %d: %s", resp.StatusCode, readForError(resp.Body))
	}

	return fmt.Errorf("request returned %d", resp.StatusCode)
}

// NOTE: val must be a pointer to a map or struct
func get[T any](ctx context.Context, c *Client, path string, val T) error {
	resp, err := c.get(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(val); err != nil {
		return fmt.Errorf("error decoding response body: %v", err)
	}

	return nil
}
