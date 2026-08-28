package agent

import (
	"github.com/bergamo17/agentic-rag-prototype/internal/openai"
)

func AvailableTools() []openai.Tool {
	return []openai.Tool{
		{
			Type: "function",
			Function: openai.FunctionSpec{
				Name:        "search_documents",
				Description: "Find the information from the uploaded document. Use this function when there is a possibility that the answer or user query corresponds with the saved documents, not a general knowledge",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "Semantic search query to search within documents",
						},
					},
					"required": []string{"query"},
				},
			},
		},
		{
			Type: "function",
			Function: openai.FunctionSpec{
				Name:        "web_search",
				Description: "Find the information from the internet or other website. Use this function when there is no possibility that the answer or user query corresponds with the saved documents",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "Search query for web search",
						},
					},
					"required": []string{"query"},
				},
			},
		},
		{
			Type: "function",
			Function: openai.FunctionSpec{
				Name:        "get_page_image",
				Description: "Re-capture specific document pages for detailed visual analysis. Use this if `search_documents` results are unclear and you need to view the page directly.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"page_number": map[string]any{
							"type":        "integer",
							"description": "Page number to view",
						},
						"document_id": map[string]any{
							"type":        "string",
							"description": "Document ID (optional, fill this if you know the ID from the `search_documents` before).",
						},
					},
					"required": []string{"page_number"},
				},
			},
		},
		{
			Type: "function",
			Function: openai.FunctionSpec{
				Name:        "list_documents",
				Description: "Displays a list of uploaded and indexed documents (file name, number of pages). Use this when a user refers to a document non-specifically (e.g., 'the document I sent yesterday') and you need to know which documents are available before performing `search_documents`.",
				Parameters: map[string]any{
					"type":       "object",
					"properties": map[string]any{},
				},
			},
		},
	}
}
