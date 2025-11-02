//go:build windows

package services

import (
	"context"
)

// TODO: replace with real WMI/Windows Service Manager calls.
// This stub returns a couple of fake rows for demonstration.
type winLister struct{}

func NewServiceLister() ServiceLister { return &winLister{} }

func (w *winLister) ListServices(ctx context.Context) ([]Service, error) {
	return []Service{
		{Name: "Spooler", DisplayName: "Print Spooler", Status: "running", StartType: "auto"},
		{Name: "W32Time", DisplayName: "Windows Time", Status: "running", StartType: "manual"},
	}, nil
}
