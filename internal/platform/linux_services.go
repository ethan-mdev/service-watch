//go:build linux

package platform

import (
	"context"
	"errors"

	"github.com/ethan-mdev/service-watch/internal/core"
)

type stubSvc struct{}

func newServiceManager() core.ServiceManager { return &stubSvc{} }

func (s *stubSvc) List(ctx context.Context) ([]core.Service, error) {
	return nil, errors.New("not implemented on linux")
}
func (s *stubSvc) Get(ctx context.Context, name string) (core.Service, error) {
	return core.Service{}, errors.New("not implemented on linux")
}
func (s *stubSvc) Start(ctx context.Context, name string) error {
	return errors.New("not implemented on linux")
}
func (s *stubSvc) Stop(ctx context.Context, name string) error {
	return errors.New("not implemented on linux")
}
func (s *stubSvc) Restart(ctx context.Context, name string) error {
	return errors.New("not implemented on linux")
}
