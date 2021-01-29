package source

import (
	"fmt"
	"strings"

	"github.com/insidieux/pinchy/internal/extension"
	"github.com/insidieux/pinchy/pkg/core"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagPrefix = `source`
)

type (
	// Factory is helper function to provide core.Source implementation.
	// Function MUST return (core.Source, func(), error) for further call with google.wire
	Factory func(*viper.Viper) (core.Source, func(), error)

	// ProviderInterface contains main information about extension.
	// Name is used for unique registration and validation
	// Flags used for provide cobra.Command flags
	// Factory used in generation main command source code with google.wire
	ProviderInterface interface {
		Name() string
		Flags() *pflag.FlagSet
		Factory() Factory
		Deprecated() bool
	}

	// ProviderList is helper custom type for handle registration and lookup for slice of ProviderInterface
	ProviderList []ProviderInterface

	provider struct {
		name       string
		flags      *pflag.FlagSet
		factory    Factory
		deprecated bool
	}
)

var (
	providerList = newProviderList()
)

// newProviderList return new empty ProviderList
func newProviderList() *ProviderList {
	return new(ProviderList)
}

// GetProviderList return internal ProviderList with pre-registered slice of ProviderInterface
func GetProviderList() ProviderList {
	return *providerList
}

// MakeFlagName is helper for generation valid flag name
func MakeFlagName(name string) string {
	return fmt.Sprintf(`%s.%s`, flagPrefix, name)
}

// Register new ProviderInterface implementation by name, flags and factory
// Register return error if ProviderInterface already registered with same name
// Register return error if ProviderInterface provide flags with incorrect name prefix
func Register(name string, flags *pflag.FlagSet, factory Factory, deprecated bool) error {
	return providerList.register(&provider{
		name:       name,
		flags:      flags,
		factory:    factory,
		deprecated: deprecated,
	})
}

// Name return ProviderInterface implementation name
func (p *provider) Name() string {
	return p.name
}

// Flags return ProviderInterface implementation flags
func (p *provider) Flags() *pflag.FlagSet {
	return p.flags
}

// Factory return ProviderInterface implementation factory
func (p *provider) Factory() Factory {
	return p.factory
}

// Deprecated return ProviderInterface implementation deprecated flag
func (p *provider) Deprecated() bool {
	return p.deprecated
}

func (pl *ProviderList) register(p ProviderInterface) error {
	if ep, _ := pl.Lookup(p.Name()); ep != nil {
		return errors.Errorf(`source provider with name "%s" has been already registered`, p.Name())
	}
	var err extension.RegisterError
	p.Flags().VisitAll(func(flag *pflag.Flag) {
		if !strings.HasPrefix(flag.Name, flagPrefix) {
			err = append(err, errors.Errorf(`flag "%s" does not contain required prefix "%s"`, flag.Name, flagPrefix))
		}
	})
	if err != nil {
		return errors.Wrapf(err, `source "%s" flags validation error`, p.Name())
	}
	*pl = append(*pl, p)
	return nil
}

// Lookup return registered ProviderInterface by name
// Lookup return error if ProviderInterface was not registered with called name
func (pl ProviderList) Lookup(name string) (ProviderInterface, error) {
	for _, p := range pl {
		if p.Name() == name {
			return p, nil
		}
	}
	return nil, errors.Errorf(`source provider with name "%s" was not registered`, name)
}

// Get return slice of registered ProviderInterface's
func (pl ProviderList) Get() []ProviderInterface {
	return pl
}
