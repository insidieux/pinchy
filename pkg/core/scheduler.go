package core

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type (
	Scheduler struct {
		ticker  *time.Ticker
		manager ManagerInterface
		logger  LoggerInterface
	}
)

func NewScheduler(ticker *time.Ticker, manager ManagerInterface, logger LoggerInterface) *Scheduler {
	return &Scheduler{
		ticker:  ticker,
		manager: manager,
		logger:  logger,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.ticker.C:
			if err := s.manager.Run(ctx); err != nil {
				s.logger.Errorln(errors.Wrap(err, `failed to process manager run`).Error())
			}
		}
	}
}
