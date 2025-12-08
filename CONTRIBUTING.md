# Contributing to alaala

Thank you for your interest in contributing to alaala! We welcome contributions from the community.

## Development Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/georgepagarigan/alaala.git
   cd alaala
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Setup Weaviate:**
   ```bash
   docker run -d \
     --name weaviate \
     -p 8080:8080 \
     -e AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED=true \
     -e PERSISTENCE_DATA_PATH=/var/lib/weaviate \
     weaviate/weaviate:latest
   ```

4. **Build:**
   ```bash
   go build -o bin/alaala ./cmd/alaala
   ```

5. **Set environment variables:**
   ```bash
   export ANTHROPIC_API_KEY="your-api-key"
   ```

## Project Structure

```
alaala/
├── cmd/alaala/          # Main application entry point
├── internal/
│   ├── mcp/             # MCP protocol implementation
│   ├── memory/          # Core memory engine
│   ├── storage/         # Database implementations (SQLite, Weaviate)
│   ├── ai/              # AI client (Claude)
│   ├── embeddings/      # Embedding service
│   └── web/             # Web UI
├── pkg/config/          # Configuration management
├── examples/            # Example configurations
└── scripts/             # Utility scripts
```

## Development Guidelines

### Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Use meaningful variable names
- Add comments for complex logic
- Keep functions small and focused

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/memory
```

### Building

```bash
# Build for current platform
go build -o bin/alaala ./cmd/alaala

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o bin/alaala-linux-amd64 ./cmd/alaala
GOOS=darwin GOARCH=arm64 go build -o bin/alaala-darwin-arm64 ./cmd/alaala
GOOS=darwin GOARCH=amd64 go build -o bin/alaala-darwin-amd64 ./cmd/alaala
```

## Areas for Contribution

### High Priority

1. **Real Embeddings Implementation**
   - Current implementation uses dummy embeddings
   - Need to integrate actual sentence-transformers
   - Options: ONNX runtime, or HTTP service

2. **Web UI**
   - Design is spec'd out (neobrutalism + Kanagawa palette)
   - Needs implementation with Go templates + HTMX
   - Memory browser, graph visualization, analytics

3. **Proper Weaviate GraphQL Implementation**
   - Current vector search is stubbed out
   - Need to properly use Weaviate Go client v4 API
   - Handle GraphQL field selection correctly

### Medium Priority

4. **Memory Graph Expansion**
   - Implement graph traversal for related memories
   - Visualization in web UI
   - Query optimization

5. **Export/Import**
   - JSON export format
   - Markdown export
   - Import validation

6. **Ollama Integration** (Documented, needs implementation)
   - HTTP client for Ollama API
   - Embeddings via Ollama (nomic-embed-text)
   - Curation via Ollama (llama3.1, mistral, etc.)
   - Configuration already documented in examples/config.yaml

### Low Priority

7. **Testing**
   - Unit tests for all packages
   - Integration tests
   - Benchmark tests

8. **Documentation**
   - API documentation
   - Architecture deep-dive
   - Usage examples

## Submitting Changes

1. **Fork the repository**

2. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes:**
   - Write clear, concise commit messages
   - Follow the code style guidelines
   - Add tests if applicable
   - Update documentation

4. **Test your changes:**
   ```bash
   go test ./...
   go build ./cmd/alaala
   ```

5. **Commit your changes:**
   ```bash
   git commit -am "Add feature: your feature description"
   ```

6. **Push to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request:**
   - Describe your changes clearly
   - Reference any related issues
   - Wait for review

## Commit Message Guidelines

- Use present tense ("Add feature" not "Added feature")
- Use imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit first line to 72 characters
- Reference issues and pull requests when relevant

Examples:
```
Add real embeddings implementation using ONNX
Fix circular import in storage package
Implement memory graph traversal
Update README with new installation instructions
```

## Code Review Process

1. All submissions require review
2. Reviews may request changes
3. Once approved, maintainers will merge

## Questions?

Feel free to open an issue for:
- Questions about contributing
- Feature requests
- Bug reports
- General discussion

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

