<script lang="ts">
    import SerieSeenOn from "$components/SerieSeenOn.svelte";
    import { SignalListItem } from "$stores";
    import * as Dialog from "$components/ui/dialog";
    import { list } from "$wails/go/models";
    let {
        item = $bindable(
            new SignalListItem(new list.ListItem()),
        ) as SignalListItem,
        ...props
    } = $props();
    let dialogOpen = $state(false);
</script>

<button onclick={() => (dialogOpen = !dialogOpen)}> {item.seenView} </button>

<Dialog.Root  bind:open={dialogOpen}>
    <!-- <Dialog.Trigger>Open</Dialog.Trigger> -->
    <Dialog.Content class="min-w-max">
        <Dialog.Header>
            <Dialog.Title>Muokkaa sivuja ja jaksoja</Dialog.Title>
        </Dialog.Header>
        <SerieSeenOn showBorder={false} bind:item />
    </Dialog.Content>
</Dialog.Root>
