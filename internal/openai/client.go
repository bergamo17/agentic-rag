package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	apiURL    = "https://api.openai.com/v1/chat/completions"
	modelName = "gpt-4o"
)

type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		APIKey:     os.Getenv("OPENAI_API_KEY"),
		HTTPClient: &http.Client{},
	}
}

type imageURL struct {
	URL string `json:"url"`
}

type contentBlock struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *imageURL `json:"image_url,omitempty"`
}

type message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type chatRequest struct {
	Model     string    `json:"model"`
	Messages  []message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

const systemPrompt = "You are an expert professional PDF analyst who gives rigorous in-depth answers."

func (c *Client) QueryVLM(query string, imageBase64 []string) (string, error) {
	content := make([]contentBlock, 0, len(imageBase64)+1)
	for _, img := range imageBase64 {
		content = append(content, contentBlock{
			Type:     "image_url",
			ImageURL: &imageURL{URL: "data:image/png;base64," + img},
		})
	}
	content = append(content, contentBlock{Type: "text", Text: query})

	reqBody := chatRequest{
		Model: modelName,
		Messages: []message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: content},
		},
		MaxTokens: 1000,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader((payload)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, string(body))
	}

	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("response kosong dari OpenAI API")
	}

	return result.Choices[0].Message.Content, nil
}
