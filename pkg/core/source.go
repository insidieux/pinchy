package core

import (
	"context"
)

type (
	// Source provides information about Services to be registered in Registry.
	// Source does not know about changes between Fetch calls, it must just return actual state of source.
	Source interface {
		Fetch(ctx context.Context) (Services, error)
	}
)
