<script lang="ts">
	import { Handle, Position, type NodeProps } from '@xyflow/svelte';
	import { axon } from '../proto/axon';

	// The `data` prop is passed automatically by SvelteFlow
	export let data: { node: axon.Node };

	// Define colors for different node types for visual distinction
	const nodeColors: Record<axon.NodeType, string> = {
		[axon.NodeType.NODE_UNKNOWN]: '#64748b',
		[axon.NodeType.START]: '#10b981',
		[axon.NodeType.END]: '#ef4444',
		[axon.NodeType.RETURN]: '#d946ef',
		[axon.NodeType.FUNC_DEF]: '#22c55e',
		[axon.NodeType.STRUCT_DEF]: '#0ea5e9',
		[axon.NodeType.CONSTANT]: '#3b82f6',
		[axon.NodeType.FUNCTION]: '#6366f1',
		[axon.NodeType.OPERATOR]: '#f97316',
		[axon.NodeType.IGNORE]: '#ec4899'
	};

	const dataPortColors: Record<string, string> = {
		int: '#00b4d4',
		string: '#d946ef',
		bool: '#ef4444',
		'[]byte': '#84cc16',
		error: '#f59e0b',
		default: '#8b5cf6'
	};

	$: node = data.node;
	$: headerColor = nodeColors[node.type] || nodeColors[axon.NodeType.NODE_UNKNOWN];

	const getPortColor = (typeName: string) => dataPortColors[typeName] || dataPortColors.default;
</script>

<div class="node-wrapper">
	<!-- Node Header -->
	<div class="node-header" style="background-color: {headerColor};">
		{node.label}
	</div>

	<!-- Node Body -->
	<div class="node-body">
		<!-- Execution Ports (conditionally rendered) -->
		{#if node.type !== axon.NodeType.START}
			<Handle type="target" position={Position.Left} id="exec_in" class="handle exec-handle" />
		{/if}
		{#if node.type !== axon.NodeType.END && node.type !== axon.NodeType.RETURN}
			<Handle type="source" position={Position.Right} id="exec_out" class="handle exec-handle" />
		{/if}

		<!-- Data Ports -->
		<div class="ports-container">
			<!-- Input Ports -->
			<div class="ports-column">
				{#each node.inputs as port}
					<div class="port">
						<Handle
							type="target"
							position={Position.Left}
							id={port.name}
							class="handle data-handle"
							style="background-color: {getPortColor(port.type_name)};"
						/>
						<span class="port-label">{port.name}</span>
					</div>
				{/each}
			</div>

			<!-- Output Ports -->
			<div class="ports-column">
				{#each node.outputs as port}
					<div class="port port-output">
						<span class="port-label">{port.name}</span>
						<Handle
							type="source"
							position={Position.Right}
							id={port.name}
							class="handle data-handle"
							style="background-color: {getPortColor(port.type_name)};"
						/>
					</div>
				{/each}
			</div>
		</div>

		<!-- Extra Info (like function implementation reference) -->
		{#if node.impl_reference}
			<div class="node-footer">
				{node.impl_reference}
			</div>
		{/if}
	</div>
</div>

<style>
	.node-wrapper {
		background-color: #2d3748;
		border: 1px solid #1a202c;
		border-radius: 8px;
		width: 220px;
		font-family: sans-serif;
		color: #e2e8f0;
		box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
	}
	.node-header {
		padding: 6px 10px;
		border-top-left-radius: 7px;
		border-top-right-radius: 7px;
		font-weight: bold;
		color: white;
	}
	.node-body {
		position: relative;
		padding: 10px;
	}
	.ports-container {
		display: flex;
		justify-content: space-between;
	}
	.ports-column {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}
	.port {
		display: flex;
		align-items: center;
		gap: 6px;
	}
	.port-output {
		justify-content: flex-end;
	}
	.port-label {
		font-size: 12px;
	}
	.handle {
		width: 12px;
		height: 12px;
		border-radius: 3px;
		border: 1px solid #1a202c;
	}
	.exec-handle {
		background-color: white;
		clip-path: polygon(0 0, 100% 50%, 0 100%); /* Arrow shape */
		border-radius: 0;
		height: 14px;
	}
	.data-handle {
		border-radius: 50%;
	}
	.node-footer {
		margin-top: 8px;
		padding-top: 8px;
		border-top: 1px solid #4a5568;
		font-size: 11px;
		color: #a0aec0;
		text-align: center;
	}
</style>