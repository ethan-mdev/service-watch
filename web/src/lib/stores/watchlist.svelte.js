export const watchlistState = $state({
  items: [],
  get numItems() {
    return this.items.length;
  },
  get numRunning() {
    return this.items.filter(item => item.service?.state === 'running').length;
  },
  get numStopped() {
    return this.items.filter(item => item.service?.state === 'stopped').length;
  },
});

// API functions for watchlist
export const watchlistAPI = {
  async fetch() {
    try {
      const response = await fetch('/v1/watchlist');
      if (response.ok) {
        const data = await response.json();
        watchlistState.items = data.items || [];
        console.log('Watchlist fetched:', watchlistState.items);
      } else {
        console.error('Failed to fetch watchlist:', response.statusText);
      }
    } catch (err) {
      console.error('Error fetching watchlist:', err);
    }
  },

  async add(serviceName) {
    try {
      console.log('Adding service to watchlist:', serviceName);
      const response = await fetch('/v1/watchlist', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ serviceName })
      });
      
      if (response.ok) {
        console.log('Service added successfully');
        await this.fetch();
      } else {
        const errorText = await response.text();
        console.error('Failed to add service:', response.statusText, errorText);
      }
    } catch (err) {
      console.error('Error adding service to watchlist:', err);
    }
  },

  async remove(serviceName) {
    try {
      console.log('Removing service from watchlist:', serviceName);
      const response = await fetch(`/v1/watchlist/${encodeURIComponent(serviceName)}`, {
        method: 'DELETE'
      });
      
      if (response.ok) {
        console.log('Service removed successfully');
        await this.fetch();
      } else {
        const errorText = await response.text();
        console.error('Failed to remove service:', response.statusText, errorText);
      }
    } catch (err) {
      console.error('Error removing service from watchlist:', err);
    }
  },

  async update(serviceName, booleanAutoRestart) {
    try {
      console.log('Updating service in watchlist:', serviceName);
      const response = await fetch(`/v1/watchlist/${encodeURIComponent(serviceName)}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ autoRestart: booleanAutoRestart })
      });
    } catch (err) {
      console.error('Error updating service in watchlist:', err);
    }
  }
};