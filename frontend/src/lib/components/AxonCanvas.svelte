<script lang="ts">
	import { SvelteFlow, Background, Controls, MiniMap, type Connection } from '@xyflow/svelte';
	import { nodes, edges } from '../store';
	import CustomNode from './CustomNode.svelte';

	// Tell SvelteFlow to use our CustomNode for nodes of type 'custom'
	const nodeTypes = { custom: CustomNode };

	// When the user connects two handles, this function is called.
	// We create a new edge and add it to our svelte store.
	const handleConnect = (event: CustomEvent<Connection>) => {
		const newEdge = { ...event.detail, id: `edge-${Math.random()}` };
		edges.update((currentEdges) => [...currentEdges, newEdge]);
	};
</script>

<div style="height: 100vh; width: 100%; background-color: #1a202c;">
	<SvelteFlow {nodeTypes} {nodes} {edges} on:connect={handleConnect} fitView>
		<Background />
		<Controls />
		<MiniMap />
	</SvelteFlow>
</div>