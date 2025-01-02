import { list } from "$wails/go/models";
// import { Sync as SyncSettings, GetSelf } from "$wails/go/main/Settings";
import {
  NewDbID,
  GetBaseItemIDs,
  GetListItemsByIDs,
  SaveListItem,
  Search,
} from "$wails/go/list/List";
import {
  objectFromGo,
  setSaveBoolVal,
  setSaveNumVal,
  setSaveStringVal,
  SignalSite,
} from "$stores";
import { type LabelValue, newDBID } from "$lib/utils/index";
import { EventsOn } from "$wails/runtime";

function setSaveTitleVal(
  val: SignalListItemTitle | undefined,
  compareVal: SignalListItemTitle | undefined,
): { val: SignalListItemTitle; hasChanges: boolean } {
  const ret = {
    val: val,
    hasChanges: false,
  };
  if (compareVal === undefined && ret.val === undefined) {
    return ret;
  }
  if (ret.val) {
    if (!compareVal) {
      ret.hasChanges = true;
    } else {
      for (const [k, v] of Object.entries(ret.val)) {
        if (ret.val[k] !== compareVal[k]) {
          ret.hasChanges = true;
        }
      }
    }
  }

  return ret;
}

function setSaveSiteVal(
  val: SignalSite | undefined,
  compareVal: SignalSite | undefined,
): { val: SignalSite | undefined; hasChanges: boolean } {
  const ret = {
    val: val,
    hasChanges: false,
  };
  if (compareVal === undefined && ret.val === undefined) {
    return ret;
  }
  if (ret.val) {
    if (!compareVal) {
      ret.hasChanges = true;
    } else {
      for (const [k, v] of Object.entries(ret.val)) {
        if (ret.val[k] !== compareVal[k]) {
          ret.hasChanges = true;
        }
      }
    }
  }
  return ret;
}

function setSaveEpisodeVal(
  val: SignalEpisodeSeen | undefined,
  compareVal: SignalEpisodeSeen | undefined,
): { val: SignalEpisodeSeen | undefined; hasChanges: boolean } {
  const ret = {
    val: val,
    hasChanges: false,
  };
  if (compareVal === undefined && ret.val === undefined) {
    return ret;
  }
  if (ret.val) {
    if (!compareVal) {
      ret.hasChanges = true;
    } else {
      for (const [k, v] of Object.entries(ret.val)) {
        if (k === "site") {
          const titleChange = setSaveSiteVal(v, compareVal.site);
          if (titleChange.hasChanges) ret.hasChanges = true;

          console.log(
            "setSave site",
            "v",
            v,
            "compareVal.site",
            compareVal.site,
            "titleChange",
            titleChange,
          );
        }
        if (ret.val[k] !== compareVal[k]) {
          ret.hasChanges = true;
        }
      }
    }
  }
  return ret;
}

export const serieTitleLanguages: LabelValue[] = [
  { label: "Englanti", value: "en" },

  { label: "Kiina (Romaji)", value: "zh_romanji" },
  { label: "Kiina", value: "zh" },
  { label: "Japani", value: "jp" },
  { label: "Japani (Romaji)", value: "jp_romaji" },
  { label: "Muu", value: "other" },
];

export const broadcastTypes: LabelValue[] = [
  { label: "Sarja (TV)", value: "TV" },
  { label: "Sarja (OVA)", value: "OVA" },
  { label: "Sarja (ONA)", value: "ONA" },
  { label: "Leffa (Teatteri)", value: "movie" },
  { label: "Leffa (TV)", value: "TV_movie" },
  { label: "Leffa (OVA)", value: "OVA_movie" },
  { label: "Leffa (ONA)", value: "ONA_movie" },
];

export class SignalListItemTitle {
  id = $state("");
  title = $state("");
  lang = $state("");
  primaryTitle = $state(false);
  itemId = $state("");
  constructor(goTitle: list.ListItemTitle = new list.ListItemTitle()) {
    this.id = goTitle.id;
    if (this.id === "") {
      NewDbID().then((newId) => (this.id = newId));
    }
    this.title = goTitle.title;
    this.lang = goTitle.lang;
    this.primaryTitle = goTitle.primaryTitle;
    this.itemId = goTitle.itemId;
    // svelte effects
    $effect.root(() => {
      $effect(() => {
        // console.log("Changes", this.seasonNum);
      });
    });
  }
  public static newDefault(
    itemId = "",
    primaryTitle = false,
  ): SignalListItemTitle {
    const newTitle = new SignalListItemTitle(new list.ListItemTitle());
    newTitle.id = newDBID();
    newTitle.itemId = itemId;
    newTitle.lang = "zh_romanji";
    newTitle.primaryTitle = primaryTitle;
    return newTitle;
  }
}

