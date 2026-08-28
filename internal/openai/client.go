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

type FunctionSpec struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type Tool struct {
	Type     string       `json:"type"`
	Function FunctionSpec `json:"function"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type Message struct {
	Role       string      `json:"role"`
	Content    interface{} `json:"content,omitempty"`
	ToolCall   []ToolCall  `json:"tool_calls,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
}

type chatRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	Tools     []Tool    `json:"tools,omitempty"`
	MaxTokens int       `json:"max_tokens"`
}

type chatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

type PageContext struct {
	DocumentID  string
	PageNumber  int
	Title       string
	ImageBase64 string
}

func (c *Client) ChatCompletion(message []Message, tools []Tool) (Message, error) {
	reqBody := chatRequest{
		Model:     modelName,
		Messages:  message,
		Tools:     tools,
		MaxTokens: 1000,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return Message{}, err
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(payload))
	if err != nil {
		return Message{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Message{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Message{}, fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, string(body))
	}

	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Message{}, fmt.Errorf("Empty response from OpenAI API")
	}

	if len(result.Choices) == 0 {
		return Message{}, fmt.Errorf("There is no choices in response")
	}

	return result.Choices[0].Message, nil
}

const systemPrompt = `You are an expert professional PDF analyst who gives rigorous in-depth answers. When relevant, mention which page number supports your claims.`

func (c *Client) QueryVLM(query string, pages []PageContext) (string, error) {
	content := make([]contentBlock, 0, len(pages)+1)
	for _, p := range pages {
		content = append(content, contentBlock{
			Type: "text",
			Text: fmt.Sprintf("The following image is Page %d of \"%s\":", p.PageNumber, p.Title),
		})
		content = append(content, contentBlock{
			Type:     "image_url",
			ImageURL: &imageURL{URL: "data:image/png;base64," + p.ImageBase64},
		})
	}
	content = append(content, contentBlock{Type: "text", Text: query})

	reqBody := chatRequest{
		Model: modelName,
		Messages: []Message{
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

	answer, ok := result.Choices[0].Message.Content.(string)
	if !ok {
		return "", fmt.Errorf("Response is not a text: %v", result.Choices[0].Message.Content)
	}

	return answer, nil
}
