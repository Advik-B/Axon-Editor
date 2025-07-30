# 🎨 Axon Editor

<p align="center">
  The official visual node graph editor for the <strong><a href="https://github.com/Advik-B/Axon">Axon Language</a></strong>.
</p>
<p align="center">
  <img src="https://img.shields.io/badge/status-in%20development-orange" alt="Status: In Development">
  <img src="https://img.shields.io/badge/Wails-v2-red" alt="Wails">
  <img src="https://img.shields.io/badge/Svelte-v4-orange" alt="Svelte">
  <img src="https://img.shields.io/badge/License-MIT-blue" alt="License: MIT">
</p>

---

**Axon Editor** is a cross-platform desktop application for visually creating, editing, and managing Axon graph files (`.ax`). It provides an intuitive drag-and-drop interface that writes standard, human-readable JSON files compatible with the [Axon CLI](https://github.com/Advik-B/Axon).

This project is the GUI, while the core compiler/transpiler logic lives in the main [Axon repository](https://github.com/Advik-B/Axon).

## ✨ Features (Planned)

-   **Visual Graph Editing**: Create, delete, and arrange nodes on an infinite canvas.
-   **Data & Execution Flow**: Connect nodes with bezier curves for both data and execution flow.
-   **Property Inspector**: Select a node to edit its properties, like labels, constant values, and implementation references.
-   **File Management**: Open and save `.ax` graph files.
-   **Live Transpilation Preview**: See the generated Go code update in real-time as you build your graph, powered by the Axon CLI.
-   **Undo/Redo**: Full history support for all graph operations.

## 🛠️ Tech Stack

-   **Backend & App Wrapper**: [**Wails**](https://wails.io/) (Go)
-   **Frontend UI**: [**Svelte**](https://svelte.dev/) with TypeScript
-   **Canvas/Rendering**: HTML5 Canvas or SVG
-   **Build Tool**: Vite

## 🚀 Getting Started (Development)

To run the editor locally for development, you'll need to set up the Wails environment.

### Prerequisites

-   [Go](https://golang.org/doc/install) (version 1.21 or later)
-   [Node.js](https://nodejs.org/) (LTS version)
-   [Wails CLI](https://wails.io/docs/gettingstarted/installation)

### Installation & Running

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/Advik-B/Axon-Editor.git
    cd Axon-Editor
    ```

2.  **Install frontend dependencies:**
    ```bash
    npm install --prefix frontend
    ```

3.  **Run in development mode:**
    This will start the application with live reloading for both the Go backend and the Svelte frontend.
    ```bash
    wails dev
    ```

4.  **Build the application:**
    To create a production-ready, distributable binary:
    ```bash
    wails build
    ```
    The output will be in the `build/bin` directory.

## 🗺️ Development Roadmap

-   [ ] **Canvas Foundation**:
    -   [ ] Setup pan and zoom controls.
    -   [ ] Render a basic grid background.
-   [ ] **Node Management**:
    -   [ ] Render nodes from a graph state.
    -   [ ] Implement node dragging.
    -   [ ] Add context menu for creating new nodes.
-   [ ] **Edge Management**:
    -   [ ] Draw bezier curves between node ports.
    -   [ ] Implement logic for creating and deleting edges.
-   [ ] **State Management**:
    -   [ ] Use a Svelte store to manage the entire graph state (`nodes`, `edges`, etc.).
    -   [ ] Implement Undo/Redo functionality.
-   [ ] **Backend Integration**:
    -   [ ] Bind Wails backend functions for saving/loading files.
    -   [ ] Implement the live transpilation panel by calling the Axon CLI.
-   [ ] **UI/UX Polish**:
    -   [ ] Create the property inspector panel.
    -   [ ] Refine colors, fonts, and node styles to match the Axon theme.

## 📄 License

This project is licensed under the **MIT License**.