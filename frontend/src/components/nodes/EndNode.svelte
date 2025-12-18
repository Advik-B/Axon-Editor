<script>
  import { Handle, Position } from '@xyflow/svelte';
  
  export let data;
  export let selected = false;
</script>

<div class="blueprint-node end-node" class:selected>
  <Handle type="target" position={Position.Left} id="exec" class="exec-handle exec-in" />
  <div class="node-title-bar">
    <span class="node-title">{data.label || 'End'}</span>
  </div>
  <div class="node-content">
    <!-- End node has no content, just execution input -->
  </div>
</div>

<style>
  .blueprint-node {
    background: #1a1a1a;
    border-radius: 4px;
    overflow: hidden;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
    min-width: 140px;
  }

  .blueprint-node.selected {
    box-shadow: 0 0 0 2px #ffa500, 0 2px 8px rgba(0, 0, 0, 0.5);
  }

  .end-node .node-title-bar {
    background: linear-gradient(to bottom, #7f2d2d 0%, #5f1f1f 100%);
    border-bottom: 1px solid #9f3d3d;
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
    min-height: 20px;
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
</style>
