//go:build windows

package platform

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"github.com/ethan-mdev/service-watch/internal/core"
)

type winSvc struct{}

func newServiceManager() core.ServiceManager { return &winSvc{} }

func (w *winSvc) List(ctx context.Context) ([]core.Service, error) {
	m, err := mgr.Connect()
	if err != nil {
		return nil, err
	}
	defer m.Disconnect()

	names, err := m.ListServices()
	if err != nil {
		return nil, err
	}

	out := make([]core.Service, 0, len(names))
	for _, name := range names {
		s, err := m.OpenService(name)
		if err != nil {
			continue
		}
		config, err := s.Config()
		if err != nil {
			s.Close()
			continue
		}
		status, err := s.Query()
		s.Close()
		if err != nil {
			continue
		}
		out = append(out, core.Service{
			Name:        name,
			DisplayName: config.DisplayName,
			State:       stateToString(status.State),
		})
	}

	return out, nil
}

func (w *winSvc) Get(ctx context.Context, name string) (core.Service, error) {
	m, err := mgr.Connect()
	if err != nil {
		return core.Service{}, err
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return core.Service{}, fmt.Errorf("service not found: %s", name)
	}
	defer s.Close()

	config, err := s.Config()
	if err != nil {
		return core.Service{}, err
	}

	status, err := s.Query()
	if err != nil {
		return core.Service{}, err
	}

	return core.Service{
		Name:        name,
		DisplayName: config.DisplayName,
		State:       stateToString(status.State),
	}, nil
}

func (w *winSvc) Start(ctx context.Context, name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("service not found: %s", name)
	}
	defer s.Close()
	if err := s.Start(); err != nil {
		return err
	}
	return waitState(ctx, s, svc.Running)
}

func (w *winSvc) Stop(ctx context.Context, name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("service not found: %s", name)
	}
	defer s.Close()
	if _, err := s.Control(svc.Stop); err != nil {
		return err
	}
	return waitState(ctx, s, svc.Stopped)
}

func (w *winSvc) Restart(ctx context.Context, name string) error {
	if err := w.Stop(ctx, name); err != nil {
		return err
	}
	select {
	case <-time.After(500 * time.Millisecond):
	case <-ctx.Done():
		return ctx.Err()
	}
	return w.Start(ctx, name)
}

func waitState(ctx context.Context, s *mgr.Service, want svc.State) error {
	tick := time.NewTicker(200 * time.Millisecond)
	defer tick.Stop()
	timeout := time.NewTimer(20 * time.Second)
	defer timeout.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout.C:
			return fmt.Errorf("timeout waiting for %v", want)
		case <-tick.C:
			st, err := s.Query()
			if err != nil {
				return err
			}
			if st.State == want {
				return nil
			}
		}
	}
}

func stateToString(state svc.State) string {
	switch state {
	case svc.Stopped:
		return "stopped"
	case svc.StartPending:
		return "start_pending"
	case svc.StopPending:
		return "stop_pending"
	case svc.Running:
		return "running"
	case svc.ContinuePending:
		return "continue_pending"
	case svc.PausePending:
		return "pause_pending"
	case svc.Paused:
		return "paused"
	default:
		return "unknown"
	}
}
