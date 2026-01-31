package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	Headers    map[string]string
}

func NewClient(baseURL string) *Client {
	return &Client{
		HTTPClient: &http.Client{},
		BaseURL:    baseURL,
		Headers:    make(map[string]string),
	}
}

func (c *Client) SetHeader(key, value string) {
	c.Headers[key] = value
}

func (c *Client) SetAuthToken(token string) {
	c.Headers["Authorization"] = token
}

func (c *Client) Request(method, endpoint string, body interface{}) ([]byte, error) {
	url := c.BaseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return responseBody, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}
