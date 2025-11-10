<script>
  import { watchlistState, watchlistAPI } from "../stores/watchlist.svelte.js";
  import { servicesAPI } from "../stores/services.svelte.js";
  import { onMount } from "svelte";

  let openMenus = $state({});

  onMount(() => {
    watchlistAPI.fetch();
  });

  function toggleMenu(serviceName) {
    openMenus[serviceName] = !openMenus[serviceName];
  }

  function closeMenu(serviceName) {
    openMenus[serviceName] = false;
  }

  async function startService(serviceName) {
    const success = await servicesAPI.start(serviceName);
    if (success) {
      await watchlistAPI.fetch(); // Refresh the list
    }
    closeMenu(serviceName);
  }

  async function stopService(serviceName) {
    const success = await servicesAPI.stop(serviceName);
    if (success) {
      await watchlistAPI.fetch(); // Refresh the list
    }
    closeMenu(serviceName);
  }

  async function restartService(serviceName) {
    const success = await servicesAPI.restart(serviceName);
    if (success) {
      await watchlistAPI.fetch(); // Refresh the list
    }
    closeMenu(serviceName);
  }

  async function toggleAutoRestart(serviceName, currentAutoRestart) {
    await watchlistAPI.update(serviceName, !currentAutoRestart);
    await watchlistAPI.fetch(); // Refresh the list
    closeMenu(serviceName);
  }
</script>

<section class="grid-12">
  <div id="watch" class="card col-span-12 lg:col-span-4">
    <div class="flex items-center justify-between mb-2">
      <h2 class="font-semibold">Services</h2>
    </div>
    <!-- Watchlist items will be dynamically added here -->
    <ul id="watchUl" class="space-y-2">
      {#each watchlistState.items as item}
        <li
          class="p-2 bg-neutral-950 rounded-md border border-neutral-800 flex items-center justify-between"
        >
          <div class="flex flex-col">
            <div class="font-medium">
              {#if item.service?.state === "running"}
                <span class="text-green-400">●</span>
              {:else if item.service?.state === "stopped"}
                <span class="text-red-400">●</span>
              {:else}
                <span class="text-gray-400">●</span>
              {/if}
              {item.serviceName}
            </div>
            <!-- Service info -->
            <div class="text-xs text-neutral-400">
              <span
                >Auto-restart {item.autoRestart ? "enabled" : "disabled"}</span
              >
            </div>
          </div>
          <div class="relative">
            <button
              class="inline-block align-middle px-2 py-0.5 border border-neutral-700 rounded-md hover:bg-neutral-800"
              >+</button
            >
            <button
              onclick={() => toggleMenu(item.serviceName)}
              class="inline-block align-middle px-2 py-0.5 border border-neutral-700 rounded-md hover:bg-neutral-800"
            >
              ⋯
            </button>

            {#if openMenus[item.serviceName]}
              <div class="absolute right-0 top-full mt-1 w-48 bg-neutral-900 border border-neutral-700 rounded-md shadow-lg z-10">
                <div class="py-1">
                  <button
                    onclick={() => startService(item.serviceName)}
                    class="w-full text-left px-4 py-2 text-sm hover:bg-neutral-800 flex items-center gap-2"
                    disabled={item.service?.state === "running"}
                  >
                    Start Service
                  </button>
                  <button
                    onclick={() => stopService(item.serviceName)}
                    class="w-full text-left px-4 py-2 text-sm hover:bg-neutral-800 flex items-center gap-2"
                    disabled={item.service?.state === "stopped"}
                  >
                    Stop Service
                  </button>
                  <button
                    onclick={() => restartService(item.serviceName)}
                    class="w-full text-left px-4 py-2 text-sm hover:bg-neutral-800 flex items-center gap-2"
                  >
                    Restart Service
                  </button>
                  <hr class="border-neutral-700 my-1">
                  <button
                    onclick={() => toggleAutoRestart(item.serviceName, item.autoRestart)}
                    class="w-full text-left px-4 py-2 text-sm hover:bg-neutral-800 flex items-center gap-2"
                  >
                    {item.autoRestart ? "Disable" : "Enable"} Auto-restart
                  </button>
                </div>
              </div>
            {/if}
          </div>
        </li>
      {/each}
    </ul>
  </div>

  <div id="logs" class="card col-span-12 lg:col-span-8">
    <div class="flex items-center justify-between mb-3">
      <h2 class="font-semibold">Event Logs</h2>
      <div class="flex items-center gap-2 flex-wrap min-w-0 flex-1 justify-end">
        <select
          id="logService"
          class="bg-neutral-950 border border-neutral-800 rounded-md px-2 py-1 text-sm shrink-0 w-32 sm:w-40"
        >
          <option value="">Any service</option>
          {#each watchlistState.items as item}
            <option value={item.serviceName}>{item.serviceName}</option>
          {/each}
        </select>
        <select
          id="logLevel"
          class="bg-neutral-950 border border-neutral-800 rounded-md px-2 py-1 text-sm shrink-0"
        >
          <option value="">Any level</option>
          <option>INFO</option>
          <option>WARN</option>
          <option>ERROR</option>
        </select>
        <select
          id="logSince"
          class="bg-neutral-950 border border-neutral-800 rounded-md px-2 py-1 text-sm shrink-0"
        >
          <option value="15">Last 15m</option>
          <option value="60">Last 60m</option>
          <option value="1440">Last 24h</option>
        </select>
        <button
          id="runQuery"
          class="px-3 py-1.5 rounded-md bg-white/10 hover:bg-white/20 border border-white/10 text-sm shrink-0 ml-2"
        >
          Run
        </button>
      </div>
    </div>
    <div
      id="logsList"
      class="h-64 overflow-auto scroll-slim space-y-1 text-sm font-mono bg-neutral-950 rounded-xl p-3 border border-neutral-800"
    ></div>
  </div>
</section>

<!-- Click outside to close menus -->
<svelte:window
  on:click={(e) => {
    if (!(e.target instanceof Element) || !e.target.closest(".relative")) {
      openMenus = {};
    }
  }}
/>
