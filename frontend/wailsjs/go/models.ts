export namespace main {
	
	export class CleanTask {
	    id: string;
	    name: string;
	    path: string;
	    mode: string;
	    retentionDays: number;
	    filePattern: string;
	    cronSpec: string;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CleanTask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.mode = source["mode"];
	        this.retentionDays = source["retentionDays"];
	        this.filePattern = source["filePattern"];
	        this.cronSpec = source["cronSpec"];
	        this.enabled = source["enabled"];
	    }
	}
	export class TaskResult {
	    success: number;
	    failed: number;
	    errors: string[];
	
	    static createFrom(source: any = {}) {
	        return new TaskResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.failed = source["failed"];
	        this.errors = source["errors"];
	    }
	}

}

