<script>
  // Properties panel for editing selected node properties
  let selectedNode = $state(null);
  let isVisible = $state(false);
  
  export function setNode(node) {
    selectedNode = node ? { ...node } : null;
    isVisible = !!node;
  }
  
  export function hide() {
    isVisible = false;
    selectedNode = null;
  }
  
  export function show() {
    if (selectedNode) {
      isVisible = true;
    }
  }
  
  function updateProperty(key, value) {
    if (selectedNode && selectedNode.data) {
      if (!selectedNode.data.config) {
        selectedNode.data.config = {};
      }
      selectedNode.data.config[key] = value;
      // Dispatch update event
      dispatchUpdate();
    }
  }
  
  function updateLabel(value) {
    if (selectedNode && selectedNode.data) {
      selectedNode.data.label = value;
      dispatchUpdate();
    }
  }
  
  function updateImplReference(value) {
    if (selectedNode && selectedNode.data) {
      selectedNode.data.impl_reference = value;
      dispatchUpdate();
    }
  }
  
  function addInput() {
    if (selectedNode && selectedNode.data) {
      if (!selectedNode.data.inputs) {
        selectedNode.data.inputs = [];
      }
      selectedNode.data.inputs = [...selectedNode.data.inputs, { name: 'new_input', type_name: 'any' }];
      dispatchUpdate();
    }
  }
  
  function addOutput() {
    if (selectedNode && selectedNode.data) {
      if (!selectedNode.data.outputs) {
        selectedNode.data.outputs = [];
      }
      selectedNode.data.outputs = [...selectedNode.data.outputs, { name: 'new_output', type_name: 'any' }];
      dispatchUpdate();
    }
  }
  
  function removePort(portType, index) {
    if (selectedNode && selectedNode.data) {
      if (portType === 'input' && selectedNode.data.inputs) {
        selectedNode.data.inputs = selectedNode.data.inputs.filter((_, i) => i !== index);
      } else if (portType === 'output' && selectedNode.data.outputs) {
        selectedNode.data.outputs = selectedNode.data.outputs.filter((_, i) => i !== index);
      }
      dispatchUpdate();
    }
  }
  
  function updatePort(portType, index, field, value) {
    if (selectedNode && selectedNode.data) {
      if (portType === 'input' && selectedNode.data.inputs) {
        selectedNode.data.inputs[index][field] = value;
        dispatchUpdate();
      } else if (portType === 'output' && selectedNode.data.outputs) {
        selectedNode.data.outputs[index][field] = value;
        dispatchUpdate();
      }
    }
  }
  
  function dispatchUpdate() {
    // This would normally dispatch a custom event to update the parent
    // For now, we'll just trigger a re-render
    selectedNode = { ...selectedNode };
  }
  
  const nodeTypeDescriptions = {
    'START': 'Entry point of the graph execution',
    'END': 'Exit point of the graph execution',
    'CONSTANT': 'A constant value that can be connected to other nodes',
    'FUNCTION': 'A function call from Go standard library or custom code',
    'OPERATOR': 'A binary operator (+, -, *, /, etc.)',
    'IGNORE': 'Explicitly discard a value',
    'RETURN': 'Return a value from the graph'
  };
</script>

