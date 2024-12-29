import './style.css'
import App from './App.svelte'
import { mount } from "svelte";
import { settingsStore } from "$stores/settingsStore.svelte";
import { listStore } from "$stores";
const app = mount(App, {
  target: document.getElementById('app') as HTMLElement
})

export default app
