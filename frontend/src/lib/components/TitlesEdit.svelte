<script lang="ts">
    import {
        SignalListItem,
        SignalListItemTitle,
        serieTitleLanguages,
    } from "$stores";
    import { getLangLabel } from "$lib/utils/index";
    import { NewDbID } from "$wails/go/list/List";
    import { list } from "$wails/go/models";
    import { Input } from "$components/ui/input";
    import { Button } from "$components/ui/button";
    import * as Select from "$components/ui/select";
    import * as Card from "$components/ui/card";
    import * as Table from "$components/ui/table";
    import { Trash2, CirclePlus } from "lucide-svelte";
    import { cn } from "$lib/utils";
    let {
        item = $bindable<SignalListItem>() as SignalListItem,
        // new SignalListItem(new list.ListItem()),
        showBorder = true,
        class: className = "",
        ...props
    } = $props();
    let titles = $derived(item.titles);

    async function deleteTitle(titleId: string) {
        if (!titleId) {
            return;
        }
        if (item.title.id === titleId) {
            console.error("cannot delete primary title");
            return;
        }
        for (let i = 0; i < item.titles.length; i++) {
            const title = item.titles[i];
            if (title.id === titleId) {
                item.titles.splice(i, 1);
                try {
                  item.saveSelf();
                } catch (e: any) {
                    console.log("Failed to save title", e);
                }
            }
        }
    }

    function newDefaultTitle(): SignalListItemTitle {
        // const defaultTitle = new SignalListItemTitle(new list.ListItemTitle());
        // const newId = await NewDbID();
        // defaultTitle.id = newId;
        // defaultTitle.title = "";
        // defaultTitle.primaryTitle = (titles.length < 1) ? true : false;
        // defaultTitle.itemId = item.id;
        // defaultTitle.lang = "zh_romanji";
        // return defaultTitle;

        return SignalListItemTitle.newDefault(item.id, titles.length > 1 ? false : true);
    }
    // async function init() {
    //     setTimeout(() => {
    //         if (item.titles.length < 1) {
    //           item.titles = [newDefaultTitle()];
    //           // item.disableSaving = true;
    //             // newDefaultTitle().then((newTitle) => {
    //             //     item.titles.push(newTitle);
    //             // });
    //         }
    //     }, 100);
    // }
    async function addNew() {
        console.log("Add new title");
        const newTitle = SignalListItemTitle.newDefault(item.id, titles.length > 1 ? false : true);
        titles.push(newTitle);
        // titles.push(await newDefaultTitle());
    }
    $effect(() => {
        // init();
    });
</script>

<!-- {#each item.titles as title, idx}
    <aside class="flex gap-2">
        <Input
            class="grow"
            bind:value={item.titles[idx].title}
            oninput={() => {
                if (item.titles[idx].primaryTitle) {
                    item.title = item.titles[idx];
                }
                item.saveSelf();
            }}
        />
        <Select.Root
            onValueChange={() => item.saveSelf()}
            bind:value={item.titles[idx].lang}
            type="single"
        >
            <Select.Trigger class="w-[15em]">
                {getLangLabel(item.titles[idx].lang)}
            </Select.Trigger>
            <Select.Content>
                {#each serieTitleLanguages as titleLang}
                    <Select.Item value={titleLang.value}
                        >{titleLang.label}</Select.Item
                    >
                {/each}
            </Select.Content>
        </Select.Root>
        <Button variant="destructive" onclick={() => deleteTitle(title.id)}
            ><Trash2 /></Button
        >
    </aside>
{/each} -->
<Card.Root class={cn(`${!showBorder && "border-0"}`)}>
    <Card.Header class="relative">
        <Button
            variant="outline"
            onclick={addNew}
            class="absolute bottom-0 right-5"
        >
            <CirclePlus /></Button
        >
        <Card.Title>Nimet</Card.Title>
        <!-- <Card.Description></Card.Description> -->
    </Card.Header>
    <Card.Content>
        <Table.Root>
            <Table.Body>
                {#if titles.length > 0}
                    {#each titles as title, idx}
                        <Table.Row>
                            <Table.Cell class="font-medium text-center">
                                <Input
                                    class="grow"
                                    placeholder="Sarjan nimi"
                                    bind:value={item.titles[idx].title}
                                    oninput={() => {
                                        if (item.titles[idx].primaryTitle) {
                                            item.title = item.titles[idx];
                                        }
                                        item.saveSelf();
                                    }}
                                />
                            </Table.Cell>
                            <Table.Cell class="">
                                <Select.Root
                                    onValueChange={() => item.saveSelf()}
                                    bind:value={item.titles[idx].lang}
                                    type="single"
                                >
                                    <Select.Trigger class="w-[15em]">
                                        {getLangLabel(item.titles[idx].lang)}
                                    </Select.Trigger>
                                    <Select.Content>
                                        {#each serieTitleLanguages as titleLang}
                                            <Select.Item value={titleLang.value}
                                                >{titleLang.label}</Select.Item
                                            >
                                        {/each}
                                    </Select.Content>
                                </Select.Root>
                            </Table.Cell>
                            <Table.Cell class="w-0 text-right">
                                <Button
                                    variant="destructive"
                                    onclick={() => deleteTitle(title.id)}
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
</Card.Root>
<!-- <Card.Footer>
        <div class="flex justify-end w-full">
            <Button>Lisää </Button>
            </div>
        </Card.Footer> -->
