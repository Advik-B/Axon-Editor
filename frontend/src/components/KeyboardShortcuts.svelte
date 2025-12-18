<script>
  // Keyboard shortcuts help panel
  let isVisible = $state(false);
  
  export function show() {
    isVisible = true;
  }
  
  export function hide() {
    isVisible = false;
  }
  
  export function toggle() {
    isVisible = !isVisible;
  }
  
  const shortcuts = [
    { category: 'File Operations', items: [
      { keys: ['Ctrl', 'N'], description: 'New graph' },
      { keys: ['Ctrl', 'O'], description: 'Open graph' },
      { keys: ['Ctrl', 'S'], description: 'Save graph' },
    ]},
    { category: 'Build & Run', items: [
      { keys: ['F5'], description: 'Build graph' },
      { keys: ['Ctrl', 'F5'], description: 'Run graph' },
      { keys: ['F7'], description: 'Validate graph' },
      { keys: ['Ctrl', 'P'], description: 'Preview code' },
    ]},
    { category: 'View', items: [
      { keys: ['Ctrl', '`'], description: 'Toggle console' },
      { keys: ['Ctrl', '0'], description: 'Fit view' },
      { keys: ['+'], description: 'Zoom in' },
      { keys: ['-'], description: 'Zoom out' },
    ]},
    { category: 'Help', items: [
      { keys: ['F1'], description: 'Show this help' },
      { keys: ['Esc'], description: 'Close dialogs' },
    ]},
  ];
</script>

{#if isVisible}
  <div class="help-overlay" onclick={() => hide()}>
    <div class="help-panel" onclick={(e) => e.stopPropagation()}>
      <div class="help-header">
        <h2>⌨️ Keyboard Shortcuts</h2>
        <button onclick={() => hide()} class="close-btn">✕</button>
      </div>
      
      <div class="help-content">
        {#each shortcuts as category}
          <div class="shortcut-category">
            <h3>{category.category}</h3>
            <div class="shortcut-list">
              {#each category.items as shortcut}
                <div class="shortcut-item">
                  <div class="shortcut-keys">
                    {#each shortcut.keys as key, i}
                      <kbd>{key}</kbd>
                      {#if i < shortcut.keys.length - 1}
                        <span class="plus">+</span>
                      {/if}
                    {/each}
                  </div>
                  <div class="shortcut-description">{shortcut.description}</div>
                </div>
              {/each}
            </div>
          </div>
        {/each}
      </div>
      
      <div class="help-footer">
        <p>Press <kbd>Esc</kbd> or click outside to close</p>
      </div>
    </div>
  </div>
{/if}

<style>
  .help-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 3000;
    animation: fadeIn 0.2s ease;
  }
  
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  
  .help-panel {
    background: #1e1e1e;
    border: 1px solid #3a3a3a;
    border-radius: 8px;
    width: 90%;
    max-width: 700px;
    max-height: 80vh;
    display: flex;
    flex-direction: column;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.7);
    animation: slideIn 0.3s ease;
  }
  
  @keyframes slideIn {
    from {
      transform: translateY(-20px);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }
  
  .help-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 24px;
    background: #252526;
    border-bottom: 1px solid #3a3a3a;
    border-radius: 8px 8px 0 0;
  }
  
  .help-header h2 {
    margin: 0;
    font-size: 20px;
    font-weight: 600;
    color: #cccccc;
  }
  
  .close-btn {
    background: transparent;
    border: none;
    color: #999;
    cursor: pointer;
    font-size: 24px;
    padding: 0 8px;
    transition: color 0.2s;
  }
  
  .close-btn:hover {
    color: #fff;
  }
  
  .help-content {
    flex: 1;
    overflow-y: auto;
    padding: 24px;
  }
  
  .shortcut-category {
    margin-bottom: 28px;
  }
  
  .shortcut-category:last-child {
    margin-bottom: 0;
  }
  
  .shortcut-category h3 {
    margin: 0 0 12px 0;
    font-size: 14px;
    font-weight: 600;
    color: #4a9eff;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  
  .shortcut-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  
  .shortcut-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    background: #252526;
    border-radius: 4px;
    border: 1px solid #3a3a3a;
  }
  
  .shortcut-keys {
    display: flex;
    gap: 4px;
    align-items: center;
    min-width: 180px;
  }
  
  kbd {
    display: inline-block;
    padding: 4px 8px;
    background: linear-gradient(to bottom, #3a3a3a 0%, #2a2a2a 100%);
    border: 1px solid #555;
    border-radius: 3px;
    font-family: 'Consolas', 'Monaco', monospace;
    font-size: 11px;
    font-weight: 600;
    color: #cccccc;
    box-shadow: 0 2px 0 #1a1a1a;
    min-width: 28px;
    text-align: center;
  }
  
  .plus {
    color: #666;
    font-size: 10px;
    font-weight: bold;
  }
  
  .shortcut-description {
    color: #999;
    font-size: 13px;
    flex: 1;
    text-align: right;
  }
  
  .help-footer {
    padding: 16px 24px;
    background: #252526;
    border-top: 1px solid #3a3a3a;
    border-radius: 0 0 8px 8px;
    text-align: center;
  }
  
  .help-footer p {
    margin: 0;
    font-size: 12px;
    color: #666;
  }
  
  .help-footer kbd {
    margin: 0 4px;
  }
  
  /* Scrollbar styling */
  .help-content::-webkit-scrollbar {
    width: 10px;
  }
  
  .help-content::-webkit-scrollbar-track {
    background: #1e1e1e;
  }
  
  .help-content::-webkit-scrollbar-thumb {
    background: #424242;
    border-radius: 5px;
  }
  
  .help-content::-webkit-scrollbar-thumb:hover {
    background: #4e4e4e;
  }
</style>
