package agent

import (
	"github.com/bergamo17/agentic-rag-prototype/internal/openai"
)

func AvailableTools(dontAllowCustomCode bool) []openai.Tool {
	tools := []openai.Tool{
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
		{
			Type: "function",
			Function: openai.FunctionSpec{
				Name:        "generate_widget",
				Description: "Generate a visual widget (chart, table, or card) to display data on the dashboard, when the user explicitly asks for the visualization or when the structured data would be clear shown visually than as a text",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"widget_type": map[string]any{
							"type":        "string",
							"enum":        []string{"chart", "table", "card"},
							"description": "The type of widget to render on the dashboard",
						},
						"title": map[string]any{
							"type":        "string",
							"description": "Short title describing what this widget shows",
						},
						"data": map[string]any{
							"type":        "string",
							"description": "JSON-encoded string containing the widget's data, structured according to the widget_type. For 'chart': {\"chartType\": \"bar\"|\"line\", \"labels\": [...], \"values\": [...]}. For 'table': {\"columns\": [...], \"rows\": [[...], ...]}. For 'card': {\"value\": \"...\", \"description\": \"...\"}.",
						},
					},
					"required": []string{"widget_type", "title", "data"},
				},
			},
		},
		{
			Type: "function",
			Function: openai.FunctionSpec{
				Name:        "create_docx",
				Description: "Create a Word (.docx) document for a daily report, internal memo, or proposal. Use this tool ONLY if the request matches one of these three document categories, using a pre-approved theme.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"theme": map[string]any{
							"type": "string",
							"enum": []string{
								"formal-report-blue", "formal-report-green", "formal-report-grey",
								"internal-memo-grey", "internal-memo-blue", "proposal-red", "proposal-blue",
							},
							"description": "The document's visual theme -- a combination of structure (formal-report/internal-memo/proposal) and pre-approve accent color.",
						},
						"title": map[string]any{
							"type":        "string",
							"description": "The document's title",
						},
						"sections": map[string]any{
							"type":        "array",
							"description": "List of content sections, in the order they should appear.",
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"type": map[string]any{
										"type":        "string",
										"enum":        []string{"heading", "paragraph", "bullet_list", "table"},
										"description": "The section type: 'heading' for a section title (use level 1 for the main title, level 2/3 for sub-sections), 'paragraph' for narrative text, 'bullet_list' for a list of points, 'table' for tabular/comparison data.",
									},
									"level": map[string]any{
										"type":        "integer",
										"description": "Heading level (1-3). Only set when type='heading'.",
									},
									"text": map[string]any{
										"type":        "string",
										"description": "The text content. Set when type='heading' or type='paragraph'",
									},
									"items": map[string]any{
										"type":        "array",
										"items":       map[string]any{"type": "string"},
										"description": "The list of bullet points. Set when type='bullet_list'",
									},
									"headers": map[string]any{
										"type":        "array",
										"items":       map[string]any{"type": "string"},
										"description": "The table's column headers. Set when type='table'.",
									},
									"rows": map[string]any{
										"type":        "array",
										"items":       map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
										"description": "The table's row data. Set when type='table'.",
									},
								},
								"required": []string{"type"},
							},
						},
					},
					"required": []string{"theme", "title", "sections"},
				},
			},
		},
	}

	if !dontAllowCustomCode {
		tools = append(tools, executePythonSchema())
	}

	return tools
}

func executePythonSchema() openai.Tool {
	return openai.Tool{
		Type: "function",
		Function: openai.FunctionSpec{
			Name:        "execute_python",
			Description: "Write and execute python code to build a custom .docx document using python-docx, for request that don't fit the standard daily-report/internal-memo/proposal templates. The script must save its output to '/workspace/output/output.docx'.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"code": map[string]any{
						"type":        "string",
						"description": "Complete python script using the python-docx library. Must import docx, build the Document object, and call doc.save('/workspace/output/output.docx').",
					},
					"title": map[string]any{
						"type":        "string",
						"description": "The generated document's title from the execute_python tool",
					},
				},
				"required": []string{"code"},
			},
		},
	}
}
