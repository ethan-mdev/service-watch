export const metricsState = $state({
    serviceFailed24h: 0,
    loading: false,
    lastFetch: null
});

export const metricsAPI = {
    async fetchServiceFailed(since = '24h') {
        try {
            metricsState.loading = true;
            const response = await fetch(`/v1/metrics?event=service_failed&since=${since}`);
            if (response.ok) {
                const data = await response.json();
                metricsState.serviceFailed24h = data.count || 0;
                metricsState.lastFetch = new Date();
                console.log('Service failed metrics fetched:', data);
            } else {
                console.error('Failed to fetch service failed metrics:', response.statusText);
            }
        } catch (err) {
            console.error('Error fetching service failed metrics:', err);
        } finally {
            metricsState.loading = false;
        }
    },

    // Periodic refresh
    startPeriodicRefresh(intervalMs = 5 * 60 * 1000) { // 5 minutes default
        return setInterval(() => {
            this.fetchServiceFailed();
        }, intervalMs);
    }
};