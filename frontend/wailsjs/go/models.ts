export namespace list {
	
	export class Site {
	    id: string;
	    url: string;
	    domainBase: string;
	    domainTopLevel: string;
	    domainProtocol: string;
	    episodeTemplate: string;
	    mainPageTempate: string;
	
	    static createFrom(source: any = {}) {
	        return new Site(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.url = source["url"];
	        this.domainBase = source["domainBase"];
	        this.domainTopLevel = source["domainTopLevel"];
	        this.domainProtocol = source["domainProtocol"];
	        this.episodeTemplate = source["episodeTemplate"];
	        this.mainPageTempate = source["mainPageTempate"];
	    }
	}
	export class EpisodeSeen {
	    id: string;
	    episodesSeen: number;
	    siteId: string;
	    site?: Site;
	    itemId: string;
	
	    static createFrom(source: any = {}) {
	        return new EpisodeSeen(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.episodesSeen = source["episodesSeen"];
	        this.siteId = source["siteId"];
	        this.site = this.convertValues(source["site"], Site);
	        this.itemId = source["itemId"];
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
	export class ListItemTitle {
	    id: string;
	    title: string;
	    lang: string;
	    primaryTitle: boolean;
	    itemId: string;
	
	    static createFrom(source: any = {}) {
	        return new ListItemTitle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.lang = source["lang"];
	        this.primaryTitle = source["primaryTitle"];
	        this.itemId = source["itemId"];
	    }
	}
	export class ListItem {
	    id: string;
	    titleId: string;
	    title?: ListItemTitle;
	    titles: ListItemTitle[];
	    type: string;
	    broadcastType: string;
	    episodesTotal: number;
	    episodesSeen: number;
	    ongoing: boolean;
	    seasonNum: number;
	    seasons: ListItem[];
	    parentItemId: string;
	    episodesSeenOn: EpisodeSeen[];
	    thubmnailImageId: string;
	    // Go type: time
	    modifiedAt: any;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new ListItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.titleId = source["titleId"];
	        this.title = this.convertValues(source["title"], ListItemTitle);
	        this.titles = this.convertValues(source["titles"], ListItemTitle);
	        this.type = source["type"];
	        this.broadcastType = source["broadcastType"];
	        this.episodesTotal = source["episodesTotal"];
	        this.episodesSeen = source["episodesSeen"];
	        this.ongoing = source["ongoing"];
	        this.seasonNum = source["seasonNum"];
	        this.seasons = this.convertValues(source["seasons"], ListItem);
	        this.parentItemId = source["parentItemId"];
	        this.episodesSeenOn = this.convertValues(source["episodesSeenOn"], EpisodeSeen);
	        this.thubmnailImageId = source["thubmnailImageId"];
	        this.modifiedAt = this.convertValues(source["modifiedAt"], null);
	        this.createdAt = this.convertValues(source["createdAt"], null);
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
	
	export class ListSearchOptions {
	    searchFields: string[];
	    orderField: string;
	    orderDirection: string;
	
	    static createFrom(source: any = {}) {
	        return new ListSearchOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.searchFields = source["searchFields"];
	        this.orderField = source["orderField"];
	        this.orderDirection = source["orderDirection"];
	    }
	}

}

export namespace main {
	
	export class Settings {
	    dbFilename: string;
	    dbFileDir: string;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dbFilename = source["dbFilename"];
	        this.dbFileDir = source["dbFileDir"];
	    }
	}

}

