<script lang="ts">
    import { Check, ChevronsUpDown, CirclePlus } from "lucide-svelte";
    import { tick } from "svelte";
    import * as Command from "$components/ui/command";
    import * as Popover from "$components/ui/popover";
    import * as AlertDialog from "$components/ui/alert-dialog";
    import { Input } from "$components/ui/input";
    import { Button } from "$components/ui/button";
    import { cn } from "$lib/utils.js";
    import { siteStore, SignalSite, type SiteSelectOption } from "$stores/siteStore.svelte";
    import { list } from "$wails/go/models";
    let {
            siteId = $bindable(""),
            selectchange = (e: SiteSelectOption) => {  },
            ...props
        } = $props();
    let open = $state(false);
    let selectValue = $state("NOSELECT");
    let triggerRef = $state<HTMLButtonElement>(null!);
    let addSiteDialogOpen = $state(false);
    let newSiteURL = $state("");


    // let selectList = $derived.by(() => {
    //   const selectOptions: SelectOption[] = [];
    //   for(const site of siteStore.sites) {
    //     selectOptions.push({
    //       label: site.getName(),
    //       siteId: site.id
    //     })
    //   }
    //   selectOptions.push({
    //     label: "Lisää uusi",
    //     siteId: "NEW"
    //   });
    //   return selectOptions;
    // });
    // const selectedValue = $derived(
    //     selectList.find((o) => o.siteId === value)
    // );
    const selectedValue = $derived(siteStore.getOptionValue(selectValue)?.label);
    // We want to refocus the trigger button when the user selects
    // an item from the list so users can continue navigating the
    // rest of the form with the keyboard.
    function closeAndFocusTrigger() {
        open = false;
        tick().then(() => {
            triggerRef.focus();
        });
    }

    async function selectChanged() {
      selectchange(siteStore.getOptionValue(selectValue));
      // selectValue = siteId;
      // selectchange(new CustomEvent("selectchange"))
      // selectchange(siteStore.getOptionValue(selectValue));
      // const storeSite = siteStore.getSiteByID(selectedValue)
      // if(storeSite) {
      //   await storeSite.saveSelf();
      //   if(storeSite.id != siteId) {
      //     siteId = storeSite.id;
      //   }
      // }
    }

    async function addSiteAccepted() {
        console.log("Accepted with value", newSiteURL);
        addSiteDialogOpen = false;
        const newSite = new SignalSite(new list.Site());
        await newSite.fromUrl(newSiteURL);
        const validRes = newSite.validate();
        if (validRes.valid) {
            await newSite.saveSelf();
            siteStore.sites.push(newSite);
            selectValue = newSite.id;
            selectChanged();
        } else {
            console.error("Not valid site", validRes);
        }
    }

    $effect(() => {
      if(siteId != "") {
        selectValue = siteId;
      }
      // if(!site?.domainBase) {
      //   site = siteStore.getOptionValue("NOSELECT") as SiteSelectOption;
      // }
    });
</script>

{#if siteStore.sites.length > 0}
    <aside class="inline-flex gap-1 w-max">
        <Popover.Root bind:open>
            <Popover.Trigger bind:ref={triggerRef}>
                {#snippet child({ props })}
                    <Button
                        variant="outline"
                        class="w-[200px] justify-between"
                        {...props}
                        role="combobox"
                        aria-expanded={open}
                    >
                        {selectedValue || "Valitse sivusto..."}
                        <ChevronsUpDown class="opacity-50" />
                    </Button>
                {/snippet}
            </Popover.Trigger>
            <Popover.Content class="w-[200px] p-0">
                <Command.Root>
                    <!-- <Command.Input oninput={(e) => { const i = e.target as HTMLInputElement; newSiteURL = i.value; }} placeholder="Etsi sivusto..." /> -->
                    <Command.List>
                        <Command.Empty>
                            <Button
                                onclick={() => {
                                    addSiteDialogOpen = true;
                                }}>Lisää <CirclePlus /></Button
                            >
                        </Command.Empty>
                        <Command.Group>
                            {#each siteStore.selectOptions as option, idx}
                                <Command.Item
                                    value={option.siteId}
                                    onSelect={() => {
                                        selectValue = option.siteId;
                                        closeAndFocusTrigger();
                                        selectChanged();
                                    }}
                                >
                                    <Check
                                        class={cn(
                                            selectValue !== option.siteId &&
                                                "text-transparent",
                                        )}
                                    />
                                    {option.label}
                                    <!-- {site.domainBase}.{site.domainTopLevel} -->
                                </Command.Item>
                            {/each}
                        </Command.Group>
                    </Command.List>
                </Command.Root>
            </Popover.Content>
        </Popover.Root>
        <Button variant="outline" onclick={() => addSiteDialogOpen=true}><CirclePlus /></Button>
        </aside>
{:else}
    <aside class="flex items-center justify-center font-semibold gap-4">
        <span class="text-lg">Ei sivua</span>
        <Button
            onclick={() => {
                addSiteDialogOpen = true;
            }}>Lisää <CirclePlus /></Button
        >
        <!-- <Button onclick={test}>Test</Button> -->
    </aside>
{/if}

<!-- Add dialog -->
<AlertDialog.Root bind:open={addSiteDialogOpen}>
    <AlertDialog.Content>
        <AlertDialog.Header>
            <AlertDialog.Title>Lisää sivu</AlertDialog.Title>
            <AlertDialog.Description>
                <Input
                    bind:value={newSiteURL}
                    onkeyup={(e) => {
                        if (e?.key === "Enter") addSiteAccepted();
                    }}
                    data-site-add
                    placeholder="https://..."
                />
            </AlertDialog.Description>
        </AlertDialog.Header>
        <AlertDialog.Footer>
            <AlertDialog.Cancel
                onclick={() => {
                    newSiteURL = "";
                }}>Peruuta</AlertDialog.Cancel
            >
            <AlertDialog.Action onclick={addSiteAccepted}>Ok</AlertDialog.Action
            >
        </AlertDialog.Footer>
    </AlertDialog.Content>
</AlertDialog.Root>
