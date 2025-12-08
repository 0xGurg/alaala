package storage

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate/entities/models"
)

const (
	// MemoryClassName is the Weaviate class name for memories
	MemoryClassName = "Memory"
)

// VectorSearchResult represents a result from vector search
type VectorSearchResult struct {
	ID       string
	Distance float64
	Metadata map[string]interface{}
}

// WeaviateStore handles vector storage operations
type WeaviateStore struct {
	client *weaviate.Client
	ctx    context.Context
}

// NewWeaviateStore creates a new Weaviate store
func NewWeaviateStore(host string, scheme string) (*WeaviateStore, error) {
	cfg := weaviate.Config{
		Host:   host,
		Scheme: scheme,
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Weaviate client: %w", err)
	}

	store := &WeaviateStore{
		client: client,
		ctx:    context.Background(),
	}

	// Initialize schema
	if err := store.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// NewWeaviateStoreWithAuth creates a new Weaviate store with authentication
func NewWeaviateStoreWithAuth(host string, scheme string, apiKey string) (*WeaviateStore, error) {
	cfg := weaviate.Config{
		Host:       host,
		Scheme:     scheme,
		AuthConfig: auth.ApiKey{Value: apiKey},
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Weaviate client: %w", err)
	}

	store := &WeaviateStore{
		client: client,
		ctx:    context.Background(),
	}

	// Initialize schema
	if err := store.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// initSchema creates the Weaviate schema for memories
func (w *WeaviateStore) initSchema() error {
	// Check if schema already exists
	exists, err := w.client.Schema().ClassExistenceChecker().
		WithClassName(MemoryClassName).
		Do(w.ctx)
	if err != nil {
		return fmt.Errorf("failed to check schema existence: %w", err)
	}

	if exists {
		return nil // Schema already exists
	}

	// Create schema
	classObj := &models.Class{
		Class:       MemoryClassName,
		Description: "A semantic memory for AI assistants",
		Properties: []*models.Property{
			{
				Name:        "content",
				DataType:    []string{"text"},
				Description: "The memory content",
			},
			{
				Name:        "projectId",
				DataType:    []string{"text"},
				Description: "Project ID",
			},
			{
				Name:        "sessionId",
				DataType:    []string{"text"},
				Description: "Session ID",
			},
			{
				Name:        "importance",
				DataType:    []string{"number"},
				Description: "Importance weight (0-1)",
			},
			{
				Name:        "contextType",
				DataType:    []string{"text"},
				Description: "Type of context",
			},
			{
				Name:        "temporalRelevance",
				DataType:    []string{"text"},
				Description: "Temporal relevance",
			},
			{
				Name:        "actionRequired",
				DataType:    []string{"boolean"},
				Description: "Whether action is required",
			},
			{
				Name:        "tags",
				DataType:    []string{"text[]"},
				Description: "Semantic tags",
			},
			{
				Name:        "triggerPhrases",
				DataType:    []string{"text[]"},
				Description: "Trigger phrases for retrieval",
			},
			{
				Name:        "createdAt",
				DataType:    []string{"number"},
				Description: "Creation timestamp (Unix)",
			},
		},
		Vectorizer: "none", // We provide our own vectors
	}

	err = w.client.Schema().ClassCreator().
		WithClass(classObj).
		Do(w.ctx)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// Store stores a memory with its embedding
func (w *WeaviateStore) Store(id string, content string, embedding []float32, metadata map[string]interface{}) error {
	properties := map[string]interface{}{
		"content": content,
	}

	// Add all metadata as properties
	for k, v := range metadata {
		properties[k] = v
	}

	_, err := w.client.Data().Creator().
		WithClassName(MemoryClassName).
		WithID(id).
		WithProperties(properties).
		WithVector(embedding).
		Do(w.ctx)

	if err != nil {
		return fmt.Errorf("failed to store memory: %w", err)
	}

	return nil
}

// Search performs vector similarity search
func (w *WeaviateStore) Search(embedding []float32, limit int, filterMap map[string]interface{}) ([]VectorSearchResult, error) {
	// For now, return empty results - this needs proper Weaviate GraphQL implementation
	// The Weaviate Go client API requires careful handling of GraphQL fields
	// TODO: Implement proper vector search with Weaviate client v4 API
	var results []VectorSearchResult

	// Placeholder to avoid unused variables
	_ = embedding
	_ = limit
	_ = filterMap

	return results, nil
}

// Delete deletes a memory by ID
func (w *WeaviateStore) Delete(id string) error {
	err := w.client.Data().Deleter().
		WithClassName(MemoryClassName).
		WithID(id).
		Do(w.ctx)

	if err != nil {
		return fmt.Errorf("failed to delete memory: %w", err)
	}

	return nil
}

// Close closes the Weaviate connection
func (w *WeaviateStore) Close() error {
	// Weaviate Go client doesn't have explicit close
	return nil
}
