package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// Client holds target Technitium API configuration details.
type Client struct {
	BaseURL    *url.URL
	Token      string
	HTTPClient *http.Client
}

// NewClient instantiates a new Technitium API Client.
func NewClient(baseURL, token string) *Client {
	u, _ := url.Parse(baseURL)
	return &Client{
		BaseURL: u,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// APIResponse represents the standard Technitium API JSON response wrapper.
type APIResponse struct {
	Status       string          `json:"status"`
	ErrorMessage string          `json:"errorMessage,omitempty"`
	Response     json.RawMessage `json:"response,omitempty"`
}

// do sends the HTTP request and handles standard Technitium API response validations.
func (c *Client) do(req *http.Request) ([]byte, error) {
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResp APIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %w, raw response: %s", err, string(bodyBytes))
	}

	if apiResp.Status != "ok" {
		msg := apiResp.ErrorMessage
		if msg == "" {
			msg = "unknown error"
		}
		return nil, fmt.Errorf("API error: %s (status: %s)", msg, apiResp.Status)
	}

	return bodyBytes, nil
}

// buildRequest constructs an http.Request with the target path, query params, and body.
func (c *Client) buildRequest(method, apiPath string, params url.Values, body io.Reader) (*http.Request, error) {
	if c.BaseURL == nil {
		return nil, errors.New("client base URL is not initialized")
	}
	u := *c.BaseURL
	u.Path = path.Join(u.Path, apiPath)
	if params != nil {
		u.RawQuery = params.Encode()
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	return req, nil
}

// AddReservation adds a DHCP lease reservation.
// Map parameters correctly per user rules:
// - 'address' (IP)
// - 'hardwareAddress' (MAC)
// - 'hostName' (Hostname)
// - 'comments' (Description/Comment)
func (c *Client) AddReservation(ip, mac, hostname, comment string) error {
	if ip == "" {
		return errors.New("IP address cannot be empty")
	}
	if mac == "" {
		return errors.New("MAC address cannot be empty")
	}
	if hostname == "" {
		return errors.New("hostname cannot be empty")
	}

	q := url.Values{}
	q.Set("address", ip)
	q.Set("hardwareAddress", mac)
	q.Set("hostName", hostname)
	if comment != "" {
		q.Set("comments", comment)
	}

	req, err := c.buildRequest(http.MethodGet, "/api/dhcp/scopes/addReservedLease", q, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req)
	return err
}

// RemoveReservation removes/deletes a DHCP lease reservation using /api/dhcp/scopes/removeReservedLease.
// It passes the colon-separated 'hardwareAddress' MAC address per user rules.
func (c *Client) RemoveReservation(mac string) error {
	if mac == "" {
		return errors.New("MAC address cannot be empty")
	}

	q := url.Values{}
	q.Set("hardwareAddress", mac)

	req, err := c.buildRequest(http.MethodGet, "/api/dhcp/scopes/removeReservedLease", q, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req)
	return err
}

// SetAppConfig pushes/configures app settings via /api/apps/config/set?name=Advanced Blocking.
func (c *Client) SetAppConfig(configJSON string) error {
	if configJSON == "" {
		return errors.New("configJSON cannot be empty")
	}

	q := url.Values{}
	q.Set("name", "Advanced Blocking")

	form := url.Values{}
	form.Set("config", configJSON)

	req, err := c.buildRequest(http.MethodPost, "/api/apps/config/set", q, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = c.do(req)
	return err
}

// GetAppConfig retrieves the current configuration of the app.
func (c *Client) GetAppConfig() (string, error) {
	q := url.Values{}
	q.Set("name", "Advanced Blocking")

	req, err := c.buildRequest(http.MethodGet, "/api/apps/config/get", q, nil)
	if err != nil {
		return "", err
	}

	bodyBytes, err := c.do(req)
	if err != nil {
		return "", err
	}

	// Response is: {"status": "ok", "response": {"config": "..."}}
	var apiResp struct {
		Response struct {
			Config interface{} `json:"config"`
		} `json:"response"`
	}
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse app config response: %w", err)
	}

	// If config is null, return empty string
	if apiResp.Response.Config == nil {
		return "", nil
	}

	// Technitium can return config as a JSON object directly or as a string.
	// Let's handle both string and object cases.
	switch v := apiResp.Response.Config.(type) {
	case string:
		return v, nil
	default:
		// Re-marshal to JSON string
		bytes, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("failed to marshal config object to string: %w", err)
		}
		return string(bytes), nil
	}
}

// FetchCurrentScope queries the active DHCP scopes/leases via /api/dhcp/scopes/list.
func (c *Client) FetchCurrentScope() ([]byte, error) {
	req, err := c.buildRequest(http.MethodGet, "/api/dhcp/scopes/list", nil, nil)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}
