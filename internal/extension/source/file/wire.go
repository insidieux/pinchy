// +build wireinject

package file

import (
	pkgFile "github.com/insidieux/pinchy/pkg/core/source/file"

	"github.com/google/wire"
	"github.com/insidieux/pinchy/pkg/core"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func NewSource(*viper.Viper) (core.Source, func(), error) {
	panic(wire.Build(
		newReader,
		wire.Bind(new(pkgFile.Reader), new(afero.Afero)),
		newPath,
		pkgFile.NewSource,
		wire.Bind(new(core.Source), new(*pkgFile.Source)),
	))
}
