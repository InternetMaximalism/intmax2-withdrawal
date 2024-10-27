package blockchain

import (
	"context"
	"intmax2-withdrawal/configs"
)

type serviceBlockchain struct {
	ctx context.Context
	cfg *configs.Config
}

func New(
	ctx context.Context,
	cfg *configs.Config,
) ServiceBlockchain {
	return &serviceBlockchain{
		ctx: ctx,
		cfg: cfg,
	}
}
