# Axon Node Graph Editor

A modern node graph editor built with Wails, Svelte 5, and Svelte-Flow, designed for visual programming and workflow design.

## Features

- **Node Graph Editor**: Interactive node-based interface powered by Svelte-Flow
- **Svelte 5**: Latest version of Svelte with improved reactivity and performance
- **Wails Framework**: Native desktop application with Go backend
- **Modern UI**: Clean, dark-themed interface optimized for node editing
- **Real-time Updates**: Live connection and node manipulation
- **Mini Map**: Navigate large graphs easily
- **Controls**: Zoom, pan, and fit-to-view controls

## Tech Stack

- **Frontend**: Svelte 5 + Vite
- **Backend**: Go (via Wails)
- **Node Graph**: Svelte-Flow (@xyflow/svelte)
- **Build Tool**: Vite

## Development

### Prerequisites

- Go 1.18 or later
- Node.js 16 or later
- npm or yarn

### Installation

1. Clone the repository:
```bash
git clone https://github.com/Advik-B/Axon-Editor.git
cd Axon-Editor
```

2. Install dependencies:
```bash
# Install Wails CLI (if not already installed)
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Install frontend dependencies
cd frontend
npm install
cd ..
```

### Running in Development Mode

Run the following command in the project directory:

```bash
wails dev
```

This will:
- Start a Vite development server with hot reload
- Launch the desktop application
- Enable dev server on http://localhost:34115 for browser-based development

You can also run the frontend separately for faster development:

```bash
cd frontend
npm run dev
```

### Building for Production

To build a redistributable, production-ready package:

```bash
wails build
```

The compiled application will be in the `build/bin` directory.

## Usage

### Creating Nodes

- Click the "Add Node" button in the toolbar to create new nodes
- Drag nodes to reposition them
- Click and drag from a node's handle to create connections

### Controls

- **Pan**: Click and drag on the background
- **Zoom**: Use mouse wheel or zoom controls
- **Reset**: Click the "Reset" button to restore the initial graph

### Mini Map

The mini map in the bottom-right corner provides an overview of your entire graph and allows quick navigation.

## Project Structure

```
.
├── main.go              # Go backend entry point
├── app.go               # Application logic
├── frontend/            # Svelte frontend
│   ├── src/
│   │   ├── App.svelte                      # Main application component
│   │   ├── components/
│   │   │   └── NodeGraphEditor.svelte      # Node graph editor component
│   │   ├── main.js                         # Frontend entry point
│   │   └── style.css                       # Global styles
│   ├── package.json     # Frontend dependencies
│   └── vite.config.js   # Vite configuration
├── build/               # Build assets and configuration
└── wails.json           # Wails configuration
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is part of the Axon project by Advik-B.

## Related Projects

- [Axon](https://github.com/Advik-B/Axon) - The main Axon project
- [Wails](https://wails.io) - Go + Web framework
- [Svelte](https://svelte.dev) - The web framework
- [Svelte-Flow](https://svelteflow.dev) - Node graph library for Svelte
