import { watchlistState, watchlistAPI } from './watchlist.svelte.js';
import { chartsState } from './charts.svelte.js';

export const sseState = $state({
  hostResources: {
    cpuPercent: 0,
    memoryMB: 0,
    totalMB: 0,
    usedPercent: 0
  }
});

let eventSource = null;

export const sseManager = {
  connect() {
    if (eventSource) return;
    
    eventSource = new EventSource('/v1/events');
    
    eventSource.onopen = () => {
      console.log('SSE connected');
    };
    
    eventSource.addEventListener('host_resources', (event) => {
      try {
        const data = JSON.parse(event.data);
        sseState.hostResources = data;
      } catch (err) {
        console.error('Error parsing host_resources:', err);
      }
    });

    eventSource.addEventListener('service_status', (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log('Service status update:', data);
        
        // Find the service in the watchlist and update its status
        const item = watchlistState.items.find(item => item.serviceName === data.serviceName);
        if (item && item.service) {
          // Update the service state in real-time
          item.service.state = data.state;
          item.service.cpuPercent = data.cpuPercent;
          item.service.memoryMB = data.memoryMB;
          item.service.uptimeSeconds = data.uptimeSec; // Note: backend sends uptimeSec
          item.service.pid = data.pid;
          
          console.log(`Updated ${data.serviceName} status to ${data.state}`);
          
          // Push data to charts if the service is pinned
          if (chartsState.isPinned(data.serviceName)) {
            const timestamp = Date.now();
            chartsState.pushData(
              data.serviceName,
              timestamp,
              data.cpuPercent || 0,
              data.memoryMB || 0,
              sseState.hostResources.totalMB || 16384 // Use total system memory
            );
            console.log(`Pushed chart data for ${data.serviceName}: CPU=${data.cpuPercent}%, MEM=${data.memoryMB}MB`);
          }
        }
      } catch (err) {
        console.error('Error parsing service_status:', err);
      }
    });

    eventSource.onerror = (error) => {
      console.error('SSE error:', error);
    };
  },
  
  disconnect() {
    if (eventSource) {
      eventSource.close();
      eventSource = null;
    }
  }
};