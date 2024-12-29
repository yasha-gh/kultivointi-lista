<script lang="ts">
	import { LoadListFile } from "$wails/go/list/ListParser"
    import Button from "$components/ui/button/button.svelte";
	async function openFileDialog() {
		try {
			await LoadListFile();
		} catch(e: any) {
			console.error("Failed to load series list file", e);
		}
	}
</script>

<Button onclick={openFileDialog}>Avaa lista tiedosto</Button>
