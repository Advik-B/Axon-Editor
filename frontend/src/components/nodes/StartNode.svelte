<script>
  import { Handle, Position } from '@xyflow/svelte';
  
  export let data;
  export let selected = false;
</script>

<div class="blueprint-node start-node" class:selected>
  <div class="node-title-bar">
    <span class="node-title">{data.label || 'Start'}</span>
  </div>
  <div class="node-content">
    <!-- Start node has no content, just execution output -->
  </div>
  <Handle type="source" position={Position.Right} id="exec" class="exec-handle" />
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

  .start-node .node-title-bar {
    background: linear-gradient(to bottom, #2d5f2d 0%, #1f441f 100%);
    border-bottom: 1px solid #3d7f3d;
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

  :global(.exec-handle) {
    width: 0;
    height: 0;
    border-style: solid;
    border-width: 6px 0 6px 10px;
    border-color: transparent transparent transparent #ffffff;
    background: transparent !important;
    border-radius: 0 !important;
  }
</style>
