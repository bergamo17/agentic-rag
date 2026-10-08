package handlers

import (
	"net/http"
	"strings"
	"time"

	db "github.com/bergamo17/agentic-rag-prototype/db/sqlc"
	"github.com/bergamo17/agentic-rag-prototype/internal/openai"
	"github.com/bergamo17/agentic-rag-prototype/util"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type conversationResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type messageResponse struct {
	ID        string                     `json:"id"`
	Role      string                     `json:"role"`
	Content   string                     `json:"content"`
	Widgets   []openai.Widget            `json:"widgets"`
	Documents []openai.GeneratedDocument `json:"documents"`
	IsPartial bool                       `json:"is_partial"`
}

type convUri struct {
	ID string `uri:"id"`
}

type renameConvRequest struct {
	Title string `json:"title" binding:"required"`
}

func (h *Handlers) ListConversations(ctx *gin.Context) {
	listConv, err := h.Q.ListConversations(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get the list of conversation"})
		return
	}

	result := make([]conversationResponse, 0, len(listConv))
	for _, r := range listConv {
		result = append(result, conversationResponse{
			ID:        util.UuidString(r.ID),
			Title:     r.Title,
			CreatedAt: r.CreatedAt.Time,
			UpdatedAt: r.UpdatedAt.Time,
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handlers) GetConversationMessages(ctx *gin.Context) {
	var req convUri

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg, err := util.ParseUUID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse conversation ID"})
		return
	}

	conv, err := h.Q.GetConversation(ctx, arg)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get the conversation"})
		return
	}

	listMessages, err := h.Q.ListMessageByConversationID(ctx, conv.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get the list of conversation"})
		return
	}

	result := make([]messageResponse, 0, len(listMessages))
	for _, m := range listMessages {
		result = append(result, messageResponse{
			ID:        util.UuidString(m.ID),
			Role:      m.Role,
			Content:   m.Content,
			Widgets:   util.DecodeJSONB[openai.Widget](m.Widgets),
			Documents: util.DecodeJSONB[openai.GeneratedDocument](m.Documents),
			IsPartial: m.IsPartial,
		})
	}

	ctx.JSON(http.StatusOK, result)

}

func (h *Handlers) RenameConversation(ctx *gin.Context) {
	var renameUri convUri

	err := ctx.ShouldBindUri(&renameUri)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	var renameRequest renameConvRequest

	err = ctx.ShouldBindJSON(&renameRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	id, err := util.ParseUUID(renameUri.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse conversation ID"})
		return
	}

	title := strings.TrimSpace(renameRequest.Title)
	if title == "" || len([]rune(title)) > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Title must be 1-100 characters"})
		return
	}

	arg := db.UpdateConversationTitleParams{
		ID:    id,
		Title: renameRequest.Title,
	}

	updatedConversation, err := h.Q.UpdateConversationTitle(ctx, arg)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the conversation"})
		return
	}

	result := &conversationResponse{
		ID:        util.UuidString(updatedConversation.ID),
		Title:     updatedConversation.Title,
		CreatedAt: updatedConversation.CreatedAt.Time,
		UpdatedAt: updatedConversation.UpdatedAt.Time,
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handlers) DeleteConversation(ctx *gin.Context) {
	var req convUri

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	id, err := util.ParseUUID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	err = h.Q.DeleteConversation(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete conversation"})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
