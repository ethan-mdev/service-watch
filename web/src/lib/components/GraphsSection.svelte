<!-- ChartsPanel.svelte -->
<script>
  import { chartsState } from "../stores/charts.svelte.js";
  import Chart from "./Chart.svelte";
</script>

{#if chartsState.pinned.length === 0}
  <div class="card text-center py-10">
    <div class="text-lg font-semibold mb-2">No charts pinned</div>
    <div class="text-sm text-neutral-400 mb-4">
      Add a service to the watchlist and click <span
        class="inline-block align-middle px-2 py-0.5 border border-neutral-700 rounded-md"
        >+</span
      > to watch its metrics.
    </div>
  </div>
{:else}
  <div class="space-y-4">
    {#each chartsState.pinned as name}
      <Chart title={name} data={chartsState.series[name]}>
        <button
          slot="actions"
          class="chip"
          on:click={() => chartsState.unpin(name)}>Remove</button
        >
      </Chart>
    {/each}
  </div>
{/if}
