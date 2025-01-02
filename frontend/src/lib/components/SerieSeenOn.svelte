<script lang="ts">
    import {
        listStore,
        SignalEpisodeSeen,
        SignalListItem,
        siteStore,
        type SiteSelectOption,
    } from "$stores";
    import { list } from "$wails/go/models";
    import { NewDbID } from "$wails/go/list/List";
    import * as Card from "$components/ui/card";
    import * as Table from "$components/ui/table";
    import SiteSelect from "$components/SiteSelect.svelte";
    import { Input } from "$components/ui/input";
    import { Button } from "$components/ui/button";
    import { CirclePlus, Trash2 } from "lucide-svelte";
    let {
        item = $bindable(
            // new SignalListItem(new list.ListItem()),
        ) as SignalListItem,
        showBorder = true,
        ...props
    } = $props();

    function siteSelected(e: SiteSelectOption, episodeSeenIdx: number) {
        console.log("Site selected", e, episodeSeenIdx);

        if (e?.siteId && e.siteId !== "NOSELECT") {
            console.log("has value", item.episodesSeenOn[episodeSeenIdx]);
            const storeSite = siteStore.getSiteByID(e.siteId);

            item.episodesSeenOn[episodeSeenIdx].site = storeSite;

            if (!storeSite) {
                item.episodesSeenOn[episodeSeenIdx].siteId = "";
            } else {
                storeSite.episodeTemplate = "Testi";
                item.episodesSeenOn[episodeSeenIdx].siteId = storeSite.id;
            }
        } else {
            item.episodesSeenOn[episodeSeenIdx].site = undefined;
            item.episodesSeenOn[episodeSeenIdx].siteId = "";
        }
        item.saveSelf(); // Reactivity broken
    }
    function episodeChanged() {
        item.saveSelf();
    }
    async function addNew() {
        const newId = await NewDbID();
        const newEp = SignalEpisodeSeen.newDefault(item.id);
        // const newEp = new SignalEpisodeSeen(new list.EpisodeSeen());
        newEp.id = newId;
        // newEp.itemId = item.id;
        newEp.siteId = "";
        newEp.site = undefined;
        item.episodesSeenOn.push(newEp);
    }
    function deleteEp(epIdx: number) {
        item.episodesSeenOn.splice(epIdx, 1);
        episodeChanged();
    }
</script>
<style>
    .episode-col {

        display: grid;
        grid-template-columns: 4rem auto;
        align-items: center;
        min-width: max-content;
    }
    </style>

<Card.Root class={`${!showBorder && "border-0"}`}>
    <Card.Header class="relative">
        <Button
            variant="outline"
            onclick={addNew}
            class="absolute bottom-0 right-5"
        >
            <CirclePlus /></Button
        >
        <Card.Title>Nähdyt jaksot</Card.Title>
        <Card.Description>Sivustot ja jaksot mitkä on nähty</Card.Description>
    </Card.Header>
    <Card.Content>
        <Table.Root>
            <Table.Body>
                {#if item?.episodesSeenOn?.length > 0}
                    {#each item.episodesSeenOn as episodeSeen, idx}
                        <Table.Row>
                            <Table.Cell class="font-medium text-center">
                                <SiteSelect
                                    selectchange={(e) => siteSelected(e, idx)}
                                    bind:siteId={episodeSeen.siteId}
                                />
                            </Table.Cell>
                            <Table.Cell class="">
                                <div class="episode-col">
                                <span>Jakso</span> <Input
                                    class="max-w-14 text-center"
                                    oninput={episodeChanged}
                                    bind:value={episodeSeen.episodesSeen}
                                />
                                </div>
                            </Table.Cell>
                            <Table.Cell class="w-0 text-right">
                                <Button
                                    variant="destructive"
                                    onclick={() => deleteEp(idx)}
                                >
                                    <Trash2 />
                                </Button>
                            </Table.Cell>
                        </Table.Row>
                    {/each}
                {/if}
            </Table.Body>
        </Table.Root>
    </Card.Content>
    <!-- <Card.Footer>
        <div class="flex justify-end w-full">
            <Button>Lisää </Button>
            </div>
        </Card.Footer> -->
</Card.Root>
