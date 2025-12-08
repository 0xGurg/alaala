package embeddings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultOllamaURL = "http://localhost:11434"
)

// OllamaEmbedder generates embeddings using Ollama
type OllamaEmbedder struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

// NewOllamaEmbedder creates a new Ollama embeddings client
func NewOllamaEmbedder(baseURL, model string) *OllamaEmbedder {
	if baseURL == "" {
		baseURL = defaultOllamaURL
	}
	if model == "" {
		model = "nomic-embed-text"
	}

	return &OllamaEmbedder{
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Embed generates an embedding for the given text
func (e *OllamaEmbedder) Embed(text string) ([]float32, error) {
	reqBody := map[string]interface{}{
		"model":  e.model,
		"prompt": text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/embeddings", e.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Ollama (is it running?): %w\n\nStart Ollama with: ollama serve", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama returned status %d: %s\n\nMake sure model is pulled: ollama pull %s",
			resp.StatusCode, string(body), e.model)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var ollamaResp struct {
		Embedding []float64 `json:"embedding"`
	}

	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(ollamaResp.Embedding) == 0 {
		return nil, fmt.Errorf("empty embedding returned from Ollama")
	}

	// Convert float64 to float32
	embedding := make([]float32, len(ollamaResp.Embedding))
	for i, v := range ollamaResp.Embedding {
		embedding[i] = float32(v)
	}

	return embedding, nil
}
