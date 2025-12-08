package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the complete application configuration
type Config struct {
	Storage    StorageConfig    `yaml:"storage"`
	AI         AIConfig         `yaml:"ai"`
	Embeddings EmbeddingsConfig `yaml:"embeddings"`
	Retrieval  RetrievalConfig  `yaml:"retrieval"`
	Web        WebConfig        `yaml:"web"`
	Logging    LoggingConfig    `yaml:"logging"`
}

// StorageConfig holds storage-related configuration
type StorageConfig struct {
	Mode       string         `yaml:"mode"` // "embedded" or "docker"
	Weaviate   WeaviateConfig `yaml:"weaviate"`
	SQLitePath string         `yaml:"sqlite_path"`
}

// WeaviateConfig holds Weaviate-specific configuration
type WeaviateConfig struct {
	EmbeddedPath string `yaml:"embedded_path"`
	DockerURL    string `yaml:"docker_url"`
}

// AIConfig holds AI provider configuration
type AIConfig struct {
	Provider      string `yaml:"provider"` // "anthropic", "openai", "ollama", "openrouter"
	APIKey        string `yaml:"api_key"`
	Model         string `yaml:"model"`
	OpenRouterURL string `yaml:"openrouter_url"` // Default: https://openrouter.ai/api/v1
}

// EmbeddingsConfig holds embeddings configuration
type EmbeddingsConfig struct {
	Provider string `yaml:"provider"` // "local", "openai", "ollama"
	Model    string `yaml:"model"`
}

// RetrievalConfig holds memory retrieval configuration
type RetrievalConfig struct {
	MaxMemories       int     `yaml:"max_memories"`
	MinImportance     float64 `yaml:"min_importance"`
	IncludeGraphDepth int     `yaml:"include_graph_depth"`
}

// WebConfig holds web UI configuration
type WebConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Host    string `yaml:"host"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level string `yaml:"level"` // "debug", "info", "warn", "error"
	File  string `yaml:"file"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	alaalaDir := filepath.Join(homeDir, ".alaala")

	return &Config{
		Storage: StorageConfig{
			Mode: "embedded",
			Weaviate: WeaviateConfig{
				EmbeddedPath: filepath.Join(alaalaDir, "weaviate"),
				DockerURL:    "http://localhost:8080",
			},
			SQLitePath: filepath.Join(alaalaDir, "alaala.db"),
		},
		AI: AIConfig{
			Provider:      "anthropic",
			Model:         "claude-3-5-sonnet-20241022",
			OpenRouterURL: "https://openrouter.ai/api/v1",
		},
		Embeddings: EmbeddingsConfig{
			Provider: "local",
			Model:    "all-MiniLM-L6-v2",
		},
		Retrieval: RetrievalConfig{
			MaxMemories:       5,
			MinImportance:     0.3,
			IncludeGraphDepth: 1,
		},
		Web: WebConfig{
			Enabled: true,
			Port:    8766,
			Host:    "localhost",
		},
		Logging: LoggingConfig{
			Level: "info",
			File:  filepath.Join(alaalaDir, "alaala.log"),
		},
	}
}

// Load reads configuration from a YAML file
func Load(path string) (*Config, error) {
	// Start with defaults
	cfg := DefaultConfig()

	// If no config file exists, return defaults
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	// Read config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Expand environment variables
	expanded := os.ExpandEnv(string(data))

	// Parse YAML
	if err := yaml.Unmarshal([]byte(expanded), cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// Save writes configuration to a YAML file
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigPath returns the default configuration file path
func GetConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".alaala", "config.yaml")
}
