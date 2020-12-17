// +build wireinject

package consul

import (
	"github.com/google/wire"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/insidieux/pinchy/pkg/core"
	"github.com/insidieux/pinchy/pkg/core/registry/consul"
	"github.com/spf13/viper"
)

func NewRegistry(*viper.Viper) (core.Registry, func(), error) {
	panic(wire.Build(
		cleanhttp.DefaultPooledTransport,
		newClientConfig,
		provideConsulClientFactory,
		newClient,
		newAgent,
		consul.NewRegistry,
		wire.Bind(new(core.Registry), new(*consul.Registry)),
	))
}
