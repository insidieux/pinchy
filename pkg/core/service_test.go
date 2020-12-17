package core

import (
	"context"
	"testing"

	"github.com/agrea/ptr"
	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestService_Validate(t *testing.T) {
	suite.Run(t, new(serviceValidateTestSuite))
}

func TestService_RegistrationID(t *testing.T) {
	suite.Run(t, new(serviceRegistrationIDTestSuite))
}

func TestServices_IDs(t *testing.T) {
	suite.Run(t, new(servicesIDsTestSuite))
}

func TestServices_Lookup(t *testing.T) {
	suite.Run(t, new(servicesLookupSuite))
}

// --- Suites ---

type serviceValidateTestSuite struct {
	suite.Suite
	service *Service
}

func (s *serviceValidateTestSuite) SetupTest() {
	s.service = new(Service)
}

func (s *serviceValidateTestSuite) TestEmptyName() {
	s.service.Address = `127.0.0.1`
	s.Error(s.service.Validate(context.Background()))
}

func (s *serviceValidateTestSuite) TestEmptyAddress() {
	s.service.Name = `service`
	s.Error(s.service.Validate(context.Background()))
}

func (s *serviceValidateTestSuite) TestValidationPassed() {
	s.service.Name = `service`
	s.service.Address = `127.0.0.1`
	s.NoError(s.service.Validate(context.Background()))
}

type serviceRegistrationIDTestSuite struct {
	suite.Suite
	service *Service
}

func (s *serviceRegistrationIDTestSuite) SetupTest() {
	s.service = new(Service)
}

func (s serviceRegistrationIDTestSuite) TestIDFromName() {
	s.service.Name = `service`
	s.Equal(s.service.Name, s.service.RegistrationID())
}

func (s serviceRegistrationIDTestSuite) TestIDFromID() {
	s.service.Name = `service`
	s.service.ID = ptr.String(`id`)
	s.Equal(*s.service.ID, s.service.RegistrationID())
}

type servicesIDsTestSuite struct {
	suite.Suite
	services Services
}

func (s *servicesIDsTestSuite) SetupTest() {
	s.services = make(Services, 0)
}

func (s *servicesIDsTestSuite) TestEmptyList() {
	s.Empty(s.services.IDs())
}

func (s *servicesIDsTestSuite) TestNonEmptyList() {
	s.services = append(s.services, &Service{Name: `name-1`}, &Service{ID: ptr.String(`id-1`)})
	ids := s.services.IDs()
	s.Len(ids, 2)
	s.Contains(ids, `name-1`)
	s.Contains(ids, `id-1`)
}

type servicesLookupSuite struct {
	suite.Suite
	services Services
}

func (s *servicesLookupSuite) SetupTest() {
	s.services = make(Services, 0)
}

func (s *servicesLookupSuite) TestNilLookup() {
	s.Nil(s.services.Lookup(`id`))
}

func (s *servicesLookupSuite) TestNotNilLookupByName() {
	expected := &Service{
		Name: `name`,
	}
	s.services = append(s.services, expected)
	got := s.services.Lookup(`name`)
	s.NotNil(got)
	s.Equal(expected, got)
}

func (s *servicesLookupSuite) TestNotNilLookupByID() {
	expected := &Service{
		ID: ptr.String(`id`),
	}
	s.services = append(s.services, expected)
	got := s.services.Lookup(`id`)
	s.NotNil(got)
	s.Equal(expected, got)
}
