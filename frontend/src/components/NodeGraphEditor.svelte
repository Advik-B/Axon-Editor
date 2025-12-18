<script>
  import { SvelteFlow, Controls, Background, MiniMap, MarkerType } from '@xyflow/svelte';
  import '@xyflow/svelte/dist/style.css';
  
  import StartNode from './nodes/StartNode.svelte';
  import EndNode from './nodes/EndNode.svelte';
  import ConstantNode from './nodes/ConstantNode.svelte';
  import FunctionNode from './nodes/FunctionNode.svelte';
  import Console from './Console.svelte';
  import CodePreview from './CodePreview.svelte';
  import PropertiesPanel from './PropertiesPanel.svelte';
  import KeyboardShortcuts from './KeyboardShortcuts.svelte';
  import ContextMenu from './ContextMenu.svelte';
  
  // Console and Code Preview refs
  let consoleRef;
  let codePreviewRef;
  let propertiesPanelRef;
  let keyboardShortcutsRef;
  let contextMenuRef;
  
  // Dynamic import of Wails functions
  let wailsAvailable = false;
  let OpenFile, SaveFile, ValidateGraph, GenerateGoCode, BuildGraph, RunGraph;
  
  // Try to load Wails runtime
  (async () => {
    try {
      const wailsModule = await import('../../wailsjs/go/main/App.js');
      OpenFile = wailsModule.OpenFile;
      SaveFile = wailsModule.SaveFile;
      ValidateGraph = wailsModule.ValidateGraph;
      GenerateGoCode = wailsModule.GenerateGoCode;
      BuildGraph = wailsModule.BuildGraph;
      RunGraph = wailsModule.RunGraph;
      wailsAvailable = true;
      consoleRef?.addMessage('Axon Editor initialized with LSP support', 'success');
    } catch (e) {
      console.warn('Wails runtime not available');
      consoleRef?.addMessage('Running in browser mode - file operations disabled', 'warning');
    }
  })();

  const nodeTypes = {
    START: StartNode,
    END: EndNode,
    CONSTANT: ConstantNode,
    FUNCTION: FunctionNode,
    OPERATOR: FunctionNode,
    IGNORE: FunctionNode,
    RETURN: FunctionNode,
  };

  let currentGraph = $state({
    id: '',
    name: 'New Graph',
    imports: ['fmt'],
    nodes: [],
    data_edges: [],
    exec_edges: []
  });

  let nodes = $state([]);
  let edges = $state([]);
  const snapGrid = [15, 15];
  let fileName = $state('Untitled');

  // Convert Axon graph to Svelte Flow format
  function axonToFlow(axonGraph) {
    currentGraph = axonGraph;
    fileName = axonGraph.name || 'Untitled';
    
    // Track node count for grid positioning
    let nodeIndex = 0;
    const gridSpacing = 250;
    const gridColumns = 3;
    
    const flowNodes = axonGraph.nodes.map(node => {
      let position;
      if (node.visual_info) {
        position = { x: node.visual_info.x, y: node.visual_info.y };
      } else {
        // Use grid-based positioning for nodes without visual info
        const col = nodeIndex % gridColumns;
        const row = Math.floor(nodeIndex / gridColumns);
        position = { x: 200 + col * gridSpacing, y: 100 + row * gridSpacing };
        nodeIndex++;
      }
      
      return {
        id: node.id,
        type: node.type,
        data: {
          label: node.label,
          type: node.type,
          inputs: node.inputs || [],
          outputs: node.outputs || [],
          impl_reference: node.impl_reference,
          config: node.config || {}
        },
        position
      };
    });

    const dataEdges = (axonGraph.data_edges || []).map((edge, i) => ({
      id: `data-${edge.from_node_id}-${edge.from_port}-${edge.to_node_id}-${edge.to_port}`,
      source: edge.from_node_id,
      sourceHandle: edge.from_port,
      target: edge.to_node_id,
      targetHandle: edge.to_port,
      type: 'default',
      style: 'stroke: #4a9eff; stroke-width: 2px;',
      markerEnd: {
        type: MarkerType.ArrowClosed,
        color: '#4a9eff',
      },
    }));

    const execEdges = (axonGraph.exec_edges || []).map((edge, i) => ({
      id: `exec-${edge.from_node_id}-${edge.to_node_id}`,
      source: edge.from_node_id,
      sourceHandle: edge.from_node_id === 'start' ? 'exec' : 'exec-out',
      target: edge.to_node_id,
      targetHandle: edge.to_node_id === 'end' ? 'exec' : 'exec-in',
      type: 'smoothstep',
      animated: true,
      style: 'stroke: #ff6b00; stroke-width: 3px;',
      markerEnd: {
        type: MarkerType.ArrowClosed,
        color: '#ff6b00',
      },
    }));

    nodes = flowNodes;
    edges = [...dataEdges, ...execEdges];
  }

  // Convert Svelte Flow format back to Axon graph
  function flowToAxon(flowNodes, flowEdges) {
    const axonNodes = flowNodes.map(node => {
      const axonNode = {
        id: node.id,
        type: node.type,
        label: node.data.label,
        visual_info: {
          x: node.position.x,
          y: node.position.y
        }
      };

      if (node.data.inputs && node.data.inputs.length > 0) {
        axonNode.inputs = node.data.inputs;
      }
      if (node.data.outputs && node.data.outputs.length > 0) {
        axonNode.outputs = node.data.outputs;
      }
      if (node.data.impl_reference) {
        axonNode.impl_reference = node.data.impl_reference;
      }
      if (node.data.config && Object.keys(node.data.config).length > 0) {
        axonNode.config = node.data.config;
      }

      return axonNode;
    });

    const dataEdges = flowEdges
      .filter(edge => edge.id.startsWith('data-'))
      .map(edge => ({
        from_node_id: edge.source,
        from_port: edge.sourceHandle,
        to_node_id: edge.target,
        to_port: edge.targetHandle
      }));

    const execEdges = flowEdges
      .filter(edge => edge.id.startsWith('exec-'))
      .map(edge => ({
        from_node_id: edge.source,
        to_node_id: edge.target
      }));

    return {
      ...currentGraph,
      nodes: axonNodes,
      data_edges: dataEdges,
      exec_edges: execEdges
    };
  }

  // Handle new connections
  function onConnect(connection) {
    // Validation 1: Node cannot connect to itself
    if (connection.source === connection.target) {
      consoleRef?.addMessage('Error: A node cannot connect to itself', 'error');
      return;
    }
    
    const isExecEdge = connection.sourceHandle?.includes('exec') || connection.targetHandle?.includes('exec');
    
    // Validation 2: For exec connections, right exec must connect to left exec of ANOTHER node
    if (isExecEdge) {
      // Source handle should be 'exec' or 'exec-out' (right side)
      // Target handle should be 'exec' or 'exec-in' (left side)
      const isSourceRight = connection.sourceHandle === 'exec' || connection.sourceHandle === 'exec-out';
      const isTargetLeft = connection.targetHandle === 'exec' || connection.targetHandle === 'exec-in';
      
      if (!isSourceRight || !isTargetLeft) {
        consoleRef?.addMessage('Error: Execution flow must go from right exec connector to left exec connector of another node', 'error');
        return;
      }
    }
    
    const newEdge = {
      id: isExecEdge 
        ? `exec-${connection.source}-${connection.target}`
        : `data-${connection.source}-${connection.sourceHandle}-${connection.target}-${connection.targetHandle}`,
      source: connection.source,
      sourceHandle: connection.sourceHandle,
      target: connection.target,
      targetHandle: connection.targetHandle,
      type: isExecEdge ? 'smoothstep' : 'default',
      animated: isExecEdge,
      style: isExecEdge ? 'stroke: #ff6b00; stroke-width: 3px;' : 'stroke: #4a9eff; stroke-width: 2px;',
      markerEnd: {
        type: MarkerType.ArrowClosed,
        color: isExecEdge ? '#ff6b00' : '#4a9eff',
      },
    };
    edges = [...edges, newEdge];
    consoleRef?.addMessage('Connection created successfully', 'info');
  }

  // Handle node drag - update positions in nodes array
  function onNodeDragStop(event) {
    // Update node positions after drag
    if (event.detail && event.detail.node) {
      const draggedNode = event.detail.node;
      nodes = nodes.map(n => {
        if (n.id === draggedNode.id) {
          return {
            ...n,
            position: draggedNode.position
          };
        }
        return n;
      });
    }
  }
  
  // Handle node selection
  function onNodeClick(event) {
    if (event.detail && event.detail.node) {
      const node = event.detail.node;
      propertiesPanelRef?.setNode(node);
      consoleRef?.addMessage(`Selected node: ${node.data.label} (${node.data.type})`, 'info');
    }
  }
  
  // Keyboard shortcuts handler
  function handleKeyDown(event) {
    // Ctrl+S - Save
    if (event.ctrlKey && event.key === 's') {
      event.preventDefault();
      handleSave();
    }
    // Ctrl+O - Open
    else if (event.ctrlKey && event.key === 'o') {
      event.preventDefault();
      handleOpen();
    }
    // Ctrl+N - New
    else if (event.ctrlKey && event.key === 'n') {
      event.preventDefault();
      handleNew();
    }
    // F5 - Build
    else if (event.key === 'F5' && !event.ctrlKey) {
      event.preventDefault();
      handleBuild();
    }
    // Ctrl+F5 - Run
    else if (event.ctrlKey && event.key === 'F5') {
      event.preventDefault();
      handleRun();
    }
    // F7 - Validate
    else if (event.key === 'F7') {
      event.preventDefault();
      validateCurrentGraph();
    }
    // Ctrl+P - Preview code
    else if (event.ctrlKey && event.key === 'p') {
      event.preventDefault();
      handlePreviewCode();
    }
    // F1 - Show help
    else if (event.key === 'F1') {
      event.preventDefault();
      keyboardShortcutsRef?.show();
    }
    // Esc - Close dialogs
    else if (event.key === 'Escape') {
      codePreviewRef?.hide();
      keyboardShortcutsRef?.hide();
      propertiesPanelRef?.hide();
    }
  }

  // Open file
  async function handleOpen() {
    if (!wailsAvailable || !OpenFile) {
      consoleRef?.addMessage('File operations not available in browser mode', 'error');
      return;
    }
    try {
      consoleRef?.addMessage('Opening file...', 'info');
      const jsonData = await OpenFile();
      const axonGraph = JSON.parse(jsonData);
      axonToFlow(axonGraph);
      consoleRef?.addMessage(`Graph "${axonGraph.name}" loaded successfully`, 'success');
      
      // Validate on load
      await validateCurrentGraph();
    } catch (err) {
      console.error('Failed to open file:', err);
      consoleRef?.addMessage('Failed to open file: ' + err.message, 'error');
    }
  }

  // Save file
  async function handleSave() {
    if (!wailsAvailable || !SaveFile) {
      consoleRef?.addMessage('File operations not available in browser mode', 'error');
      return;
    }
    try {
      const axonGraph = flowToAxon(nodes, edges);
      const jsonData = JSON.stringify(axonGraph, null, 2);
      await SaveFile(jsonData);
      consoleRef?.addMessage(`Graph "${axonGraph.name}" saved successfully`, 'success');
    } catch (err) {
      console.error('Failed to save file:', err);
      consoleRef?.addMessage('Failed to save file: ' + err.message, 'error');
    }
  }
  
  // Validate graph
  async function validateCurrentGraph() {
    if (!wailsAvailable || !ValidateGraph) {
      consoleRef?.addMessage('Validation not available', 'warning');
      return;
    }
    try {
      const axonGraph = flowToAxon(nodes, edges);
      const jsonData = JSON.stringify(axonGraph);
      consoleRef?.addMessage('Validating graph...', 'info');
      
      const diagnostics = await ValidateGraph(jsonData);
      
      if (!diagnostics || diagnostics.length === 0) {
        consoleRef?.addMessage('✓ Graph validation passed - no errors found', 'success');
      } else {
        consoleRef?.addMessage(`Found ${diagnostics.length} issue(s):`, 'warning');
        diagnostics.forEach(diag => {
          const severity = diag.severity || 'info';
          consoleRef?.addMessage(`  ${diag.message} [${diag.code || 'validation'}]`, severity);
        });
      }
    } catch (err) {
      console.error('Validation failed:', err);
      consoleRef?.addMessage('Validation failed: ' + err.message, 'error');
    }
  }
  
  // Preview generated code
  async function handlePreviewCode() {
    if (!wailsAvailable || !GenerateGoCode) {
      consoleRef?.addMessage('Code generation not available', 'error');
      return;
    }
    try {
      const axonGraph = flowToAxon(nodes, edges);
      const jsonData = JSON.stringify(axonGraph);
      
      codePreviewRef?.setGenerating(true);
      codePreviewRef?.show();
      consoleRef?.addMessage('Generating Go code...', 'info');
      
      const goCode = await GenerateGoCode(jsonData);
      
      codePreviewRef?.setCode(goCode);
      codePreviewRef?.setGenerating(false);
      consoleRef?.addMessage('Code generated successfully', 'success');
    } catch (err) {
      console.error('Code generation failed:', err);
      codePreviewRef?.setGenerating(false);
      consoleRef?.addMessage('Code generation failed: ' + err.message, 'error');
    }
  }
  
  // Build graph
  async function handleBuild() {
    if (!wailsAvailable || !BuildGraph) {
      consoleRef?.addMessage('Build not available', 'error');
      return;
    }
    try {
      const axonGraph = flowToAxon(nodes, edges);
      const jsonData = JSON.stringify(axonGraph);
      
      consoleRef?.addMessage('Building graph with Axon transpiler...', 'info');
      
      const buildOutput = await BuildGraph(jsonData);
      
      consoleRef?.addMessage('Build output:', 'info');
      buildOutput.split('\n').forEach(line => {
        if (line.trim()) {
          consoleRef?.addMessage(line, 'output');
        }
      });
    } catch (err) {
      console.error('Build failed:', err);
      consoleRef?.addMessage('Build failed: ' + err.message, 'error');
    }
  }
  
  // Run graph
  async function handleRun() {
    if (!wailsAvailable || !RunGraph) {
      consoleRef?.addMessage('Run not available', 'error');
      return;
    }
    try {
      const axonGraph = flowToAxon(nodes, edges);
      const jsonData = JSON.stringify(axonGraph);
      
      consoleRef?.addMessage('Running graph...', 'info');
      
      const runOutput = await RunGraph(jsonData);
      
      consoleRef?.addMessage('Execution output:', 'info');
      runOutput.split('\n').forEach(line => {
        if (line.trim()) {
          consoleRef?.addMessage(line, 'output');
        }
      });
    } catch (err) {
      console.error('Run failed:', err);
      consoleRef?.addMessage('Run failed: ' + err.message, 'error');
    }
  }

  // Create new graph
  function handleNew() {
    const newGraph = {
      id: 'new-graph-' + Date.now(),
      name: 'New Graph',
      imports: ['fmt'],
      nodes: [
        {
          id: 'start',
          type: 'START',
          label: 'Start',
          visual_info: { x: 100, y: 250 }
        },
        {
          id: 'end',
          type: 'END',
          label: 'End',
          visual_info: { x: 600, y: 250 }
        }
      ],
      data_edges: [],
      exec_edges: []
    };
    axonToFlow(newGraph);
  }
  
  // Context menu handlers
  function handlePaneContextMenu(event) {
    const e = event.detail?.event || event;
    if (e.preventDefault) e.preventDefault();
    
    const clientX = e.clientX || 0;
    const clientY = e.clientY || 0;
    
    const menuItems = [
      { icon: '➕', label: 'Add Constant', action: () => handleAddNode('CONSTANT') },
      { icon: '⚙️', label: 'Add Function', action: () => handleAddNode('FUNCTION') },
      { icon: '🔢', label: 'Add Operator', action: () => handleAddNode('OPERATOR') },
      { separator: true },
      { icon: '📄', label: 'New Graph', action: handleNew, shortcut: 'Ctrl+N' },
      { icon: '📁', label: 'Open', action: handleOpen, shortcut: 'Ctrl+O' },
      { icon: '💾', label: 'Save', action: handleSave, shortcut: 'Ctrl+S' },
      { separator: true },
      { icon: '✓', label: 'Validate', action: validateCurrentGraph, shortcut: 'F7' },
      { icon: '🔨', label: 'Build', action: handleBuild, shortcut: 'F5' },
      { icon: '▶️', label: 'Run', action: handleRun, shortcut: 'Ctrl+F5' },
    ];
    contextMenuRef?.show(clientX, clientY, menuItems);
  }
  
  function handleNodeContextMenu(event) {
    const e = event.detail?.event || event;
    if (e.preventDefault) e.preventDefault();
    
    const node = event.detail?.node || event.node;
    if (!node) return;
    
    const clientX = e.clientX || 0;
    const clientY = e.clientY || 0;
    
    // Don't allow deleting START or END nodes
    const canDelete = node.data?.type !== 'START' && node.data?.type !== 'END';
    
    const menuItems = [
      { icon: '⚙️', label: 'Properties', action: (n) => propertiesPanelRef?.setNode(n) },
      { icon: '📋', label: 'Duplicate', action: (n) => duplicateNode(n) },
      { separator: true },
      { icon: '🗑️', label: 'Delete', action: (n) => deleteNode(n), danger: true, disabled: !canDelete },
    ];
    contextMenuRef?.show(clientX, clientY, menuItems, { node });
  }
  
  function duplicateNode(node) {
    if (!node) return;
    const newNodeId = `node-${Date.now()}`;
    const newNode = {
      ...node,
      id: newNodeId,
      position: { x: node.position.x + 50, y: node.position.y + 50 },
      data: { ...node.data }
    };
    nodes = [...nodes, newNode];
    consoleRef?.addMessage(`Node duplicated: ${node.data.label}`, 'info');
  }
  
  function deleteNode(node) {
    if (!node) return;
    if (node.data.type === 'START' || node.data.type === 'END') {
      consoleRef?.addMessage('Cannot delete START or END nodes', 'error');
      return;
    }
    
    // Remove node
    nodes = nodes.filter(n => n.id !== node.id);
    
    // Remove connected edges
    edges = edges.filter(e => e.source !== node.id && e.target !== node.id);
    
    consoleRef?.addMessage(`Node deleted: ${node.data.label}`, 'info');
  }

  // Add node
  let nextNodePosition = { x: 350, y: 200 };
  const nodeSpacing = 250;
  
  function handleAddNode(nodeType) {
    const newNodeId = `node-${Date.now()}`;
    let newNode = {
      id: newNodeId,
      type: nodeType,
      data: {
        label: nodeType,
        type: nodeType,
      },
      position: { ...nextNodePosition }
    };
    
    // Move position for next node in a cascading pattern
    nextNodePosition.x += 50;
    nextNodePosition.y += 50;
    if (nextNodePosition.x > 800) {
      nextNodePosition.x = 350;
      nextNodePosition.y += nodeSpacing;
    }

    switch (nodeType) {
      case 'CONSTANT':
        newNode.data.outputs = [{ name: 'out', type_name: 'int' }];
        newNode.data.config = { value: '0' };
        break;
      case 'FUNCTION':
        newNode.data.inputs = [{ name: 'a', type_name: 'any' }];
        newNode.data.outputs = [];
        newNode.data.impl_reference = 'fmt.Println';
        break;
      case 'OPERATOR':
        newNode.data.inputs = [
          { name: 'a', type_name: 'int' },
          { name: 'b', type_name: 'int' }
        ];
        newNode.data.outputs = [{ name: 'out', type_name: 'int' }];
        newNode.data.config = { op: '+' };
        break;
    }

    nodes = [...nodes, newNode];
  }

  // Initialize with example
  handleNew();
