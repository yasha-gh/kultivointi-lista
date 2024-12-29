import { SignalListItemTitle } from "$stores"
import { list } from "$wails/go/models";

export type objectFromGoResults<T> = {
  val: T;
  titles?: Map<string, SignalListItemTitle>;
};

export function objectFromGo<T>(value: object): objectFromGoResults<T> {
  const o: any = {};
  let titles: Map<string, SignalListItemTitle> | undefined;
  for (const [k, v] of Object.entries(value)) {
    switch (k) {
      case "titles":
        if (Array.isArray(v)) {
          for (const title of v as list.ListItemTitle[]) {
            if (!title?.id || title.id == "") {
              console.error("No item id", title);
              continue;
            }
            if (!titles) {
              titles = new Map<string, SignalListItemTitle>();
            }
            titles.set(title.id, title);
          }
          o[k] = v;
        } else {
          o[k] = [];
        }
        break;
      default:
        o[k] = v;
        break;
    }
  }
  return {
    val: o as T,
    titles: titles,
  } as objectFromGoResults<T>;
}

export function setSaveNumVal(
  val: string | number,
  compareVal: number
): { val: number, hasChanges: boolean } {
  const ret = {
    val: 0,
    hasChanges: false
  }
  // let numVal = Number(val);
  ret.val = Number(val);
  if (Number.isNaN(ret.val)) {
    ret.val = 0;
  }
  if (compareVal !== ret.val) {
    ret.hasChanges = true;
  }
  return ret;
}

export function setSaveStringVal(
  val: string,
  compareVal: string
): { val: string, hasChanges: boolean } {
  const ret = {
    val: val,
    hasChanges: false
  }

  if (!val) {
    ret.val = "";
  }
  if (compareVal !== ret.val) {
    ret.hasChanges = true;
  }
  return ret;
}

export function setSaveBoolVal(
  val: string | number | boolean | undefined,
  compareVal: string | number | boolean | undefined
): { val: boolean, hasChanges: boolean } {
  console.log("set save bool", val)
  const ret = {
    val: val,
    hasChanges: false
  }
  if(val === true || val === false) {
    ret.val = val;
  }
  if(val === "true") {
    ret.val = true;
  }
  if(val === "false") {
    ret.val == false;
  }
  if(val === 1) {
    ret.val = true;
  }
  if(val === 0) {
    ret.val = false;
  }
  if (!val) {
    ret.val = false;
  }
  if (compareVal !== ret.val) {
    ret.hasChanges = true;
  }
  return ret;
}
