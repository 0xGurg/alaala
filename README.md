# alaala 🧠

> _"alaala" (Tagalog for "memory") - A semantic memory system for AI assistants_

A high-performance Go implementation of a semantic memory system that enables AI assistants to maintain context across sessions using the Model Context Protocol (MCP). Built with Weaviate for vector search, SQLite for metadata, and Claude AI for intelligent memory curation.

## ✨ Features

- **MCP Protocol Integration** - Works seamlessly with Cursor, Claude Desktop, and other MCP-compatible clients
- **Hybrid Memory Injection** - Auto-inject context at session start + dynamic updates on each prompt + on-demand searches
- **AI-Powered Curation** - Claude or Ollama analyzes conversation transcripts to extract meaningful insights
- **Local AI Support** - Use Ollama for completely private, local memory curation and embeddings
- **Memory Graph** - Memories can reference and relate to each other (references, supersedes, related_to)
- **Multi-Project Workspaces** - Automatic project isolation with separate memory spaces
- **Semantic Search** - Vector similarity search with importance weighting and trigger phrase matching
- **Export/Import** - Backup and share memories in JSON format
- **Web UI** (Coming Soon) - Beautiful neobrutalist interface with Kanagawa color palette

## 🚀 Quick Start

### Installation

#### Option 1: Download Binary (Recommended)

```bash
# macOS (ARM64)
curl -L https://github.com/georgepagarigan/alaala/releases/latest/download/alaala-darwin-arm64 -o alaala
chmod +x alaala
sudo mv alaala /usr/local/bin/

# macOS (AMD64)
curl -L https://github.com/georgepagarigan/alaala/releases/latest/download/alaala-darwin-amd64 -o alaala
chmod +x alaala
sudo mv alaala /usr/local/bin/

# Linux
curl -L https://github.com/georgepagarigan/alaala/releases/latest/download/alaala-linux-amd64 -o alaala
chmod +x alaala
sudo mv alaala /usr/local/bin/
```

#### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/georgepagarigan/alaala.git
cd alaala

# Build
go build -o bin/alaala ./cmd/alaala

# Install (optional)
sudo mv bin/alaala /usr/local/bin/
```

### Setup Weaviate (Required for Vector Search)

#### Using Docker (Recommended)

```bash
# Start Weaviate container
docker run -d \
  --name weaviate \
  -p 8080:8080 \
  -e QUERY_DEFAULTS_LIMIT=25 \
  -e AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED=true \
  -e PERSISTENCE_DATA_PATH='/var/lib/weaviate' \
  weaviate/weaviate:latest
```

#### Or use embedded mode (experimental)

Set `storage.mode: embedded` in your config file.

### Configuration

1. **Initialize your first project:**

```bash
cd /path/to/your/project
alaala init
```

This creates `.alaala-project.json` in your project directory.

2. **Create configuration file:**

The default config is created at `~/.alaala/config.yaml`. Customize it:

```yaml
storage:
  mode: docker  # or "embedded"
  weaviate:
    docker_url: http://localhost:8080
  sqlite_path: ~/.alaala/alaala.db

ai:
  provider: anthropic  # or "ollama" for local AI
  api_key: ${ANTHROPIC_API_KEY}  # not needed for ollama
  model: claude-3-5-sonnet-20241022  # or "llama3.1" for ollama
  ollama_url: http://localhost:11434  # if using ollama

embeddings:
  provider: local  # or "ollama" for local embeddings
  model: all-MiniLM-L6-v2  # or "nomic-embed-text" for ollama
  ollama_url: http://localhost:11434  # if using ollama

retrieval:
  max_memories: 5
  min_importance: 0.3
  include_graph_depth: 1

web:
  enabled: true
  port: 8766
  host: localhost

logging:
  level: info
  file: ~/.alaala/alaala.log
