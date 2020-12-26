package consul

import (
	"net/http"

	pkgConsul "github.com/insidieux/pinchy/pkg/core/registry/consul"

	"github.com/hashicorp/consul/api"
	"github.com/insidieux/pinchy/internal/extension/registry"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	registryName = `consul`

	flagConsulAddress = `address`

	defaultCommonTag = `pinchy`
)

type (
	client interface {
		Agent() *api.Agent
	}
	factory func(*api.Config) (*api.Client, error)
)

func init() {
	set := pflag.NewFlagSet(registryName, pflag.ExitOnError)
	set.String(registry.MakeFlagName(flagConsulAddress), `127.0.0.1:8500`, `Consul http api address`)
	if err := registry.Register(registryName, set, NewRegistry); err != nil {
		panic(err)
	}
}

func provideClientConfig(v *viper.Viper, transport *http.Transport) (*api.Config, error) {
	flag := registry.MakeFlagName(flagConsulAddress)
	address := v.GetString(flag)
	if address == `` {
		return nil, errors.Errorf(`Flag "%s" is required`, flag)
	}

	cfg := api.DefaultConfig()
	cfg.Address = address
	cfg.Transport = transport
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

func provideAgent(c client) pkgConsul.Agent {
	return c.Agent()
}
