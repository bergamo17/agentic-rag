package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bergamo17/agentic-rag-prototype/internal/agent"
	"github.com/bergamo17/agentic-rag-prototype/internal/handlers"
	mlservice "github.com/bergamo17/agentic-rag-prototype/internal/mlservices"
	"github.com/bergamo17/agentic-rag-prototype/internal/openai"
	"github.com/bergamo17/agentic-rag-prototype/internal/websearch"
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
	webClient := websearch.NewClient(os.Getenv("TAVILY_API_KEY"))

	if len(os.Args) > 1 && os.Args[1] == "cli" {
		runCLI(mlClient, openAIClient, webClient)
		return
	}

	h := handlers.New(mlClient, openAIClient, webClient)

	router := gin.Default()
	router.POST("/documents", h.EmbedDocument)
	router.POST("/chat", h.Chat)
	router.POST("/chat/agent", h.ChatAgent)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Agentic RAG API (Go) run in: %s", port)
	router.Run(":" + port)
}

func runCLI(mlClient *mlservice.Client, openaiClient *openai.Client, webClient *websearch.Client) {
	fmt.Println("===Agentic AI CLI====")
	fmt.Println("Type your question, or 'exit' for quitting.")

	messages := []openai.Message{
		{
			Role:    "system",
			Content: "You are a helpful assistant that can search documents and the web to answer questions.",
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "exit" || input == "quit" {
			break
		}

		messages = append(messages, openai.Message{
			Role:    "user",
			Content: input,
		})

		answer, pages, err := agent.AgentLoop(openaiClient, webClient, mlClient, messages)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("\nAgent: %s\n", answer)
		if len(pages) > 0 {
			fmt.Println("\n[Referenced pages]")
			for _, p := range pages {
				fmt.Printf("  - %s, page %d (doc: %s)\n", p.Title, p.PageNumber, p.DocumentID)
			}
		}

		messages = append(messages, openai.Message{
			Role:    "assistant",
			Content: answer,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Printf("scanner error: %v", err)
	}
}
