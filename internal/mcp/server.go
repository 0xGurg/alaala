package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/georgepagarigan/alaala/internal/memory"
)

// Server implements the MCP (Model Context Protocol) server
type Server struct {
	engine   *memory.Engine
	curator  *memory.Curator
	reader   *bufio.Reader
	writer   io.Writer
	handlers map[string]RequestHandler
}

// RequestHandler handles MCP requests
type RequestHandler func(params json.RawMessage) (interface{}, error)

// NewServer creates a new MCP server
func NewServer(engine *memory.Engine, curator *memory.Curator) *Server {
	server := &Server{
		engine:   engine,
		curator:  curator,
		reader:   bufio.NewReader(os.Stdin),
		writer:   os.Stdout,
		handlers: make(map[string]RequestHandler),
	}

	server.registerHandlers()
	return server
}

// registerHandlers registers all MCP request handlers
func (s *Server) registerHandlers() {
	// Tool handlers
	s.handlers["tools/list"] = s.handleListTools
	s.handlers["tools/call"] = s.handleCallTool
	
	// Resource handlers
	s.handlers["resources/list"] = s.handleListResources
	s.handlers["resources/read"] = s.handleReadResource
	
	// Prompt handlers
	s.handlers["prompts/list"] = s.handleListPrompts
	s.handlers["prompts/get"] = s.handleGetPrompt
	
	// Server info
	s.handlers["initialize"] = s.handleInitialize
}

// Run starts the MCP server
func (s *Server) Run() error {
	fmt.Fprintln(os.Stderr, "MCP server started, waiting for requests...")

	for {
		// Read JSON-RPC request from stdin
		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read request: %w", err)
		}

		// Parse request
		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.sendError(nil, -32700, "Parse error", err)
			continue
		}

		// Handle request
		s.handleRequest(&req)
	}

	return nil
}

// handleRequest processes a single JSON-RPC request
func (s *Server) handleRequest(req *JSONRPCRequest) {
	handler, ok := s.handlers[req.Method]
	if !ok {
		s.sendError(req.ID, -32601, "Method not found", nil)
		return
	}

	result, err := handler(req.Params)
	if err != nil {
		s.sendError(req.ID, -32603, "Internal error", err)
		return
	}

	s.sendResult(req.ID, result)
}

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(params json.RawMessage) (interface{}, error) {
	return map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools":     map[string]bool{},
			"resources": map[string]bool{},
			"prompts":   map[string]bool{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "alaala",
			"version": "0.1.0",
		},
	}, nil
}

// sendResult sends a successful JSON-RPC response
func (s *Server) sendResult(id interface{}, result interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

	s.sendResponse(&resp)
}

// sendError sends an error JSON-RPC response
func (s *Server) sendError(id interface{}, code int, message string, data interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	s.sendResponse(&resp)
}

// sendResponse sends a JSON-RPC response
func (s *Server) sendResponse(resp *JSONRPCResponse) {
	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal response: %v\n", err)
		return
	}

	fmt.Fprintf(s.writer, "%s\n", data)
}

// JSON-RPC types

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC 2.0 error
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

