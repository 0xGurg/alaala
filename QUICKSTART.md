# Quick Start Guide

Get alaala running in under 5 minutes!

## Prerequisites

- Docker installed (for Weaviate)
- **Choose one AI provider:**
  - Anthropic API key (Claude - best quality), **OR**
  - OpenRouter API key (multiple models - flexible), **OR**
  - Ollama installed (local AI - private) - [ollama.ai](https://ollama.ai)

## Step 1: Install alaala

### Option A: Download Binary (Recommended)

```bash
# macOS (ARM64)
curl -L https://github.com/0xGurg/alaala/releases/latest/download/alaala_darwin_arm64.tar.gz | tar xz
sudo mv alaala /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/0xGurg/alaala/releases/latest/download/alaala_darwin_amd64.tar.gz | tar xz
sudo mv alaala /usr/local/bin/

# Linux (AMD64)
curl -L https://github.com/0xGurg/alaala/releases/latest/download/alaala_linux_amd64.tar.gz | tar xz
sudo mv alaala /usr/local/bin/

# Linux (ARM64)
curl -L https://github.com/0xGurg/alaala/releases/latest/download/alaala_linux_arm64.tar.gz | tar xz
sudo mv alaala /usr/local/bin/

# Verify installation
alaala version
```

### Option B: Build from Source

```bash
git clone https://github.com/0xGurg/alaala.git
cd alaala
go build -o bin/alaala ./cmd/alaala
sudo mv bin/alaala /usr/local/bin/
```

## Step 2: Setup Weaviate

```bash
./scripts/setup-weaviate.sh
```

This will start a Weaviate container on `http://localhost:8080`.

## Step 3: Configure

Choose your AI provider:

### Option A: Anthropic Claude (Cloud, Best Quality)

```bash
# Set your Anthropic API key
export ANTHROPIC_API_KEY="sk-ant-..."
```

### Option B: OpenRouter (Cloud, Multiple Models)

```bash
# Get API key from https://openrouter.ai
export OPENROUTER_API_KEY="sk-or-v1-..."

# Edit ~/.alaala/config.yaml and set:
# ai:
#   provider: openrouter
#   model: anthropic/claude-3.5-sonnet
#   # Or try other models:
#   # model: openai/gpt-4-turbo
#   # model: meta-llama/llama-3.1-70b-instruct
#   # model: google/gemini-pro-1.5
```

### Option C: Ollama (Local, Private)

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
      "command": "/usr/local/bin/alaala",
      "args": ["serve"],
      "env": {
        "OPENROUTER_API_KEY": "sk-or-v1-..."
      }
    }
  }
}
```

**Note:** Change to `ANTHROPIC_API_KEY` if using Claude, or remove `env` entirely if using Ollama (local).

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

