package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/georgepagarigan/alaala/internal/memory"
)

// Resource represents an MCP resource
type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType"`
}

// handleListResources returns the list of available resources
func (s *Server) handleListResources(params json.RawMessage) (interface{}, error) {
	resources := []Resource{
		{
			URI:         "memory://session-context",
			Name:        "Session Context",
			Description: "Current session context with relevant memories",
			MimeType:    "text/plain",
		},
		{
			URI:         "memory://project-memories",
			Name:        "Project Memories",
			Description: "All memories for the current project",
			MimeType:    "application/json",
		},
	}

	return map[string]interface{}{
		"resources": resources,
	}, nil
}

// handleReadResource reads a resource
func (s *Server) handleReadResource(params json.RawMessage) (interface{}, error) {
	var req struct {
		URI string `json:"uri"`
	}

	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("invalid read resource params: %w", err)
	}

	switch req.URI {
	case "memory://session-context":
		return s.resourceSessionContext()
	case "memory://project-memories":
		return s.resourceProjectMemories()
	default:
		return nil, fmt.Errorf("unknown resource URI: %s", req.URI)
	}
}

// resourceSessionContext provides session context
func (s *Server) resourceSessionContext() (interface{}, error) {
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

	// Format as text
	text := formatSessionPrimer(primer)

	return map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"uri":      "memory://session-context",
				"mimeType": "text/plain",
				"text":     text,
			},
		},
	}, nil
}

// resourceProjectMemories provides all project memories
func (s *Server) resourceProjectMemories() (interface{}, error) {
	// Get current project
	projectID, err := s.getCurrentProjectID()
	if err != nil {
		return nil, err
	}

	// Search for all memories (high limit)
	results, err := s.engine.SearchMemories(&memory.SearchQuery{
		Query:         "",
		ProjectID:     projectID,
		Limit:         100,
		MinImportance: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get project memories: %w", err)
	}

	// Format memories
	var memories []map[string]interface{}
	for _, result := range results {
		memories = append(memories, map[string]interface{}{
			"id":          result.Memory.ID,
			"content":     result.Memory.Content,
			"importance":  result.Memory.Importance,
			"tags":        result.Memory.SemanticTags,
			"contextType": result.Memory.ContextType,
			"createdAt":   result.Memory.CreatedAt,
		})
	}

	data, err := json.Marshal(memories)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"uri":      "memory://project-memories",
				"mimeType": "application/json",
				"text":     string(data),
			},
		},
	}, nil
}

// Helper functions

func formatSessionPrimer(primer *memory.SessionPrimer) string {
	text := fmt.Sprintf("# Session Context for %s\n\n", primer.ProjectName)

	if primer.LastSessionDate != nil {
		text += fmt.Sprintf("Last session: %s\n\n", primer.TimeSinceLastSession)
	} else {
		text += "This is the first session for this project.\n\n"
	}

	if len(primer.TopMemories) > 0 {
		text += "## Key Memories:\n\n"
		for i, mem := range primer.TopMemories {
			text += fmt.Sprintf("%d. %s\n", i+1, mem.Content)
			if len(mem.SemanticTags) > 0 {
				text += fmt.Sprintf("   Tags: %v\n", mem.SemanticTags)
			}
			text += "\n"
		}
	}

	if len(primer.UnresolvedItems) > 0 {
		text += "## Unresolved Items:\n\n"
		for i, mem := range primer.UnresolvedItems {
			text += fmt.Sprintf("%d. %s\n\n", i+1, mem.Content)
		}
	}

	return text
}

