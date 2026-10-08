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
)

const (
	maxFiles        = 5
	maxFileBytes    = 10 << 20
	maxRequestBytes = 50 << 20
	maxTextRunes    = 40000
)

var allowedImageMIME = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/gif":  true,
	"image/webp": true,
}

type attachmentResult struct {
	Text    string              // isi file teks, siap ditempel ke prompt
	Images  []openai.ImageInput // gambar untuk content blocks
	Skipped []string            // file yang tidak didukung / gagal dibaca
}

func processAttachments(files []*multipart.FileHeader) attachmentResult {
	var res attachmentResult
	var text strings.Builder

	for _, fh := range files {
		if fh.Size > maxFileBytes {
			res.Skipped = append(res.Skipped, fh.Filename+" (terlalu besar)")
			continue
		}

		f, err := fh.Open()
		if err != nil {
			res.Skipped = append(res.Skipped, fh.Filename+" (gagal dibuka)")
			continue
		}
		data, err := io.ReadAll(io.LimitReader(f, maxFileBytes+1))
		f.Close()
		if err != nil || len(data) > maxFileBytes {
			res.Skipped = append(res.Skipped, fh.Filename+" (gagal dibaca)")
			continue
		}

		// Jangan percaya Content-Type dari klien; deteksi dari isi file.
		mime := http.DetectContentType(data)
		ext := strings.ToLower(filepath.Ext(fh.Filename))

		switch {
		case allowedImageMIME[mime]:
			res.Images = append(res.Images, openai.ImageInput{
				MIME:   mime,
				Base64: base64.StdEncoding.EncodeToString(data),
			})

		case (ext == ".txt" || ext == ".md") && utf8.Valid(data):
			content := string(data)
			if r := []rune(content); len(r) > maxTextRunes {
				content = string(r[:maxTextRunes]) + "\n...[dipotong]"
			}
			fmt.Fprintf(&text, "\n\n[File: %s]\n%s", fh.Filename, content)

		default:
			res.Skipped = append(res.Skipped, fh.Filename+" (format belum didukung)")
		}
	}

	res.Text = text.String()
	return res
}
