// nolint:unused
package withdrawal

import (
	"context"
	"intmax2-withdrawal/configs"
	"intmax2-withdrawal/internal/logger"
	withdrawalAggregator "intmax2-withdrawal/internal/use_cases/withdrawal_aggregator"
	ucWithdrawalAggregator "intmax2-withdrawal/pkg/use_cases/withdrawal_aggregator"
)

//go:generate mockgen -destination=mock_command.go -package=withdrawal -source=commands.go

type Commands interface {
	WithdrawalAggregator(
		ctx context.Context,
		cfg *configs.Config,
		log logger.Logger,
		db SQLDriverApp,
		sb ServiceBlockchain,
	) withdrawalAggregator.UseCaseWithdrawalAggregator
}

type commands struct{}

func newCommands() Commands {
	return &commands{}
}

func (c *commands) WithdrawalAggregator(
	ctx context.Context,
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
	sb ServiceBlockchain,
) withdrawalAggregator.UseCaseWithdrawalAggregator {
	return ucWithdrawalAggregator.New(ctx, cfg, log, db, sb)
}
