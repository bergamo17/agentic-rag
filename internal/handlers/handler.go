package handlers

import (
	"context"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	db "github.com/bergamo17/agentic-rag-prototype/db/sqlc"
	"github.com/bergamo17/agentic-rag-prototype/internal/agent"
	"github.com/bergamo17/agentic-rag-prototype/internal/docbuilder"
	mlservice "github.com/bergamo17/agentic-rag-prototype/internal/mlservices"
	"github.com/bergamo17/agentic-rag-prototype/internal/openai"
	"github.com/bergamo17/agentic-rag-prototype/internal/websearch"
	"github.com/bergamo17/agentic-rag-prototype/util"
)

type Handlers struct {
	ML     *mlservice.Client
	OpenAI *openai.Client
	Web    *websearch.Client
	Doc    *docbuilder.Client
	Sand   *docbuilder.SandboxClient
	Q      db.Querier
}

type pageRef struct {
	DocumentID string `json:"document_id"`
	Title      string `json:"title"`
	PageNumber int    `json:"page_number"`
}

// const (
// 	historyLimit  = 20
// 	maxTitleRunes = 50
// )

const agentSystemPrompt = `Kamu adalah asisten AI yang membantu menjawab pertanyaan menggunakan dokumen yang tersedia, pencarian web, dan alat visualisasi.

ATURAN PENTING setelah memanggil generate_widget, create_docx, atau execute_python:
- JANGAN PERNAH menyertakan markdown image (![...](...)), link file, atau path file apapun di jawaban teks kamu, walau kamu tahu nama filenya.
- Hasil widget/dokumen dirender otomatis dan terpisah oleh antarmuka -- kamu tidak perlu dan tidak boleh mereferensikannya dengan syntax markdown apapun.
- Cukup rujuk secara natural, misalnya "Berikut tabelnya:" atau "Saya sudah buatkan grafiknya." tanpa embed apapun.
- Jangan pernah mengarang nama file atau path yang tidak muncul di tool result.
- Gunakan execute_python dengan output_format yang sesuai permintaan user: "pdf" untuk laporan siap cetak, "xlsx" untuk data tabular/hitungan, "md" untuk catatan/dokumentasi teks, "docx" untuk dokumen Word kustom.
- Jika user meminta lebih dari satu format, panggil execute_python sekali per format.

Gunakan search_documents untuk pertanyaan yang mungkin terjawab dari dokumen yang diunggah, web_search untuk informasi umum/terkini, dan generate_widget saat data akan lebih jelas ditampilkan secara visual.`

func New(mlClient *mlservice.Client, openaiClient *openai.Client, webClient *websearch.Client, docClient *docbuilder.Client, sandboxClient *docbuilder.SandboxClient, queries db.Querier) *Handlers {
	return &Handlers{
		ML:     mlClient,
		OpenAI: openaiClient,
		Web:    webClient,
		Doc:    docClient,
		Sand:   sandboxClient,
		Q:      queries,
	}
}

func (h *Handlers) EmbedDocument(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file tidak ditemukan"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := h.ML.EmbedDocument(file.Filename, fileBytes)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

type ChatRequest struct {
	Query          string                  `json:"query" form:"query"`
	K              int                     `json:"k" form:"k"`
	ConversationID string                  `json:"conversation_id" form:"conversation_id"`
	Files          []*multipart.FileHeader `json:"-" form:"files"`
}

func (h *Handlers) Chat(c *gin.Context) {
	var req ChatRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.K == 0 {
		req.K = 3
	}

	pages, err := h.ML.Retrieve(req.Query, req.K)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if len(pages) == 0 {
		c.JSON(http.StatusOK, gin.H{"answer": "There is no document in the index", "pages": []any{}})
		return
	}

	pageContext := make([]openai.PageContext, len(pages))
	pageMeta := make([]gin.H, len(pages))
	for i, p := range pages {
		pageContext[i] = openai.PageContext{
			PageNumber:  p.PageNumber,
			Title:       p.Title,
			ImageBase64: p.ImageBase64,
		}
		pageMeta[i] = gin.H{
			"title":       p.Title,
			"page_number": p.PageNumber,
			"page_image":  p.ImageBase64,
		}
	}

	answer, err := h.OpenAI.QueryVLM(req.Query, pageContext)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"answer": answer,
		"pages":  pageMeta,
	})
}