</script>

<div class="graph-container" onkeydown={handleKeyDown} tabindex="0">
  <div class="toolbar">
    <div class="toolbar-left">
      <h1>Axon Editor</h1>
      <span class="file-name">{fileName}</span>
    </div>
    <div class="toolbar-center">
      <button onclick={handleNew}>📄 New</button>
      <button onclick={handleOpen}>📁 Open</button>
      <button onclick={handleSave}>💾 Save</button>
      <span class="toolbar-separator"></span>
      <button onclick={validateCurrentGraph}>✓ Validate</button>
      <button onclick={handlePreviewCode}>👁️ Preview Code</button>
      <button onclick={handleBuild}>🔨 Build</button>
      <button onclick={handleRun}>▶️ Run</button>
    </div>
    <div class="toolbar-right">
      <button onclick={() => keyboardShortcutsRef?.show()} class="help-btn" title="Keyboard Shortcuts (F1)">
        ⌨️
      </button>
      <div class="dropdown">
        <button class="dropdown-btn">Add Node ▼</button>
        <div class="dropdown-content">
          <button onclick={() => handleAddNode('CONSTANT')}>Constant</button>
          <button onclick={() => handleAddNode('FUNCTION')}>Function</button>
          <button onclick={() => handleAddNode('OPERATOR')}>Operator</button>
        </div>
      </div>
    </div>
  </div>
  
  <div class="flow-wrapper">
    <SvelteFlow
      {nodes}
      {edges}
      {nodeTypes}
      {snapGrid}
      onconnect={onConnect}
      onnodedragstop={onNodeDragStop}
      onnodeclick={onNodeClick}
      onnodecontextmenu={handleNodeContextMenu}
      onpanecontextmenu={handlePaneContextMenu}
      fitView
    >
      <Controls />
      <Background 
        gap={20}
      />
      <MiniMap 
        nodeColor={(node) => {
          if (node.type === 'START') return '#2d5f2d';
          if (node.type === 'END') return '#7f2d2d';
          if (node.type === 'CONSTANT') return '#5d2d7f';
          if (node.type === 'OPERATOR') return '#7f5f2d';
          return '#2d5f7f';
        }}
        maskColor="rgba(0, 0, 0, 0.7)"
      />
    </SvelteFlow>
  </div>
  
  <div class="info-panel">
    <div class="legend">
      <h3>Legend</h3>
      <div class="legend-item">
        <div class="legend-icon exec-icon"></div>
        <span>Execution Flow</span>
      </div>
      <div class="legend-item">
        <div class="legend-icon data-icon"></div>
        <span>Data Flow</span>
      </div>
    </div>
  </div>
