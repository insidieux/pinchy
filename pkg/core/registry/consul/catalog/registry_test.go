package catalog

import (
	"context"
	"testing"

	"github.com/agrea/ptr"
	"github.com/hashicorp/consul/api"
	"github.com/insidieux/pinchy/pkg/core"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestNewRegistry(t *testing.T) {
	suite.Run(t, new(newRegistryTestSuite))
}

func TestRegistry_Fetch(t *testing.T) {
	suite.Run(t, new(registryFetchTestSuite))
}

func TestRegistry_Deregister(t *testing.T) {
	suite.Run(t, new(registryDeregisterTestSuite))
}

func TestRegistry_Register(t *testing.T) {
	suite.Run(t, new(registryRegisterTestSuite))
}

func TestRegistry_WithLogger(t *testing.T) {
	suite.Run(t, new(registryWithLoggerTestSuite))
}

// --- Suites ---

type newRegistryTestSuite struct {
	suite.Suite
}

func (s *newRegistryTestSuite) TestNewRegistry() {
	got := NewRegistry(nil, ``)
	s.Implements((*core.Registry)(nil), got)
	s.Equal(&Registry{nil, nil, ``}, got)
}

type registryFetchTestSuite struct {
	suite.Suite
	catalog  *MockCatalog
	registry *Registry
}

func (s *registryFetchTestSuite) SetupTest() {
	s.catalog = new(MockCatalog)
	s.registry = NewRegistry(s.catalog, `test`)
	s.registry.logger, _ = test.NewNullLogger()
}

func (s *registryFetchTestSuite) TestErrorCatalogServicesFetch() {
	s.catalog.On(`Services`, mock.Anything).Return(nil, nil, errors.New(`expected error`))

	s.registry.catalog = s.catalog
	services, err := s.registry.Fetch(context.Background())
	s.Nil(services)
	s.EqualError(err, `failed to fetch registered services info: expected error`)
}

func (s *registryFetchTestSuite) TestErrorCatalogServiceFetch() {
	s.catalog.On(`Services`, mock.Anything).Return(map[string][]string{`name`: nil}, nil, nil)
	s.catalog.On(`Service`, `name`, mock.Anything, mock.Anything).Return(nil, nil, errors.New(`expected error`))

	s.registry.catalog = s.catalog
	services, err := s.registry.Fetch(context.Background())
	s.Nil(services)
	s.EqualError(err, `failed to fetch registered service info: expected error`)
}

func (s *registryFetchTestSuite) TestSuccess() {
	expectedTags := []string{`tags`}
	expectedMeta := map[string]string{`key`: `value`}

	s.catalog.On(`Services`, mock.Anything).Return(map[string][]string{`name`: nil}, nil, nil)
	s.catalog.On(`Service`, `name`, mock.Anything, mock.Anything).Return([]*api.CatalogService{
		{
			ServiceName:    `name`,
			ServiceAddress: `127.0.0.1`,
			ServiceID:      `id`,
			ServiceTags:    expectedTags,
			ServiceMeta:    expectedMeta,
			ServicePort:    80,
			Node:           `node-1`,
			Datacenter:     `dc-1`,
			Address:        `127.0.0.1`,
			NodeMeta:       expectedMeta,
		},
	}, nil, nil)

	fetchedServices, err := s.registry.Fetch(context.Background())

	s.NoError(err)
	s.Equal(core.Services{
		&core.Service{
			Name:    `name`,
			Address: `127.0.0.1`,
			ID:      ptr.String(`id`),
			Tags:    &expectedTags,
			Meta:    &expectedMeta,
			Port:    ptr.Int(80),
			Node: &core.Node{
				Node:       `node-1`,
				Address:    `127.0.0.1`,
				Datacenter: ptr.String(`dc-1`),
				NodeMeta:   &expectedMeta,
			},
		},
	}, fetchedServices)

}

type registryDeregisterTestSuite struct {
	suite.Suite
	catalog  *MockCatalog
	registry *Registry
	service  *core.Service
}

func (s *registryDeregisterTestSuite) SetupTest() {
	s.catalog = new(MockCatalog)
	s.registry = NewRegistry(s.catalog, `test`)
	s.registry.logger, _ = test.NewNullLogger()
	s.service = &core.Service{
		Name:    `service`,
		Address: `127.0.0.1`,
		Node: &core.Node{
			Node:    `node-1`,
			Address: `127.0.0.1`,
		},
	}
}

func (s *registryDeregisterTestSuite) TestErrorServiceValidation() {
	err := s.registry.Deregister(context.Background(), &core.Service{
		Name: s.service.Name,
	})
	s.Error(err)
	s.Contains(err.Error(), `service has validation error before deregister`)
}

func (s *registryDeregisterTestSuite) TestErrorServiceCustomValidation() {
	var err error
	err = s.registry.Deregister(context.Background(), &core.Service{
		Name:    s.service.Name,
		Address: s.service.Address,
	})
	s.Error(err)
	s.Contains(err.Error(), `service field "Node" is required and cannot be empty`)

	err = s.registry.Deregister(context.Background(), &core.Service{
		Name:    s.service.Name,
		Address: s.service.Address,
		Node:    &core.Node{},
	})
	s.Error(err)
	s.Contains(err.Error(), `service field "Node.Node" is required and cannot be empty`)

	err = s.registry.Deregister(context.Background(), &core.Service{
		Name:    s.service.Name,
		Address: s.service.Address,
		Node: &core.Node{
			Node: s.service.Node.Node,
		},
	})
	s.Error(err)
	s.Contains(err.Error(), `service field "Node.Address" is required and cannot be empty`)
}

func (s *registryDeregisterTestSuite) TestErrorCatalogDeregister() {
	s.catalog.On(`Deregister`, mock.Anything, mock.Anything).Return(nil, errors.New(`expected error`))

	err := s.registry.Deregister(context.Background(), s.service)
	s.EqualError(err, `failed deregister service by service id "service": expected error`)
}

func (s *registryDeregisterTestSuite) TestSuccess() {
	s.catalog.On(`Deregister`, &api.CatalogDeregistration{
		Node:      s.service.Node.Node,
		ServiceID: s.service.RegistrationID(),
	}, mock.Anything).Return(nil, nil)

	err := s.registry.Deregister(context.Background(), s.service)
	s.NoError(err)
}

type registryRegisterTestSuite struct {
	suite.Suite
	catalog  *MockCatalog
	registry *Registry
}

func (s *registryRegisterTestSuite) SetupTest() {
	s.catalog = new(MockCatalog)
	s.registry = NewRegistry(s.catalog, `test`)
	s.registry.logger, _ = test.NewNullLogger()
}

func (s *registryRegisterTestSuite) TestErrorServiceValidation() {
	err := s.registry.Register(context.Background(), &core.Service{
		Name: `name`,
	})
	s.Error(err)
	s.Contains(err.Error(), `service has validation error before registration`)
}

func (s *registryRegisterTestSuite) TestErrorCatalogRegister() {
	s.catalog.On(`Register`, mock.Anything, mock.Anything).Return(nil, errors.New(`expected error`))

	err := s.registry.Register(context.Background(), &core.Service{
		Name:    `name`,
		Address: `127.0.0.1`,
		Node: &core.Node{
			Node:    `node-1`,
			Address: `127.0.0.1`,
		},
	})
	s.EqualError(err, `failed register service by service id "name": expected error`)
}

func (s *registryRegisterTestSuite) TestSuccess() {
	s.catalog.On(`Register`, mock.Anything, mock.Anything).Return(nil, nil)

	expectedTags := []string{`tags`}
	expectedMeta := map[string]string{`key`: `value`}
	err := s.registry.Register(context.Background(), &core.Service{
		Name:    `name`,
		Address: `127.0.0.1`,
		ID:      ptr.String(`id`),
		Tags:    &expectedTags,
		Meta:    &expectedMeta,
		Port:    ptr.Int(80),
		Node: &core.Node{
			Node:    `node-1`,
			Address: `127.0.0.1`,
		},
	})
	s.NoError(err)
}

type registryWithLoggerTestSuite struct {
	suite.Suite
}

func (s *registryWithLoggerTestSuite) TestWithLogger() {
	logger, _ := test.NewNullLogger()
	src := NewRegistry(nil, ``)
	src.WithLogger(logger)
}

// --- Mocks ---

// MockCatalog is an autogenerated mock type for the Catalog type
type MockCatalog struct {
	mock.Mock
}

// Deregister provides a mock function with given fields: _a0, _a1
func (_m *MockCatalog) Deregister(_a0 *api.CatalogDeregistration, _a1 *api.WriteOptions) (*api.WriteMeta, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *api.WriteMeta
	if rf, ok := ret.Get(0).(func(*api.CatalogDeregistration, *api.WriteOptions) *api.WriteMeta); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.WriteMeta)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*api.CatalogDeregistration, *api.WriteOptions) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: _a0, _a1
