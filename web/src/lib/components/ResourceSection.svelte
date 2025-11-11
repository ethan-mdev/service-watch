<script>
  import { watchlistState, watchlistAPI } from "../stores/watchlist.svelte.js";
  import { servicesAPI } from "../stores/services.svelte.js";
  import { sseState, sseManager } from "../stores/sse.svelte.js";
  import { metricsState, metricsAPI } from "../stores/metrics.svelte.js";
  import { onMount, onDestroy } from "svelte";

  let metricsInterval;
  let serviceFailed24h = $state(0);
  let showServiceModal = $state(false);
  let availableServices = $state([]);
  let filteredServices = $state([]);
  let searchTerm = $state("");

  onMount(async () => {
    await watchlistAPI.fetch();
    sseManager.connect();

    // Fetch failed services count on mount
    serviceFailed24h = await metricsAPI.getServiceFailed("24h");

    // Set up periodic refresh (every 5 minutes) for failed services
    metricsInterval = metricsAPI.startPeriodicRefresh(
      { event: "service_failed", since: "24h" },
      5 * 60 * 1000,
    );

    // Update local state when new data comes in
    setInterval(
      async () => {
        serviceFailed24h = await metricsAPI.getServiceFailed("24h");
      },
      5 * 60 * 1000,
    );
  });

  onDestroy(() => {
    sseManager.disconnect();
    if (metricsInterval) {
      clearInterval(metricsInterval);
    }
  });

  async function handleAddService() {
    // Fetch available services
    availableServices = await servicesAPI.fetchAvailable();
    filteredServices = availableServices.filter(
      (service) =>
        !watchlistState.items.some((item) => item.serviceName === service.name),
    );
    showServiceModal = true;
  }

  function filterServices() {
    if (searchTerm.trim() === "") {
      filteredServices = availableServices.filter(
        (service) =>
          !watchlistState.items.some(
            (item) => item.serviceName === service.name,
          ),
      );
    } else {
      filteredServices = availableServices.filter(
        (service) =>
          (service.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
            service.displayName
              .toLowerCase()
              .includes(searchTerm.toLowerCase())) &&
          !watchlistState.items.some(
            (item) => item.serviceName === service.name,
          ),
      );
    }
  }

  async function selectService(serviceName) {
    await watchlistAPI.add(serviceName);
    showServiceModal = false;
    searchTerm = "";
  }

  function closeModal() {
    showServiceModal = false;
    searchTerm = "";
  }

  function stopPropagation(event) {
    event.stopPropagation();
  }

  function handleKeydown(event) {
    if (event.key === "Escape") {
      closeModal();
    }
  }

  function handleModalKeydown(event) {
    event.stopPropagation();
  }

  // Filter services when search term changes
  $effect(() => {
    if (searchTerm !== undefined) filterServices();
  });
</script>

