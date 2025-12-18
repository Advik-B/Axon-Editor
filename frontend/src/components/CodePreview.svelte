<script>
  import { onMount, onDestroy } from 'svelte';
  import * as monaco from 'monaco-editor';
  
  // Monaco Editor instance
  let editor;
  let editorContainer;
  let code = $state('');
  let isVisible = $state(false);
  let isGenerating = $state(false);
  
  export function setCode(newCode) {
    code = newCode;
    if (editor && newCode) {
      editor.setValue(newCode);
    }
  }
  
  export function show() {
    isVisible = true;
    // Initialize editor after DOM is ready
    setTimeout(() => initializeEditor(), 50);
  }
  
  export function hide() {
    isVisible = false;
    if (editor) {
      editor.dispose();
      editor = null;
    }
  }
  
  export function toggle() {
    if (isVisible) {
      hide();
    } else {
      show();
    }
  }
  
  export function setGenerating(value) {
    isGenerating = value;
  }
  
  function initializeEditor() {
    if (!editorContainer || editor) return;
    
    // Configure Monaco Editor
    editor = monaco.editor.create(editorContainer, {
      value: code || '',
      language: 'go',
      theme: 'vs-dark',
      readOnly: true,
      automaticLayout: true,
      minimap: {
        enabled: true
      },
      scrollBeyondLastLine: false,
      fontSize: 13,
      fontFamily: "'Consolas', 'Monaco', 'Courier New', monospace",
      lineNumbers: 'on',
      renderWhitespace: 'selection',
      folding: true,
      bracketPairColorization: {
        enabled: true
      },
      stickyScroll: {
        enabled: true
      }
    });
  }
  
  function copyToClipboard() {
    if (editor) {
      const value = editor.getValue();
      navigator.clipboard.writeText(value).then(() => {
        alert('Code copied to clipboard!');
      }).catch(err => {
        console.error('Failed to copy:', err);
      });
    }
  }
  
  function downloadCode() {
    if (editor) {
      const value = editor.getValue();
      const blob = new Blob([value], { type: 'text/plain' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'generated.go';
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    }
  }
  
  onDestroy(() => {
    if (editor) {
      editor.dispose();
    }
  });
</script>

{#if isVisible}
  <div class="code-preview-overlay" onclick={() => hide()} role="dialog" aria-modal="true" onkeydown={(e) => e.key === 'Escape' && hide()}>
    <div class="code-preview-panel" onclick={(e) => e.stopPropagation()} role="document">
      <div class="code-preview-header">
        <div class="code-preview-title">
          <span class="preview-icon">📄</span>
          <span>Generated Go Code</span>
          <span class="monaco-badge">Monaco Editor</span>
        </div>
        <div class="code-preview-actions">
          <button onclick={() => downloadCode()} class="preview-btn" title="Download code">
            💾 Download
          </button>
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
          <div bind:this={editorContainer} class="monaco-container"></div>
        {:else}
          <div class="code-empty">
            <p>No code generated yet</p>
            <p class="hint">Build or preview your graph to see the generated Go code</p>
          </div>
        {/if}
      </div>
      
      <div class="code-preview-footer">
        <span class="code-stats">{code.split('\n').length} lines</span>
        <span class="editor-info">Monaco Editor • Go Language Support • Read-only</span>
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
    max-width: 1200px;
    max-height: 85vh;
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
  
  .monaco-badge {
    font-size: 10px;
    background: #007acc;
    color: white;
    padding: 2px 6px;
    border-radius: 3px;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.5px;
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
    overflow: hidden;
    background: #1e1e1e;
    position: relative;
  }
  
  .monaco-container {
    width: 100%;
    height: 100%;
    min-height: 400px;
  }
  
  .code-loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 80px 40px;
    color: #999;
  }
  
  .spinner {
    width: 50px;
    height: 50px;
    border: 4px solid #3a3a3a;
    border-top-color: #007acc;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 20px;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
  
  .code-empty {
    text-align: center;
    padding: 80px 40px;
    color: #666;
  }
  
  .code-empty p {
    margin: 0 0 8px 0;
    font-size: 16px;
  }
  
  .hint {
    font-size: 13px;
    color: #555;
  }
  
  .code-preview-footer {
    padding: 8px 16px;
    background: #252526;
    border-top: 1px solid #3a3a3a;
    border-radius: 0 0 8px 8px;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .code-stats {
    font-size: 11px;
    color: #888;
  }
  
  .editor-info {
    font-size: 10px;
    color: #666;
  }
</style>