func (_m *MockCatalog) Register(_a0 *api.CatalogRegistration, _a1 *api.WriteOptions) (*api.WriteMeta, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *api.WriteMeta
	if rf, ok := ret.Get(0).(func(*api.CatalogRegistration, *api.WriteOptions) *api.WriteMeta); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.WriteMeta)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*api.CatalogRegistration, *api.WriteOptions) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Service provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockCatalog) Service(_a0 string, _a1 string, _a2 *api.QueryOptions) ([]*api.CatalogService, *api.QueryMeta, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []*api.CatalogService
	if rf, ok := ret.Get(0).(func(string, string, *api.QueryOptions) []*api.CatalogService); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*api.CatalogService)
		}
	}

	var r1 *api.QueryMeta
	if rf, ok := ret.Get(1).(func(string, string, *api.QueryOptions) *api.QueryMeta); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*api.QueryMeta)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string, *api.QueryOptions) error); ok {
		r2 = rf(_a0, _a1, _a2)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Services provides a mock function with given fields: _a0
func (_m *MockCatalog) Services(_a0 *api.QueryOptions) (map[string][]string, *api.QueryMeta, error) {
	ret := _m.Called(_a0)

	var r0 map[string][]string
	if rf, ok := ret.Get(0).(func(*api.QueryOptions) map[string][]string); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]string)
		}
	}

	var r1 *api.QueryMeta
	if rf, ok := ret.Get(1).(func(*api.QueryOptions) *api.QueryMeta); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*api.QueryMeta)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*api.QueryOptions) error); ok {
		r2 = rf(_a0)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
