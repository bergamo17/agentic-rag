package server

import (
	"github.com/bergamo17/agentic-rag-prototype/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	h      *handlers.Handlers
}

func New(h *handlers.Handlers, allowedOrigins []string) *Server {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	s := &Server{
		router: r,
		h:      h,
	}

	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.router.POST("/documents", s.h.EmbedDocument)
	s.router.POST("/chat", s.h.Chat)
	s.router.POST("/chat/agent", s.h.ChatAgent)
	s.router.GET("/documents/download", s.h.DownloadDocument)
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
