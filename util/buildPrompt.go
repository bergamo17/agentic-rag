package util

import (
	_ "embed"
	"html/template"
	"strings"
)

//go:embed prompts/AGENTS.md
var agentPromptTmpl string

type PromptData struct {
	CompanyName, Now, UserName, Department string
}

func BuildSystemPrompt(d PromptData) (string, error) {
	t, err := template.New("AGENTS").Parse(agentPromptTmpl)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	if err := t.Execute(&sb, d); err != nil {
		return "", err
	}
	return sb.String(), nil
}
