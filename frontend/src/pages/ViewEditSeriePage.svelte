<script lang="ts">
    import { getPathLast } from "$lib/utils/index";
    import { SignalListItem } from "$stores/listStore.svelte";
    import { list } from "$wails/go/models";
    import { listStore } from "$stores";
    import Serie from "$components/Serie.svelte";
    let itemId = $state(getPathLast());
    let item = $state(new SignalListItem(new list.ListItem()));

    </script>
<!-- <h1>
{#if item?.title?.title}
    {item.title.title}
    {:else}
    Uusi sarja
    {/if}
</h1> -->
<Serie bind:itemId={itemId} bind:item={item} />
