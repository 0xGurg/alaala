package storage

// GraphTraverser handles memory relationship traversal
type GraphTraverser struct {
	sqlStore *SQLiteStore
}

// NewGraphTraverser creates a new graph traverser
func NewGraphTraverser(sqlStore *SQLiteStore) *GraphTraverser {
	return &GraphTraverser{
		sqlStore: sqlStore,
	}
}

// ExpandMemories performs BFS traversal of memory relationships
// Returns additional memory IDs to include, up to the specified depth
func (g *GraphTraverser) ExpandMemories(seedIDs []string, depth int) ([]string, error) {
	if depth == 0 || len(seedIDs) == 0 {
		return []string{}, nil
	}

	visited := make(map[string]bool)
	var result []string

	// Mark seed IDs as visited
	for _, id := range seedIDs {
		visited[id] = true
	}

	// BFS traversal
	currentLevel := seedIDs
	for currentDepth := 0; currentDepth < depth; currentDepth++ {
		if len(currentLevel) == 0 {
			break
		}

		var nextLevel []string

		// Get relationships for all IDs in current level
		for _, memID := range currentLevel {
			rels, err := g.sqlStore.GetRelationships(memID)
			if err != nil {
				continue // Skip on error, don't fail entire traversal
			}

			for _, rel := range rels {
				// Add related memory if not visited
				relatedID := rel.ToMemoryID
				if relatedID == memID {
					// This is an incoming relationship, use FromMemoryID
					relatedID = rel.FromMemoryID
				}

				if !visited[relatedID] {
					visited[relatedID] = true
					result = append(result, relatedID)
					nextLevel = append(nextLevel, relatedID)
				}
			}
		}

		currentLevel = nextLevel
	}

	return result, nil
}
