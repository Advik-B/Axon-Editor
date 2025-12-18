package main

import (
	"axon-editor/lsp"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx       context.Context
	lspServer *lsp.AxonLSP
	workDir   string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// Setup work directory
	workDir, err := os.Getwd()
	if err != nil {
		workDir = os.TempDir()
	}
	a.workDir = workDir
	
	// Initialize LSP server with gopls integration
	lspServer, err := lsp.NewAxonLSP(ctx, workDir)
	if err != nil {
		runtime.LogWarning(ctx, fmt.Sprintf("Failed to initialize LSP server: %v", err))
		// Continue without LSP if gopls is not available
	} else {
		a.lspServer = lspServer
		runtime.LogInfo(ctx, "LSP server initialized with gopls support")
	}
}

// AxonGraph represents the complete Axon graph structure
type AxonGraph struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Imports   []string    `json:"imports"`
	Nodes     []AxonNode  `json:"nodes"`
	DataEdges []DataEdge  `json:"data_edges"`
	ExecEdges []ExecEdge  `json:"exec_edges"`
	Comments  []Comment   `json:"comments,omitempty"`
}

// AxonNode represents a node in the Axon graph
type AxonNode struct {
	ID            string            `json:"id"`
	Type          string            `json:"type"`
	Label         string            `json:"label"`
	Inputs        []Port            `json:"inputs,omitempty"`
	Outputs       []Port            `json:"outputs,omitempty"`
	ImplReference string            `json:"impl_reference,omitempty"`
	Config        map[string]string `json:"config,omitempty"`
	VisualInfo    *VisualInfo       `json:"visual_info,omitempty"`
}

// Port represents a connection point for data on a node
type Port struct {
	Name     string `json:"name"`
	TypeName string `json:"type_name"`
}

// DataEdge represents a data dependency between two nodes
type DataEdge struct {
	FromNodeID string `json:"from_node_id"`
	FromPort   string `json:"from_port"`
	ToNodeID   string `json:"to_node_id"`
	ToPort     string `json:"to_port"`
}

// ExecEdge represents an execution dependency between two nodes
type ExecEdge struct {
	FromNodeID string `json:"from_node_id"`
	ToNodeID   string `json:"to_node_id"`
}

// Comment stores documentation that can be attached to nodes
type Comment struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

// VisualInfo stores the position of a node
type VisualInfo struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width,omitempty"`
	Height float64 `json:"height,omitempty"`
}

// OpenFile opens a file dialog and loads an Axon graph
func (a *App) OpenFile() (string, error) {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Open Axon Graph",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Axon Files (*.ax)",
				Pattern:     "*.ax",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})
	if err != nil {
		return "", err
	}
	if filePath == "" {
		return "", fmt.Errorf("no file selected")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(data), nil
}

// SaveFile saves an Axon graph to a file
func (a *App) SaveFile(graphJSON string) error {
	filePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save Axon Graph",
		DefaultFilename: "graph.ax",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Axon Files (*.ax)",
				Pattern:     "*.ax",
			},
		},
	})
	if err != nil {
		return err
	}
	if filePath == "" {
		return fmt.Errorf("no file selected")
	}

	// Validate JSON before saving
	var graph AxonGraph
	if err := json.Unmarshal([]byte(graphJSON), &graph); err != nil {
		return fmt.Errorf("invalid graph JSON: %w", err)
	}

	// Pretty print JSON
	prettyJSON, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}

	if err := os.WriteFile(filePath, prettyJSON, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ValidateGraph validates an Axon graph using LSP
func (a *App) ValidateGraph(graphJSON string) ([]lsp.Diagnostic, error) {
	if a.lspServer == nil {
		return nil, fmt.Errorf("LSP server not initialized")
	}
	
	diagnostics, _, err := a.lspServer.ValidateGraphWithGopls(graphJSON)
	return diagnostics, err
}

// GenerateGoCode generates Go code from an Axon graph
func (a *App) GenerateGoCode(graphJSON string) (string, error) {
	var graph AxonGraph
	if err := json.Unmarshal([]byte(graphJSON), &graph); err != nil {
		return "", fmt.Errorf("invalid graph JSON: %w", err)
	}
	
	if a.lspServer == nil {
		return "", fmt.Errorf("LSP server not initialized")
	}
	
	// Convert to LSP graph format
	lspGraph := &lsp.AxonGraph{
		ID:        graph.ID,
		Name:      graph.Name,
		Imports:   graph.Imports,
		Nodes:     convertNodesToLSP(graph.Nodes),
		DataEdges: convertDataEdgesToLSP(graph.DataEdges),
		ExecEdges: convertExecEdgesToLSP(graph.ExecEdges),
	}
	
	return a.lspServer.GenerateGoCode(lspGraph)
}

// GetCompletions returns code completion suggestions
func (a *App) GetCompletions(context string) ([]lsp.CompletionItem, error) {
	if a.lspServer == nil {
		return nil, fmt.Errorf("LSP server not initialized")
	}
	
	return a.lspServer.GetCompletions(context), nil
}

// GetHoverInfo returns hover information for a node
func (a *App) GetHoverInfo(graphID, nodeID string) (string, error) {
	if a.lspServer == nil {
		return "", fmt.Errorf("LSP server not initialized")
	}
	
	return a.lspServer.GetHoverInfo(graphID, nodeID)
}

// BuildGraph builds the Axon graph using the actual Axon CLI transpiler
func (a *App) BuildGraph(graphJSON string) (string, error) {
	var graph AxonGraph
	if err := json.Unmarshal([]byte(graphJSON), &graph); err != nil {
		return "", fmt.Errorf("invalid graph JSON: %w", err)
	}
	
	// Create temporary build directory
	buildDir := filepath.Join(a.workDir, "build", "tmp")
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create build directory: %w", err)
	}
	
	// Write the graph as .ax file
	axFile := filepath.Join(buildDir, fmt.Sprintf("%s.ax", graph.ID))
	prettyJSON, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal graph: %w", err)
	}
	
	if err := os.WriteFile(axFile, prettyJSON, 0644); err != nil {
		return "", fmt.Errorf("failed to write .ax file: %w", err)
	}
	
	// Check if Axon CLI is available
	axonPath, err := exec.LookPath("axon")
	if err != nil {
		return "", fmt.Errorf("Axon CLI not found. Please install: go install github.com/Advik-B/Axon@latest")
	}
	
	// Run Axon build command
	cmd := exec.CommandContext(a.ctx, axonPath, "build", axFile)
	cmd.Dir = buildDir
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return string(output), fmt.Errorf("Axon build failed: %w\n%s", err, string(output))
	}
	
	// Now compile the generated Go code
	// The Axon CLI typically outputs to an "out" directory
	outDir := filepath.Join(buildDir, "out")
	if _, statErr := os.Stat(outDir); os.IsNotExist(statErr) {
		// Try current directory if "out" doesn't exist
		outDir = buildDir
	}
	
	// Determine executable name based on platform
	execName := "executable"
	if filepath.Ext(execName) == "" && os.Getenv("GOOS") == "windows" {
		execName += ".exe"
	}
	
	// Compile with Go
	goCmd := exec.CommandContext(a.ctx, "go", "build", "-o", execName, "./...")
	goCmd.Dir = outDir
	goBuildOutput, err := goCmd.CombinedOutput()
	
	combinedOutput := string(output) + "\n" + string(goBuildOutput)
	
	if err != nil {
		return combinedOutput, fmt.Errorf("Go compilation failed: %w\n%s", err, string(goBuildOutput))
	}
	
	return combinedOutput + "\nBuild successful!", nil
}