</div>

<!-- Console Component -->
<Console bind:this={consoleRef} />

<!-- Code Preview Component -->
<CodePreview bind:this={codePreviewRef} />

<!-- Properties Panel Component -->
<PropertiesPanel bind:this={propertiesPanelRef} />

<!-- Keyboard Shortcuts Help -->
<KeyboardShortcuts bind:this={keyboardShortcutsRef} />

<!-- Context Menu -->
<ContextMenu bind:this={contextMenuRef} />

<style>
  .graph-container {
    width: 100%;
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: #0d0d0d;
    outline: none; /* Remove focus outline */
  }
  
  .graph-container:focus {
    outline: none;
  }

  .toolbar {
    background: linear-gradient(to bottom, #1e1e1e 0%, #181818 100%);
    padding: 0.75rem 1rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid #2a2a2a;
    min-height: 60px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  }

  .toolbar-left {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .toolbar h1 {
    margin: 0;
    color: #ffffff;
    font-size: 1.5rem;
    font-weight: 600;
    letter-spacing: 0.5px;
  }

  .file-name {
    color: #999999;
    font-size: 0.9rem;
  }

  .toolbar-center, .toolbar-right {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }
  
  .toolbar-separator {
    width: 1px;
    height: 24px;
    background: #3a3a3a;
    margin: 0 4px;
  }

  button {
    padding: 0.5rem 1.2rem;
    background: linear-gradient(to bottom, #2d5f7f 0%, #1f445f 100%);
    color: white;
    border: 1px solid #3d7f9f;
    border-radius: 3px;
    cursor: pointer;
    font-size: 0.85rem;
    font-weight: 500;
    transition: all 0.2s;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
  }

  button:hover {
    background: linear-gradient(to bottom, #3d6f8f 0%, #2f546f 100%);
    border-color: #4d8faf;
  }

  button:active {
    box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.5);
  }

  .dropdown {
    position: relative;
    display: inline-block;
  }

  .dropdown-btn {
    background: linear-gradient(to bottom, #2d7f5f 0%, #1f5f44 100%);
    border-color: #3d9f7f;
  }

  .dropdown-btn:hover {
    background: linear-gradient(to bottom, #3d8f6f 0%, #2f6f54 100%);
    border-color: #4daf8f;
  }

  .dropdown-content {
    display: none;
    position: absolute;
    right: 0;
    background-color: #1e1e1e;
    min-width: 160px;
    box-shadow: 0px 8px 16px 0px rgba(0,0,0,0.5);
    z-index: 1000;
    border-radius: 3px;
    border: 1px solid #2a2a2a;
    overflow: hidden;
  }

  .dropdown-content button {
    width: 100%;
    text-align: left;
    padding: 0.75rem 1rem;
    background: #1e1e1e;
    color: white;
    border: none;
    border-radius: 0;
    border-bottom: 1px solid #2a2a2a;
    box-shadow: none;
  }

  .dropdown-content button:last-child {
    border-bottom: none;
  }

  .dropdown-content button:hover {
    background: #2a2a2a;
  }

  .dropdown:hover .dropdown-content {
    display: block;
  }
  
  .help-btn {
    padding: 0.5rem;
    width: 38px;
    background: linear-gradient(to bottom, #7f5f2d 0%, #5f4f1d 100%);
    border-color: #9f7f3d;
  }
  
  .help-btn:hover {
    background: linear-gradient(to bottom, #8f6f3d 0%, #6f5f2d 100%);
    border-color: #af8f4d;
  }

  .flow-wrapper {
    flex: 1;
    position: relative;
    background: #0d0d0d;
  }

  .info-panel {
    position: absolute;
    top: 80px;
    right: 20px;
    background: rgba(30, 30, 30, 0.95);
    padding: 1rem;
    border-radius: 4px;
    border: 1px solid #2a2a2a;
    z-index: 100;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
  }

  .legend h3 {
    margin: 0 0 0.75rem 0;
    color: #ffffff;
    font-size: 0.9rem;
    font-weight: 600;
  }

  .legend-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.6rem;
    color: #d0d0d0;
    font-size: 0.85rem;
  }

  .legend-icon {
    width: 20px;
    height: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .exec-icon::before {
    content: '';
    width: 0;
    height: 0;
    border-style: solid;
    border-width: 5px 0 5px 8px;
    border-color: transparent transparent transparent #ffffff;
  }

  .data-icon {
    width: 12px;
    height: 12px;
    background: #1b9e77;
    border: 2px solid #ffffff;
    border-radius: 50%;
  }
</style>
