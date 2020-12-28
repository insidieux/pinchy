package internal

import (
	"strings"
	"time"

	"github.com/insidieux/pinchy/internal/extension/registry"
	"github.com/insidieux/pinchy/internal/extension/source"
	"github.com/insidieux/pinchy/pkg/core"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Provider for viper.Viper bound to current command pflag.FlagSet
func provideViper(set *pflag.FlagSet) (*viper.Viper, error) {
	v := viper.New()
	v.SetEnvPrefix(strings.ToUpper(name))
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	v.AutomaticEnv()
	if err := v.BindPFlags(set); err != nil {
		return nil, errors.Wrap(err, `failed to bind command line arguments`)
	}
	return v, nil
}

// Provider for logrus.Level
func provideLoggerLevel(commandViper *viper.Viper) (logrus.Level, error) {
	level := commandViper.GetString(`logger.level`)
	if level == `` {
		level = logrus.InfoLevel.String()
	}
	return logrus.ParseLevel(level)
}

// Provider for logrus.FieldLogger
func provideLogger(level logrus.Level) core.LoggerInterface {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetLevel(level)
	return logger
}

// Provider for core.Registry
func provideRegistry(commandViper *viper.Viper, factory registry.Factory, logger core.LoggerInterface) (core.Registry, func(), error) {
	r, cleanup, err := factory(commandViper)
	if r != nil {
		lr, ok := r.(core.Loggable)
		if ok {
			lr.WithLogger(logger)
		}
	}
	return r, cleanup, err
}

// Provider for core.Source
func provideSource(commandViper *viper.Viper, factory source.Factory, logger core.LoggerInterface) (core.Source, func(), error) {
	s, cleanup, err := factory(commandViper)
	if s != nil {
		ls, ok := s.(core.Loggable)
		if ok {
			ls.WithLogger(logger)
		}
	}
	return s, cleanup, err
}

// Provider for core.Source
func provideManagerExitOnError(commandViper *viper.Viper) core.ManagerExitOnError {
	return core.ManagerExitOnError(commandViper.GetBool(`manager.exit-on-error`))
}

// Provider for time.Ticker
func provideTicker(commandViper *viper.Viper) *time.Ticker {
	return time.NewTicker(commandViper.GetDuration(`scheduler.interval`))
}
