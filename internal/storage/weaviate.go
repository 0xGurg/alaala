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
	// Build near vector argument
	nearVector := w.client.GraphQL().NearVectorArgBuilder().
		WithVector(embedding)

	// Build the query
	query := w.client.GraphQL().Get().
		WithClassName(MemoryClassName).
		WithNearVector(nearVector).
		WithLimit(limit)

	// Add filters if provided
	if projectID, ok := filterMap["project_id"].(string); ok && projectID != "" {
		// Simple project filter - just query and parse results manually
		// More complex filters can be added later
		_ = projectID // Will use in manual filtering below
	}

	// Execute the query - we need to get the raw response
	result, err := query.Do(w.ctx)
	if err != nil {
		return nil, fmt.Errorf("weaviate query failed: %w", err)
	}

	// Parse results
	var searchResults []VectorSearchResult

	// Extract data from GraphQL response
	if result.Data == nil {
		return searchResults, nil
	}

	getData, ok := result.Data["Get"].(map[string]interface{})
	if !ok {
		return searchResults, nil
	}

	memories, ok := getData[MemoryClassName].([]interface{})
	if !ok {
		return searchResults, nil
	}

	// Parse each memory result
	for _, item := range memories {
		memData, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// Try to get ID and certainty/distance from _additional
		id := ""
		distance := 0.0

		if additional, ok := memData["_additional"].(map[string]interface{}); ok {
			if idVal, ok := additional["id"].(string); ok {
				id = idVal
			}
			// Weaviate might return "certainty" or "distance"
			if distVal, ok := additional["distance"].(float64); ok {
				distance = distVal
			} else if certVal, ok := additional["certainty"].(float64); ok {
				distance = 1.0 - certVal // Convert certainty to distance
			}
		}

		if id == "" {
			continue
		}

		// Apply project filter if specified (manual filtering)
		if projectID, ok := filterMap["project_id"].(string); ok && projectID != "" {
			if projID, ok := memData["projectId"].(string); ok && projID != projectID {
				continue // Skip if project doesn't match
			}
		}

		// Apply importance filter if specified
		if minImp, ok := filterMap["importance_gte"].(float64); ok {
			if imp, ok := memData["importance"].(float64); ok && imp < minImp {
				continue // Skip if importance too low
			}
		}

		searchResults = append(searchResults, VectorSearchResult{
			ID:       id,
			Distance: distance,
			Metadata: memData,
		})
	}

	return searchResults, nil
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
