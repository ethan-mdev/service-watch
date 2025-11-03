package core

import "context"

// ServiceManager abstracts OS service control.
type ServiceManager interface {
	List(ctx context.Context) ([]Service, error)
	Get(ctx context.Context, name string) (Service, error)
	Start(ctx context.Context, name string) error
	Stop(ctx context.Context, name string) error
	Restart(ctx context.Context, name string) error
}
