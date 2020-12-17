// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/insidieux/pinchy/internal/extension/registry"
	"github.com/insidieux/pinchy/internal/extension/source"
	"github.com/insidieux/pinchy/pkg/core"
	"github.com/spf13/pflag"
)

var (
	managerWireSet = wire.NewSet(
		provideViper,
		wire.NewSet(
			provideLoggerLevel,
			provideLogger,
		),
		provideRegistry,
		provideSource,
		provideManagerExitOnError,
		core.NewManager,
	)
	schedulerWireSet = wire.NewSet(
		provideTicker,
		managerWireSet,
		core.NewScheduler,
	)
)

func newManager(_ *pflag.FlagSet, _ source.Factory, _ registry.Factory) (core.ManagerInterface, func(), error) {
	panic(wire.Build(
		managerWireSet,
	))
}

func newScheduler(_ *pflag.FlagSet, _ source.Factory, _ registry.Factory) (*core.Scheduler, func(), error) {
	panic(wire.Build(
		schedulerWireSet,
	))
}
