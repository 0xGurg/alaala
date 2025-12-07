package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore handles SQLite operations for metadata storage
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLite store
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	store := &SQLiteStore{db: db}

	// Initialize schema
	if err := store.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// initSchema creates the database schema
func (s *SQLiteStore) initSchema() error {
	schema := `
	-- Projects table
	CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		path TEXT NOT NULL UNIQUE,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	-- Sessions table
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		started_at DATETIME NOT NULL,
		ended_at DATETIME,
		duration_seconds INTEGER,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
	);

	-- Memories table (metadata only, vectors in Weaviate)
	CREATE TABLE IF NOT EXISTS memories (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		session_id TEXT,
		content TEXT NOT NULL,
		importance REAL NOT NULL DEFAULT 0.5,
		context_type TEXT,
		temporal_relevance TEXT,
		action_required BOOLEAN DEFAULT FALSE,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
		FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE SET NULL
	);

	-- Memory tags (many-to-many)
	CREATE TABLE IF NOT EXISTS memory_tags (
		memory_id TEXT NOT NULL,
		tag TEXT NOT NULL,
		PRIMARY KEY (memory_id, tag),
		FOREIGN KEY (memory_id) REFERENCES memories(id) ON DELETE CASCADE
	);

	-- Memory trigger phrases
	CREATE TABLE IF NOT EXISTS memory_triggers (
		memory_id TEXT NOT NULL,
		phrase TEXT NOT NULL,
		PRIMARY KEY (memory_id, phrase),
		FOREIGN KEY (memory_id) REFERENCES memories(id) ON DELETE CASCADE
	);

	-- Memory relationships (graph)
	CREATE TABLE IF NOT EXISTS memory_relationships (
		from_memory_id TEXT NOT NULL,
		to_memory_id TEXT NOT NULL,
		relationship_type TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		PRIMARY KEY (from_memory_id, to_memory_id, relationship_type),
		FOREIGN KEY (from_memory_id) REFERENCES memories(id) ON DELETE CASCADE,
		FOREIGN KEY (to_memory_id) REFERENCES memories(id) ON DELETE CASCADE
	);

	-- Indexes for performance
	CREATE INDEX IF NOT EXISTS idx_memories_project ON memories(project_id);
	CREATE INDEX IF NOT EXISTS idx_memories_session ON memories(session_id);
	CREATE INDEX IF NOT EXISTS idx_memories_importance ON memories(importance);
	CREATE INDEX IF NOT EXISTS idx_memories_created ON memories(created_at);
	CREATE INDEX IF NOT EXISTS idx_sessions_project ON sessions(project_id);
	CREATE INDEX IF NOT EXISTS idx_sessions_started ON sessions(started_at);
	`

	_, err := s.db.Exec(schema)
	return err
}

