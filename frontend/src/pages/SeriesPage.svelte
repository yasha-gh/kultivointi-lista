<script lang="ts">
    import OpenListFile from "$lib/components/OpenListFile.svelte";
    import Button from "$lib/components/ui/button/button.svelte";
    import Input from "$lib/components/ui/input/input.svelte";
    import { ScrollArea } from "$components/ui/scroll-area";
    import { listStore } from "$stores/listStore.svelte";
    import ListTableView from "$components/ListTableView.svelte";
    import ListSearchBar from "$components/ListSearchBar.svelte";
    async function getList() {
        // await listStore.GetMainList();
    }
    function removeLast() {
        const idx = listStore.listIDs.length - 1;
        listStore.listIDs.splice(idx);
    }
    let searchInputVal = $state("");
</script>

<h1 class="text-2xl font-semibold">Kultivointi lista</h1>
<ListSearchBar bind:value={searchInputVal} />
<ScrollArea class="p-4">
    <section class="main-content">
        <!-- <Button onclick={getList}>Get list</Button>
        <Button onclick={removeLast}>Remove last</Button>
        <Button onclick={() => listStore.GetPage(2)}>GetPage</Button> -->
        {#if listStore.totalPages > 0}
            <ListTableView />
        {:else}
            <OpenListFile />
        {/if}
    </section>
</ScrollArea>

<style>
    .main-content {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }
</style>
