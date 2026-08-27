package main

import (
	"log"
	"os"

	"github.com/bergamo17/agentic-rag-prototype/internal/handlers"
	mlservice "github.com/bergamo17/agentic-rag-prototype/internal/mlservices"
	"github.com/bergamo17/agentic-rag-prototype/internal/openai"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env is not found")
	}

	mlServiceURL := os.Getenv("ML_SERVICE_URL")
	if mlServiceURL == "" {
		mlServiceURL = "http://localhost:8001"
		// mlServiceURL = "https://zm35wvtq-8001.use2.devtunnels.ms"
	}

	mlClient := mlservice.NewClient(mlServiceURL)
	openAIClient := openai.NewClient()
	h := handlers.New(mlClient, openAIClient)

	router := gin.Default()
	router.POST("/documents", h.EmbedDocument)
	router.POST("/chat", h.Chat)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Agentic RAG API (Go) run in: %s", port)
	router.Run(":" + port)
}
