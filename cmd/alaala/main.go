package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/0xGurg/alaala/internal/ai"
	"github.com/0xGurg/alaala/internal/embeddings"
	"github.com/0xGurg/alaala/internal/mcp"
	"github.com/0xGurg/alaala/internal/memory"
	"github.com/0xGurg/alaala/internal/storage"
	"github.com/0xGurg/alaala/pkg/config"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "serve":
		serveMCP()
	case "init":
		initProject()
	case "version":
		printVersion()
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`alaala - Semantic memory system for AI assistants

Usage:
  alaala <command> [options]

Commands:
  serve      Start the MCP server (for Cursor/Claude Desktop integration)
  init       Initialize a new project with .alaala-project.json
  version    Print version information
  help       Show this help message

Examples:
  # Start MCP server for Cursor
  alaala serve

  # Initialize project
  alaala init

Installation:
  brew tap 0xGurg/distillery && brew install alaala

Uninstallation:
  brew uninstall alaala && brew untap 0xGurg/distillery

For more information, visit: https://github.com/0xGurg/alaala
`)
}

func serveMCP() {
	// Load configuration
	cfg, err := config.Load(config.GetConfigPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Loaded config from: %s\n", config.GetConfigPath())
	fmt.Fprintf(os.Stderr, "Weaviate URL: %s\n", cfg.Storage.WeaviateURL)
	fmt.Fprintf(os.Stderr, "AI provider: %s\n", cfg.AI.Provider)

	// Initialize storage
	sqlStore, err := initSQLiteStore(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SQLite: %v\n", err)
		os.Exit(1)
	}
	defer sqlStore.Close()

	weaviateStore, err := initWeaviateStore(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize Weaviate: %v\n", err)
		os.Exit(1)
	}
	defer weaviateStore.Close()

	// Initialize embeddings
	embedder, err := initEmbeddings(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize embeddings: %v\n", err)
		os.Exit(1)
	}

	// Initialize memory engine
	engine := memory.NewEngine(sqlStore, weaviateStore, embedder)

	// Initialize AI client
	aiClient, err := initAIClient(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize AI client: %v\n", err)
		os.Exit(1)
	}

	// Initialize curator
	curator := memory.NewCurator(engine, aiClient)

	// Start MCP server
	mcpServer := mcp.NewServer(engine, curator)

	fmt.Fprintf(os.Stderr, "MCP server ready\n")

	if err := mcpServer.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
		os.Exit(1)
	}
}

func initProject() {
	fmt.Println("Initializing alaala project...")

	// Create .alaala-project.json
	projectFile := ".alaala-project.json"
	if _, err := os.Stat(projectFile); err == nil {
		fmt.Printf("Project already initialized (%s exists)\n", projectFile)
		return
	}

	cwd, _ := os.Getwd()
	projectName := filepath.Base(cwd)

	projectConfig := fmt.Sprintf(`{
  "name": "%s",
  "created": "%s",
  "version": "1"
}
`, projectName, time.Now().Format(time.RFC3339))

	if err := os.WriteFile(projectFile, []byte(projectConfig), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create project file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created %s\n", projectFile)
	fmt.Println("Project initialized successfully!")

	// Create default config if it doesn't exist
	cfgPath := config.GetConfigPath()
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		cfg := config.DefaultConfig()
		if err := cfg.Save(cfgPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to create default config: %v\n", err)
		} else {
			fmt.Printf("Created default config at %s\n", cfgPath)
		}
	}
}

// Initialization helper functions

func initSQLiteStore(cfg *config.Config) (*storage.SQLiteStore, error) {
	// Ensure directory exists
	dir := filepath.Dir(cfg.Storage.SQLitePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return storage.NewSQLiteStore(cfg.Storage.SQLitePath)
}

func initWeaviateStore(cfg *config.Config) (*storage.WeaviateStore, error) {
	// Parse Weaviate URL
	url := cfg.Storage.WeaviateURL
	if url == "" {
		url = "http://localhost:8080"
	}

	host := "localhost:8080"
	scheme := "http"

	// Simple URL parsing
	if len(url) > 8 && url[:8] == "https://" {
		scheme = "https"
		host = url[8:]
	} else if len(url) > 7 && url[:7] == "http://" {
		scheme = "http"
		host = url[7:]
	}

	return storage.NewWeaviateStore(host, scheme)
}

func initEmbeddings(cfg *config.Config) (*embeddings.Client, error) {
	return embeddings.NewClient(cfg.Embeddings.Provider, cfg.Embeddings.Model)
}

func initAIClient(cfg *config.Config) (memory.AIClient, error) {
	switch cfg.AI.Provider {
	case "anthropic":
		apiKey := cfg.AI.APIKey
		if apiKey == "" {
			apiKey = os.Getenv("ANTHROPIC_API_KEY")
		}
		if apiKey == "" {
			return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
		}
		return ai.NewClaudeClient(apiKey, cfg.AI.Model), nil
	case "openrouter":
		apiKey := cfg.AI.APIKey
		if apiKey == "" {
			apiKey = os.Getenv("OPENROUTER_API_KEY")
		}
		if apiKey == "" {
			return nil, fmt.Errorf("OPENROUTER_API_KEY not set")
		}
		return ai.NewOpenRouterClient(apiKey, cfg.AI.Model, cfg.AI.OpenRouterURL), nil
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", cfg.AI.Provider)
	}
}

func printVersion() {
	fmt.Printf("alaala version %s\n", version)
	if commit != "none" {
		fmt.Printf("  commit: %s\n", commit)
	}
	if date != "unknown" {
		fmt.Printf("  built:  %s\n", date)
	}
}
