export const metricsState = $state({
    loading: false,
    lastFetch: null
});

export const metricsAPI = {
    async fetchMetrics(params = {}) {
        try {
            metricsState.loading = true;
            
            // Build query string from parameters
            const queryParams = new URLSearchParams();
            if (params.event) queryParams.set('event', params.event);
            if (params.service) queryParams.set('service', params.service);
            if (params.limit) queryParams.set('limit', params.limit.toString());
            if (params.since) queryParams.set('since', params.since);
            if (params.level) queryParams.set('level', params.level);
            
            const queryString = queryParams.toString();
            const url = `/v1/metrics${queryString ? '?' + queryString : ''}`;
            
            console.log('Fetching metrics:', url);
            
            const response = await fetch(url);
            if (response.ok) {
                const data = await response.json();
                metricsState.lastFetch = new Date();
                
                console.log('Metrics fetched:', data);
                return data;
            } else {
                console.error('Failed to fetch metrics:', response.statusText);
                return null;
            }
        } catch (err) {
            console.error('Error fetching metrics:', err);
            return null;
        } finally {
            metricsState.loading = false;
        }
    },

    // Convenience methods for common queries
    async getServiceFailed(since = '24h') {
        const data = await this.fetchMetrics({ event: 'service_failed', since });
        return data?.count || 0;
    },

    async getServiceLogs(params = {}) {
        // For log queries, return the full data array
        const data = await this.fetchMetrics(params);
        console.log('Raw metrics data:', data);
        // The API returns { count: X, items: [...] }
        if (data && Array.isArray(data.items)) {
            return data.items;
        } else if (Array.isArray(data)) {
            return data;
        }
        return [];
    },

    async getServiceRestarts(service, since = '24h') {
        const data = await this.fetchMetrics({ event: 'restart_success', service, since });
        return data?.count || 0;
    },

    // Periodic refresh
    startPeriodicRefresh(params, intervalMs = 5 * 60 * 1000) {
        return setInterval(() => {
            this.fetchMetrics(params);
        }, intervalMs);
    }
};