package core

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestNewManager(t *testing.T) {
	suite.Run(t, new(newManagerTestSuite))
}

func TestManager_Run(t *testing.T) {
	suite.Run(t, new(managerRunTestSuite))
}

func TestManagerError_Add(t *testing.T) {
	suite.Run(t, new(managerErrorAddTestSuite))
}

func TestManagerError_String(t *testing.T) {
	suite.Run(t, new(managerErrorStringTestSuite))
}

func TestManagerError_Error(t *testing.T) {
	suite.Run(t, new(managerErrorErrorTestSuite))
}

func TestManagerError_HasErrors(t *testing.T) {
	suite.Run(t, new(managerErrorHasErrorsTestSuite))
}

// --- Suites ---

type newManagerTestSuite struct {
	suite.Suite
}

func (s *newManagerTestSuite) TestNewManager() {
	s.Equal(
		&Manager{nil, nil, nil, true},
		NewManager(nil, nil, nil, true),
	)
}

type managerRunTestSuite struct {
	suite.Suite
	manager *Manager
}

func (s *managerRunTestSuite) SetupTest() {
	s.manager = new(Manager)
	s.manager.logger, _ = test.NewNullLogger()
}

func (s *managerRunTestSuite) TestErrorFetchFromSource() {
	ctx := context.Background()
	sourceMock := new(MockSource)
	sourceMock.On(`Fetch`, ctx).Return(nil, errors.New(`expected error`))

	s.manager.source = sourceMock
	err := s.manager.Run(ctx)
	s.Error(err)
	s.EqualError(err, `failed to fetch services from source: expected error`)
}

func (s *managerRunTestSuite) TestErrorFetchFromRegistry() {
	ctx := context.Background()
	sourceMock := new(MockSource)
	sourceMock.On(`Fetch`, ctx).Return(Services{}, nil)
	registryMock := new(MockRegistry)
	registryMock.On(`Fetch`, ctx).Return(nil, errors.New(`expected error`))

	s.manager.source = sourceMock
	s.manager.registry = registryMock

	err := s.manager.Run(ctx)
	s.Error(err)
	s.EqualError(err, `failed to fetch services from registry: expected error`)
}

func (s *managerRunTestSuite) TestErrorDeregisterOrphan() {
	ctx := context.Background()
	sourceMock := new(MockSource)
	sourceMock.On(`Fetch`, ctx).Return(Services{{Name: `service-1`}}, nil)
	registryMock := new(MockRegistry)
	registryMock.On(`Fetch`, ctx).Return(Services{{Name: `service-2`}}, nil)
	registryMock.On(`Deregister`, ctx, `service-2`).Return(errors.New(`expected error`))

	s.manager.source = sourceMock
	s.manager.registry = registryMock
	s.manager.exitOnError = true

	err := s.manager.Run(ctx)
	s.Error(err)
	s.EqualError(err, `failed to deregister services: failed to deregister service "service-2" from registry: expected error`)
}

func (s *managerRunTestSuite) TestErrorRegister() {
	ctx := context.Background()
	service := &Service{Name: `service-1`}
	sourceMock := new(MockSource)
	sourceMock.On(`Fetch`, ctx).Return(Services{service}, nil)
	registryMock := new(MockRegistry)
	registryMock.On(`Fetch`, ctx).Return(Services{}, nil)
	registryMock.On(`Register`, ctx, service).Return(errors.New(`expected error`))

	s.manager.source = sourceMock
	s.manager.registry = registryMock
	s.manager.exitOnError = true

	err := s.manager.Run(ctx)
	s.Error(err)
	s.EqualError(err, `failed to register services: failed to register service "service-1" in registry: expected error`)
}

func (s *managerRunTestSuite) TestSuccess() {
	ctx := context.Background()
	sourceMock := new(MockSource)
	sourceMock.On(`Fetch`, ctx).Return(Services{{Name: `service-1`}}, nil)
	registryMock := new(MockRegistry)
	registryMock.On(`Fetch`, ctx).Return(Services{{Name: `service-2`}}, nil)
	registryMock.On(`Deregister`, ctx, `service-2`).Return(nil)
	registryMock.On(`Register`, ctx, mock.Anything).Return(nil)

	s.manager.source = sourceMock
	s.manager.registry = registryMock
	s.NoError(s.manager.Run(ctx))
}

type managerErrorAddTestSuite struct {
	suite.Suite
	err managerError
}

func (s *managerErrorAddTestSuite) TestAdd() {
	s.err = append(s.err, errors.New(`expected error`))
	s.Len(s.err, 1)
	s.Equal(`expected error`, s.err[0].Error())
}

type managerErrorStringTestSuite struct {
	suite.Suite
	err managerError
}

func (s *managerErrorStringTestSuite) TestEmptyString() {
	s.Empty(s.err.String())
}

func (s *managerErrorStringTestSuite) TestNonEmptyString() {
	s.err = append(s.err, errors.New(`expected error 1`))
	s.err = append(s.err, errors.New(`expected error 2`))
	s.NotEmpty(s.err.String())
	s.Equal(`expected error 1; expected error 2`, s.err.String())
}

type managerErrorErrorTestSuite struct {
	suite.Suite
	err managerError
}

func (s *managerErrorErrorTestSuite) TestEmptyError() {
	s.Empty(s.err.Error())
}

func (s *managerErrorErrorTestSuite) TestNonEmptyError() {
	s.err = append(s.err, errors.New(`expected error 1`))
	s.err = append(s.err, errors.New(`expected error 2`))
	s.NotEmpty(s.err.Error())
	s.Equal(`expected error 1; expected error 2`, s.err.Error())
}

type managerErrorHasErrorsTestSuite struct {
	suite.Suite
	err managerError
}

func (s *managerErrorHasErrorsTestSuite) TestEmptyError() {
	s.False(s.err.HasErrors())
}

func (s *managerErrorHasErrorsTestSuite) TestNonEmptyError() {
	s.err = append(s.err, errors.New(`expected error 1`))
	s.err = append(s.err, errors.New(`expected error 2`))
	s.True(s.err.HasErrors())
}

// --- Mocks ---

// MockRegistry is an autogenerated mock type for the Registry type
type MockRegistry struct {
	mock.Mock
}

// Deregister provides a mock function with given fields: ctx, serviceID
func (_m *MockRegistry) Deregister(ctx context.Context, serviceID string) error {
	ret := _m.Called(ctx, serviceID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, serviceID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Fetch provides a mock function with given fields: ctx
func (_m *MockRegistry) Fetch(ctx context.Context) (Services, error) {
	ret := _m.Called(ctx)

	var r0 Services
	if rf, ok := ret.Get(0).(func(context.Context) Services); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Services)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: ctx, service
func (_m *MockRegistry) Register(ctx context.Context, service *Service) error {
	ret := _m.Called(ctx, service)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *Service) error); ok {
		r0 = rf(ctx, service)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSource is an autogenerated mock type for the Source type
type MockSource struct {
	mock.Mock
}

// Fetch provides a mock function with given fields: ctx
func (_m *MockSource) Fetch(ctx context.Context) (Services, error) {
	ret := _m.Called(ctx)

	var r0 Services
	if rf, ok := ret.Get(0).(func(context.Context) Services); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Services)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
