package consul

import (
	pkgConsul "github.com/insidieux/pinchy/pkg/core/registry/consul"

	"github.com/hashicorp/consul/api"
	"github.com/insidieux/pinchy/internal/extension/registry"
	"github.com/insidieux/pinchy/pkg/core/registry/consul/agent"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	registryName        = `consul`
	registryAgentName   = `consul-agent`
	registryCatalogName = `consul-catalog`

	flagConsulAddress = `address`
	flagTag           = `tag`
)

type (
	client interface {
		Agent() *api.Agent
		Catalog() *api.Catalog
	}
	factory func(*api.Config) (*api.Client, error)
)

func init() {
	set := pflag.NewFlagSet(registryName, pflag.ExitOnError)
	set.String(registry.MakeFlagName(flagConsulAddress), `127.0.0.1:8500`, `Consul http api address`)
	set.String(registry.MakeFlagName(flagTag), `pinchy`, `Common service tag added for all registered service`)
	// register deprecated consul agent registry
	if err := registry.Register(registryName, set, NewAgentRegistry, true); err != nil {
		panic(err)
	}
	// register new consul agent registry
	if err := registry.Register(registryAgentName, set, NewAgentRegistry, false); err != nil {
		panic(err)
	}
	// register new consul catalog registry
	if err := registry.Register(registryCatalogName, set, NewCatalogRegistry, false); err != nil {
		panic(err)
	}
}

func provideTag(v *viper.Viper) (pkgConsul.Tag, error) {
	flag := registry.MakeFlagName(flagTag)
	tag := v.GetString(flag)
	if tag == `` {
		return ``, errors.Errorf(`Flag "%s" is required`, flag)
	}
	return pkgConsul.Tag(tag), nil
}

func provideClientConfig(v *viper.Viper) (*api.Config, error) {
	flag := registry.MakeFlagName(flagConsulAddress)
	address := v.GetString(flag)
	if address == `` {
		return nil, errors.Errorf(`Flag "%s" is required`, flag)
	}

	cfg := api.DefaultConfig()
	cfg.Address = address
	return cfg, nil
}

func provideConsulClientFactory() factory {
	return api.NewClient
}

func provideClient(cfg *api.Config, factory factory) (client, error) {
	c, err := factory(cfg)
	if err != nil {
		return nil, errors.Wrap(err, `failed to create consul client`)
	}
	return c, nil
}

func provideAgent(c client) agent.Agent {
	return c.Agent()
}

func provideCatalog(c client) *api.Catalog {
	return c.Catalog()
}
