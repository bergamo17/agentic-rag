package websearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(apikey string) *Client {
	return &Client{APIKey: apikey, HTTPClient: &http.Client{}}
}

type SearchResult struct {
	Title   string `json:"title"`
	Url     string `json:"url"`
	Content string `json:"content"`
}

type searchRequest struct {
	APIKey     string `json:"api_key"`
	Query      string `json:"query"`
	MaxResults int    `json:"max_results"`
}

type searchResponse struct {
	Result []SearchResult `json:"results"`
}

func (c *Client) Search(query string) ([]SearchResult, error) {
	payload, _ := json.Marshal(searchRequest{
		APIKey:     c.APIKey,
		Query:      query,
		MaxResults: 5,
	})

	resp, err := c.HTTPClient.Post("https://api.tavily.com/search", "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("web_search error (%d): %s", resp.StatusCode, string(body))
	}

	var result searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Result, nil
}
