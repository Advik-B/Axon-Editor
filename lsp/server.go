package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// AxonLSP represents the Language Server for Axon that uses gopls underneath
type AxonLSP struct {
	ctx         context.Context
	graphs      map[string]*AxonGraph
	graphMutex  sync.RWMutex
	diagnostics map[string][]Diagnostic
	goplsClient *GoplsClient
	workDir     string
}

// GoplsClient manages communication with gopls
type GoplsClient struct {
	cmd    *exec.Cmd
	stdin  *os.File
	stdout *os.File
	stderr *os.File
}

// AxonGraph represents an Axon graph for analysis
type AxonGraph struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Imports   []string    `json:"imports"`
	Nodes     []AxonNode  `json:"nodes"`
	DataEdges []DataEdge  `json:"data_edges"`
	ExecEdges []ExecEdge  `json:"exec_edges"`
}

// AxonNode represents a node in the graph
type AxonNode struct {
	ID            string            `json:"id"`
	Type          string            `json:"type"`
	Label         string            `json:"label"`
	Inputs        []Port            `json:"inputs,omitempty"`
	Outputs       []Port            `json:"outputs,omitempty"`
	ImplReference string            `json:"impl_reference,omitempty"`
	Config        map[string]string `json:"config,omitempty"`
}

// Port represents a data port on a node
type Port struct {
	Name     string `json:"name"`
	TypeName string `json:"type_name"`
}

// DataEdge represents a data connection
type DataEdge struct {
	FromNodeID string `json:"from_node_id"`
	FromPort   string `json:"from_port"`
	ToNodeID   string `json:"to_node_id"`
	ToPort     string `json:"to_port"`
}

// ExecEdge represents an execution connection
type ExecEdge struct {
	FromNodeID string `json:"from_node_id"`
	ToNodeID   string `json:"to_node_id"`
}

// Diagnostic represents an error or warning
type Diagnostic struct {
	NodeID   string `json:"node_id"`
	Severity string `json:"severity"` // "error", "warning", "info"
	Message  string `json:"message"`
	Code     string `json:"code,omitempty"`
	Line     int    `json:"line,omitempty"`
	Column   int    `json:"column,omitempty"`
}

// NewAxonLSP creates a new LSP server instance with gopls integration
func NewAxonLSP(ctx context.Context, workDir string) (*AxonLSP, error) {
	lsp := &AxonLSP{
		ctx:         ctx,
		graphs:      make(map[string]*AxonGraph),
		diagnostics: make(map[string][]Diagnostic),
		workDir:     workDir,
	}

	// Initialize gopls client
	client, err := NewGoplsClient(ctx, workDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gopls: %w", err)
	}
	lsp.goplsClient = client

	return lsp, nil
}

// NewGoplsClient starts a gopls language server instance
func NewGoplsClient(ctx context.Context, workDir string) (*GoplsClient, error) {
	// Check if gopls is available
	goplsPath, err := exec.LookPath("gopls")
	if err != nil {
		return nil, fmt.Errorf("gopls not found in PATH. Please install: go install golang.org/x/tools/gopls@latest")
	}

	cmd := exec.CommandContext(ctx, goplsPath, "-mode=stdio")
	cmd.Dir = workDir

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start gopls: %w", err)
	}

	return &GoplsClient{
		cmd:    cmd,
		stdin:  stdin.(*os.File),
		stdout: stdout.(*os.File),
		stderr: stderr.(*os.File),
	}, nil
}

// Close stops the gopls client
func (gc *GoplsClient) Close() error {
	if gc.cmd != nil && gc.cmd.Process != nil {
		return gc.cmd.Process.Kill()
	}
	return nil
}

