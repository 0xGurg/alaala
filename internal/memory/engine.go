package memory

import (
	"fmt"
	"time"

	"github.com/0xGurg/alaala/internal/storage"
	"github.com/google/uuid"
)

// Engine is the core memory management system
type Engine struct {
	sqlStore    *storage.SQLiteStore
	vectorStore VectorStore
	embedder    Embedder
}

// VectorStore is an interface for vector database operations
type VectorStore interface {
	Store(id string, content string, embedding []float32, metadata map[string]interface{}) error
	Search(embedding []float32, limit int, filters map[string]interface{}) ([]storage.VectorSearchResult, error)
	Delete(id string) error
}

// Embedder is an interface for generating embeddings
type Embedder interface {
	Embed(text string) ([]float32, error)
}

// NewEngine creates a new memory engine
func NewEngine(sqlStore *storage.SQLiteStore, vectorStore VectorStore, embedder Embedder) *Engine {
	return &Engine{
		sqlStore:    sqlStore,
		vectorStore: vectorStore,
		embedder:    embedder,
	}
}

// CreateMemory creates a new memory
func (e *Engine) CreateMemory(mem *Memory) error {
	// Generate ID if not provided
	if mem.ID == "" {
		mem.ID = uuid.New().String()
	}

	// Generate embedding
	embedding, err := e.embedder.Embed(mem.Content)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Store in SQLite
	sqlMemory := &storage.Memory{
		ID:                mem.ID,
		ProjectID:         mem.ProjectID,
		SessionID:         stringPtr(mem.SessionID),
		Content:           mem.Content,
		Importance:        mem.Importance,
		ContextType:       stringPtr(string(mem.ContextType)),
		TemporalRelevance: stringPtr(string(mem.TemporalRelevance)),
		ActionRequired:    mem.ActionRequired,
		Tags:              mem.SemanticTags,
		TriggerPhrases:    mem.TriggerPhrases,
	}

	if err := e.sqlStore.CreateMemory(sqlMemory); err != nil {
		return fmt.Errorf("failed to store memory in SQLite: %w", err)
	}

	// Store in vector database
	metadata := map[string]interface{}{
		"project_id":         mem.ProjectID,
		"importance":         mem.Importance,
		"context_type":       string(mem.ContextType),
		"temporal_relevance": string(mem.TemporalRelevance),
		"action_required":    mem.ActionRequired,
		"tags":               mem.SemanticTags,
		"trigger_phrases":    mem.TriggerPhrases,
		"created_at":         mem.CreatedAt.Unix(),
	}

	if err := e.vectorStore.Store(mem.ID, mem.Content, embedding, metadata); err != nil {
		return fmt.Errorf("failed to store memory in vector database: %w", err)
	}

	mem.CreatedAt = time.Now()
	mem.UpdatedAt = mem.CreatedAt

	return nil
}

// GetMemory retrieves a memory by ID
func (e *Engine) GetMemory(id string) (*Memory, error) {
	sqlMemory, err := e.sqlStore.GetMemory(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory: %w", err)
	}
	if sqlMemory == nil {
		return nil, nil
	}

	return e.sqlMemoryToMemory(sqlMemory), nil
}

