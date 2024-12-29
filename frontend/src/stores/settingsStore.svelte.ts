import { main } from "$wails/go/models";
import { Sync as SyncSettings, GetSelf } from "$wails/go/main/Settings";
import { EventsOn } from "$wails/runtime";

export class SettingsStore {
  #dbFileDir = $state("")
  #dbFilename = $state("")
  $settings: main.Settings // Non directly reactive underlying object shared with backend

  // ** Local fields  ** //

  // this.dbFileDir
  get dbFileDir(): string {
    if(!this?.$settings?.dbFileDir) {
      return this.#dbFileDir;
    };
    return this.$settings.dbFileDir
  }
  set dbFileDir(value: string) {
    this.sync("dbFileDir", value);
    this.#dbFileDir = value;
    this.$settings.dbFileDir = value;
  }

  // this.dbFileDir
  get dbFilename(): string {
    if(!this?.$settings?.dbFilename) {
      return this.#dbFilename;
    };
    return this.$settings.dbFilename;
  }
  set dbFilename(value: string) {
    this.sync("dbFilename", value);
    this.#dbFilename = value;
    this.$settings.dbFilename = value;
  }

  constructor() {
    console.log("SettingsStore: init");
    GetSelf().then((currentSettings) => {
      console.log("Get self on store init", currentSettings);
      this.$settings = currentSettings;
      this.dbFileDir = currentSettings.dbFileDir;
      this.dbFilename = currentSettings.dbFilename;
    });

    this.onSync(); // Event listener from backend

    // svelte effects
    $effect.root(() => {
      $effect(() => {
        console.log("Changes", this.dbFilename);

        // // From SettingsStore to Backend sync
        // this.sync("dbFileDir", this.dbFileDir);
        // this.sync("dbFilename", this.dbFilename);
      });
    });
  }

  // 2 Way sync event with backend, must be filtered on both ends
  // $settings is not reactive so it's fine to assign as is
  async onSync() {
    console.log("onSync init");
    EventsOn("settings_sync", (newSettings: main.Settings) => {
      console.log("Sync event received", newSettings);
      this.$settings = main.Settings.createFrom(newSettings);

      // Update local fields if not same value
      if (this.dbFilename != newSettings.dbFilename) {
        this.dbFilename = newSettings.dbFilename;
      }
      if (this.dbFileDir != newSettings.dbFileDir) {
        this.dbFileDir = newSettings.dbFileDir;
      }
    })
  }

  // Sync SettingsStore value to backend (Syncs whole object if value changed)
  sync(field: string, value: any) {
    if(!this.$settings) {
      console.log("Cannot sync, $settings not initialized", value);
      return
    };

    // Check / set fields if valid
    switch(field) {
      case "dbFilename":
        if(this.$settings.dbFilename == value || value == "") return;
        this.$settings.dbFilename = value;
        break;
      case "dbFileDir":
        if(this.$settings.dbFileDir == value || value == "") return;
        this.$settings.dbFileDir = value;
        break;
      default:
        console.error("SettingsStore: trying to sync unknown field", field);
        return
    }

    // Sync changes to backend
    try {
      SyncSettings(this.$settings);
    } catch(e: any) {
      console.error("Failed to sync settings", e);
    }
  }

}

export const settingsStore = new SettingsStore();