// GenerateGoCode generates Go code from an Axon graph using the actual Axon transpiler
func (lsp *AxonLSP) GenerateGoCode(graph *AxonGraph) (string, error) {
	// Save the graph to a temporary .ax file
	tmpDir := filepath.Join(lsp.workDir, "tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create tmp directory: %w", err)
	}

	axFile := filepath.Join(tmpDir, fmt.Sprintf("%s.ax", graph.ID))
	
	// Marshal graph to JSON
	graphJSON, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal graph: %w", err)
	}
	
	// Write .ax file
	if err := os.WriteFile(axFile, graphJSON, 0644); err != nil {
		return "", fmt.Errorf("failed to write .ax file: %w", err)
	}

	// Check if Axon CLI is available
	axonPath, err := exec.LookPath("axon")
	if err != nil {
		// Axon CLI not found, try to find it in common locations or fall back to simple generation
		return lsp.generateGoCodeFallback(graph)
	}

	// Run Axon transpiler: axon build <file.ax>
	// The Axon CLI typically outputs the generated Go code to stdout or a specific directory
	cmd := exec.CommandContext(lsp.ctx, axonPath, "build", axFile)
	cmd.Dir = tmpDir
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("axon transpiler failed: %w\nOutput: %s", err, string(output))
	}

	// The Axon CLI typically creates an output directory with the generated code
	// Read the generated main.go file
	outputFile := filepath.Join(tmpDir, "out", "main.go")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		// Try alternative output location
		outputFile = filepath.Join(tmpDir, "main.go")
	}

	generatedCode, err := os.ReadFile(outputFile)
	if err != nil {
		// If we can't read the output file, return the transpiler output as-is
		// It might contain the generated code or error messages
		return string(output), nil
	}

	return string(generatedCode), nil
}

// generateGoCodeFallback is a simplified fallback when Axon CLI is not available
func (lsp *AxonLSP) generateGoCodeFallback(graph *AxonGraph) (string, error) {
	code := "package main\n\n"
	code += "// WARNING: Generated using fallback transpiler\n"
	code += "// Install Axon CLI for proper transpilation: go install github.com/Advik-B/Axon@latest\n\n"

	// Add imports
	if len(graph.Imports) > 0 {
		code += "import (\n"
		for _, imp := range graph.Imports {
			code += fmt.Sprintf("\t\"%s\"\n", imp)
		}
		code += ")\n\n"
	}

	code += "func main() {\n"
	code += "\t// Generated from Axon graph: " + graph.Name + "\n"
	code += "\t// This is a simplified fallback - use Axon CLI for full transpilation\n\n"
	
	// Simple node-to-code generation (placeholder)
	for _, node := range graph.Nodes {
		switch node.Type {
		case "START":
			code += "\t// START\n"
		case "END":
			code += "\t// END\n"
		case "CONSTANT":
			if val, ok := node.Config["value"]; ok {
				code += fmt.Sprintf("\t%s := %s\n", node.Label, val)
			}
		case "FUNCTION":
			if node.ImplReference != "" {
				code += fmt.Sprintf("\t%s()\n", node.ImplReference)
			}
		case "OPERATOR":
			if op, ok := node.Config["op"]; ok {
				code += fmt.Sprintf("\t// Operator: %s\n", op)
			}
		}
	}

	code += "}\n"
	return code, nil
}

// ValidateGraphWithGopls validates an Axon graph using both Axon rules and gopls
func (lsp *AxonLSP) ValidateGraphWithGopls(graphJSON string) ([]Diagnostic, string, error) {
	var graph AxonGraph
	if err := json.Unmarshal([]byte(graphJSON), &graph); err != nil {
		return nil, "", fmt.Errorf("invalid JSON: %w", err)
	}

	lsp.graphMutex.Lock()
	lsp.graphs[graph.ID] = &graph
	lsp.graphMutex.Unlock()

	// First, perform Axon-specific validation
	axonDiagnostics := lsp.validateAxonRules(&graph)

	// Generate Go code from the graph
	goCode, err := lsp.GenerateGoCode(&graph)
	if err != nil {
		return axonDiagnostics, "", fmt.Errorf("failed to generate Go code: %w", err)
	}

	// Write generated code to a temporary file for gopls analysis
	tmpDir := filepath.Join(lsp.workDir, "tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return axonDiagnostics, goCode, fmt.Errorf("failed to create tmp directory: %w", err)
	}

	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("%s.go", graph.ID))
	if err := os.WriteFile(tmpFile, []byte(goCode), 0644); err != nil {
		return axonDiagnostics, goCode, fmt.Errorf("failed to write tmp file: %w", err)
	}

	// TODO: Use gopls to validate the generated Go code
	// This would involve sending LSP requests to gopls and parsing responses
	// For now, we return Axon diagnostics and the generated code

	lsp.graphMutex.Lock()
	lsp.diagnostics[graph.ID] = axonDiagnostics
	lsp.graphMutex.Unlock()

	return axonDiagnostics, goCode, nil
}

