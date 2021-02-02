package catalog

import (
	"context"
	"fmt"

	"github.com/agrea/ptr"
	"github.com/hashicorp/consul/api"
	"github.com/insidieux/pinchy/pkg/core"
	"github.com/insidieux/pinchy/pkg/core/registry/consul"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type (
	// Catalog interface provide common function for work with Consul HTTP API /v1/catalog
	Catalog interface {
		Services(*api.QueryOptions) (map[string][]string, *api.QueryMeta, error)
		Service(string, string, *api.QueryOptions) ([]*api.CatalogService, *api.QueryMeta, error)
		Deregister(*api.CatalogDeregistration, *api.WriteOptions) (*api.WriteMeta, error)
		Register(*api.CatalogRegistration, *api.WriteOptions) (*api.WriteMeta, error)
	}

	// Registry is implementation of core.Registry interface
	Registry struct {
		catalog Catalog
		logger  core.LoggerInterface
		tag     consul.Tag
	}
)

// NewRegistry provide Registry as core.Registry implementation
func NewRegistry(catalog Catalog, tag consul.Tag) *Registry {
	return &Registry{
		catalog: catalog,
		tag:     tag,
	}
}

// Fetch make request for Catalog.Services plus Catalog.Service and try to cast result to core.Services
func (r *Registry) Fetch(ctx context.Context) (core.Services, error) {
	r.logger.Infoln(`Fetch registered services from catalog`)
	query := &api.QueryOptions{
		Filter: fmt.Sprintf(`("%s" in Tags)`, r.tag),
	}
	query = query.WithContext(ctx)
	names, _, err := r.catalog.Services(query)
	if err != nil {
		return nil, errors.Wrap(err, `failed to fetch registered services info`)
	}

	r.logger.Infoln(`Prepare registered services list`)
	result := make([]*core.Service, 0)
	for name := range names {
		opts := &api.QueryOptions{}
		opts = opts.WithContext(ctx)
		items, _, err := r.catalog.Service(name, string(r.tag), opts)
		if err != nil {
			return nil, errors.Wrap(err, `failed to fetch registered service info`)
		}

		for _, item := range items {
			service := &core.Service{
				Name:    item.ServiceName,
				Address: item.ServiceAddress,
				ID:      ptr.String(item.ServiceID),
			}
			if item.ServicePort != 0 {
				service.Port = ptr.Int(item.ServicePort)
			}
			if len(item.ServiceTags) > 0 {
				service.Tags = &item.ServiceTags
			}
			if len(item.ServiceMeta) > 0 {
				service.Meta = &item.ServiceMeta
			}
			service.Node = &core.Node{
				Node:       item.Node,
				Address:    item.Address,
				Datacenter: ptr.String(item.Datacenter),
				NodeMeta:   &item.NodeMeta,
			}
			result = append(result, service)
		}

	}
	return result, nil
}

// Deregister make request for Catalog.Deregister by core.Service RegistrationID
func (r *Registry) Deregister(ctx context.Context, service *core.Service) error {
	r.logger.Infof(`Validate service "%s"`, service.RegistrationID())
	if err := service.Validate(ctx, r.validateService); err != nil {
		return errors.Wrap(err, `service has validation error before deregister`)
	}

	r.logger.Infof(`Send service deregister catalog request for service "%s"`, service.RegistrationID())
	opts := &api.WriteOptions{}
	opts = opts.WithContext(ctx)
	_, err := r.catalog.Deregister(
		&api.CatalogDeregistration{
			ServiceID: service.RegistrationID(),
			Node:      service.Node.Node,
		},
		opts,
	)

	if err != nil {
		return errors.Wrapf(err, `failed deregister service by service id "%s"`, service.RegistrationID())
	}
	return nil
}

// Register make request for Catalog.Register for core.Service
func (r *Registry) Register(ctx context.Context, service *core.Service) error {
	r.logger.Infof(`Validate service "%s"`, service.RegistrationID())
	if err := service.Validate(ctx, r.validateService); err != nil {
		return errors.Wrap(err, `service has validation error before registration`)
	}

	cr := &api.CatalogRegistration{
		Node:    service.Node.Node,
		Address: service.Node.Address,
		Service: &api.AgentService{
			Kind:              api.ServiceKindTypical,
			Service:           service.Name,
			Address:           service.Address,
			Weights:           api.AgentWeights{},
			EnableTagOverride: true,
		},
	}
	if service.ID != nil {
		cr.Service.ID = *service.ID
	}
	if service.Port != nil {
		cr.Service.Port = *service.Port
	}
	cr.Service.Tags = append(cr.Service.Tags, string(r.tag))
	if service.Tags != nil {
		cr.Service.Tags = append(cr.Service.Tags, *service.Tags...)
	}
	cr.Service.Tags = funk.UniqString(cr.Service.Tags)
	if service.Meta != nil {
		cr.Service.Meta = *service.Meta
	}

	r.logger.Infof(`Send service register catalog request for service "%s"`, service.RegistrationID())
	opts := &api.WriteOptions{}
	opts = opts.WithContext(ctx)
	if _, err := r.catalog.Register(cr, opts); err != nil {
		return errors.Wrapf(err, `failed register service by service id "%s"`, service.RegistrationID())
	}
	return nil
}

// WithLogger is implementation of core.Loggable interface
func (r *Registry) WithLogger(logger core.LoggerInterface) {
	r.logger = logger
}

// validateService is implementation of core.ValidationFunc func
func (r *Registry) validateService(_ context.Context, service *core.Service) error {
	if service.Node == nil {
		return errors.New(`service field "Node" is required and cannot be empty`)
	}
	if service.Node.Node == `` {
		return errors.New(`service field "Node.Node" is required and cannot be empty`)
	}
	if service.Node.Address == `` {
		return errors.Errorf(`service field "Node.Address" is required and cannot be empty`)
	}
	return nil
}
