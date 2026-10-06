package agent

import (
	"strings"

	"github.com/bergamo17/agentic-rag-prototype/internal/openai"
)

var customFormatKeywords = []string{
	"pdf", "docx", "xlsx", "markdown", "md", "spreadsheet", "docs", "excel",
}

func ClassifyRequest(userMessage string) bool {
	standardKeywords := []string{
		"formal report", "daily report",
		"internal memo", "memo internal",
		"proposal",
	}

	lower := strings.ToLower(userMessage)

	for _, kw := range customFormatKeywords {
		if strings.Contains(lower, kw) {
			return false
		}
	}

	for _, kw := range standardKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

func getLastUserMessage(message []openai.Message) string {
	for i := len(message) - 1; i >= 0; i-- {
		if message[i].Role == "user" {
			if content, ok := message[i].Content.(string); ok {
				return content
			}
		}
	}
	return ""
}
