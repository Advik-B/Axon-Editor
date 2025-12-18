# Axon Editor

The official visual node graph IDE for [Axon](https://github.com/Advik-B/Axon) - a visual, node-based programming language that transpiles to Go.

## Overview

Axon Editor is a full-featured desktop IDE built with Wails (Go + Web) and Svelte 5 that provides a comprehensive development environment for creating and editing Axon visual programs. It features LSP integration with gopls, integrated console, code preview, real-time validation, and direct integration with the Axon CLI transpiler. Design complex programs using a node-based workflow and save them as `.ax` files that can be transpiled to Go code.

## Features

### Core IDE Features

- **🎯 Visual Node Editor**: Intuitive drag-and-drop interface for creating Axon programs
- **📝 LSP Integration**: Language Server Protocol support with gopls for Go code intelligence
- **🔍 Real-time Validation**: Instant feedback on graph structure, type mismatches, and errors
- **💻 Integrated Console**: Built-in console for build output, errors, and execution results
- **👁️ Code Preview**: Real-time Go code preview with syntax highlighting
- **🔨 Build & Run**: Direct integration with Axon CLI for transpilation and execution
- **⚙️ Properties Panel**: Edit node properties, ports, and configuration on-the-fly
- **⌨️ Keyboard Shortcuts**: Comprehensive keyboard shortcuts for efficient workflow

### Node Types

Support for all Axon node types:
- **START/END**: Program entry and exit points
- **CONSTANT**: Static values with configurable types
- **FUNCTION**: Function calls from Go standard library (e.g., `fmt.Println`, `os.ReadFile`)
- **OPERATOR**: Binary operations (+, -, *, /, type casting)
- **IGNORE**: Explicit value discard for error handling
- **RETURN**: Return values from graphs

### Graph Features

- **Dual Edge System**:
  - **Data Edges** (Blue): Represent data flow between nodes
  - **Execution Edges** (Orange): Control program execution order
- **Type Checking**: Automatic validation of port type compatibility
- **Cycle Detection**: Prevents infinite loops in execution flow
- **Unreachable Node Detection**: Warns about nodes disconnected from START
- **MiniMap**: Overview of large graphs for easy navigation
- **Zoom & Pan**: Smooth navigation with fit-to-view support
- **Grid Snapping**: Align nodes to a grid for clean layouts

### IDE Capabilities

- **File Operations**: Open, edit, and save `.ax` files
- **Validation**: Real-time graph validation with detailed diagnostics
- **Build System**: Transpile to Go and compile with error reporting
- **Execution**: Run generated programs and view output in console
- **Code Generation**: Uses actual Axon CLI transpiler (with fallback)
- **gopls Integration**: Go language server for code intelligence
- **Syntax Highlighting**: Color-coded Go code preview
- **Message Filtering**: Filter console by errors, warnings, info, and output

### User Interface

- **Modern Dark Theme**: Professional IDE-style interface
- **Collapsible Console**: Bottom panel with auto-scroll and filtering
- **Modal Dialogs**: Code preview and help overlays
- **Toolbar**: Quick access to common operations
- **Status Indicators**: Visual feedback for operations
- **Keyboard Navigation**: Full keyboard shortcut support

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Ctrl+N` | New graph |
| `Ctrl+O` | Open file |
| `Ctrl+S` | Save file |
| `F7` | Validate graph |
| `Ctrl+P` | Preview generated code |
| `F5` | Build graph |
| `Ctrl+F5` | Run graph |
| `F1` | Show keyboard shortcuts help |
| `Esc` | Close dialogs |

## Tech Stack

- **Frontend**: Svelte 5 + Svelte-Flow + Vite
- **Backend**: Go (via Wails v2.11.0)
- **Node Graph**: @xyflow/svelte
- **Desktop Framework**: Wails
- **LSP**: gopls integration for Go code intelligence
- **Transpiler**: Axon CLI integration

## Prerequisites

- Go 1.24 or later
- Node.js 16 or later
- npm or yarn
- **Axon CLI** (recommended): `go install github.com/Advik-B/Axon@latest`
- **gopls** (recommended): `go install golang.org/x/tools/gopls@latest`

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
- Initialize LSP server with gopls (if installed)
- Enable live reload for both frontend and backend changes
- Display console output for debugging

### Frontend Development

You can also run the frontend separately for faster iteration:

```bash
cd frontend
npm run dev
```

**Note**: Running only the frontend won't have access to:
- File system operations (Open/Save)
- LSP validation and code generation
- Build and run functionality
- Axon CLI integration

### Working with the LSP

The editor integrates with gopls for Go code intelligence. To ensure full functionality:

1. Install gopls: `go install golang.org/x/tools/gopls@latest`
2. Install Axon CLI: `go install github.com/Advik-B/Axon@latest`
3. The LSP server will automatically start when the application launches
4. Validation uses both Axon-specific rules and gopls for Go code checking

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

1. Click **📄 New** (or press `Ctrl+N`) to create a blank graph with START and END nodes
2. Use **Add Node ▼** dropdown to add:
   - **Constants**: Static values with configurable types
   - **Functions**: Go standard library functions (e.g., `fmt.Println`, `os.ReadFile`)
   - **Operators**: Mathematical operations or type casting

### Connecting Nodes

- **Execution Flow** (Orange): Drag from orange exec handles to control execution order
- **Data Flow** (Blue): Drag from blue data handles to pass data between nodes
- The editor validates type compatibility and warns about mismatches

### Editing Nodes

1. **Click a node** to select it and open the Properties Panel
2. In the Properties Panel, you can:
   - Change the node label
   - Edit function references
   - Modify input/output ports
   - Add or remove ports
   - Configure node-specific settings

### IDE Workflow

1. **Design**: Create your visual program using nodes and edges
2. **Validate** (`F7`): Check for errors, type mismatches, and structural issues
3. **Preview** (`Ctrl+P`): View the generated Go code with syntax highlighting
4. **Build** (`F5`): Transpile to Go using Axon CLI and compile
5. **Run** (`Ctrl+F5`): Execute your program and view output in console

### File Operations

- **Open** (`Ctrl+O`): Load existing `.ax` files
- **Save** (`Ctrl+S`): Export your graph as a `.ax` file compatible with Axon CLI
- The console displays confirmation messages and any errors

### Using the Console

- Located at the bottom of the screen
- Filter messages by type: All, Output, Errors, Warnings, Info, Success
- Click 🗑️ to clear console
- Click ▼/▲ to collapse/expand
- Auto-scrolls to latest messages

### Example Workflow: Building a Simple Addition Program

1. Create a **CONSTANT** node with value `5`, label it `x`
2. Create another **CONSTANT** with value `3`, label it `y`
3. Add an **OPERATOR** node, set operation to `+`, label it `sum`
4. Connect **data edges**: `x.out` → `sum.a`, `y.out` → `sum.b`
5. Add a **FUNCTION** node with `fmt.Println`, label it `print`
6. Connect **data edge**: `sum.out` → `print.a`
7. Connect **execution edges**: START → sum → print → END
8. Press `F7` to validate (should pass)
9. Press `Ctrl+P` to preview the generated Go code
10. Press `F5` to build
11. Press `Ctrl+F5` to run and see output: `8`
12. Save as `add.ax` with `Ctrl+S`

### Advanced Features

- **Keyboard Shortcuts**: Press `F1` to view all available shortcuts
- **Properties Panel**: Click any node to edit its properties
- **Type Validation**: The LSP validates Go types and warns about incompatibilities
- **Error Detection**: Cycle detection, unreachable nodes, and missing connections
- **Code Intelligence**: gopls provides Go-aware validation of generated code

## Project Structure

```
.
├── main.go                   # Go backend entry point
├── app.go                    # Application logic, file operations, LSP integration
├── lsp/
│   └── server.go            # LSP server with gopls integration
├── frontend/                 # Svelte frontend
│   ├── src/
│   │   ├── App.svelte                        # Main app component
│   │   ├── components/
│   │   │   ├── NodeGraphEditor.svelte        # Main editor component
│   │   │   ├── Console.svelte                # Integrated console
│   │   │   ├── CodePreview.svelte            # Code preview modal
│   │   │   ├── PropertiesPanel.svelte        # Node properties editor
│   │   │   ├── KeyboardShortcuts.svelte      # Help panel
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

This editor provides seamless integration with the Axon CLI transpiler:

### Direct Integration

The editor includes built-in integration with Axon CLI:
- Press `F5` to build (transpile + compile)
- Press `Ctrl+F5` to run your program
- View all output in the integrated console
- Automatic error reporting and diagnostics

### Using Axon CLI Separately

You can also use the Axon CLI directly with saved `.ax` files:

```bash
# Build your visual program to Go code
axon build myprogram.ax

# Preview the graph structure
axon preview myprogram.ax

# Run the generated Go code
cd out
go run main.go
```

### LSP and Code Intelligence

The editor uses gopls to provide:
- Real-time type checking for generated Go code
- Validation of Go standard library function calls
- Type compatibility checking between nodes
- Error detection before building

### Fallback Mode

If Axon CLI is not installed, the editor provides a basic fallback transpiler for simple graphs. For full functionality, install Axon CLI:

```bash
go install github.com/Advik-B/Axon@latest
```

Learn more about Axon at: https://github.com/Advik-B/Axon

## Troubleshooting

### LSP Not Working

If validation or code preview isn't working:
1. Install gopls: `go install golang.org/x/tools/gopls@latest`
2. Ensure gopls is in your PATH
3. Restart the editor

### Build Fails

If builds fail with "Axon CLI not found":
1. Install Axon CLI: `go install github.com/Advik-B/Axon@latest`
2. Verify installation: `axon --version`
3. Ensure Go bin directory is in PATH

### Generated Code Has Errors

The editor provides two levels of validation:
1. **Axon validation**: Checks graph structure, cycles, unreachable nodes
2. **Go validation**: Uses gopls to check generated code

Run validation (`F7`) before building to catch issues early.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

Areas for contribution:
- Additional node types
- More keyboard shortcuts
- Undo/redo functionality
- Node templates and snippets
- Enhanced debugging features
- Testing infrastructure

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
