<script>
  import { writable, get } from 'svelte/store';
  import { SvelteFlow, Controls, Background, MiniMap, MarkerType } from '@xyflow/svelte';
  import '@xyflow/svelte/dist/style.css';
  import { OpenFile, SaveFile } from '../../wailsjs/go/main/App.js';
  
  import StartNode from './nodes/StartNode.svelte';
  import EndNode from './nodes/EndNode.svelte';
  import ConstantNode from './nodes/ConstantNode.svelte';
  import FunctionNode from './nodes/FunctionNode.svelte';

  const nodeTypes = {
    START: StartNode,
    END: EndNode,
    CONSTANT: ConstantNode,
    FUNCTION: FunctionNode,
    OPERATOR: FunctionNode,
    IGNORE: FunctionNode,
    RETURN: FunctionNode,
  };

  let currentGraph = {
    id: '',
    name: 'New Graph',
    imports: ['fmt'],
    nodes: [],
    data_edges: [],
    exec_edges: []
  };

  const nodes = writable([]);
  const edges = writable([]);
  const snapGrid = [15, 15];
  let fileName = 'Untitled';

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

    nodes.set(flowNodes);
    edges.set([...dataEdges, ...execEdges]);
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
    edges.update((eds) => {
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
      return [...eds, newEdge];
    });
  }

  // Handle node drag
  function onNodeDragStop(event) {
    // Position updates are handled automatically by Svelte Flow
  }

  // Open file
  async function handleOpen() {
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
    try {
      const currentNodes = get(nodes);
      const currentEdges = get(edges);
      
      const axonGraph = flowToAxon(currentNodes, currentEdges);
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
    nodes.update((nds) => {
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

      return [...nds, newNode];
    });
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
      <Background />
      <MiniMap nodeColor="#4a9eff" />
    </SvelteFlow>
  </div>
  
  <div class="info-panel">
    <div class="legend">
      <h3>Legend</h3>
      <div class="legend-item">
        <div class="legend-line exec-line"></div>
        <span>Execution Flow</span>
      </div>
      <div class="legend-item">
        <div class="legend-line data-line"></div>
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
    background: #1a1a1a;
  }

  .toolbar {
    background: #2a2a2a;
    padding: 0.75rem 1rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid #3a3a3a;
    min-height: 60px;
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
  }

  .file-name {
    color: #aaaaaa;
    font-size: 0.9rem;
  }

  .toolbar-center, .toolbar-right {
    display: flex;
    gap: 0.5rem;
  }

  button {
    padding: 0.5rem 1rem;
    background: #4a9eff;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
    transition: background 0.2s;
  }

  button:hover {
    background: #3a8eef;
  }

  .dropdown {
    position: relative;
    display: inline-block;
  }

  .dropdown-btn {
    background: #2ecc71;
  }

  .dropdown-btn:hover {
    background: #27ae60;
  }

  .dropdown-content {
    display: none;
    position: absolute;
    right: 0;
    background-color: #2a2a2a;
    min-width: 160px;
    box-shadow: 0px 8px 16px 0px rgba(0,0,0,0.2);
    z-index: 1000;
    border-radius: 4px;
    overflow: hidden;
  }

  .dropdown-content button {
    width: 100%;
    text-align: left;
    padding: 0.75rem 1rem;
    background: #2a2a2a;
    color: white;
    border-radius: 0;
  }

  .dropdown-content button:hover {
    background: #3a3a3a;
  }

  .dropdown:hover .dropdown-content {
    display: block;
  }

  .flow-wrapper {
    flex: 1;
    position: relative;
  }

  .info-panel {
    position: absolute;
    top: 80px;
    right: 20px;
    background: rgba(42, 42, 42, 0.95);
    padding: 1rem;
    border-radius: 8px;
    border: 1px solid #3a3a3a;
    z-index: 100;
  }

  .legend h3 {
    margin: 0 0 0.75rem 0;
    color: #ffffff;
    font-size: 0.9rem;
  }

  .legend-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
    color: #cccccc;
    font-size: 0.85rem;
  }

  .legend-line {
    width: 30px;
    height: 3px;
    border-radius: 2px;
  }

  .exec-line {
    background: #ff6b00;
  }

  .data-line {
    background: #4a9eff;
  }
</style>
