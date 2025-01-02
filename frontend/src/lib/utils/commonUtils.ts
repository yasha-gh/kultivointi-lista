import { serieTitleLanguages, SignalListItem } from "$stores/listStore.svelte";
import { GetListItemsByIDs } from "$wails/go/list/List";
import { list } from "$wails/go/models";
import { customAlphabet } from 'nanoid'
import { listStore } from "$stores";

export function getPathLast(): string {
  const path = `${window.location}`.split("#");
  let id = "";
  if (path.length > 1) {
    const pathParts = path[path.length - 1].split("/");
    id = pathParts[pathParts.length - 1];
  }
  return id;
}

export function getLangLabel(value: any): string {
  for (const l of serieTitleLanguages) {
    if (l.value === value) {
      return l.label;
    }
  }

  console.log("kieli", value);
  return "Kieltä ei löytynyt";
}

export function newDBID(): string {
  const customNanoID = customAlphabet("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 15)
  return customNanoID();
}
