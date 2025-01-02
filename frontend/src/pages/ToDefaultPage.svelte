<script lang="ts">
  import { push as goTo } from "svelte-spa-router";
  $effect(() => {
    goTo("#/series")
  });
  </script>
