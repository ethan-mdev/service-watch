<section class="grid-12">
  <div class="card col-span-12 lg:col-span-7">
    <h3 class="font-semibold mb-4">Host</h3>
    <div class="grid grid-cols-2 gap-8 place-items-center">
      <div class="flex flex-col items-center justify-center">
        <svg class="gauge" viewBox="0 0 140 140">
          <circle class="bg" cx="70" cy="70" r="60" fill="none" stroke-width="14"/>
          <circle id="cpuArc" class="fg" cx="70" cy="70" r="60" fill="none" stroke-width="14"
                  stroke-dasharray="377" stroke-dashoffset="377" stroke-linecap="round"
                  transform="rotate(-90 70 70)"/>
          <text id="cpuTxt" x="70" y="75" text-anchor="middle">CPU 0%</text>
        </svg>
      </div>
      <div class="flex flex-col items-center justify-center">
        <svg class="gauge" viewBox="0 0 140 140">
          <circle class="bg" cx="70" cy="70" r="60" fill="none" stroke-width="14"/>
          <circle id="memArc" class="fg" cx="70" cy="70" r="60" fill="none" stroke-width="14"
                  stroke-dasharray="377" stroke-dashoffset="377" stroke-linecap="round"
                  transform="rotate(-90 70 70)"/>
          <text id="memTxt" x="70" y="75" text-anchor="middle">Mem 0%</text>
        </svg>
      </div>
    </div>
  </div>

  <div class="card col-span-12 lg:col-span-5">
    <div class="flex items-center justify-between">
      <h3 class="font-semibold">Add Service</h3>
      <button class="chip">+ Add</button>
    </div>
    <div class="mt-4 grid grid-cols-2 gap-3">
      <div class="rounded-xl border border-neutral-800 bg-neutral-950 p-3">
        <div class="text-xs subtle">Watched</div>
        <div id="kpiWatched" class="text-2xl font-semibold mt-1">-</div>
      </div>
      <div class="rounded-xl border border-neutral-800 bg-neutral-950 p-3">
        <div class="text-xs subtle">Alerts (24h)</div>
        <div id="kpiAlerts" class="text-2xl font-semibold mt-1">-</div>
      </div>
      <div class="rounded-xl border border-neutral-800 bg-neutral-950 p-3">
        <div class="text-xs subtle">Running</div>
        <div id="kpiRunning" class="text-2xl font-semibold mt-1">-</div>
      </div>
      <div class="rounded-xl border border-neutral-800 bg-neutral-950 p-3">
        <div class="text-xs subtle">Stopped</div>
        <div id="kpiStopped" class="text-2xl font-semibold mt-1">-</div>
      </div>
    </div>
  </div>
</section>
