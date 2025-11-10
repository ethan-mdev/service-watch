// src/stores/charts.svelte.js
export const chartsState = $state({
  pinned: [],                       
  series: {},                        
  maxPoints: 600,

  isPinned(n) { 
    return this.pinned.includes(n); 
  },
  
  toggle(n) { 
    if (this.isPinned(n)) {
      this.unpin(n);
    } else {
      this.pin(n);
    }
  },
  
  pin(n) { 
    console.log('Pinning service:', n);
    if (!this.isPinned(n)) {
      this.pinned.push(n);
      this.ensure(n);
    }
  },
  
  unpin(n) { 
    console.log('Unpinning service:', n);
    const index = this.pinned.indexOf(n);
    if (index > -1) {
      this.pinned.splice(index, 1);
    }
  },

  ensure(n) {
    if (!this.series[n]) {
      this.series[n] = { t:[], cpu:[], mem:[], memIsMB:false };
    }
  },

  pushData(n, t, cpuPct, memMB, totalSystemMB) {
    this.ensure(n);
    const s = this.series[n];
    s.memIsMB = false;
    // Convert milliseconds to seconds for uPlot
    const timestamp = typeof t === 'number' ? t / 1000 : Date.now() / 1000;
    s.t.push(timestamp); 
    s.cpu.push(cpuPct ?? 0); 
    
    // Convert memory MB to percentage of total system memory
    const memPercent = totalSystemMB > 0 ? ((memMB ?? 0) / totalSystemMB) * 100 : 0;
    s.mem.push(memPercent);
    
    console.log(`Chart data for ${n}: timestamp=${timestamp}, cpu=${cpuPct}%, mem=${memMB}MB (${memPercent.toFixed(1)}%)`);
    
    const over = s.t.length - this.maxPoints;
    if (over > 0) { 
      s.t.splice(0, over); 
      s.cpu.splice(0, over); 
      s.mem.splice(0, over); 
    }
  },

  // Add some initial data points when pinning a service
  initializeWithCurrentData(serviceName, currentData, totalSystemMB = 16384) {
    this.ensure(serviceName);
    const now = Date.now() / 1000;
    const s = this.series[serviceName];
    
    const memPercent = totalSystemMB > 0 ? ((currentData?.memoryMB ?? 0) / totalSystemMB) * 100 : 0;
    
    // Add just a single current point to start
    s.t.push(now);
    s.cpu.push(currentData?.cpuPercent ?? 0);
    s.mem.push(memPercent);
    
    console.log(`Initialized chart for ${serviceName} with current data point`);
  }
});