export class SignalEpisodeSeen {
  id = $state("");
  #episodesSeen = $state(0);
  get episodesSeen() {
    return this.#episodesSeen;
  }
  set episodesSeen(val: string | number) {
    val = Number(val);
    if (Number.isNaN(val)) {
      val = 0;
    }
    this.#episodesSeen = Number(val);
  }
  siteId = $state("");
  itemId = $state("");
  site = $state<SignalSite | undefined>(undefined);

  public static newDefault(itemId = ""): SignalEpisodeSeen {
    const newEpSeen = new SignalEpisodeSeen(new list.EpisodeSeen());
    newEpSeen.id = newDBID();
    newEpSeen.itemId = itemId;
    newEpSeen.episodesSeen = 0;
    return newEpSeen;
  }

  constructor(goEpisodeSeen: list.EpisodeSeen = new list.EpisodeSeen()) {
    this.id = goEpisodeSeen.id;
    if (this.id === "") {
      NewDbID().then((newId) => (this.id = newId));
    }
    this.episodesSeen = goEpisodeSeen.episodesSeen;
    this.siteId = goEpisodeSeen.siteId;
    this.itemId = goEpisodeSeen.itemId;
    // console.log(
    //   "init site",
    //   goEpisodeSeen?.site?.domainBase,
    //   goEpisodeSeen.site,
    // );
    if (goEpisodeSeen?.site) {
      this.site = new SignalSite(goEpisodeSeen.site);
    } else {
      this.site = undefined;
    }
    // console.log("this site", this.site);
    // svelte effects
    $effect.root(() => {
      $effect(() => {
        // console.log("Changes", this.seasonNum);
      });
    });
  }
}

// ** MAIN ITEM ** //
export class SignalListItem {
  id = $state("");

