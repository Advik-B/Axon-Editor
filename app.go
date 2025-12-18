package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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

	if err := os.WriteFile(filePath, prettyJSON, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
