package main

import (
	"context"
	"fmt"
	"github.com/Advik-B/Axon/parser" // Assuming Axon is a local module or in GOPATH
	"github.com/Advik-B/Axon/pkg/axon"
	"github.com/Advik-B/Axon/transpiler"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// LoadGraphFromFile uses the core parser to load a graph and returns the typed struct.
// Wails will automatically serialize this to JSON for the frontend.
func (a *App) LoadGraphFromFile(filePath string) (*axon.Graph, error) {
	fmt.Printf("Loading graph from: %s\n", filePath)
	graph, err := parser.LoadGraphFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load graph: %w", err)
	}
	return graph, nil
}

// SaveGraphToFile takes a graph object from the frontend and uses the core writer.
func (a *App) SaveGraphToFile(filePath string, graph *axon.Graph) error {
	fmt.Printf("Saving graph to: %s\n", filePath)
	return parser.SaveGraphToFile(graph, filePath)
}

// TranspileGraph uses the core transpiler to generate Go code from the graph object.
func (a *App) TranspileGraph(graph *axon.Graph) (string, error) {
	fmt.Println("Transpiling graph in memory...")
	code, err := transpiler.Transpile(graph)
	if err != nil {
		return "", fmt.Errorf("transpilation failed: %w", err)
	}
	return code, nil
}