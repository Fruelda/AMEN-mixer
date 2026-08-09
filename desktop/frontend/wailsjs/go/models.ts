export namespace models {
	
	export class Channel {
	    id: number;
	    name: string;
	    app: string;
	    volume: number;
	    muted: boolean;
	    connected: boolean;
	    selected: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Channel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.app = source["app"];
	        this.volume = source["volume"];
	        this.muted = source["muted"];
	        this.connected = source["connected"];
	        this.selected = source["selected"];
	    }
	}

}

