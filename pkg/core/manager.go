package core

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/thoas/go-funk"
)

type (
	// ManagerInterface is main unit between Source and Registry.
	// Implementation must provide full cycle for fetch services from Source and register them into Registry.
	ManagerInterface interface {
		Run(ctx context.Context) error
	}

	// Manager is built-in ManagerInterface implementation.
	Manager struct {
		source      Source
		registry    Registry
		logger      LoggerInterface
		exitOnError ManagerExitOnError
	}

	// ManagerExitOnError provide information how to handle errors and panics during manager.Run process.
	ManagerExitOnError bool

	managerError []error
)

// NewManager provider built-in ManagerInterface implementation
func NewManager(source Source, registry Registry, logger LoggerInterface, exitOnError ManagerExitOnError) ManagerInterface {
	return &Manager{
		source:      source,
		registry:    registry,
		logger:      logger,
		exitOnError: exitOnError,
	}
}

// Run contains next steps
// - Call Source.Fetch
// - Call Registry.Fetch
// - Check orphan Services fetched from Registry
// - Remove orphan Services
// - Register Services fetched from Source
func (m *Manager) Run(ctx context.Context) error {
	m.logger.Infoln(`Fetching services from source`)
	incoming, err := m.source.Fetch(ctx)
	if err != nil {
		return errors.Wrap(err, `failed to fetch services from source`)
	}

	m.logger.Infoln(`Fetching services from registry`)
	registered, err := m.registry.Fetch(ctx)
	if err != nil {
		return errors.Wrap(err, `failed to fetch services from registry`)
	}

	m.logger.Infoln(`Checking difference between registered services and incoming list`)
	orphan := m.findOrphan(incoming, registered)
	if len(orphan) > 0 {
		m.logger.Infof(`Deleting %d orphan services`, len(orphan))
		if err := m.deregisterServices(ctx, orphan); err != nil {
			err := errors.Wrap(err, `failed to deregister services`)
			m.logger.Error(err.Error())
			if m.exitOnError {
				return err
			}
		}
	}

	if len(incoming) > 0 {
		m.logger.Infoln(`Registering services in registry`)
		if err := m.registerServices(ctx, incoming); err != nil {
			err := errors.Wrap(err, `failed to register services`)
			m.logger.Error(err.Error())
			if m.exitOnError {
				return err
			}
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
		if err := m.registry.Deregister(ctx, service); err != nil {
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
