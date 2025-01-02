<script lang="ts">
    import { NewDbID, GetListItemsByIDs } from "$wails/go/list/List";
    import SerieTitle from "$components/SerieTitle.svelte";
    import SerieSeenOn from "$components/SerieSeenOn.svelte";
    import BroadcastTypeSelect from "$components/BroadcastTypeSelect.svelte";
    import * as Card from "$components/ui/card";
    import { Input } from "$components/ui/input";
    import { Label } from "$components/ui/label";
    import { list } from "$wails/go/models";
    import {
        broadcastTypes,
        listStore,
        SignalListItem,
    } from "$stores/listStore.svelte";
    import { type LabelValue } from "$lib/utils/index";
    // let item = $state(new SignalListItem(new list.ListItem()));
    let displayError = $state("");
    let {
        // itemId = $bindable(""),
        newItem = $bindable(false),
        disableTitle = false,
        item = $bindable(
            // new SignalListItem(new list.ListItem()),
        ) as SignalListItem,
        ...props
    } = $props();
    function broadcastTypeSelected(bt: LabelValue) {
        item.broadcastType = bt.value;
        // console.log("item bt", item.broadcastType);
        // console.log("bt", bt);
    }
    $effect(() => {
        // if (
        //     itemId.trim() == "" ||
        //     itemId.trim().toUpperCase() == "NEW" ||
        //     newItem == true
        // ) {
        //     newItem = true;
        //     NewDbID().then((newId) => {
        //         itemId = newId;
        //         item.disableSaving = true;
        //         item.id = newId;
        //     });
        // } else {
        //     // const ret = GetItemById(itemId, item);
        //     // item = ret.item;
        //     // if(ret.err) {
        //     //   displayError = ret.err;
        //     // }
        //     try {
        //         let listItem = listStore.currentPageItems.find(
        //             (ci) => ci.id === itemId,
        //         );
        //         if (listItem) {
        //             item = listItem;
        //             console.log("current", item);
        //         } else {
        //             GetListItemsByIDs([itemId]).then((items) => {
        //                 console.log(items);
        //                 if (items.length != 1 || !items[0]) {
        //                     displayError = "Virhe sarjaa hakiessa";
        //                 } else {
        //                     item = new SignalListItem(items[0]);

        //                     console.log("from DB", item);
        //                 }
        //             });
        //         }
        //     } catch (e: any) {
        //         console.error("Failed to get list items from DB", e);
        //     }
        // }
    });
    $inspect(item.broadcastType);
</script>

{#if displayError === ""}
    <section class="serie-component flex flex-col gap-4 p-8">
        {#if !disableTitle}
            <SerieTitle bind:item />
        {/if}
        <Card.Root>
            <Card.Header>
                <Card.Title>Perustiedot</Card.Title>
                <!-- <Card.Description>Card Description</Card.Description> -->
            </Card.Header>
            <Card.Content>
                <div class="grid grid-cols-2 gap-4">
                    <div class="flex flex-col space-y-1.5">
                        <Label for="seasonNum">Kausi numero</Label>
                        <Input
                            id="seasonNum"
                            class="text-center"
                            pattern="[0-9]"
                            bind:value={item.seasonNum}
                            placeholder="kausi numero"
                        />
                    </div>
                    <div class="flex flex-col space-y-1.5">
                        <Label for="broadcastType">Tyyppi</Label>
                        <BroadcastTypeSelect
                            id="broadcastType"
                            onselect={broadcastTypeSelected}
                            broadcastType={item.broadcastType}
                        />
                    </div>
                </div>
            </Card.Content>
            <!-- <Card.Footer>
        <p>Card Footer</p>
      </Card.Footer> -->
        </Card.Root>

        <SerieSeenOn bind:item />
        <!-- <h2>
    {#if item?.title?.title}
    {item.title.title}
    {:else}
    Uusi sarja
    {/if}
    </h2> -->
    </section>
{:else}
    <section class="serie-component serie-error">
        <span class="text-red-500">{displayError}</span>
    </section>
{/if}
