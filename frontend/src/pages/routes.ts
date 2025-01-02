// import DevicesPage from "$pages/settings/DevicesPage.svelte";
import FrontPage from "$pages/FrontPage.svelte";
import SeriesPage from "$pages/SeriesPage.svelte";
import ViewEditSeriePage from "$pages/ViewEditSeriePage.svelte";
import ToDefaultPage from "$pages/ToDefaultPage.svelte";

export const routes = {
    // Exact path
    '/': ToDefaultPage,
    '/series': SeriesPage,
    '/serie/:itemId': ViewEditSeriePage,
    // '/settings/devices': DevicesPage,
    // '/user': UserPage,
    // Using named parameters, with last being optional
    // '/author/:first/:last?': Author,

    // Wildcard parameter
    // '/book/*': Book,

    // Catch-all
    // This is optional, but if present it must be the last
    // '*': NotFound,
}
