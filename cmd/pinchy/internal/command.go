package internal

import (
	"fmt"
	"time"

	"github.com/insidieux/pinchy/internal/extension/registry"
	"github.com/insidieux/pinchy/internal/extension/source"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	name = `pinchy`
)

// NewCommand provide root cobra.Command
// Root command contains pre-generated subcommands for all registered core.Source and core.Registry
func NewCommand(version string) *cobra.Command {
	rootCommand := &cobra.Command{
		Use:     name,
		Version: version,
	}
	rootCommand.SetOut(logrus.New().Out)

	for _, sourceProvider := range source.GetProviderList() {
		sourceCmd := &cobra.Command{
			Use:   sourceProvider.Name(),
			Short: fmt.Sprintf(`Fetch data from source "%s"`, sourceProvider.Name()),
		}
		for _, registryProvider := range registry.GetProviderList() {
			registryCmd := &cobra.Command{
				Use:   registryProvider.Name(),
				Short: fmt.Sprintf(`Save data in registry "%s"`, registryProvider.Name()),
			}
			onceCommand := &cobra.Command{
				Use:   `once`,
				Short: `Run main process only once: sync and return result`,
				RunE: func(cmd *cobra.Command, args []string) error {
					manager, cleanup, err := newManager(cmd.Flags(), sourceProvider.Factory(), registryProvider.Factory())
					if cleanup != nil {
						cleanup()
					}
					if err != nil {
						return errors.Wrap(err, `failed to bootstrap manager`)
					}
					return manager.Run(cmd.Context())
				},
			}
			watchCommand := &cobra.Command{
				Use:   `watch`,
				Short: `Run main process as daemon: sync repeatedly with constant interval`,
				RunE: func(cmd *cobra.Command, args []string) error {
					sc, cleanup, err := newScheduler(cmd.Flags(), sourceProvider.Factory(), registryProvider.Factory())
					if cleanup != nil {
						cleanup()
					}
					if err != nil {
						return errors.Wrap(err, `failed to bootstrap scheduler`)
					}
					sc.Run(cmd.Context())
					return nil
				},
			}
			watchCommand.Flags().Duration(`scheduler.interval`, time.Minute, `Interval between manager runs (1s, 1m, 5m, 1h and others)`)
			registryCmd.PersistentFlags().Bool(`manager.continue-on-error`, false, `Omit errors during process manager`)
			registryCmd.PersistentFlags().Bool(`manager.exit-on-error`, false, `Stop manager process on first error and by pass it to command line`)
			_ = registryCmd.PersistentFlags().MarkDeprecated(`manager.continue-on-error`, `Flag "manager.continue-on-error" is deprecated, use "manager.exit-on-error" instead`)
			_ = registryCmd.PersistentFlags().MarkHidden(`manager.continue-on-error`)
			registryCmd.PersistentFlags().AddFlagSet(registryProvider.Flags())
			registryCmd.AddCommand(onceCommand)
			registryCmd.AddCommand(watchCommand)
			sourceCmd.AddCommand(registryCmd)
		}
		sourceCmd.PersistentFlags().AddFlagSet(sourceProvider.Flags())
		rootCommand.AddCommand(sourceCmd)
	}
	rootCommand.PersistentFlags().String(`logger.level`, logrus.InfoLevel.String(), `Log level`)
	return rootCommand
}
