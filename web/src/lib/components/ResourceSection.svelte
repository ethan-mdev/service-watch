<script>
  import { watchlistState, watchlistAPI } from '../stores/watchlist.svelte.js';
  import { sseState, sseManager } from '../stores/sse.svelte.js';
  import { metricsState, metricsAPI } from '../stores/metrics.svelte.js';
  import { onMount, onDestroy } from 'svelte';

  let metricsInterval;
  let serviceFailed24h = $state(0);

  onMount(async () => {
    await watchlistAPI.fetch();
    sseManager.connect();
    
    // Fetch failed services count on mount
    serviceFailed24h = await metricsAPI.getServiceFailed('24h');
    
    // Set up periodic refresh (every 5 minutes) for failed services
    metricsInterval = metricsAPI.startPeriodicRefresh(
      { event: 'service_failed', since: '24h' },
      5 * 60 * 1000
    );
    
    // Update local state when new data comes in
    setInterval(async () => {
      serviceFailed24h = await metricsAPI.getServiceFailed('24h');
    }, 5 * 60 * 1000);
  });

  onDestroy(() => {
    sseManager.disconnect();
    if (metricsInterval) {
      clearInterval(metricsInterval);
    }
  });

  async function handleAddService() {
    const serviceName = prompt('Enter service name to watch:');
    if (serviceName?.trim()) {
      await watchlistAPI.add(serviceName.trim());
    }
  }
</script>

<section class="grid-12">
  <div class="card col-span-12 lg:col-span-7">
    <h3 class="font-semibold mb-4">Host</h3>
    
    <div class="grid grid-cols-2 gap-8 place-items-center">
      <div class="flex flex-col items-center justify-center">
        <svg class="gauge" viewBox="0 0 140 140">
          <circle class="bg" cx="70" cy="70" r="60" fill="none" stroke-width="14"/>
          <circle id="cpuArc" class="fg" cx="70" cy="70" r="60" fill="none" stroke-width="14"
                  stroke-dasharray="377" stroke-dashoffset={377 - (377 * (sseState.hostResources.cpuPercent || 0) / 100)} stroke-linecap="round"
                  transform="rotate(-90 70 70)"/>
          <text id="cpuTxt" x="70" y="75" text-anchor="middle">CPU {(sseState.hostResources.cpuPercent || 0).toFixed(1)}%</text>
        </svg>
      </div>
      <div class="flex flex-col items-center justify-center">
        <svg class="gauge" viewBox="0 0 140 140">
          <circle class="bg" cx="70" cy="70" r="60" fill="none" stroke-width="14"/>
          <circle id="memArc" class="fg" cx="70" cy="70" r="60" fill="none" stroke-width="14"
                  stroke-dasharray="377" stroke-dashoffset={377 - (377 * (sseState.hostResources.usedPercent || 0) / 100)} stroke-linecap="round"
                  transform="rotate(-90 70 70)"/>
          <text id="memTxt" x="70" y="75" text-anchor="middle">Mem {(sseState.hostResources.usedPercent || 0).toFixed(1)}%</text>
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
        <div class="text-xs">Failed (24h) {metricsState.loading ? '-' : ''}</div>
        <div class="text-2xl font-semibold mt-1">{serviceFailed24h}</div>
      </div>
      <div class="rounded-xl border border-neutral-800 bg-neutral-950 p-3">
        <div class="text-xs">Running</div>
        <div class="text-2xl font-semibold mt-1">{watchlistState.numRunning}</div>
      </div>
      <div class="rounded-xl border border-neutral-800 bg-neutral-950 p-3">
        <div class="text-xs">Stopped</div>
        <div class="text-2xl font-semibold mt-1">{watchlistState.numStopped}</div>
      </div>
    </div>
  </div>
</section>