// Project represents a project in the database
type Project struct {
	ID        string
	Name      string
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Session represents a session in the database
type Session struct {
	ID              string
	ProjectID       string
	StartedAt       time.Time
	EndedAt         *time.Time
	DurationSeconds *int
}

// Memory represents memory metadata in the database
type Memory struct {
	ID                 string
	ProjectID          string
	SessionID          *string
	Content            string
	Importance         float64
	ContextType        *string
	TemporalRelevance  *string
	ActionRequired     bool
	Tags               []string
	TriggerPhrases     []string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// MemoryRelationship represents a relationship between memories
type MemoryRelationship struct {
	FromMemoryID     string
	ToMemoryID       string
	RelationshipType string
	CreatedAt        time.Time
}

// CreateProject creates a new project
func (s *SQLiteStore) CreateProject(project *Project) error {
	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO projects (id, name, path, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, project.ID, project.Name, project.Path, project.CreatedAt, project.UpdatedAt)

	return err
}

// GetProject retrieves a project by ID
func (s *SQLiteStore) GetProject(id string) (*Project, error) {
	var project Project
	err := s.db.QueryRow(`
		SELECT id, name, path, created_at, updated_at
		FROM projects WHERE id = ?
	`, id).Scan(&project.ID, &project.Name, &project.Path, &project.CreatedAt, &project.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &project, nil
}

// GetProjectByPath retrieves a project by path
func (s *SQLiteStore) GetProjectByPath(path string) (*Project, error) {
	var project Project
	err := s.db.QueryRow(`
		SELECT id, name, path, created_at, updated_at
		FROM projects WHERE path = ?
	`, path).Scan(&project.ID, &project.Name, &project.Path, &project.CreatedAt, &project.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &project, nil
}

// CreateSession creates a new session
func (s *SQLiteStore) CreateSession(session *Session) error {
	_, err := s.db.Exec(`
		INSERT INTO sessions (id, project_id, started_at, ended_at, duration_seconds)
		VALUES (?, ?, ?, ?, ?)
	`, session.ID, session.ProjectID, session.StartedAt, session.EndedAt, session.DurationSeconds)

	return err
}

// UpdateSession updates a session
func (s *SQLiteStore) UpdateSession(session *Session) error {
	_, err := s.db.Exec(`
		UPDATE sessions 
		SET ended_at = ?, duration_seconds = ?
		WHERE id = ?
	`, session.EndedAt, session.DurationSeconds, session.ID)

	return err
}

// GetSession retrieves a session by ID
func (s *SQLiteStore) GetSession(id string) (*Session, error) {
	var session Session
	err := s.db.QueryRow(`
		SELECT id, project_id, started_at, ended_at, duration_seconds
		FROM sessions WHERE id = ?
	`, id).Scan(&session.ID, &session.ProjectID, &session.StartedAt, &session.EndedAt, &session.DurationSeconds)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetLastSession retrieves the most recent session for a project
func (s *SQLiteStore) GetLastSession(projectID string) (*Session, error) {
	var session Session
	err := s.db.QueryRow(`
		SELECT id, project_id, started_at, ended_at, duration_seconds
		FROM sessions 
		WHERE project_id = ? 
		ORDER BY started_at DESC 
		LIMIT 1
	`, projectID).Scan(&session.ID, &session.ProjectID, &session.StartedAt, &session.EndedAt, &session.DurationSeconds)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// CreateMemory creates a new memory with tags and trigger phrases
func (s *SQLiteStore) CreateMemory(memory *Memory) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	memory.CreatedAt = now
	memory.UpdatedAt = now

	// Insert memory
	_, err = tx.Exec(`
		INSERT INTO memories (id, project_id, session_id, content, importance, 
			context_type, temporal_relevance, action_required, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, memory.ID, memory.ProjectID, memory.SessionID, memory.Content, memory.Importance,
		memory.ContextType, memory.TemporalRelevance, memory.ActionRequired,
		memory.CreatedAt, memory.UpdatedAt)
	if err != nil {
		return err
	}

	// Insert tags
	for _, tag := range memory.Tags {
		_, err = tx.Exec(`INSERT INTO memory_tags (memory_id, tag) VALUES (?, ?)`, memory.ID, tag)
		if err != nil {
			return err
		}
	}

	// Insert trigger phrases
	for _, phrase := range memory.TriggerPhrases {
		_, err = tx.Exec(`INSERT INTO memory_triggers (memory_id, phrase) VALUES (?, ?)`, memory.ID, phrase)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetMemory retrieves a memory by ID with its tags and trigger phrases
func (s *SQLiteStore) GetMemory(id string) (*Memory, error) {
	var memory Memory
	err := s.db.QueryRow(`
		SELECT id, project_id, session_id, content, importance,
			context_type, temporal_relevance, action_required, created_at, updated_at
		FROM memories WHERE id = ?
	`, id).Scan(&memory.ID, &memory.ProjectID, &memory.SessionID, &memory.Content,
		&memory.Importance, &memory.ContextType, &memory.TemporalRelevance,
		&memory.ActionRequired, &memory.CreatedAt, &memory.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Load tags
	rows, err := s.db.Query(`SELECT tag FROM memory_tags WHERE memory_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		memory.Tags = append(memory.Tags, tag)
	}

	// Load trigger phrases
	rows, err = s.db.Query(`SELECT phrase FROM memory_triggers WHERE memory_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var phrase string
		if err := rows.Scan(&phrase); err != nil {
			return nil, err
		}
		memory.TriggerPhrases = append(memory.TriggerPhrases, phrase)
	}

	return &memory, nil
}

// CreateRelationship creates a relationship between two memories
func (s *SQLiteStore) CreateRelationship(rel *MemoryRelationship) error {
	rel.CreatedAt = time.Now()

	_, err := s.db.Exec(`
		INSERT INTO memory_relationships (from_memory_id, to_memory_id, relationship_type, created_at)
		VALUES (?, ?, ?, ?)
	`, rel.FromMemoryID, rel.ToMemoryID, rel.RelationshipType, rel.CreatedAt)

	return err
}

// GetRelationships retrieves all relationships for a memory
func (s *SQLiteStore) GetRelationships(memoryID string) ([]MemoryRelationship, error) {
	rows, err := s.db.Query(`
		SELECT from_memory_id, to_memory_id, relationship_type, created_at
		FROM memory_relationships
		WHERE from_memory_id = ? OR to_memory_id = ?
	`, memoryID, memoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relationships []MemoryRelationship
	for rows.Next() {
		var rel MemoryRelationship
		if err := rows.Scan(&rel.FromMemoryID, &rel.ToMemoryID, &rel.RelationshipType, &rel.CreatedAt); err != nil {
			return nil, err
		}
		relationships = append(relationships, rel)
	}

	return relationships, nil
}

