// API functions for service control
export const servicesAPI = {
  async start(serviceName) {
    try {
      const response = await fetch(
        `/v1/services/${encodeURIComponent(serviceName)}/start`,
        {
          method: "POST",
        }
      );
      if (response.ok) {
        console.log('Service started successfully:', serviceName);
        return true;
      } else {
        const errorText = await response.text();
        console.error('Failed to start service:', response.statusText, errorText);
        return false;
      }
    } catch (err) {
      console.error("Error starting service:", err);
      return false;
    }
  },

  async stop(serviceName) {
    try {
      const response = await fetch(
        `/v1/services/${encodeURIComponent(serviceName)}/stop`,
        {
          method: "POST",
        }
      );
      if (response.ok) {
        console.log('Service stopped successfully:', serviceName);
        return true;
      } else {
        const errorText = await response.text();
        console.error('Failed to stop service:', response.statusText, errorText);
        return false;
      }
    } catch (err) {
      console.error("Error stopping service:", err);
      return false;
    }
  },

  async restart(serviceName) {
    try {
      const response = await fetch(
        `/v1/services/${encodeURIComponent(serviceName)}/restart`,
        {
          method: "POST",
        }
      );
      if (response.ok) {
        console.log('Service restarted successfully:', serviceName);
        return true;
      } else {
        const errorText = await response.text();
        console.error('Failed to restart service:', response.statusText, errorText);
        return false;
      }
    } catch (err) {
      console.error("Error restarting service:", err);
      return false;
    }
  }
};