<script>
  import { onMount, onDestroy } from 'svelte';
  import uPlot from 'uplot';
  import 'uplot/dist/uPlot.min.css';
  export let title = '';
  export let data;         // { t, cpu, mem, memIsMB:true }
  export let height = 300;
  let el, u, ro;

  function build() {
    u = new uPlot({
      width: el.clientWidth, height,
      scales: { 
        x: { time: true }, 
        y: { 
          auto: true,
          range: (u, dataMin, dataMax) => {
            // Ensure we always have at least a 2% range for visibility
            const minRange = 2;
            const padding = 0.1; // 10% padding
            
            if (dataMax - dataMin < minRange) {
              const center = (dataMax + dataMin) / 2;
              return [Math.max(0, center - minRange/2), center + minRange/2];
            }
            
            const range = dataMax - dataMin;
            const paddedMin = Math.max(0, dataMin - range * padding);
            const paddedMax = dataMax + range * padding;
            
            return [paddedMin, paddedMax];
          }
        }
      },
      series: [
        { label: 'Time' },
        { label: 'CPU %', stroke: '#22d3ee', width: 2, scale: 'y' },
        { label: 'Memory %', stroke: '#a78bfa', width: 2, scale: 'y' }
      ],
      axes: [
        { 
          stroke: '#6b6b6b', 
          grid: { stroke: '#262626' }
        },
        { 
          stroke: '#6b6b6b', 
          grid: { stroke: '#262626' },
          values: (u, vals) => vals.map(v => v.toFixed(1) + '%')
        }
      ],
      cursor: { drag: { x: true, y: false } }, 
      legend: { show: true }
    }, [data.t || [], data.cpu || [], data.mem || []], el);

    ro = new ResizeObserver(()=>u.setSize({ width: el.clientWidth, height }));
    ro.observe(el);
  }
  onMount(build);
  onDestroy(()=>{ ro?.disconnect(); u?.destroy(); });
  $: if (u && data) {
    u.setData([data.t || [], data.cpu || [], data.mem || []]); 
  }
</script>

<div class="card chart-card">
  <div class="flex items-center justify-between mb-2">
    <div class="font-semibold">{title}</div>
    <slot name="actions" />
  </div>
  <div bind:this={el} style="width:100%"></div>
</div>
