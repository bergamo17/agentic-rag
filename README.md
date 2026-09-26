# Agentic AI

An agentic AI system: a Go service runs the agent
loop and orchestrates tools, a Python ML service handles document embedding and
vector search, and a Next.js dashboard visualizes the results — all in a single
monorepo.

## Features

- **Agentic loop** — Go (Gin) service that drives a tool-calling agent loop
  (OpenAI function-calling style schema)
- **Document retrieval** — Python FastAPI service using ColQwen2 embeddings and
  Qdrant for vector search
- **Widget-based responses** — agent can respond with structured widgets (tables,
  charts) in addition to plain text, rendered by the dashboard
- **Interactive dashboard** — Next.js app with a chat interface, a Zustand store for
  state, and typed (TypeScript) data contracts shared with the backend
- **Document generation tool** — a `docbuilder` package exposing `create_docx`
  (pre-approved templates) and `execute_python` (custom document generation) tools,
  run inside an isolated, network-disabled Docker sandbox per request
- **Dual interface** — usable both as a CLI and over HTTP

## Architecture

```
.
├── (Go service)        # Agent loop, tool orchestration (Gin)
│   ├── agent
|   |   ├── loop.go          # Core agentic loop
|   |   ├── tools.go         # Tool definitions
|   ├── docbuilder
|   |   ├── client.go        # Function to build document when using generate doc tool with existed document template
|   |   ├── sandbox.go       # Docker-isolated sandbox runner for document generation using AI's python script
|   ├── handler
|   |   ├── handlers.go      # HTTP handlers
|   ├── mlservice
|   |   ├── client.go        # Go service to call the embedding and retrieval function
|   ├── openai
|   |   ├── client.go        # LLM model as the brain
|   ├── webseach
|   |   ├── client.go        # Web search function when using web search tool
├── cmd 
│   └── main.go
├── (Python ML service)  # FastAPI — embeddings & vector search
│   └── model-services             # ColQwen2 embeddings, Qdrant integration
└── dashboard/           # Next.js frontend
    ├── chat-window.tsx  # Chat UI
    └── ...               # Widget components, Zustand store, TS types
```

<!-- TODO: adjust folder names/paths to match your actual repo layout -->

## Tech Stack

| Layer            | Tools                                                |
|-------------------|-------------------------------------------------------|
| Agent / Backend   | Go, Gin                                               |
| ML Service        | Python, FastAPI, ColQwen2, Qdrant                     |
| Frontend          | Next.js, TypeScript, Zustand                          |
| Document Gen      | Docker (isolated sandbox, `--network none`)           |
| Interfaces        | CLI, HTTP                                             |

## Getting Started

### Prerequisites

- Go 1.2x+
- Python 3.1x+
- Docker
- A running Qdrant instance
- Node.js (for the dashboard)

### Setup

```bash
git clone https://github.com/bergamo17/agentic-rag.git
cd agentic-rag

# start the Python ML service
cd ml-service && pip install -r requirements.txt && uvicorn main:app --reload

# start Qdrant (if not already running)
docker run -p 6333:6333 qdrant/qdrant

# start the Go agent service
go run ./cmd/main.go

# start the dashboard
cd dashboard && npm install && npm run dev
```

<!-- TODO: confirm actual folder names, commands, and ports used in your repo -->

### Configuration

```
QDRANT_URL=http://localhost:6333
ML_SERVICE_URL=http://localhost:8000
OPENAI_API_KEY=<your-key>   # or ANTHROPIC_API_KEY, depending on model used
```

<!-- TODO: replace with your actual env vars -->

## Usage

The agent can be run in two modes:

- **CLI** — for local testing and debugging the agent loop directly
- **HTTP** — for use with the Next.js dashboard, which renders chat responses and
  structured widgets (tables, charts)

<!-- TODO: add actual CLI command / example usage -->

## Known Issues / In Progress

- CORS configuration between the Next.js dashboard and the Go backend needs to be
  finalized
- Dashboard (widget rendering, chat UI) is under active development

## License

<!-- TODO: add a license, or remove this section if none -->