func (h *Handlers) ChatAgent(c *gin.Context) {
	var req ChatRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var conv db.Conversation

	if req.ConversationID == "" {
		conversation, err := h.Q.CreateConversation(c, util.MakeTitle(req.Query))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create conversation."})
			return
		}
		conv = conversation
	} else {
		convId, err := util.ParseUUID(req.ConversationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id."})
			return
		}

		existing, err := h.Q.GetConversation(c, convId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get chat history."})
			return
		}

		conv = existing
	}

	usrMessage := db.InsertMessageParams{
		ConversationID: conv.ID,
		Role:           "user",
		Content:        req.Query,
	}

	userMessage, err := h.Q.InsertMessage(c, usrMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the message/chat."})
		return
	}

	listMessage := db.ListRecentMessagesParams{
		ConversationID: conv.ID,
		Limit:          20,
	}

	recentMessages, err := h.Q.ListRecentMessages(c, listMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the chat history on this conversation."})
		return
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	sys, err := util.BuildSystemPrompt(util.PromptData{
		CompanyName: "Ascendiz",
		Now:         time.Now().In(loc).Format("Monday, 2 January 2006 15:04"),
		UserName:    "Ariel",
		Department:  "Technology",
	})

	messages := make([]openai.Message, 0, len(recentMessages)+1)
	messages = append(messages, openai.Message{Role: "system", Content: sys})
	for i := len(recentMessages) - 1; i >= 0; i-- {
		messages = append(messages, openai.Message{
			Role:    recentMessages[i].Role,
			Content: recentMessages[i].Content,
		})
	}

	answer, pages, widgets, docs, isPartial, err := agent.AgentLoop(h.OpenAI, h.Web, h.ML, h.Doc, h.Sand, c, messages)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	pageMeta := make([]gin.H, len(pages))
	pageRefs := make([]pageRef, len(pages))
	for i, p := range pages {
		pageMeta[i] = gin.H{
			"document_id": p.DocumentID,
			"title":       p.Title,
			"page_number": p.PageNumber,
			"page_image":  p.ImageBase64,
		}
		pageRefs[i] = pageRef{DocumentID: p.DocumentID, Title: p.Title, PageNumber: p.PageNumber}
	}

	// widgetMeta := make([]gin.H, len(widgets))
	// for i, w := range widgets {
	// 	widgetMeta[i] = gin.H{
	// 		"widget_type": w.WidgetType,
	// 		"title":       w.Title,
	// 		"data":        w.Data,
	// 	}
	// }

	docMeta := make([]gin.H, len(docs))
	for i, d := range docs {
		docMeta[i] = gin.H{
			"title":       d.Title,
			"theme":       d.Theme,
			"output_path": d.OutputPath,
			"format":      d.Format,
		}
	}

	saveContext := context.WithoutCancel(c)

	assistantMessage, err := h.Q.InsertMessage(saveContext, db.InsertMessageParams{
		ConversationID: conv.ID,
		Role:           "assistant",
		Content:        answer,
		Widgets:        util.ToJSONB(widgets, len(widgets)),
		Documents:      util.ToJSONB(docMeta, len(docMeta)),
		Pages:          util.ToJSONB(pageRefs, len(pageRefs)),
		IsPartial:      isPartial,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the AI answer."})
		return
	}

	if _, err := h.Q.UpdateConversationTime(saveContext, conv.ID); err != nil {
		log.Printf("failed to update conversation time: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"conversation_id": util.UuidString(conv.ID),
		"title":           conv.Title,
		"user_message_id": util.UuidString(userMessage.ID),
		"message_id":      util.UuidString(assistantMessage.ID),
		"answer":          answer,
		"pages":           pageMeta,
		"widgets":         widgets,
		"documents":       docMeta,
		"is_partial":      isPartial,
	})

}

var isAllowedPath = []string{
	"sandbox/output/",
	"generated-docs/",
}

func (h *Handlers) DownloadDocument(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path is required"})
		return
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path"})
		return
	}

	if !isPathAllowed(absPath) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.FileAttachment(absPath, filepath.Base(absPath))
}

func isPathAllowed(absPath string) bool {
	for _, dir := range isAllowedPath {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		if strings.HasPrefix(absPath, absDir) {
			return true
		}
	}
	return false
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
