package memory

import "time"

// ContextType represents the type of context for a memory
type ContextType string

const (
	ContextTypeTechnicalImplementation ContextType = "TECHNICAL_IMPLEMENTATION"
	ContextTypeArchitecture            ContextType = "ARCHITECTURE"
	ContextTypeDecision                ContextType = "DECISION"
	ContextTypeBreakthrough            ContextType = "BREAKTHROUGH"
	ContextTypeRelationship            ContextType = "RELATIONSHIP"
	ContextTypeUnresolved              ContextType = "UNRESOLVED"
	ContextTypeMilestone               ContextType = "MILESTONE"
	ContextTypePreference              ContextType = "PREFERENCE"
)

// TemporalRelevance represents how long a memory stays relevant
type TemporalRelevance string

const (
	TemporalRelevancePersistent TemporalRelevance = "persistent"
	TemporalRelevanceSession    TemporalRelevance = "session"
	TemporalRelevanceTemporary  TemporalRelevance = "temporary"
)

// RelationshipType represents the type of relationship between memories
type RelationshipType string

const (
	RelationshipTypeReferences RelationshipType = "references"
	RelationshipTypeSupersedes RelationshipType = "supersedes"
	RelationshipTypeRelatedTo  RelationshipType = "related_to"
	RelationshipTypeConflicts  RelationshipType = "conflicts"
	RelationshipTypeExpands    RelationshipType = "expands"
)

// Memory represents a complete memory with all its metadata
type Memory struct {
	ID                string
	ProjectID         string
	SessionID         string
	Content           string
	Importance        float64
	SemanticTags      []string
	ContextType       ContextType
	TriggerPhrases    []string
	QuestionTypes     []string
	TemporalRelevance TemporalRelevance
	ActionRequired    bool
	Reasoning         string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Relationships     []Relationship
}

// Relationship represents a connection between memories
type Relationship struct {
	ToMemoryID string
	Type       RelationshipType
	CreatedAt  time.Time
}

// SearchQuery represents a memory search request
type SearchQuery struct {
	Query             string
	ProjectID         string
	Limit             int
	MinImportance     float64
	ContextTypes      []ContextType
	IncludeGraphDepth int
}

// SearchResult represents a memory search result with scoring
type SearchResult struct {
	Memory          *Memory
	SimilarityScore float64
	RelevanceScore  float64
	TriggerMatched  bool
}

// SessionPrimer represents contextual information injected at session start
type SessionPrimer struct {
	ProjectName          string
	LastSessionDate      *time.Time
	TimeSinceLastSession string
	LastSessionSummary   string
	TopMemories          []*Memory
	UnresolvedItems      []*Memory
}

// CurationRequest represents a request to curate memories from a transcript
type CurationRequest struct {
	ProjectID  string
	SessionID  string
	Transcript string
}

// CurationResponse represents the result of memory curation
type CurationResponse struct {
	Memories      []*Memory
	Relationships []struct {
		FromID string
		ToID   string
		Type   RelationshipType
	}
	Summary string
}
