package blockchain

import (
	"context"
	"errors"
	"fmt"
	errorsB "intmax2-withdrawal/internal/blockchain/errors"
	"intmax2-withdrawal/internal/open_telemetry"
	"strings"

	"github.com/prodadidb/go-validation"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var ErrScrollChainIDInvalid = fmt.Errorf(errorsB.ErrScrollChainIDInvalidStr, ScrollMainNetChainID, ScrollSepoliaChainID)

type ChainIDType string

const (
	ScrollMainNetChainID ChainIDType = "534352"
	ScrollSepoliaChainID ChainIDType = "534351"
)

type ChainLinkEvmJSONRPC string

const (
	ScrollMainNetChainLinkEvmJSONRPC ChainLinkEvmJSONRPC = "https://rpc.scroll.io"
	ScrollSepoliaChainLinkEvmJSONRPC ChainLinkEvmJSONRPC = "https://sepolia-rpc.scroll.io"
)

func (sb *serviceBlockchain) scrollNetworkChainIDValidator() error {
	return validation.Validate(sb.cfg.Blockchain.ScrollNetworkChainID,
		validation.Required,
		validation.In(
			string(ScrollMainNetChainID), string(ScrollSepoliaChainID),
		),
	)
}

func (sb *serviceBlockchain) ScrollNetworkChainLinkEvmJSONRPC(ctx context.Context) (string, error) {
	const (
		hName                   = "ServiceBlockchain func:ScrollNetworkChainLinkEvmJSONRPC"
		scrollNetworkChainIDKey = "scroll_network_chain_id"
		emptyKey                = ""
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(scrollNetworkChainIDKey, sb.cfg.Blockchain.ScrollNetworkChainID),
		))
	defer span.End()

	err := sb.scrollNetworkChainIDValidator()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return emptyKey, errors.Join(ErrScrollChainIDInvalid, err)
	}

	if strings.EqualFold(sb.cfg.Blockchain.ScrollNetworkChainID, string(ScrollMainNetChainID)) {
		return string(ScrollMainNetChainLinkEvmJSONRPC), nil
	}

	return string(ScrollSepoliaChainLinkEvmJSONRPC), nil
}