```

3. **Set your AI provider:**

**Option A: Using Anthropic Claude (Cloud)**
```bash
export ANTHROPIC_API_KEY="your-api-key-here"
```

**Option B: Using Ollama (Local)**
```bash
# Install Ollama from https://ollama.ai
ollama pull llama3.1
ollama pull nomic-embed-text
# Update config.yaml to use provider: ollama
```

### MCP Configuration

#### For Cursor

Add this to your Cursor settings (Cursor Settings > Features > Model Context Protocol):

```json
{
  "mcpServers": {
    "alaala": {
      "command": "/usr/local/bin/alaala",
      "args": ["serve"],
      "env": {
        "ANTHROPIC_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

#### For Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "alaala": {
      "command": "/usr/local/bin/alaala",
      "args": ["serve"]
    }
  }
}
```

## 📖 Usage

### Basic Commands

```bash
# Start MCP server (for Cursor/Claude Desktop)
alaala serve

# Start web UI
alaala web

# Initialize project
alaala init

# Show version
alaala version
```

### Using with Cursor

Once configured, alaala runs automatically in the background. The AI can:

1. **Automatic Context Injection** - Relevant memories are injected at session start and updated on each prompt
2. **Search Memories** - Use the `search_memories` tool:
   ```
   Search memories about authentication
   ```
3. **Save Important Insights** - Use the `save_memory` tool:
   ```
   Remember that I prefer JWT tokens over session cookies
   ```
4. **Curate Sessions** - After a conversation, the AI can call `curate_session` to extract key insights

### MCP Tools Available

| Tool | Description | Example |
|------|-------------|---------|
| `search_memories` | Search for relevant memories | Find memories about "database schema" |
| `save_memory` | Manually save a memory | Save "Project uses PostgreSQL 15" |
| `curate_session` | Extract memories from transcript | Analyze this conversation |
| `list_projects` | List all projects | Show all my projects |

### MCP Resources

| Resource | Description |
|----------|-------------|
| `memory://session-context` | Current session context with relevant memories |
| `memory://project-memories` | All memories for the current project |

## 🏗️ Architecture

```
alaala/
├── cmd/alaala/          # CLI entry point
├── internal/
│   ├── mcp/             # MCP protocol server
│   ├── memory/          # Core memory engine
│   ├── storage/         # SQLite + Weaviate
│   ├── ai/              # Claude AI client
│   ├── embeddings/      # Embedding service
│   └── web/             # Web UI (coming soon)
├── pkg/config/          # Configuration
└── examples/            # Example configs
```

### How It Works

1. **Session Start** - alaala injects a session primer with:
   - Last session timestamp
   - Top relevant memories
   - Unresolved items

2. **During Conversation** - On each prompt:
   - Dynamic memory resource updates with relevant context
   - AI can search for specific memories
   - AI can save important insights

3. **Session End** - Optionally:
   - AI analyzes full transcript
   - Extracts meaningful memories with metadata
   - Creates relationship graph between memories

### Memory Structure

Each memory contains:

```go
{
  content: "User prefers functional programming style",
  importance: 0.9,
  semanticTags: ["preference", "coding-style"],
  contextType: "PREFERENCE",
  triggerPhrases: ["coding style", "how to write code"],
  questionTypes: ["what style does user prefer"],
  temporalRelevance: "persistent",
  actionRequired: false,
  reasoning: "Important for future code suggestions"
}
```

## 🎨 Upcoming Features

- [x] Core memory engine with MCP
- [x] SQLite + Weaviate integration
- [x] AI-powered curation (Claude + Ollama)
- [x] **Ollama support** for local AI (documented, needs implementation)
- [ ] **Web UI** with neobrutalism design
- [ ] **Real embeddings** (currently using dummy embeddings)
- [ ] **Memory graph visualization**
- [ ] **Export/import functionality**
- [ ] **Cross-project search**
- [ ] **Homebrew formula** for easy installation

## 🛠️ Development

```bash
# Clone repository
git clone https://github.com/georgepagarigan/alaala.git
cd alaala

# Install dependencies
go mod download

# Run tests (coming soon)
go test ./...

# Build
go build -o bin/alaala ./cmd/alaala

# Run
./bin/alaala serve
```

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📝 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- Inspired by [RLabs-Inc/memory](https://github.com/RLabs-Inc/memory)
- Built with [Weaviate](https://weaviate.io/) for vector search
- Powered by [Claude](https://www.anthropic.com/) for AI curation
- Designed with [Kanagawa](https://github.com/rebelot/kanagawa.nvim) color palette

## 📞 Support

- Issues: [GitHub Issues](https://github.com/georgepagarigan/alaala/issues)
- Documentation: [Full docs](https://github.com/georgepagarigan/alaala/tree/main/docs)

---

Made with ❤️ by [George Pagarigan](https://github.com/georgepagarigan)