  #titleId = $state("");
  get titleId() {
    return this.#titleId;
  }
  set titleId(val: string) {
    const check = setSaveStringVal(val, this.#titleId);
    this.#titleId = check.val;
    if (check.hasChanges) {
      this.saveSelf();
    }
  }

  #title = $state<SignalListItemTitle>(
    new SignalListItemTitle(new list.ListItemTitle()),
  );
  get title() {
    return this.#title;
  }
  set title(val: SignalListItemTitle) {
    const check = setSaveTitleVal(val, this.#title);
    this.#title = check.val;
    console.log("title set", val, check.val);
    if (check.hasChanges) {
      this.saveSelf();
    }
  }

  #titles = $state<SignalListItemTitle[]>([]);
  get titles() {
    return this.#titles;
  }
  set titles(newTitles: SignalListItemTitle[]) {
    const titleList: SignalListItemTitle[] = [];
    let hasChanges = false;
    for (const title of newTitles) {
      const currentTitle = this.#titles.find((t) => t.id === title.id);
      if (currentTitle) {
        const newTitleCheck = setSaveTitleVal(title, currentTitle);
        if (newTitleCheck.hasChanges) hasChanges = true;
        if (newTitleCheck.val) {
          titleList.push(newTitleCheck.val);
        }
      }
    }
    this.#titles = titleList;
    if (hasChanges) {
      this.saveSelf();
    }
  }

  #type = $state("");
  get type() {
    return this.#type;
  }
  set type(val: string) {
    const check = setSaveStringVal(val, this.#type);
    if (check.hasChanges) {
      this.saveSelf();
    }
  }

  #broadcastType = $state("");
  get broadcastType() {
    return this.#broadcastType;
  }
  set broadcastType(val: string) {
    const check = setSaveStringVal(val, this.#broadcastType);
    this.#broadcastType = check.val;
    if (check.hasChanges) {
      this.saveSelf();
    }
  }

  episodesTotal = $state(0);
  episodesSeen = $state(0);

  #ongoing = $state(false);
  get ongoing() {
    return this.#ongoing;
  }
  set ongoing(val: boolean) {
    const check = setSaveBoolVal(val, this.#ongoing);
    this.#ongoing = check.val;
    if (check.hasChanges) {
      this.saveSelf();
    }
  }

  #seasonNum = $state(0);
  get seasonNum() {
    return this.#seasonNum;
  }
  set seasonNum(val: string | number) {
    const check = setSaveNumVal(val, this.#seasonNum);
    this.#seasonNum = check.val;
    if (check.hasChanges) {
      this.saveSelf();
    }
  }
  // #seasons = z.array(z.any()).default([]),
  #parentItemId = $state("");
  get parentItemId() {
    return this.#parentItemId;
  }
  set parentItemId(val: string) {
    const check = setSaveStringVal(val, this.#parentItemId);
    this.#parentItemId = check.val;
    if (check.hasChanges) {
      this.saveSelf();
    }
  }

  // TODO Fix reactivity
  #episodesSeenOn = $state<SignalEpisodeSeen[]>([]);
  get episodesSeenOn() {
    // console.log("episodesSeenOn getter");
    return this.#episodesSeenOn;
  }
  set episodesSeenOn(val: SignalEpisodeSeen[]) {
    console.log("episodesSeenOn setter", val); // Never called
    for (const seen of val) {
      const compareValIdx = this.#episodesSeenOn.findIndex(
        (e) => e.id === seen.id,
      );
      const compareVal =
        compareValIdx === -1 ? undefined : this.#episodesSeenOn[compareValIdx];
      if (compareVal) {
        this.#episodesSeenOn[compareValIdx] = seen;
      } else {
        this.#episodesSeenOn.push(seen);
      }
      // const check = setSaveEpisodeVal(seen, compareVal);
      // if(check.val) {
      //   newSeens.push(check.val);
      // }
    }
  }

  #thubmnailImageId = $state("");
  get thubmnailImageId() {
    return this.#thubmnailImageId;
  }
  set thubmnailImageId(val: string) {
    const check = setSaveStringVal(val, this.#thubmnailImageId);
    this.#thubmnailImageId = check.val;
    if (check.hasChanges) {
      this.saveSelf();
    }
  }

  // JS Only fields
  seenView = $derived.by(() => {
    if (this.episodesSeenOn.length > 0) {
      const seenar = this.episodesSeenOn.map((e) => e.episodesSeen.toString());
      return seenar.join("/");
    } else {
      return "0";
    }
  });
  disableSaving = $state(false);

  public static newDefault(disableSaving = false): SignalListItem {
    const newListItem = new SignalListItem(new list.ListItem());
    newListItem.disableSaving = disableSaving;
    newListItem.id = newDBID();
    const primaryTitle = SignalListItemTitle.newDefault(newListItem.id, true);
    console.log("New default item", "title itemId", primaryTitle.itemId, "item ID", newListItem.id);
    newListItem.titles.push(primaryTitle);
    newListItem.title = newListItem.titles[0];
    newListItem.broadcastType = "ONA";
    const newEpSeen = SignalEpisodeSeen.newDefault(newListItem.id);
    newListItem.episodesSeenOn.push(newEpSeen);
    newListItem.ongoing = true;
    newListItem.seasonNum = 1;
    newListItem.type = "base";
    newListItem.thubmnailImageId = "";

    return newListItem;
  }
  constructor(goItem: list.ListItem = new list.ListItem()) {
    // console.log("Signal list item", goItem.id);
    this.id = goItem.id;
    if (this.id === "") {
      NewDbID().then((newId) => (this.id = newId));
    }
    this.#titleId = goItem.titleId;
    if (!goItem.titles) {
      goItem.titles = [];
    }
    for (const currentTitle of goItem.titles) {
      if (currentTitle.primaryTitle === true) {
        this.#title = new SignalListItemTitle(currentTitle);
      }
      this.#titles.push(new SignalListItemTitle(currentTitle));
    }
    if (this.#titles.length === 0) {
      this.#title = new SignalListItemTitle(new list.ListItemTitle());
    }
    this.#type = goItem.type;
    this.#broadcastType = goItem.broadcastType;
    this.episodesTotal = goItem.episodesTotal;
    this.episodesSeen = goItem.episodesSeen;
    this.#ongoing = goItem.ongoing;
    this.#seasonNum = goItem.seasonNum;

    this.#parentItemId = goItem.parentItemId;
    if (Array.isArray(goItem.episodesSeenOn)) {
      // const newSeenOn: SignalEpisodeSeen[] = [];
      for (const seenOn of goItem.episodesSeenOn) {
        // newSeenOn.push(new SignalEpisodeSeen(seenOn))
        this.#episodesSeenOn.push(new SignalEpisodeSeen(seenOn));
      }
      // this.#episodesSeenOn = newSeenOn;
    }
    this.#thubmnailImageId = goItem.thubmnailImageId;

    // svelte effects
    $effect.root(() => {
      $effect(() => {
        // console.log("Changes", this.seasonNum);
      });
    });
  }

  public static fromID = async (itemId: string): Promise<SignalListItem | undefined> => {
    if(!itemId) {
      console.error("From id, no itemId provided", itemId);
    }
    const listItems = await GetListItemsByIDs([itemId]);
    if(listItems.length !== 1) {
      console.error("DB result is not 1", listItems.length);
      return undefined;
    }
    return new SignalListItem(listItems[0]);
  }
  saveSelf = async () => {
    if (!this.disableSaving || !this?.id) {
      const isSuccess = await SaveListItem(list.ListItem.createFrom(this));
      console.log("Saving success", isSuccess);
    } else {
      console.warn("Saving is disabled");
    }
  };
}

/*** LIST STORE ***/
export class ListStore {
  listIDs = $state<string[]>([]);
  currentPageItems = $state<SignalListItem[]>([]);
  currentPage = $state(0);
  perPage = $state(50);
  totalPages = $derived(Math.ceil(this.listIDs.length / this.perPage));
  searchResults = $state<SignalListItem[]>([]);
  searchPerformed = $state(false);
  searchCurrentPage = $state(0);
  searchTotalPages = $state(0);
  searchTotalResults = $derived(this.searchResults.length);

  constructor() {
    console.log("ListStore: init");
    // this.GetMainList();
    // Timeout to prevent stackoverflow on app load maybe?
    setTimeout(() => {
      this.getListIDs().then(() => {
        this.GetPage(1);
      });
    }, 1000);

    // this.onSync(); // Event listener from backend

    // svelte effects
    $effect.root(() => {
      $effect(() => {
        // console.log("Changes", this.list);
      });
    });
  }
  async getListIDs() {
    this.listIDs = await GetBaseItemIDs();
  }
  async GetPage(pageNum: number) {
    if (pageNum < 1 || pageNum > this.totalPages) {
      console.error(
        "GetPage: Invalid pagenum",
        "pageNum",
        pageNum,
        "totalPages",
        this.totalPages,
      );
    }
    const from = (pageNum - 1) * this.perPage;
    let to = from + this.perPage;
    if (to > this.listIDs.length) {
      to = this.listIDs.length - 1;
    }
    const getIDs: string[] = [];
    for (let i = from; i < to; i++) {
      getIDs.push(this.listIDs[i]);
    }
    const goListItems = await GetListItemsByIDs(getIDs);
    const newCurrentListItems: SignalListItem[] = [];
    for (const goItem of goListItems) {
      newCurrentListItems.push(new SignalListItem(goItem));
    }
    this.currentPageItems = newCurrentListItems;
    this.currentPage = pageNum;
  }

  async nextPage() {
    if (this.currentPage < this.totalPages) {
      await this.GetPage(this.currentPage + 1);
    }
  }
  async prevPage() {
    if (this.currentPage > 0) {
      await this.GetPage(this.currentPage - 1);
    }
  }

  findItemPage(itemId: string): { found: boolean; pageNum: number } {
    let currentPage = 1;
    let pageItem = 0;
    let found = false;
    for (const currentId of this.listIDs) {
      if (pageItem > this.perPage) {
        pageItem = 0;
        currentPage++;
      }
      if (currentId === itemId) {
        found = true;
        break;
      }
      pageItem++;
    }
    return {
      found,
      pageNum: currentPage,
    };
  }

  async refreshList(gotoItemPage: boolean = false, toPageItemID = "") {
    await this.getListIDs();
    let getPage = 1;
    if (gotoItemPage) {
      const itemPage = this.findItemPage(toPageItemID);
      if (itemPage.found) {
        getPage = itemPage.pageNum;
      }
    }
    await this.GetPage(getPage);
  }

  async search(query: string) {
    try {
      const options = new list.ListSearchOptions();
      options.orderDirection = "DESC";
      options.orderField = "title";
      options.searchFields = ["title"];
      const results = await Search(query, options);
      const listItems: SignalListItem[] = [];
      for (const item of results) {
        listItems.push(new SignalListItem(item));
      }
      this.searchResults = listItems;
      this.searchPerformed = true;
      this.searchTotalPages = 1;
      console.log("Search results", this.searchResults);
    } catch (e: any) {
      console.error("error while searching", e);
    }
  }
  clearSearch() {
    this.searchResults = [];
    this.searchPerformed = false;
    this.searchCurrentPage = 0;
    this.searchTotalPages = 0;
  }

  async onListItemTitlePush() {
    console.log("list item title push event listener init");
    EventsOn("list_item_title_push", (newTitle: list.ListItemTitle) => {
      if (!newTitle?.id || newTitle.id == "") {
        console.error("list_item_title_push: no ID specified");
        return;
      }
      console.log("Push event received", newTitle);
      const { val, titles } = objectFromGo<SignalListItemTitle>(newTitle);
      if (titles) {
        const title = titles.get(newTitle.id);
        if (title) {
          // this.setListItemTitle(title);
        }
      }
    });
  }
}

export const listStore = new ListStore();
