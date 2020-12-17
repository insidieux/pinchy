package core

import (
	"context"
)

type (
	// Registry is endpoint/upstream for storing information about Services fetched from Source.
	// Registry must implements 3 methods to make store clean: fetch current state, deregister orphans, register/update incoming
	Registry interface {
		Fetch(ctx context.Context) (Services, error)
		Register(ctx context.Context, service *Service) error
		Deregister(ctx context.Context, serviceID string) error
	}
)
