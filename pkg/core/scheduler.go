package core

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type (
	// Scheduler is simple wrapper for ManagerInterface.
	// Scheduler provides possibility to call Manager.Run with constant interval passed to time.Ticker.
	Scheduler struct {
		ticker  *time.Ticker
		manager ManagerInterface
		logger  LoggerInterface
	}
)

// NewScheduler provides Scheduler with predefined time.Ticker, ManagerInterface and LoggerInterface
func NewScheduler(ticker *time.Ticker, manager ManagerInterface, logger LoggerInterface) *Scheduler {
	return &Scheduler{
		ticker:  ticker,
		manager: manager,
		logger:  logger,
	}
}

// Run start listen time.Ticker ticks until context.Context canceled.
func (s *Scheduler) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.ticker.C:
			s.logger.Infoln(`Running manager`)
			if err := s.manager.Run(ctx); err != nil {
				s.logger.Errorln(errors.Wrap(err, `failed to process manager run`).Error())
			}
		}
	}
}
