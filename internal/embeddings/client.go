package embeddings

import (
	"fmt"
)

// Client handles text embedding generation
type Client struct {
	provider string
	model    string
}

// NewClient creates a new embeddings client
func NewClient(provider, model string) (*Client, error) {
	return &Client{
		provider: provider,
		model:    model,
	}, nil
}

// Embed generates an embedding vector for the given text
func (c *Client) Embed(text string) ([]float32, error) {
	switch c.provider {
	case "local":
		return c.embedLocal(text)
	default:
		return nil, fmt.Errorf("unknown embeddings provider: %s (only 'local' is supported)", c.provider)
	}
}

// embedLocal generates embeddings using a local model
func (c *Client) embedLocal(text string) ([]float32, error) {
	// TODO: Implement actual local embeddings using sentence-transformers
	// For now, return a dummy embedding vector
	// This should be replaced with actual model inference

	// Dummy 384-dimensional vector (typical for all-MiniLM-L6-v2)
	embedding := make([]float32, 384)

	// Simple hash-based fake embedding for development
	hash := simpleHash(text)
	for i := 0; i < 384; i++ {
		embedding[i] = float32((hash+i)%100) / 100.0
	}

	return embedding, nil
}

// simpleHash creates a simple hash of a string (for dummy embeddings)
func simpleHash(s string) int {
	h := 0
	for i := 0; i < len(s); i++ {
		h = 31*h + int(s[i])
	}
	if h < 0 {
		h = -h
	}
	return h
}
