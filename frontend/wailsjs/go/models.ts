export namespace lsp {
	
	export class CompletionItem {
	    label: string;
	    kind: string;
	    detail: string;
	    insert_text: string;
	
	    static createFrom(source: any = {}) {
	        return new CompletionItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.label = source["label"];
	        this.kind = source["kind"];
	        this.detail = source["detail"];
	        this.insert_text = source["insert_text"];
	    }
	}
	export class Diagnostic {
	    node_id: string;
	    severity: string;
	    message: string;
	    code?: string;
	    line?: number;
	    column?: number;
	
	    static createFrom(source: any = {}) {
	        return new Diagnostic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.node_id = source["node_id"];
	        this.severity = source["severity"];
	        this.message = source["message"];
	        this.code = source["code"];
	        this.line = source["line"];
	        this.column = source["column"];
	    }
	}

}