// RunGraph runs the built Axon graph
func (a *App) RunGraph(graphJSON string) (string, error) {
	var graph AxonGraph
	if err := json.Unmarshal([]byte(graphJSON), &graph); err != nil {
		return "", fmt.Errorf("invalid graph JSON: %w", err)
	}
	
	// Build first
	buildOutput, err := a.BuildGraph(graphJSON)
	if err != nil {
		return buildOutput, err
	}
	
	// Determine the executable location
	buildDir := filepath.Join(a.workDir, "build", "tmp")
	outDir := filepath.Join(buildDir, "out")
	if _, statErr := os.Stat(outDir); os.IsNotExist(statErr) {
		outDir = buildDir
	}
	
	// Determine executable name
	execName := "executable"
	if filepath.Ext(execName) == "" && os.Getenv("GOOS") == "windows" {
		execName += ".exe"
	}
	
	// Run the built executable
	cmd := exec.CommandContext(a.ctx, filepath.Join(outDir, execName))
	cmd.Dir = outDir
	output, err := cmd.CombinedOutput()
	
	fullOutput := buildOutput + "\n--- Program Output ---\n" + string(output)
	
	if err != nil {
		return fullOutput, fmt.Errorf("execution failed: %w\n%s", err, string(output))
	}
	
	return fullOutput, nil
}

// Helper functions to convert between app and LSP types
func convertNodesToLSP(nodes []AxonNode) []lsp.AxonNode {
	result := make([]lsp.AxonNode, len(nodes))
	for i, node := range nodes {
		result[i] = lsp.AxonNode{
			ID:            node.ID,
			Type:          node.Type,
			Label:         node.Label,
			Inputs:        convertPortsToLSP(node.Inputs),
			Outputs:       convertPortsToLSP(node.Outputs),
			ImplReference: node.ImplReference,
			Config:        node.Config,
		}
	}
	return result
}

func convertPortsToLSP(ports []Port) []lsp.Port {
	result := make([]lsp.Port, len(ports))
	for i, port := range ports {
		result[i] = lsp.Port{
			Name:     port.Name,
			TypeName: port.TypeName,
		}
	}
	return result
}

func convertDataEdgesToLSP(edges []DataEdge) []lsp.DataEdge {
	result := make([]lsp.DataEdge, len(edges))
	for i, edge := range edges {
		result[i] = lsp.DataEdge{
			FromNodeID: edge.FromNodeID,
			FromPort:   edge.FromPort,
			ToNodeID:   edge.ToNodeID,
			ToPort:     edge.ToPort,
		}
	}
	return result
}

func convertExecEdgesToLSP(edges []ExecEdge) []lsp.ExecEdge {
	result := make([]lsp.ExecEdge, len(edges))
	for i, edge := range edges {
		result[i] = lsp.ExecEdge{
			FromNodeID: edge.FromNodeID,
			ToNodeID:   edge.ToNodeID,
		}
	}
	return result
}
