<script lang="ts">
    import * as Pagination from "$lib/components/ui/pagination/index.js";
    import { ChevronRight, ChevronLeft } from "lucide-svelte";
    import { listStore } from "$stores";
    const count = $derived(listStore.listIDs.length);
    $inspect(listStore.currentPage);
</script>

{#if !listStore.searchPerformed}
    <Pagination.Root
        bind:page={listStore.currentPage}
        {count}
        perPage={listStore.perPage}
        siblingCount={3}
    >
        {#snippet children({ pages, currentPage })}
            <Pagination.Content>
                <Pagination.Item>
                    <Pagination.PrevButton onclick={() => listStore.prevPage()}>
                        <ChevronLeft class="size-4" />
                        <span class="hidden sm:block">Edellinen</span>
                    </Pagination.PrevButton>
                </Pagination.Item>
                {#each pages as page (page.key)}
                    {#if page.type === "ellipsis"}
                        <Pagination.Item>
                            <Pagination.Ellipsis />
                        </Pagination.Item>
                    {:else}
                        <Pagination.Item
                            isVisible={currentPage === listStore.currentPage}
                        >
                            <Pagination.Link
                                onclick={() => listStore.GetPage(page.value)}
                                {page}
                                isActive={currentPage === page.value}
                            >
                                {page.value}
                            </Pagination.Link>
                        </Pagination.Item>
                    {/if}
                {/each}
                <Pagination.Item>
                    <Pagination.NextButton onclick={() => listStore.nextPage()}>
                        <span class="hidden sm:block">Seuraava</span>
                        <ChevronRight class="size-4" />
                    </Pagination.NextButton>
                </Pagination.Item>
            </Pagination.Content>
        {/snippet}
    </Pagination.Root>
{/if}
