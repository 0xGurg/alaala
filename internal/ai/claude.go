package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	claudeAPIURL = "https://api.anthropic.com/v1/messages"
	apiVersion   = "2023-06-01"
)

// ClaudeClient handles interactions with Claude API for memory curation
type ClaudeClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewClaudeClient creates a new Claude API client
func NewClaudeClient(apiKey string, model string) *ClaudeClient {
	if model == "" {
		model = "claude-3-5-sonnet-20241022"
	}

	return &ClaudeClient{
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{},
	}
}

// CurateMemories analyzes a transcript and extracts meaningful memories
func (c *ClaudeClient) CurateMemories(req *CurationRequest) (*CurationResponse, error) {
	prompt := c.buildCurationPrompt(req.Transcript)

	// Call Claude API
	response, err := c.callClaude(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to call Claude API: %w", err)
	}

	// Parse the response
	curationResp, err := c.parseCurationResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse curation response: %w", err)
	}

	return curationResp, nil
}

// buildCurationPrompt creates the prompt for memory curation
func (c *ClaudeClient) buildCurationPrompt(transcript string) string {
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
func (c *ClaudeClient) parseCurationResponse(response string) (*CurationResponse, error) {
	var curation CurationResponse
	
	// Extract JSON from response (Claude might include explanatory text)
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

// claudeRequest represents a request to Claude API
type claudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []claudeMessage `json:"messages"`
}

// claudeMessage represents a message in the conversation
type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// claudeResponse represents Claude's response
type claudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
}

// callClaude makes an API call to Claude
func (c *ClaudeClient) callClaude(prompt string) (string, error) {
	reqBody := claudeRequest{
		Model:     c.model,
		MaxTokens: 4096,
		Messages: []claudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", claudeAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude")
	}

	return claudeResp.Content[0].Text, nil
}

// Helper functions

func findJSONStart(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == '{' {
			return i
		}
	}
	return -1
}

func findJSONEnd(s string) int {
	depth := 0
	start := -1
	
	for i := 0; i < len(s); i++ {
		if s[i] == '{' {
			if start == -1 {
				start = i
			}
			depth++
		} else if s[i] == '}' {
			depth--
			if depth == 0 && start != -1 {
				return i
			}
		}
	}
	
	return -1
}

