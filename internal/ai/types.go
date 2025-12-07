package ai

// CurationRequest represents a request to curate memories
type CurationRequest struct {
	Transcript string
	ProjectID  string
	SessionID  string
}

// CurationResponse represents the AI's curated memories
type CurationResponse struct {
	Memories      []CuratedMemory       `json:"memories"`
	Relationships []MemoryRelationship  `json:"relationships"`
	Summary       string                `json:"summary"`
}

// CuratedMemory represents a memory extracted by the AI
type CuratedMemory struct {
	Content           string   `json:"content"`
	Importance        float64  `json:"importance_weight"`
	SemanticTags      []string `json:"semantic_tags"`
	ContextType       string   `json:"context_type"`
	TriggerPhrases    []string `json:"trigger_phrases"`
	QuestionTypes     []string `json:"question_types"`
	TemporalRelevance string   `json:"temporal_relevance"`
	ActionRequired    bool     `json:"action_required"`
	Reasoning         string   `json:"reasoning"`
}

// MemoryRelationship represents a relationship between memories
type MemoryRelationship struct {
	FromIndex int    `json:"from_index"`
	ToIndex   int    `json:"to_index"`
	Type      string `json:"type"`
}

