package core

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
	"github.com/thoas/go-funk"
)

var (
	validation = validator.New()
)

type (
	Service struct {
		Name    string             `json:"," validate:"required"`
		Address string             `json:"," validate:"required"`
		ID      *string            `json:",omitempty"`
		Port    *int               `json:",omitempty"`
		Tags    *[]string          `json:",omitempty"`
		Meta    *map[string]string `json:",omitempty"`
	}
	Services []*Service
)

func (s *Service) Validate(ctx context.Context) error {
	return validation.StructCtx(ctx, s)
}

func (s *Service) RegistrationID() string {
	id := s.Name
	if s.ID != nil {
		id = *s.ID
	}
	return id
}

func (s Services) IDs() []string {
	ids := funk.Map(s, func(service *Service) string {
		return service.RegistrationID()
	})
	return cast.ToStringSlice(ids)
}

func (s Services) Lookup(id string) *Service {
	for _, service := range s {
		if service.RegistrationID() == id {
			return service
		}
	}
	return nil
}
