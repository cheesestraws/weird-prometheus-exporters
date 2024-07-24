package truenas

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL string
	User    string
	Pass    string

	httpClient *http.Client
}

func NewClient(baseURL string, user string, pass string, httpClient *http.Client) *Client {
	c := &Client{
		BaseURL: baseURL,
		User:    user,
		Pass:    pass,
	}

	if httpClient != nil {
		c.httpClient = httpClient
	} else {
		c.httpClient = http.DefaultClient
	}

	return c
}

// top level function not method because type parameters
func BasicGet[A any](c *Client, ctx context.Context, endpoint string) (A, error) {
	var Zero A

	url, err := url.JoinPath(c.BaseURL, "/api/v2.0/", endpoint)
	if err != nil {
		return Zero, err
	}

	htreq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Zero, err
	}

	htreq.SetBasicAuth(c.User, c.Pass)

	resp, err := c.httpClient.Do(htreq)
	if err != nil {
		return Zero, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Zero, fmt.Errorf("bad http status %v", resp.StatusCode)
	}

	respjson, err := io.ReadAll(resp.Body)
	if err != nil {
		return Zero, err
	}

	var result A
	err = json.Unmarshal(respjson, &result)
	if err != nil {
		return Zero, err
	}

	return result, nil
}

func (c *Client) AlertList(ctx context.Context) ([]Alert, error) {
	return BasicGet[[]Alert](c, ctx, "/alert/list")
}

func (c *Client) Pools(ctx context.Context) ([]Pool, error) {
	return BasicGet[[]Pool](c, ctx, "/pool")
}

func (c *Client) CloudSyncs(ctx context.Context) ([]CloudSync, error) {
	return BasicGet[[]CloudSync](c, ctx, "/cloudsync")
}