// validateAxonRules performs Axon-specific validation
func (lsp *AxonLSP) validateAxonRules(graph *AxonGraph) []Diagnostic {
	diagnostics := []Diagnostic{}

	// Validate START node exists
	hasStart := false
	hasEnd := false
	for _, node := range graph.Nodes {
		if node.Type == "START" {
			hasStart = true
		}
		if node.Type == "END" {
			hasEnd = true
		}
	}

	if !hasStart {
		diagnostics = append(diagnostics, Diagnostic{
			NodeID:   "",
			Severity: "error",
			Message:  "Graph must have a START node",
			Code:     "missing_start",
		})
	}

	if !hasEnd {
		diagnostics = append(diagnostics, Diagnostic{
			NodeID:   "",
			Severity: "error",
			Message:  "Graph must have an END node",
			Code:     "missing_end",
		})
	}

	// Build node lookup map
	nodeMap := make(map[string]*AxonNode)
	for i := range graph.Nodes {
		nodeMap[graph.Nodes[i].ID] = &graph.Nodes[i]
	}

	// Validate data edges
	for _, edge := range graph.DataEdges {
		fromNode, fromExists := nodeMap[edge.FromNodeID]
		toNode, toExists := nodeMap[edge.ToNodeID]

		if !fromExists {
			diagnostics = append(diagnostics, Diagnostic{
				NodeID:   edge.FromNodeID,
				Severity: "error",
				Message:  fmt.Sprintf("Data edge references non-existent source node: %s", edge.FromNodeID),
				Code:     "invalid_edge",
			})
			continue
		}

		if !toExists {
			diagnostics = append(diagnostics, Diagnostic{
				NodeID:   edge.ToNodeID,
				Severity: "error",
				Message:  fmt.Sprintf("Data edge references non-existent target node: %s", edge.ToNodeID),
				Code:     "invalid_edge",
			})
			continue
		}

		// Validate port exists
		fromPortExists := false
		var fromPortType string
		for _, port := range fromNode.Outputs {
			if port.Name == edge.FromPort {
				fromPortExists = true
				fromPortType = port.TypeName
				break
			}
		}

		if !fromPortExists {
			diagnostics = append(diagnostics, Diagnostic{
				NodeID:   edge.FromNodeID,
				Severity: "error",
				Message:  fmt.Sprintf("Node %s does not have output port: %s", edge.FromNodeID, edge.FromPort),
				Code:     "invalid_port",
			})
		}

		toPortExists := false
		var toPortType string
		for _, port := range toNode.Inputs {
			if port.Name == edge.ToPort {
				toPortExists = true
				toPortType = port.TypeName
				break
			}
		}

		if !toPortExists {
			diagnostics = append(diagnostics, Diagnostic{
				NodeID:   edge.ToNodeID,
				Severity: "error",
				Message:  fmt.Sprintf("Node %s does not have input port: %s", edge.ToNodeID, edge.ToPort),
				Code:     "invalid_port",
			})
		}

		// Type checking using Go type system rules
		if fromPortExists && toPortExists {
			if !lsp.isTypeCompatible(fromPortType, toPortType) {
				diagnostics = append(diagnostics, Diagnostic{
					NodeID:   edge.ToNodeID,
					Severity: "warning",
					Message:  fmt.Sprintf("Type mismatch: connecting %s to %s", fromPortType, toPortType),
					Code:     "type_mismatch",
				})
			}
		}
	}

	// Validate execution edges
	for _, edge := range graph.ExecEdges {
		if _, exists := nodeMap[edge.FromNodeID]; !exists {
			diagnostics = append(diagnostics, Diagnostic{
				NodeID:   edge.FromNodeID,
				Severity: "error",
				Message:  fmt.Sprintf("Execution edge references non-existent source node: %s", edge.FromNodeID),
				Code:     "invalid_edge",
			})
		}

		if _, exists := nodeMap[edge.ToNodeID]; !exists {
			diagnostics = append(diagnostics, Diagnostic{
				NodeID:   edge.ToNodeID,
				Severity: "error",
				Message:  fmt.Sprintf("Execution edge references non-existent target node: %s", edge.ToNodeID),
				Code:     "invalid_edge",
			})
		}
	}

	// Check for unreachable nodes
	reachable := lsp.findReachableNodes(graph, nodeMap)
	for _, node := range graph.Nodes {
		if node.Type != "START" && node.Type != "END" && !reachable[node.ID] {
			diagnostics = append(diagnostics, Diagnostic{
				NodeID:   node.ID,
				Severity: "warning",
				Message:  fmt.Sprintf("Node %s is unreachable from START", node.Label),
				Code:     "unreachable_node",
			})
		}
	}

	// Check for execution cycles
	if lsp.hasCycle(graph) {
		diagnostics = append(diagnostics, Diagnostic{
			NodeID:   "",
			Severity: "error",
			Message:  "Execution graph contains cycles",
			Code:     "execution_cycle",
		})
	}

	return diagnostics
}

