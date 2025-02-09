package withdrawal_aggregator

import (
	"context"
	"fmt"
	"intmax2-withdrawal/configs"
	"intmax2-withdrawal/internal/logger"
	"intmax2-withdrawal/internal/open_telemetry"

	ucWithdrawalAggregator "intmax2-withdrawal/internal/use_cases/withdrawal_aggregator"
	service "intmax2-withdrawal/internal/withdrawal_service"
)

// uc describes use case
type uc struct {
	cfg *configs.Config
	log logger.Logger
	db  SQLDriverApp
	sb  ServiceBlockchain
}

func New(
	ctx context.Context,
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
	sb ServiceBlockchain,
) ucWithdrawalAggregator.UseCaseWithdrawalAggregator {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
		sb:  sb,
	}
}

func (u *uc) Do(ctx context.Context) (err error) {
	const (
		hName = "UseCase WithdrawalAggregator"
	)

	u.log.Infof("Starting WithdrawalAggregator")

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	defer func() {
		if r := recover(); r != nil {
			const msg = "exec of withdrawal aggregator error occurred: %w"
			err = fmt.Errorf(msg, fmt.Errorf("%+v", r))
			open_telemetry.MarkSpanError(spanCtx, err)
		} else {
			u.log.Infof("Completed WithdrawalAggregator")
		}
	}()

	service.WithdrawalAggregator(spanCtx, u.cfg, u.log, u.db, u.sb)

	return err
}
