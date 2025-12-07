# Quick Start Guide

Get alaala running in under 5 minutes!

## Prerequisites

- Go 1.21+ installed
- Docker installed (for Weaviate)
- **Either:**
  - Anthropic API key (for cloud AI), **OR**
  - Ollama installed (for local AI) - [ollama.ai](https://ollama.ai)

## Step 1: Clone and Build

```bash
git clone https://github.com/georgepagarigan/alaala.git
cd alaala
go build -o bin/alaala ./cmd/alaala
```

## Step 2: Setup Weaviate

```bash
./scripts/setup-weaviate.sh
```

This will start a Weaviate container on `http://localhost:8080`.

## Step 3: Configure

Choose your AI provider:

### Option A: Anthropic Claude (Cloud-based)

```bash
# Set your Anthropic API key
export ANTHROPIC_API_KEY="sk-ant-..."
```

### Option B: Ollama (Local, Private)

```bash
# Install Ollama from https://ollama.ai
# Then pull the models
ollama pull llama3.1
ollama pull nomic-embed-text

# Edit ~/.alaala/config.yaml and set:
# ai:
#   provider: ollama
#   model: llama3.1
# embeddings:
#   provider: ollama
#   model: nomic-embed-text
```

### Initialize Your Project

```bash
# Initialize your project
cd /path/to/your/coding/project
~/alaala/bin/alaala init
```

This creates `.alaala-project.json` and `~/.alaala/config.yaml`.

## Step 4: Configure MCP in Cursor

1. Open Cursor Settings
2. Go to Features > Model Context Protocol
3. Add this configuration:

```json
{
  "mcpServers": {
    "alaala": {
      "command": "/Users/yourusername/alaala/bin/alaala",
      "args": ["serve"],
      "env": {
        "ANTHROPIC_API_KEY": "sk-ant-..."
      }
    }
  }
}
```

**Replace `/Users/yourusername/alaala/bin/alaala` with your actual path!**

## Step 5: Restart Cursor

Close and reopen Cursor. Alaala should now be running in the background!

## Testing It Works

In Cursor, try these commands:

1. **Check if alaala is running:**
   ```
   Can you list your available tools?
   ```
   You should see `search_memories`, `save_memory`, etc.

2. **Save a memory:**
   ```
   Remember that I prefer functional programming style
   ```

3. **Search for memories:**
   ```
   What do you remember about my coding preferences?
   ```

## Troubleshooting

### "Weaviate not accessible"

```bash
# Check if Weaviate is running
docker ps | grep weaviate

# If not, start it
docker start weaviate
```

### "Failed to curate memories"

Check your API key:
```bash
echo $ANTHROPIC_API_KEY
```

### "Project not found"

Make sure you ran `alaala init` in your project directory:
```bash
cd /path/to/your/project
alaala init
```

### MCP server not starting

Check Cursor logs or stderr output. Common issues:
- Incorrect binary path
- Missing API key
- Weaviate not running

## What's Next?

- Read the full [README.md](README.md)
- Check [STATUS.md](STATUS.md) for known limitations
- See [CONTRIBUTING.md](CONTRIBUTING.md) to help improve alaala
- Star the repo if you find it useful! ⭐

## Current Limitations

**Note:** This is an early version with some limitations:

- ❌ Vector search uses dummy embeddings (needs real implementation)
- ❌ Semantic similarity doesn't work properly yet
- ❌ No web UI yet
- ✅ But: MCP integration, AI curation, and storage all work!

See [STATUS.md](STATUS.md) for full details.

