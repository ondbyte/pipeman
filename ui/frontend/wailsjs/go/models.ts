export namespace main {
	
	export class Arg {
	    Name: string;
	
	    static createFrom(source: any = {}) {
	        return new Arg(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	    }
	}

}

