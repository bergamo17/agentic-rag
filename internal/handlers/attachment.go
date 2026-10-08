package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/bergamo17/agentic-rag-prototype/internal/openai"
	"github.com/bergamo17/agentic-rag-prototype/util"
)

const (
	maxFiles        = 5
	maxFileBytes    = 10 << 20
	maxRequestBytes = 50 << 20
	maxTextRunes    = 40000
	maxDocxXMLBytes = 20 << 20
)

var allowedImageMIME = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/gif":  true,
	"image/webp": true,
}

type attachmentResult struct {
	Text      string
	Images    []openai.ImageInput
	Documents []openai.DocumentInput
	Skipped   []string
}

func processAttachments(files []*multipart.FileHeader) attachmentResult {
	var res attachmentResult
	var text strings.Builder

	for _, fh := range files {
		if fh.Size > maxFileBytes {
			res.Skipped = append(res.Skipped, fh.Filename+" (too large)")
			continue
		}

		f, err := fh.Open()
		if err != nil {
			res.Skipped = append(res.Skipped, fh.Filename+" (failed to open)")
			continue
		}
		data, err := io.ReadAll(io.LimitReader(f, maxFileBytes+1))
		f.Close()
		if err != nil || len(data) > maxFileBytes {
			res.Skipped = append(res.Skipped, fh.Filename+" (failed to read)")
			continue
		}

		mime := http.DetectContentType(data)
		ext := strings.ToLower(filepath.Ext(fh.Filename))

		switch {
		case allowedImageMIME[mime]:
			res.Images = append(res.Images, openai.ImageInput{
				MIME:   mime,
				Base64: base64.StdEncoding.EncodeToString(data),
			})

		case mime == "application/pdf":
			if !openai.SupportsPDF() {
				res.Skipped = append(res.Skipped, fh.Filename+" (PDF dinonaktifkan: LLM_PDF_INPUT=false)")
				continue
			}
			res.Documents = append(res.Documents, openai.DocumentInput{
				Name:   fh.Filename,
				MIME:   mime,
				Base64: base64.StdEncoding.EncodeToString(data),
			})

		case ext == ".docx" && mime == "application/zip":
			content, err := util.ExtractDocxText(data)
			if err != nil || strings.TrimSpace(content) == "" {
				res.Skipped = append(res.Skipped, fh.Filename+" (docx kosong / tidak bisa dibaca)")
				continue
			}
			util.AppendTextFile(&text, fh.Filename, content)

		case (ext == ".txt" || ext == ".md") && utf8.Valid(data):
			content := string(data)
			if r := []rune(content); len(r) > maxTextRunes {
				content = string(r[:maxTextRunes]) + "\n...[dipotong]"
			}
			fmt.Fprintf(&text, "\n\n[File: %s]\n%s", fh.Filename, content)

		default:
			res.Skipped = append(res.Skipped, fh.Filename+" (not supporting this type of format)")
		}
	}

	res.Text = text.String()
	return res
}
