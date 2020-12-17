package registry

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
	flagPrefix = `registry`
)

type (
	Factory           func(*viper.Viper) (core.Registry, func(), error)
	ProviderInterface interface {
		Name() string
		Flags() *pflag.FlagSet
		Factory() Factory
	}
	ProviderList []ProviderInterface
	provider     struct {
		name    string
		flags   *pflag.FlagSet
		factory Factory
	}
)

var (
	providerList = newProviderList()
)

func newProviderList() *ProviderList {
	return new(ProviderList)
}

func GetProviderList() ProviderList {
	return *providerList
}

func MakeFlagName(name string) string {
	return fmt.Sprintf(`%s.%s`, flagPrefix, name)
}

func Register(name string, flags *pflag.FlagSet, factory Factory) error {
	return providerList.register(&provider{
		name:    name,
		flags:   flags,
		factory: factory,
	})
}

func (p *provider) Name() string {
	return p.name
}

func (p *provider) Flags() *pflag.FlagSet {
	return p.flags
}

func (p *provider) Factory() Factory {
	return p.factory
}

func (pl *ProviderList) register(p ProviderInterface) error {
	if ep, _ := pl.Lookup(p.Name()); ep != nil {
		return errors.Errorf(`registry provider with name "%s" has been already registered`, p.Name())
	}
	var err extension.RegisterError
	p.Flags().VisitAll(func(flag *pflag.Flag) {
		if !strings.HasPrefix(flag.Name, flagPrefix) {
			err = append(err, errors.Errorf(`flag "%s" does not contain required prefix "%s"`, flag.Name, flagPrefix))
		}
	})
	if err != nil {
		return errors.Wrapf(err, `registry "%s" flags validation error`, p.Name())
	}
	*pl = append(*pl, p)
	return nil
}

func (pl ProviderList) Lookup(name string) (ProviderInterface, error) {
	for _, p := range pl {
		if p.Name() == name {
			return p, nil
		}
	}
	return nil, errors.Errorf(`registry provider with name "%s" was not registered`, name)
}

func (pl ProviderList) Get() []ProviderInterface {
	return pl
}
