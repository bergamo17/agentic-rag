package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bergamo17/agentic-rag-prototype/util"
)

const anthropicVersion = "2023-06-01"

type anthropicSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

type anthropicBlock struct {
	Type      string           `json:"type"`
	Text      string           `json:"text,omitempty"`
	Title     string           `json:"title,omitempty"`
	Source    *anthropicSource `json:"source,omitempty"`
	ID        string           `json:"id,omitempty"`
	Name      string           `json:"name,omitempty"`
	Input     json.RawMessage  `json:"input,omitempty"`
	ToolUseID string           `json:"tool_use_id,omitempty"`
	Content   string           `json:"content,omitempty"`
}

type anthropicMessage struct {
	Role    string           `json:"role"`
	Content []anthropicBlock `json:"content"`
}

type anthropicTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	InputSchema map[string]any `json:"input_schema"`
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []anthropicMessage `json:"messages"`
	Tools     []anthropicTool    `json:"tools,omitempty"`
}

type anthropicResponse struct {
	Content    []anthropicBlock `json:"content"`
	StopReason string           `json:"stop_reason"`
	Usage      util.Usage       `json:"usage"`
}

func toAnthropic(msgs []Message) (string, []anthropicMessage) {
	var systems []string
	var out []anthropicMessage

	add := func(role string, blocks ...anthropicBlock) {
		if len(blocks) == 0 {
			return
		}
		if n := len(out); n > 0 && out[n-1].Role == role {
			out[n-1].Content = append(out[n-1].Content, blocks...)
			return
		}
		out = append(out, anthropicMessage{Role: role, Content: blocks})
	}

	for _, m := range msgs {
		switch m.Role {
		case "system":
			if s := contentText(m.Content); s != "" {
				systems = append(systems, s)
			}
		case "tool":
			add("user", anthropicBlock{
				Type:      "tool_result",
				ToolUseID: m.ToolCallID,
				Content:   contentText(m.Content),
			})
		case "assistant":
			var blocks []anthropicBlock
			if t := contentText(m.Content); strings.TrimSpace(t) != "" {
				blocks = append(blocks, anthropicBlock{Type: "text", Text: t})
			}
			for _, tc := range m.ToolCall {
				args := strings.TrimSpace(tc.Function.Arguments)
				if args == "" || !json.Valid([]byte(args)) {
					args = "{}"
				}
				blocks = append(blocks, anthropicBlock{
					Type:  "tool_use",
					ID:    tc.ID,
					Name:  tc.Function.Name,
					Input: json.RawMessage(args),
				})
			}
			add("assistant", blocks...)
		default: // "user"
			add("user", userBlocks(m.Content)...)
		}
	}

	for len(out) > 0 && out[0].Role != "user" {
		out = out[1:]
	}

	return strings.Join(systems, "\n"), out
}

func (c *Client) chatAnthropic(ctx context.Context, messages []Message, tools []Tool) (message Message, usage util.Usage, err error) {
	system, msgs := toAnthropic(messages)
	if len(msgs) == 0 {
		return Message{}, util.Usage{}, fmt.Errorf("no user message to send to the LLM")
	}

	reqBody := anthropicRequest{
		Model:     modelName,
		MaxTokens: maxTokens,
		System:    system,
		Messages:  msgs,
	}

	for _, t := range tools {
		schema := t.Function.Parameters
		if schema == nil {
			schema = map[string]any{"type": "object", "properties": map[string]any{}}
		}
		reqBody.Tools = append(reqBody.Tools, anthropicTool{
			Name:        t.Function.Name,
			Description: t.Function.Description,
			InputSchema: schema,
		})
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return Message{}, util.Usage{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(payload))
	if err != nil {
		return Message{}, util.Usage{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", anthropicVersion)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Message{}, util.Usage{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Message{}, util.Usage{}, fmt.Errorf("Anthropic API error (%d): %s", resp.StatusCode, string(body))
	}

	var result anthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Message{}, util.Usage{TotalTokens: result.Usage.TotalTokens}, fmt.Errorf("failed to decode Anthropic response: %w", err)
	}

	message = Message{Role: "assistant"}
	usage = util.Usage{
		TotalTokens: result.Usage.TotalTokens,
	}
	var text strings.Builder
	for _, b := range result.Content {
		switch b.Type {
		case "text":
			text.WriteString(b.Text)
		case "tool_use":
			tc := ToolCall{ID: b.ID, Type: "function"}
			tc.Function.Name = b.Name
			tc.Function.Arguments = string(b.Input)
			message.ToolCall = append(message.ToolCall, tc)
		}
	}
	if text.Len() > 0 {
		message.Content = text.String()
	}
	if text.Len() == 0 && len(message.ToolCall) == 0 {
		return Message{}, usage, fmt.Errorf("empty response from Anthropic API (stop_reason=%q)", result.StopReason)
	}

	return message, usage, nil
}

func contentText(v interface{}) string {
	switch c := v.(type) {
	case nil:
		return ""
	case string:
		return c
	case []contentBlock:
		var parts []string
		for _, b := range c {
			if b.Type == "text" && b.Text != "" {
				parts = append(parts, b.Text)
			}
		}
		return strings.Join(parts, "\n")
	default:
		return fmt.Sprint(c)
	}
}

func userBlocks(v interface{}) []anthropicBlock {
	switch c := v.(type) {
	case string:
		if strings.TrimSpace(c) == "" {
			return nil
		}
		return []anthropicBlock{{Type: "text", Text: c}}
	case []contentBlock:
		var blocks []anthropicBlock
		for _, b := range c {
			switch {
			case b.Type == "image_url" && b.ImageURL != nil:
				if mediaType, data, ok := parseDataURI(b.ImageURL.URL); ok {
					blocks = append(blocks, anthropicBlock{
						Type:   "image",
						Source: &anthropicSource{Type: "base64", MediaType: mediaType, Data: data},
					})
				}
			case b.Doc != nil:
				blocks = append(blocks, anthropicBlock{
					Type:   "document",
					Title:  b.Doc.Name,
					Source: &anthropicSource{Type: "base64", MediaType: b.Doc.MIME, Data: b.Doc.Base64},
				})
			case b.Text != "":
				blocks = append(blocks, anthropicBlock{Type: "text", Text: b.Text})
			}
		}
		return blocks
	}
	return nil
}

func parseDataURI(uri string) (mediaType, data string, ok bool) {
	rest, found := strings.CutPrefix(uri, "data:")
	if !found {
		return "", "", false
	}
	meta, data, found := strings.Cut(rest, ",")
	if !found {
		return "", "", false
	}
	mediaType, _, _ = strings.Cut(meta, ";")
	if mediaType == "" {
		mediaType = "image/png"
	}
	return mediaType, data, true
}