<section class="grid-12">
  <div class="card col-span-12 lg:col-span-7">
    <h3 class="font-semibold mb-4">Host</h3>

    <div class="grid grid-cols-2 gap-8 place-items-center">
      <div class="flex flex-col items-center justify-center">
        <svg class="gauge" viewBox="0 0 140 140">
          <circle
            class="bg"
            cx="70"
            cy="70"
            r="60"
            fill="none"
            stroke-width="14"
          />
          <circle
            id="cpuArc"
            class="fg"
            cx="70"
            cy="70"
            r="60"
            fill="none"
            stroke-width="14"
            stroke-dasharray="377"
            stroke-dashoffset={377 -
              (377 * (sseState.hostResources.cpuPercent || 0)) / 100}
            stroke-linecap="round"
            transform="rotate(-90 70 70)"
          />
          <text id="cpuTxt" x="70" y="75" text-anchor="middle"
            >CPU {(sseState.hostResources.cpuPercent || 0).toFixed(1)}%</text
          >
        </svg>
      </div>
      <div class="flex flex-col items-center justify-center">
        <svg class="gauge" viewBox="0 0 140 140">
          <circle
            class="bg"
            cx="70"
            cy="70"
            r="60"
            fill="none"
            stroke-width="14"
          />
          <circle
            id="memArc"
            class="fg"
            cx="70"
            cy="70"
            r="60"
            fill="none"
            stroke-width="14"
            stroke-dasharray="377"
            stroke-dashoffset={377 -
              (377 * (sseState.hostResources.usedPercent || 0)) / 100}
            stroke-linecap="round"
            transform="rotate(-90 70 70)"
          />
          <text id="memTxt" x="70" y="75" text-anchor="middle"
            >Mem {(sseState.hostResources.usedPercent || 0).toFixed(1)}%</text
          >
        </svg>
      </div>
    </div>
  </div>

  <div class="card col-span-12 lg:col-span-5">
    <div class="flex items-center justify-between">
      <h3 class="font-semibold">Watchlist</h3>
      <button class="chip" onclick={handleAddService}>+ Add</button>
    </div>
    <div class="mt-4 grid grid-cols-2 gap-3">
      <div class="rounded-xl border border-neutral-800 bg-neutral-950 p-3">
        <div class="text-xs">Watched</div>
        <div class="text-2xl font-semibold mt-1">{watchlistState.numItems}</div>
      </div>
      <div class="rounded-xl border border-neutral-800 bg-neutral-950 p-3">
        <div class="text-xs">
          Failed (24h) {metricsState.loading ? "-" : ""}
        </div>
        <div class="text-2xl font-semibold mt-1">{serviceFailed24h}</div>
      </div>
      <div class="rounded-xl border border-neutral-800 bg-neutral-950 p-3">
        <div class="text-xs">Running</div>
        <div class="text-2xl font-semibold mt-1">
          {watchlistState.numRunning}
        </div>
      </div>
      <div class="rounded-xl border border-neutral-800 bg-neutral-950 p-3">
        <div class="text-xs">Stopped</div>
        <div class="text-2xl font-semibold mt-1">
          {watchlistState.numStopped}
        </div>
      </div>
    </div>
  </div>
</section>

<!-- Service Selection Modal -->
{#if showServiceModal}
  <div
    class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
    role="dialog"
    aria-modal="true"
    aria-labelledby="modal-title"
    tabindex="-1"
    onclick={closeModal}
    onkeydown={handleKeydown}
  >
    <section
      class="bg-neutral-900 border border-neutral-700 rounded-lg p-6 w-full max-w-md max-h-[80vh] flex flex-col"
      role="button"
      tabindex="0"
      onclick={stopPropagation}
      onkeydown={handleModalKeydown}
    >
      <div class="flex items-center justify-between mb-6">
        <h3 id="modal-title" class="text-lg font-semibold">
          Add Service to Watchlist
        </h3>
        <button
          class="text-neutral-400 hover:text-white text-xl w-8 h-8 flex items-center justify-center"
          onclick={closeModal}
          aria-label="Close modal">×</button
        >
      </div>

      <input
        type="text"
        placeholder="Search services..."
        bind:value={searchTerm}
        class="w-full px-4 py-3 bg-neutral-950 border border-neutral-700 rounded-md mb-4 focus:outline-none focus:border-neutral-500 text-base"
        aria-label="Search services"
      />

      <div class="flex-1 overflow-y-auto -mx-2 px-2">
        <div class="space-y-2">
          {#each filteredServices as service}
            <button
              class="w-full text-left px-4 py-3 hover:bg-neutral-800 rounded-md flex items-start justify-between group transition-colors"
              onclick={() => selectService(service.name)}
              aria-label="Add {service.name} to watchlist"
            >
              <div class="flex-1 min-w-0">
                <div class="font-medium text-base truncate">{service.name}</div>
                <div class="text-sm text-neutral-400 truncate mt-1">
                  {service.displayName}
                </div>
              </div>
              <div class="ml-3 shrink-0">
                <span
                  class="text-xs px-3 py-1 rounded-full {service.state ===
                  'running'
                    ? 'bg-green-900 text-green-300'
                    : 'bg-red-900 text-red-300'}"
                >
                  {service.state}
                </span>
              </div>
            </button>
          {/each}

          {#if filteredServices.length === 0}
            <div class="text-center py-8 text-neutral-400">
              {searchTerm
                ? "No services found matching your search"
                : "No new services to add"}
            </div>
          {/if}
        </div>
      </div>
    </section>
  </div>
{/if}
