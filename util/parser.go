package util

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgtype"
)

const maxTextRunes = 40000
const maxDocxXMLBytes = 20 << 20

func ParseUUID(s string) (pgtype.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: id, Valid: true}, nil
}

func UuidString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return uuid.UUID(u.Bytes).String()
}

func AppendTextFile(sb *strings.Builder, name, content string) {
	if r := []rune(content); len(r) > maxTextRunes {
		content = string(r[:maxTextRunes]) + "\n...[dipotong]"
	}
	fmt.Fprintf(sb, "\n\n[File: %s]\n%s", name, content)
}

func ExtractDocxText(data []byte) (string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}
	for _, zf := range zr.File {
		if zf.Name != "word/document.xml" {
			continue
		}
		rc, err := zf.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()
		return parseDocxXML(io.LimitReader(rc, maxDocxXMLBytes))
	}
	return "", errors.New("word/document.xml tidak ditemukan")
}

func parseDocxXML(r io.Reader) (string, error) {
	dec := xml.NewDecoder(r)
	var sb strings.Builder
	var inText, inTabs bool
	cellDepth := 0

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			if sb.Len() > 0 { // XML terpotong oleh batas ukuran: pakai yang sudah terbaca
				break
			}
			return "", err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "t":
				inText = true
			case "tabs":
				inTabs = true // definisi tab stop, bukan karakter tab
			case "tab":
				if !inTabs {
					sb.WriteByte('\t')
				}
			case "br", "cr":
				sb.WriteByte('\n')
			case "tc":
				cellDepth++
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "t":
				inText = false
			case "tabs":
				inTabs = false
			case "p":
				if cellDepth > 0 {
					sb.WriteByte(' ')
				} else {
					sb.WriteByte('\n')
				}
			case "tc":
				cellDepth--
				sb.WriteString(" | ")
			case "tr":
				sb.WriteByte('\n')
			}
		case xml.CharData:
			if inText {
				sb.Write(t)
			}
		}
	}

	return strings.TrimSpace(sb.String()), nil
}
