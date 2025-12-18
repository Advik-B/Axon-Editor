package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

// LSPMessage represents a JSON-RPC 2.0 message
type LSPMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *LSPError       `json:"error,omitempty"`
}

// LSPError represents a JSON-RPC error
type LSPError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// LSPClient manages JSON-RPC communication with gopls
type LSPClient struct {
	stdin      io.WriteCloser
	stdout     io.ReadCloser
	ctx        context.Context
	nextID     atomic.Int32
	pending    map[int]chan LSPMessage
	pendingMux sync.RWMutex
}

// NewLSPClient creates a new LSP client
func NewLSPClient(ctx context.Context, stdin io.WriteCloser, stdout io.ReadCloser) *LSPClient {
	client := &LSPClient{
		stdin:   stdin,
		stdout:  stdout,
		ctx:     ctx,
		pending: make(map[int]chan LSPMessage),
	}
	
	// Start reading responses
	go client.readLoop()
	
	return client
}

// readLoop continuously reads messages from gopls
func (c *LSPClient) readLoop() {
	reader := bufio.NewReader(c.stdout)
	
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}
		
		// Read Content-Length header
		var contentLength int
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Printf("LSP read error: %v\n", err)
				}
				return
			}
			
			if line == "\r\n" {
				break
			}
			
			var length int
			if _, err := fmt.Sscanf(line, "Content-Length: %d\r\n", &length); err == nil {
				contentLength = length
			}
		}
		
		if contentLength == 0 {
			continue
		}
		
		// Read message body
		body := make([]byte, contentLength)
		if _, err := io.ReadFull(reader, body); err != nil {
			fmt.Printf("LSP read body error: %v\n", err)
			return
		}
		
		// Parse message
		var msg LSPMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			fmt.Printf("LSP parse error: %v\n", err)
			continue
		}
		
		// Handle response
		if msg.ID > 0 {
			c.pendingMux.RLock()
			ch, ok := c.pending[msg.ID]
			c.pendingMux.RUnlock()
			
			if ok {
				ch <- msg
				close(ch)
				
				c.pendingMux.Lock()
				delete(c.pending, msg.ID)
				c.pendingMux.Unlock()
			}
		}
	}
}

// Call sends a request to gopls and waits for response
func (c *LSPClient) Call(method string, params interface{}) (json.RawMessage, error) {
	id := int(c.nextID.Add(1))
	
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}
	
	msg := LSPMessage{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  paramsJSON,
	}
	
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	
	// Create response channel
	respChan := make(chan LSPMessage, 1)
	c.pendingMux.Lock()
	c.pending[id] = respChan
	c.pendingMux.Unlock()
	
	// Send message
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(msgJSON))
	if _, err := c.stdin.Write([]byte(header)); err != nil {
		return nil, fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := c.stdin.Write(msgJSON); err != nil {
		return nil, fmt.Errorf("failed to write message: %w", err)
	}
	
	// Wait for response
	select {
	case resp := <-respChan:
		if resp.Error != nil {
			return nil, fmt.Errorf("LSP error %d: %s", resp.Error.Code, resp.Error.Message)
		}
		return resp.Result, nil
	case <-c.ctx.Done():
		return nil, c.ctx.Err()
	}
}

// Notify sends a notification to gopls (no response expected)
func (c *LSPClient) Notify(method string, params interface{}) error {
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}
	
	msg := LSPMessage{
		JSONRPC: "2.0",
		Method:  method,
		Params:  paramsJSON,
	}
	
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(msgJSON))
	if _, err := c.stdin.Write([]byte(header)); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := c.stdin.Write(msgJSON); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	
	return nil
}

// Initialize sends the initialize request to gopls
func (c *LSPClient) Initialize(rootURI string) error {
	params := map[string]interface{}{
		"processId": nil,
		"rootUri":   rootURI,
		"capabilities": map[string]interface{}{
			"textDocument": map[string]interface{}{
				"publishDiagnostics": map[string]bool{
					"relatedInformation": true,
				},
			},
		},
	}
	
	_, err := c.Call("initialize", params)
	if err != nil {
		return fmt.Errorf("initialize failed: %w", err)
	}
	
	// Send initialized notification
	return c.Notify("initialized", map[string]interface{}{})
}

// DidOpen notifies gopls that a document was opened
func (c *LSPClient) DidOpen(uri, languageID, text string) error {
	params := map[string]interface{}{
		"textDocument": map[string]interface{}{
			"uri":        uri,
			"languageId": languageID,
			"version":    1,
			"text":       text,
		},
	}
	
	return c.Notify("textDocument/didOpen", params)
}

// DidChange notifies gopls that a document changed
func (c *LSPClient) DidChange(uri, text string, version int) error {
	params := map[string]interface{}{
		"textDocument": map[string]interface{}{
			"uri":     uri,
			"version": version,
		},
		"contentChanges": []map[string]interface{}{
			{"text": text},
		},
	}
	
	return c.Notify("textDocument/didChange", params)
}

// Hover requests hover information for a position
func (c *LSPClient) Hover(uri string, line, character int) (string, error) {
	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": uri,
		},
		"position": map[string]int{
			"line":      line,
			"character": character,
		},
	}
	
	result, err := c.Call("textDocument/hover", params)
	if err != nil {
		return "", err
	}
	
	var hover struct {
		Contents struct {
			Kind  string `json:"kind"`
			Value string `json:"value"`
		} `json:"contents"`
	}
	
	if err := json.Unmarshal(result, &hover); err != nil {
		return "", fmt.Errorf("failed to parse hover: %w", err)
	}
	
	return hover.Contents.Value, nil
}

// Close shuts down the LSP client
func (c *LSPClient) Close() error {
	c.Notify("shutdown", nil)
	c.Notify("exit", nil)
	return c.stdin.Close()
}
