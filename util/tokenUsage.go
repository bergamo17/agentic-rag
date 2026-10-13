package util

type Usage struct {
	TotalTokens int `json:"total_tokens"`
}

// func Total(u Usage) int {
// 	return u.PromptTokens + u.CompletionTokens
// }

// func Add(u *Usage, o Usage) {
// 	u.PromptTokens += o.PromptTokens
// 	u.CompletionTokens += o.CompletionTokens
// }
