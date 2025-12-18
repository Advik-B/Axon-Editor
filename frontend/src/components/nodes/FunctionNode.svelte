<script>
  import { Handle, Position } from '@xyflow/svelte';
  
  export let data;
  export let selected = false;
  
  $: inputs = data.inputs || [];
  $: outputs = data.outputs || [];
  $: hasExec = data.type === 'FUNCTION' || data.type === 'OPERATOR';
</script>

<div class="axon-node function-node" class:selected>
  {#if hasExec}
    <Handle type="target" position={Position.Left} id="exec-in" class="exec-handle exec-in" />
  {/if}
  
  <div class="node-header">
    <span class="node-type">{data.type || 'FUNCTION'}</span>
  </div>
  <div class="node-body">
    <div class="node-label">{data.label || 'Function'}</div>
    {#if data.impl_reference}
      <div class="node-impl">{data.impl_reference}</div>
    {/if}
    {#if data.config?.op}
      <div class="node-operator">{data.config.op}</div>
    {/if}
  </div>
  
  <div class="node-ports-container">
    <div class="node-ports inputs">
      {#each inputs as input, i}
        <div class="port input-port">
          <Handle 
            type="target" 
            position={Position.Left} 
            id={input.name}
            class="data-handle"
          />
          <span class="port-name">{input.name}</span>
          <span class="port-type">{input.type_name}</span>
        </div>
      {/each}
    </div>
    
    <div class="node-ports outputs">
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
  
  {#if hasExec}
    <Handle type="source" position={Position.Right} id="exec-out" class="exec-handle exec-out" />
  {/if}
</div>

<style>
  .axon-node {
    padding: 10px;
    border-radius: 8px;
    border: 2px solid #4a9eff;
    background: #2a2a2a;
    min-width: 180px;
    color: #ffffff;
    position: relative;
  }

  .axon-node.selected {
    border-color: #ffaa00;
    box-shadow: 0 0 10px rgba(255, 170, 0, 0.5);
  }

  .function-node {
    border-color: #4a9eff;
    background: linear-gradient(135deg, #1a2a3a 0%, #2a2a2a 100%);
  }

  .node-header {
    font-size: 0.7rem;
    font-weight: bold;
    color: #4a9eff;
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

  .node-impl {
    font-size: 0.75rem;
    color: #aaaaaa;
    margin-top: 4px;
    font-family: monospace;
  }

  .node-operator {
    font-size: 1.2rem;
    color: #ffaa00;
    margin-top: 4px;
    font-weight: bold;
  }

  .node-ports-container {
    display: flex;
    justify-content: space-between;
  }

  .node-ports {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .node-ports.inputs {
    align-items: flex-start;
  }

  .node-ports.outputs {
    align-items: flex-end;
  }

  .port {
    display: flex;
    align-items: center;
    gap: 4px;
    position: relative;
  }

  .input-port {
    padding-left: 12px;
  }

  .output-port {
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

  :global(.exec-handle) {
    width: 12px;
    height: 12px;
    background: #ff6b00;
    border: 2px solid #ffffff;
  }

  :global(.exec-in) {
    top: 10px;
  }

  :global(.exec-out) {
    top: 10px;
  }
</style>
