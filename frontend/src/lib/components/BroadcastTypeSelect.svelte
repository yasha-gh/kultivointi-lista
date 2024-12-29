<script lang="ts">
    import { broadcastTypes } from "$stores";
    import { type LabelValue } from "$lib/utils/index";
    import * as Select from "$components/ui/select";
    let value = $state("");
    let {
      id,
      // value = $bindable(),
      broadcastType = "",
      onselect = (val: LabelValue) => {},
      ...props
    } = $props();
    const triggerContent = $derived(
        broadcastTypes.find((f) => f.value === value)?.label ??
            "Valitse tyyppi",
    );
    function valueChanged(e: string) {
      broadcastType = value;
      console.log("value changed", e);
      onselect(broadcastTypes.find((b) => b.value === value));
    }
    $effect(() => {
      if(broadcastType) {
        value = broadcastType;
      }
    });
    $inspect(value);
</script>

<Select.Root type="single" name="favoriteFruit" onValueChange={valueChanged} bind:value={value}>
    <Select.Trigger id={id} class="">
        {triggerContent}
    </Select.Trigger>
    <Select.Content>
        <Select.Group>
            <Select.GroupHeading>Tyyppi</Select.GroupHeading>
            {#each broadcastTypes as broadcastType}
                <Select.Item
                    value={broadcastType.value}
                    label={broadcastType.label}
                >
                    {broadcastType.label}
                </Select.Item>
            {/each}
        </Select.Group>
    </Select.Content>
</Select.Root>
