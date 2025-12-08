package memory

import (
	"fmt"

	"github.com/0xGurg/alaala/internal/ai"
	"github.com/google/uuid"
)

// Curator handles AI-powered memory curation
type Curator struct {
	engine   *Engine
	aiClient AIClient
}

// AIClient is an interface for AI-powered curation
type AIClient interface {
	CurateMemories(req *ai.CurationRequest) (*ai.CurationResponse, error)
}

// NewCurator creates a new curator
func NewCurator(engine *Engine, aiClient AIClient) *Curator {
	return &Curator{
		engine:   engine,
		aiClient: aiClient,
	}
}

// CurateSession curates memories from a session transcript
func (c *Curator) CurateSession(projectID, sessionID, transcript string) (*CurationResponse, error) {
	// Call AI to extract memories
	aiReq := &ai.CurationRequest{
		Transcript: transcript,
		ProjectID:  projectID,
		SessionID:  sessionID,
	}

	aiResp, err := c.aiClient.CurateMemories(aiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to curate memories with AI: %w", err)
	}

	// Convert AI memories to our memory format and store them
	var memories []*Memory
	memoryIDs := make([]string, len(aiResp.Memories))

	for i, curatedMem := range aiResp.Memories {
		mem := &Memory{
			ID:                uuid.New().String(),
			ProjectID:         projectID,
			SessionID:         sessionID,
			Content:           curatedMem.Content,
			Importance:        curatedMem.Importance,
			SemanticTags:      curatedMem.SemanticTags,
			ContextType:       ContextType(curatedMem.ContextType),
			TriggerPhrases:    curatedMem.TriggerPhrases,
			QuestionTypes:     curatedMem.QuestionTypes,
			TemporalRelevance: TemporalRelevance(curatedMem.TemporalRelevance),
			ActionRequired:    curatedMem.ActionRequired,
			Reasoning:         curatedMem.Reasoning,
		}

		// Store memory
		if err := c.engine.CreateMemory(mem); err != nil {
			return nil, fmt.Errorf("failed to store memory: %w", err)
		}

		memories = append(memories, mem)
		memoryIDs[i] = mem.ID
	}

	// Store relationships
	var relationships []struct {
		FromID string
		ToID   string
		Type   RelationshipType
	}

	for _, rel := range aiResp.Relationships {
		if rel.FromIndex < 0 || rel.FromIndex >= len(memoryIDs) ||
			rel.ToIndex < 0 || rel.ToIndex >= len(memoryIDs) {
			continue // Invalid indices
		}

		fromID := memoryIDs[rel.FromIndex]
		toID := memoryIDs[rel.ToIndex]
		relType := RelationshipType(rel.Type)

		// TODO: Store relationship in database
		// For now, just add to response
		relationships = append(relationships, struct {
			FromID string
			ToID   string
			Type   RelationshipType
		}{
			FromID: fromID,
			ToID:   toID,
			Type:   relType,
		})
	}

	return &CurationResponse{
		Memories:      memories,
		Relationships: relationships,
		Summary:       aiResp.Summary,
	}, nil
}
