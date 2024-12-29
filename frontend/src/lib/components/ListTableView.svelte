<script lang="ts">
    import * as Table from "$components/ui/table";
    import { Checkbox } from "$components/ui/checkbox";
    import { Input } from "$components/ui/input";
    import { listStore } from "$stores";
    import ListTablePagination from "$components/ListTablePagination.svelte";
    import { SignalListItem } from "$stores/listStore.svelte";
    import { push as goTo } from "svelte-spa-router";
    import SerieTitle from "$components/SerieTitle.svelte";
    import SeenOnPopupEdit from "$components/SeenOnPopupEdit.svelte";
    import { cn } from "$lib/utils";
    const list = $derived.by(() =>
        listStore.searchPerformed
            ? listStore.searchResults
            : listStore.currentPageItems,
    );
    const currentPage = $derived.by(() =>
        listStore.searchPerformed
            ? listStore.searchCurrentPage
            : listStore.currentPage,
    );
    const totalPages = $derived.by(() =>
        listStore.searchPerformed
            ? listStore.searchTotalPages
            : listStore.totalPages,
    );
    let editSeason = $state(false);

    function openSeriePage(item: SignalListItem) {
        console.log("go to serie page", item.title.title);
        const serieURL = `#/serie/${item.id}`;
        goTo(serieURL);
    }
</script>

{#if listStore?.listIDs}
    <Table.Root>
        <Table.Caption>
            {#if listStore.searchPerformed}
                Haku tuloksia: {listStore.searchTotalResults}
            {:else}
                Listassa: {listStore.listIDs.length} - Sivuja: {listStore.totalPages}
            {/if}
        </Table.Caption>
        <!-- <Table.Caption>Listassa</Table.Caption> -->
        <Table.Header>
            <Table.Row>
                <Table.Head class="w-full text-center">Nimi</Table.Head>
                <!-- <Table.Head>ID</Table.Head> -->
                <Table.Head class="text-center min-w-10">Kausi</Table.Head>
                <Table.Head class="text-center">Jakso</Table.Head>
                <Table.Head class="text-center">Ongoing</Table.Head>
                <!-- <Table.Head>Method</Table.Head> -->
                <!-- <Table.Head class="text-right">Amount</Table.Head> -->
            </Table.Row>
        </Table.Header>
        <Table.Body>
            {#if totalPages > 0}
                {#each list as item, idx (item.id)}
                    <Table.Row ondblclick={() => openSeriePage(item)}>
                        <Table.Cell class="text-left ">
                            <SerieTitle
                                class="text-lg font-normal max-w-max inline line-clamp-1"
                                bind:item={list[idx]}
                            />
                        </Table.Cell>

                        <Table.Cell class={cn("p-3 h-full")}>
                            {#if item?.seasonNum || item?.seasonNum === 0}
                                <!-- {item.seasonNum} -->
                                <button
                                    class="relative w-full h-full"
                                    onclick={(e) => {
                                        //editSeason = !editSeason
                                        const target =
                                            e.target as HTMLButtonElement;
                                        const input = target.querySelector(
                                            "input",
                                        ) as HTMLInputElement;
                                        input.dataset.showInput = "true";
                                        input.focus();
                                        console.log(
                                            "input",
                                            input,
                                            "target",
                                            target,
                                        );
                                    }}
                                >
                                    <span class="pointer-events-none">
                                        {item.seasonNum}
                                    </span>
                                    <Input
                                        data-edit-season={idx}
                                        data-show-input="false"
                                        class="absolute text-center top-0 left-0 w-full h-full"
                                        pattern="[0-9]"
                                        onblur={(e) => {
                                            const target =
                                                e.target as HTMLInputElement;
                                            target.dataset.showInput = "false";
                                        }}
                                        bind:value={list[idx].seasonNum}
                                    />
                                </button>
                            {:else}
                                0
                            {/if}
                        </Table.Cell>
                        <Table.Cell>
                            <SeenOnPopupEdit bind:item={list[idx]} />
                            <!-- {#if item?.seenView}
                                {item?.seenView}
                            {:else}
                                0
                            {/if} -->
                        </Table.Cell>
                        <Table.Cell>
                            {#if item?.ongoing === true || item?.ongoing === false}
                                <Checkbox
                                    bind:checked={listStore.currentPageItems[
                                        idx
                                    ].ongoing}
                                />
                            {/if}
                        </Table.Cell>
                    </Table.Row>
                {/each}
            {/if}
            <!-- <Table.Row>
            <Table.Cell class="font-medium">INV001</Table.Cell>
            <Table.Cell>Paid</Table.Cell>
            <Table.Cell>Credit Card</Table.Cell>
            <Table.Cell class="text-right">$250.00</Table.Cell>
            </Table.Row> -->
        </Table.Body>
    </Table.Root>
    {#if totalPages > 1}
        <ListTablePagination />
    {/if}
{/if}

<style>
    :global(input[data-show-input="false"]) {
        display: none;
    }
</style>
