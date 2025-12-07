# Go Memory System with MCP Integration

## Overview

Create "alaala" (Tagalog for "memory") - a performant Go-based semantic memory system that enables AI assistants to maintain context across sessions using the Model Context Protocol (MCP).

## Architecture

### Core Components

**1. MCP Server (Go)**

- Implements MCP protocol (stdio transport)
- **Hybrid Memory Injection**:
  - Session start: Auto-inject session primer via MCP prompt
  - Every user prompt: Dynamic MCP resource provides relevant memories (like original hooks)
  - On-demand: AI can call search_memories() tool for deeper searches
- Handles project isolation and session management
- Built with Go for high performance and easy deployment

**2. Storage Layer**

- **Weaviate**: Vector embeddings + semantic search (configurable: embedded or Docker)
- **SQLite**: Metadata, sessions, projects, memory relationships
- Support for memory graphs (memories reference each other)

**3. AI Integration**

- Anthropic Claude API for memory curation
- Local embedding models (sentence-transformers compatible)
- Support for Ollama as fallback for embeddings

**4. Web UI**

- Simple dashboard for browsing memories
- Project management interface
- Memory graph visualization
- Export/import functionality

## Project Structure

```
alaala/
├── cmd/
│   └── alaala/
│       └── main.go                 # CLI entry point
├── internal/
│   ├── mcp/
│   │   ├── server.go              # MCP server implementation
│   │   ├── resources.go           # MCP resources (session context, etc.)
│   │   ├── tools.go               # MCP tools (search, save, curate)
│   │   └── prompts.go             # MCP prompts (session primer)
│   ├── memory/
│   │   ├── engine.go              # Core memory engine
│   │   ├── curator.go             # AI-powered curation
│   │   └── retrieval.go           # Smart retrieval strategies
│   ├── storage/
│   │   ├── weaviate.go            # Weaviate client
│   │   ├── sqlite.go              # SQLite operations
│   │   └── graph.go               # Memory relationship graph
│   ├── embeddings/
│   │   ├── client.go              # Embedding service interface
│   │   └── transformers.go        # Local model integration
│   ├── ai/
│   │   ├── claude.go              # Claude API client
│   │   └── types.go               # Curation types
│   └── web/
│       ├── server.go              # Web UI server
│       ├── handlers.go            # API handlers
│       └── templates/             # HTML templates
├── pkg/
│   └── config/
│       └── config.go              # Configuration management
├── web/
│   ├── static/                    # CSS, JS
│   └── templates/                 # HTML templates
├── scripts/
│   └── setup-weaviate.sh          # Docker setup script
├── examples/
│   └── cursor-mcp-config.json     # Example MCP config
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

## Implementation Plan

### Phase 1: Foundation

1. **Project Setup**

   - Initialize Go module as `alaala`
   - Setup dependencies: Weaviate Go client, SQLite driver, MCP SDK
   - Create basic project structure

2. **Storage Layer**

   - Implement Weaviate integration (embedded + Docker modes)
   - Setup SQLite schema (projects, sessions, memories, relationships)
   - Create storage abstraction layer

3. **Configuration**

   - YAML/JSON config file support
   - Environment variable overrides
   - Project detection (.alaala-project.json)

### Phase 2: Core Memory Engine

1. **Memory Types & Metadata**

   - Define memory structure (content, importance, tags, trigger_phrases, etc.)
   - Implement semantic embeddings integration
   - Create memory graph relationships

2. **Retrieval System**

   - Vector similarity search with Weaviate
   - Metadata filtering and scoring
   - Two-stage retrieval (obligatory + scored)
   - Trigger phrase matching

3. **Session Management**

   - Session tracking and temporal context
   - Last session analysis
   - Project isolation

### Phase 3: MCP Integration

1. **MCP Protocol Implementation**

   - Stdio transport server
   - Resource handlers (session-context, project memories)
   - Tool implementations (search, save, curate, export)
   - Prompt templates (session primer)

2. **AI Curation**

   - Claude API integration for transcript analysis
   - Memory extraction with structured metadata
   - Batch curation processing
   - Importance weighting and semantic tagging

### Phase 4: Enhanced Features

1. **Memory Graph**

   - Relationship types (references, supersedes, related_to)
   - Graph traversal for context expansion
   - Memory evolution tracking

2. **Multi-Project Workspaces**

   - Automatic project detection
   - Cross-project memory search (optional)
   - Project-specific configuration

3. **Export/Import**

   - JSON export format
   - Backup and restore functionality
   - Memory sharing between projects

4. **Web UI**

   - **Design**: Neobrutalism style (bold borders, chunky shadows, playful)
   - **Framework**: Tailwind CSS + shadcn/ui components
   - **Colors**: Kanagawa palette (fujiWhite, surimiOrange, autumnRed, waveBlue, etc.)
   - **Features**:
     - Memory browser with semantic search
     - Project switcher with stats
     - Interactive graph visualization (D3.js with Kanagawa theme)
     - Export/import interface
     - Analytics dashboard (memory count, importance distribution)
     - Responsive design with neobrutalist cards and buttons

### Phase 5: Deployment & Documentation

1. **Build & Distribution**

   - Cross-platform binaries (Darwin, Linux, Windows)
   - Homebrew formula (macOS)
   - Installation scripts
   - Docker image (optional)

2. **Documentation**

   - **README.md** with:
     - Project overview and philosophy
     - Installation methods (Homebrew, binary download, build from source)
     - Quick start guide (< 5 minutes to first memory)
     - MCP configuration for Cursor/Claude Desktop with examples
     - Basic usage examples (searching, saving, curating)
     - Web UI screenshots and walkthrough
     - Configuration options
     - Troubleshooting common issues
   - **USAGE.md** - Detailed usage patterns and best practices
   - **API.md** - MCP tools specification and API reference
   - **ARCHITECTURE.md** - Deep-dive into system design

## Key Files to Create

**Core Engine:**

- `internal/memory/engine.go` - Main memory engine with CRUD operations
- `internal/memory/curator.go` - AI-powered curation logic
- `internal/memory/retrieval.go` - Smart retrieval strategies

**MCP Layer:**

- `internal/mcp/server.go` - MCP protocol server
- `internal/mcp/tools.go` - Tool definitions and handlers
- `internal/mcp/resources.go` - Resource providers

**Storage:**

- `internal/storage/weaviate.go` - Weaviate client wrapper
- `internal/storage/sqlite.go` - SQLite operations
- `internal/storage/graph.go` - Memory graph logic

**AI Integration:**

- `internal/ai/claude.go` - Claude API client for curation

## Configuration Example

```yaml
# ~/.alaala/config.yaml
storage:
  mode: embedded  # embedded | docker
  weaviate:
    embedded_path: ~/.alaala/weaviate
    docker_url: http://localhost:8080
  sqlite_path: ~/.alaala/alaala.db

