<script lang="ts">
    import { ScrollArea } from "$components/ui/scroll-area";
    import { Button } from "$components/ui/button";
    import * as Card from "$components/ui/card";
    import { getPathLast } from "$lib/utils/index";
    import { SignalListItem } from "$stores/listStore.svelte";
    import Serie from "$components/Serie.svelte";
    import { DeleteListItem } from "$wails/go/list/List";
    import { Trash2 } from "lucide-svelte";
    import { listStore } from "$stores";
    let itemId = $state(getPathLast());
    let item = $state<SignalListItem | undefined>() as SignalListItem | undefined;

    async function deleteListItem() {
      if(item) {
        const success = await DeleteListItem(item.id)
        if(success) {
          item = undefined;
          await listStore.refreshList()
        } else {
          console.error("Failed to delete list item")
        }
      }
    }

    function init() {
        if (!item && itemId) {
            SignalListItem.fromID(itemId).then((itemFromID) => {
                if (itemFromID) {
                    item = itemFromID;
                }
            });
        }
    }
    $effect(() => {
        init();
    });
</script>

<!-- <h1>
{#if item?.title?.title}
    {item.title.title}
    {:else}
    Uusi sarja
    {/if}
</h1> -->
<ScrollArea>
    {#if item}
        <Serie bind:item />
        <div class="p-8 pt-0">
        <Card.Root>
            <Card.Content class="flex gap-2 justify-end p-4">
                <Button onclick={deleteListItem} variant="destructive"><Trash2 /> Poista sarja</Button>
            </Card.Content>
        </Card.Root>
        </div>
    {:else}
        <div>Sarjaa ei ole valittu</div>
    {/if}
</ScrollArea>
