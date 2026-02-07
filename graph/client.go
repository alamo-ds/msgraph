package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/alamo-ds/msgraph/auth"
	"github.com/alamo-ds/msgraph/env"
	"github.com/s-hammon/p"
	"golang.org/x/oauth2"
)

const (
	DefaultBaseURL        = "https://graph.microsoft.com/v1.0"
	DefaultTimeoutSeconds = 5
)

type Client struct {
	BaseURL string
	c       *http.Client
	adCfg   auth.AzureADConfig

	tokenCache map[string]*oauth2.Token
	tokenMu    sync.Mutex

	eTagCache map[string]string
	// NOTE: if we expect the eTag to change for a resource,
	// NOTE: then this can become just a sync.Mutex.
	eTagMu sync.RWMutex
}

func NewClient(azureADCfg ...auth.AzureADConfig) *Client {
	cfg := env.LoadConfigFile()
	if len(azureADCfg) != 0 {
		cfg = azureADCfg[0]
		env.WriteConfigFile(cfg)
	}

	eTagCache := env.LoadCacheFile().ETags
	if eTagCache == nil {
		eTagCache = make(map[string]string)
		env.WriteCacheFile("eTag", eTagCache)
	}

	client := &Client{
		BaseURL: DefaultBaseURL,
		c: &http.Client{
			Timeout: DefaultTimeoutSeconds * time.Second,
		},
		adCfg:      cfg,
		tokenCache: make(map[string]*oauth2.Token),
		eTagCache:  eTagCache,
	}

	return client
}

func (c *Client) Close() {
	env.WriteCacheFile("eTags", c.eTagCache)
}

func (c *Client) cacheKey() string {
	return p.Format("%s:%s", c.adCfg.TenantID, c.adCfg.ClientID)
}

func (c *Client) getToken() (*oauth2.Token, bool) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	token, ok := c.tokenCache[c.cacheKey()]
	return token, ok
}

func (c *Client) putToken(token *oauth2.Token) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	c.tokenCache[c.cacheKey()] = token
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

func (c *Client) refreshToken() (*oauth2.Token, error) {
	resp, err := auth.GetAccessToken(c.adCfg)
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{
		AccessToken: resp.Token,
		TokenType:   resp.TokenType,
		Expiry:      time.Now().Add(time.Second * time.Duration(resp.ExpiresIn)),
	}
	if token.AccessToken == "" {
		return nil, errors.New("server response missing access_token")
	}

	c.putToken(token)
	return token, nil
}

func (c *Client) get(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := c.newAuthRequest(ctx, http.MethodGet, path, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.c.Do(req)
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

	req, err := c.newAuthRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("If-Match", eTag)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, requestErr(resp)
	}

	return resp, nil
}

func (c *Client) post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := c.newAuthRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %v", err)
	}
	if resp.StatusCode != 201 {
		return nil, requestErr(resp)
	}

	return resp, nil
}

func (c *Client) newAuthRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	token, ok := c.getToken()
	if !ok || time.Now().After(token.Expiry) {
		var err error
		token, err = c.refreshToken()
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return nil, makeReqErr(err)
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	return req, nil
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
