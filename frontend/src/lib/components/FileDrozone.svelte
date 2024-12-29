<script lang="ts">

function dropHandler(ev: any) {
	console.log("File(s) dropped", ev);

	// Prevent default behavior (Prevent file from being opened)
	ev.preventDefault();

	if (ev.dataTransfer.items) {
		// Use DataTransferItemList interface to access the file(s)
		[...ev.dataTransfer.items].forEach((item, i) => {
			// If dropped items aren't files, reject them
			if (item.kind === "file") {
				const file = item.getAsFile();
				console.log(`… file[${i}].name = ${file.name}`);
			}
		});
	} else {
		// Use DataTransfer interface to access the file(s)
		[...ev.dataTransfer.files].forEach((file, i) => {
			console.log(`… file[${i}].name = ${file.name}`);
		});
	}
}

function dragOverHandler(ev: any) {
	// console.log("File(s) in drop zone", ev);

	// Prevent default behavior (Prevent file from being opened)
	ev.preventDefault();
}
function onFileSelected(e: Event) {
	console.log("File input file selected", e.target);
	const target = e.target as HTMLInputElement 
	if(!target?.files || target.files.length < 1) {
	}
}
let fileInput: HTMLInputElement
function onClick() {
	fileInput.click();	
} 
</script>
<style>
.drop-zone {
	aspect-ratio: 16 / 8;	
	padding: 2rem;
	box-sizing: border-box;
	display: flex;	
	justify-content: stretch;
	div {
		display: flex;
		align-items: center;
		justify-content:center;
		border: 1px dashed white;
		border-radius: 1rem;
		width: 100%;
		font-size: 1.4em;
	}
}
.drop-input {
	/* visibility: hidden; */
	display: none;
}
</style>

<div
	class="drop-zone"
	ondrop={dropHandler}
	ondragover={dragOverHandler}
	onclick={onClick}
	onkeyup={onClick}
	role="button"
	tabindex="0"
>
	<input 
		class="drop-input"
		bind:this={fileInput}
		onchange={onFileSelected}
		type="file" 
	/>
	<div>Tiputa lista tänne tai klikkaa ja valitse</div>
</div>