ai:
  provider: anthropic
  api_key: ${ANTHROPIC_API_KEY}
  model: claude-3-5-sonnet-20241022
  
embeddings:
  provider: local  # local | openai
  model: all-MiniLM-L6-v2

retrieval:
  max_memories: 5
  min_importance: 0.3
  include_graph_depth: 1  # Follow relationships 1 level deep

web:
  enabled: true
  port: 8766

logging:
  level: info
  file: ~/.alaala/alaala.log
```

## MCP Tools Specification

**search_memories**

- Input: `query` (string), `limit` (int), `project_id` (optional)
- Output: Array of memories with metadata
- Use: AI searches for relevant context

**save_memory**

- Input: `content` (string), `importance` (float), `tags` ([]string), `context_type` (string)
- Output: Memory ID
- Use: AI saves important insights during conversation

**curate_session**

- Input: `transcript` (string), `session_id` (string)
- Output: Array of curated memories
- Use: End-of-session batch curation

**list_projects**

- Output: Array of projects with stats
- Use: Project navigation

**export_memories**

- Input: `project_id` (string), `format` (json|markdown)
- Output: Exported data
- Use: Backup, sharing

## Success Criteria

- MCP server works with Cursor out of the box
- Memory retrieval < 100ms for typical queries
- AI curation produces high-quality, relevant memories
- Web UI provides intuitive memory browsing
- Single-binary deployment (no dependencies except optional Docker)
- Comprehensive documentation and examples

## Todo List

1. **setup-project** - Initialize Go module, dependencies, and project structure
2. **storage-layer** - Implement Weaviate and SQLite storage with graph support
3. **memory-engine** - Build core memory engine with retrieval and embeddings
4. **mcp-server** - Implement MCP protocol server with resources, tools, prompts
5. **ai-curation** - Integrate Claude API for transcript-based memory curation
6. **enhanced-features** - Add memory graph, multi-project, export/import features
7. **web-ui** - Build web dashboard for memory browsing and visualization
8. **deployment** - Create build scripts, binaries, and installation tools
9. **documentation** - Write comprehensive docs, examples, and MCP config guides

