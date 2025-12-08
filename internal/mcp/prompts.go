package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/0xGurg/alaala/internal/memory"
)

// Prompt represents an MCP prompt
type Prompt struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Arguments   []map[string]interface{} `json:"arguments,omitempty"`
}

// handleListPrompts returns the list of available prompts
func (s *Server) handleListPrompts(params json.RawMessage) (interface{}, error) {
	prompts := []Prompt{
		{
			Name:        "session_primer",
			Description: "Session primer with temporal context and relevant memories",
			Arguments:   []map[string]interface{}{},
		},
	}

	return map[string]interface{}{
		"prompts": prompts,
	}, nil
}

// handleGetPrompt gets a prompt
func (s *Server) handleGetPrompt(params json.RawMessage) (interface{}, error) {
	var req struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("invalid get prompt params: %w", err)
	}

	switch req.Name {
	case "session_primer":
		return s.promptSessionPrimer()
	default:
		return nil, fmt.Errorf("unknown prompt: %s", req.Name)
	}
}

// promptSessionPrimer generates the session primer prompt
func (s *Server) promptSessionPrimer() (interface{}, error) {
	// Get current project
	projectID, err := s.getCurrentProjectID()
	if err != nil {
		return nil, err
	}

	// Get session primer
	primer, err := s.engine.GetSessionPrimer(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session primer: %w", err)
	}

	// Format as prompt
	text := formatSessionPrimerAsPrompt(primer)

	return map[string]interface{}{
		"description": "Session context and relevant memories",
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": map[string]interface{}{
					"type": "text",
					"text": text,
				},
			},
		},
	}, nil
}

// Helper functions

func formatSessionPrimerAsPrompt(primer *memory.SessionPrimer) string {
	text := "# Session Context\n\n"
	text += fmt.Sprintf("Project: %s\n\n", primer.ProjectName)

	if primer.LastSessionDate != nil {
		text += fmt.Sprintf("Time since last session: %s\n\n", primer.TimeSinceLastSession)

		if primer.LastSessionSummary != "" {
			text += fmt.Sprintf("Last session summary: %s\n\n", primer.LastSessionSummary)
		}
	} else {
		text += "This is the first session for this project.\n\n"
	}

	if len(primer.TopMemories) > 0 {
		text += "## Relevant Context\n\n"
		text += "Here are the most relevant memories for this session:\n\n"

		for i, mem := range primer.TopMemories {
			text += fmt.Sprintf("%d. **%s**\n", i+1, mem.Content)

			if len(mem.SemanticTags) > 0 {
				text += fmt.Sprintf("   - Tags: %v\n", mem.SemanticTags)
			}

			if mem.ContextType != "" {
				text += fmt.Sprintf("   - Type: %s\n", mem.ContextType)
			}

			text += "\n"
		}
	}

	if len(primer.UnresolvedItems) > 0 {
		text += "## Unresolved Items\n\n"
		text += "These items need attention:\n\n"

		for i, mem := range primer.UnresolvedItems {
			text += fmt.Sprintf("%d. %s\n\n", i+1, mem.Content)
		}
	}

	text += "\n---\n\n"
	text += "Memories will surface naturally as we converse. You can search for specific memories or save important insights as we work together.\n"

	return text
}
