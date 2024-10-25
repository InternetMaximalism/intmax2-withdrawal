package withdrawal

import (
	"context"
)

//go:generate mockgen -destination=mock_blockchain_service.go -package=withdrawal -source=blockchain_service.go

type ServiceBlockchain interface {
	ChainSB
}

type ChainSB interface {
	ScrollNetworkChainLinkEvmJSONRPC(ctx context.Context) (string, error)
}