// isTypeCompatible checks if two Go types are compatible
func (lsp *AxonLSP) isTypeCompatible(fromType, toType string) bool {
	// Any type is compatible with anything
	if fromType == "any" || toType == "any" {
		return true
	}
	
	// interface{} is compatible with anything
	if fromType == "interface{}" || toType == "interface{}" {
		return true
	}
	
	// Exact match
	if fromType == toType {
		return true
	}
	
	// TODO: Add more sophisticated type compatibility checking
	// This could use gopls's type information
	
	return false
}

// findReachableNodes performs BFS to find all nodes reachable from START
func (lsp *AxonLSP) findReachableNodes(graph *AxonGraph, nodeMap map[string]*AxonNode) map[string]bool {
	reachable := make(map[string]bool)
	queue := []string{}

	// Find START node
	for _, node := range graph.Nodes {
		if node.Type == "START" {
			queue = append(queue, node.ID)
			reachable[node.ID] = true
			break
		}
	}

	// BFS
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, edge := range graph.ExecEdges {
			if edge.FromNodeID == current && !reachable[edge.ToNodeID] {
				reachable[edge.ToNodeID] = true
				queue = append(queue, edge.ToNodeID)
			}
		}
	}

	return reachable
}

// hasCycle detects cycles in the execution graph using DFS
func (lsp *AxonLSP) hasCycle(graph *AxonGraph) bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var dfs func(nodeID string) bool
	dfs = func(nodeID string) bool {
		visited[nodeID] = true
		recStack[nodeID] = true

		for _, edge := range graph.ExecEdges {
			if edge.FromNodeID == nodeID {
				if !visited[edge.ToNodeID] {
					if dfs(edge.ToNodeID) {
						return true
					}
				} else if recStack[edge.ToNodeID] {
					return true
				}
			}
		}

		recStack[nodeID] = false
		return false
	}

	for _, node := range graph.Nodes {
		if !visited[node.ID] {
			if dfs(node.ID) {
				return true
			}
		}
	}

	return false
}

