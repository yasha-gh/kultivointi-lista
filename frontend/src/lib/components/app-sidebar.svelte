<script lang="ts">
import { Settings, Globe, MonitorPlay } from "lucide-svelte";
import * as Sidebar from "$lib/components/ui/sidebar/index.js";

// Menu items.
const items = [

	// {
	// 	title: "Devices",
	// 	url: "#/settings/devices",
	// 	icon: House,
	// },
	{
		title: "Lista",
		url: "#/series",
		icon: MonitorPlay,
	},
	{
		title: "Sivut",
		url: "#/",
		icon: Globe,
	},
	{
		title: "Asetukset",
		url: "#",
		icon: Settings,
	},
];
</script>

<Sidebar.Root>
	<Sidebar.Content>
		<Sidebar.Group>
			<Sidebar.GroupLabel>Kultivointi lista</Sidebar.GroupLabel>
			<Sidebar.GroupContent>
				<Sidebar.Menu>
					{#each items as item (item.title)}
						<Sidebar.MenuItem>
							<Sidebar.MenuButton>
								{#snippet child({ props })}
									<a href={item.url} {...props}>
										<item.icon />
										<span>{item.title}</span>
									</a>
								{/snippet}
							</Sidebar.MenuButton>
						</Sidebar.MenuItem>
					{/each}
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group>
	</Sidebar.Content>
</Sidebar.Root>
