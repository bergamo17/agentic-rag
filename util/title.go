package util

import "strings"

const maxTitleRunes = 50

func MakeTitle(q string) string {
	q = strings.Join(strings.Fields(q), " ")
	r := []rune(q)
	if len(r) > maxTitleRunes {
		return string(r[:maxTitleRunes]) + "…"
	}
	return q
}
