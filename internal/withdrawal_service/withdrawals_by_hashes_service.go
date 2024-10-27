package withdrawal_service

import (
	"context"
	"fmt"
	"intmax2-withdrawal/configs"
	"intmax2-withdrawal/internal/logger"
	postWithdrawalsByHashes "intmax2-withdrawal/internal/use_cases/post_withdrawals_by_hashes"
	mDBApp "intmax2-withdrawal/pkg/sql_db/db_app/models"
)

func PostWithdrawalsByHashes(
	ctx context.Context,
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
	input *postWithdrawalsByHashes.UCPostWithdrawalsByHashesInput,
) (*[]mDBApp.Withdrawal, error) {
	withdrawals, err := db.WithdrawalsByHashes(input.TransferHashes)
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawals by hashes: %w", err)
	}
	return withdrawals, nil
}
