package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	mlservice "github.com/bergamo17/agentic-rag-prototype/internal/mlservices"
	"github.com/bergamo17/agentic-rag-prototype/internal/openai"
	"github.com/bergamo17/agentic-rag-prototype/internal/websearch"
)

func executeTool(mlClient *mlservice.Client, webClient *websearch.Client, tc openai.ToolCall) (string, []openai.PageContext, error) {
	switch tc.Function.Name {

	case "search_documents":
		var args struct {
			Query string `json:"query"`
		}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
			return "", nil, err
		}
		pages, err := mlClient.Retrieve(args.Query, 3)
		if err != nil {
			return "", nil, err
		}
		pageContext := make([]openai.PageContext, len(pages))
		for i, p := range pages {
			pageContext[i] = openai.PageContext{
				DocumentID:  p.DocumentID,
				PageNumber:  p.PageNumber,
				Title:       p.Title,
				ImageBase64: p.ImageBase64,
			}
		}
		return fmt.Sprintf("Founded %d relevant pages.", len(pages)), pageContext, nil

	case "web_search":
		var args struct {
			Query string `json:"query"`
		}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
			return "", nil, err
		}

		results, err := webClient.Search(args.Query)
		if err != nil {
			return "", nil, err
		}

		var summary strings.Builder
		summary.WriteString("Searching result:\n")
		for _, r := range results {
			summary.WriteString(fmt.Sprintf("- %s: %s\n  %s\n", r.Title, r.Url, r.Content))
		}

		return summary.String(), nil, nil

	case "get_page_image":
		var args struct {
			DocumentID string `json:"document_id"`
			PageNumber int    `json:"page_number"`
		}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
			return "", nil, err
		}

		page, err := mlClient.GetPage(args.DocumentID, args.PageNumber)
		if err != nil {
			return "", nil, err
		}

		pageContext := openai.PageContext{
			DocumentID:  page.DocumentID,
			PageNumber:  page.PageNumber,
			Title:       page.Title,
			ImageBase64: page.ImageBase64,
		}

		return fmt.Sprintf("Page %d from %s is retrieved successfully", page.PageNumber, page.Title), []openai.PageContext{pageContext}, nil

	case "list_documents":
		docs, err := mlClient.ListDocuments()
		if err != nil {
			return "", nil, err
		}

		if len(docs) == 0 {
			return "There is not existing documents.", nil, nil
		}

		summary := "Document available:\n"
		for _, d := range docs {
			summary += fmt.Sprintf("- %s (ID: %s, %d halaman)\n", d.Title, d.DocumentID, d.NumPages)
		}

		return summary, nil, nil

	default:
		return "", nil, fmt.Errorf("Unknown tool: %s", tc.Function.Name)
	}

}

func AgentLoop(
	opeaiClient *openai.Client,
	webClient *websearch.Client,
	mlClient *mlservice.Client,
	messages []openai.Message,
) (string, []openai.PageContext, bool, error) {
	const maxItterations = 5
	var usedPages []openai.PageContext

	for i := 0; i < maxItterations; i++ {
		log.Printf("Itteration- %d started", i)
		resp, err := opeaiClient.ChatCompletion(messages, AvailableTools())
		if err != nil {
			return "", nil, false, err
		}

		if len(resp.ToolCall) == 0 {
			answer, _ := resp.Content.(string)
			return answer, usedPages, false, nil
		}

		messages = append(messages, resp)

		for _, tc := range resp.ToolCall {
			log.Printf("Tool called: %s | Arguments: %s", tc.Function.Name, tc.Function.Arguments)
			result, pages, err := executeTool(mlClient, webClient, tc)
			if err != nil {
				result = fmt.Sprintf("Error executing tool: %s", err)
			}
			log.Printf("Called Tool result: %s", result)

			if len(pages) > 0 {
				usedPages = append(usedPages, pages...)
			}

			messages = append(messages, openai.Message{
				Role:       "tool",
				Content:    result,
				ToolCallID: tc.ID,
			})
		}
	}

	messages = append(messages, openai.Message{
		Role:     "system",
		Content:  "You have reached the maximum limit for search steps. Do not call the tool again. Based on all the information you have gathered so far, provide the best answer you can formulate now, and indicate which parts may still be incomplete.",
		ToolCall: []openai.ToolCall{},
	})

	finalResp, err := opeaiClient.ChatCompletion(messages, []openai.Tool{})
	if err != nil {
		return "", nil, false, err
	}

	finalAnswer, ok := finalResp.Content.(string)
	if !ok {
		return "Sorry, unable to compile a summary of the answer.", usedPages, false, nil
	}

	return finalAnswer, usedPages, true, nil
}
