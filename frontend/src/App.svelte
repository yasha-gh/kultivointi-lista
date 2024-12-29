<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import AppSidebar from "$lib/components/app-sidebar.svelte";
  import { ModeWatcher } from "mode-watcher";
  import Router from 'svelte-spa-router';
  import { routes } from "$pages/routes";

  import "./app.pcss";
</script>

<ModeWatcher />
<Sidebar.Provider>
  <AppSidebar />
  <main>
    <Sidebar.Trigger />
    <Router {routes}/>
    <!-- {@render children?.()} -->
  </main>
</Sidebar.Provider>
<style>
main {
    display: flex;
    flex-direction: column;
    max-height: 100vh;
    width: 100%;
}
</style>
