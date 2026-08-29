package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bergamo17/agentic-rag-prototype/internal/agent"
	mlservice "github.com/bergamo17/agentic-rag-prototype/internal/mlservices"
	"github.com/bergamo17/agentic-rag-prototype/internal/openai"
	"github.com/bergamo17/agentic-rag-prototype/internal/websearch"
)

type Handlers struct {
	ML     *mlservice.Client
	OpenAI *openai.Client
	Web    *websearch.Client
}

func New(mlClient *mlservice.Client, openaiClient *openai.Client, webClient *websearch.Client) *Handlers {
	return &Handlers{
		ML:     mlClient,
		OpenAI: openaiClient,
		Web:    webClient,
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
	Query string `json:"query"`
	K     int    `json:"k"`
}

func (h *Handlers) Chat(c *gin.Context) {
	var req ChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messages := []openai.Message{
		{
			Role:    "system",
			Content: "You are a helpful assistant that can search documents and the web to answer questions.",
		},
		{
			Role:    "user",
			Content: req.Query,
		},
	}

	answer, pages, isPartial, err := agent.AgentLoop(h.OpenAI, h.Web, h.ML, messages)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	pageMeta := make([]gin.H, len(pages))
	for i, p := range pages {
		pageMeta[i] = gin.H{
			"document_id": p.DocumentID,
			"title":       p.Title,
			"page_number": p.PageNumber,
			"page_image":  p.ImageBase64,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"answer":     answer,
		"pages":      pageMeta,
		"is_partial": isPartial,
	})
}
