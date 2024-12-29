<script lang="ts">

import Button from "$lib/components/ui/button/button.svelte";
import { Input } from "$components/ui/input";
import { Save } from "$wails/go/main/Settings";
import { settingsStore as settings } from "$stores";
import { ScrollArea } from "$components/ui/scroll-area"
async function test(e: any) {
  await Save();
}


</script>

<h1>Etusivu</h1>
<ScrollArea class="p-4">
    <Input bind:value={settings.dbFilename} />
    <Button onclick={test}>Test</Button>
    <div>DBFile: {settings.dbFileDir}/{settings.dbFilename}</div>
</ScrollArea>
