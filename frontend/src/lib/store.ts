import { writable } from 'svelte/store';
import { type Node, type Edge } from '@xyflow/svelte';
import { Graph, Node as AxonNode, DataEdge, ExecEdge } from './proto/axon.pb'; // Import generated TS types

// The main Svelte store for our application
export const nodes = writable<Node[]>([]);
export const edges = writable<Edge[]>([]);

// The raw Axon graph object, our single source of truth
export const axonGraph = writable<Graph>(Graph.create());

/**
 * Converts an Axon Graph (from Protobuf) into nodes and edges for xyflow.
 * @param graph The Axon graph object.
 */
export function loadFromAxonGraph(graph: Graph) {
    const xyNodes: Node[] = graph.nodes.map((axonNode: AxonNode) => {
        // Find or create visual info
        const position = axonNode.visualInfo ?? { x: 0, y: 0 };

        return {
            id: axonNode.id,
            type: 'custom', // We'll use a single custom node component
            position: { x: position.x, y: position.y },
            data: {
                // Pass the raw Axon node data to our Svelte component
                node: axonNode
            }
        };
    });

    const dataEdges: Edge[] = graph.dataEdges.map((edge: DataEdge) => ({
        id: `data-${edge.fromNodeId}-${edge.fromPort}-${edge.toNodeId}-${edge.toPort}`,
        source: edge.fromNodeId,
        target: edge.toNodeId,
        sourceHandle: edge.fromPort,
        targetHandle: edge.toPort,
        type: 'default',
        style: 'stroke: #00b4d4; stroke-width: 2;', // Example styling
        animated: false
    }));

    const execEdges: Edge[] = graph.execEdges.map((edge: ExecEdge) => ({
        id: `exec-${edge.fromNodeId}-${edge.toNodeId}`,
        source: edge.fromNodeId,
        target: edge.toNodeId,
        sourceHandle: 'exec_out', // Standardized handle IDs
        targetHandle: 'exec_in',
        type: 'default',
        style: 'stroke: #fff; stroke-width: 2.5;',
        animated: true
    }));

    nodes.set(xyNodes);
    edges.set([...dataEdges, ...execEdges]);
    axonGraph.set(graph);
}

// TODO: Create a function `updateAxonGraph()` that converts xyflow state back to an Axon Graph
// This would be called before saving the file.