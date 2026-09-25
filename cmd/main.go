package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bergamo17/agentic-rag-prototype/internal/agent"
	"github.com/bergamo17/agentic-rag-prototype/internal/docbuilder"
	"github.com/bergamo17/agentic-rag-prototype/internal/handlers"
	mlservice "github.com/bergamo17/agentic-rag-prototype/internal/mlservices"
	"github.com/bergamo17/agentic-rag-prototype/internal/openai"
	"github.com/bergamo17/agentic-rag-prototype/internal/websearch"

	"github.com/gin-contrib/cors"
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

	docPythonPath := os.Getenv("DOC_BUILDER_PYTHON_PATH")
	if docPythonPath == "" {
		docPythonPath = "python3"
	}

	docScriptPath := os.Getenv("DOC_BUILDER_SCRIPT_PATH")
	if docScriptPath == "" {
		docScriptPath = "document-services/document_builder.py"
	}

	sandboxInputPath := os.Getenv("SANBOX_INPUT_DIR")
	if sandboxInputPath == "" {
		sandboxInputPath = "sandbox/input"
	}

	sandboxOutputPath := os.Getenv("SANDBOX_OUTPUT_DIR")
	if sandboxOutputPath == "" {
		sandboxOutputPath = "sandbox/output"
	}

	mlClient := mlservice.NewClient(mlServiceURL)
	openAIClient := openai.NewClient()
	webClient := websearch.NewClient(os.Getenv("TAVILY_API_KEY"))
	docClient := docbuilder.NewClient(docPythonPath, docScriptPath)
	sandboxClient := docbuilder.NewSandboxClient(sandboxInputPath, sandboxOutputPath)

	if len(os.Args) > 1 && os.Args[1] == "cli" {
		runCLI(mlClient, openAIClient, webClient, docClient, sandboxClient)
		return
	}

	h := handlers.New(mlClient, openAIClient, webClient, docClient, sandboxClient)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.POST("/documents", h.EmbedDocument)
	router.POST("/chat", h.Chat)
	router.POST("/chat/agent", h.ChatAgent)
	router.GET("/documents/download", h.DownloadDocument)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Agentic RAG API (Go) run in: %s", port)
	router.Run(":" + port)
}

func runCLI(mlClient *mlservice.Client, openaiClient *openai.Client, webClient *websearch.Client, docClient *docbuilder.Client, sanboxClient *docbuilder.SandboxClient) {
	fmt.Println("===Agentic AI CLI====")
	fmt.Println("Type your question, or 'exit' for quitting.")

	ctx := context.Background()

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

		answer, pages, widgets, docs, isPartial, err := agent.AgentLoop(openaiClient, webClient, mlClient, docClient, sanboxClient, ctx, messages)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if isPartial == true {
			fmt.Printf("\nAgent: [Disclaimer] The final Answer is partially completed > %s\n", answer)
		} else {
			fmt.Printf("\nAgent: %s\n", answer)
		}

		if len(pages) > 0 {
			fmt.Println("\n[Referenced pages]")
			for _, p := range pages {
				fmt.Printf("  - %s, page %d (doc: %s)\n", p.Title, p.PageNumber, p.DocumentID)
			}
		}

		if len(widgets) > 0 {
			fmt.Println("\n[Generated widgets]")
			for _, w := range widgets {
				fmt.Printf("  - [%s] %s\n", w.WidgetType, w.Title)
				fmt.Printf("    data: %s\n", w.Data)
			}
		}

		if len(docs) > 0 {
			fmt.Println("\n[Generated docs]")
			for _, d := range docs {
				fmt.Printf("  - %s (theme: %s)\n", d.Title, d.Theme)
				fmt.Printf("    saved to: %s\n", d.OutputPath)
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
