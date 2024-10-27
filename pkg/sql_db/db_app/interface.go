package db_app

import (
	"context"
	"intmax2-withdrawal/pkg/sql_db/db_app/models"

	"github.com/dimiro1/health"
)

type SQLDb interface {
	GenericCommands
	ServiceCommands
	Withdrawals
}

type GenericCommands interface {
	Begin(ctx context.Context) (interface{}, error)
	Rollback()
	Commit() error
	Exec(ctx context.Context, input interface{}, executor func(d interface{}, input interface{}) error) (err error)
}

type ServiceCommands interface {
	Migrator(ctx context.Context, command string) (step int, err error)
	Check(ctx context.Context) health.Health
}

type Withdrawals interface {
	CreateWithdrawal(
		id string,
		transferData *models.TransferData,
		transferMerkleProof *models.TransferMerkleProof,
		transaction *models.Transaction,
		txMerkleProof *models.TxMerkleProof,
		transferHash string,
		blockNumber int64,
		blockHash string,
		enoughBalanceProof *models.EnoughBalanceProof,
	) (*models.Withdrawal, error)
	UpdateWithdrawalsStatus(ids []string, status models.WithdrawalStatus) error
	WithdrawalByID(id string) (*models.Withdrawal, error)
	WithdrawalsByHashes(transferHashes []string) (*[]models.Withdrawal, error)
	WithdrawalsByStatus(status models.WithdrawalStatus, limit *int) (*[]models.Withdrawal, error)
}
