//go:build !windows && !linux && !darwin
// +build !windows,!linux,!darwin

package services

// Fallback for unsupported platforms
func NewServiceLister() ServiceLister {
	panic("NewServiceLister not implemented for this platform")
}
