package ai

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

// OllamaClient handles interactions with Ollama API for memory curation
type OllamaClient struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

// NewOllamaClient creates a new Ollama API client
func NewOllamaClient(baseURL, model string) *OllamaClient {
	if baseURL == "" {
		baseURL = defaultOllamaURL
	}
	if model == "" {
		model = "llama3.1"
	}

	return &OllamaClient{
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 300 * time.Second, // Ollama can be slow on CPU
		},
	}
}

// CurateMemories analyzes a transcript and extracts meaningful memories
func (c *OllamaClient) CurateMemories(req *CurationRequest) (*CurationResponse, error) {
	prompt := c.buildCurationPrompt(req.Transcript)

	// Call Ollama API
	response, err := c.callOllama(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to call Ollama API: %w", err)
	}

	// Parse the response
	curationResp, err := c.parseCurationResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse curation response: %w", err)
	}

	return curationResp, nil
}

// buildCurationPrompt creates the prompt for memory curation
func (c *OllamaClient) buildCurationPrompt(transcript string) string {
	return fmt.Sprintf(`You are a memory curator for an AI assistant. Your task is to analyze the following conversation transcript and extract the most important, meaningful memories that should be preserved.

For each memory, provide:
- content: A clear, concise statement of the memory
- importance_weight: A float between 0 and 1 indicating importance
- semantic_tags: Keywords that describe the memory
- context_type: One of: TECHNICAL_IMPLEMENTATION, ARCHITECTURE, DECISION, BREAKTHROUGH, RELATIONSHIP, UNRESOLVED, MILESTONE, PREFERENCE
- trigger_phrases: Phrases that should trigger recall of this memory
- question_types: Types of questions this memory would help answer
- temporal_relevance: "persistent", "session", or "temporary"
- action_required: Boolean indicating if follow-up action is needed
- reasoning: Why this memory is worth preserving

Also identify relationships between memories (references, supersedes, related_to, etc.)

Respond ONLY with valid JSON in this format:
{
  "memories": [
    {
      "content": "...",
      "importance_weight": 0.9,
      "semantic_tags": ["tag1", "tag2"],
      "context_type": "TECHNICAL_IMPLEMENTATION",
      "trigger_phrases": ["phrase1", "phrase2"],
      "question_types": ["how does X work", "what is Y"],
      "temporal_relevance": "persistent",
      "action_required": false,
      "reasoning": "..."
    }
  ],
  "relationships": [
    {
      "from_index": 0,
      "to_index": 1,
      "type": "references"
    }
  ],
  "summary": "Brief summary of the session"
}

TRANSCRIPT:
%s

Remember: Only extract memories that are genuinely worth preserving. Quality over quantity.`, transcript)
}

// parseCurationResponse parses the AI's JSON response
func (c *OllamaClient) parseCurationResponse(response string) (*CurationResponse, error) {
	var curation CurationResponse

	// Extract JSON from response (might include explanatory text)
	jsonStart := findJSONStart(response)
	jsonEnd := findJSONEnd(response)

	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	if err := json.Unmarshal([]byte(jsonStr), &curation); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &curation, nil
}

// ollamaRequest represents a request to Ollama API
type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format,omitempty"`
}

// ollamaResponse represents Ollama's response
type ollamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

// callOllama makes an API call to Ollama
func (c *OllamaClient) callOllama(prompt string) (string, error) {
	reqBody := ollamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
		Format: "json", // Request JSON format response
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ollama (is it running?): %w\n\nStart Ollama with: ollama serve\nPull model with: ollama pull %s", err, c.model)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama returned status %d: %s\n\nMake sure model is pulled: ollama pull %s",
			resp.StatusCode, string(body), c.model)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if ollamaResp.Response == "" {
		return "", fmt.Errorf("empty response from Ollama")
	}

	return ollamaResp.Response, nil
}
