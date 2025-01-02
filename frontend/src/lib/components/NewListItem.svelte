<script lang="ts">
    import { SignalListItem, SignalListItemTitle, listStore } from "$stores";
    import { list } from "$wails/go/models";
    import { NewDbID } from "$wails/go/list/List";
    import * as Dialog from "$lib/components/ui/dialog";
    import { Button } from "$components/ui/button";
    import { ScrollArea } from "$components/ui/scroll-area";
    import Serie from "$components/Serie.svelte";
    import SerieTitle from "./SerieTitle.svelte";
    import TitlesEdit from "$components/TitlesEdit.svelte";
    import { Save } from "lucide-svelte";
    let {
        item = $bindable<SignalListItem>() as SignalListItem,
        // SignalListItem.newDefault(true)
        showBorder = true,
        open = $bindable(false),
        ...props
    } = $props();
    let listItem = $state(item);

    async function saveItem() {
        console.log("save item", $state.snapshot(listItem));
        // for (const title of listItem.titles) {
        //     console.log("Title", title.itemId, listItem.id);
        // }
        // for (const seen of listItem.episodesSeenOn) {
        //     console.log("seen", seen.itemId, listItem.id);
        // }
        listItem.disableSaving = false;
        await listItem.saveSelf();
        await listStore.refreshList(true, listItem.id);
        open = false;
        listItem = SignalListItem.newDefault(true);
    }

    function init() {
        if (!listItem) {
            console.log("init: Creating new default list item");
            listItem = SignalListItem.newDefault(true);
        }
    }
    $effect(() => {
        init();
    });
</script>

<Dialog.Root bind:open>
    <Dialog.Content class="min-w-max max-h-[95vh]">
        <Dialog.Header>
            <Dialog.Title class="pb-4">Uusi sarja / leffa</Dialog.Title>
        </Dialog.Header>
        <ScrollArea class="new-item-dialog-body">
            <div class="px-8">
                <TitlesEdit bind:item={listItem} />
            </div>
            <Serie disableTitle={true} bind:item={listItem} />
        </ScrollArea>
        <Dialog.Footer class="items-center pr-8">
            <Button
                onclick={() => {
                    open = false;
                }}
                variant="outline">Peruuta</Button
            >
            <Button onclick={saveItem}><Save /> Tallenna</Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

<style>
    :global(.new-item-dialog-body) {
        --_dialog-header-height: 60px;
        --_dialog-footer-height: 120px;
        max-height: calc(
            95vh - var(--_dialog-footer-height) - var(--_dialog-header-height)
        );
    }
</style>
