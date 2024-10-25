package post_withdrawals_by_hashes

import (
	"context"
	"fmt"
	"intmax2-withdrawal/configs"
	"intmax2-withdrawal/internal/logger"
	"intmax2-withdrawal/internal/open_telemetry"

	postWithdrwalsByHashes "intmax2-withdrawal/internal/use_cases/post_withdrawals_by_hashes"
	service "intmax2-withdrawal/internal/withdrawal_service"
	mDBApp "intmax2-withdrawal/pkg/sql_db/db_app/models"
)

// uc describes use case
type uc struct {
	cfg *configs.Config
	log logger.Logger
	db  SQLDriverApp
}

func New(cfg *configs.Config, log logger.Logger, db SQLDriverApp) postWithdrwalsByHashes.UseCasePostWithdrawalsByHashes {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(ctx context.Context, input *postWithdrwalsByHashes.UCPostWithdrawalsByHashesInput) (*[]mDBApp.Withdrawal, error) {
	const (
		hName = "UseCase PostWithdrawalsByHashes"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCInputEmpty)
		return nil, ErrUCInputEmpty
	}

	withdrawals, err := service.PostWithdrawalsByHashes(ctx, u.cfg, u.log, u.db, input)
	if err != nil {
		return nil, fmt.Errorf("failed to post withdrawals by hashes: %w", err)
	}

	return withdrawals, nil
}
