package core

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/thoas/go-funk"
)

type (
	ManagerInterface interface {
		Run(ctx context.Context) error
	}
	Manager struct {
		source      Source
		registry    Registry
		logger      LoggerInterface
		exitOnError ManagerExitOnError
	}
	ManagerExitOnError bool
	managerError       []error
)

func NewManager(source Source, registry Registry, logger LoggerInterface, exitOnError ManagerExitOnError) ManagerInterface {
	return &Manager{
		source:      source,
		registry:    registry,
		logger:      logger,
		exitOnError: exitOnError,
	}
}

func (m *Manager) Run(ctx context.Context) error {
	incoming, err := m.source.Fetch(ctx)
	if err != nil {
		return errors.Wrap(err, `failed to fetch services from source`)
	}

	registered, err := m.registry.Fetch(ctx)
	if err != nil {
		return errors.Wrap(err, `failed to fetch services from registry`)
	}

	orphan := m.findOrphan(incoming, registered)
	if len(orphan) > 0 {
		if err := m.deregisterServices(ctx, orphan); err != nil {
			err := errors.Wrap(err, `failed to deregister services`)
			m.logger.Error(err.Error())
			if m.exitOnError {
				return err
			}
		}
	}
	if err := m.registerServices(ctx, incoming); err != nil {
		err := errors.Wrap(err, `failed to register services`)
		m.logger.Error(err.Error())
		if m.exitOnError {
			return err
		}
	}
	return nil
}

func (m *Manager) findOrphan(incoming Services, registered Services) Services {
	_, right := funk.Difference(incoming.IDs(), registered.IDs())
	orphan := make(Services, 0)
	for _, id := range cast.ToStringSlice(right) {
		if service := registered.Lookup(id); service != nil {
			orphan = append(orphan, service)
		}
	}
	return orphan
}

func (m *Manager) deregisterServices(ctx context.Context, services Services) error {
	me := new(managerError)
	for _, service := range services {
		if err := m.registry.Deregister(ctx, service.RegistrationID()); err != nil {
			me.Add(errors.Wrapf(err, `failed to deregister service "%s" from registry`, service.RegistrationID()))
			continue
		}
	}
	if me.HasErrors() {
		return me
	}
	return nil
}

func (m *Manager) registerServices(ctx context.Context, services Services) error {
	me := new(managerError)
	for _, service := range services {
		if err := m.registry.Register(ctx, service); err != nil {
			me.Add(errors.Wrapf(err, `failed to register service "%s" in registry`, service.RegistrationID()))
			continue
		}
	}
	if me.HasErrors() {
		return me
	}
	return nil
}

func (e *managerError) Add(err error) {
	*e = append(*e, err)
}

func (e *managerError) HasErrors() bool {
	return len(*e) > 0
}

func (e *managerError) String() string {
	var slice []string
	for _, err := range *e {
		slice = append(slice, err.Error())
	}
	return strings.Join(slice, `; `)
}

func (e *managerError) Error() string {
	return e.String()
}
