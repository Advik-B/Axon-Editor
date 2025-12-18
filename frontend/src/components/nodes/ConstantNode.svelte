<script>
  import { Handle, Position } from '@xyflow/svelte';
  
  export let data;
  export let selected = false;
  
  $: outputs = data.outputs || [];
</script>

<div class="blueprint-node constant-node" class:selected>
  <div class="node-title-bar">
    <span class="node-title">{data.label || 'Constant'}</span>
  </div>
  <div class="node-content">
    {#if data.config?.value}
      <div class="constant-value">{data.config.value}</div>
    {/if}
    <div class="node-pins">
      {#each outputs as output, i}
        <div class="pin-row output-pin">
          <span class="pin-label">{output.name}</span>
          <span class="pin-type">{output.type_name}</span>
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
</div>

<style>
  .blueprint-node {
    background: #1a1a1a;
    border-radius: 4px;
    overflow: hidden;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
    min-width: 160px;
  }

  .blueprint-node.selected {
    box-shadow: 0 0 0 2px #ffa500, 0 2px 8px rgba(0, 0, 0, 0.5);
  }

  .constant-node .node-title-bar {
    background: linear-gradient(to bottom, #5d2d7f 0%, #3d1f5f 100%);
    border-bottom: 1px solid #7d3d9f;
  }

  .node-title-bar {
    padding: 6px 12px;
    font-size: 0.85rem;
    font-weight: 600;
    color: #ffffff;
    text-align: center;
    letter-spacing: 0.3px;
  }

  .node-content {
    padding: 8px 12px;
  }

  .constant-value {
    font-size: 1rem;
    color: #e0e0e0;
    font-family: 'Consolas', 'Courier New', monospace;
    text-align: center;
    padding: 4px 8px;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 2px;
    margin-bottom: 8px;
  }

  .node-pins {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .pin-row {
    display: flex;
    align-items: center;
    gap: 6px;
    position: relative;
    padding-right: 16px;
    justify-content: flex-end;
  }

  .pin-label {
    font-size: 0.8rem;
    color: #d0d0d0;
  }

  .pin-type {
    font-size: 0.75rem;
    color: #888888;
    font-family: 'Consolas', monospace;
  }

  :global(.data-handle) {
    width: 12px;
    height: 12px;
    background: #1b9e77;
    border: 2px solid #ffffff;
    border-radius: 50%;
  }
</style>
