<script>
  import { Handle, Position } from '@xyflow/svelte';
  
  export let data;
  export let selected = false;
</script>

<div class="axon-node end-node" class:selected>
  <Handle type="target" position={Position.Left} id="exec" class="exec-handle" />
  <div class="node-header">
    <span class="node-type">END</span>
  </div>
  <div class="node-body">
    <div class="node-label">{data.label || 'End'}</div>
  </div>
</div>

<style>
  .axon-node {
    padding: 10px;
    border-radius: 8px;
    border: 2px solid #4a9eff;
    background: #2a2a2a;
    min-width: 120px;
    color: #ffffff;
  }

  .axon-node.selected {
    border-color: #ffaa00;
    box-shadow: 0 0 10px rgba(255, 170, 0, 0.5);
  }

  .end-node {
    border-color: #ff0000;
    background: linear-gradient(135deg, #3a1a1a 0%, #2a2a2a 100%);
  }

  .node-header {
    font-size: 0.75rem;
    font-weight: bold;
    color: #ff0000;
    margin-bottom: 5px;
    text-align: center;
  }

  .node-body {
    text-align: center;
  }

  .node-label {
    font-size: 0.9rem;
    color: #ffffff;
  }

  :global(.exec-handle) {
    width: 12px;
    height: 12px;
    background: #ff6b00;
    border: 2px solid #ffffff;
  }
</style>
