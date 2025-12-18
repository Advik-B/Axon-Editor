<script>
  import { Handle, Position } from '@xyflow/svelte';
  
  export let data;
  export let selected = false;
  
  $: inputs = data.inputs || [];
  $: outputs = data.outputs || [];
  $: hasExec = data.type === 'FUNCTION' || data.type === 'OPERATOR';
  $: isOperator = data.type === 'OPERATOR';
</script>

<div class="blueprint-node function-node" class:selected class:operator-node={isOperator}>
  {#if hasExec}
    <Handle type="target" position={Position.Left} id="exec-in" class="exec-handle exec-in" style="top: 18px;" />
  {/if}
  
  <div class="node-title-bar">
    <span class="node-title">{data.label || 'Function'}</span>
  </div>
  <div class="node-content">
    {#if data.impl_reference}
      <div class="function-ref">{data.impl_reference}</div>
    {/if}
    {#if data.config?.op}
      <div class="operator-symbol">{data.config.op}</div>
    {/if}
    
    <div class="pins-container">
      {#if inputs.length > 0 || outputs.length > 0}
        <div class="pins-grid">
          <div class="pins-column inputs">
            {#each inputs as input, i}
              <div class="pin-row input-pin">
                <Handle 
                  type="target" 
                  position={Position.Left} 
                  id={input.name}
                  class="data-handle"
                />
                <span class="pin-label">{input.name}</span>
                <span class="pin-type">{input.type_name}</span>
              </div>
            {/each}
          </div>
          
          <div class="pins-column outputs">
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
      {/if}
    </div>
  </div>
  
  {#if hasExec}
    <Handle type="source" position={Position.Right} id="exec-out" class="exec-handle exec-out" style="top: 18px;" />
  {/if}
</div>

<style>
  .blueprint-node {
    background: #1a1a1a;
    border-radius: 4px;
    overflow: hidden;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
    min-width: 200px;
    position: relative;
  }

  .blueprint-node.selected {
    box-shadow: 0 0 0 2px #ffa500, 0 2px 8px rgba(0, 0, 0, 0.5);
  }

  .function-node .node-title-bar {
    background: linear-gradient(to bottom, #2d5f7f 0%, #1f445f 100%);
    border-bottom: 1px solid #3d7f9f;
  }

  .operator-node .node-title-bar {
    background: linear-gradient(to bottom, #7f5f2d 0%, #5f441f 100%);
    border-bottom: 1px solid #9f7f3d;
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

  .function-ref {
    font-size: 0.75rem;
    color: #aaaaaa;
    font-family: 'Consolas', monospace;
    text-align: center;
    margin-bottom: 8px;
  }

  .operator-symbol {
    font-size: 1.5rem;
    color: #ffa500;
    font-weight: bold;
    text-align: center;
    margin-bottom: 8px;
  }

  .pins-container {
    margin-top: 4px;
  }

  .pins-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 8px;
  }

  .pins-column {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .pins-column.inputs {
    align-items: flex-start;
  }

  .pins-column.outputs {
    align-items: flex-end;
  }

  .pin-row {
    display: flex;
    align-items: center;
    gap: 6px;
    position: relative;
  }

  .input-pin {
    padding-left: 16px;
  }

  .output-pin {
    padding-right: 16px;
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

  :global(.exec-handle.exec-in) {
    width: 0;
    height: 0;
    border-style: solid;
    border-width: 6px 10px 6px 0;
    border-color: transparent #ffffff transparent transparent;
    background: transparent !important;
    border-radius: 0 !important;
  }

  :global(.exec-handle.exec-out) {
    width: 0;
    height: 0;
    border-style: solid;
    border-width: 6px 0 6px 10px;
    border-color: transparent transparent transparent #ffffff;
    background: transparent !important;
    border-radius: 0 !important;
  }
</style>
