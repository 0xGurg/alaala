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
	defaultOpenRouterURL = "https://openrouter.ai/api/v1"
)

// OpenRouterClient handles interactions with OpenRouter API for memory curation
// OpenRouter uses OpenAI-compatible API format
type OpenRouterClient struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

// NewOpenRouterClient creates a new OpenRouter API client
func NewOpenRouterClient(apiKey string, model string, baseURL string) *OpenRouterClient {
	if baseURL == "" {
		baseURL = defaultOpenRouterURL
	}
	if model == "" {
		model = "anthropic/claude-3.5-sonnet"
	}

	return &OpenRouterClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // OpenRouter can be slow for some models
		},
	}
}

// CurateMemories analyzes a transcript and extracts meaningful memories
func (c *OpenRouterClient) CurateMemories(req *CurationRequest) (*CurationResponse, error) {
	prompt := c.buildCurationPrompt(req.Transcript)

	// Call OpenRouter API
	response, err := c.callOpenRouter(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenRouter API: %w", err)
	}

	// Parse the response
	curationResp, err := c.parseCurationResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse curation response: %w", err)
	}

	return curationResp, nil
}

// buildCurationPrompt creates the prompt for memory curation
func (c *OpenRouterClient) buildCurationPrompt(transcript string) string {
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
func (c *OpenRouterClient) parseCurationResponse(response string) (*CurationResponse, error) {
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

// openRouterRequest represents a request to OpenRouter API (OpenAI-compatible format)
type openRouterRequest struct {
	Model    string                   `json:"model"`
	Messages []openRouterMessage      `json:"messages"`
	MaxTokens int                     `json:"max_tokens,omitempty"`
}

// openRouterMessage represents a message in the conversation
type openRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openRouterResponse represents OpenRouter's response (OpenAI-compatible format)
type openRouterResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// callOpenRouter makes an API call to OpenRouter with retry logic
func (c *OpenRouterClient) callOpenRouter(prompt string) (string, error) {
	var lastErr error
	maxRetries := 3
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			time.Sleep(backoff)
		}
		
		response, err := c.makeRequest(prompt)
		if err == nil {
			return response, nil
		}
		
		lastErr = err
		
		// Don't retry on certain errors
		if !c.shouldRetry(err) {
			return "", err
		}
	}
	
	return "", fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

// makeRequest performs a single API request
func (c *OpenRouterClient) makeRequest(prompt string) (string, error) {
	reqBody := openRouterRequest{
		Model: c.model,
		Messages: []openRouterMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens: 4096,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/chat/completions", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("HTTP-Referer", "https://github.com/georgepagarigan/alaala") // Optional but recommended
	req.Header.Set("X-Title", "alaala") // Optional but recommended

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var openRouterResp openRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if openRouterResp.Error != nil {
		return "", c.formatAPIError(openRouterResp.Error, resp.StatusCode)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	if len(openRouterResp.Choices) == 0 {
		return "", fmt.Errorf("empty response from OpenRouter")
	}

	return openRouterResp.Choices[0].Message.Content, nil
}

// shouldRetry determines if an error is retryable
func (c *OpenRouterClient) shouldRetry(err error) bool {
	errStr := err.Error()
	
	// Retry on rate limits
	if contains(errStr, "rate limit") || contains(errStr, "429") {
		return true
	}
	
	// Retry on temporary errors
	if contains(errStr, "timeout") || contains(errStr, "connection") {
		return true
	}
	
	// Retry on server errors (5xx)
	if contains(errStr, "500") || contains(errStr, "502") || contains(errStr, "503") {
		return true
	}
	
	// Don't retry on client errors (4xx except 429)
	return false
}

// formatAPIError creates a helpful error message
func (c *OpenRouterClient) formatAPIError(apiErr *struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}, statusCode int) error {
	baseMsg := fmt.Sprintf("OpenRouter API error: %s", apiErr.Message)
	
	// Add helpful suggestions based on error type
	switch {
	case contains(apiErr.Code, "invalid_api_key") || contains(apiErr.Code, "authentication"):
		return fmt.Errorf("%s\n\nPlease check your OPENROUTER_API_KEY environment variable", baseMsg)
	
	case contains(apiErr.Code, "rate_limit") || statusCode == 429:
		return fmt.Errorf("%s\n\nYou've hit the rate limit. The request will be retried automatically", baseMsg)
	
	case contains(apiErr.Message, "model") && contains(apiErr.Message, "not found"):
		return fmt.Errorf("%s\n\nModel '%s' is not available. Try: anthropic/claude-3.5-sonnet, openai/gpt-4-turbo, or meta-llama/llama-3.1-70b-instruct", 
			baseMsg, c.model)
	
	case contains(apiErr.Code, "insufficient_quota") || contains(apiErr.Message, "credits"):
		return fmt.Errorf("%s\n\nInsufficient credits. Please add credits to your OpenRouter account", baseMsg)
	
	default:
		return fmt.Errorf("%s (type: %s, code: %s)", baseMsg, apiErr.Type, apiErr.Code)
	}
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	s = toLowerSimple(s)
	substr = toLowerSimple(substr)
	return len(s) >= len(substr) && 
		   (s == substr || findSubstring(s, substr) != -1)
}

func toLowerSimple(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func findSubstring(haystack, needle string) int {
	if len(needle) > len(haystack) {
		return -1
	}
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}