// GetCompletions returns available node types and Go-aware suggestions
func (lsp *AxonLSP) GetCompletions(context string) []CompletionItem {
	completions := []CompletionItem{
		{Label: "CONSTANT", Kind: "node", Detail: "Constant value node", InsertText: "CONSTANT"},
		{Label: "FUNCTION", Kind: "node", Detail: "Function call node", InsertText: "FUNCTION"},
		{Label: "OPERATOR", Kind: "node", Detail: "Operator node (+, -, *, /)", InsertText: "OPERATOR"},
		{Label: "IGNORE", Kind: "node", Detail: "Ignore/discard value", InsertText: "IGNORE"},
		{Label: "RETURN", Kind: "node", Detail: "Return value", InsertText: "RETURN"},
	}

	// Add Go standard library function completions
	if context == "impl_reference" {
		completions = append(completions, []CompletionItem{
			{Label: "fmt.Println", Kind: "function", Detail: "Print to console", InsertText: "fmt.Println"},
			{Label: "fmt.Printf", Kind: "function", Detail: "Formatted print", InsertText: "fmt.Printf"},
			{Label: "fmt.Sprintf", Kind: "function", Detail: "Format string", InsertText: "fmt.Sprintf"},
			{Label: "os.ReadFile", Kind: "function", Detail: "Read file contents", InsertText: "os.ReadFile"},
			{Label: "os.WriteFile", Kind: "function", Detail: "Write file contents", InsertText: "os.WriteFile"},
			{Label: "strings.ToUpper", Kind: "function", Detail: "Convert to uppercase", InsertText: "strings.ToUpper"},
			{Label: "strings.ToLower", Kind: "function", Detail: "Convert to lowercase", InsertText: "strings.ToLower"},
			{Label: "strings.Split", Kind: "function", Detail: "Split string", InsertText: "strings.Split"},
			{Label: "strings.Join", Kind: "function", Detail: "Join strings", InsertText: "strings.Join"},
		}...)
	}

	// TODO: Integrate with gopls to get Go-aware completions from the generated code

	return completions
}

// CompletionItem represents a code completion suggestion
type CompletionItem struct {
	Label      string `json:"label"`
	Kind       string `json:"kind"`
	Detail     string `json:"detail"`
	InsertText string `json:"insert_text"`
}

// GetHoverInfo returns information about a node with Go type info
func (lsp *AxonLSP) GetHoverInfo(graphID, nodeID string) (string, error) {
	lsp.graphMutex.RLock()
	graph, exists := lsp.graphs[graphID]
	lsp.graphMutex.RUnlock()

	if !exists {
		return "", fmt.Errorf("graph not found: %s", graphID)
	}

	for _, node := range graph.Nodes {
		if node.ID == nodeID {
			info := fmt.Sprintf("**%s** (%s)\n\n", node.Label, node.Type)
			
			if len(node.Inputs) > 0 {
				info += "**Inputs:**\n"
				for _, port := range node.Inputs {
					info += fmt.Sprintf("- %s: `%s`\n", port.Name, port.TypeName)
				}
			}
			
			if len(node.Outputs) > 0 {
				info += "**Outputs:**\n"
				for _, port := range node.Outputs {
					info += fmt.Sprintf("- %s: `%s`\n", port.Name, port.TypeName)
				}
			}
			
			if node.ImplReference != "" {
				info += fmt.Sprintf("\n**Implementation:** `%s`\n", node.ImplReference)
				// TODO: Use gopls to get detailed documentation for the Go function
			}
			
			if len(node.Config) > 0 {
				info += "\n**Configuration:**\n"
				for k, v := range node.Config {
					info += fmt.Sprintf("- %s: %s\n", k, v)
				}
			}
			
			return info, nil
		}
	}

	return "", fmt.Errorf("node not found: %s", nodeID)
}

// Close shuts down the LSP server and gopls client
func (lsp *AxonLSP) Close() error {
	if lsp.goplsClient != nil {
		return lsp.goplsClient.Close()
	}
	return nil
}
