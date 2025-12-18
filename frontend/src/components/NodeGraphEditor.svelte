<script>
  import { SvelteFlow, Controls, Background, MiniMap, MarkerType } from '@xyflow/svelte';
  import '@xyflow/svelte/dist/style.css';
  
  import StartNode from './nodes/StartNode.svelte';
  import EndNode from './nodes/EndNode.svelte';
  import ConstantNode from './nodes/ConstantNode.svelte';
  import FunctionNode from './nodes/FunctionNode.svelte';
  
  // Dynamic import of Wails functions
  let wailsAvailable = false;
  let OpenFile, SaveFile;
  
  // Try to load Wails runtime
  (async () => {
    try {
      const { OpenFile: of, SaveFile: sf } = await import('../../wailsjs/go/main/App.js');
      OpenFile = of;
      SaveFile = sf;
      wailsAvailable = true;
    } catch (e) {
      console.warn('Wails runtime not available');
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
    const isExecEdge = connection.sourceHandle?.includes('exec') || connection.targetHandle?.includes('exec');
    
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
  }

  // Handle node drag
  function onNodeDragStop(event) {
    // Position updates are handled automatically by Svelte Flow
  }

  // Open file
  async function handleOpen() {
    if (!wailsAvailable || !OpenFile) {
      alert('File operations are not available. Please run in Wails desktop mode.');
      return;
    }
    try {
      const jsonData = await OpenFile();
      const axonGraph = JSON.parse(jsonData);
      axonToFlow(axonGraph);
    } catch (err) {
      console.error('Failed to open file:', err);
      alert('Failed to open file: ' + err.message);
    }
  }

  // Save file
  async function handleSave() {
    if (!wailsAvailable || !SaveFile) {
      alert('File operations are not available. Please run in Wails desktop mode.');
      return;
    }
    try {
      const axonGraph = flowToAxon(nodes, edges);
      const jsonData = JSON.stringify(axonGraph, null, 2);
      await SaveFile(jsonData);
    } catch (err) {
      console.error('Failed to save file:', err);
      alert('Failed to save file: ' + err.message);
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

<div class="graph-container">
  <div class="toolbar">
    <div class="toolbar-left">
      <h1>Axon Editor</h1>
      <span class="file-name">{fileName}</span>
    </div>
    <div class="toolbar-center">
      <button on:click={handleNew}>New</button>
      <button on:click={handleOpen}>Open</button>
      <button on:click={handleSave}>Save</button>
    </div>
    <div class="toolbar-right">
      <div class="dropdown">
        <button class="dropdown-btn">Add Node ▼</button>
        <div class="dropdown-content">
          <button on:click={() => handleAddNode('CONSTANT')}>Constant</button>
          <button on:click={() => handleAddNode('FUNCTION')}>Function</button>
          <button on:click={() => handleAddNode('OPERATOR')}>Operator</button>
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
      on:connect={onConnect}
      on:nodedragstop={onNodeDragStop}
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

<style>
  .graph-container {
    width: 100%;
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: #0d0d0d;
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
