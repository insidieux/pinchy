// +build wireinject

package consul

import (
	"github.com/google/wire"
	"github.com/insidieux/pinchy/pkg/core"
	"github.com/insidieux/pinchy/pkg/core/registry/consul/agent"
	"github.com/insidieux/pinchy/pkg/core/registry/consul/catalog"
	"github.com/spf13/viper"
)

var (
	wireSet = wire.NewSet(
		provideClientConfig,
		provideConsulClientFactory,
		provideClient,
		provideTag,
	)
)

func NewAgentRegistry(*viper.Viper) (core.Registry, func(), error) {
	panic(wire.Build(
		wireSet,
		provideAgent,
		agent.NewRegistry,
		wire.Bind(new(core.Registry), new(*agent.Registry)),
	))
}

func NewCatalogRegistry(*viper.Viper) (core.Registry, func(), error) {
	panic(wire.Build(
		wireSet,
		provideCatalog,
		catalog.NewRegistry,
		wire.Bind(new(core.Registry), new(*catalog.Registry)),
	))
}
