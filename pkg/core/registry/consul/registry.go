package consul

import (
	"context"

	"github.com/agrea/ptr"
	"github.com/hashicorp/consul/api"
	"github.com/insidieux/pinchy/pkg/core"
	"github.com/pkg/errors"
)

type (
	// Agent interface provide common function for work with Consul HTTP API
	Agent interface {
		Services() (map[string]*api.AgentService, error)
		ServiceRegister(service *api.AgentServiceRegistration) error
		ServiceDeregister(serviceID string) error
	}

	// Registry is implementation of core.Registry interface
	Registry struct {
		agent Agent
	}
)

// NewRegistry provide Registry as core.Registry implementation
func NewRegistry(agent Agent) *Registry {
	return &Registry{
		agent: agent,
	}
}

// Fetch make request for Agent.Services and try to cast result to core.Services
func (r *Registry) Fetch(_ context.Context) (core.Services, error) {
	registered, err := r.agent.Services()
	if err != nil {
		return nil, errors.Wrap(err, `failed to fetch registered services info`)
	}
	result := make([]*core.Service, 0)
	for _, item := range registered {
		service := &core.Service{
			Name:    item.Service,
			Address: item.Address,
			ID:      ptr.String(item.ID),
		}
		if item.Port != 0 {
			service.Port = ptr.Int(item.Port)
		}
		if len(item.Tags) > 0 {
			service.Tags = &item.Tags
		}
		if len(item.Meta) > 0 {
			service.Meta = &item.Meta
		}
		result = append(result, service)
	}
	return result, nil
}

// Deregister make request for Agent.ServiceDeregister by core.Service RegistrationID
func (r *Registry) Deregister(_ context.Context, serviceID string) error {
	if err := r.agent.ServiceDeregister(serviceID); err != nil {
		return errors.Wrapf(err, `failed deregister service by service id "%s"`, serviceID)
	}
	return nil
}

// Register make request for Agent.ServiceRegister for core.Service
func (r *Registry) Register(ctx context.Context, service *core.Service) error {
	if err := service.Validate(ctx); err != nil {
		return errors.Wrap(err, `service has validation error before registration`)
	}
	asr := &api.AgentServiceRegistration{
		Kind:    api.ServiceKindTypical,
		Name:    service.Name,
		Address: service.Address,
	}
	if service.ID != nil {
		asr.ID = *service.ID
	}
	if service.Port != nil {
		asr.Port = *service.Port
	}
	if service.Tags != nil {
		asr.Tags = *service.Tags
	}
	if service.Meta != nil {
		asr.Meta = *service.Meta
	}
	if err := r.agent.ServiceRegister(asr); err != nil {
		return errors.Wrapf(err, `failed register service by service id "%s"`, service.RegistrationID())
	}
	return nil
}
