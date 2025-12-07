# Testing OpenRouter Integration

This guide helps you test alaala with different OpenRouter models.

## Prerequisites

1. Get an OpenRouter API key from [openrouter.ai](https://openrouter.ai)
2. Add credits to your account (most models cost < $0.01 per request)
3. Build alaala: `go build -o bin/alaala ./cmd/alaala`

## Basic Testing

### 1. Set API Key

```bash
export OPENROUTER_API_KEY="sk-or-v1-..."
```

### 2. Create Test Config

Create `~/.alaala/config.yaml`:

```yaml
storage:
  mode: docker
  weaviate:
    docker_url: http://localhost:8080
  sqlite_path: ~/.alaala/alaala.db

ai:
  provider: openrouter
  model: anthropic/claude-3.5-sonnet  # Start with this
  openrouter_url: https://openrouter.ai/api/v1

embeddings:
  provider: local
  model: all-MiniLM-L6-v2

retrieval:
  max_memories: 5
  min_importance: 0.3
  include_graph_depth: 1

web:
  enabled: true
  port: 8766
  host: localhost

logging:
  level: debug  # Use debug for testing
  file: ~/.alaala/alaala.log
```

### 3. Initialize Project

```bash
cd /tmp/test-alaala
alaala init
```

### 4. Test with Cursor/Claude Desktop

Add to your MCP configuration:

```json
{
  "mcpServers": {
    "alaala": {
      "command": "/path/to/alaala/bin/alaala",
      "args": ["serve"],
      "env": {
        "OPENROUTER_API_KEY": "sk-or-v1-..."
      }
    }
  }
}
```

## Testing Different Models

### Test 1: Claude 3.5 Sonnet (Best Quality)

**Config:**
```yaml
ai:
  provider: openrouter
  model: anthropic/claude-3.5-sonnet
```

**Expected:** High-quality memory extraction with detailed reasoning

**Cost:** ~$3 input / ~$15 output per 1M tokens

### Test 2: GPT-4 Turbo (Fast & Reliable)

**Config:**
```yaml
ai:
  provider: openrouter
  model: openai/gpt-4-turbo
```

**Expected:** Fast responses, good quality

**Cost:** ~$10 input / ~$30 output per 1M tokens

### Test 3: Llama 3.1 70B (Cost-Effective)

**Config:**
```yaml
ai:
  provider: openrouter
  model: meta-llama/llama-3.1-70b-instruct
```

**Expected:** Good quality at lower cost

**Cost:** ~$0.50 input / ~$0.80 output per 1M tokens

### Test 4: Gemini Pro 1.5 (Fast & Cheap)

**Config:**
```yaml
ai:
  provider: openrouter
  model: google/gemini-pro-1.5
```

**Expected:** Very fast, good for high-volume

**Cost:** ~$1.25 input / ~$5 output per 1M tokens

### Test 5: Mistral Large (European Alternative)

**Config:**
```yaml
ai:
  provider: openrouter
  model: mistralai/mistral-large
```

**Expected:** Good quality, EU-based

**Cost:** ~$2 input / ~$6 output per 1M tokens

## Error Handling Tests

### Test Rate Limiting

Make rapid successive requests to trigger rate limits. The client should:
- Automatically retry with exponential backoff
- Show clear error messages
- Eventually succeed

### Test Invalid API Key

```bash
export OPENROUTER_API_KEY="invalid-key"
```

**Expected Error:**
```
OpenRouter API error: Invalid API key

Please check your OPENROUTER_API_KEY environment variable
```

### Test Invalid Model

**Config:**
```yaml
ai:
  model: nonexistent/model
```

**Expected Error:**
```
OpenRouter API error: Model 'nonexistent/model' is not available

Try: anthropic/claude-3.5-sonnet, openai/gpt-4-turbo, or meta-llama/llama-3.1-70b-instruct
```

### Test Insufficient Credits

If you run out of credits:

**Expected Error:**
```
OpenRouter API error: Insufficient credits

Please add credits to your OpenRouter account
```

## Comparing Models

### Quality Test

Create a complex conversation transcript and compare memory extraction across models:

1. Use Claude 3.5 Sonnet (baseline)
2. Use GPT-4 Turbo
3. Use Llama 3.1 70B
4. Use Gemini Pro 1.5

**Compare:**
- Number of memories extracted
- Importance weights assigned
- Quality of semantic tags
- Accuracy of context types
- Reasoning quality

### Performance Test

Measure response times:

```bash
time alaala serve
# Then trigger curation requests
```

**Expected order (fastest to slowest):**
1. Gemini Pro 1.5 (~2-5s)
2. GPT-4 Turbo (~3-8s)
3. Llama 3.1 70B (~5-10s)
4. Claude 3.5 Sonnet (~5-12s)

### Cost Test

Track costs in OpenRouter dashboard:
- Check cost per request
- Calculate cost per memory extracted
- Find best value for your use case

## Manual Integration Test

### Scenario: Complete Session

1. Start a new coding session
2. Have a conversation about a technical topic
3. End the session and trigger curation
4. Verify memories were extracted correctly
5. Start a new session
6. Verify session primer includes relevant memories

### Expected Results

**After Curation:**
- 3-5 high-quality memories extracted
- Each with importance > 0.5
- Relevant semantic tags
- Appropriate context types
- Clear reasoning

**Session Primer:**
- Shows last session timestamp
- Includes top 3 relevant memories
- Formatted clearly

## Troubleshooting

### OpenRouter Connection Issues

```bash
# Test OpenRouter API directly
curl https://openrouter.ai/api/v1/models \
  -H "Authorization: Bearer $OPENROUTER_API_KEY"
```

### Check Logs

```bash
tail -f ~/.alaala/alaala.log
```

Look for:
- API request/response details
- Error messages
- Retry attempts

### Verify Configuration

```bash
cat ~/.alaala/config.yaml
```

Ensure:
- `provider: openrouter`
- Valid API key (or env var set)
- Correct model name
- Proper URL (optional)

## Success Criteria

✅ OpenRouter client successfully connects
✅ Memory curation works with all tested models
✅ Error messages are clear and actionable
✅ Retry logic handles rate limits
✅ Different models show appropriate quality/cost trade-offs
✅ Documentation is accurate

## Reporting Issues

If you find issues:

1. Check logs: `~/.alaala/alaala.log`
2. Verify API key: `echo $OPENROUTER_API_KEY`
3. Test with `anthropic/claude-3.5-sonnet` (most reliable)
4. Report with: model name, error message, and logs

## Next Steps

After testing:
- Choose your preferred model
- Set it in production config
- Monitor costs in OpenRouter dashboard
- Adjust model based on quality/cost needs

