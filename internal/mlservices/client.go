package mlservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL, HTTPClient: &http.Client{}}
}

type EmbedResult struct {
	Title    string `json:"title"`
	NumPages int    `json:"num_pages"`
}

func (c *Client) EmbedDocument(filename string, fileBytes []byte) (*EmbedResult, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	if _, err := part.Write(fileBytes); err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, c.BaseURL+"/embed", &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Model-service error (%d): %s", resp.StatusCode, string(body))
	}

	var result EmbedResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type RetrievedPage struct {
	Title       string `json:"title"`
	PageNumber  int    `json:"page_number"`
	ImageBase64 string `json:"image_base64"`
}

type retrieveResponse struct {
	Pages []RetrievedPage `json:"pages"`
}

func (c *Client) Retrieve(query string, k int) ([]RetrievedPage, error) {
	payload, _ := json.Marshal(map[string]interface{}{"query": query, "k": k})

	resp, err := c.HTTPClient.Post(c.BaseURL+"/retrieve", "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("model-service error (%d): %s", resp.StatusCode, string(body))
	}

	var result retrieveResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Pages, nil
}
