<script>
  // Context menu component for node graph editor
  let isVisible = $state(false);
  let position = $state({ x: 0, y: 0 });
  let menuItems = $state([]);
  let clickedNode = $state(null);
  let clickedEdge = $state(null);
  
  export function show(x, y, items, context = {}) {
    position = { x, y };
    menuItems = items;
    clickedNode = context.node || null;
    clickedEdge = context.edge || null;
    isVisible = true;
  }
  
  export function hide() {
    isVisible = false;
    clickedNode = null;
    clickedEdge = null;
  }
  
  function handleItemClick(item) {
    if (item.action) {
      item.action(clickedNode, clickedEdge);
    }
    hide();
  }
  
  function handleClickOutside(event) {
    if (isVisible) {
      hide();
    }
  }
</script>

<svelte:window onclick={handleClickOutside} />

{#if isVisible}
  <div 
    class="context-menu" 
    style="left: {position.x}px; top: {position.y}px;"
    onclick={(e) => e.stopPropagation()}
    role="menu"
  >
    {#each menuItems as item}
      {#if item.separator}
        <div class="menu-separator"></div>
      {:else}
        <button
          class="menu-item {item.danger ? 'danger' : ''}"
          onclick={() => handleItemClick(item)}
          disabled={item.disabled}
        >
          {#if item.icon}
            <span class="menu-icon">{item.icon}</span>
          {/if}
          <span class="menu-label">{item.label}</span>
          {#if item.shortcut}
            <span class="menu-shortcut">{item.shortcut}</span>
          {/if}
        </button>
      {/if}
    {/each}
  </div>
{/if}

<style>
  .context-menu {
    position: fixed;
    background: #1e1e1e;
    border: 1px solid #3a3a3a;
    border-radius: 4px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.5);
    z-index: 10000;
    min-width: 200px;
    padding: 4px 0;
    animation: menuSlideIn 0.15s ease;
  }
  
  @keyframes menuSlideIn {
    from {
      opacity: 0;
      transform: translateY(-5px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
  
  .menu-item {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    background: transparent;
    border: none;
    color: #cccccc;
    font-size: 13px;
    cursor: pointer;
    text-align: left;
    transition: background 0.1s;
  }
  
  .menu-item:hover:not(:disabled) {
    background: #2a2a2a;
  }
  
  .menu-item:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
  
  .menu-item.danger {
    color: #f48771;
  }
  
  .menu-item.danger:hover:not(:disabled) {
    background: #4a2a2a;
  }
  
  .menu-icon {
    font-size: 14px;
    width: 16px;
    text-align: center;
  }
  
  .menu-label {
    flex: 1;
  }
  
  .menu-shortcut {
    font-size: 11px;
    color: #666;
    font-family: monospace;
  }
  
  .menu-separator {
    height: 1px;
    background: #3a3a3a;
    margin: 4px 0;
  }
</style>
