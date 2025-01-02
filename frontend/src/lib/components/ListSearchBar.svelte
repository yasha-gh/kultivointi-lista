<script lang="ts">
    import { Input } from "$lib/components/ui/input/index.js";
    import { listStore } from "$stores";
    import { Search } from "lucide-svelte";
    import { cn } from "$lib/utils";
    let { value = $bindable(), ...props } = $props();

    const debounceTime = 300;
    const minChars = 1;
    let timeout;
    async function search() {
        clearTimeout(timeout);
        timeout = setTimeout(() => {
            if (value.length >= minChars) {
                console.log("Change", value);
                listStore.search(value);
            }
            if (value.length === 0) {
                listStore.clearSearch();
            }
        }, debounceTime);
    }
</script>

<aside class={cn(`flex gap-2 relative p-2 ${props.class}`)}>
    <Search
        class="text-muted-foreground absolute left-6 top-[50%] h-4 w-4 translate-y-[-50%]"
    />
    <Input bind:value oninput={search} class="pl-10" placeholder="Hae listasta" />
</aside>
