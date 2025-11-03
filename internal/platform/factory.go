package platform

import "github.com/ethan-mdev/service-watch/internal/core"

// MakeServiceManager returns the OS-specific implementation.
func MakeServiceManager() core.ServiceManager {
	return newServiceManager() // provided by windows_services.go or linux_services.go
}
