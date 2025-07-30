export namespace axon {
	
	export class Comment {
	    id?: string;
	    content?: string;
	
	    static createFrom(source: any = {}) {
	        return new Comment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.content = source["content"];
	    }
	}
	export class DataEdge {
	    from_node_id?: string;
	    from_port?: string;
	    to_node_id?: string;
	    to_port?: string;
	
	    static createFrom(source: any = {}) {
	        return new DataEdge(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.from_node_id = source["from_node_id"];
	        this.from_port = source["from_port"];
	        this.to_node_id = source["to_node_id"];
	        this.to_port = source["to_port"];
	    }
	}
	export class ExecEdge {
	    from_node_id?: string;
	    to_node_id?: string;
	
	    static createFrom(source: any = {}) {
	        return new ExecEdge(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.from_node_id = source["from_node_id"];
	        this.to_node_id = source["to_node_id"];
	    }
	}
	export class VisualInfo {
	    x?: number;
	    y?: number;
	    width?: number;
	    height?: number;
	
	    static createFrom(source: any = {}) {
	        return new VisualInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	        this.width = source["width"];
	        this.height = source["height"];
	    }
	}
	export class Port {
	    name?: string;
	    type_name?: string;
	
	    static createFrom(source: any = {}) {
	        return new Port(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type_name = source["type_name"];
	    }
	}
	export class Node {
	    id?: string;
	    type?: number;
	    label?: string;
	    inputs?: Port[];
	    outputs?: Port[];
	    impl_reference?: string;
	    config?: Record<string, string>;
	    visual_info?: VisualInfo;
	    comment_ids?: string[];
	
	    static createFrom(source: any = {}) {
	        return new Node(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.label = source["label"];
	        this.inputs = this.convertValues(source["inputs"], Port);
	        this.outputs = this.convertValues(source["outputs"], Port);
	        this.impl_reference = source["impl_reference"];
	        this.config = source["config"];
	        this.visual_info = this.convertValues(source["visual_info"], VisualInfo);
	        this.comment_ids = source["comment_ids"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Graph {
	    id?: string;
	    name?: string;
	    imports?: string[];
	    nodes?: Node[];
	    data_edges?: DataEdge[];
	    exec_edges?: ExecEdge[];
	    comments?: Comment[];
	
	    static createFrom(source: any = {}) {
	        return new Graph(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.imports = source["imports"];
	        this.nodes = this.convertValues(source["nodes"], Node);
	        this.data_edges = this.convertValues(source["data_edges"], DataEdge);
	        this.exec_edges = this.convertValues(source["exec_edges"], ExecEdge);
	        this.comments = this.convertValues(source["comments"], Comment);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	

}

