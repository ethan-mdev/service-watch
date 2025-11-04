<script>
  import { onMount, onDestroy } from 'svelte';

  let services = $state([]);
  let events = $state([]);
  let connected = $state(false);
  let eventSource = null;

  onMount(async () => {
    // Load initial data
    await loadWatchlist();

    // Connect to SSE
    eventSource = new EventSource('/v1/events');

    eventSource.onopen = () => {
      connected = true;
    };

    eventSource.onerror = () => {
      connected = false;
    };

    eventSource.addEventListener('connected', (e) => {
      addEvent('Connected to SSE');
    });

    eventSource.addEventListener('service_restarting', (e) => {
      const data = JSON.parse(e.data);
      addEvent(`Attempting to restart ${data.service_name}`);
    });

    eventSource.addEventListener('service_restart_success', (e) => {
      const data = JSON.parse(e.data);
      addEvent(`✓ ${data.service_name} restarted (count: ${data.restart_count})`);
      loadWatchlist();
    });

    eventSource.addEventListener('service_restart_failed', (e) => {
      const data = JSON.parse(e.data);
      addEvent(`✗ Failed to restart ${data.service_name}: ${data.error}`);
    });

    eventSource.addEventListener('service_failed', (e) => {
      const data = JSON.parse(e.data);
      addEvent(`✗ ${data.service_name} exceeded max restart attempts`);
    });
  });

  onDestroy(() => {
    if (eventSource) {
      eventSource.close();
      connected = false;
    }
  });

  async function loadWatchlist() {
    const res = await fetch('/v1/watchlist');
    const data = await res.json();
    services = data.items || [];
  }

  function addEvent(message) {
    const timestamp = new Date().toLocaleTimeString();
    events = [`[${timestamp}] ${message}`, ...events].slice(0, 20);
  }
</script>

<h1>Service Watch</h1>
<p>Status: {connected ? '🟢 Connected' : '🔴 Disconnected'}</p>

<h2>Watchlist</h2>
<table border="1">
  <thead>
    <tr>
      <th>Service</th>
      <th>State</th>
      <th>CPU %</th>
      <th>Memory MB</th>
      <th>Restart Count</th>
      <th>Auto-Restart</th>
    </tr>
  </thead>
  <tbody>
    {#each services as item}
      <tr>
        <td>{item.serviceName}</td>
        <td>{item.service?.state || 'unknown'}</td>
        <td>{item.service?.cpuPercent?.toFixed(1) || '-'}</td>
        <td>{item.service?.memoryMB?.toFixed(1) || '-'}</td>
        <td>{item.restartCount}</td>
        <td>{item.autoRestart ? 'Yes' : 'No'}</td>
      </tr>
    {/each}
  </tbody>
</table>

<h2>Live Events</h2>
<ul>
  {#each events as event}
    <li>{event}</li>
  {/each}
</ul>