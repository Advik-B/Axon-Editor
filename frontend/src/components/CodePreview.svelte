<script>
  // Code Preview component for displaying generated Go code with syntax highlighting
  let code = $state('');
  let isVisible = $state(false);
  let isGenerating = $state(false);
  
  export function setCode(newCode) {
    code = newCode;
  }
  
  export function show() {
    isVisible = true;
  }
  
  export function hide() {
    isVisible = false;
  }
  
  export function toggle() {
    isVisible = !isVisible;
  }
  
  export function setGenerating(value) {
    isGenerating = value;
  }
  
  function copyToClipboard() {
    navigator.clipboard.writeText(code).then(() => {
      alert('Code copied to clipboard!');
    }).catch(err => {
      console.error('Failed to copy:', err);
    });
  }
  
  // Simple syntax highlighting for Go code
  function highlightGo(code) {
    if (!code) return '';
    
    // Keywords
    let highlighted = code.replace(
      /\b(package|import|func|var|const|type|struct|interface|return|if|else|for|range|switch|case|default|break|continue|go|defer|chan|select|map)\b/g,
      '<span class="keyword">$1</span>'
    );
    
    // Types
    highlighted = highlighted.replace(
      /\b(int|int8|int16|int32|int64|uint|uint8|uint16|uint32|uint64|float32|float64|string|bool|byte|rune|error|interface{}|any)\b/g,
      '<span class="type">$1</span>'
    );
    
    // Strings
    highlighted = highlighted.replace(
      /"([^"\\]|\\.)*"/g,
      '<span class="string">$&</span>'
    );
    
    // Comments
    highlighted = highlighted.replace(
      /\/\/.*/g,
      '<span class="comment">$&</span>'
    );
    
    // Function calls
    highlighted = highlighted.replace(
      /\b([a-zA-Z_][a-zA-Z0-9_]*)\s*\(/g,
      '<span class="function">$1</span>('
    );
    
    return highlighted;
  }
  
  const highlightedCode = $derived(highlightGo(code));
</script>

{#if isVisible}
  <div class="code-preview-overlay" onclick={() => hide()}>
    <div class="code-preview-panel" onclick={(e) => e.stopPropagation()}>
      <div class="code-preview-header">
        <div class="code-preview-title">
          <span class="preview-icon">📄</span>
          <span>Generated Go Code</span>
        </div>
        <div class="code-preview-actions">
          <button onclick={() => copyToClipboard()} class="preview-btn" title="Copy to clipboard">
            📋 Copy
          </button>
          <button onclick={() => hide()} class="preview-btn close-btn" title="Close">
            ✕
          </button>
        </div>
      </div>
      
      <div class="code-preview-content">
        {#if isGenerating}
          <div class="code-loading">
            <div class="spinner"></div>
            <p>Generating code...</p>
          </div>
        {:else if code}
          <pre class="code-block"><code>{@html highlightedCode}</code></pre>
        {:else}
          <div class="code-empty">
            <p>No code generated yet</p>
            <p class="hint">Build or preview your graph to see the generated Go code</p>
          </div>
        {/if}
      </div>
      
      <div class="code-preview-footer">
        <span class="code-stats">{code.split('\n').length} lines</span>
      </div>
    </div>
  </div>
{/if}

<style>
  .code-preview-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
    animation: fadeIn 0.2s ease;
  }
  
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  
  .code-preview-panel {
    background: #1e1e1e;
    border: 1px solid #3a3a3a;
    border-radius: 8px;
    width: 90%;
    max-width: 900px;
    max-height: 80vh;
    display: flex;
    flex-direction: column;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
    animation: slideUp 0.3s ease;
  }
  
  @keyframes slideUp {
    from {
      transform: translateY(20px);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }
  
  .code-preview-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    background: #252526;
    border-bottom: 1px solid #3a3a3a;
    border-radius: 8px 8px 0 0;
  }
  
  .code-preview-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    font-weight: 600;
    color: #cccccc;
  }
  
  .preview-icon {
    font-size: 18px;
  }
  
  .code-preview-actions {
    display: flex;
    gap: 8px;
  }
  
  .preview-btn {
    background: #3a3a3a;
    color: #cccccc;
    border: 1px solid #555;
    padding: 6px 12px;
    border-radius: 4px;
    font-size: 12px;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .preview-btn:hover {
    background: #454545;
    border-color: #666;
  }
  
  .close-btn {
    background: #c94f4f;
    border-color: #d66;
  }
  
  .close-btn:hover {
    background: #d66;
  }
  
  .code-preview-content {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    background: #1e1e1e;
  }
  
  .code-loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 40px;
    color: #999;
  }
  
  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid #3a3a3a;
    border-top-color: #4a9eff;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 16px;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
  
  .code-empty {
    text-align: center;
    padding: 40px;
    color: #666;
  }
  
  .hint {
    font-size: 12px;
    color: #555;
    margin-top: 8px;
  }
  
  .code-block {
    margin: 0;
    padding: 16px;
    background: #252526;
    border-radius: 4px;
    border: 1px solid #3a3a3a;
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    font-size: 13px;
    line-height: 1.6;
    overflow-x: auto;
  }
  
  .code-block code {
    color: #d4d4d4;
  }
  
  .code-block :global(.keyword) {
    color: #569cd6;
    font-weight: 500;
  }
  
  .code-block :global(.type) {
    color: #4ec9b0;
  }
  
  .code-block :global(.string) {
    color: #ce9178;
  }
  
  .code-block :global(.comment) {
    color: #6a9955;
    font-style: italic;
  }
  
  .code-block :global(.function) {
    color: #dcdcaa;
  }
  
  .code-preview-footer {
    padding: 8px 16px;
    background: #252526;
    border-top: 1px solid #3a3a3a;
    border-radius: 0 0 8px 8px;
  }
  
  .code-stats {
    font-size: 11px;
    color: #888;
  }
  
  /* Scrollbar styling */
  .code-preview-content::-webkit-scrollbar {
    width: 10px;
  }
  
  .code-preview-content::-webkit-scrollbar-track {
    background: #1e1e1e;
  }
  
  .code-preview-content::-webkit-scrollbar-thumb {
    background: #424242;
    border-radius: 5px;
  }
  
  .code-preview-content::-webkit-scrollbar-thumb:hover {
    background: #4e4e4e;
  }
  
  .code-block::-webkit-scrollbar {
    height: 8px;
  }
  
  .code-block::-webkit-scrollbar-track {
    background: #252526;
  }
  
  .code-block::-webkit-scrollbar-thumb {
    background: #424242;
    border-radius: 4px;
  }
</style>
