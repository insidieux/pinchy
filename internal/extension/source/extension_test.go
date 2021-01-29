package source

import (
	"fmt"
	"testing"

	"github.com/insidieux/pinchy/pkg/core"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func Test_newProviderList(t *testing.T) {
	suite.Run(t, new(newProviderListTestSuite))
}

func TestGetProviderList(t *testing.T) {
	suite.Run(t, new(getProviderListTestSuite))
}

func TestMakeFlagName(t *testing.T) {
	suite.Run(t, new(makeFlagNameTestSuite))
}

func TestProviderList_Get(t *testing.T) {
	suite.Run(t, new(providerListGetTestSuite))
}

func TestProviderList_Lookup(t *testing.T) {
	suite.Run(t, new(providerListLookupTestSuite))
}

func TestProviderList_register(t *testing.T) {
	suite.Run(t, new(providerListRegisterTestSuite))
}

func TestProvider_Factory(t *testing.T) {
	suite.Run(t, new(providerFactoryTestSuite))
}

func TestProvider_Flags(t *testing.T) {
	suite.Run(t, new(providerFlagsTestSuite))
}

func TestProvider_Name(t *testing.T) {
	suite.Run(t, new(providerNameTestSuite))
}

func TestProvider_Deprecated(t *testing.T) {
	suite.Run(t, new(providerDeprecatedTestSuite))
}

func TestRegister(t *testing.T) {
	suite.Run(t, new(registerTestSuite))
}

// --- Suites ---

type newProviderListTestSuite struct {
	suite.Suite
}

func (s *newProviderListTestSuite) TestNewProviderList() {
	got := newProviderList()
	s.Equal(new(ProviderList), got)
}

type getProviderListTestSuite struct {
	suite.Suite
}

func (s *newProviderListTestSuite) TestGetProviderList() {
	got := GetProviderList()
	s.Equal(*newProviderList(), got)
}

type makeFlagNameTestSuite struct {
	suite.Suite
}

func (s *makeFlagNameTestSuite) TestMakeFlagName() {
	s.Equal(fmt.Sprintf(`%s.%s`, flagPrefix, `flag`), MakeFlagName(`flag`))
}

type providerListGetTestSuite struct {
	suite.Suite
}

func (s *providerListGetTestSuite) TestGet() {
	s.Equal(*new([]ProviderInterface), newProviderList().Get())
}

type providerListLookupTestSuite struct {
	suite.Suite
	list *ProviderList
}

func (s *providerListLookupTestSuite) SetupTest() {
	s.list = newProviderList()
}

func (s *providerListLookupTestSuite) TestErrorProviderWasNotRegistered() {
	provider, err := s.list.Lookup(`provider`)
	s.Nil(provider)
	s.Error(err)
	s.EqualError(err, `source provider with name "provider" was not registered`)
}

func (s *providerListLookupTestSuite) TestSuccess() {
	registeredProvider := new(MockProviderInterface)
	registeredProvider.On(`Name`).Return(`provider`)
	*s.list = append(*s.list, registeredProvider)

	provider, err := s.list.Lookup(`provider`)
	s.NoError(err)
	s.NotNil(provider)
	s.Equal(registeredProvider, provider)
}

type providerListRegisterTestSuite struct {
	suite.Suite
	list     *ProviderList
	provider *MockProviderInterface
}

func (s *providerListRegisterTestSuite) SetupTest() {
	s.list = newProviderList()
	s.provider = new(MockProviderInterface)
	s.provider.On(`Name`).Return(`provider`)
}

func (s *providerListRegisterTestSuite) TestErrorProviderHasBeenRegistered() {
	beforeRegisterProvider := new(MockProviderInterface)
	beforeRegisterProvider.On(`Name`).Return(`provider`)
	*s.list = append(*s.list, beforeRegisterProvider)

	err := s.list.register(s.provider)
	s.Error(err)
	s.EqualError(err, `source provider with name "provider" has been already registered`)
}

func (s *providerListRegisterTestSuite) TestErrorProviderHasFlagsWithoutRequiredPrefix() {
	flags := pflag.NewFlagSet(`provider`, pflag.ExitOnError)
	flags.String(`wrong.flag`, `value`, `usage`)
	s.provider.On(`Flags`).Return(flags)

	err := s.list.register(s.provider)
	s.Error(err)
	s.EqualError(
		err,
		fmt.Sprintf(
			`source "%s" flags validation error: flag "%s" does not contain required prefix "%s"`,
			`provider`,
			`wrong.flag`,
			flagPrefix,
		),
	)
}

func (s *providerListRegisterTestSuite) TestSuccess() {
	flags := pflag.NewFlagSet(`provider`, pflag.ExitOnError)
	flags.String(`source.flag`, `value`, `usage`)
	s.provider.On(`Flags`).Return(flags)

	err := s.list.register(s.provider)
	s.NoError(err)
}

type providerFactoryTestSuite struct {
	suite.Suite
}

func (s *providerFactoryTestSuite) TestSuccess() {
	factory := new(MockFactory)
	p := &provider{
		factory: factory.Execute,
	}
	s.IsType(*new(Factory), p.Factory())
}

type providerFlagsTestSuite struct {
	suite.Suite
}

func (s *providerFlagsTestSuite) TestSuccess() {
	p := &provider{
		flags: pflag.CommandLine,
	}
	s.Equal(pflag.CommandLine, p.Flags())
}

type providerNameTestSuite struct {
	suite.Suite
}

func (s *providerNameTestSuite) TestSuccess() {
	p := &provider{
		name: `name`,
	}
	s.Equal(`name`, p.Name())
}

type providerDeprecatedTestSuite struct {
	suite.Suite
}

func (s *providerDeprecatedTestSuite) TestSuccess() {
	p := &provider{
		deprecated: false,
	}
	s.False(p.Deprecated())
}

type registerTestSuite struct {
	suite.Suite
	provider *MockProviderInterface
}

func (s *registerTestSuite) SetupTest() {
	providerList = newProviderList()
}

func (s *registerTestSuite) TestErrorProviderHasBeenRegistered() {
	*providerList = append(*providerList, &provider{name: `provider`})

	err := Register(`provider`, nil, nil, false)
	s.Error(err)
	s.EqualError(err, `source provider with name "provider" has been already registered`)
}

func (s *registerTestSuite) TestErrorProviderHasFlagsWithoutRequiredPrefix() {
	flags := pflag.NewFlagSet(`provider`, pflag.ExitOnError)
	flags.String(`wrong.flag`, `value`, `usage`)

	err := Register(`provider`, flags, nil, false)
	s.Error(err)
	s.EqualError(
		err,
		fmt.Sprintf(
			`source "%s" flags validation error: flag "%s" does not contain required prefix "%s"`,
			`provider`,
			`wrong.flag`,
			flagPrefix,
		),
	)
}

func (s *registerTestSuite) TestSuccess() {
	flags := pflag.NewFlagSet(`provider`, pflag.ExitOnError)
	flags.String(`source.flag`, `value`, `usage`)

	err := Register(`provider`, flags, nil, false)
	s.NoError(err)
}

// --- Mocks ---

// MockFactory is an autogenerated mock type for the Factory type
type MockFactory struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *MockFactory) Execute(_a0 *viper.Viper) (core.Source, func(), error) {
	ret := _m.Called(_a0)

	var r0 core.Source
	if rf, ok := ret.Get(0).(func(*viper.Viper) core.Source); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.Source)
		}
	}

	var r1 func()
	if rf, ok := ret.Get(1).(func(*viper.Viper) func()); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(func())
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*viper.Viper) error); ok {
		r2 = rf(_a0)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockProviderInterface is an autogenerated mock type for the ProviderInterface type
type MockProviderInterface struct {
	mock.Mock
}

// Deprecated provides a mock function with given fields:
func (_m *MockProviderInterface) Deprecated() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Factory provides a mock function with given fields:
func (_m *MockProviderInterface) Factory() Factory {
	ret := _m.Called()

	var r0 Factory
	if rf, ok := ret.Get(0).(func() Factory); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Factory)
		}
	}

	return r0
}

// Flags provides a mock function with given fields:
func (_m *MockProviderInterface) Flags() *pflag.FlagSet {
	ret := _m.Called()

	var r0 *pflag.FlagSet
	if rf, ok := ret.Get(0).(func() *pflag.FlagSet); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pflag.FlagSet)
		}
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *MockProviderInterface) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
