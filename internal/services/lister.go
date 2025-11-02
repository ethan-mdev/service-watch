package services

import "context"

type Service struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	Status      string `json:"status,omitempty"`    // "running", "stopped", etc.
	StartType   string `json:"startType,omitempty"` // "auto", "manual", "disabled"
}

type ServiceLister interface {
	ListServices(ctx context.Context) ([]Service, error)
}
