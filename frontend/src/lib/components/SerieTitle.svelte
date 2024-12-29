<script lang="ts">
    import {
        SignalListItem,
        SignalListItemTitle,
        listStore,
        serieTitleLanguages,
    } from "$stores/listStore.svelte";
    import { NewDbID } from "$wails/go/list/List";
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import { Input } from "$components/ui/input";
    import * as Select from "$components/ui/select";
    import * as Card from "$components/ui/card";
    import { Button } from "$components/ui/button";
    import { getLangLabel } from "$lib/utils/index";
    import { list } from "$wails/go/models";
    import { Trash2, CirclePlus } from "lucide-svelte";
    import { cn } from "$lib/utils";
    let {
        item = $bindable<SignalListItem>(
            new SignalListItem(new list.ListItem()),
        ) as SignalListItem,
        editMode = $bindable(false),
        titleTag = "h2",
        class: className = "",
        ...props
    } = $props();
    let dialogOpen = $state(false);
    function startEdit() {
        dialogOpen = true;
        console.log("start edit", dialogOpen);
    }
    let newTitle = $state(newDefaultTitle());

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
        const defaultTitle = new SignalListItemTitle(new list.ListItemTitle());
        NewDbID().then((id) => {
            defaultTitle.id = id;
        });
        defaultTitle.title = "";
        defaultTitle.primaryTitle = false;
        defaultTitle.itemId = item.id;
        defaultTitle.lang = "zh_romanji";
        return defaultTitle;
    }

    async function addTitle() {
        try {
            newTitle.itemId = item.id;
            if (!newTitle.title) {
                console.error("Item title not set");
                return;
            }
            item.titles.push(newTitle);
            item.saveSelf();

            newTitle = newDefaultTitle();
        } catch (e: any) {
            console.error("Failed to add item title", e);
        }
    }

    // $effect(() => {
    //     if (newTitle.lang === "") newTitle.lang = "zh_romanji";
    // });
</script>

<svelte:element
    this={titleTag}
    class={cn(`serie-title text-primary text-4xl ${className}`)}
    onclick={startEdit}
    role="button"
    tabindex="0"
>
    {#if item.titles.length > 0}
        {#each item.titles as title, idx}
            {#if title.primaryTitle}
                <strong>{title.title}</strong>
            {:else}
                {title.title}
            {/if}
            {#if idx < item.titles.length - 1}
                &nbsp;/&nbsp;
            {/if}
        {/each}
    {/if}
</svelte:element>
<Dialog.Root bind:open={dialogOpen}>
    <Dialog.Content>
        <Dialog.Header>
            <Dialog.Title class="pb-4">Muokkaa nimiä</Dialog.Title>
            <Dialog.Description class="gap-2 flex flex-col">
                {#each item.titles as title, idx}
                    <aside class="flex gap-2">
                        <Input
                            class="grow"
                            bind:value={item.titles[idx].title}
                            oninput={() => {
                              if(item.titles[idx].primaryTitle) {
                                item.title = item.titles[idx];
                              }
                              item.saveSelf()}
                            }
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
                        <Button
                            variant="destructive"
                            onclick={() => deleteTitle(title.id)}
                            ><Trash2 /></Button
                        >
                    </aside>
                {/each}

                <h3
                    class="pt-4 text-sm text-primary font-semibold leading-none tracking-tight"
                >
                    Lisää uusi
                </h3>
                <aside class="flex gap-2 pt-4">
                    <Input
                        class="grow"
                        placeholder="Uusi sarjan nimi"
                        bind:value={newTitle.title}
                    />
                    <Select.Root bind:value={newTitle.lang} type="single">
                        <Select.Trigger class="w-[20em]"
                            >{getLangLabel(newTitle.lang)}</Select.Trigger
                        >
                        <Select.Content>
                            {#each serieTitleLanguages as titleLang}
                                <Select.Item value={titleLang.value}
                                    >{titleLang.label}</Select.Item
                                >
                            {/each}
                        </Select.Content>
                    </Select.Root>
                    <Button variant="default" onclick={addTitle}>
                        Lisää <CirclePlus />
                    </Button>
                </aside>
            </Dialog.Description>
        </Dialog.Header>
    </Dialog.Content>
</Dialog.Root>
