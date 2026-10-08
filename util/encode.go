package util

import (
	"encoding/json"
	"log"
)

func ToJSONB(v any, n int) []byte {
	if n == 0 {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return b
}

func DecodeJSONB[T any](raw []byte) []T {
	out := []T{}
	if len(raw) == 0 {
		return out
	}
	if err := json.Unmarshal(raw, &out); err != nil || out == nil {
		if err != nil {
			log.Printf("failed to decode jsonb: %v", err)
		}
		return []T{}
	}
	return out
}
