<script lang="ts">
	import { onMount } from 'svelte';
	import AxonCanvas from './lib/components/AxonCanvas.svelte';
	import { axonGraph, loadFromAxonGraph, updateAxonGraph } from './lib/store';
	import { axon } from './lib/proto/axon';

	// --- Mock Backend & Data ---
	// This simulates what we would get from the Wails backend.
	let transpiledCode = '// Click "Transpile" to see Go code here';

	// Mock data for the 'add.ax' example
	const mockAddGraph: axon.Graph = {
		id: 'basic-addition-v2',
		name: 'Add Numbers with Execution Flow',
		imports: ['fmt'],
		nodes: [
			{ id: 'start', type: axon.NodeType.START, label: 'Start', inputs: [], outputs: [], impl_reference: '', config: new Map(), comment_ids: [] },
			{ id: 'end', type: axon.NodeType.END, label: 'End', inputs: [], outputs: [], impl_reference: '', config: new Map(), comment_ids: [] },
			{
				id: 'const1',
				type: axon.NodeType.CONSTANT,
				label: 'x',
				outputs: [{ name: 'out', type_name: 'int' }],
				config: new Map([['value', '5']]),
				inputs: [], impl_reference: '', comment_ids: []
			},
			{
				id: 'const2',
				type: axon.NodeType.CONSTANT,
				label: 'y',
				outputs: [{ name: 'out', type_name: 'int' }],
				config: new Map([['value', '3']]),
				inputs: [], impl_reference: '', comment_ids: []
			},
			{
				id: 'sum',
				type: axon.NodeType.OPERATOR,
				label: 'z',
				inputs: [
					{ name: 'a', type_name: 'int' },
					{ name: 'b', type_name: 'int' }
				],
				outputs: [{ name: 'out', type_name: 'int' }],
				config: new Map([['op', '+']]),
				impl_reference: '', comment_ids: []
			},
			{
				id: 'printer',
				type: axon.NodeType.FUNCTION,
				label: 'PrintResult',
				impl_reference: 'fmt.Println',
				inputs: [{ name: 'a', type_name: 'int' }],
				outputs: [], config: new Map(), comment_ids: []
			}
		],
		data_edges: [
			{ from_node_id: 'const1', from_port: 'out', to_node_id: 'sum', to_port: 'a' },
			{ from_node_id: 'const2', from_port: 'out', to_node_id: 'sum', to_port: 'b' },
			{ from_node_id: 'sum', from_port: 'out', to_node_id: 'printer', to_port: 'a' }
		],
		exec_edges: [
			{ from_node_id: 'start', to_node_id: 'sum' },
			{ from_node_id: 'sum', to_node_id: 'printer' },
			{ from_node_id: 'printer', to_node_id: 'end' }
		],
		comments: []
	};

	// --- Component Logic ---
	// On component mount, load our mock data into the stores
	onMount(() => {
		handleLoad();
	});

	function handleLoad() {
		// In a real app, this would call a Wails function:
		// const graph = await Wails.LoadGraphFromFile(...)
		const graph = mockAddGraph;
		loadFromAxonGraph(graph);
		transpiledCode = '// Graph loaded. Click "Transpile" to generate code.';
	}

	async function handleSave() {
		updateAxonGraph(); // Sync visual state to the axonGraph store
		const graphToSave = $axonGraph;
		console.log('Saving graph:', graphToSave);
		// In a real app, this would call a Wails function:
		// await Wails.SaveGraphToFile('path/to/file.ax', graphToSave);
		alert('Graph saved to console! (See developer tools)');
	}

	async function handleTranspile() {
		updateAxonGraph(); // Sync visual state first!
		const graphToTranspile = $axonGraph;
		console.log('Transpiling graph:', graphToTranspile);
		// In a real app, this would call a Wails function:
		// transpiledCode = await Wails.TranspileGraph(graphToTranspile);
		transpiledCode = `package main\n\nimport (\n\t"fmt"\n)\n\nfunc main() {\n\tx := 5\n\ty := 3\n\tz := x + y\n\tfmt.Println(z)\n}\n`;
	}
</script>

<main>
	<div class="canvas-container">
		<AxonCanvas />
	</div>

	<div class="sidebar">
		<div class="controls">
			<h3>Axon Editor</h3>
			<button on:click={handleLoad}>Load Mock</button>
			<button on:click={handleSave}>Save</button>
			<button on:click={handleTranspile}>Transpile</button>
		</div>
		<div class="code-panel">
			<h4>Generated Go Code</h4>
			<pre>{transpiledCode}</pre>
		</div>
	</div>
</main>

<style>
	:global(body) {
		margin: 0;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell,
			'Open Sans', 'Helvetica Neue', sans-serif;
	}
	main {
		display: flex;
		height: 100vh;
	}
	.canvas-container {
		flex-grow: 1;
	}
	.sidebar {
		width: 350px;
		background-color: #1f2937;
		color: #d1d5db;
		display: flex;
		flex-direction: column;
		border-left: 1px solid #4b5563;
	}
	.controls {
		padding: 1rem;
		border-bottom: 1px solid #4b5563;
	}
	.controls h3 {
		margin-top: 0;
	}
	.controls button {
		background-color: #4f46e5;
		color: white;
		border: none;
		padding: 8px 16px;
		border-radius: 6px;
		cursor: pointer;
		margin-right: 8px;
	}
	.controls button:hover {
		background-color: #4338ca;
	}
	.code-panel {
		padding: 1rem;
		flex-grow: 1;
		overflow-y: auto;
	}
	.code-panel h4 {
		margin-top: 0;
	}
	pre {
		background-color: #111827;
		padding: 1rem;
		border-radius: 6px;
		white-space: pre-wrap;
		word-wrap: break-word;
		font-size: 13px;
	}
</style>