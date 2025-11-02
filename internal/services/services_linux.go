//go:build linux

package services

import (
	"context"
	"errors"
)

// TODO: use systemd dbus client (github.com/coreos/go-systemd/v22/dbus)
type linuxLister struct{}

func NewServiceLister() ServiceLister { return &linuxLister{} }

func (l *linuxLister) ListServices(ctx context.Context) ([]Service, error) {
	return nil, errors.New("service listing not implemented on Linux yet")
}
