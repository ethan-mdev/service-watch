<section class="grid-12">
  <div id="watch" class="card col-span-12 lg:col-span-4">
    <div class="flex items-center justify-between mb-2">
      <h2 class="text-lg font-semibold">Watchlist</h2>
    </div>
    <p class="text-sm text-neutral-400 mb-4">Add a service to the watchlist to monitor its metrics.</p>
    <!-- Watchlist items will be dynamically added here -->
    <ul id="watchUl" class="space-y-2 hidden"></ul>
  </div>

  <div id="logs" class="card col-span-12 lg:col-span-8">
    <div class="flex items-center justify-between mb-3">
      <div class="flex items-center gap-2">
        <select
          id="logService"
          class="bg-neutral-950 border border-neutral-800 rounded-md px-1 py-1 text-sm min-w-40"
        >
          <option value="">Any service</option>
        </select>
        <select
          id="logLevel"
          class="bg-neutral-950 border border-neutral-800 rounded-md px-1 py-1 text-sm"
        >
          <option value="">Any level</option>
          <option>INFO</option>
          <option>WARN</option>
          <option>ERROR</option>
        </select>
        <select
          id="logSince"
          class="bg-neutral-950 border border-neutral-800 rounded-md px-1 py-1 text-sm"
        >
          <option value="15">Last 15m</option>
          <option value="60">Last 60m</option>
          <option value="1440">Last 24h</option>
        </select>
        <input
          id="logText"
          placeholder="contains..."
          class="bg-neutral-950 border border-neutral-800 rounded-md px-1 py-1 text-sm"
        />
        <button
          id="runQuery"
          class="px-3 py-1.5 rounded-md bg-white/10 hover:bg-white/20 border border-white/10 text-sm"
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
