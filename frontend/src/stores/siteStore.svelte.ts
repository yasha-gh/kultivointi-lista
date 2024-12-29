import { list } from "$wails/go/models";
// import { Sync as SyncSettings, GetSelf } from "$wails/go/main/Settings";
import { NewDbID, SaveSite, GetSites } from "$wails/go/list/List";
import { setSaveStringVal } from "$stores";
import { type LabelValue } from "$lib/utils/index";
import { EventsOn } from "$wails/runtime";

type SiteValidateResult = {
  valid: boolean;
  errors: string[];
};

// ** MAIN SITE ** //
export class SignalSite {
  // #id = $state("");
  // get id() {
  //   return this.#id;
  // }
  // set id(val: string) {
  //   this.#id = val;
  // }
  // #url = $state("");
  // get url() {
  //   return this.#url;
  // }
  // set url(val: string) {
  //   this.#url = val;
  // }

  id = $state("");
  url = $state("");
  domainBase = $state("");
  domainTopLevel = $state("");
  domainProtocol = $state("");
  episodeTemplate = $state("");
  mainPageTempate = $state("");

  disableSaving = $state(false);
  constructor(goSite: list.Site) {
    // this.#id = goItem?.id;
    // if (!this.#id) {
    //   NewDbID()
    //     .then((newId) => {
    //       this.#id = newId;
    //     })
    //     .catch((e) => console.log("Failed to get DBID"));
    // }
    // this.#url = goItem.url;
    this.id = goSite.id;
    if (!this.id) {
      NewDbID().then((newId) => {
        this.id = newId;
      });
    }
    this.url = goSite.url;
    this.domainBase = goSite.domainBase;
    this.domainTopLevel = goSite.domainTopLevel;
    this.domainProtocol = goSite.domainProtocol;
    this.episodeTemplate = goSite.episodeTemplate;
    this.mainPageTempate = goSite.mainPageTempate;
    // $effect.root(() => {
    //   $effect(() => {
    //     // console.log("Changes", this.seasonNum);
    //   });
    // });
  }

  fromUrl = async (url: string) => {
    if (!url) {
      console.error("no URL provided");
      return;
    }
    console.log("From URL id on init", this.id);
    if (!this?.id) {
      this.id = await NewDbID();
      // NewDbID().then((newId) => {
      //   this.id = newId;
      // });
    }
    const protoSplit = url.split("://");
    let urlBody = "";
    switch (protoSplit.length) {
      case 1:
        this.domainProtocol = "https";
        urlBody = protoSplit[0];
        break;
      case 2:
        this.domainProtocol = protoSplit[0];
        urlBody = protoSplit[1];
      default:
        this.domainProtocol = "https";
        break;
    }
    if (!urlBody) {
      console.error(
        "failed to parse url body from url",
        "url",
        url,
        "protoSplit",
        protoSplit,
      );
      return;
    }
    const domainBodySplit = urlBody.split("/");
    if (!domainBodySplit) {
      console.error("Failed to get domain from urlBody");
      return;
    }
    let domain = "";
    if (domainBodySplit.length === 1) {
      domain = domainBodySplit[0];
    } else {
      domain = domainBodySplit[domainBodySplit.length - 1];
    }
    if (!domain) {
      console.error(
        "Failed to get domain from urlBody",
        "urlBody",
        urlBody,
        "domainBodySplit",
        domainBodySplit,
      );
    }

    const tooPartExtensions = [
      ".co.uk",
      ".com.au",
      ".net.au",
      ".org.au",
      ".gov.au",
      ".edu.au",
      ".ac.uk",
    ];
    let tldFound = false;
    for (const tld of tooPartExtensions) {
      if (domain.includes(tld)) {
        tldFound = true;
        domain = domain.replace(tld, "");
        this.domainBase = domain;
        this.domainTopLevel = tld;
        break;
      }
    }
    if (!tldFound) {
      const domainParts = domain.split(".");
      if (domainParts.length > 1) {
        this.domainTopLevel = domainParts[domainParts.length - 1];
        const tld = `.${this.domainTopLevel}`;
        domain = domain.replace(tld, "");
        this.domainBase = domain;
      }
    }
    this.setURL();
    console.log("After URL parsing", this);
  };

  setURL = () => {
    this.url = "";
    if (this.domainProtocol) {
      this.url = `${this.domainProtocol}://`;
    }
    this.url = `${this.url}${this.domainBase}.${this.domainTopLevel}`;
  };

  getName = () => {
    return `${this.domainBase}.${this.domainTopLevel}`;
  };

  validate = () => {
    const res: SiteValidateResult = {
      valid: true,
      errors: [],
    };
    if (!this.id) {
      res.errors.push("ID Puuttuu");
      res.valid = false;
    }
    if (!this?.domainBase) {
      res.errors.push("Domain nimi puuttuu");
      res.valid = false;
    }
    if (!this?.domainTopLevel) {
      res.errors.push("Domain pääte puuttuu");
      res.valid = false;
    }
    if (!res.valid) {
      console.log("Not valid site", this);
    }
    return res;
  };

  saveSelf = async () => {
    if (!this.disableSaving) {
      console.log("save self param", list.Site.createFrom(this));
      try {
        await SaveSite(list.Site.createFrom(this));
      } catch (e) {
        console.error("failed to save site", e);
      }
    }
  };
}

export type SiteSelectOption = {
  label: string;
  siteId: string;
};
/*** SITE STORE ***/
export class SiteStore {
  sites = $state<SignalSite[]>([]);
  selectOptions = $derived.by(() => {
    const selectOptions: SiteSelectOption[] = [
      {
        label: "Ei valintaa",
        siteId: "NOSELECT",
      },
    ];
    for (const site of siteStore.sites) {
      selectOptions.push({
        label: site.getName(),
        siteId: site.id,
      });
    }
    return selectOptions;
  });
  constructor() {
    setTimeout(() => {
      this.GetSites();
    }, 1000);

    // svelte effects
    $effect.root(() => {
      $effect(() => {
        // console.log("Changes", this.list);
      });
    });
  }

  getOptionValue(siteId: string) {
    const option = this.selectOptions.find((o) => o.siteId === siteId);
    return (option) ? option : {
      label: "Ei valintaa",
      siteId: "NOSELECT",
    };
  }
  getSiteByID(siteId: string) {
    return this.sites.find((s) => s.id === siteId);
  }
  async GetSites() {
    try {
      const goSites = await GetSites();
      for (const goSite of goSites) {
        const site = new SignalSite(goSite);
        this.sites.push(site);
      }
    } catch (e: any) {
      console.error("Failed to get sites", e);
    }
  }
}

export const siteStore = new SiteStore();
