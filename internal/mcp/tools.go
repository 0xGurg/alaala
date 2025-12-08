package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/0xGurg/alaala/internal/memory"
)

// Tool represents an MCP tool
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// handleListTools returns the list of available tools
func (s *Server) handleListTools(params json.RawMessage) (interface{}, error) {
	tools := []Tool{
		{
			Name:        "search_memories",
			Description: "Search for relevant memories based on a query",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "The search query",
					},
					"limit": map[string]interface{}{
						"type":        "number",
						"description": "Maximum number of memories to return",
						"default":     5,
					},
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "Project ID to search within (optional)",
					},
					"min_importance": map[string]interface{}{
						"type":        "number",
						"description": "Minimum importance threshold (0-1)",
						"default":     0.3,
					},
				},
				"required": []string{"query"},
			},
		},
		{
			Name:        "save_memory",
			Description: "Save a new memory",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"content": map[string]interface{}{
						"type":        "string",
						"description": "The memory content",
					},
					"importance": map[string]interface{}{
						"type":        "number",
						"description": "Importance weight (0-1)",
						"default":     0.5,
					},
					"tags": map[string]interface{}{
						"type":        "array",
						"description": "Semantic tags",
						"items":       map[string]string{"type": "string"},
					},
					"context_type": map[string]interface{}{
						"type":        "string",
						"description": "Context type (TECHNICAL_IMPLEMENTATION, ARCHITECTURE, etc.)",
					},
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "Project ID",
					},
				},
				"required": []string{"content", "project_id"},
			},
		},
		{
			Name:        "curate_session",
			Description: "Curate memories from a session transcript",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"transcript": map[string]interface{}{
						"type":        "string",
						"description": "The conversation transcript",
					},
					"session_id": map[string]interface{}{
						"type":        "string",
						"description": "Session ID",
					},
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "Project ID",
					},
				},
				"required": []string{"transcript", "project_id"},
			},
		},
		{
			Name:        "list_projects",
			Description: "List all projects",
			InputSchema: map[string]interface{}{
				"type": "object",
			},
		},
	}

	return map[string]interface{}{
		"tools": tools,
	}, nil
}

// handleCallTool executes a tool
func (s *Server) handleCallTool(params json.RawMessage) (interface{}, error) {
	var req struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}

	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("invalid tool call params: %w", err)
	}

	switch req.Name {
	case "search_memories":
		return s.toolSearchMemories(req.Arguments)
	case "save_memory":
		return s.toolSaveMemory(req.Arguments)
	case "curate_session":
		return s.toolCurateSession(req.Arguments)
	case "list_projects":
		return s.toolListProjects(req.Arguments)
	default:
		return nil, fmt.Errorf("unknown tool: %s", req.Name)
	}
}

// toolSearchMemories implements the search_memories tool
func (s *Server) toolSearchMemories(args json.RawMessage) (interface{}, error) {
	var params struct {
		Query         string  `json:"query"`
		Limit         int     `json:"limit"`
		ProjectID     string  `json:"project_id"`
		MinImportance float64 `json:"min_importance"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	// Default values
	if params.Limit == 0 {
		params.Limit = 5
	}
	if params.MinImportance == 0 {
		params.MinImportance = 0.3
	}

	// Get current project if not specified
	if params.ProjectID == "" {
		projectID, err := s.getCurrentProjectID()
		if err != nil {
			return nil, err
		}
		params.ProjectID = projectID
	}

	// Search memories
	query := &memory.SearchQuery{
		Query:         params.Query,
		ProjectID:     params.ProjectID,
		Limit:         params.Limit,
		MinImportance: params.MinImportance,
	}

	results, err := s.engine.SearchMemories(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search memories: %w", err)
	}

	// Format results
	var memories []map[string]interface{}
	for _, result := range results {
		memories = append(memories, map[string]interface{}{
			"id":               result.Memory.ID,
			"content":          result.Memory.Content,
			"importance":       result.Memory.Importance,
			"tags":             result.Memory.SemanticTags,
			"context_type":     result.Memory.ContextType,
			"similarity_score": result.SimilarityScore,
			"relevance_score":  result.RelevanceScore,
			"trigger_matched":  result.TriggerMatched,
			"created_at":       result.Memory.CreatedAt,
		})
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": formatMemoriesAsText(memories),
			},
		},
	}, nil
}

// toolSaveMemory implements the save_memory tool
func (s *Server) toolSaveMemory(args json.RawMessage) (interface{}, error) {
	var params struct {
		Content     string   `json:"content"`
		Importance  float64  `json:"importance"`
		Tags        []string `json:"tags"`
		ContextType string   `json:"context_type"`
		ProjectID   string   `json:"project_id"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	// Default importance
	if params.Importance == 0 {
		params.Importance = 0.5
	}

	// Create memory
	mem := &memory.Memory{
		ProjectID:    params.ProjectID,
		Content:      params.Content,
		Importance:   params.Importance,
		SemanticTags: params.Tags,
		ContextType:  memory.ContextType(params.ContextType),
	}

	if err := s.engine.CreateMemory(mem); err != nil {
		return nil, fmt.Errorf("failed to create memory: %w", err)
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Memory saved successfully with ID: %s", mem.ID),
			},
		},
	}, nil
}

// toolCurateSession implements the curate_session tool
func (s *Server) toolCurateSession(args json.RawMessage) (interface{}, error) {
	var params struct {
		Transcript string `json:"transcript"`
		SessionID  string `json:"session_id"`
		ProjectID  string `json:"project_id"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	// Curate memories
	result, err := s.curator.CurateSession(params.ProjectID, params.SessionID, params.Transcript)
	if err != nil {
		return nil, fmt.Errorf("failed to curate session: %w", err)
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Curated %d memories from session. Summary: %s", len(result.Memories), result.Summary),
			},
		},
	}, nil
}

// toolListProjects implements the list_projects tool
func (s *Server) toolListProjects(args json.RawMessage) (interface{}, error) {
	// TODO: Implement project listing
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": "Project listing not yet implemented",
			},
		},
	}, nil
}

// Helper functions

func (s *Server) getCurrentProjectID() (string, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	// Look for .alaala-project.json
	projectFile := ".alaala-project.json"
	if _, err := os.Stat(projectFile); err != nil {
		// Create a new project
		projectName := filepath.Base(cwd)
		project, err := s.engine.GetOrCreateProject(projectName, cwd)
		if err != nil {
			return "", err
		}
		return project.ID, nil
	}

	// Read project file
	var projectConfig struct {
		Name string `json:"name"`
	}
	data, err := os.ReadFile(projectFile)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(data, &projectConfig); err != nil {
		return "", err
	}

	// Get or create project
	project, err := s.engine.GetOrCreateProject(projectConfig.Name, cwd)
	if err != nil {
		return "", err
	}

	return project.ID, nil
}

func formatMemoriesAsText(memories []map[string]interface{}) string {
	if len(memories) == 0 {
		return "No memories found."
	}

	result := fmt.Sprintf("Found %d relevant memories:\n\n", len(memories))
	for i, mem := range memories {
		result += fmt.Sprintf("%d. %s\n", i+1, mem["content"])
		result += fmt.Sprintf("   Importance: %.2f | Relevance: %.2f\n", mem["importance"], mem["relevance_score"])
		if tags, ok := mem["tags"].([]string); ok && len(tags) > 0 {
			result += fmt.Sprintf("   Tags: %v\n", tags)
		}
		result += "\n"
	}

	return result
}
