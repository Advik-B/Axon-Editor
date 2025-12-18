<script>
  import { Handle, Position } from '@xyflow/svelte';
  
  export let data;
  export let selected = false;
  
  $: outputs = data.outputs || [];
</script>

<div class="axon-node constant-node" class:selected>
  <div class="node-header">
    <span class="node-type">CONSTANT</span>
  </div>
  <div class="node-body">
    <div class="node-label">{data.label || 'Constant'}</div>
    {#if data.config?.value}
      <div class="node-value">{data.config.value}</div>
    {/if}
  </div>
  <div class="node-ports">
    {#each outputs as output, i}
      <div class="port output-port">
        <span class="port-name">{output.name}</span>
        <span class="port-type">{output.type_name}</span>
        <Handle 
          type="source" 
          position={Position.Right} 
          id={output.name}
          class="data-handle"
        />
      </div>
    {/each}
  </div>
</div>

<style>
  .axon-node {
    padding: 10px;
    border-radius: 8px;
    border: 2px solid #4a9eff;
    background: #2a2a2a;
    min-width: 150px;
    color: #ffffff;
  }

  .axon-node.selected {
    border-color: #ffaa00;
    box-shadow: 0 0 10px rgba(255, 170, 0, 0.5);
  }

  .constant-node {
    border-color: #9c27b0;
    background: linear-gradient(135deg, #2a1a3a 0%, #2a2a2a 100%);
  }

  .node-header {
    font-size: 0.7rem;
    font-weight: bold;
    color: #9c27b0;
    margin-bottom: 5px;
    text-align: center;
  }

  .node-body {
    text-align: center;
    margin-bottom: 8px;
  }

  .node-label {
    font-size: 0.9rem;
    color: #ffffff;
    font-weight: 600;
  }

  .node-value {
    font-size: 0.8rem;
    color: #aaaaaa;
    margin-top: 4px;
    font-family: monospace;
  }

  .node-ports {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .port {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 4px;
    position: relative;
    padding-right: 12px;
  }

  .port-name {
    font-size: 0.75rem;
    color: #cccccc;
  }

  .port-type {
    font-size: 0.7rem;
    color: #888888;
    font-family: monospace;
  }

  :global(.data-handle) {
    width: 10px;
    height: 10px;
    background: #4a9eff;
    border: 2px solid #ffffff;
  }
</style>
