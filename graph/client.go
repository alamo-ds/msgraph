package graph

import (
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/s-hammon/msgraph/auth"
	"github.com/s-hammon/msgraph/env"
	"github.com/s-hammon/p"
	"golang.org/x/oauth2"
)

const defaultBaseURL = "https://graph.microsoft.com/v1.0"

type Client struct {
	BaseURL string
	c       *http.Client
	adCfg   auth.AzureADConfig

	tokenCache map[string]*oauth2.Token
	mu         sync.Mutex
}

func NewClient(azureADCfg ...auth.AzureADConfig) *Client {
	cfg := env.LoadDotFile()
	if len(azureADCfg) != 0 {
		cfg = azureADCfg[0]
		env.WriteDotFile(cfg)
	}

	client := &Client{
		BaseURL:    defaultBaseURL,
		c:          &http.Client{},
		adCfg:      cfg,
		tokenCache: make(map[string]*oauth2.Token),
	}

	return client
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	token, ok := c.getToken()
	if !ok || time.Now().After(token.Expiry) {
		var err error
		token, err = c.refreshToken()
		if err != nil {
			return nil, err
		}
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	return c.c.Do(req)
}

func (c *Client) cacheKey() string {
	return p.Format("%s:%s", c.adCfg.TenantID, c.adCfg.ClientID)
}

func (c *Client) getToken() (*oauth2.Token, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	token, ok := c.tokenCache[c.cacheKey()]
	return token, ok
}

func (c *Client) putToken(token *oauth2.Token) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.tokenCache[c.cacheKey()] = token
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

// Use this only if JoinPath will not throw an error
func joinPath(base string, elem ...string) string {
	u, _ := url.JoinPath(base, elem...)
	return u
}
