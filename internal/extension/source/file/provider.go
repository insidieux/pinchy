package file

import (
	pkgFile "github.com/insidieux/pinchy/pkg/core/source/file"

	"github.com/insidieux/pinchy/internal/extension/source"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	sourceName = `file`

	flagFilePath = `path`
)

func init() {
	set := pflag.NewFlagSet(sourceName, pflag.ExitOnError)
	set.String(source.MakeFlagName(flagFilePath), `$HOME/services.yml`, `YML file config path`)

	if err := source.Register(sourceName, set, NewSource, false); err != nil {
		panic(err)
	}
}

func provideReader() afero.Afero {
	return afero.Afero{
		Fs: afero.NewReadOnlyFs(afero.NewOsFs()),
	}
}

func providePath(v *viper.Viper) (pkgFile.Path, error) {
	flag := source.MakeFlagName(flagFilePath)
	path := v.GetString(flag)
	if path == `` {
		return ``, errors.Errorf(`flag "%s" is required`, flag)
	}
	return pkgFile.Path(path), nil
}
