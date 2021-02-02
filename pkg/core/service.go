package core

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/thoas/go-funk"
)

type (
	// Node contains info about host/node/server, hosting service. Used for catalog registration in consul.
	Node struct {
		Node       string             `json:","`
		Address    string             `json:","`
		Datacenter *string            `json:",omitempty"`
		NodeMeta   *map[string]string `json:",omitempty"`
	}

	// Service contains all the necessary information for further registration in Registry
	Service struct {
		Name    string             `json:","`
		Address string             `json:","`
		ID      *string            `json:",omitempty"`
		Port    *int               `json:",omitempty"`
		Tags    *[]string          `json:",omitempty"`
		Meta    *map[string]string `json:",omitempty"`
		Node    *Node              `json:",omitempty"`
	}

	// Services is simple helper for hold slice of Service's
	Services []*Service

	// ValidationFunc is a func providing additional validation, called by Source or Registry
	ValidationFunc func(context.Context, *Service) error
)

// Validate process validation to check required fields for Service, such as Service.Name and Service.Address
// Also there is a possibility to pass your own additional checks
func (s *Service) Validate(ctx context.Context, checks ...ValidationFunc) error {
	if s.Name == `` {
		return errors.New(`service field "name" is required and cannot be empty`)
	}
	if s.Address == `` {
		return errors.Errorf(`service "%s" field "address" is required and cannot be empty`, s.Name)
	}
	for _, check := range checks {
		if err := check(ctx, s); err != nil {
			return errors.Wrapf(err, `service "%s" custom check failed`, s.Name)
		}
	}
	return nil
}

// RegistrationID generate identification for registration in Registry.
func (s *Service) RegistrationID() string {
	id := s.Name
	if s.ID != nil {
		id = *s.ID
	}
	return id
}

// IDs return slice of Service.RegistrationID
func (s Services) IDs() []string {
	ids := funk.Map(s, func(service *Service) string {
		return service.RegistrationID()
	})
	return cast.ToStringSlice(ids)
}

// Lookup return Service by Service.RegistrationID, if found
func (s Services) Lookup(id string) *Service {
	for _, service := range s {
		if service.RegistrationID() == id {
			return service
		}
	}
	return nil
}