// SearchMemories searches for relevant memories
func (e *Engine) SearchMemories(query *SearchQuery) ([]*SearchResult, error) {
	// Generate embedding for query
	queryEmbedding, err := e.embedder.Embed(query.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Build filters
	filters := map[string]interface{}{
		"project_id": query.ProjectID,
	}
	if query.MinImportance > 0 {
		filters["importance_gte"] = query.MinImportance
	}

	// Search vector database
	limit := query.Limit
	if limit == 0 {
		limit = 5
	}

	vectorResults, err := e.vectorStore.Search(queryEmbedding, limit*2, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to search vector database: %w", err)
	}

	// Convert to search results and score
	var results []*SearchResult
	for _, vr := range vectorResults {
		// Get full memory from SQLite
		mem, err := e.GetMemory(vr.ID)
		if err != nil {
			continue
		}
		if mem == nil {
			continue
		}

		// Calculate similarity score (1 - normalized distance)
		similarityScore := 1.0 - vr.Distance

		// Check for trigger phrase matches
		triggerMatched := e.checkTriggerMatch(query.Query, mem.TriggerPhrases)

		// Calculate relevance score
		relevanceScore := e.calculateRelevanceScore(mem, similarityScore, triggerMatched)

		results = append(results, &SearchResult{
			Memory:          mem,
			SimilarityScore: similarityScore,
			RelevanceScore:  relevanceScore,
			TriggerMatched:  triggerMatched,
		})
	}

	// Sort by relevance score
	sortByRelevance(results)

	// Limit results
	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// GetOrCreateProject gets or creates a project based on path
func (e *Engine) GetOrCreateProject(name string, path string) (*storage.Project, error) {
	// Try to get existing project
	project, err := e.sqlStore.GetProjectByPath(path)
	if err != nil {
		return nil, err
	}

	// Create if doesn't exist
	if project == nil {
		project = &storage.Project{
			ID:   uuid.New().String(),
			Name: name,
			Path: path,
		}
		if err := e.sqlStore.CreateProject(project); err != nil {
			return nil, err
		}
	}

	return project, nil
}

// CreateSession creates a new session
func (e *Engine) CreateSession(projectID string) (*storage.Session, error) {
	session := &storage.Session{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		StartedAt: time.Now(),
	}

	if err := e.sqlStore.CreateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

// EndSession ends a session
func (e *Engine) EndSession(sessionID string) error {
	session, err := e.sqlStore.GetSession(sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	now := time.Now()
	session.EndedAt = &now
	duration := int(now.Sub(session.StartedAt).Seconds())
	session.DurationSeconds = &duration

	return e.sqlStore.UpdateSession(session)
}

// GetSessionPrimer generates a session primer for context injection
func (e *Engine) GetSessionPrimer(projectID string) (*SessionPrimer, error) {
	project, err := e.sqlStore.GetProject(projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found: %s", projectID)
	}

	primer := &SessionPrimer{
		ProjectName: project.Name,
	}

	// Get last session
	lastSession, err := e.sqlStore.GetLastSession(projectID)
	if err != nil {
		return nil, err
	}

	if lastSession != nil && lastSession.EndedAt != nil {
		primer.LastSessionDate = lastSession.EndedAt
		timeSince := time.Since(*lastSession.EndedAt)
		primer.TimeSinceLastSession = formatDuration(timeSince)
	}

	// Get top memories (high importance, recent)
	topMemories, err := e.SearchMemories(&SearchQuery{
		Query:         project.Name, // Use project name as general query
		ProjectID:     projectID,
		Limit:         3,
		MinImportance: 0.7,
	})
	if err == nil && len(topMemories) > 0 {
		for _, result := range topMemories {
			primer.TopMemories = append(primer.TopMemories, result.Memory)
		}
	}

	return primer, nil
}

// Helper functions

func (e *Engine) sqlMemoryToMemory(sqlMem *storage.Memory) *Memory {
	mem := &Memory{
		ID:             sqlMem.ID,
		ProjectID:      sqlMem.ProjectID,
		Content:        sqlMem.Content,
		Importance:     sqlMem.Importance,
		SemanticTags:   sqlMem.Tags,
		TriggerPhrases: sqlMem.TriggerPhrases,
		ActionRequired: sqlMem.ActionRequired,
		CreatedAt:      sqlMem.CreatedAt,
		UpdatedAt:      sqlMem.UpdatedAt,
	}

	if sqlMem.SessionID != nil {
		mem.SessionID = *sqlMem.SessionID
	}
	if sqlMem.ContextType != nil {
		mem.ContextType = ContextType(*sqlMem.ContextType)
	}
	if sqlMem.TemporalRelevance != nil {
		mem.TemporalRelevance = TemporalRelevance(*sqlMem.TemporalRelevance)
	}

	return mem
}

func (e *Engine) checkTriggerMatch(query string, triggers []string) bool {
	// TODO: Implement sophisticated trigger matching
	// For now, simple substring match
	queryLower := toLower(query)
	for _, trigger := range triggers {
		if contains(queryLower, toLower(trigger)) {
			return true
		}
	}
	return false
}

func (e *Engine) calculateRelevanceScore(mem *Memory, similarity float64, triggerMatched bool) float64 {
	score := similarity * 0.6     // Base semantic similarity (60%)
	score += mem.Importance * 0.3 // Importance weight (30%)

	if triggerMatched {
		score += 0.2 // Trigger match boost (20%)
	}

	// Boost for action required
	if mem.ActionRequired {
		score += 0.1
	}

	// Normalize to 0-1
	if score > 1.0 {
		score = 1.0
	}

	return score
}

func sortByRelevance(results []*SearchResult) {
	// Simple bubble sort for now
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].RelevanceScore > results[i].RelevanceScore {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func toLower(s string) string {
	// Simple ASCII lowercase
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

func contains(haystack, needle string) bool {
	if len(needle) > len(haystack) {
		return false
	}
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		mins := int(d.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	days := int(d.Hours() / 24)
	if days == 1 {
		return "1 day ago"
	}
	if days < 7 {
		return fmt.Sprintf("%d days ago", days)
	}
	weeks := days / 7
	if weeks == 1 {
		return "1 week ago"
	}
	if weeks < 4 {
		return fmt.Sprintf("%d weeks ago", weeks)
	}
	months := days / 30
	if months == 1 {
		return "1 month ago"
	}
	return fmt.Sprintf("%d months ago", months)
}
