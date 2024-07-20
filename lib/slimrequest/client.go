package slimrequest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	c := &Client{
		BaseURL: baseURL,
	}

	if httpClient != nil {
		c.httpClient = httpClient
	} else {
		c.httpClient = http.DefaultClient
	}

	return c
}

// top level function not method because type parameters
func Do[A any](c *Client, ctx context.Context, r Request) (A, error) {
	var Zero A

	reqjson, err := json.Marshal(r)
	if err != nil {
		return Zero, err
	}
	body := bytes.NewBuffer(reqjson)
	url, err := url.JoinPath(c.BaseURL, "/jsonrpc.js")
	if err != nil {
		return Zero, err
	}

	htreq, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return Zero, err
	}

	resp, err := c.httpClient.Do(htreq)
	if err != nil {
		return Zero, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Zero, errors.New("bad http status")
	}

	respjson, err := io.ReadAll(resp.Body)
	if err != nil {
		return Zero, err
	}

	var overallResult req[A]
	err = json.Unmarshal(respjson, &overallResult)
	if err != nil {
		return Zero, err
	}

	return overallResult.Result, nil
}

func (c *Client) ServerStatus(ctx context.Context) (ServerStatus, error) {
	req := NewRequest("0", []string{"serverstatus"})
	return Do[ServerStatus](c, ctx, req)
}

func (c *Client) ExtendedServerStatus(ctx context.Context) (ExtendedServerStatus, error) {
	req := NewRequest("0", []string{"serverstatus", "0"})
	return Do[ExtendedServerStatus](c, ctx, req)
}

func (c *Client) PlayerStatus(ctx context.Context, playerID string) (PlayerStatus, error) {
	req := NewRequest(playerID, []string{"status"})
	ps, err := Do[PlayerStatus](c, ctx, req)
	ps.PlayerID = playerID
	
	return ps, err
}

