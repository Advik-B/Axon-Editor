<script>
  import { writable } from 'svelte/store';
  import { SvelteFlow, Controls, Background, MiniMap } from '@xyflow/svelte';
  import '@xyflow/svelte/dist/style.css';

  // Initial nodes
  const initialNodes = [
    {
      id: '1',
      type: 'input',
      data: { label: 'Input Node' },
      position: { x: 250, y: 5 }
    },
    {
      id: '2',
      data: { label: 'Default Node' },
      position: { x: 100, y: 100 }
    },
    {
      id: '3',
      type: 'output',
      data: { label: 'Output Node' },
      position: { x: 400, y: 100 }
    }
  ];

  // Initial edges
  const initialEdges = [
    { id: 'e1-2', source: '1', target: '2', animated: true },
    { id: 'e2-3', source: '2', target: '3' }
  ];

  const nodes = writable(initialNodes);
  const edges = writable(initialEdges);
  const snapGrid = [15, 15];

  // Handle new connections
  function onConnect(connection) {
    edges.update((eds) => {
      const newEdge = {
        ...connection,
        id: `e${connection.source}-${connection.target}`,
        animated: true
      };
      return [...eds, newEdge];
    });
  }

  // Handle node drag stop
  function onNodeDragStop(event) {
    console.log('Node drag stopped', event.detail);
  }
</script>

<div class="graph-container">
  <div class="toolbar">
    <h1>Axon Node Graph Editor</h1>
    <div class="toolbar-buttons">
      <button on:click={() => {
        const newNodeId = String(Date.now());
        nodes.update((nds) => [
          ...nds,
          {
            id: newNodeId,
            data: { label: 'New Node' },
            position: { x: Math.random() * 400, y: Math.random() * 400 }
          }
        ]);
      }}>
        Add Node
      </button>
      <button on:click={() => {
        nodes.set(initialNodes);
        edges.set(initialEdges);
      }}>
        Reset
      </button>
    </div>
  </div>
  
  <div class="flow-wrapper">
    <SvelteFlow
      {nodes}
      {edges}
      {snapGrid}
      on:connect={onConnect}
      on:nodedragstop={onNodeDragStop}
      fitView
    >
      <Controls />
      <Background />
      <MiniMap />
    </SvelteFlow>
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
    padding: 1rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid #3a3a3a;
  }

  .toolbar h1 {
    margin: 0;
    color: #ffffff;
    font-size: 1.5rem;
    font-weight: 600;
  }

  .toolbar-buttons {
    display: flex;
    gap: 0.5rem;
  }

  .toolbar-buttons button {
    padding: 0.5rem 1rem;
    background: #4a9eff;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
    transition: background 0.2s;
  }

  .toolbar-buttons button:hover {
    background: #3a8eef;
  }

  .flow-wrapper {
    flex: 1;
    position: relative;
  }
</style>
