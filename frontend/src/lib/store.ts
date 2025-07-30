import { get, writable } from 'svelte/store';
import { type Node as XYNode, type Edge as XYEdge, type OnConnect, type Connection } from '@xyflow/svelte';
import { axon } from './proto/axon'; // Import generated TS types

// The main Svelte stores for xyflow's visual state
export const nodes = writable<XYNode[]>([]);
export const edges = writable<XYEdge[]>([]);

// The raw Axon graph object, our single source of truth.
// Correctly initialized with the create() factory method.
export const axonGraph = writable<axon.Graph>(axon.Graph.create());

/**
 * Converts an Axon Graph (from Protobuf) into nodes and edges for xyflow.
 * This function should be called when a file is loaded.
 * @param graph The Axon graph object from the backend.
 */
export function loadFromAxonGraph(graph: axon.Graph) {
	const xyNodes: XYNode[] = graph.nodes.map((axonNode: axon.Node) => {
		// Find or create visual info, defaulting to (0,0)
		const position = axonNode.visual_info ?? { x: 0, y: 0 };

		return {
			id: axonNode.id,
			type: 'custom', // We'll use a single custom node component
			position: { x: position.x, y: position.y },
			data: {
				// Pass the raw Axon node data directly to our Svelte component
				node: axonNode
			}
		};
	});

	const dataEdges: XYEdge[] = graph.data_edges.map((edge: axon.DataEdge) => ({
		id: `data-${edge.from_node_id}-${edge.from_port}-${edge.to_node_id}-${edge.to_port}`,
		source: edge.from_node_id,
		target: edge.to_node_id,
		sourceHandle: edge.from_port,
		targetHandle: edge.to_port,
		type: 'default',
		style: 'stroke: #00b4d4; stroke-width: 2;', // Data Edge Color
		animated: false
	}));

	const execEdges: XYEdge[] = graph.exec_edges.map((edge: axon.ExecEdge) => ({
		id: `exec-${edge.from_node_id}-${edge.to_node_id}`,
		source: edge.from_node_id,
		target: edge.to_node_id,
		sourceHandle: 'exec_out', // Standardized handle IDs
		targetHandle: 'exec_in',
		type: 'default',
		style: 'stroke: #fff; stroke-width: 2.5;', // Exec Edge Color
		animated: true
	}));

	nodes.set(xyNodes);
	edges.set([...dataEdges, ...execEdges]);
	axonGraph.set(graph);
}

/**
 * Syncs the state from xyflow (nodes, edges) back into the main axonGraph store.
 * This should be called before saving the file or running the transpiler.
 */
export function updateAxonGraph() {
	const currentNodes = get(nodes);
	const currentEdges = get(edges);
	const currentGraph = get(axonGraph);

	// Create a new graph object to avoid mutation issues.
	const newGraph = axon.Graph.clone(currentGraph);

	// 1. Update node positions from xyflow state
	const nodeMap = new Map<string, axon.Node>();
	newGraph.nodes.forEach((n) => nodeMap.set(n.id, n));

	for (const xyNode of currentNodes) {
		const axonNode = nodeMap.get(xyNode.id);
		if (axonNode) {
			if (!axonNode.visual_info) {
				axonNode.visual_info = { x: 0, y: 0, width: 0, height: 0 };
			}
			axonNode.visual_info.x = xyNode.position.x;
			axonNode.visual_info.y = xyNode.position.y;
		}
	}

	// 2. Rebuild edges from scratch from xyflow state
	newGraph.data_edges = [];
	newGraph.exec_edges = [];

	for (const edge of currentEdges) {
		if (!edge.source || !edge.target || !edge.sourceHandle || !edge.targetHandle) continue;

		if (edge.sourceHandle === 'exec_out') {
			// This is an Execution Edge
			newGraph.exec_edges.push(
				axon.ExecEdge{
					from_node_id: edge.source,
					to_node_id: edge.target
				});
		} else {
			// This is a Data Edge
			newGraph.data_edges.push(
				axon.DataEdge.create({
					from_node_id: edge.source,
					from_port: edge.sourceHandle,
					to_node_id: edge.target,
					to_port: edge.targetHandle
				})
			);
		}
	}

	// 3. Update the main store with the fully synced graph
	axonGraph.set(newGraph);
}