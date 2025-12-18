<script>
  // Console component for displaying output, errors, and messages
  let consoleOutput = $state([]);
  let isExpanded = $state(true);
  let filter = $state('all'); // 'all', 'errors', 'warnings', 'info', 'output'
  
  // Export functions to be called from parent
  export function addMessage(message, type = 'info') {
    consoleOutput = [...consoleOutput, {
      id: Date.now() + Math.random(),
      timestamp: new Date(),
      message,
      type, // 'error', 'warning', 'info', 'output', 'success'
    }];
    
    // Auto-scroll to bottom
    setTimeout(() => {
      const consoleElement = document.getElementById('console-content');
      if (consoleElement) {
        consoleElement.scrollTop = consoleElement.scrollHeight;
      }
    }, 10);
  }
  
  export function clear() {
    consoleOutput = [];
  }
  
  function formatTime(date) {
    return date.toLocaleTimeString('en-US', { hour12: false });
  }
  
  function getFilteredOutput() {
    if (filter === 'all') return consoleOutput;
    return consoleOutput.filter(item => item.type === filter);
  }
  
  const filtered = $derived(getFilteredOutput());
  
  function getTypeIcon(type) {
    switch (type) {
      case 'error': return '❌';
      case 'warning': return '⚠️';
      case 'success': return '✅';
      case 'output': return '▶';
      default: return 'ℹ️';
    }
  }
  
  function getTypeClass(type) {
    switch (type) {
      case 'error': return 'console-error';
      case 'warning': return 'console-warning';
      case 'success': return 'console-success';
      case 'output': return 'console-output';
      default: return 'console-info';
    }
  }
</script>

<div class="console-container" class:expanded={isExpanded}>
  <div class="console-header">
    <div class="console-title">
      <span class="console-icon">🖥️</span>
      <span>Console</span>
      <span class="console-count">{filtered.length} messages</span>
    </div>
    <div class="console-controls">
      <select bind:value={filter} class="filter-select">
        <option value="all">All</option>
        <option value="output">Output</option>
        <option value="error">Errors</option>
        <option value="warning">Warnings</option>
        <option value="info">Info</option>
        <option value="success">Success</option>
      </select>
      <button onclick={() => clear()} class="console-btn" title="Clear console">
        🗑️
      </button>
      <button onclick={() => isExpanded = !isExpanded} class="console-btn" title="Toggle console">
        {isExpanded ? '▼' : '▲'}
      </button>
    </div>
  </div>
  
  {#if isExpanded}
    <div id="console-content" class="console-content">
      {#if filtered.length === 0}
        <div class="console-empty">
          <p>No messages to display</p>
        </div>
      {:else}
        {#each filtered as item (item.id)}
          <div class="console-line {getTypeClass(item.type)}">
            <span class="console-time">[{formatTime(item.timestamp)}]</span>
            <span class="console-icon-inline">{getTypeIcon(item.type)}</span>
            <span class="console-message">{item.message}</span>
          </div>
        {/each}
      {/if}
    </div>
  {/if}
</div>

<style>
  .console-container {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    background: #1e1e1e;
    border-top: 1px solid #3a3a3a;
    z-index: 1000;
    display: flex;
    flex-direction: column;
    max-height: 40vh;
    transition: max-height 0.3s ease;
  }
  
  .console-container:not(.expanded) {
    max-height: 40px;
  }
  
  .console-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    background: #252526;
    border-bottom: 1px solid #3a3a3a;
  }
  
  .console-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    font-weight: 500;
    color: #cccccc;
  }
  
  .console-icon {
    font-size: 16px;
  }
  
  .console-count {
    font-size: 11px;
    color: #888;
    background: #3a3a3a;
    padding: 2px 6px;
    border-radius: 3px;
  }
  
  .console-controls {
    display: flex;
    gap: 8px;
    align-items: center;
  }
  
  .filter-select {
    background: #3a3a3a;
    color: #cccccc;
    border: 1px solid #555;
    padding: 4px 8px;
    border-radius: 3px;
    font-size: 12px;
    cursor: pointer;
  }
  
  .filter-select:hover {
    background: #454545;
  }
  
  .console-btn {
    background: transparent;
    border: none;
    color: #cccccc;
    cursor: pointer;
    padding: 4px 8px;
    border-radius: 3px;
    font-size: 14px;
    transition: background 0.2s;
  }
  
  .console-btn:hover {
    background: #3a3a3a;
  }
  
  .console-content {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    font-size: 12px;
    line-height: 1.5;
  }
  
  .console-empty {
    text-align: center;
    padding: 20px;
    color: #666;
  }
  
  .console-line {
    display: flex;
    gap: 8px;
    padding: 2px 4px;
    margin: 1px 0;
    border-radius: 2px;
  }
  
  .console-line:hover {
    background: #2a2a2a;
  }
  
  .console-time {
    color: #666;
    flex-shrink: 0;
  }
  
  .console-icon-inline {
    flex-shrink: 0;
  }
  
  .console-message {
    flex: 1;
    white-space: pre-wrap;
    word-break: break-word;
  }
  
  .console-error {
    color: #f48771;
  }
  
  .console-warning {
    color: #dcdcaa;
  }
  
  .console-success {
    color: #4ec9b0;
  }
  
  .console-output {
    color: #cccccc;
  }
  
  .console-info {
    color: #9cdcfe;
  }
  
  /* Scrollbar styling */
  .console-content::-webkit-scrollbar {
    width: 10px;
  }
  
  .console-content::-webkit-scrollbar-track {
    background: #1e1e1e;
  }
  
  .console-content::-webkit-scrollbar-thumb {
    background: #424242;
    border-radius: 5px;
  }
  
  .console-content::-webkit-scrollbar-thumb:hover {
    background: #4e4e4e;
  }
</style>
