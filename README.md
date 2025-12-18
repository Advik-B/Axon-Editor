# Axon Editor

The official visual node graph editor for [Axon](https://github.com/Advik-B/Axon) - a visual, node-based programming language that transpiles to Go.

## Overview

Axon Editor is a desktop application built with Wails (Go + Web) and Svelte 5 that provides an intuitive interface for creating and editing Axon visual programs. It allows you to design complex programs using a node-based workflow and save them as `.ax` files that can be transpiled to Go code using the Axon CLI.

## Features

- **Visual Node Editor**: Drag-and-drop interface for creating Axon programs
- **Custom Node Types**: Support for all Axon node types:
  - START/END: Program entry and exit points
  - CONSTANT: Static values
  - FUNCTION: Function calls (e.g., `fmt.Println`, `os.ReadFile`)
  - OPERATOR: Binary operations (+, -, *, /, etc.)
  - IGNORE: Explicit value discard
- **Dual Edge System**:
  - **Data Edges** (Blue): Represent data flow between nodes
  - **Execution Edges** (Orange): Control program execution order
- **File Operations**: Open, edit, and save `.ax` files
- **Modern UI**: Clean, dark-themed interface optimized for node editing
- **Real-time Updates**: Live connection and node manipulation

## Tech Stack

- **Frontend**: Svelte 5 + Svelte-Flow + Vite
- **Backend**: Go (via Wails v2.11.0)
- **Node Graph**: @xyflow/svelte
- **Desktop Framework**: Wails

## Prerequisites

- Go 1.24 or later
- Node.js 16 or later
- npm or yarn

For building desktop applications, you also need:
- Linux: `libgtk-3-dev`, `libwebkit2gtk-4.0-dev`
- macOS: Xcode Command Line Tools
- Windows: Build tools for Windows

## Installation

1. **Clone the repository**:
```bash
git clone https://github.com/Advik-B/Axon-Editor.git
cd Axon-Editor
```

2. **Install Wails CLI** (if not already installed):
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

3. **Install frontend dependencies**:
```bash
cd frontend
npm install
cd ..
```

## Development

### Running in Development Mode

To run the editor in development mode with hot reload:

```bash
wails dev
```

This will:
- Start a Vite development server with hot module replacement
- Launch the desktop application window
- Enable live reload for both frontend and backend changes

### Frontend Development

You can also run the frontend separately for faster iteration:

```bash
cd frontend
npm run dev
```

Note: Running only the frontend won't have access to file system operations (Open/Save).

## Building

To build a production-ready application:

```bash
wails build
```

The compiled application will be in the `build/bin` directory.

### Platform-Specific Builds

For cross-platform builds:
```bash
# macOS
wails build -platform darwin/universal

# Windows
wails build -platform windows/amd64

# Linux
wails build -platform linux/amd64
```

## Usage

### Creating a New Graph

1. Click **New** to create a blank graph with START and END nodes
2. Use **Add Node** dropdown to add:
   - Constants (values)
   - Functions (like `fmt.Println`)
   - Operators (+, -, *, /, etc.)

### Connecting Nodes

- **Execution Flow** (Orange): Drag from orange exec handles to control execution order
- **Data Flow** (Blue): Drag from blue data handles to pass data between nodes

### Editing Nodes

Click on a node to select it. Node properties like label, type, and configuration can be seen in the node itself.

### Saving and Opening

- **Open**: Load existing `.ax` files
- **Save**: Export your graph as a `.ax` file compatible with Axon CLI

### Example Workflow

1. Create a CONSTANT node with value `5`
2. Create another CONSTANT with value `3`
3. Add an OPERATOR node with operation `+`
4. Connect data edges from constants to operator inputs
5. Add a FUNCTION node with `fmt.Println`
6. Connect operator output to function input
7. Connect execution edges: START → OPERATOR → FUNCTION → END
8. Save as `add.ax`
9. Use Axon CLI to transpile: `axon build add.ax`

## Project Structure

```
.
├── main.go                   # Go backend entry point
├── app.go                    # Application logic and file operations
├── frontend/                 # Svelte frontend
│   ├── src/
│   │   ├── App.svelte                        # Main app component
│   │   ├── components/
│   │   │   ├── NodeGraphEditor.svelte        # Main editor component
│   │   │   └── nodes/                        # Custom Axon node components
│   │   │       ├── StartNode.svelte
│   │   │       ├── EndNode.svelte
│   │   │       ├── ConstantNode.svelte
│   │   │       └── FunctionNode.svelte
│   │   ├── main.js                           # Frontend entry
│   │   └── style.css                         # Global styles
│   ├── package.json                          # Frontend dependencies
│   └── vite.config.js                        # Vite configuration
├── build/                    # Build assets and binaries
├── examples/                 # Example .ax files
└── wails.json               # Wails configuration
```

## Axon File Format

Axon files (`.ax`) are JSON-based and human-readable. Example:

```json
{
  "id": "basic-addition",
  "name": "Add Numbers",
  "imports": ["fmt"],
  "nodes": [
    { "id": "start", "type": "START", "label": "Start" },
    { "id": "const1", "type": "CONSTANT", "label": "x", 
      "outputs": [{"name": "out", "type_name": "int"}],
      "config": {"value": "5"} },
    ...
  ],
  "data_edges": [
    { "from_node_id": "const1", "from_port": "out", 
      "to_node_id": "sum", "to_port": "a" }
  ],
  "exec_edges": [
    { "from_node_id": "start", "to_node_id": "sum" }
  ]
}
```

## Integration with Axon

This editor creates `.ax` files that work seamlessly with the Axon CLI:

```bash
# Build your visual program to Go code
axon build myprogram.ax

# Preview the graph structure
axon preview myprogram.ax

# Run the generated Go code
go run out/main.go
```

Learn more about Axon at: https://github.com/Advik-B/Axon

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is part of the Axon project and shares its MIT License.

## Credits

- **Axon**: Visual programming language by [Advik-B](https://github.com/Advik-B)
- **Wails**: Go + Web framework - https://wails.io
- **Svelte**: Reactive web framework - https://svelte.dev
- **Svelte-Flow**: Node graph library - https://svelteflow.dev

## Related Projects

- [Axon](https://github.com/Advik-B/Axon) - The core Axon transpiler and CLI
- [Wails](https://wails.io) - Build desktop apps using Go & Web Technologies
- [Svelte](https://svelte.dev) - Cybernetically enhanced web apps
