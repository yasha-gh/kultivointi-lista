<script lang="ts">
    import { LoadListFile } from "$wails/go/list/ListParser";
    import Button from "$components/ui/button/button.svelte";
    let importStatus = "idle";
    let importError = "";
    async function openFileDialog() {
        try {
            importStatus = "importing";
            await LoadListFile();
            location.reload();
        } catch (e: any) {
            console.error("Failed to load series list file", e);
            importStatus = "error";
            importError = "Listan tuonti ei onnistunut";
        }
    }
</script>

{#if importStatus == "idle"}
    <Button onclick={openFileDialog}>Avaa lista tiedosto</Button>
{/if}
{#if importStatus == "error"}
    <div class="flex items-center justify-center">
        <span class="text-red-500">{importError} </span>
    </div>
{/if}
{#if importStatus == "importing"}
    <div class="flex items-center justify-center">
        Listaa ladataan tiedostosta
    </div>
{/if}
