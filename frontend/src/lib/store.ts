import { get, writable } from 'svelte/store';
import { type Node as XYNode, type Edge as XYEdge } from '@xyflow/svelte';
// The import is now correct based on your file structure.
import { axon } from './proto/axon';

// The main Svelte stores for xyflow's visual state
export const nodes = writable<XYNode[]>([]);
export const edges = writable<XYEdge[]>([]);

// The raw Axon graph object, our single source of truth.
// This is the correct way to initialize it using the 'new' constructor.
export const axonGraph = writable<axon.Graph>(new axon.Graph());

/**
 * Converts an Axon Graph (from Protobuf) into nodes and edges for xyflow.
 * This function should be called when a file is loaded.
 * @param graph The Axon graph object from the backend.
 */
export function loadFromAxonGraph(graph: axon.Graph) {
	const xyNodes: XYNode[] = graph.nodes.map((axonNode: axon.Node) => {
		const position = axonNode.visual_info ?? { x: 0, y: 0 };

		return {
			id: axonNode.id,
			type: 'custom',
			position: { x: position.x, y: position.y },
			data: { node: axonNode }
		};
	});

	const dataEdges: XYEdge[] = graph.data_edges.map((edge: axon.DataEdge) => ({
		id: `data-${edge.from_node_id}-${edge.from_port}-${edge.to_node_id}-${edge.to_port}`,
		source: edge.from_node_id,
		target: edge.to_node_id,
		sourceHandle: edge.from_port,
		targetHandle: edge.to_port,
		type: 'default',
		style: 'stroke: #00b4d4; stroke-width: 2;',
		animated: false
	}));

	const execEdges: XYEdge[] = graph.exec_edges.map((edge: axon.ExecEdge) => ({
		id: `exec-${edge.from_node_id}-${edge.to_node_id}`,
		source: edge.from_node_id,
		target: edge.to_node_id,
		sourceHandle: 'exec_out',
		targetHandle: 'exec_in',
		type: 'default',
		style: 'stroke: #fff; stroke-width: 2.5;',
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

	// To "clone", we convert the current graph to a plain object and create a new instance from it.
	const newGraph = axon.Graph.fromObject(currentGraph.toObject());

	// 1. Update node positions from xyflow state
	const nodeMap = new Map<string, axon.Node>();
	newGraph.nodes.forEach((n) => nodeMap.set(n.id, n));

	for (const xyNode of currentNodes) {
		const axonNode = nodeMap.get(xyNode.id);
		if (axonNode) {
			if (!axonNode.visual_info) {
				// Create the optional sub-message using the 'new' constructor.
				axonNode.visual_info = new axon.VisualInfo();
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

		if (edge.sourceHandle === 'exec_out' && edge.targetHandle === 'exec_in') {
			// This is an Execution Edge
			newGraph.exec_edges.push(
				new axon.ExecEdge({
					from_node_id: edge.source,
					to_node_id: edge.target
				})
			);
		} else {
			// This is a Data Edge
			newGraph.data_edges.push(
				new axon.DataEdge({
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