{#if isVisible && selectedNode}
  <div class="properties-panel">
    <div class="properties-header">
      <h3>Node Properties</h3>
      <button onclick={() => hide()} class="close-btn">✕</button>
    </div>
    
    <div class="properties-content">
      <!-- Node Type Info -->
      <div class="property-section">
        <h4>Type: {selectedNode.data.type}</h4>
        <p class="type-description">{nodeTypeDescriptions[selectedNode.data.type] || 'Custom node type'}</p>
      </div>
      
      <!-- Node ID -->
      <div class="property-group">
        <label>ID</label>
        <input type="text" value={selectedNode.id} readonly disabled />
      </div>
      
      <!-- Label -->
      {#if selectedNode.data.type !== 'START' && selectedNode.data.type !== 'END'}
        <div class="property-group">
          <label>Label</label>
          <input 
            type="text" 
            value={selectedNode.data.label} 
            oninput={(e) => updateLabel(e.target.value)}
          />
        </div>
      {/if}
      
      <!-- Implementation Reference (for FUNCTION nodes) -->
      {#if selectedNode.data.type === 'FUNCTION'}
        <div class="property-group">
          <label>Function Reference</label>
          <input 
            type="text" 
            value={selectedNode.data.impl_reference || ''} 
            oninput={(e) => updateImplReference(e.target.value)}
            placeholder="e.g., fmt.Println"
          />
          <small>Go function to call (e.g., fmt.Println, strings.ToUpper)</small>
        </div>
      {/if}
      
      <!-- Configuration -->
      {#if selectedNode.data.config && Object.keys(selectedNode.data.config).length > 0}
        <div class="property-section">
          <h4>Configuration</h4>
          {#each Object.entries(selectedNode.data.config) as [key, value]}
            <div class="property-group">
              <label>{key}</label>
              <input 
                type="text" 
                value={value} 
                oninput={(e) => updateProperty(key, e.target.value)}
              />
            </div>
          {/each}
        </div>
      {/if}
      
      <!-- Inputs -->
      {#if selectedNode.data.inputs || selectedNode.data.type === 'FUNCTION' || selectedNode.data.type === 'OPERATOR'}
        <div class="property-section">
          <div class="section-header">
            <h4>Input Ports</h4>
            <button onclick={addInput} class="add-btn">+ Add</button>
          </div>
          {#if selectedNode.data.inputs && selectedNode.data.inputs.length > 0}
            {#each selectedNode.data.inputs as input, i}
              <div class="port-item">
                <input 
                  type="text" 
                  value={input.name} 
                  oninput={(e) => updatePort('input', i, 'name', e.target.value)}
                  placeholder="Port name"
                  class="port-name"
                />
                <input 
                  type="text" 
                  value={input.type_name} 
                  oninput={(e) => updatePort('input', i, 'type_name', e.target.value)}
                  placeholder="Type"
                  class="port-type"
                />
                <button onclick={() => removePort('input', i)} class="remove-btn">🗑️</button>
              </div>
            {/each}
          {:else}
            <p class="empty-list">No input ports</p>
          {/if}
        </div>
      {/if}
      
      <!-- Outputs -->
      {#if selectedNode.data.outputs || selectedNode.data.type === 'CONSTANT' || selectedNode.data.type === 'OPERATOR'}
        <div class="property-section">
          <div class="section-header">
            <h4>Output Ports</h4>
            <button onclick={addOutput} class="add-btn">+ Add</button>
          </div>
          {#if selectedNode.data.outputs && selectedNode.data.outputs.length > 0}
            {#each selectedNode.data.outputs as output, i}
              <div class="port-item">
                <input 
                  type="text" 
                  value={output.name} 
                  oninput={(e) => updatePort('output', i, 'name', e.target.value)}
                  placeholder="Port name"
                  class="port-name"
                />
                <input 
                  type="text" 
                  value={output.type_name} 
                  oninput={(e) => updatePort('output', i, 'type_name', e.target.value)}
                  placeholder="Type"
                  class="port-type"
                />
                <button onclick={() => removePort('output', i)} class="remove-btn">🗑️</button>
              </div>
            {/each}
          {:else}
            <p class="empty-list">No output ports</p>
          {/if}
        </div>
      {/if}
    </div>
  </div>
{/if}

<style>
  .properties-panel {
    position: fixed;
    right: 20px;
    top: 80px;
    width: 320px;
    max-height: calc(100vh - 180px);
    background: #1e1e1e;
    border: 1px solid #3a3a3a;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.5);
    display: flex;
    flex-direction: column;
    z-index: 1500;
  }
  
  .properties-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    background: #252526;
    border-bottom: 1px solid #3a3a3a;
    border-radius: 8px 8px 0 0;
  }
  
  .properties-header h3 {
    margin: 0;
    font-size: 14px;
    font-weight: 600;
    color: #cccccc;
  }
  
  .close-btn {
    background: transparent;
    border: none;
    color: #999;
    cursor: pointer;
    font-size: 18px;
    padding: 0 4px;
    transition: color 0.2s;
  }
  
  .close-btn:hover {
    color: #fff;
  }
  
  .properties-content {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
  }
  
  .property-section {
    margin-bottom: 20px;
  }
  
  .property-section h4 {
    margin: 0 0 8px 0;
    font-size: 13px;
    font-weight: 600;
    color: #4a9eff;
  }
  
  .type-description {
    margin: 0 0 12px 0;
    font-size: 11px;
    color: #888;
    font-style: italic;
  }
  
  .property-group {
    margin-bottom: 12px;
  }
  
  .property-group label {
    display: block;
    margin-bottom: 4px;
    font-size: 11px;
    color: #999;
    font-weight: 500;
    text-transform: uppercase;
  }
  
  .property-group input {
    width: 100%;
    padding: 6px 8px;
    background: #252526;
    border: 1px solid #3a3a3a;
    border-radius: 3px;
    color: #cccccc;
    font-size: 12px;
    font-family: 'Consolas', 'Monaco', monospace;
  }
  
  .property-group input:focus {
    outline: none;
    border-color: #4a9eff;
  }
  
  .property-group input:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .property-group small {
    display: block;
    margin-top: 4px;
    font-size: 10px;
    color: #666;
  }
  
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }
  
  .add-btn {
    background: #2d7f5f;
    border: 1px solid #3d9f7f;
    color: #fff;
    padding: 4px 8px;
    border-radius: 3px;
    font-size: 11px;
    cursor: pointer;
    transition: background 0.2s;
  }
  
  .add-btn:hover {
    background: #3d8f6f;
  }
  
  .port-item {
    display: flex;
    gap: 6px;
    margin-bottom: 6px;
    align-items: center;
  }
  
  .port-name {
    flex: 2;
  }
  
  .port-type {
    flex: 1;
  }
  
  .remove-btn {
    background: transparent;
    border: none;
    color: #c94f4f;
    cursor: pointer;
    font-size: 14px;
    padding: 4px;
    transition: color 0.2s;
  }
  
  .remove-btn:hover {
    color: #ff6b6b;
  }
  
  .empty-list {
    margin: 8px 0;
    font-size: 11px;
    color: #666;
    font-style: italic;
  }
  
  /* Scrollbar styling */
  .properties-content::-webkit-scrollbar {
    width: 8px;
  }
  
  .properties-content::-webkit-scrollbar-track {
    background: #1e1e1e;
  }
  
  .properties-content::-webkit-scrollbar-thumb {
    background: #424242;
    border-radius: 4px;
  }
  
  .properties-content::-webkit-scrollbar-thumb:hover {
    background: #4e4e4e;
  }
</style>
