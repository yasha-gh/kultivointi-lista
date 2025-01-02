<script lang="ts">
    import OpenListFile from "$lib/components/OpenListFile.svelte";
    import Button from "$lib/components/ui/button/button.svelte";
    import Input from "$lib/components/ui/input/input.svelte";
    import { CirclePlus } from "lucide-svelte";
    import { ScrollArea } from "$components/ui/scroll-area";
    import { listStore } from "$stores/listStore.svelte";
    import ListTableView from "$components/ListTableView.svelte";
    import ListSearchBar from "$components/ListSearchBar.svelte";
    import NewListItem from "$lib/components/NewListItem.svelte";

    let searchInputVal = $state("");
    let newDialogOpen = $state(false);
</script>

<h1 class="text-2xl font-semibold">Kultivointi lista</h1>
<section class="page-head flex items-center p-4 gap-4">
    <ListSearchBar class="grow" bind:value={searchInputVal} />
    <Button onclick={() => newDialogOpen = !newDialogOpen}>Uusi <CirclePlus /></Button>
</section>
<NewListItem bind:open={newDialogOpen} />
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
